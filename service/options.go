package service

import (
	"context"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/server"
	"github.com/fananchong/v-micro/transport"
	"github.com/micro/cli"
)

// Options options
type Options struct {
	// For the Command Line itself
	Name        string
	Description string
	Version     string

	// Flags and Action
	Flags  []cli.Flag
	Action func(*cli.Context)

	// We need pointers to things so we can swap them out if needed.
	Registry  *registry.Registry
	Selector  *selector.Selector
	Transport *transport.Transport
	Client    *client.Client
	Server    *server.Server

	Clients    map[string]func(...client.Option) client.Client
	Registries map[string]func(...registry.Option) registry.Registry
	Selectors  map[string]func(...selector.Option) selector.Selector
	Servers    map[string]func(...server.Option) server.Server
	Transports map[string]func(...transport.Option) transport.Transport

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option option
type Option func(o *Options)

// Name Command line Name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Description Command line Description
func Description(d string) Option {
	return func(o *Options) {
		o.Description = d
	}
}

// Version Command line Version
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Selector selector
func Selector(s *selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// Registry regitstry
func Registry(r *registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Transport transport
func Transport(t *transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Client client
func Client(c *client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Server server
func Server(s *server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// NewClient New client func
func NewClient(name string, b func(...client.Option) client.Client) Option {
	return func(o *Options) {
		o.Clients[name] = b
	}
}

// NewRegistry New registry func
func NewRegistry(name string, r func(...registry.Option) registry.Registry) Option {
	return func(o *Options) {
		o.Registries[name] = r
	}
}

// NewSelector New selector func
func NewSelector(name string, s func(...selector.Option) selector.Selector) Option {
	return func(o *Options) {
		o.Selectors[name] = s
	}
}

// NewServer New server func
func NewServer(name string, s func(...server.Option) server.Server) Option {
	return func(o *Options) {
		o.Servers[name] = s
	}
}

// NewTransport New transport func
func NewTransport(name string, t func(...transport.Option) transport.Transport) Option {
	return func(o *Options) {
		o.Transports[name] = t
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
