package rpc

import (
	"github.com/fananchong/v-micro/codec"
)

type rpcRequest struct {
	service     string
	method      string
	contentType string
	codec       codec.Codec
	header      map[string]string
	body        []byte
	rawBody     interface{}
}

func (r *rpcRequest) Codec() codec.Reader {
	return r.codec
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

func (r *rpcRequest) Header() map[string]string {
	return r.header
}

func (r *rpcRequest) Body() interface{} {
	return r.rawBody
}
