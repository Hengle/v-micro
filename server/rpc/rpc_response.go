package rpc

import (
	"github.com/fananchong/v-micro/codec"
)

type rpcResponse struct {
	codec codec.Codec
}

func (r *rpcResponse) Codec() codec.Writer {
	return r.codec
}
