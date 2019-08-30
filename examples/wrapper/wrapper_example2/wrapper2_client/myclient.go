package main

import (
	"context"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/examples/hello/proto"
	"github.com/fananchong/v-micro/selector"
)

// 演示客户端，每次选择指定的 comet 上发

var myClient proto.GreeterService

func myStart() (err error) {
	log.Info("call myStart ...")
	myClient = proto.NewGreeterService("wrapper2_server", new(Greeter), service.Client())
	return
}

type cliWrapper struct {
	client.Client
}

// 根据元数据，每次选择指定的 comet 上发
func (wrapper *cliWrapper) Call(ctx context.Context, req client.Request, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	sid := md["SERVER_ID"]
	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(selector.FilterLabel("SERVER_ID", sid)),
	))
	return wrapper.Client.Call(ctx, req, callOptions...)
}

func newWrapper(c client.Client) client.Client {
	return &cliWrapper{c}
}
