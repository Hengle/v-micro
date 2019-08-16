package transport

import (
	"context"
)

// Options common options
type Options struct {
	// Receive buffer size
	RecvBufSize int

	// Send buffer size
	SendBufSize int

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// DialOptions dial options
type DialOptions struct {
	OnClose func()

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// ListenOptions listen options
type ListenOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option option
type Option func(*Options)

// DialOption dial option
type DialOption func(*DialOptions)

// ListenOption listen option
type ListenOption func(*ListenOptions)

// RecvBufSize set the receive buffer size
func RecvBufSize(size int) Option {
	return func(o *Options) {
		o.RecvBufSize = size
	}
}

// SendBufSize set the send buffer size
func SendBufSize(size int) Option {
	return func(o *Options) {
		o.SendBufSize = size
	}
}

// OnClose on close callback
func OnClose(f func()) DialOption {
	return func(o *DialOptions) {
		o.OnClose = f
	}
}
