package common

import (
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/server"
	"github.com/fananchong/v-micro/transport"
)

// NewOptions server 插件都要初始化 options ，因此把它拎出来，放这里
func NewOptions(opt ...server.Option) server.Options {
	opts := server.Options{
		Codecs:   make(map[string]codec.NewCodec),
		Metadata: map[string]string{},
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Registry == nil {
		opts.Registry = registry.DefaultRegistry
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = server.DefaultRegisterCheck
	}

	if len(opts.Address) == 0 {
		opts.Address = server.DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = server.DefaultName
	}

	if len(opts.ID) == 0 {
		opts.ID = server.DefaultID
	}

	if len(opts.Version) == 0 {
		opts.Version = server.DefaultVersion
	}

	return opts
}
