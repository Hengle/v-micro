package micro

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/server"
	cmd "github.com/fananchong/v-micro/service"
)

type service struct {
	opts Options
	cmd  cmd.Cmd
	once sync.Once
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)

	options.Client = &cmd.ClientWrapper{
		options.Client,
		metadata.Metadata{
			"Micro-From-Service": options.Server.Options().Name,
		},
	}

	return &service{
		opts: options,
		cmd:  cmd.DefaultCmd,
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

	s.once.Do(func() {
		// Initialise the command flags, overriding new service
		_ = s.cmd.Init(
			cmd.Flags(s.opts.Flags...),
			cmd.Action(s.opts.Action),
			cmd.Registry(&s.opts.Registry),
			cmd.Transport(&s.opts.Transport),
			cmd.Client(&s.opts.Client),
			cmd.Server(&s.opts.Server),
		)
	})
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
