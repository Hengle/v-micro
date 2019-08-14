package common

import (
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
)

// defaultContentType default content type
var defaultContentType = "application/protobuf"

// InitOptions new
func InitOptions(opts *client.Options, options ...client.Option) {
	for _, o := range options {
		o(opts)
	}

	if opts.Registry == nil {
		log.Fatal("Init fail. Registry is nil.")
	}

	if opts.Transport == nil {
		log.Fatal("Init fail. Transport is nil.")
	}

	if opts.Selector == nil {
		log.Fatal("Init fail. Selector is nil.")
	}

	if len(opts.ContentType) == 0 {
		opts.ContentType = defaultContentType
	}
}
