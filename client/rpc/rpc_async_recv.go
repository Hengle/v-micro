package rpc

import (
	"context"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	hcodec "github.com/fananchong/v-micro/internal/codec"
	"github.com/fananchong/v-micro/transport"
)

func (r *rpcClient) AsyncRecv(nodeID string, cli transport.Client) {
	go r.asyncRecv(nodeID, cli)
}

func (r *rpcClient) asyncRecv(nodeID string, cli transport.Client) {
	for {
		var msg transport.Message
		if err := cli.Recv(&msg); err != nil {
			log.Error(err)
			return
		}

		// strip our headers
		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}

		// create new context
		ctx := metadata.NewContext(context.Background(), hdr)

		// we use this Content-Type header to identify the codec needed
		ct := msg.Header[metadata.CONTENTTYPE]
		var cf codec.NewCodec
		var err error
		if cf, err = r.newCodec(ct); err != nil {
			log.Error(err)
			continue
		}

		rcodec0 := newRPCCodec(msg.Body, cli, cf)
		rcodec1 := newRPCCodec([]byte(hcodec.GetHeader("Micro-RD", msg.Header)), cli, cf)

		// internal request
		request := &rpcRequest{
			method:      hcodec.GetHeader(metadata.METHOD, msg.Header),
			contentType: ct,
			codec:       []codec.Codec{rcodec0, rcodec1},
		}

		// serve the actual request using the request router
		if err := r.router.ServeRequest(ctx, request); err != nil {
			log.Error(err)
			continue
		}
	}
}
