package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	proto "github.com/fananchong/v-micro/examples/server/proto"
)

var service micro.Service

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (c *Greeter) Hello(ctx context.Context, rsp *proto.Response) {
	log.Info("Received Greeter.Hello Response:%s", rsp.GetMsg())
}

func test() (err error) {
	// Use the generated client stub
	cl := proto.NewGreeterService("greeter", new(Greeter), service.Client())

	// Make request
	err = cl.Hello(context.Background(), &proto.Request{
		Name: "John",
	})
	return
}

func main() {
	service = micro.NewService(
		micro.Name("client"),
		micro.AfterStart(test),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
