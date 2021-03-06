// Code generated by protoc-gen-go.
// source: dcache.proto
// DO NOT EDIT!

/*
Package dcache is a generated protocol buffer package.

It is generated from these files:
	dcache.proto

It has these top-level messages:
	GetRequest
	GetResponse
	SetRequest
	SetResponse
	DelRequest
	DelResponse
*/
package dcache

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetRequest struct {
	Group string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *GetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetResponse struct {
	Status    bool    `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Value     []byte  `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	MinuteQps float64 `protobuf:"fixed64,3,opt,name=minute_qps,json=minuteQps" json:"minute_qps,omitempty"`
}

func (m *GetResponse) Reset()                    { *m = GetResponse{} }
func (m *GetResponse) String() string            { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()               {}
func (*GetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *GetResponse) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *GetResponse) GetMinuteQps() float64 {
	if m != nil {
		return m.MinuteQps
	}
	return 0
}

type SetRequest struct {
	Group string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *SetRequest) Reset()                    { *m = SetRequest{} }
func (m *SetRequest) String() string            { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()               {}
func (*SetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SetRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *SetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SetRequest) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type SetResponse struct {
	Status    bool    `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	MinuteQps float64 `protobuf:"fixed64,2,opt,name=minute_qps,json=minuteQps" json:"minute_qps,omitempty"`
}

func (m *SetResponse) Reset()                    { *m = SetResponse{} }
func (m *SetResponse) String() string            { return proto.CompactTextString(m) }
func (*SetResponse) ProtoMessage()               {}
func (*SetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SetResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *SetResponse) GetMinuteQps() float64 {
	if m != nil {
		return m.MinuteQps
	}
	return 0
}

type DelRequest struct {
	Group string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
}

func (m *DelRequest) Reset()                    { *m = DelRequest{} }
func (m *DelRequest) String() string            { return proto.CompactTextString(m) }
func (*DelRequest) ProtoMessage()               {}
func (*DelRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DelRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *DelRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type DelResponse struct {
	Status    bool    `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Value     []byte  `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	MinuteQps float64 `protobuf:"fixed64,3,opt,name=minute_qps,json=minuteQps" json:"minute_qps,omitempty"`
}

func (m *DelResponse) Reset()                    { *m = DelResponse{} }
func (m *DelResponse) String() string            { return proto.CompactTextString(m) }
func (*DelResponse) ProtoMessage()               {}
func (*DelResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DelResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *DelResponse) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *DelResponse) GetMinuteQps() float64 {
	if m != nil {
		return m.MinuteQps
	}
	return 0
}

func init() {
	proto.RegisterType((*GetRequest)(nil), "dcache.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "dcache.GetResponse")
	proto.RegisterType((*SetRequest)(nil), "dcache.SetRequest")
	proto.RegisterType((*SetResponse)(nil), "dcache.SetResponse")
	proto.RegisterType((*DelRequest)(nil), "dcache.DelRequest")
	proto.RegisterType((*DelResponse)(nil), "dcache.DelResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GroupCache service

type GroupCacheClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Del(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*DelResponse, error)
}

type groupCacheClient struct {
	cc *grpc.ClientConn
}

func NewGroupCacheClient(cc *grpc.ClientConn) GroupCacheClient {
	return &groupCacheClient{cc}
}

func (c *groupCacheClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := grpc.Invoke(ctx, "/dcache.GroupCache/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupCacheClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := grpc.Invoke(ctx, "/dcache.GroupCache/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupCacheClient) Del(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*DelResponse, error) {
	out := new(DelResponse)
	err := grpc.Invoke(ctx, "/dcache.GroupCache/Del", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GroupCache service

type GroupCacheServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Del(context.Context, *DelRequest) (*DelResponse, error)
}

func RegisterGroupCacheServer(s *grpc.Server, srv GroupCacheServer) {
	s.RegisterService(&_GroupCache_serviceDesc, srv)
}

func _GroupCache_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupCacheServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcache.GroupCache/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupCacheServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupCache_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupCacheServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcache.GroupCache/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupCacheServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupCache_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupCacheServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcache.GroupCache/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupCacheServer).Del(ctx, req.(*DelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GroupCache_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dcache.GroupCache",
	HandlerType: (*GroupCacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _GroupCache_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _GroupCache_Set_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _GroupCache_Del_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dcache.proto",
}

func init() { proto.RegisterFile("dcache.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x92, 0xc1, 0x4a, 0x03, 0x31,
	0x18, 0x84, 0x9b, 0x2e, 0x2e, 0x76, 0xda, 0x83, 0x44, 0x91, 0x45, 0x10, 0x4a, 0x4e, 0x3d, 0x15,
	0xd1, 0xbe, 0x81, 0x0b, 0x0b, 0xde, 0x4c, 0x6e, 0x5e, 0x64, 0xad, 0x3f, 0x2a, 0xae, 0x4d, 0xda,
	0x24, 0x82, 0xaf, 0xe3, 0x93, 0x4a, 0x92, 0x42, 0xb6, 0xc5, 0x83, 0x0a, 0xde, 0x32, 0x3f, 0x7c,
	0xf9, 0x27, 0x33, 0xc1, 0xe4, 0x71, 0xd9, 0x2e, 0x9f, 0x69, 0x6e, 0x36, 0xda, 0x69, 0x5e, 0x26,
	0x25, 0x16, 0x40, 0x43, 0x4e, 0xd2, 0xda, 0x93, 0x75, 0xfc, 0x04, 0x07, 0x4f, 0x1b, 0xed, 0x4d,
	0xc5, 0xa6, 0x6c, 0x36, 0x92, 0x49, 0xf0, 0x23, 0x14, 0xaf, 0xf4, 0x51, 0x0d, 0xe3, 0x2c, 0x1c,
	0xc5, 0x1d, 0xc6, 0x91, 0xb2, 0x46, 0xaf, 0x2c, 0xf1, 0x53, 0x94, 0xd6, 0xb5, 0xce, 0xdb, 0xc8,
	0x1d, 0xca, 0xad, 0x0a, 0xd7, 0xbd, 0xb7, 0x9d, 0xa7, 0x88, 0x4e, 0x64, 0x12, 0xfc, 0x1c, 0x78,
	0x7b, 0x59, 0x79, 0x47, 0xf7, 0x6b, 0x63, 0xab, 0x62, 0xca, 0x66, 0x4c, 0x8e, 0xd2, 0xe4, 0xd6,
	0x58, 0x71, 0x03, 0xa8, 0x5f, 0x3b, 0xca, 0xab, 0x8a, 0xde, 0x2a, 0x51, 0x63, 0xac, 0x7e, 0xe0,
	0x73, 0xd7, 0xd1, 0x70, 0xdf, 0xd1, 0x02, 0xa8, 0xa9, 0xfb, 0x43, 0x46, 0x91, 0xfa, 0x87, 0x8c,
	0x2e, 0x3f, 0x19, 0xd0, 0x84, 0xbd, 0xd7, 0xa1, 0x44, 0x7e, 0x81, 0xa2, 0x21, 0xc7, 0xf9, 0x7c,
	0x5b, 0x71, 0x6e, 0xf4, 0xec, 0x78, 0x67, 0x96, 0xbc, 0x88, 0x41, 0x20, 0x54, 0x9f, 0x50, 0xdf,
	0x10, 0x6a, 0x9f, 0xa8, 0xa9, 0xcb, 0x44, 0x4e, 0x24, 0x13, 0xbd, 0xf7, 0x8a, 0xc1, 0x43, 0x19,
	0x7f, 0xda, 0xd5, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0xb5, 0xa1, 0x24, 0x79, 0x02, 0x00,
	0x00,
}
