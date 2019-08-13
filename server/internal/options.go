package common

import (
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/server"
	"github.com/fananchong/v-micro/transport"
)

// InitOptions server 插件都要初始化 options ，因此把它拎出来，放这里
func InitOptions(opts *server.Options, opt ...server.Option) {
	for _, o := range opt {
		o(opts)
	}

	if opts.Registry == nil {
		opts.Registry = registry.DefaultRegistry
	} else {
		registry.DefaultRegistry = opts.Registry
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	} else {
		transport.DefaultTransport = opts.Transport
	}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = server.DefaultRegisterCheck
	} else {
		server.DefaultRegisterCheck = opts.RegisterCheck
	}

	if len(opts.Address) == 0 {
		opts.Address = server.DefaultAddress
	} else {
		server.DefaultAddress = opts.Address
	}

	if len(opts.Name) == 0 {
		opts.Name = server.DefaultName
	} else {
		server.DefaultName = opts.Name
	}

	if len(opts.ID) == 0 {
		opts.ID = server.DefaultID
	} else {
		server.DefaultID = opts.ID
	}

	if len(opts.Version) == 0 {
		opts.Version = server.DefaultVersion
	} else {
		server.DefaultVersion = opts.Version
	}
}
