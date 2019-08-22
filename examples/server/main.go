package main

import (
	"context"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	proto "github.com/fananchong/v-micro/examples/server/proto"
)

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (s *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("Received Greeter.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
	)

	service.Init()

	// Register Handlers
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
