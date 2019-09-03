package rpc

import (
	"context"
	"fmt"

	"github.com/fananchong/v-micro/client"
	common "github.com/fananchong/v-micro/client/internal"
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/metadata"
	"github.com/fananchong/v-micro/connector"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/transport"
)

type rpcClient struct {
	opts   client.Options
	router *router
}

func newRPCClient(opts ...client.Option) client.Client {
	rc := &rpcClient{
		opts: client.Options{
			Codecs: make(map[string]codec.NewCodec),
		},
		router: newRPCRouter(),
	}

	common.InitOptions(&rc.opts, opts...)

	c := client.Client(rc)

	// wrap in reverse
	for i := len(rc.opts.Wrappers); i > 0; i-- {
		c = rc.opts.Wrappers[i-1](c)
	}

	return c
}

func (r *rpcClient) newCodec(contentType string) (codec.NewCodec, error) {
	if c, ok := r.opts.Codecs[contentType]; ok {
		return c, nil
	}
	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (r *rpcClient) call(ctx context.Context, node *registry.Node, req client.Request, opts client.CallOptions) (err error) {
	msg := &codec.Message{
		Method: req.Method(),
		Header: make(map[string]string),
	}

	md, ok := metadata.FromContext(ctx)
	if ok {
		for k, v := range md {
			msg.Header[k] = v
		}
	}

	// set the content type for the request
	msg.Header[metadata.CONTENTTYPE] = req.ContentType()

	var cf codec.NewCodec
	if cf, err = r.newCodec(req.ContentType()); err != nil {
		return
	}

	var cli transport.Client
	if cli, err = r.Options().Connector.Get(node); err != nil {
		return
	}

	cc := newRPCCodec(nil, cli, cf)

	if err = cc.Write(msg, req.Body()); err != nil {
		return
	}

	return
}

func (r *rpcClient) Init(opts ...client.Option) error {
	common.InitOptions(&r.opts, opts...)

	connectorOpts := []connector.Option{
		connector.OnConnected(r.AsyncRecv),
	}
	return r.opts.Connector.Init(connectorOpts...)
}

func (r *rpcClient) Options() client.Options {
	return r.opts
}

func (r *rpcClient) Handle(h interface{}) (err error) {
	return r.router.Handle(h)
}

func (r *rpcClient) next(request client.Request, opts client.CallOptions) (selector.Next, error) {
	service := request.Service()

	// get next nodes from the selector
	next, err := r.opts.Selector.Select(service, opts.SelectOptions...)
	if err != nil {
		err = fmt.Errorf("service %s: %v", service, err.Error())
		return nil, err
	}

	return next, nil
}

func (r *rpcClient) NewRequest(service, method string, req interface{}) client.Request {
	return newRequest(service, method, req, r.opts.ContentType)
}

func (r *rpcClient) Call(ctx context.Context, request client.Request, opts ...client.CallOption) (err error) {
	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := r.next(request, callOpts)
	if err != nil {
		log.Error(err)
		return err
	}

	// make copy of call method
	rcall := r.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	call := func() (error, bool) {
		// select next node
		node, err := next()
		service := request.Service()
		if err != nil {
			return fmt.Errorf("service %s: %s", service, err.Error()), false
		}
		// make the call
		err = rcall(ctx, node, request, callOpts)
		return err, true
	}

	var canTry bool
	for {
		if err, canTry = call(); err != nil && canTry {
			log.Errorf("call fail and try again, err:%s", err.Error())
			continue
		} else {
			break
		}
	}
	return
}

func (r *rpcClient) Broadcast(ctx context.Context, request client.Request, opts ...client.CallOption) {
	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	opt := client.WithSelectOption(selector.WithStrategy(selector.StatefulRoundRobin))
	opt(&callOpts)

	if md, ok := metadata.FromContext(ctx); ok {
		md["Micro-Broadcast"] = "true"
	} else {
		md := make(metadata.Metadata)
		md["Micro-Broadcast"] = "true"
		ctx = metadata.NewContext(ctx, md)
	}

	next, err := r.next(request, callOpts)
	if err != nil {
		return
	}

	// make copy of call method
	rcall := r.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	for {
		if node, _ := next(); node != nil {
			// make the call
			if err := rcall(ctx, node, request, callOpts); err != nil {
				log.Error(err)
			}
		} else {
			break
		}
	}
}

func (r *rpcClient) String() string {
	return "rpc"
}

// NewClient Creates a new client with the options passed in
func NewClient(opts ...client.Option) client.Client {
	return newRPCClient(opts...)
}
