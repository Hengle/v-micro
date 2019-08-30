package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
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
	// Use the generated client stub
	cl := proto.NewGreeterService("wrapper1_server", new(Greeter), service.Client())

	// Make request
	_ = cl.Hello(context.Background(), &proto.Request{
		Name: "John",
	})
	return
}

func main() {
	service = micro.NewService(
		micro.Name("wrapper1_client"),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
