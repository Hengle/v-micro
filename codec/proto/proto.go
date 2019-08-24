// Package proto provides a proto codec
package proto

import (
	"io"
	"io/ioutil"

	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
	"github.com/golang/protobuf/proto"
)

type codecImpl struct {
	Conn io.ReadWriteCloser
}

func (c *codecImpl) ReadHeader(m *codec.Message) error {
	return nil
}

func (c *codecImpl) ReadBody(b interface{}) error {
	if b == nil {
		return nil
	}
	buf, err := ioutil.ReadAll(c.Conn)
	if err != nil {
		log.Error(err)
		return err
	}
	return proto.Unmarshal(buf, b.(proto.Message))
}

func (c *codecImpl) Write(m *codec.Message, b interface{}) error {
	p, ok := b.(proto.Message)
	if !ok {
		return nil
	}
	buf, err := proto.Marshal(p)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = c.Conn.Write(buf)
	return err
}

func (c *codecImpl) Close() error {
	return c.Conn.Close()
}

func (c *codecImpl) String() string {
	return "proto"
}

// NewCodec new
func NewCodec(c io.ReadWriteCloser) codec.Codec {
	return &codecImpl{
		Conn: c,
	}
}
