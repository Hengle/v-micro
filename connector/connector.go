package connector

import (
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/transport"
)

// Connector Get available connections based on registration information
type Connector interface {
	Init(opts ...Option) error
	Options() Options
	Get(node *registry.Node) (transport.Client, error)
	Close() error
	String() string
}
