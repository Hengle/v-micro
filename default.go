package micro

import (
	"github.com/fananchong/v-micro/transport"
	"github.com/fananchong/v-micro/transport/gotcp"
)

var (
	// DefaultTransport default transport obj
	DefaultTransport transport.Transport = gotcp.NewTransport()
)
