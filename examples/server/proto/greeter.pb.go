// Code generated by protoc-gen-go. DO NOT EDIT.
// source: greeter.proto

package mypackage

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/fananchong/v-micro/client"
	server "github.com/fananchong/v-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_e585294ab3f34af5, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Response struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_e585294ab3f34af5, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "mypackage.Request")
	proto.RegisterType((*Response)(nil), "mypackage.Response")
}

func init() { proto.RegisterFile("greeter.proto", fileDescriptor_e585294ab3f34af5) }

var fileDescriptor_e585294ab3f34af5 = []byte{
	// 138 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x2f, 0x4a, 0x4d,
	0x2d, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xcc, 0xad, 0x2c, 0x48, 0x4c,
	0xce, 0x4e, 0x4c, 0x4f, 0x55, 0x92, 0xe5, 0x62, 0x0f, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11,
	0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3,
	0x95, 0x64, 0xb8, 0x38, 0x82, 0x52, 0x8b, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x04, 0xb8, 0x98,
	0x73, 0x8b, 0xd3, 0xa1, 0xd2, 0x20, 0xa6, 0x91, 0x2d, 0x17, 0xbb, 0x3b, 0xc4, 0x60, 0x21, 0x23,
	0x2e, 0x56, 0x8f, 0xd4, 0x9c, 0x9c, 0x7c, 0x21, 0x21, 0x3d, 0xb8, 0xe1, 0x7a, 0x50, 0x93, 0xa5,
	0x84, 0x51, 0xc4, 0x20, 0xc6, 0x29, 0x31, 0x24, 0xb1, 0x81, 0x5d, 0x63, 0x0c, 0x08, 0x00, 0x00,
	0xff, 0xff, 0xe4, 0xec, 0xf3, 0x94, 0x9e, 0x00, 0x00, 0x00,
}

// Client API for Greeter service

type GreeterService interface {
	Hello(ctx context.Context, req *Request, opts ...client.CallOption) error
}

type GreeterCallback interface {
	Hello(ctx context.Context, rsp *Response) error
}

type greeterService struct {
	c    client.Client
	name string
}

func NewGreeterService(name string, hdcb GreeterCallback, c client.Client) GreeterService {
	if c == nil {
		panic("client is nil")
	}
	if len(name) == 0 {
		panic("name is nil")
	}
	if err := c.Handle(hdcb); err != nil {
		panic(err)
	}
	return &greeterService{
		c:    c,
		name: name,
	}
}

func (c *greeterService) Hello(ctx context.Context, req *Request, opts ...client.CallOption) error {
	r := c.c.NewRequest(c.name, "Greeter.Hello", req)
	err := c.c.Call(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// Server API for Greeter service

type GreeterHandler interface {
	Hello(ctx context.Context, req *Request, rsp *Response) error
}

func RegisterGreeterHandler(s server.Server, hdlr GreeterHandler) error {
	return s.Handle(hdlr)
}
