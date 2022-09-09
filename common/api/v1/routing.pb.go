// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing.proto

package v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Routing struct {
	// 规则所属服务以及命名空间
	Service   *wrappers.StringValue `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Namespace *wrappers.StringValue `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// 每个服务可以配置多条入站或者出站规则
	// 对于每个请求，从上到下依次匹配，若命中则终止
	Inbounds             []*Route              `protobuf:"bytes,3,rep,name=inbounds,proto3" json:"inbounds,omitempty"`
	Outbounds            []*Route              `protobuf:"bytes,4,rep,name=outbounds,proto3" json:"outbounds,omitempty"`
	Ctime                *wrappers.StringValue `protobuf:"bytes,5,opt,name=ctime,proto3" json:"ctime,omitempty"`
	Mtime                *wrappers.StringValue `protobuf:"bytes,6,opt,name=mtime,proto3" json:"mtime,omitempty"`
	Revision             *wrappers.StringValue `protobuf:"bytes,7,opt,name=revision,proto3" json:"revision,omitempty"`
	ServiceToken         *wrappers.StringValue `protobuf:"bytes,8,opt,name=service_token,proto3" json:"service_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Routing) Reset()         { *m = Routing{} }
func (m *Routing) String() string { return proto.CompactTextString(m) }
func (*Routing) ProtoMessage()    {}
func (*Routing) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_562225d32d6ce1a1, []int{0}
}
func (m *Routing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Routing.Unmarshal(m, b)
}
func (m *Routing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Routing.Marshal(b, m, deterministic)
}
func (dst *Routing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Routing.Merge(dst, src)
}
func (m *Routing) XXX_Size() int {
	return xxx_messageInfo_Routing.Size(m)
}
func (m *Routing) XXX_DiscardUnknown() {
	xxx_messageInfo_Routing.DiscardUnknown(m)
}

var xxx_messageInfo_Routing proto.InternalMessageInfo

func (m *Routing) GetService() *wrappers.StringValue {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *Routing) GetNamespace() *wrappers.StringValue {
	if m != nil {
		return m.Namespace
	}
	return nil
}

func (m *Routing) GetInbounds() []*Route {
	if m != nil {
		return m.Inbounds
	}
	return nil
}

func (m *Routing) GetOutbounds() []*Route {
	if m != nil {
		return m.Outbounds
	}
	return nil
}

func (m *Routing) GetCtime() *wrappers.StringValue {
	if m != nil {
		return m.Ctime
	}
	return nil
}

func (m *Routing) GetMtime() *wrappers.StringValue {
	if m != nil {
		return m.Mtime
	}
	return nil
}

func (m *Routing) GetRevision() *wrappers.StringValue {
	if m != nil {
		return m.Revision
	}
	return nil
}

func (m *Routing) GetServiceToken() *wrappers.StringValue {
	if m != nil {
		return m.ServiceToken
	}
	return nil
}

type Route struct {
	// 如果匹配Source规则，按照Destination路由
	// 多个Source之间的关系为或
	Sources              []*Source      `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
	Destinations         []*Destination `protobuf:"bytes,2,rep,name=destinations,proto3" json:"destinations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_562225d32d6ce1a1, []int{1}
}
func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (dst *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(dst, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetSources() []*Source {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *Route) GetDestinations() []*Destination {
	if m != nil {
		return m.Destinations
	}
	return nil
}

type Source struct {
	// 主调方服务以及命名空间
	Service   *wrappers.StringValue `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Namespace *wrappers.StringValue `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// 主调方服务实例标签或者请求标签
	// value支持正则匹配
	Metadata             map[string]*MatchString `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Source) Reset()         { *m = Source{} }
func (m *Source) String() string { return proto.CompactTextString(m) }
func (*Source) ProtoMessage()    {}
func (*Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_562225d32d6ce1a1, []int{2}
}
func (m *Source) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Source.Unmarshal(m, b)
}
func (m *Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Source.Marshal(b, m, deterministic)
}
func (dst *Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Source.Merge(dst, src)
}
func (m *Source) XXX_Size() int {
	return xxx_messageInfo_Source.Size(m)
}
func (m *Source) XXX_DiscardUnknown() {
	xxx_messageInfo_Source.DiscardUnknown(m)
}

var xxx_messageInfo_Source proto.InternalMessageInfo

func (m *Source) GetService() *wrappers.StringValue {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *Source) GetNamespace() *wrappers.StringValue {
	if m != nil {
		return m.Namespace
	}
	return nil
}

