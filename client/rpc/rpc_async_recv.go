package rpc

import (
	"context"
	"sync"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	hcodec "github.com/fananchong/v-micro/internal/codec"
	"github.com/fananchong/v-micro/transport"
)

type rpcAsyncRecv struct {
	rc      *rpcClient
	clients sync.Map
}

func (r *rpcAsyncRecv) SetRPCClient(rc *rpcClient) {
	r.rc = rc
}

func (r *rpcAsyncRecv) Join(name string, cli transport.Client) {
LOOP:
	if old, ok := r.clients.LoadOrStore(name, cli); ok {
		if old == cli {
			return
		}
		old.(transport.Client).Close()
		r.clients.Delete(name)
		goto LOOP
	}
	go r.doRecv(cli)
}

func (r *rpcAsyncRecv) doRecv(cli transport.Client) {
	for {
		var msg transport.Message
		if err := cli.Recv(&msg); err != nil {
			log.Error(err)
			// TODO
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
			msg.Header["Content-Type"] = r.rc.Options().ContentType
			ct = r.rc.Options().ContentType
		}
		var cf codec.NewCodec
		var err error
		if cf, err = r.rc.newCodec(ct); err != nil {
			log.Error(err)
			continue
		}

		rcodec := newRPCCodec(&msg, cli, cf)

		// internal request
		request := &rpcRequest{
			service:     hcodec.GetHeader("Micro-Service", msg.Header),
			method:      hcodec.GetHeader("Micro-Method", msg.Header),
			contentType: ct,
			codec:       rcodec,
			body:        msg.Body,
		}

		// serve the actual request using the request router
		if err := r.rc.router.ServeRequest(ctx, request); err != nil {
			log.Error(err)
			continue
		}
	}
}
