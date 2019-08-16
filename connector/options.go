package connector

import (
	"context"

	"github.com/fananchong/v-micro/transport"
)

// Options options
type Options struct {
	Transport transport.Transport

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option used to initialise the connector
type Option func(*Options)

// Transport mechanism for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}
