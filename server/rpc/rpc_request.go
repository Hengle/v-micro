package rpc

import (
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/transport"
)

type requestImpl struct {
	service     string
	method      string
	contentType string
	socket      transport.Socket
	codec       codec.Codec
	header      map[string]string
	body        []byte
	rawBody     interface{}
}

func (r *requestImpl) Codec() codec.Reader {
	return r.codec
}

func (r *requestImpl) ContentType() string {
	return r.contentType
}

func (r *requestImpl) Service() string {
	return r.service
}

func (r *requestImpl) Method() string {
	return r.method
}

func (r *requestImpl) Header() map[string]string {
	return r.header
}

func (r *requestImpl) Body() interface{} {
	return r.rawBody
}

func (r *requestImpl) Read() ([]byte, error) {
	// got a body
	if r.body != nil {
		b := r.body
		r.body = nil
		return b, nil
	}

	var msg transport.Message
	err := r.socket.Recv(&msg)
	if err != nil {
		return nil, err
	}
	r.header = msg.Header

	return msg.Body, nil
}
