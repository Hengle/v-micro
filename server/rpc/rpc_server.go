package rpc

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	hcodec "github.com/fananchong/v-micro/internal/codec"
	mnet "github.com/fananchong/v-micro/internal/net"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/server"
	common "github.com/fananchong/v-micro/server/internal"
	"github.com/fananchong/v-micro/transport"
)

type rpcServer struct {
	router     *router
	exit       chan chan error
	exitAccept chan int
	opts       server.Options
	registered int64
	ts         transport.Listener
}

func newRPCServer(opts ...server.Option) server.Server {
	s := &rpcServer{
		opts: server.Options{
			Codecs:   make(map[string]codec.NewCodec),
			Metadata: map[string]string{},
		},
		router:     newRPCRouter(),
		exit:       make(chan chan error),
		exitAccept: make(chan int),
	}
	common.InitOptions(&s.opts, opts...)
	return s
}

func (s *rpcServer) Options() server.Options {
	opts := s.opts
	return opts
}

func (s *rpcServer) Init(opts ...server.Option) error {
	common.InitOptions(&s.opts, opts...)
	s.router.hdlrWrappers = s.opts.HdlrWrappers
	return nil
}

func (s *rpcServer) Handle(h interface{}) error {
	return s.router.Handle(h)
}

func (s *rpcServer) Start() (err error) {
	config := s.Options()

	// start listening on the transport
	if s.ts, err = config.Transport.Listen(config.Address); err != nil {
		log.Error(err)
		return
	}

	log.Infof("Transport [%s] Listening on %s", config.Transport.String(), s.ts.Addr())

	// swap address
	s.opts.Address = s.ts.Addr()

	// use RegisterCheck func before register
	if err = s.opts.RegisterCheck(s.opts.Context); err != nil {
		log.Infof("Server %s-%s register check error: %s", config.Name, config.ID, err)
	} else {
		// announce self to the world
		if err = s.register(); err != nil {
			log.Infof("Server %s-%s register error: %s", config.Name, config.ID, err)
		}
	}

	// start acceptting on the listener
	go s.accept()

	// ttl for register
	go s.registerTTL()

	return nil
}

func (s *rpcServer) Stop() error {
	if s.ts != nil {
		s.ts.Close()
	}
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *rpcServer) String() string {
	return "rpc"
}

func (s *rpcServer) register() (err error) {
	// parse address for host, port
	config := s.Options()
	addr, port, err := common.ExtractAdvertiseIP(config)
	if err != nil {
		log.Error(err)
		return err
	}

	// make copy of metadata
	md := make(metadata.Metadata)
	for k, v := range config.Metadata {
		md[k] = v
	}

	// register service
	node := &registry.Node{
		ID:       config.Name + "-" + config.ID,
		Address:  mnet.HostPort(addr, port),
		Metadata: md,
	}

	node.Metadata["transport"] = config.Transport.String()
	node.Metadata["server"] = s.String()
	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["PID"] = strconv.Itoa(os.Getpid())

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	if atomic.LoadInt64(&s.registered) == 0 {
		log.Infof("Registry [%s] Registering node: %s", config.Registry.String(), node.ID)
	}

	// create registry options
	rOpts := []registry.RegisterOption{registry.RegisterTTL(config.RegisterTTL)}

	if err := config.Registry.Register(service, rOpts...); err != nil {
		log.Error(err)
		return err
	}

	// already registered? don't need to register subscribers
	registered := atomic.SwapInt64(&s.registered, 1)
	if registered != 0 {
		return nil
	}

	return
}

func (s *rpcServer) deregister() (err error) {
	config := s.Options()
	addr, port, err := common.ExtractAdvertiseIP(config)
	if err != nil {
		log.Error(err)
		return err
	}

	node := &registry.Node{
		ID:      config.Name + "-" + config.ID,
		Address: mnet.HostPort(addr, port),
	}

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	log.Infof("Registry [%s] Deregistering node: %s", config.Registry.String(), node.ID)
	if err := config.Registry.Deregister(service); err != nil {
		log.Error(err)
		return err
	}

	registered := atomic.SwapInt64(&s.registered, 0)
	if registered == 0 {
		return
	}

	return
}

func (s *rpcServer) accept() {
	for {
		// listen for connections
		err := s.ts.Accept(s.serveConn)

		select {
		// check if we're supposed to exit
		case <-s.exitAccept:
			return
		// check the error and backoff
		default:
			if err != nil {
				log.Infof("Accept error: %v", err)
				time.Sleep(time.Second)
				continue
			}
		}
		// no error just exit
		return
	}
}

// serveConn serves a single connection
func (s *rpcServer) serveConn(sock transport.Socket) {
	defer func() {
		// close socket
		sock.Close()

		if r := recover(); r != nil {
			log.Info("panic recovered: ", r)
			log.Info(string(debug.Stack()))
		}
	}()

	for {
		var msg transport.Message
		if err := sock.Recv(&msg); err != nil {
			log.Error(err)
			return
		}

		// strip our headers
		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}

		// create new context
		ctx := metadata.NewContext(context.Background(), hdr)

		// we use this Content-Type header to identify the codec needed
		ct := msg.Header[metadata.CONTENTTYPE]
		var cf codec.NewCodec
		var err error
		if cf, err = s.newCodec(ct); err != nil {
			log.Error(err)
			continue
		}

		rcodec := newRPCCodec(msg.Body, sock, cf)

		// internal request
		request := &rpcRequest{
			id:          hcodec.GetHeader(metadata.MSGID, msg.Header),
			method:      hcodec.GetHeader(metadata.METHOD, msg.Header),
			contentType: ct,
			codec:       rcodec,
			header:      msg.Header,
			body:        msg.Body,
		}

		// internal response
		response := &rpcResponse{
			codec: rcodec,
		}

		// serve the actual request using the request router
		if err := s.router.ServeRequest(ctx, request, response); err != nil {
			log.Error(err)
			continue
		}
	}
}

func (s *rpcServer) registerTTL() {
	t := new(time.Ticker)

	// only process if it exists
	if s.opts.RegisterInterval > time.Duration(0) {
		// new ticker
		t = time.NewTicker(s.opts.RegisterInterval)
	}

	config := s.Options()

	var ch chan error
Loop:
	for {
		select {
		// register self on interval
		case <-t.C:
			if err := s.opts.RegisterCheck(s.opts.Context); err != nil && atomic.LoadInt64(&s.registered) != 0 {
				log.Infof("Server %s-%s register check error: %s, deregister it", config.Name, config.ID, err)
				// deregister self in case of error
				if err := s.deregister(); err != nil {
					log.Infof("Server %s-%s deregister error: %s", config.Name, config.ID, err)
				}
			} else {
				if err := s.register(); err != nil {
					log.Infof("Server %s-%s register error: %s", config.Name, config.ID, err)
				}
			}
		// wait for exit
		case ch = <-s.exit:
			close(s.exitAccept)
			t.Stop()
			break Loop
		}
	}

	// deregister self
	var err error
	if err = s.deregister(); err != nil {
		log.Infof("Server %s-%s deregister error: %s", config.Name, config.ID, err)
	}
	ch <- err
}

func (s *rpcServer) newCodec(contentType string) (codec.NewCodec, error) {
	if cf, ok := s.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

// NewServer returns a new server with options passed in
func NewServer(opts ...server.Option) server.Server {
	return newRPCServer(opts...)
}
