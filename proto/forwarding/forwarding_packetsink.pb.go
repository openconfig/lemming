// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: proto/forwarding/forwarding_packetsink.proto

package forwarding

import (
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

type PacketInjectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PortId        *PortId                `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	ContextId     *ContextId             `protobuf:"bytes,2,opt,name=context_id,json=contextId,proto3" json:"context_id,omitempty"`
	Bytes         []byte                 `protobuf:"bytes,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
	Action        PortAction             `protobuf:"varint,4,opt,name=action,proto3,enum=forwarding.PortAction" json:"action,omitempty"`
	Preprocesses  []*ActionDesc          `protobuf:"bytes,7,rep,name=preprocesses,proto3" json:"preprocesses,omitempty"`
	StartHeader   PacketHeaderId         `protobuf:"varint,10,opt,name=start_header,json=startHeader,proto3,enum=forwarding.PacketHeaderId" json:"start_header,omitempty"`
	ParsedFields  []*PacketFieldBytes    `protobuf:"bytes,11,rep,name=parsed_fields,json=parsedFields,proto3" json:"parsed_fields,omitempty"`
	Debug         bool                   `protobuf:"varint,12,opt,name=debug,proto3" json:"debug,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketInjectRequest) Reset() {
	*x = PacketInjectRequest{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketInjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketInjectRequest) ProtoMessage() {}

func (x *PacketInjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketInjectRequest.ProtoReflect.Descriptor instead.
func (*PacketInjectRequest) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{0}
}

func (x *PacketInjectRequest) GetPortId() *PortId {
	if x != nil {
		return x.PortId
	}
	return nil
}

func (x *PacketInjectRequest) GetContextId() *ContextId {
	if x != nil {
		return x.ContextId
	}
	return nil
}

func (x *PacketInjectRequest) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *PacketInjectRequest) GetAction() PortAction {
	if x != nil {
		return x.Action
	}
	return PortAction_PORT_ACTION_UNSPECIFIED
}

func (x *PacketInjectRequest) GetPreprocesses() []*ActionDesc {
	if x != nil {
		return x.Preprocesses
	}
	return nil
}

func (x *PacketInjectRequest) GetStartHeader() PacketHeaderId {
	if x != nil {
		return x.StartHeader
	}
	return PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED
}

func (x *PacketInjectRequest) GetParsedFields() []*PacketFieldBytes {
	if x != nil {
		return x.ParsedFields
	}
	return nil
}

func (x *PacketInjectRequest) GetDebug() bool {
	if x != nil {
		return x.Debug
	}
	return false
}

type PacketInjectResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketInjectResponse) Reset() {
	*x = PacketInjectResponse{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketInjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketInjectResponse) ProtoMessage() {}

func (x *PacketInjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketInjectResponse.ProtoReflect.Descriptor instead.
func (*PacketInjectResponse) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{1}
}

type PacketSinkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ContextId     *ContextId             `protobuf:"bytes,1,opt,name=context_id,json=contextId,proto3" json:"context_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketSinkRequest) Reset() {
	*x = PacketSinkRequest{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketSinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketSinkRequest) ProtoMessage() {}

func (x *PacketSinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketSinkRequest.ProtoReflect.Descriptor instead.
func (*PacketSinkRequest) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{2}
}

func (x *PacketSinkRequest) GetContextId() *ContextId {
	if x != nil {
		return x.ContextId
	}
	return nil
}

type PacketSinkPacketInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bytes         []byte                 `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
	PortId        *PortId                `protobuf:"bytes,2,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	Ingress       *PortId                `protobuf:"bytes,3,opt,name=ingress,proto3" json:"ingress,omitempty"`
	Egress        *PortId                `protobuf:"bytes,4,opt,name=egress,proto3" json:"egress,omitempty"`
	ParsedFields  []*PacketFieldBytes    `protobuf:"bytes,5,rep,name=parsed_fields,json=parsedFields,proto3" json:"parsed_fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketSinkPacketInfo) Reset() {
	*x = PacketSinkPacketInfo{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketSinkPacketInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketSinkPacketInfo) ProtoMessage() {}

func (x *PacketSinkPacketInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketSinkPacketInfo.ProtoReflect.Descriptor instead.
func (*PacketSinkPacketInfo) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{3}
}

func (x *PacketSinkPacketInfo) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *PacketSinkPacketInfo) GetPortId() *PortId {
	if x != nil {
		return x.PortId
	}
	return nil
}

func (x *PacketSinkPacketInfo) GetIngress() *PortId {
	if x != nil {
		return x.Ingress
	}
	return nil
}

func (x *PacketSinkPacketInfo) GetEgress() *PortId {
	if x != nil {
		return x.Egress
	}
	return nil
}

func (x *PacketSinkPacketInfo) GetParsedFields() []*PacketFieldBytes {
	if x != nil {
		return x.ParsedFields
	}
	return nil
}

type PacketSinkPortInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Port          *PortDesc              `protobuf:"bytes,1,opt,name=port,proto3" json:"port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketSinkPortInfo) Reset() {
	*x = PacketSinkPortInfo{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketSinkPortInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketSinkPortInfo) ProtoMessage() {}

func (x *PacketSinkPortInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketSinkPortInfo.ProtoReflect.Descriptor instead.
func (*PacketSinkPortInfo) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{4}
}

func (x *PacketSinkPortInfo) GetPort() *PortDesc {
	if x != nil {
		return x.Port
	}
	return nil
}

type PacketSinkResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Resp:
	//
	//	*PacketSinkResponse_Packet
	//	*PacketSinkResponse_Port
	Resp          isPacketSinkResponse_Resp `protobuf_oneof:"resp"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PacketSinkResponse) Reset() {
	*x = PacketSinkResponse{}
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PacketSinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketSinkResponse) ProtoMessage() {}

func (x *PacketSinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarding_forwarding_packetsink_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketSinkResponse.ProtoReflect.Descriptor instead.
func (*PacketSinkResponse) Descriptor() ([]byte, []int) {
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP(), []int{5}
}

func (x *PacketSinkResponse) GetResp() isPacketSinkResponse_Resp {
	if x != nil {
		return x.Resp
	}
	return nil
}

func (x *PacketSinkResponse) GetPacket() *PacketSinkPacketInfo {
	if x != nil {
		if x, ok := x.Resp.(*PacketSinkResponse_Packet); ok {
			return x.Packet
		}
	}
	return nil
}

func (x *PacketSinkResponse) GetPort() *PacketSinkPortInfo {
	if x != nil {
		if x, ok := x.Resp.(*PacketSinkResponse_Port); ok {
			return x.Port
		}
	}
	return nil
}

type isPacketSinkResponse_Resp interface {
	isPacketSinkResponse_Resp()
}

type PacketSinkResponse_Packet struct {
	Packet *PacketSinkPacketInfo `protobuf:"bytes,1,opt,name=packet,proto3,oneof"`
}

type PacketSinkResponse_Port struct {
	Port *PacketSinkPortInfo `protobuf:"bytes,2,opt,name=port,proto3,oneof"`
}

func (*PacketSinkResponse_Packet) isPacketSinkResponse_Resp() {}

func (*PacketSinkResponse_Port) isPacketSinkResponse_Resp() {}

var File_proto_forwarding_forwarding_packetsink_proto protoreflect.FileDescriptor

var file_proto_forwarding_forwarding_packetsink_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69,
	0x6e, 0x67, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x1a, 0x28, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x66, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x92, 0x03, 0x0a, 0x13, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x49, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b,
	0x0a, 0x07, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72,
	0x74, 0x49, 0x64, 0x52, 0x06, 0x70, 0x6f, 0x72, 0x74, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0a, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x49, 0x64, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72,
	0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x73, 0x63, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x65, 0x73, 0x12, 0x3d, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x12, 0x41, 0x0a, 0x0d, 0x70, 0x61, 0x72, 0x73, 0x65, 0x64, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x73, 0x65, 0x64, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x22, 0x16, 0x0a, 0x14, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x6e, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x49, 0x0a, 0x11, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x6e,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x66,
	0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x49, 0x64, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x64, 0x22, 0xf6,
	0x01, 0x0a, 0x14, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x2b, 0x0a,
	0x07, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72, 0x74,
	0x49, 0x64, 0x52, 0x06, 0x70, 0x6f, 0x72, 0x74, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x07, 0x69, 0x6e,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x64, 0x52,
	0x07, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2a, 0x0a, 0x06, 0x65, 0x67, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x64, 0x52, 0x06, 0x65, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x41, 0x0a, 0x0d, 0x70, 0x61, 0x72, 0x73, 0x65, 0x64, 0x5f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x66, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x73, 0x65,
	0x64, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x3e, 0x0a, 0x12, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x53, 0x69, 0x6e, 0x6b, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x28, 0x0a,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x66, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x44, 0x65, 0x73,
	0x63, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x8e, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a,
	0x0a, 0x06, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x48, 0x00, 0x52, 0x06, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x34, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x6e, 0x6b,
	0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x42, 0x06, 0x0a, 0x04, 0x72, 0x65, 0x73, 0x70, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x6c, 0x65, 0x6d, 0x6d, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_forwarding_forwarding_packetsink_proto_rawDescOnce sync.Once
	file_proto_forwarding_forwarding_packetsink_proto_rawDescData = file_proto_forwarding_forwarding_packetsink_proto_rawDesc
)

func file_proto_forwarding_forwarding_packetsink_proto_rawDescGZIP() []byte {
	file_proto_forwarding_forwarding_packetsink_proto_rawDescOnce.Do(func() {
		file_proto_forwarding_forwarding_packetsink_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_forwarding_forwarding_packetsink_proto_rawDescData)
	})
	return file_proto_forwarding_forwarding_packetsink_proto_rawDescData
}

var file_proto_forwarding_forwarding_packetsink_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_forwarding_forwarding_packetsink_proto_goTypes = []any{
	(*PacketInjectRequest)(nil),  // 0: forwarding.PacketInjectRequest
	(*PacketInjectResponse)(nil), // 1: forwarding.PacketInjectResponse
	(*PacketSinkRequest)(nil),    // 2: forwarding.PacketSinkRequest
	(*PacketSinkPacketInfo)(nil), // 3: forwarding.PacketSinkPacketInfo
	(*PacketSinkPortInfo)(nil),   // 4: forwarding.PacketSinkPortInfo
	(*PacketSinkResponse)(nil),   // 5: forwarding.PacketSinkResponse
	(*PortId)(nil),               // 6: forwarding.PortId
	(*ContextId)(nil),            // 7: forwarding.ContextId
	(PortAction)(0),              // 8: forwarding.PortAction
	(*ActionDesc)(nil),           // 9: forwarding.ActionDesc
	(PacketHeaderId)(0),          // 10: forwarding.PacketHeaderId
	(*PacketFieldBytes)(nil),     // 11: forwarding.PacketFieldBytes
	(*PortDesc)(nil),             // 12: forwarding.PortDesc
}
var file_proto_forwarding_forwarding_packetsink_proto_depIdxs = []int32{
	6,  // 0: forwarding.PacketInjectRequest.port_id:type_name -> forwarding.PortId
	7,  // 1: forwarding.PacketInjectRequest.context_id:type_name -> forwarding.ContextId
	8,  // 2: forwarding.PacketInjectRequest.action:type_name -> forwarding.PortAction
	9,  // 3: forwarding.PacketInjectRequest.preprocesses:type_name -> forwarding.ActionDesc
	10, // 4: forwarding.PacketInjectRequest.start_header:type_name -> forwarding.PacketHeaderId
	11, // 5: forwarding.PacketInjectRequest.parsed_fields:type_name -> forwarding.PacketFieldBytes
	7,  // 6: forwarding.PacketSinkRequest.context_id:type_name -> forwarding.ContextId
	6,  // 7: forwarding.PacketSinkPacketInfo.port_id:type_name -> forwarding.PortId
	6,  // 8: forwarding.PacketSinkPacketInfo.ingress:type_name -> forwarding.PortId
	6,  // 9: forwarding.PacketSinkPacketInfo.egress:type_name -> forwarding.PortId
	11, // 10: forwarding.PacketSinkPacketInfo.parsed_fields:type_name -> forwarding.PacketFieldBytes
	12, // 11: forwarding.PacketSinkPortInfo.port:type_name -> forwarding.PortDesc
	3,  // 12: forwarding.PacketSinkResponse.packet:type_name -> forwarding.PacketSinkPacketInfo
	4,  // 13: forwarding.PacketSinkResponse.port:type_name -> forwarding.PacketSinkPortInfo
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_forwarding_forwarding_packetsink_proto_init() }
func file_proto_forwarding_forwarding_packetsink_proto_init() {
	if File_proto_forwarding_forwarding_packetsink_proto != nil {
		return
	}
	file_proto_forwarding_forwarding_action_proto_init()
	file_proto_forwarding_forwarding_common_proto_init()
	file_proto_forwarding_forwarding_port_proto_init()
	file_proto_forwarding_forwarding_packetsink_proto_msgTypes[5].OneofWrappers = []any{
		(*PacketSinkResponse_Packet)(nil),
		(*PacketSinkResponse_Port)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_forwarding_forwarding_packetsink_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_forwarding_forwarding_packetsink_proto_goTypes,
		DependencyIndexes: file_proto_forwarding_forwarding_packetsink_proto_depIdxs,
		MessageInfos:      file_proto_forwarding_forwarding_packetsink_proto_msgTypes,
	}.Build()
	File_proto_forwarding_forwarding_packetsink_proto = out.File
	file_proto_forwarding_forwarding_packetsink_proto_rawDesc = nil
	file_proto_forwarding_forwarding_packetsink_proto_goTypes = nil
	file_proto_forwarding_forwarding_packetsink_proto_depIdxs = nil
}
