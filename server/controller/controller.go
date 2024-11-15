package controller

import (
	"doh-server/server/model"
	"doh-server/server/service"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// API列表

var svc = service.NewService()

// DoHQuery RFC 8484 https://datatracker.ietf.org/doc/html/rfc8484
// 提供 DoH 查询
func DoHQuery(c *gin.Context) {
	var body []byte
	var err error
	switch c.Request.Method {
	case http.MethodGet:
		dnsQuery := c.Query("dns")
		dnsQueryLen := len(dnsQuery)
		body = make([]byte, dnsQueryLen, dnsQueryLen+8)
		_, err = base64.RawURLEncoding.Decode(body, []byte(dnsQuery))
	case http.MethodPost:
		body, err = io.ReadAll(c.Request.Body)
	}
	if err != nil {
		logrus.Errorf("read body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg:": "bad request"})
		return
	}
	msg := new(dns.Msg)
	err = msg.Unpack(body)
	if err != nil {
		logrus.Errorf("unpack body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg:": "bad request"})
		return
	}
	res, err := svc.Query(model.NewDoHRequest(model.NewDoHResponse()), msg)
	if err != nil {
		logrus.Errorf("query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg:": "server error"})
		return
	}
	c.Data(http.StatusOK, "application/dns-message", res.([]byte))
}

// DJAQuery 提供 JSON API DNS 查询
func DJAQuery(c *gin.Context) {
	req := model.NewDJARequest(model.NewDJAResponse())
	err := c.ShouldBindQuery(req)
	if err != nil {
		logrus.Errorf("bind query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg:": "bad request"})
		return
	}
	msg := new(dns.Msg)
	msg.SetQuestion(req.Name+".", req.Type)
	res, err := svc.Query(req, msg)
	if err != nil {
		logrus.Errorf("query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg:": "server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DNSQuery 提供 DNS 查询
func DNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	res, err := svc.Query(model.NewDNSRequest(model.NewDNSResponse()), r)
	if err != nil {
		logrus.Errorf("query: %v", err)
		return
	}
	m := res.(*dns.Msg)
	err = w.WriteMsg(m)
	if err != nil {
		logrus.Errorf("write msg: %v", err)
	}
}
