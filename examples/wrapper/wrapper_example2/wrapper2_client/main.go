package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/examples/hello/proto"
)

var service micro.Service

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (c *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) {
	log.Infof("Received Greeter.Hello Response:%s", rsp.GetMsg())
}

func start() (err error) {
	log.Info("call start ...")
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"SERVER_ID": "1",
	})
	_ = myClient.Hello(ctx, &proto.Request{
		Name: "John1",
	})

	ctx = metadata.NewContext(context.Background(), map[string]string{
		"SERVER_ID": "2",
	})
	_ = myClient.Hello(ctx, &proto.Request{
		Name: "John2",
	})

	return
}

func main() {
	service = micro.NewService(
		micro.Name("wrapper2_client"),
		micro.WrapClient(newWrapper),
		micro.AfterStart(myStart),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
