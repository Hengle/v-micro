package singleconnection

import (
	"github.com/fananchong/v-micro/connector"
)

type singleConnection struct {
	opts connector.Options
}

func (s *singleConnection) Init(opts ...connector.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	return nil
}

func (s *singleConnection) Options() connector.Options {
	return s.opts
}

func (s *singleConnection) Close() error {
	return nil
}

func (s *singleConnection) String() string {
	return "singleconnection"
}

// NewConnector new
func NewConnector(opts ...connector.Option) connector.Connector {
	return nil
}
