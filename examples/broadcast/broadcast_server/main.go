package main

import (
	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/broadcast/proto"
)

var service micro.Service

func main() {
	service = micro.NewService(
		micro.Name("broadcast_server"),
	)

	service.Init()

	// Register Handlers
	proto.RegisterSayHandler(service.Server(), new(Say))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
