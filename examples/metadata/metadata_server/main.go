package main

import (
	"context"
	"fmt"

	micro "github.com/fananchong/v-micro"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/examples/hello/proto"
)

// Greeter Greeter
type Greeter struct{}

// Hello Greeter.Hello
func (s *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	md, _ := metadata.FromContext(ctx)
	for k, v := range md {
		log.Infof("k: %s, v: %s", k, v)
	}
	rsp.Msg = fmt.Sprintf("Hello %s, Account:%s, ID:%s", req.Name, md["Account"], md["ID"])
	return nil
}

var service micro.Service

func beforeStart() error {
	// 服务器元数据，会通过 Registry 组件，发布。可以在 Selector 组件中给 Filter 使用
	// 更详细例子，看后续的 Selector 或 Filter 例子
	// 这里仅演示下，如何设置服务器元数据
	md := service.Server().Options().Metadata
	md["SERVER_ID"] = service.Server().Options().ID
	return nil
}

func main() {
	service = micro.NewService(
		micro.Name("metadata_server"),
		micro.BeforeStart(beforeStart),
	)

	service.Init()

	// Register Handlers
	_ = proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
