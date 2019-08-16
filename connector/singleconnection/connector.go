package singleconnection

import (
	"fmt"
	"sync"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/connector"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/transport"
)

type singleConnection struct {
	opts  connector.Options
	conns sync.Map
	mutex sync.Mutex
}

func (s *singleConnection) Init(opts ...connector.Option) (err error) {
	for _, o := range opts {
		o(&s.opts)
	}
	return
}

func (s *singleConnection) Options() connector.Options {
	return s.opts
}

func (s *singleConnection) Get(node *registry.Node) (transport.Client, error) {
	if c, ok := s.conns.Load(node.ID); ok {
		return c.(transport.Client), nil
	}
	return s.createConnect(node)
}

func (s *singleConnection) createConnect(node *registry.Node) (transport.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 1. If it already exists, return directly
	if c, ok := s.conns.Load(node.ID); ok {
		return c.(transport.Client), nil
	}

	log.Infof("Connect node, id:%s address:%s metadata:%v", node.ID, node.Address, node.Metadata)

	// 2. Dial a connect
	dOpts := []transport.DialOption{
		transport.OnClose(func() { s.onConnectClose(node.ID) }),
	}
	c, err := s.opts.Transport.Dial(node.Address, dOpts...)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}

	// 3. Add cache
	s.conns.Store(node.ID, c)

	// 4. Return
	return c, nil
}

func (s *singleConnection) onConnectClose(id string) {
	log.Infof("Disconnect node, id:%s", id)
	s.conns.Delete(id)
}

func (s *singleConnection) Close() (err error) {
	var ids []string
	s.conns.Range(func(key interface{}, value interface{}) bool {
		ids = append(ids, key.(string))
		value.(transport.Client).Close()
		return true
	})
	for _, id := range ids {
		s.conns.Delete(id)
	}
	return
}

func (s *singleConnection) String() string {
	return "singleconnection"
}

// NewConnector new
func NewConnector(opts ...connector.Option) connector.Connector {
	s := &singleConnection{}
	for _, o := range opts {
		o(&s.opts)
	}
	return s
}
