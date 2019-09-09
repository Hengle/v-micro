package rpc

import (
	"context"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	hcodec "github.com/fananchong/v-micro/internal/codec"
	"github.com/fananchong/v-micro/transport"
)

type requestInfo struct {
	nodeID string
	req    client.Request
}

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

		rcodec := newRPCCodec(msg.Body, cli, cf)

		id := hcodec.GetHeader(metadata.MSGID, msg.Header)
		if request, ok := r.requests.Load(id); ok {
			r.requests.Delete(id)
			req := (request.(*requestInfo)).req.(*rpcRequest)
			req.codec = rcodec
			// serve the actual request using the request router
			if err := r.router.ServeRequest(ctx, req); err != nil {
				log.Error(err)
				continue
			}
		} else {
			log.Errorf("No find msg, msg id: %s", id)
		}
	}
}

func (r *rpcClient) OnTcpClose(nodeID string) {
	r.requests.Range(func(id, request interface{}) bool {
		if (request.(*requestInfo)).nodeID == nodeID {
			r.requests.Delete(id)
		}
		return true
	})
}
