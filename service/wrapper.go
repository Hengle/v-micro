package service

import (
	"context"

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/common/metadata"
)

// ClientWrapper client wrapper
type ClientWrapper struct {
	client.Client
	Headers metadata.Metadata
}

func (c *ClientWrapper) setHeaders(ctx context.Context) context.Context {
	// copy metadata
	mda, _ := metadata.FromContext(ctx)
	md := metadata.Copy(mda)

	// set headers
	for k, v := range c.Headers {
		if _, ok := md[k]; !ok {
			md[k] = v
		}
	}

	return metadata.NewContext(ctx, md)
}

// Call call
func (c *ClientWrapper) Call(ctx context.Context, req client.Request, opts ...client.CallOption) error {
	ctx = c.setHeaders(ctx)
	return c.Client.Call(ctx, req, opts...)
}
