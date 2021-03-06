package micro

import (
	"context"
	"time"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/connector"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/server"
	"github.com/fananchong/v-micro/transport"
	"github.com/micro/cli"
)

// Options options
type Options struct {
	Logger    log.Logger
	Client    client.Client
	Server    server.Server
	Registry  registry.Registry
	Transport transport.Transport

	// Flags and Action
	Flags  []cli.Flag
	Action func(*cli.Context)

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option option
type Option func(*Options)

// Logger logger
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Client client
func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Server server
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		// Update Client and Server
		_ = o.Client.Init(client.Registry(r))
		_ = o.Server.Init(server.Registry(r))
		// Update Selector
		_ = o.Client.Options().Selector.Init(selector.Registry(r))
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		_ = o.Client.Init(client.Selector(s))
	}
}

// Connector sets the connector for the service client
func Connector(ct connector.Connector) Option {
	return func(o *Options) {
		_ = o.Client.Init(client.Connector(ct))
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		// Update Client and Server
		_ = o.Client.Init(client.Transport(t))
		_ = o.Server.Init(server.Transport(t))
		_ = o.Client.Options().Connector.Init(connector.Transport(t))
	}
}

// Convenience options

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.Address(addr))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.Metadata(md))
	}
}

// Flags flags
func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}

// Action action
func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Action = a
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		_ = o.Server.Init(server.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		_ = o.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		_ = o.Server.Init(wrappers...)
	}
}

// Before and Afters

// BeforeStart BeforeStart
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop BeforeStop
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart AfterStart
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop AfterStop
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
