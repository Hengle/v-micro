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

func test() (err error) {
	// Use the generated client stub
	cl := proto.NewGreeterService("hello_server", new(Greeter), service.Client())

	// Make request
	err = cl.Hello(context.Background(), &proto.Request{
		Name: "John",
	})
	return
}

func main() {
	service = micro.NewService(
		micro.Name("hello_client"),
		micro.AfterStart(test),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
