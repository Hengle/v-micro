package rpc

import (
	"github.com/fananchong/v-micro/codec"
)

type rpcRequest struct {
	service     string
	method      string
	contentType string
	codec       []codec.Codec
	body        interface{}
}

func newRequest(service, method string, request interface{}, contentType string) *rpcRequest {
	return &rpcRequest{
		service:     service,
		method:      method,
		body:        request,
		contentType: contentType,
	}
}

func (r *rpcRequest) ContentType() string {
	return r.contentType
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

func (r *rpcRequest) Codec(index int) codec.Reader {
	return r.codec[index]
}
