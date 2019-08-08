package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/fananchong/v-micro/client"
	common "github.com/fananchong/v-micro/client/internal"
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/selector"
)

type rpcClient struct {
	once sync.Once
	opts client.Options
}

func newRPCClient(opt ...client.Option) client.Client {
	opts := common.NewOptions(opt...)

	rc := &rpcClient{
		once: sync.Once{},
		opts: opts,
	}

	c := client.Client(rc)

	// wrap in reverse
	for i := len(opts.Wrappers); i > 0; i-- {
		c = opts.Wrappers[i-1](c)
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
	// TODO
	return
}

func (r *rpcClient) Init(opts ...client.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

func (r *rpcClient) Options() client.Options {
	return r.opts
}

func (r *rpcClient) Handle(interface{}) (err error) {
	// TODO
	return
}

func (r *rpcClient) next(request client.Request, opts client.CallOptions) (selector.Next, error) {
	service := request.Service()

	// get next nodes from the selector
	next, err := r.opts.Selector.Select(service, opts.SelectOptions...)
	if err != nil && err == selector.ErrNotFound {
		err = fmt.Errorf("service %s: %v", service, err.Error())
		log.Error(err)
		return nil, err
	} else if err != nil {
		err = fmt.Errorf("error selecting %s node: %v", service, err.Error())
		log.Error(err)
		return nil, err
	}

	return next, nil
}

func (r *rpcClient) Call(ctx context.Context, request client.Request, opts ...client.CallOption) (err error) {
	// TODO
	return
}

func (r *rpcClient) String() string {
	return "rpc"
}
