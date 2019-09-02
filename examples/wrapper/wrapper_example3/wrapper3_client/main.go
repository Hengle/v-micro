package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/hello/proto"
	"github.com/fananchong/v-micro/registry"
)

var service micro.Service

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (c *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) {
	log.Infof("Received Greeter.Hello Response:%s", rsp.GetMsg())
}

func start() (err error) {
	// Use the generated client stub
	cl := proto.NewGreeterService("wrapper3_server", new(Greeter), service.Client())

	// Make request
	_ = cl.Hello(context.Background(), &proto.Request{
		Name: "John",
	})
	return
}

func logWrapper(fn client.CallFunc) client.CallFunc {
	return func(ctx context.Context, node *registry.Node, req client.Request, opts client.CallOptions) error {
		log.Infof("[Log Wrapper] Before serving request method: %v", req.Method())
		err := fn(ctx, node, req, opts)
		log.Infof("[Log Wrapper] After serving request")
		return err
	}
}

func main() {
	service = micro.NewService(
		micro.Name("wrapper3_client"),
		micro.WrapCall(logWrapper),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
