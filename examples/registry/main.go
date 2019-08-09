package main

import (
	"time"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/selector"
)

func test() (err error) {
	s := selector.DefaultSelector
	if _, err = s.Select("greeter"); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.AfterStart(test),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
