package rpc

import (
	"reflect"

	"github.com/fananchong/v-micro/server"
)

type handlerImpl struct {
	name    string
	handler interface{}
}

func newHandler(handler interface{}) server.Handler {
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()
	return &handlerImpl{
		name:    name,
		handler: handler,
	}
}

func (r *handlerImpl) Name() string {
	return r.name
}

func (r *handlerImpl) Handler() interface{} {
	return r.handler
}
