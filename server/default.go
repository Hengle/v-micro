package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fananchong/v-micro/log"
	"github.com/google/uuid"
)

var (
	// DefaultAddress default address
	DefaultAddress = ":0"
	// DefaultName default name
	DefaultName = "server"
	// DefaultVersion default version
	DefaultVersion = "latest"
	// DefaultID default id
	DefaultID = uuid.New().String()
	// DefaultRegisterCheck default register check
	DefaultRegisterCheck = func(context.Context) error { return nil }
)

// DefaultServer default server
var DefaultServer Server

// DefaultOptions returns config options for the default service
func DefaultOptions() Options {
	return DefaultServer.Options()
}

// Init initialises the default server with options passed in
func Init(opt ...Option) {
	if DefaultServer == nil {
		panic("")
	}
	DefaultServer.Init(opt...)
}

// NewHandler creates a new handler interface using the default server
// Handlers are required to be a public object with public
// endpoints. Call to a service endpoint such as Foo.Bar expects
// the type:
//
//	type Foo struct {}
//	func (f *Foo) Bar(ctx, req, rsp) error {
//		return nil
//	}
//
func NewHandler(h interface{}) Handler {
	return DefaultServer.NewHandler(h)
}

// Handle registers a handler interface with the default server to
// handle inbound requests
func Handle(h Handler) error {
	return DefaultServer.Handle(h)
}

// Run starts the default server and waits for a kill
// signal before exiting. Also registers/deregisters the server
func Run() error {
	if err := Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Infof("Received signal %s", <-ch)

	return Stop()
}

// Start starts the default server
func Start() error {
	config := DefaultServer.Options()
	log.Infof("Starting server %s id %s", config.Name, config.ID)
	return DefaultServer.Start()
}

// Stop stops the default server
func Stop() error {
	log.Infof("Stopping server")
	return DefaultServer.Stop()
}

// String returns name of Server implementation
func String() string {
	return DefaultServer.String()
}
