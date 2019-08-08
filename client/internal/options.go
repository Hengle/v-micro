package common

import (
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/transport"
)

// NewOptions new
func NewOptions(options ...client.Option) client.Options {
	opts := client.Options{
		Codecs: make(map[string]codec.NewCodec),
		CallOptions: client.CallOptions{
			Backoff: client.DefaultBackoff,
			Retry:   client.DefaultRetry,
			Retries: client.DefaultRetries,
		},
	}

	for _, o := range options {
		o(&opts)
	}

	if len(opts.ContentType) == 0 {
		opts.ContentType = client.DefaultContentType
	}

	if opts.Registry == nil {
		opts.Registry = registry.DefaultRegistry
	}

	if opts.Selector == nil {
		opts.Selector = selector.DefaultSelector
	}

	if opts.Transport == nil {
		opts.Transport = transport.DefaultTransport
	}

	return opts
}
