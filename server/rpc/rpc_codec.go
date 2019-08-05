package rpc

import (
	"bytes"
	"net/rpc"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/codec/proto"
	"github.com/fananchong/v-micro/transport"
)

var (
	// DefaultContentType default content type
	DefaultContentType = "application/protobuf"

	// DefaultCodecs default codecs
	DefaultCodecs = map[string]codec.NewCodec{
		"application/protobuf": proto.NewCodec,
	}
)

type rcpCodec struct {
	socket transport.Socket
	codec  codec.Codec
	req    *transport.Message
	buf    *readWriteCloser
}

func (c *rcpCodec) ReadRequestHeader(r *rpc.Request) error {
	// r.ServiceMethod
	// var n uint16
	// binary.Read(c.rwc, binary.LittleEndian, &n)
	// temp := make([]byte, n)
	// c.rwc.Read(temp[:])
	// r.ServiceMethod = string(temp)
	return nil
}

func (c *rcpCodec) ReadRequestBody(body interface{}) error {
	// body string
	// var n uint16
	// binary.Read(c.rwc, binary.LittleEndian, &n)
	// temp := make([]byte, n)
	// c.rwc.Read(temp[:])
	// *body.(*string) = string(temp)
	return nil
}

func (c *rcpCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	// write r.ServiceMethod
	// binary.Write(c.rwc, binary.LittleEndian, uint16(len(r.ServiceMethod)))
	// c.rwc.Write([]byte(r.ServiceMethod))
	// write r.Error
	// binary.Write(c.rwc, binary.LittleEndian, uint16(len(r.Error)))
	// c.rwc.Write([]byte(r.Error))
	// write body
	// data := []byte(*body.(*string))
	// binary.Write(c.rwc, binary.LittleEndian, uint16(len(data)))
	// c.rwc.Write(data)
	return nil
}

func (c *rcpCodec) Close() (err error) {
	// return c.rwc.Close()
	return
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
