package rpc

import (
	"bytes"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/codec/proto"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/internal/buffer"
	hcodec "github.com/fananchong/v-micro/internal/codec"
	"github.com/fananchong/v-micro/transport"
)

type rpcCodec struct {
	client transport.Client
	codec  codec.Codec
	req    *transport.Message
	buf    *buffer.ReadWriteCloser
}

var (
	// DefaultCodecs default codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"application/protobuf": proto.NewCodec,
	}
)

func newRPCCodec(req *transport.Message, client transport.Client, c codec.NewCodec) codec.Codec {
	rwc := &buffer.ReadWriteCloser{
		RBuf: bytes.NewBuffer(req.Body),
		WBuf: bytes.NewBuffer(nil),
	}
	r := &rpcCodec{
		buf:    rwc,
		client: client,
		codec:  c(rwc),
		req:    req,
	}
	return r
}

func (c *rpcCodec) Write(m *codec.Message, body interface{}) (err error) {
	c.buf.WBuf.Reset()

	// create header
	if m.Header == nil {
		m.Header = map[string]string{}
	}

	// copy original header
	for k, v := range c.req.Header {
		m.Header[k] = v
	}

	// set the headers
	hcodec.SetHeaders(m, m)

	// write to codec
	if err = c.codec.Write(m, body); err != nil {
		log.Error(err)
		return
	}

	// set body
	m.Body = c.buf.WBuf.Bytes()

	// create new transport message
	msg := transport.Message{
		Header: m.Header,
		Body:   m.Body,
	}
	// send the request
	if err = c.client.Send(&msg); err != nil {
		log.Error(err)
		return
	}
	return
}

func (c *rpcCodec) ReadHeader(r *codec.Message) (err error) {
	return
}

func (c *rpcCodec) ReadBody(b interface{}) (err error) {
	// don't read empty body
	if len(c.req.Body) == 0 {
		return nil
	}

	// read body
	if err = c.codec.ReadBody(b); err != nil {
		log.Error(err)
		return
	}
	return
}

func (c *rpcCodec) Close() error {
	c.buf.Close()
	c.codec.Close()
	return nil
}

func (c *rpcCodec) String() string {
	return "rpc"
}
