package main

import (
	"context"
	"fmt"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/hello/proto"
	"github.com/fananchong/v-micro/server"
)

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (s *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Infof("Received Greeter.Hello request. Name:%s", req.Name)
	rsp.Msg = fmt.Sprintf("Hello %s", req.Name)
	return nil
}

func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Infof("[Log Wrapper] Before serving request method: %v", req.Method())
		err := fn(ctx, req, rsp)
		log.Infof("[Log Wrapper] After serving request")
		return err
	}
}

var service micro.Service

func main() {
	service = micro.NewService(
		micro.Name("wrapper1_server"),
		micro.WrapHandler(logWrapper),
	)

	service.Init()

	// Register Handlers
	_ = proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
