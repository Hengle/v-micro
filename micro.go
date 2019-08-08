// Package micro is a pluggable framework for microservices
package micro

import (
	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/server"
)

// Service is an interface that wraps the lower level libraries
// within go-micro. Its a convenience method for building
// and initialising services.
type Service interface {
	Init(...Option)
	Options() Options
	Client() client.Client
	Server() server.Server
	Run() error
	String() string
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return newService(opts...)
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(s server.Server, h interface{}) error {
	return s.Handle(h)
}
