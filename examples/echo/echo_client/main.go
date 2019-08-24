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
var r = rand.New(rand.NewSource(time.Now().Unix()))
var num = uint32(r.Int31n(1000))
var cl proto.EchoService

// Echo Echo
type Echo struct{}

// Echo Echo.Echo
func (e *Echo) Echo(ctx context.Context, rsp *proto.Response) {
	if num != rsp.Num {
		log.Fatalf("data error!")
	}
	num = uint32(r.Int31n(1000))
	cl.Echo(context.Background(), &proto.Request{
		Num: num,
	})
}

func start() (err error) {
	// Use the generated client stub
	cl = proto.NewEchoService("echo_server", new(Echo), service.Client())

	// Make request
	err = cl.Echo(context.Background(), &proto.Request{
		Num: num,
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
