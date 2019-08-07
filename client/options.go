package client

import (
	"context"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/transport"
)

// Options options
type Options struct {
	// Plugged interfaces
	Codecs    map[string]codec.NewCodec
	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	// Middleware for client
	Wrappers []Wrapper

	// Default Call Options
	CallOptions CallOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// CallOptions call options
type CallOptions struct {
	SelectOptions []selector.SelectOption

	// Address of remote hosts	// TODO
	Address []string
	// Backoff func
	Backoff BackoffFunc
	// Check if retriable func
	Retry RetryFunc
	// Number of Call attempts
	Retries int

	// Middleware for low level call func
	CallWrappers []CallWrapper

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

func newOptions(options ...Option) Options {
	opts := Options{
		Codecs: make(map[string]codec.NewCodec),
		CallOptions: CallOptions{
			Backoff: DefaultBackoff,
			Retry:   DefaultRetry,
			Retries: DefaultRetries,
		},
	}

	for _, o := range options {
		o(&opts)
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

// Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c codec.NewCodec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// Registry to find nodes for a given service
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Transport to use for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Selector Select is used to select a node to route a request to
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// Wrap Adds a Wrapper to a list of options passed into the client
func Wrap(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

// WrapCall Adds a Wrapper to the list of CallFunc wrappers
func WrapCall(cw ...CallWrapper) Option {
	return func(o *Options) {
		o.CallOptions.CallWrappers = append(o.CallOptions.CallWrappers, cw...)
	}
}

// Backoff is used to set the backoff function used
// when retrying Calls
func Backoff(fn BackoffFunc) Option {
	return func(o *Options) {
		o.CallOptions.Backoff = fn
	}
}

// Retries Number of retries when making the request.
// Should this be a Call Option?
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// Retry sets the retry function to be used when re-trying.
func Retry(fn RetryFunc) Option {
	return func(o *Options) {
		o.CallOptions.Retry = fn
	}
}

// Call Options

// WithAddress sets the remote addresses to use rather than using service discovery
func WithAddress(a ...string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

// WithSelectOption select option
func WithSelectOption(so ...selector.SelectOption) CallOption {
	return func(o *CallOptions) {
		o.SelectOptions = append(o.SelectOptions, so...)
	}
}

// WithCallWrapper is a CallOption which adds to the existing CallFunc wrappers
func WithCallWrapper(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, cw...)
	}
}

// WithBackoff is a CallOption which overrides that which
// set in Options.CallOptions
func WithBackoff(fn BackoffFunc) CallOption {
	return func(o *CallOptions) {
		o.Backoff = fn
	}
}

// WithRetry is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetry(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.Retry = fn
	}
}

// WithRetries is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}
