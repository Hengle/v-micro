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
	buf    *buffer.ReadWriteCloser
}

var (
	// DefaultCodecs default codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"":                     proto.NewCodec, // default
		"application/protobuf": proto.NewCodec,
	}
)

func newRPCCodec(rbuf []byte, client transport.Client, c codec.NewCodec) codec.Codec {
	rwc := &buffer.ReadWriteCloser{
		RBuf: bytes.NewBuffer(rbuf),
		WBuf: bytes.NewBuffer(nil),
	}
	r := &rpcCodec{
		buf:    rwc,
		client: client,
		codec:  c(rwc),
	}
	return r
}

func (c *rpcCodec) Write(m *codec.Message, body interface{}) (err error) {
	c.buf.WBuf.Reset()

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