func (m *Source) GetMetadata() map[string]*MatchString {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Destination struct {
	// 被调方服务以及命名空间
	Service   *wrappers.StringValue `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Namespace *wrappers.StringValue `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// 被调方服务实例标签
	// value支持正则匹配
	Metadata map[string]*MatchString `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// 根据服务名和服务实例metadata筛选符合条件的服务实例子集
	// 服务实例子集可以设置优先级和权重
	// 优先级：整型，范围[0, 9]，最高优先级为0
	// 权重：整型
	// 先按优先级路由，如果存在高优先级，不会使用低优先级
	// 如果存在优先级相同的子集，再按权重分配
	// 优先级和权重可以都不设置/设置一个/设置两个
	// 如果部分设置优先级，部分没有设置，认为没有设置的优先级最低
	// 如果部分设置权重，部分没有设置，认为没有设置的权重为0
	// 如果全部没有设置权重，认为权重相同
	Priority *wrappers.UInt32Value `protobuf:"bytes,4,opt,name=priority,proto3" json:"priority,omitempty"`
	Weight   *wrappers.UInt32Value `protobuf:"bytes,5,opt,name=weight,proto3" json:"weight,omitempty"`
	// 将请求转发到代理服务
	Transfer *wrappers.StringValue `protobuf:"bytes,6,opt,name=transfer,proto3" json:"transfer,omitempty"`
	// 是否对该set执行隔离，隔离后，不会再分配流量
	Isolate              *wrappers.BoolValue `protobuf:"bytes,7,opt,name=isolate,proto3" json:"isolate,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Destination) Reset()         { *m = Destination{} }
func (m *Destination) String() string { return proto.CompactTextString(m) }
func (*Destination) ProtoMessage()    {}
func (*Destination) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_562225d32d6ce1a1, []int{3}
}
func (m *Destination) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Destination.Unmarshal(m, b)
}
func (m *Destination) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Destination.Marshal(b, m, deterministic)
}
func (dst *Destination) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Destination.Merge(dst, src)
}
func (m *Destination) XXX_Size() int {
	return xxx_messageInfo_Destination.Size(m)
}
func (m *Destination) XXX_DiscardUnknown() {
	xxx_messageInfo_Destination.DiscardUnknown(m)
}

var xxx_messageInfo_Destination proto.InternalMessageInfo

func (m *Destination) GetService() *wrappers.StringValue {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *Destination) GetNamespace() *wrappers.StringValue {
	if m != nil {
		return m.Namespace
	}
	return nil
}

func (m *Destination) GetMetadata() map[string]*MatchString {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Destination) GetPriority() *wrappers.UInt32Value {
	if m != nil {
		return m.Priority
	}
	return nil
}

func (m *Destination) GetWeight() *wrappers.UInt32Value {
	if m != nil {
		return m.Weight
	}
	return nil
}

func (m *Destination) GetTransfer() *wrappers.StringValue {
	if m != nil {
		return m.Transfer
	}
	return nil
}

func (m *Destination) GetIsolate() *wrappers.BoolValue {
	if m != nil {
		return m.Isolate
	}
	return nil
}

func init() {
	proto.RegisterType((*Routing)(nil), "v1.Routing")
	proto.RegisterType((*Route)(nil), "v1.Route")
	proto.RegisterType((*Source)(nil), "v1.Source")
	proto.RegisterMapType((map[string]*MatchString)(nil), "v1.Source.MetadataEntry")
	proto.RegisterType((*Destination)(nil), "v1.Destination")
	proto.RegisterMapType((map[string]*MatchString)(nil), "v1.Destination.MetadataEntry")
}

func init() { proto.RegisterFile("routing.proto", fileDescriptor_routing_562225d32d6ce1a1) }

