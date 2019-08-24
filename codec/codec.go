// Package codec is an interface for encoding messages
package codec

import (
	"io"
)

// NewCodec Takes in a connection/buffer and returns a new Codec
type NewCodec func(io.ReadWriteCloser) Codec

// Codec encodes/decodes various types of messages used within v-micro.
// ReadHeader and ReadBody are called in pairs to read requests/responses
// from the connection. Close is called when finished with the
// connection. ReadBody may be called with a nil argument to force the
// body to be read and discarded.
type Codec interface {
	Reader
	Writer
	Close() error
	String() string
}

// Reader reader
type Reader interface {
	ReadHeader(*Message) error
	ReadBody(interface{}) error
}

// Writer writer
type Writer interface {
	Write(*Message, interface{}) error
}

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

// Message represents detailed information about
// the communication, likely followed by the body.
// In the case of an error, body may be nil.
type Message struct {
	Service string // 服务
	Method  string // 方法名， class:method

	// The values read from the socket
	Header map[string]string
	Body   []byte
}
