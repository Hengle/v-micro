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
	md, _ := metadata.FromContext(ctx)
	for k, v := range md {
		log.Infof("k: %s, v: %s", k, v)
	}
	log.Infof("Received Greeter.Hello Response:%s", rsp.GetMsg())
}

func start() (err error) {
	// Use the generated client stub
	cl := proto.NewGreeterService("metadata_server", new(Greeter), service.Client())

	// Set arbitrary headers in context
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Account": "john",
		"ID":      "1",
	})

	// Make request
	_ = cl.Hello(ctx, &proto.Request{
		Name: "John",
	})
	return
}

func main() {
	service = micro.NewService(
		micro.Name("metadata_client"),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
