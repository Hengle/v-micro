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

	"github.com/fananchong/v-micro/client"
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/log"
)

var (
	// A value sent as a placeholder for the server's response value when the server
	// receives an invalid request. It is never decoded by the client since the Response
	// contains an error when it is used.
	invalidRequest = struct{}{}

	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

type methodType struct {
	sync.Mutex  // protects counters
	method      reflect.Method
	ArgType     reflect.Type
	ContextType reflect.Type
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

type request struct {
	msg  *codec.Message
	next *request // for free list in Server
}

type response struct {
	msg  *codec.Message
	next *response // for free list in Server
}

// router represents an RPC router.
type router struct {
	name         string
	mu           sync.Mutex // protects the serviceMap
	serviceMap   map[string]*service
	reqLock      sync.Mutex // protects freeReq
	freeReq      *request
	respLock     sync.Mutex // protects freeResp
	freeResp     *response
	hdlrWrappers []client.HandlerWrapper
}

func newRPCRouter() *router {
	return &router{
		serviceMap: make(map[string]*service),
	}
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
	var argType, contextType reflect.Type

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}

	switch mtype.NumIn() {
	case 3:
		// method that takes a context
		argType = mtype.In(2)
		contextType = mtype.In(1)
	default:
		log.Info("method", mname, "of", mtype, "has wrong number of ins:", mtype.NumIn())
		return nil
	}

	// First arg need not be a pointer.
	if !isExportedOrBuiltinType(argType) {
		log.Info(mname, "argument type not exported:", argType)
		return nil
	}

	return &methodType{method: method, ArgType: argType, ContextType: contextType}
}

func (router *router) sendResponse(sending sync.Locker, req *request, reply interface{}, cc codec.Writer, last bool) error {
	msg := new(codec.Message)
	msg.Type = codec.Response
	msg.Service = req.msg.Service
	msg.Method = req.msg.Method
	resp := router.getResponse()
	resp.msg = msg
	sending.Lock()
	err := cc.Write(resp.msg, reply)
	sending.Unlock()
	router.freeResponse(resp)
	return err
}

func (s *service) call(ctx context.Context, router *router, mtype *methodType, argv reflect.Value) error {
	function := mtype.method.Func
	fn := func(ctx context.Context, rsp interface{}) error {
		function.Call([]reflect.Value{s.rcvr, mtype.prepareContext(ctx), reflect.ValueOf(rsp)})
		return nil
	}

	// wrap the handler
	for i := len(router.hdlrWrappers); i > 0; i-- {
		fn = router.hdlrWrappers[i-1](fn)
	}

	// execute handler
	fn(ctx, argv.Interface())

	return nil
}

func (m *methodType) prepareContext(ctx context.Context) reflect.Value {
	if contextv := reflect.ValueOf(ctx); contextv.IsValid() {
		return contextv
	}
	return reflect.Zero(m.ContextType)
}

func (router *router) getRequest() *request {
	router.reqLock.Lock()
	req := router.freeReq
	if req == nil {
		req = new(request)
	} else {
		router.freeReq = req.next
		*req = request{}
	}
	router.reqLock.Unlock()
	return req
}

func (router *router) freeRequest(req *request) {
	router.reqLock.Lock()
	req.next = router.freeReq
	router.freeReq = req
	router.reqLock.Unlock()
}

func (router *router) getResponse() *response {
	router.respLock.Lock()
	resp := router.freeResp
	if resp == nil {
		resp = new(response)
	} else {
		router.freeResp = resp.next
		*resp = response{}
	}
	router.respLock.Unlock()
	return resp
}

func (router *router) freeResponse(resp *response) {
	router.respLock.Lock()
	resp.next = router.freeResp
	router.freeResp = resp
	router.respLock.Unlock()
}

func (router *router) readRequest(r *rpcRequest) (service *service, mtype *methodType, argv reflect.Value, keepReading bool, err error) {
	cc := r.codec
	service, mtype, keepReading, err = router.readHeader(r)
	if err != nil {
		log.Error(err)
		if !keepReading {
			return
		}
		// discard body
		cc.ReadBody(nil)
		return
	}

	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = cc.ReadBody(argv.Interface()); err != nil {
		log.Error(err)
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}
	return
}

func (router *router) readHeader(r *rpcRequest) (service *service, mtype *methodType, keepReading bool, err error) {
	// We read the header successfully. If we see an error now,
	// we can still recover and move on to the next request.
	keepReading = true

	serviceMethod := strings.Split(r.Method(), ".")
	if len(serviceMethod) != 2 {
		err = errors.New("rpc: service/method request ill-formed: " + r.Method())
		return
	}
	// Look up the request.
	router.mu.Lock()
	service = router.serviceMap[serviceMethod[0]]
	router.mu.Unlock()
	if service == nil {
		err = errors.New("rpc: can't find service " + serviceMethod[0])
		return
	}
	mtype = service.method[serviceMethod[1]]
	if mtype == nil {
		err = errors.New("rpc: can't find method " + serviceMethod[1])
	}
	return
}

func (router *router) Handle(h interface{}) error {
	name := reflect.Indirect(reflect.ValueOf(h)).Type().Name()

	router.mu.Lock()
	defer router.mu.Unlock()
	if router.serviceMap == nil {
		router.serviceMap = make(map[string]*service)
	}

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
	if _, present := router.serviceMap[name]; present {
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
	router.serviceMap[s.name] = s
	return nil
}

func (router *router) ServeRequest(ctx context.Context, r *rpcRequest) error {
	defer func() {
		if re := recover(); re != nil {
			log.Info("panic recovered: ", re)
			log.Info(string(debug.Stack()))
		}
	}()

	service, mtype, argv, keepReading, err := router.readRequest(r)
	if err != nil {
		log.Error(err)
		if !keepReading {
			return err
		}
		return err
	}
	return service.call(ctx, router, mtype, argv)
}
