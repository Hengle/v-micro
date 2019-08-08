package client

import "github.com/fananchong/v-micro/codec"

type rpcResponse struct {
	header map[string]string
	body   []byte
	codec  codec.Codec
}

func (r *rpcResponse) Codec() codec.Reader {
	return r.codec
}

func (r *rpcResponse) Header() map[string]string {
	return r.header
}
