package common

import (
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
)

// defaultContentType default content type
var defaultContentType = "application/protobuf"

// InitOptions new
func InitOptions(options *client.Options, opts ...client.Option) {
	for _, o := range opts {
		o(options)
	}

	if options.Registry == nil {
		log.Fatal("Init fail. Registry is nil.")
	}

	if options.Transport == nil {
		log.Fatal("Init fail. Transport is nil.")
	}

	if options.Selector == nil {
		log.Fatal("Init fail. Selector is nil.")
	}

	if len(options.ContentType) == 0 {
		options.ContentType = defaultContentType
	}
}
