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

var (
	// defaultContentType default content type
	defaultContentType = "application/protobuf"

	// DefaultCodecs default codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"application/protobuf": proto.NewCodec,
	}
)

type rpcCodec struct {
	socket transport.Socket
	codec  codec.Codec
	buf    *buffer.ReadWriteCloser
}

func newRPCCodec(rbuf []byte, socket transport.Socket, c codec.NewCodec) codec.Codec {
	rwc := &buffer.ReadWriteCloser{
		RBuf: bytes.NewBuffer(rbuf),
		WBuf: bytes.NewBuffer(nil),
	}
	r := &rpcCodec{
		buf:    rwc,
		codec:  c(rwc),
		socket: socket,
	}
	return r
}

func (c *rpcCodec) ReadHeader(r *codec.Message) (err error) {
	return
}

func (c *rpcCodec) ReadBody(b interface{}) error {
	return c.codec.ReadBody(b)
}

func (c *rpcCodec) Write(m *codec.Message, b interface{}) error {
	c.buf.WBuf.Reset()

	hcodec.SetHeaders(m, m)

	// write the body to codec
	if err := c.codec.Write(m, b); err != nil {
		c.buf.WBuf.Reset()
		log.Error(err)
		return err
	}

	// send on the socket
	return c.socket.Send(&transport.Message{
		Header: m.Header,
		Body:   c.buf.WBuf.Bytes(),
	})
}

func (c *rpcCodec) Close() (err error) {
	c.buf.Close()
	c.codec.Close()
	return nil
}

func (c *rpcCodec) String() string {
	return "rpc"
}
