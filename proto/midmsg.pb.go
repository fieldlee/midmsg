// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/midmsg.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type NetReqInfo struct {
	M_Body               []byte   `protobuf:"bytes,1,opt,name=m_Body,json=mBody,proto3" json:"m_Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetReqInfo) Reset()         { *m = NetReqInfo{} }
func (m *NetReqInfo) String() string { return proto.CompactTextString(m) }
func (*NetReqInfo) ProtoMessage()    {}
func (*NetReqInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_7304c339219e1145, []int{0}
}

func (m *NetReqInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetReqInfo.Unmarshal(m, b)
}
func (m *NetReqInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetReqInfo.Marshal(b, m, deterministic)
}
func (m *NetReqInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetReqInfo.Merge(m, src)
}
func (m *NetReqInfo) XXX_Size() int {
	return xxx_messageInfo_NetReqInfo.Size(m)
}
func (m *NetReqInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NetReqInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NetReqInfo proto.InternalMessageInfo

func (m *NetReqInfo) GetM_Body() []byte {
	if m != nil {
		return m.M_Body
	}
	return nil
}

type NetRspInfo struct {
	M_Resp               []byte   `protobuf:"bytes,1,opt,name=m_Resp,json=mResp,proto3" json:"m_Resp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetRspInfo) Reset()         { *m = NetRspInfo{} }
func (m *NetRspInfo) String() string { return proto.CompactTextString(m) }
func (*NetRspInfo) ProtoMessage()    {}
func (*NetRspInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_7304c339219e1145, []int{1}
}

func (m *NetRspInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetRspInfo.Unmarshal(m, b)
}
func (m *NetRspInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetRspInfo.Marshal(b, m, deterministic)
}
func (m *NetRspInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetRspInfo.Merge(m, src)
}
func (m *NetRspInfo) XXX_Size() int {
	return xxx_messageInfo_NetRspInfo.Size(m)
}
func (m *NetRspInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NetRspInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NetRspInfo proto.InternalMessageInfo

func (m *NetRspInfo) GetM_Resp() []byte {
	if m != nil {
		return m.M_Resp
	}
	return nil
}

func init() {
	proto.RegisterType((*NetReqInfo)(nil), "proto.NetReqInfo")
	proto.RegisterType((*NetRspInfo)(nil), "proto.NetRspInfo")
}

func init() { proto.RegisterFile("proto/midmsg.proto", fileDescriptor_7304c339219e1145) }

var fileDescriptor_7304c339219e1145 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0xcd, 0x4c, 0xc9, 0x2d, 0x4e, 0xd7, 0x03, 0x73, 0x84, 0x58, 0xc1, 0x94, 0x92,
	0x32, 0x17, 0x97, 0x5f, 0x6a, 0x49, 0x50, 0x6a, 0xa1, 0x67, 0x5e, 0x5a, 0xbe, 0x90, 0x28, 0x17,
	0x5b, 0x6e, 0xbc, 0x53, 0x7e, 0x4a, 0xa5, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x6b, 0x2e,
	0x88, 0x03, 0x53, 0x54, 0x5c, 0x80, 0x50, 0x14, 0x94, 0x5a, 0x5c, 0x00, 0x57, 0x04, 0xe2, 0x18,
	0xe5, 0x72, 0x71, 0xf9, 0x66, 0xa6, 0x04, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x0a, 0xe9, 0x71,
	0xb1, 0x04, 0x57, 0xe6, 0x25, 0x0b, 0x09, 0x42, 0xac, 0xd3, 0x43, 0x58, 0x22, 0x85, 0x2c, 0x04,
	0x31, 0x52, 0x89, 0x41, 0x48, 0x9f, 0x8b, 0xd5, 0xb1, 0x98, 0x04, 0x0d, 0x46, 0xf6, 0x5c, 0xbc,
	0xce, 0x39, 0x99, 0xa9, 0x79, 0x25, 0x48, 0x36, 0x3a, 0x27, 0xe6, 0xe4, 0x10, 0x6b, 0x40, 0x12,
	0x1b, 0x58, 0xcc, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x36, 0xb5, 0x1a, 0xc6, 0x1d, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MidServiceClient is the client API for MidService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MidServiceClient interface {
	Sync(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error)
	Async(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error)
}

type midServiceClient struct {
	cc *grpc.ClientConn
}

func NewMidServiceClient(cc *grpc.ClientConn) MidServiceClient {
	return &midServiceClient{cc}
}

func (c *midServiceClient) Sync(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error) {
	out := new(NetRspInfo)
	err := c.cc.Invoke(ctx, "/proto.MidService/Sync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *midServiceClient) Async(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error) {
	out := new(NetRspInfo)
	err := c.cc.Invoke(ctx, "/proto.MidService/Async", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MidServiceServer is the server API for MidService service.
type MidServiceServer interface {
	Sync(context.Context, *NetReqInfo) (*NetRspInfo, error)
	Async(context.Context, *NetReqInfo) (*NetRspInfo, error)
}

// UnimplementedMidServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMidServiceServer struct {
}

func (*UnimplementedMidServiceServer) Sync(ctx context.Context, req *NetReqInfo) (*NetRspInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (*UnimplementedMidServiceServer) Async(ctx context.Context, req *NetReqInfo) (*NetRspInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Async not implemented")
}

func RegisterMidServiceServer(s *grpc.Server, srv MidServiceServer) {
	s.RegisterService(&_MidService_serviceDesc, srv)
}

func _MidService_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetReqInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MidServiceServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MidService/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MidServiceServer).Sync(ctx, req.(*NetReqInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _MidService_Async_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetReqInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MidServiceServer).Async(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MidService/Async",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MidServiceServer).Async(ctx, req.(*NetReqInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _MidService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MidService",
	HandlerType: (*MidServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _MidService_Sync_Handler,
		},
		{
			MethodName: "Async",
			Handler:    _MidService_Async_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/midmsg.proto",
}

// ClientServiceClient is the client API for ClientService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClientServiceClient interface {
	Call(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error)
}

type clientServiceClient struct {
	cc *grpc.ClientConn
}

func NewClientServiceClient(cc *grpc.ClientConn) ClientServiceClient {
	return &clientServiceClient{cc}
}

func (c *clientServiceClient) Call(ctx context.Context, in *NetReqInfo, opts ...grpc.CallOption) (*NetRspInfo, error) {
	out := new(NetRspInfo)
	err := c.cc.Invoke(ctx, "/proto.ClientService/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServiceServer is the server API for ClientService service.
type ClientServiceServer interface {
	Call(context.Context, *NetReqInfo) (*NetRspInfo, error)
}

// UnimplementedClientServiceServer can be embedded to have forward compatible implementations.
type UnimplementedClientServiceServer struct {
}

func (*UnimplementedClientServiceServer) Call(ctx context.Context, req *NetReqInfo) (*NetRspInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}

func RegisterClientServiceServer(s *grpc.Server, srv ClientServiceServer) {
	s.RegisterService(&_ClientService_serviceDesc, srv)
}

func _ClientService_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetReqInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServiceServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClientService/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServiceServer).Call(ctx, req.(*NetReqInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClientService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClientService",
	HandlerType: (*ClientServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _ClientService_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/midmsg.proto",
}