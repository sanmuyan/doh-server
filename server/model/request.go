package model

type Request interface {
	Response() Response
	Inbound() string
}

type DJARequest struct {
	res  Response
	Name string `form:"name" binding:"required"`
	Type uint16 `form:"type"`
}

func NewDJARequest(res Response) *DJARequest {
	return &DJARequest{
		Type: 1,
		res:  res,
	}
}

func (r *DJARequest) Response() Response {
	return r.res
}

func (r *DJARequest) Inbound() string {
	return "DJA"
}

type DoHRequest struct {
	res Response
}

func NewDoHRequest(res Response) *DoHRequest {
	return &DoHRequest{
		res: res,
	}
}

func (r *DoHRequest) Response() Response {
	return r.res
}

func (r *DoHRequest) Inbound() string {
	return "DoH"
}

type DNSRequest struct {
	res Response
}

func NewDNSRequest(res Response) *DNSRequest {
	return &DNSRequest{
		res: res,
	}
}

func (r *DNSRequest) Response() Response {
	return r.res
}

func (r *DNSRequest) Inbound() string {
	return "DNS"
}
