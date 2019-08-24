package main

import (
	"context"
	"math/rand"
	"time"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/echo/proto"
)

var service micro.Service
var cl proto.EchoService
var r = rand.New(rand.NewSource(time.Now().Unix()))

// Echo Echo
type Echo struct{}

// Echo Echo.Echo
func (e *Echo) Echo(ctx context.Context, req *proto.Request, rsp *proto.Response) {
	if req.Num != rsp.Num {
		log.Fatalf("data error!")
	}
	cl.Echo(context.Background(), &proto.Request{
		Num: uint32(r.Int31n(1000)),
	})
}

func start() (err error) {
	// Use the generated client stub
	cl = proto.NewEchoService("echo_server", new(Echo), service.Client())

	// Make request
	err = cl.Echo(context.Background(), &proto.Request{
		Num: uint32(r.Int31n(1000)),
	})
	return
}

func main() {
	service = micro.NewService(
		micro.Name("echo_client"),
		micro.AfterStart(start),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
