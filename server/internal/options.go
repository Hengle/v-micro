package common

import (
	"context"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/server"
	"github.com/google/uuid"
)

// InitOptions server 插件都要初始化 options ，因此把它拎出来，放这里
func InitOptions(options *server.Options, opts ...server.Option) {
	for _, o := range opts {
		o(options)
	}

	if options.Registry == nil {
		log.Fatal("Init fail. Registry is nil.")
	}

	if options.Transport == nil {
		log.Fatal("Init fail. Transport is nil.")
	}

	if options.RegisterCheck == nil {
		options.RegisterCheck = func(context.Context) error { return nil }
	}

	if len(options.Address) == 0 {
		options.Address = ":0"
	}

	if len(options.Name) == 0 {
		options.Name = "server"
	}

	if len(options.ID) == 0 {
		options.ID = uuid.New().String()
	}

	if len(options.Version) == 0 {
		options.Version = "latest"
	}
}
