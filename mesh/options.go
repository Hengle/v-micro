package mesh

import (
	"context"

	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/transport"
)

// Options options
type Options struct {
	Registry  registry.Registry
	Transport transport.Transport

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(opt ...Option) Options {
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}

	if opts.Registry == nil {
		opts.Registry = registry.DefaultRegistry
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	}

	return opts
}

// Registry used for discovery
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Transport mechanism for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Option option
type Option func(*Options)
