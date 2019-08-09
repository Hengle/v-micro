package main

import (
	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
)

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
