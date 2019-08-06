// Package server is an interface for a micro server
package server

import (
	"github.com/fananchong/v-micro/codec"
)

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
	// Service name requested
	Service() string
	// The action requested
	Method() string
	// Content type provided
	ContentType() string
	// Header of the request
	Header() map[string]string
	// Body is the initial decoded value
	Body() interface{}
	// Read the undecoded request body
	Read() ([]byte, error)
	// The encoded message stream
	Codec() codec.Reader
}

// Response is the response writer for unencoded messages
type Response interface {
	// Encoded writer
	Codec() codec.Writer
	// Write the header
	WriteHeader(map[string]string)
	// write a response directly to the client
	Write([]byte) error
}
