package main

import (
	"context"
	"time"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
)

// 暂未代码生成，测试用（未完成，测试用）

type testRequest struct {
}

// The service to call
func (req *testRequest) Service() string { return "greeter" }

// The action to take
func (req *testRequest) Method() string { return "Greeter.Hello" }

// The content type
func (req *testRequest) ContentType() string { return "application/protobuf" }

// The unencoded request body
func (req *testRequest) Body() interface{} { return "" }

var service micro.Service

func test() (err error) {
	cl := service.Client()
	err = cl.Call(context.Background(), &testRequest{})
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func main() {
	service = micro.NewService(
		micro.Name("client"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.AfterStart(test),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
