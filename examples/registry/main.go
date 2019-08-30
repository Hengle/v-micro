package main

import (
	"time"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
)

var service micro.Service

func test() (err error) {
	// 这个是显式调用感兴趣的服务，来触发 registry 功能
	// 实际上 client.Call client.Broadcast 内部实现就会调用以下的代码；并按需连接微服务
	// 因此 test 内的代码通常使用者不会碰到
	slt := service.Client().Options().Selector
	if _, err = slt.Select("reg"); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	service = micro.NewService(
		micro.Name("reg"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.AfterStart(test),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
