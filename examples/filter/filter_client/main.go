package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/hello/proto"
	"github.com/fananchong/v-micro/selector"
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
	cl := proto.NewGreeterService("filter_server", new(Greeter), service.Client())

	callOptions := []client.CallOption{client.WithSelectOption(
		selector.WithFilter(selector.FilterLabel("SERVER_ID", "2")),
	)}

	// Make request
	_ = cl.Hello(context.Background(), &proto.Request{
		Name: "John",
	}, callOptions...)
	return
}

func main() {
	service = micro.NewService(
		micro.Name("filter_client"),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
