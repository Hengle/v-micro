package main

import (
	"context"
	"sync/atomic"
	"time"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/examples/echo/proto"
)

// Echo Echo
type Echo struct{}

// Echo Echo.Echo
func (e *Echo) Echo(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Num = req.Num
	atomic.AddInt64(&counter, 1)
	return nil
}

var counter int64

func statistics() {
	var t = time.Now().UnixNano() / 1e6
	for {
		select {
		case <-time.After(time.Second * 5):
			now := time.Now().UnixNano() / 1e6
			v := atomic.SwapInt64(&counter, 0)
			log.Info("count: ", float64(v)/float64((now-t)/1000), "/s")
			t = now
		}
	}
}

func main() {
	service := micro.NewService(
		micro.Name("echo_server"),
		micro.AfterStart(func() (err error) { go statistics(); return }),
	)

	service.Init()

	// Register Handlers
	proto.RegisterEchoHandler(service.Server(), new(Echo))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
