package cache

import (
	"context"
	"time"

	"github.com/fananchong/v-micro/selector"
)

type ttlKey struct{}

// TTL Set the registry cache ttl
func TTL(t time.Duration) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ttlKey{}, t)
	}
}
