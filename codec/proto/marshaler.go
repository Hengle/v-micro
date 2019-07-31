package proto

import (
	"github.com/golang/protobuf/proto"
)

type marshalerImpl struct{}

func (marshalerImpl) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (marshalerImpl) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (marshalerImpl) String() string {
	return "proto"
}
