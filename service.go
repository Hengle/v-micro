package micro

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/server"
	cmd "github.com/fananchong/v-micro/service"
)

type service struct {
	opts Options
	cmd  cmd.Cmd
}

func newService(opt ...Option) Service {
	// 1. Initialize the default plugin
	c := cmd.NewCmd()
	opts := Options{
		Logger:    *c.Options().Logger,
		Client:    *c.Options().Client,
		Server:    *c.Options().Server,
		Registry:  *c.Options().Registry,
		Transport: *c.Options().Transport,
		Action:    c.Options().Action,
		Context:   context.Background(),
	}

	// 2. Replace the default plugin with options
	for _, o := range opt {
		o(&opts)
	}

	// 3. After instantiation, service.Init is called,
	//    and the command line parameters are internally parsed
	//    to dynamically replace the plugin.
	return &service{
		opts: opts,
		cmd:  c,
	}
}

// Init initialises options. Additionally it calls service.Init
// which parses command line flags. service.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	// Initialise the command flags, overriding new service
	_ = s.cmd.Init(
		// The plugin can be replaced by an option in the code,
		// so reassign it to ensure consistency
		cmd.ID(s.opts.Server.Options().ID),
		cmd.Name(s.opts.Server.Options().Name),
		cmd.Version(s.opts.Server.Options().Version),
		cmd.Flags(s.opts.Flags...),
		cmd.Action(s.opts.Action),
		cmd.Logger(&s.opts.Logger),
		cmd.Registry(&s.opts.Registry),
		cmd.Transport(&s.opts.Transport),
		cmd.Client(&s.opts.Client),
		cmd.Server(&s.opts.Server),
	)

	s.opts.Client = &cmd.ClientWrapper{
		Client: s.opts.Client,
		Headers: metadata.Metadata{
			"Micro-From-Service": s.opts.Server.Options().Name,
		},
	}
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) String() string {
	return "micro"
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}
