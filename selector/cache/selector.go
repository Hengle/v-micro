package cache

import (
	"time"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/selector"
)

type registrySelector struct {
	so selector.Options
	rc Cache
}

func (c *registrySelector) newCache() Cache {
	ropts := []Option{}
	if c.so.Context != nil {
		if t, ok := c.so.Context.Value(ttlKey{}).(time.Duration); ok {
			ropts = append(ropts, WithTTL(t))
		}
	}
	return New(c.so.Registry, ropts...)
}

func (c *registrySelector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&c.so)
	}

	c.rc.Stop()
	c.rc = c.newCache()

	return nil
}

func (c *registrySelector) Options() selector.Options {
	return c.so
}

func (c *registrySelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	sopts := selector.SelectOptions{
		Strategy: c.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := c.rc.GetService(service)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services), nil
}

// Close stops the watcher and destroys the cache
func (c *registrySelector) Close() error {
	c.rc.Stop()

	return nil
}

func (c *registrySelector) String() string {
	return "registry"
}

// NewSelector new
func NewSelector(opts ...selector.Option) selector.Selector {
	sopts := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		log.Fatal("Init fail. Registry is nil.")
	}

	s := &registrySelector{
		so: sopts,
	}
	s.rc = s.newCache()

	return s
}
