// Package metadata is a way of defining message headers
package metadata

import (
	"context"
)

type metaKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

// Copy copy metadata
func Copy(md Metadata) Metadata {
	cmd := make(Metadata)
	for k, v := range md {
		cmd[k] = v
	}
	return cmd
}

// FromContext from context
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metaKey{}).(Metadata)
	return md, ok
}

// NewContext new
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metaKey{}, md)
}
