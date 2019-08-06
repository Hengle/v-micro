package rpc

import (
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/transport"
)

type responseImpl struct {
	header map[string]string
	socket transport.Socket
	codec  codec.Codec
}

func (r *responseImpl) Codec() codec.Writer {
	return r.codec
}

func (r *responseImpl) WriteHeader(hdr map[string]string) {
	for k, v := range hdr {
		r.header[k] = v
	}
}

func (r *responseImpl) Write(b []byte) error {
	return r.socket.Send(&transport.Message{
		Header: r.header,
		Body:   b,
	})
}
