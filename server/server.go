// Package server is an interface for a micro server
package server

import "context"

// Server is a simple micro server abstraction
type Server interface {
	Options() Options
	Init(...Option) error
	Handle(interface{}) error
	Start() error
	Stop() error
	String() string
}

// Request is a synchronous request interface
type Request interface {
	ID() string
	// The action requested
	Method() string
	// Content type provided
	ContentType() string
	// Header of the request
	Header() map[string]string
	// Body is the initial decoded value
	Body() interface{}
}

// HandlerFunc represents a single method of a handler. It's used primarily
// for the wrappers. What's handed to the actual method is the concrete
// request and response types.
type HandlerFunc func(ctx context.Context, req Request, rsp interface{}) error

// HandlerWrapper wraps the HandlerFunc and returns the equivalent
type HandlerWrapper func(HandlerFunc) HandlerFunc
