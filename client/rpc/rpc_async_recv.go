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

		// set local/remote ips
		hdr["Local"] = cli.Local()
		hdr["Remote"] = cli.Remote()

		// create new context
		ctx := metadata.NewContext(context.Background(), hdr)

		// we use this Content-Type header to identify the codec needed
		ct := msg.Header["Content-Type"]
		// no content type
		if len(ct) == 0 {
			msg.Header["Content-Type"] = r.Options().ContentType
			ct = r.Options().ContentType
		}
		var cf codec.NewCodec
		var err error
		if cf, err = r.newCodec(ct); err != nil {
			log.Error(err)
			continue
		}

		rcodec := newRPCCodec(cli, cf)

		// internal request
		request := &rpcRequest{
			service:     hcodec.GetHeader("Micro-Service", msg.Header),
			method:      hcodec.GetHeader("Micro-Method", msg.Header),
			contentType: ct,
			codec:       rcodec,
			body:        msg.Body,
		}

		// serve the actual request using the request router
		if err := r.router.ServeRequest(ctx, request); err != nil {
			log.Error(err)
			continue
		}
	}
}
