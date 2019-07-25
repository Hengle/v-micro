// Package memory provides an in-memory registry
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/fananchong/v-micro/registry"
	"github.com/google/uuid"
)

type registryImpl struct {
	options registry.Options

	sync.RWMutex
	Services map[string][]*registry.Service
	Watchers map[string]*watcherImpl
}

var (
	timeout = time.Millisecond * 10
)

func (m *registryImpl) watch(r *registry.Result) {
	var watchers []*watcherImpl

	m.RLock()
	for _, w := range m.Watchers {
		watchers = append(watchers, w)
	}
	m.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			m.Lock()
			delete(m.Watchers, w.id)
			m.Unlock()
		default:
			select {
			case w.res <- r:
			case <-time.After(timeout):
			}
		}
	}
}

func (m *registryImpl) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&m.options)
	}

	// add services
	m.Lock()
	for k, v := range getServices(m.options.Context) {
		s := m.Services[k]
		m.Services[k] = addServices(s, v)
	}
	m.Unlock()
	return nil
}

func (m *registryImpl) Options() registry.Options {
	return m.options
}

func (m *registryImpl) GetService(service string) ([]*registry.Service, error) {
	m.RLock()
	s, ok := m.Services[service]
	if !ok || len(s) == 0 {
		m.RUnlock()
		return nil, registry.ErrNotFound
	}
	m.RUnlock()
	return s, nil
}

func (m *registryImpl) ListServices() ([]*registry.Service, error) {
	m.RLock()
	var services []*registry.Service
	for _, service := range m.Services {
		services = append(services, service...)
	}
	m.RUnlock()
	return services, nil
}

func (m *registryImpl) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	go m.watch(&registry.Result{Action: "update", Service: s})

	m.Lock()
	services := addServices(m.Services[s.Name], []*registry.Service{s})
	m.Services[s.Name] = services
	m.Unlock()
	return nil
}

func (m *registryImpl) Deregister(s *registry.Service) error {
	go m.watch(&registry.Result{Action: "delete", Service: s})

	m.Lock()
	services := delServices(m.Services[s.Name], []*registry.Service{s})
	m.Services[s.Name] = services
	m.Unlock()
	return nil
}

func (m *registryImpl) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	w := &watcherImpl{
		exit: make(chan bool),
		res:  make(chan *registry.Result),
		id:   uuid.New().String(),
		wo:   wo,
	}

	m.Lock()
	m.Watchers[w.id] = w
	m.Unlock()
	return w, nil
}

func (m *registryImpl) String() string {
	return "memory"
}

// NewRegistry new
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	services := getServices(options.Context)
	if services == nil {
		services = make(map[string][]*registry.Service)
	}

	return &registryImpl{
		options:  options,
		Services: services,
		Watchers: make(map[string]*watcherImpl),
	}
}
