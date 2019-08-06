package rpc

import (
	"bytes"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/codec/proto"
	"github.com/fananchong/v-micro/transport"
	"github.com/pkg/errors"
)

var (
	// DefaultContentType default content type
	DefaultContentType = "application/protobuf"

	// DefaultCodecs default codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"application/protobuf": proto.NewCodec,
	}
)

type rpcCodec struct {
	socket transport.Socket
	codec  codec.Codec
	req    *transport.Message
	buf    *readWriteCloser
}

func newRPCCodec(req *transport.Message, socket transport.Socket, c codec.NewCodec) codec.Codec {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(req.Body),
		wbuf: bytes.NewBuffer(nil),
	}
	r := &rpcCodec{
		buf:    rwc,
		codec:  c(rwc),
		req:    req,
		socket: socket,
	}
	return r
}

func (c *rpcCodec) ReadHeader(r *codec.Message, t codec.MessageType) (err error) {
	// the initial message
	m := codec.Message{
		Header: c.req.Header,
		Body:   c.req.Body,
	}

	// set some internal things
	getHeaders(&m)

	// read header via codec
	if err = c.codec.ReadHeader(&m, codec.Request); err != nil {
		return
	}

	// set message
	*r = m

	return
}

func (c *rpcCodec) ReadBody(b interface{}) error {
	// don't read empty body
	if len(c.req.Body) == 0 {
		return nil
	}
	// decode the usual way
	return c.codec.ReadBody(b)
}

func (c *rpcCodec) Write(r *codec.Message, b interface{}) error {
	c.buf.wbuf.Reset()

	// create a new message
	m := &codec.Message{
		Service: r.Service,
		Method:  r.Method,
		ID:      r.ID,
		Error:   r.Error,
		Type:    r.Type,
		Header:  r.Header,
	}

	if m.Header == nil {
		m.Header = map[string]string{}
	}

	setHeaders(m, r)

	// the body being sent
	var body []byte

	// write the body to codec
	if err := c.codec.Write(m, b); err != nil {
		c.buf.wbuf.Reset()

		// write an error if it failed
		m.Error = errors.Wrapf(err, "Unable to encode body").Error()
		m.Header["Micro-Error"] = m.Error
		// no body to write
		if err := c.codec.Write(m, nil); err != nil {
			return err
		}
	}

	// set the body
	body = c.buf.wbuf.Bytes()

	// Set content type if theres content
	if len(body) > 0 {
		m.Header["Content-Type"] = c.req.Header["Content-Type"]
	}

	// send on the socket
	return c.socket.Send(&transport.Message{
		Header: m.Header,
		Body:   body,
	})
}

func (c *rpcCodec) Close() (err error) {
	// return c.rwc.Close()
	return
}

func (c *rpcCodec) String() string {
	return "rpc"
}

// ======================= readWriteCloser =======================

type readWriteCloser struct {
	wbuf *bytes.Buffer
	rbuf *bytes.Buffer
}

func (rwc *readWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.rbuf.Read(p)
}

func (rwc *readWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.wbuf.Write(p)
}

func (rwc *readWriteCloser) Close() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()
	return nil
}

// ======================= misc =======================

func getHeader(hdr string, md map[string]string) string {
	if hd := md[hdr]; len(hd) > 0 {
		return hd
	}
	return ""
}

func getHeaders(m *codec.Message) {
	m.ID = getHeader("Micro-Id", m.Header)
	m.Method = getHeader("Micro-Method", m.Header)
	m.Service = getHeader("Micro-Service", m.Header)
}

func setHeaders(m, r *codec.Message) {
	set := func(hdr, v string) {
		if len(v) == 0 {
			return
		}
		m.Header[hdr] = v
	}

	// set headers
	set("Micro-Id", r.ID)
	set("Micro-Service", r.Service)
	set("Micro-Method", r.Method)
	set("Micro-Error", r.Error)
}