var fileDescriptor_routing_562225d32d6ce1a1 = []byte{
	// 468 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x93, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0x86, 0xb1, 0x15, 0x5b, 0xf6, 0xb8, 0xa6, 0x65, 0x4f, 0x8b, 0x69, 0x4b, 0x30, 0x0d, 0xcd,
	0x49, 0x21, 0xb6, 0x29, 0x69, 0x8e, 0xa1, 0x3d, 0x14, 0x9a, 0x8b, 0x42, 0x7b, 0x2d, 0x6b, 0x79,
	0xa2, 0x2c, 0x91, 0x76, 0xc5, 0xee, 0x48, 0xc1, 0xb7, 0xbe, 0x53, 0xdf, 0xaa, 0x4f, 0x51, 0xb4,
	0x5a, 0xdb, 0x71, 0x82, 0x41, 0x87, 0x42, 0x6e, 0xd2, 0xce, 0xf7, 0xcf, 0xcc, 0xfe, 0x3b, 0x03,
	0x63, 0xa3, 0x4b, 0x92, 0x2a, 0x8d, 0x0a, 0xa3, 0x49, 0xb3, 0x6e, 0x75, 0x3e, 0x79, 0x9f, 0x6a,
	0x9d, 0x66, 0x78, 0xe6, 0x4e, 0x96, 0xe5, 0xed, 0xd9, 0x83, 0x11, 0x45, 0x81, 0xc6, 0x36, 0xcc,
	0x64, 0x94, 0xeb, 0x15, 0x66, 0xcd, 0xcf, 0xf4, 0x4f, 0x00, 0x61, 0xdc, 0xa4, 0x60, 0x9f, 0x20,
	0xb4, 0x68, 0x2a, 0x99, 0x20, 0xef, 0x1c, 0x77, 0x4e, 0x47, 0xb3, 0xb7, 0x51, 0x93, 0x2a, 0xda,
	0xa4, 0x8a, 0x6e, 0xc8, 0x48, 0x95, 0xfe, 0x14, 0x59, 0x89, 0xf1, 0x06, 0x66, 0x97, 0x30, 0x54,
	0x22, 0x47, 0x5b, 0x88, 0x04, 0x79, 0xb7, 0x85, 0x72, 0x87, 0xb3, 0x13, 0x18, 0x48, 0xb5, 0xd4,
	0xa5, 0x5a, 0x59, 0x1e, 0x1c, 0x07, 0xa7, 0xa3, 0xd9, 0x30, 0xaa, 0xce, 0xa3, 0xba, 0x25, 0x8c,
	0xb7, 0x21, 0xf6, 0x11, 0x86, 0xba, 0x24, 0xcf, 0x1d, 0x3d, 0xe5, 0x76, 0x31, 0x36, 0x83, 0x5e,
	0x42, 0x32, 0x47, 0xde, 0x6b, 0xd1, 0x47, 0x83, 0xd6, 0x9a, 0xdc, 0x69, 0xfa, 0x6d, 0x34, 0x0e,
	0x65, 0x17, 0x30, 0x30, 0x58, 0x49, 0x2b, 0xb5, 0xe2, 0x61, 0x0b, 0xd9, 0x96, 0x66, 0x57, 0x30,
	0xf6, 0xc6, 0xfd, 0x22, 0x7d, 0x8f, 0x8a, 0x0f, 0x5a, 0xc8, 0xf7, 0x25, 0xd3, 0x25, 0xf4, 0xdc,
	0xcd, 0xd9, 0x07, 0x08, 0xad, 0x2e, 0x4d, 0x82, 0x96, 0x77, 0x9c, 0x2b, 0x50, 0xbb, 0x72, 0xe3,
	0x8e, 0xe2, 0x4d, 0x88, 0xcd, 0xe1, 0xd5, 0x0a, 0x2d, 0x49, 0x25, 0x48, 0x6a, 0x65, 0x79, 0xd7,
	0xa1, 0xaf, 0x6b, 0xf4, 0xcb, 0xee, 0x3c, 0xde, 0x83, 0xa6, 0xbf, 0xbb, 0xd0, 0x6f, 0x12, 0xbd,
	0xc8, 0x60, 0x2c, 0x60, 0x90, 0x23, 0x89, 0x95, 0x20, 0xe1, 0x07, 0x83, 0xef, 0xae, 0x16, 0x5d,
	0xfb, 0xd0, 0x57, 0x45, 0x66, 0x1d, 0x6f, 0xc9, 0xc9, 0x77, 0x18, 0xef, 0x85, 0xd8, 0x1b, 0x08,
	0xee, 0x71, 0xed, 0xda, 0x1e, 0xc6, 0xf5, 0x27, 0x3b, 0x81, 0x5e, 0x55, 0x17, 0xf3, 0x0d, 0x39,
	0x17, 0xae, 0x05, 0x25, 0x77, 0x4d, 0x23, 0x71, 0x13, 0xbd, 0xec, 0x5e, 0x74, 0xa6, 0x7f, 0x03,
	0x18, 0x3d, 0x32, 0xe8, 0x45, 0x7c, 0xf8, 0xfc, 0xcc, 0x87, 0x77, 0x4f, 0xde, 0xed, 0x90, 0x19,
	0xf5, 0x8c, 0x16, 0x46, 0x6a, 0x23, 0x69, 0xcd, 0x8f, 0x0e, 0x54, 0xfd, 0xf1, 0x4d, 0xd1, 0x7c,
	0xe6, 0x67, 0x74, 0x43, 0xb3, 0x05, 0xf4, 0x1f, 0x50, 0xa6, 0x77, 0x74, 0x70, 0x8d, 0x1e, 0xeb,
	0x3c, 0x5b, 0xd7, 0x23, 0x23, 0x94, 0xbd, 0x45, 0xd3, 0x6a, 0x95, 0xb6, 0x34, 0x5b, 0x40, 0x28,
	0xad, 0xce, 0x04, 0xa1, 0x5f, 0xa6, 0xc9, 0x33, 0xe1, 0x95, 0xd6, 0x99, 0xb7, 0xd5, 0xa3, 0xff,
	0xf7, 0xb1, 0x97, 0x7d, 0x57, 0x6a, 0xfe, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xb3, 0x0a, 0x3d, 0x7c,
	0x52, 0x05, 0x00, 0x00,
}
