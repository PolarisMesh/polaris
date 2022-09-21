// Code generated by protoc-gen-go. DO NOT EDIT.
// source: request_v2.proto

package v2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DiscoverRequest_DiscoverRequestType int32

const (
	DiscoverRequest_UNKNOWN         DiscoverRequest_DiscoverRequestType = 0
	DiscoverRequest_ROUTING         DiscoverRequest_DiscoverRequestType = 1
	DiscoverRequest_CIRCUIT_BREAKER DiscoverRequest_DiscoverRequestType = 2
)

var DiscoverRequest_DiscoverRequestType_name = map[int32]string{
	0: "UNKNOWN",
	1: "ROUTING",
	2: "CIRCUIT_BREAKER",
}
var DiscoverRequest_DiscoverRequestType_value = map[string]int32{
	"UNKNOWN":         0,
	"ROUTING":         1,
	"CIRCUIT_BREAKER": 2,
}

func (x DiscoverRequest_DiscoverRequestType) String() string {
	return proto.EnumName(DiscoverRequest_DiscoverRequestType_name, int32(x))
}
func (DiscoverRequest_DiscoverRequestType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_request_v2_da411c9ba009335d, []int{1, 0}
}

type Service struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Revision             string   `protobuf:"bytes,3,opt,name=revision,proto3" json:"revision,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_request_v2_da411c9ba009335d, []int{0}
}
func (m *Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service.Unmarshal(m, b)
}
func (m *Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service.Marshal(b, m, deterministic)
}
func (dst *Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service.Merge(dst, src)
}
func (m *Service) XXX_Size() int {
	return xxx_messageInfo_Service.Size(m)
}
func (m *Service) XXX_DiscardUnknown() {
	xxx_messageInfo_Service.DiscardUnknown(m)
}

var xxx_messageInfo_Service proto.InternalMessageInfo

func (m *Service) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Service) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Service) GetRevision() string {
	if m != nil {
		return m.Revision
	}
	return ""
}

type DiscoverRequest struct {
	Type                 DiscoverRequest_DiscoverRequestType `protobuf:"varint,1,opt,name=type,proto3,enum=v2.DiscoverRequest_DiscoverRequestType" json:"type,omitempty"`
	Serivce              *Service                            `protobuf:"bytes,2,opt,name=serivce,proto3" json:"serivce,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *DiscoverRequest) Reset()         { *m = DiscoverRequest{} }
func (m *DiscoverRequest) String() string { return proto.CompactTextString(m) }
func (*DiscoverRequest) ProtoMessage()    {}
func (*DiscoverRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_request_v2_da411c9ba009335d, []int{1}
}
func (m *DiscoverRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoverRequest.Unmarshal(m, b)
}
func (m *DiscoverRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoverRequest.Marshal(b, m, deterministic)
}
func (dst *DiscoverRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoverRequest.Merge(dst, src)
}
func (m *DiscoverRequest) XXX_Size() int {
	return xxx_messageInfo_DiscoverRequest.Size(m)
}
func (m *DiscoverRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoverRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoverRequest proto.InternalMessageInfo

func (m *DiscoverRequest) GetType() DiscoverRequest_DiscoverRequestType {
	if m != nil {
		return m.Type
	}
	return DiscoverRequest_UNKNOWN
}

func (m *DiscoverRequest) GetSerivce() *Service {
	if m != nil {
		return m.Serivce
	}
	return nil
}

func init() {
	proto.RegisterType((*Service)(nil), "v2.Service")
	proto.RegisterType((*DiscoverRequest)(nil), "v2.DiscoverRequest")
	proto.RegisterEnum("v2.DiscoverRequest_DiscoverRequestType", DiscoverRequest_DiscoverRequestType_name, DiscoverRequest_DiscoverRequestType_value)
}

func init() { proto.RegisterFile("request_v2.proto", fileDescriptor_request_v2_da411c9ba009335d) }

var fileDescriptor_request_v2_da411c9ba009335d = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x4a, 0x2d, 0x2c,
	0x4d, 0x2d, 0x2e, 0x89, 0x2f, 0x33, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x33,
	0x52, 0x0a, 0xe7, 0x62, 0x0f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0x12, 0xe2, 0x62, 0xc9,
	0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85, 0x64, 0xb8, 0x38,
	0x41, 0x74, 0x71, 0x41, 0x62, 0x72, 0xaa, 0x04, 0x13, 0x58, 0x02, 0x21, 0x20, 0x24, 0xc5, 0xc5,
	0x51, 0x94, 0x5a, 0x96, 0x59, 0x9c, 0x99, 0x9f, 0x27, 0xc1, 0x0c, 0x96, 0x84, 0xf3, 0x95, 0x8e,
	0x33, 0x72, 0xf1, 0xbb, 0x64, 0x16, 0x27, 0xe7, 0x97, 0xa5, 0x16, 0x05, 0x41, 0x6c, 0x16, 0xb2,
	0xe6, 0x62, 0x29, 0xa9, 0x2c, 0x80, 0xd8, 0xc0, 0x67, 0xa4, 0xae, 0x57, 0x66, 0xa4, 0x87, 0xa6,
	0x04, 0x9d, 0x1f, 0x52, 0x59, 0x90, 0x1a, 0x04, 0xd6, 0x24, 0xa4, 0xca, 0xc5, 0x5e, 0x9c, 0x5a,
	0x94, 0x59, 0x06, 0x75, 0x08, 0xb7, 0x11, 0x37, 0x48, 0x3f, 0xd4, 0xf1, 0x41, 0x30, 0x39, 0x25,
	0x2f, 0x2e, 0x61, 0x2c, 0x66, 0x08, 0x71, 0x73, 0xb1, 0x87, 0xfa, 0x79, 0xfb, 0xf9, 0x87, 0xfb,
	0x09, 0x30, 0x80, 0x38, 0x41, 0xfe, 0xa1, 0x21, 0x9e, 0x7e, 0xee, 0x02, 0x8c, 0x42, 0xc2, 0x5c,
	0xfc, 0xce, 0x9e, 0x41, 0xce, 0xa1, 0x9e, 0x21, 0xf1, 0x4e, 0x41, 0xae, 0x8e, 0xde, 0xae, 0x41,
	0x02, 0x4c, 0x4a, 0x2c, 0x1c, 0xec, 0x02, 0xdc, 0x5e, 0x2c, 0x1c, 0xcc, 0x02, 0xac, 0x49, 0x6c,
	0xe0, 0xd0, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x23, 0xae, 0x3a, 0x13, 0x41, 0x01, 0x00,
	0x00,
}
