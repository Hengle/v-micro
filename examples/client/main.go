package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
)

// 暂未代码生成，测试用

type testRequest struct {
}

// The service to call
func (req *testRequest) Service() string { return "Greeter" }

// The action to take
func (req *testRequest) Method() string { return "Hello" }

// The content type
func (req *testRequest) ContentType() string { return "application/protobuf" }

// The unencoded request body
func (req *testRequest) Body() interface{} { return "" }

func main() {
	// create a new service
	service := micro.NewService(
		micro.Name("client"),
	)

	// parse command line flags
	service.Init()

	cl := service.Client()

	err := cl.Call(context.Background(), &testRequest{})
	if err != nil {
		log.Error(err)
		return
	}
}
