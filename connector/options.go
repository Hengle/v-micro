package connector

import (
	"context"
)

// Options options
type Options struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option used to initialise the connector
type Option func(*Options)
