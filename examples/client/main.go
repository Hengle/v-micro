package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	proto "github.com/fananchong/v-micro/examples/server/proto"
)

var service micro.Service

func test() (err error) {
	// Use the generated client stub
	cl := proto.NewGreeterService("greeter", service.Client())

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
