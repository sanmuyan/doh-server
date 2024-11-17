package model

import (
	"github.com/miekg/dns"
)

type Response interface {
	Encode(*dns.Msg) (any, error)
}

type Question struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
}

type Answer struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	Data string `json:"data"`
	TTL  uint32 `json:"TTL"`
}

// DJAResponse https://developers.google.com/speed/public-dns/docs/doh/json
type DJAResponse struct {
	Status   int      `json:"Status"`
	TC       bool     `json:"TC"`
	RD       bool     `json:"RD"`
	RA       bool     `json:"RA"`
	AD       bool     `json:"AD"`
	CD       bool     `json:"CD"`
	Question Question `json:"Question"`
	Answer   []Answer `json:"Answer,omitempty"`
}

func NewDJAResponse() *DJAResponse {
	return &DJAResponse{}
}

// Encode 根据不同的 DNS 类型返回相应数据
func (r *DJAResponse) Encode(msg *dns.Msg) (any, error) {
	r.Status = msg.Rcode
	r.TC = msg.Truncated
	r.RD = msg.RecursionDesired
	r.RA = msg.RecursionAvailable
	r.AD = msg.AuthenticatedData
	r.CD = msg.CheckingDisabled
	r.Question = Question{
		Name: msg.Question[0].Name,
		Type: msg.Question[0].Qtype,
	}
	for _, ans := range msg.Answer {
		r.answer(ans)
	}
	for _, ans := range msg.Ns {
		r.answer(ans)
	}
	for _, ans := range msg.Extra {
		r.answer(ans)
	}
	return r, nil
}

func (r *DJAResponse) answer(ans dns.RR) {
	switch t := ans.(type) {
	case *dns.A:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.A.String(),
			TTL:  t.Hdr.Ttl,
		})
	case *dns.CNAME:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.Target,
			TTL:  t.Hdr.Ttl,
		})
	case *dns.AAAA:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.AAAA.String(),
			TTL:  t.Hdr.Ttl,
		})
	case *dns.MX:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.Mx,
			TTL:  t.Hdr.Ttl,
		})
	case *dns.TXT:
		for _, v := range t.Txt {
			r.Answer = append(r.Answer, Answer{
				Name: t.Hdr.Name,
				Type: t.Hdr.Rrtype,
				Data: v,
				TTL:  t.Hdr.Ttl,
			})
		}
	case *dns.SOA:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.Ns,
			TTL:  t.Hdr.Ttl,
		})
	case *dns.SRV:
		r.Answer = append(r.Answer, Answer{
			Name: t.Hdr.Name,
			Type: t.Hdr.Rrtype,
			Data: t.Target,
			TTL:  t.Hdr.Ttl,
		})
	default:
		r.Answer = append(r.Answer, Answer{
			Name: ans.Header().Name,
			Type: ans.Header().Rrtype,
			Data: ans.String(),
			TTL:  ans.Header().Ttl,
		})
	}
}

type DoHResponse struct {
}

func NewDoHResponse() *DoHResponse {
	return &DoHResponse{}
}

// Encode 直接返回二进制数据
func (r *DoHResponse) Encode(msg *dns.Msg) (any, error) {
	return msg.Pack()
}

type DNSResponse struct {
}

func NewDNSResponse() *DNSResponse {
	return &DNSResponse{}
}

// Encode 直接返回消息
func (r *DNSResponse) Encode(msg *dns.Msg) (any, error) {
	return msg, nil
}
