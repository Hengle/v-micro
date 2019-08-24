package rpc

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Meh, we need to get rid of this shit

import (
	"context"
	"errors"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/fananchong/v-micro/common/log"
)

type methodType struct {
	method      reflect.Method
	ArgType     reflect.Type
	ReplyType   reflect.Type
	ContextType reflect.Type
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

// router represents an RPC router.
type router struct {
	name       string
	serviceMap sync.Map
}

func newRPCRouter() *router {
	return &router{}
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// prepareMethod returns a methodType for the provided method or nil
// in case if the method was unsuitable.
func prepareMethod(method reflect.Method) *methodType {
	mtype := method.Type
	mname := method.Name
	var replyType, argType, contextType reflect.Type

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}

	switch mtype.NumIn() {
	case 4:
		// method that takes a context
		argType = mtype.In(2)
		replyType = mtype.In(3)
		contextType = mtype.In(1)
	default:
		log.Info("method", mname, "of", mtype, "has wrong number of ins:", mtype.NumIn())
		return nil
	}

	if argType.Kind() != reflect.Ptr {
		log.Info("method", mname, "arg type not a pointer:", argType)
		return nil
	}

	// First arg need not be a pointer.
	if !isExportedOrBuiltinType(argType) {
		log.Info(mname, "argument type not exported:", argType)
		return nil
	}

	if replyType.Kind() != reflect.Ptr {
		log.Info("method", mname, "reply type not a pointer:", replyType)
		return nil
	}

	// Reply type must be exported.
	if !isExportedOrBuiltinType(replyType) {
		log.Info("method", mname, "reply type not exported:", replyType)
		return nil
	}

	return &methodType{method: method, ArgType: argType, ReplyType: replyType, ContextType: contextType}
}

func (s *service) call(ctx context.Context, router *router, mtype *methodType, argv, replyv reflect.Value) error {
	function := mtype.method.Func
	fn := func(ctx context.Context, req, rsp interface{}) error {
		function.Call([]reflect.Value{s.rcvr, mtype.prepareContext(ctx), reflect.ValueOf(req), reflect.ValueOf(rsp)})
		return nil
	}

	// execute handler
	fn(ctx, argv.Interface(), replyv.Interface())

	return nil
}

func (m *methodType) prepareContext(ctx context.Context) reflect.Value {
	if contextv := reflect.ValueOf(ctx); contextv.IsValid() {
		return contextv
	}
	return reflect.Zero(m.ContextType)
}

func (router *router) readRequest(r *rpcRequest) (service *service, mtype *methodType, argv, replyv reflect.Value, err error) {
	service, mtype, err = router.readHeader(r)
	if err != nil {
		log.Error(err)
		return
	}

	// Decode the argument value.
	cc0 := r.Codec(0)
	replyv = reflect.New(mtype.ReplyType.Elem())
	if err = cc0.ReadBody(replyv.Interface()); err != nil {
		log.Error(err)
		return
	}

	cc1 := r.Codec(1)
	argv = reflect.New(mtype.ArgType.Elem())
	if err = cc1.ReadBody(argv.Interface()); err != nil {
		log.Error(err)
		return
	}

	return
}

func (router *router) readHeader(r *rpcRequest) (s *service, mtype *methodType, err error) {
	serviceMethod := strings.Split(r.Method(), ".")
	if len(serviceMethod) != 2 {
		err = errors.New("rpc: service/method request ill-formed: " + r.Method())
		return
	}
	// Look up the request.
	tmpV, ok := router.serviceMap.Load(serviceMethod[0])
	if !ok {
		err = errors.New("rpc: can't find service " + serviceMethod[0])
		return
	}
	s = tmpV.(*service)
	mtype = s.method[serviceMethod[1]]
	if mtype == nil {
		err = errors.New("rpc: can't find method " + serviceMethod[1])
	}
	return
}

func (router *router) Handle(h interface{}) error {
	name := reflect.Indirect(reflect.ValueOf(h)).Type().Name()

	if len(name) == 0 {
		return errors.New("rpc.Handle: handler has no name")
	}
	if !isExported(name) {
		return errors.New("rpc.Handle: type " + name + " is not exported")
	}

	s := new(service)
	s.typ = reflect.TypeOf(h)
	s.rcvr = reflect.ValueOf(h)

	// check name
	if _, present := router.serviceMap.Load(name); present {
		return errors.New("rpc.Handle: service already defined: " + name)
	}

	s.name = name
	s.method = make(map[string]*methodType)

	// Install the methods
	for m := 0; m < s.typ.NumMethod(); m++ {
		method := s.typ.Method(m)
		if mt := prepareMethod(method); mt != nil {
			s.method[method.Name] = mt
		}
	}

	// Check there are methods
	if len(s.method) == 0 {
		return errors.New("rpc Register: type " + s.name + " has no exported methods of suitable type")
	}

	// save handler
	router.serviceMap.Store(s.name, s)
	return nil
}

func (router *router) ServeRequest(ctx context.Context, r *rpcRequest) error {
	defer func() {
		if re := recover(); re != nil {
			log.Info("panic recovered: ", re)
			log.Info(string(debug.Stack()))
		}
	}()

	service, mtype, argv, replyv, err := router.readRequest(r)
	if err != nil {
		log.Error(err)
		return err
	}
	return service.call(ctx, router, mtype, argv, replyv)
}
