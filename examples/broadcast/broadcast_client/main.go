package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/broadcast/proto"
)

var service micro.Service
var cl proto.SayService

func start() (err error) {
	// Use the generated client stub
	cl = proto.NewSayService("broadcast_server", new(Say), service.Client())

	// Make request
	err = cl.Ping(context.Background(), &proto.Request{
		Name: "ABC",
	})

	cl.BroadcastHello(context.Background(), &proto.Request{Name: "1"})
	cl.BroadcastHello(context.Background(), &proto.Request{Name: "2"})
	cl.BroadcastHello(context.Background(), &proto.Request{Name: "3"})

	return
}

func main() {
	service = micro.NewService(
		micro.Name("broadcast_client"),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
