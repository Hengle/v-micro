// Package client is an interface for an RPC client
package client

import (
	"context"

	"github.com/fananchong/v-micro/registry"
)

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidirectional streaming of requests.
type Client interface {
	Init(...Option) error
	Options() Options
	Handle(interface{}) error
	Call(ctx context.Context, req Request, opts ...CallOption) error
	String() string
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	// The service to call
	Service() string
	// The action to take
	Method() string
	// The content type
	ContentType() string
	// The unencoded request body
	Body() interface{}
}

// CallFunc represents the individual call func
type CallFunc func(ctx context.Context, node *registry.Node, req Request, opts CallOptions) error

// CallWrapper is a low level wrapper for the CallFunc
type CallWrapper func(CallFunc) CallFunc

// Wrapper wraps a client and returns a client
type Wrapper func(Client) Client
