package connector

import (
	"context"

	"github.com/fananchong/v-micro/transport"
)

// Options options
type Options struct {
	Transport transport.Transport

	OnConnected func(nodeID string, cli transport.Client)
	OnClose     func(nodeID string)

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

// OnConnected OnConnected
func OnConnected(h func(nodeID string, cli transport.Client)) Option {
	return func(o *Options) {
		o.OnConnected = h
	}
}

// OnClose OnClose
func OnClose(h func(nodeID string)) Option {
	return func(o *Options) {
		o.OnClose = h
	}
}
