package service

import (
	"context"
	"crypto/tls"
	"doh-server/pkg/config"
	"doh-server/server/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/sanmuyan/xpkg/xrequest"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"time"
)

// 接口逻辑

type Service struct {
	ctx context.Context
}

func NewService() *Service {
	return &Service{
		ctx: context.Background(),
	}
}

var Client *dns.Client

var CacheStore *cache.Cache

func (s *Service) LoadCache(questionID string) (*dns.Msg, bool) {
	if !config.Conf.Cache {
		return nil, false
	}
	if data, ok := CacheStore.Get(questionID); ok {
		return data.(*dns.Msg), true
	}
	return nil, false
}

func (s *Service) StoreCache(questionID string, msg *dns.Msg) {
	if !config.Conf.Cache {
		return
	}
	CacheStore.Set(questionID, msg, cache.DefaultExpiration)
}

// Query 查询上游 DNS
func (s *Service) Query(req model.Request, reqMsg *dns.Msg) (any, error) {
	if len(reqMsg.Question) == 0 {
		return nil, errors.New("empty question")
	}
	// 先查询缓存，如果缓存命中，直接返回
	questionID := s.buildQuestionID(reqMsg)
	if cacheMsg, ok := s.LoadCache(questionID); ok {
		logrus.Infof("cache request: %s %s", req.Inbound(), questionID)
		logrus.Infof("cache response: %s %s %s", req.Inbound(), questionID, s.buildAnswerID(cacheMsg))
		// 缓存需要设置 Reply 否则客户端可能不接受
		return req.Response().Encode(cacheMsg.Copy().SetReply(reqMsg))
	}
	var err error
	var resMsg *dns.Msg
	logrus.Infof("server request: %s %s:%s %s", req.Inbound(), config.Conf.UpstreamNet, config.Conf.UpstreamServer, questionID)
	switch config.Conf.UpstreamNet {
	case "doh":
		resMsg, err = s.queryDoH(reqMsg)
	case "udp", "tcp", "tcp-tls":
		resMsg, _, err = Client.Exchange(reqMsg, config.Conf.UpstreamServer)
	default:
		return nil, errors.New("unsupported upstream net")
	}
	if err != nil {
		return nil, err
	}
	if resMsg.Rcode != dns.RcodeSuccess {
		logrus.Warnf("rcode not success: %s", dns.RcodeToString[resMsg.Rcode])
	}
	go s.StoreCache(questionID, resMsg)
	logrus.Infof("server response: %s %s %s", req.Inbound(), questionID, s.buildAnswerID(resMsg))
	return req.Response().Encode(resMsg)
}

func (s *Service) buildQuestionID(msg *dns.Msg) string {
	question := msg.Question[0]
	return fmt.Sprintf("%s %s %s", dns.Type(question.Qtype).String(), dns.Class(question.Qclass).String(), question.Name)
}

func (s *Service) buildAnswerID(msg *dns.Msg) string {
	return string(xutil.RemoveError(json.Marshal(msg.Answer)))
}

func (s *Service) queryDoH(reqMsg *dns.Msg) (*dns.Msg, error) {
	buf, err := reqMsg.Pack()
	if err != nil {
		return nil, err
	}
	req := xrequest.Request{
		Config: &xrequest.Options{
			Method:  "POST",
			URL:     config.Conf.UpstreamServer,
			Timeout: config.Conf.UpstreamTimeout,
			Body:    buf,
			Head: map[string]string{
				"Content-Type": "application/dns-message",
			},
		},
	}
	resp, err := req.Request()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	resMsg := new(dns.Msg)
	err = resMsg.Unpack(resp.Body)
	if err != nil {
		return nil, err
	}
	return resMsg, nil
}

// InitClient
// 初始化上游 DNS 客户端
// 初始化缓存客户端
func InitClient() {
	Client = new(dns.Client)
	Client.Timeout = time.Duration(config.Conf.UpstreamTimeout) * time.Second
	Client.Net = config.Conf.UpstreamNet
	if config.Conf.UpstreamNet == "tcp-tls" {
		Client.TLSConfig = &tls.Config{}
	}
	CacheStore = cache.New(time.Duration(config.Conf.CacheTTL)*time.Second, time.Duration(config.Conf.CacheTTL*10)*time.Second)
}
