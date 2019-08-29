package main

import (
	"context"

	"github.com/fananchong/protoc-gen-vmicro/micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/broadcast/proto"
)

// Say Say
type Say struct{}

// BroadcastHello Say.BroadcastHello
func (s *Say) BroadcastHello(ctx context.Context, req *proto.Request, rsp *micro.NoReply) error {
	log.Infof("BroadcastHello, recv req: %v", *req)
	return nil
}

// Ping Say.Ping
func (s *Say) Ping(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Infof("Ping, recv req: %v", *req)
	rsp.Msg = "Hello " + req.Name + ". I am " + service.Server().Options().ID
	return nil
}
