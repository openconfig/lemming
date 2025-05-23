// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: fault/proto/test/test.proto

package test

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	mi := &file_fault_proto_test_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fault_proto_test_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_fault_proto_test_test_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type PingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	mi := &file_fault_proto_test_test_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fault_proto_test_test_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_fault_proto_test_test_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_fault_proto_test_test_proto protoreflect.FileDescriptor

var file_fault_proto_test_test_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x22, 0x1f, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x20, 0x0a, 0x0c, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x83, 0x01, 0x0a, 0x04,
	0x50, 0x69, 0x6e, 0x67, 0x12, 0x3a, 0x0a, 0x05, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x12, 0x16, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x69, 0x6e,
	0x67, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x3f, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x16, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30,
	0x01, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6f, 0x70, 0x65, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x6c, 0x65, 0x6d, 0x6d, 0x69,
	0x6e, 0x67, 0x2f, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fault_proto_test_test_proto_rawDescOnce sync.Once
	file_fault_proto_test_test_proto_rawDescData = file_fault_proto_test_test_proto_rawDesc
)

func file_fault_proto_test_test_proto_rawDescGZIP() []byte {
	file_fault_proto_test_test_proto_rawDescOnce.Do(func() {
		file_fault_proto_test_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_fault_proto_test_test_proto_rawDescData)
	})
	return file_fault_proto_test_test_proto_rawDescData
}

var file_fault_proto_test_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fault_proto_test_test_proto_goTypes = []any{
	(*PingRequest)(nil),  // 0: test.ping.PingRequest
	(*PingResponse)(nil), // 1: test.ping.PingResponse
}
var file_fault_proto_test_test_proto_depIdxs = []int32{
	0, // 0: test.ping.Ping.Unary:input_type -> test.ping.PingRequest
	0, // 1: test.ping.Ping.Stream:input_type -> test.ping.PingRequest
	1, // 2: test.ping.Ping.Unary:output_type -> test.ping.PingResponse
	1, // 3: test.ping.Ping.Stream:output_type -> test.ping.PingResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fault_proto_test_test_proto_init() }
func file_fault_proto_test_test_proto_init() {
	if File_fault_proto_test_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fault_proto_test_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fault_proto_test_test_proto_goTypes,
		DependencyIndexes: file_fault_proto_test_test_proto_depIdxs,
		MessageInfos:      file_fault_proto_test_test_proto_msgTypes,
	}.Build()
	File_fault_proto_test_test_proto = out.File
	file_fault_proto_test_test_proto_rawDesc = nil
	file_fault_proto_test_test_proto_goTypes = nil
	file_fault_proto_test_test_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PingClient is the client API for Ping service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PingClient interface {
	Unary(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Stream(ctx context.Context, opts ...grpc.CallOption) (Ping_StreamClient, error)
}

type pingClient struct {
	cc grpc.ClientConnInterface
}

func NewPingClient(cc grpc.ClientConnInterface) PingClient {
	return &pingClient{cc}
}

func (c *pingClient) Unary(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/test.ping.Ping/Unary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pingClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Ping_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ping_serviceDesc.Streams[0], "/test.ping.Ping/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &pingStreamClient{stream}
	return x, nil
}

type Ping_StreamClient interface {
	Send(*PingRequest) error
	Recv() (*PingResponse, error)
	grpc.ClientStream
}

type pingStreamClient struct {
	grpc.ClientStream
}

func (x *pingStreamClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pingStreamClient) Recv() (*PingResponse, error) {
	m := new(PingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PingServer is the server API for Ping service.
type PingServer interface {
	Unary(context.Context, *PingRequest) (*PingResponse, error)
	Stream(Ping_StreamServer) error
}

// UnimplementedPingServer can be embedded to have forward compatible implementations.
type UnimplementedPingServer struct {
}

func (*UnimplementedPingServer) Unary(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unary not implemented")
}
func (*UnimplementedPingServer) Stream(Ping_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

func RegisterPingServer(s *grpc.Server, srv PingServer) {
	s.RegisterService(&_Ping_serviceDesc, srv)
}

func _Ping_Unary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServer).Unary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.ping.Ping/Unary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServer).Unary(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ping_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PingServer).Stream(&pingStreamServer{stream})
}

type Ping_StreamServer interface {
	Send(*PingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type pingStreamServer struct {
	grpc.ServerStream
}

func (x *pingStreamServer) Send(m *PingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pingStreamServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Ping_serviceDesc = grpc.ServiceDesc{
	ServiceName: "test.ping.Ping",
	HandlerType: (*PingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unary",
			Handler:    _Ping_Unary_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Ping_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "fault/proto/test/test.proto",
}
