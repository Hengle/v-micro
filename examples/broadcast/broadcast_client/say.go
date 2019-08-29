package main

import (
	"context"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/broadcast/proto"
)

// Say Say
type Say struct{}

// Ping Say.Ping
func (s *Say) Ping(ctx context.Context, req *proto.Request, rsp *proto.Response) {
	log.Infof("Ping, recv rsp:%v", *rsp)
}
