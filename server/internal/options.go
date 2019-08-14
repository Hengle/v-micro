package common

import (
	"context"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/server"
	"github.com/google/uuid"
)

// InitOptions server 插件都要初始化 options ，因此把它拎出来，放这里
func InitOptions(opts *server.Options, opt ...server.Option) {
	for _, o := range opt {
		o(opts)
	}

	if opts.Registry == nil {
		log.Fatal("Init fail. Registry is nil.")
	}

	if opts.Transport == nil {
		log.Fatal("Init fail. Transport is nil.")
	}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = func(context.Context) error { return nil }
	}

	if len(opts.Address) == 0 {
		opts.Address = ":0"
	}

	if len(opts.Name) == 0 {
		opts.Name = "server"
	}

	if len(opts.ID) == 0 {
		opts.ID = uuid.New().String()
	}

	if len(opts.Version) == 0 {
		opts.Version = "latest"
	}
}
