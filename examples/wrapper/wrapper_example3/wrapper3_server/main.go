package main

import (
	"context"
	"fmt"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/hello/proto"
)

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (s *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Infof("Received Greeter.Hello request. Name:%s", req.Name)
	rsp.Msg = fmt.Sprintf("Hello %s", req.Name)
	return nil
}

var service micro.Service

func main() {
	service = micro.NewService(
		micro.Name("wrapper3_server"),
	)

	service.Init()

	// Register Handlers
	_ = proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
