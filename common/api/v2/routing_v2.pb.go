// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing_v2.proto

package v2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RoutingPolicy int32

const (
	// Route by rule rule => RuleRoutingConfig
	RoutingPolicy_RulePolicy RoutingPolicy = 0
	// Route by destination metadata
	RoutingPolicy_MetadataPolicy RoutingPolicy = 1
)

var RoutingPolicy_name = map[int32]string{
	0: "RulePolicy",
	1: "MetadataPolicy",
}
var RoutingPolicy_value = map[string]int32{
	"RulePolicy":     0,
	"MetadataPolicy": 1,
}

func (x RoutingPolicy) String() string {
	return proto.EnumName(RoutingPolicy_name, int32(x))
}
func (RoutingPolicy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{0}
}

type MetadataFailover_FailoverRange int32

const (
	// ALL return all instances
	MetadataFailover_ALL MetadataFailover_FailoverRange = 0
	// OTHERS retuen without thie labels instances
	MetadataFailover_OTHERS MetadataFailover_FailoverRange = 1
	// OTHER_KEYS return other instances which match keys
	MetadataFailover_OTHER_KEYS MetadataFailover_FailoverRange = 2
)

var MetadataFailover_FailoverRange_name = map[int32]string{
	0: "ALL",
	1: "OTHERS",
	2: "OTHER_KEYS",
}
var MetadataFailover_FailoverRange_value = map[string]int32{
	"ALL":        0,
	"OTHERS":     1,
	"OTHER_KEYS": 2,
}

func (x MetadataFailover_FailoverRange) String() string {
	return proto.EnumName(MetadataFailover_FailoverRange_name, int32(x))
}
func (MetadataFailover_FailoverRange) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{4, 0}
}

// label type for gateway request
type SourceMatch_Type int32

const (
	// custom arguments
	SourceMatch_CUSTOM SourceMatch_Type = 0
	// method, match the http post/get/put/delete or grpc method
	SourceMatch_METHOD SourceMatch_Type = 1
	// header, match the http header, dubbo attachment, grpc header
	SourceMatch_HEADER SourceMatch_Type = 2
	// query, match the http query, dubbo argument
	SourceMatch_QUERY SourceMatch_Type = 3
	// caller host ip
	SourceMatch_CALLER_IP SourceMatch_Type = 4
	// path, math the http url
	SourceMatch_PATH SourceMatch_Type = 5
)

var SourceMatch_Type_name = map[int32]string{
	0: "CUSTOM",
	1: "METHOD",
	2: "HEADER",
	3: "QUERY",
	4: "CALLER_IP",
	5: "PATH",
}
var SourceMatch_Type_value = map[string]int32{
	"CUSTOM":    0,
	"METHOD":    1,
	"HEADER":    2,
	"QUERY":     3,
	"CALLER_IP": 4,
	"PATH":      5,
}

func (x SourceMatch_Type) String() string {
	return proto.EnumName(SourceMatch_Type_name, int32(x))
}
func (SourceMatch_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{9, 0}
}

// FlowStaining 流量染色
type FlowStaining struct {
	Id string `protobuf:"bytes,10,opt,name=id,proto3" json:"id,omitempty"`
	// flow statining rule name
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// rules staining rules
	Rules                []*StaineRule `protobuf:"bytes,3,rep,name=rules,proto3" json:"rules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *FlowStaining) Reset()         { *m = FlowStaining{} }
func (m *FlowStaining) String() string { return proto.CompactTextString(m) }
func (*FlowStaining) ProtoMessage()    {}
func (*FlowStaining) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{0}
}
func (m *FlowStaining) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FlowStaining.Unmarshal(m, b)
}
func (m *FlowStaining) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FlowStaining.Marshal(b, m, deterministic)
}
func (dst *FlowStaining) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FlowStaining.Merge(dst, src)
}
func (m *FlowStaining) XXX_Size() int {
	return xxx_messageInfo_FlowStaining.Size(m)
}
func (m *FlowStaining) XXX_DiscardUnknown() {
	xxx_messageInfo_FlowStaining.DiscardUnknown(m)
}

var xxx_messageInfo_FlowStaining proto.InternalMessageInfo

func (m *FlowStaining) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FlowStaining) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FlowStaining) GetRules() []*StaineRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

type StaineRule struct {
	// Traffic matching rules
	Arguments []*SourceMatch `protobuf:"bytes,1,rep,name=arguments,proto3" json:"arguments,omitempty"`
	// Staining label
	Labels []*StaineLabel `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty"`
	// Stain Label
	Priority uint32 `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`
	// rule is enabled
	Enable bool `protobuf:"varint,4,opt,name=enable,proto3" json:"enable,omitempty"`
	// Set the percentage of traffic that needs to be dyed
	StainePercent        uint32   `protobuf:"varint,5,opt,name=stainePercent,proto3" json:"stainePercent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StaineRule) Reset()         { *m = StaineRule{} }
func (m *StaineRule) String() string { return proto.CompactTextString(m) }
func (*StaineRule) ProtoMessage()    {}
func (*StaineRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{1}
}
func (m *StaineRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StaineRule.Unmarshal(m, b)
}
func (m *StaineRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StaineRule.Marshal(b, m, deterministic)
}
func (dst *StaineRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StaineRule.Merge(dst, src)
}
func (m *StaineRule) XXX_Size() int {
	return xxx_messageInfo_StaineRule.Size(m)
}
func (m *StaineRule) XXX_DiscardUnknown() {
	xxx_messageInfo_StaineRule.DiscardUnknown(m)
}

var xxx_messageInfo_StaineRule proto.InternalMessageInfo

func (m *StaineRule) GetArguments() []*SourceMatch {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *StaineRule) GetLabels() []*StaineLabel {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *StaineRule) GetPriority() uint32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *StaineRule) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func (m *StaineRule) GetStainePercent() uint32 {
	if m != nil {
		return m.StainePercent
	}
	return 0
}

type StaineLabel struct {
	Key                  string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StaineLabel) Reset()         { *m = StaineLabel{} }
func (m *StaineLabel) String() string { return proto.CompactTextString(m) }
func (*StaineLabel) ProtoMessage()    {}
func (*StaineLabel) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{2}
}
func (m *StaineLabel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StaineLabel.Unmarshal(m, b)
}
func (m *StaineLabel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StaineLabel.Marshal(b, m, deterministic)
}
func (dst *StaineLabel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StaineLabel.Merge(dst, src)
}
func (m *StaineLabel) XXX_Size() int {
	return xxx_messageInfo_StaineLabel.Size(m)
}
func (m *StaineLabel) XXX_DiscardUnknown() {
	xxx_messageInfo_StaineLabel.DiscardUnknown(m)
}

var xxx_messageInfo_StaineLabel proto.InternalMessageInfo

func (m *StaineLabel) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *StaineLabel) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// configuration root for route
type Routing struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// route rule name
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// namespace namingspace of routing rules
	Namespace string `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Enable this router
	Enable bool `protobuf:"varint,4,opt,name=enable,proto3" json:"enable,omitempty"`
	// Router type
	RoutingPolicy RoutingPolicy `protobuf:"varint,5,opt,name=routing_policy,proto3,enum=v2.RoutingPolicy" json:"routing_policy,omitempty"`
	// Routing configuration for router
	RoutingConfig *any.Any `protobuf:"bytes,6,opt,name=routing_config,proto3" json:"routing_config,omitempty"`
	// revision routing version
	Revision string `protobuf:"bytes,7,opt,name=revision,proto3" json:"revision,omitempty"`
	// ctime create time of the rules
	Ctime string `protobuf:"bytes,8,opt,name=ctime,proto3" json:"ctime,omitempty"`
	// mtime modify time of the rules
	Mtime string `protobuf:"bytes,9,opt,name=mtime,proto3" json:"mtime,omitempty"`
	// etime enable time of the rules
	Etime string `protobuf:"bytes,10,opt,name=etime,proto3" json:"etime,omitempty"`
	// priority rules priority
	Priority             uint32   `protobuf:"varint,11,opt,name=priority,proto3" json:"priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Routing) Reset()         { *m = Routing{} }
func (m *Routing) String() string { return proto.CompactTextString(m) }
func (*Routing) ProtoMessage()    {}
func (*Routing) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{3}
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

func (m *Routing) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Routing) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Routing) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Routing) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func (m *Routing) GetRoutingPolicy() RoutingPolicy {
	if m != nil {
		return m.RoutingPolicy
	}
	return RoutingPolicy_RulePolicy
}

func (m *Routing) GetRoutingConfig() *any.Any {
	if m != nil {
		return m.RoutingConfig
	}
	return nil
}

func (m *Routing) GetRevision() string {
	if m != nil {
		return m.Revision
	}
	return ""
}

func (m *Routing) GetCtime() string {
	if m != nil {
		return m.Ctime
	}
	return ""
}

func (m *Routing) GetMtime() string {
	if m != nil {
		return m.Mtime
	}
	return ""
}

func (m *Routing) GetEtime() string {
	if m != nil {
		return m.Etime
	}
	return ""
}

func (m *Routing) GetPriority() uint32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

type MetadataFailover struct {
	// failover_range metadata route bottom type
	FailoverRange MetadataFailover_FailoverRange `protobuf:"varint,1,opt,name=failover_range,json=failoverRange,proto3,enum=v2.MetadataFailover_FailoverRange" json:"failover_range,omitempty"`
	// only use to failover_range == OTHER_KEYS
	Labels               map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MetadataFailover) Reset()         { *m = MetadataFailover{} }
func (m *MetadataFailover) String() string { return proto.CompactTextString(m) }
func (*MetadataFailover) ProtoMessage()    {}
func (*MetadataFailover) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{4}
}
func (m *MetadataFailover) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataFailover.Unmarshal(m, b)
}
func (m *MetadataFailover) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataFailover.Marshal(b, m, deterministic)
}
func (dst *MetadataFailover) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataFailover.Merge(dst, src)
}
func (m *MetadataFailover) XXX_Size() int {
	return xxx_messageInfo_MetadataFailover.Size(m)
}
func (m *MetadataFailover) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataFailover.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataFailover proto.InternalMessageInfo

func (m *MetadataFailover) GetFailoverRange() MetadataFailover_FailoverRange {
	if m != nil {
		return m.FailoverRange
	}
	return MetadataFailover_ALL
}

func (m *MetadataFailover) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

// MetadataRoutingConfig metadata routing configuration
type MetadataRoutingConfig struct {
	// service
	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// namespace
	Namespace string            `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Labels    map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// When metadata not found, it will fall back to the
	Failover             *MetadataFailover `protobuf:"bytes,4,opt,name=failover,proto3" json:"failover,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MetadataRoutingConfig) Reset()         { *m = MetadataRoutingConfig{} }
func (m *MetadataRoutingConfig) String() string { return proto.CompactTextString(m) }
func (*MetadataRoutingConfig) ProtoMessage()    {}
func (*MetadataRoutingConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{5}
}
func (m *MetadataRoutingConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataRoutingConfig.Unmarshal(m, b)
}
func (m *MetadataRoutingConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataRoutingConfig.Marshal(b, m, deterministic)
}
func (dst *MetadataRoutingConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataRoutingConfig.Merge(dst, src)
}
func (m *MetadataRoutingConfig) XXX_Size() int {
	return xxx_messageInfo_MetadataRoutingConfig.Size(m)
}
func (m *MetadataRoutingConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataRoutingConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataRoutingConfig proto.InternalMessageInfo

func (m *MetadataRoutingConfig) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *MetadataRoutingConfig) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *MetadataRoutingConfig) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *MetadataRoutingConfig) GetFailover() *MetadataFailover {
	if m != nil {
		return m.Failover
	}
	return nil
}

// RuleRoutingConfig routing configuration
type RuleRoutingConfig struct {
	// source source info
	Sources []*Source `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
	// destination destinations info
	Destinations         []*Destination `protobuf:"bytes,2,rep,name=destinations,proto3" json:"destinations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RuleRoutingConfig) Reset()         { *m = RuleRoutingConfig{} }
func (m *RuleRoutingConfig) String() string { return proto.CompactTextString(m) }
func (*RuleRoutingConfig) ProtoMessage()    {}
func (*RuleRoutingConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{6}
}
func (m *RuleRoutingConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuleRoutingConfig.Unmarshal(m, b)
}
func (m *RuleRoutingConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuleRoutingConfig.Marshal(b, m, deterministic)
}
func (dst *RuleRoutingConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuleRoutingConfig.Merge(dst, src)
}
func (m *RuleRoutingConfig) XXX_Size() int {
	return xxx_messageInfo_RuleRoutingConfig.Size(m)
}
func (m *RuleRoutingConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RuleRoutingConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RuleRoutingConfig proto.InternalMessageInfo

func (m *RuleRoutingConfig) GetSources() []*Source {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *RuleRoutingConfig) GetDestinations() []*Destination {
	if m != nil {
		return m.Destinations
	}
	return nil
}

type Source struct {
	// Main tuning service and namespace
	Service   string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Master Control Service Example Tag or Request Label
	// Value supports regular matching
	Arguments            []*SourceMatch `protobuf:"bytes,3,rep,name=arguments,proto3" json:"arguments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Source) Reset()         { *m = Source{} }
func (m *Source) String() string { return proto.CompactTextString(m) }
func (*Source) ProtoMessage()    {}
func (*Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{7}
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

func (m *Source) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Source) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Source) GetArguments() []*SourceMatch {
	if m != nil {
		return m.Arguments
	}
	return nil
}

type Destination struct {
	// Templated service and namespace
	Service   string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Templated service example label
	// Value supports regular matching
	Labels map[string]*MatchString `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// According to the service name and service instance Metadata Filter the
	// qualified service instance subset Service instance subset can set priority
	// and weight Priority: integer, range [0, 9], the highest priority is 0
	// Weight: Integer
	// Press priority routing, if there is high priority, low priority will not
	// use If there is a subset of the same priority, then assign by weight
	// Priority and weight can be not set / set up one / set two
	// If the section is set priority, some are not set, it is considered that the
	// priority is not set. If the part is set, some is not set, it is considered
	// that the weight is not set to 0 If you have no weight, you think the weight
	// is the same
	Priority uint32 `protobuf:"varint,4,opt,name=priority,proto3" json:"priority,omitempty"`
	Weight   uint32 `protobuf:"varint,5,opt,name=weight,proto3" json:"weight,omitempty"`
	// Forward requests to proxy service
	Transfer string `protobuf:"bytes,6,opt,name=transfer,proto3" json:"transfer,omitempty"`
	// Whether to isolate the SET, after isolation, no traffic will be allocated
	Isolate              bool     `protobuf:"varint,7,opt,name=isolate,proto3" json:"isolate,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Destination) Reset()         { *m = Destination{} }
func (m *Destination) String() string { return proto.CompactTextString(m) }
func (*Destination) ProtoMessage()    {}
func (*Destination) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{8}
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

func (m *Destination) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Destination) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Destination) GetLabels() map[string]*MatchString {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Destination) GetPriority() uint32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *Destination) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *Destination) GetTransfer() string {
	if m != nil {
		return m.Transfer
	}
	return ""
}

func (m *Destination) GetIsolate() bool {
	if m != nil {
		return m.Isolate
	}
	return false
}

type SourceMatch struct {
	Type SourceMatch_Type `protobuf:"varint,1,opt,name=type,proto3,enum=v2.SourceMatch_Type" json:"type,omitempty"`
	// header key or query key
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// header value or query value
	Value                *MatchString `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SourceMatch) Reset()         { *m = SourceMatch{} }
func (m *SourceMatch) String() string { return proto.CompactTextString(m) }
func (*SourceMatch) ProtoMessage()    {}
func (*SourceMatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_routing_v2_eb1407801dce78dd, []int{9}
}
func (m *SourceMatch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SourceMatch.Unmarshal(m, b)
}
func (m *SourceMatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SourceMatch.Marshal(b, m, deterministic)
}
func (dst *SourceMatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SourceMatch.Merge(dst, src)
}
func (m *SourceMatch) XXX_Size() int {
	return xxx_messageInfo_SourceMatch.Size(m)
}
func (m *SourceMatch) XXX_DiscardUnknown() {
	xxx_messageInfo_SourceMatch.DiscardUnknown(m)
}

var xxx_messageInfo_SourceMatch proto.InternalMessageInfo

func (m *SourceMatch) GetType() SourceMatch_Type {
	if m != nil {
		return m.Type
	}
	return SourceMatch_CUSTOM
}

func (m *SourceMatch) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SourceMatch) GetValue() *MatchString {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*FlowStaining)(nil), "v2.FlowStaining")
	proto.RegisterType((*StaineRule)(nil), "v2.StaineRule")
	proto.RegisterType((*StaineLabel)(nil), "v2.StaineLabel")
	proto.RegisterType((*Routing)(nil), "v2.Routing")
	proto.RegisterType((*MetadataFailover)(nil), "v2.MetadataFailover")
	proto.RegisterMapType((map[string]string)(nil), "v2.MetadataFailover.LabelsEntry")
	proto.RegisterType((*MetadataRoutingConfig)(nil), "v2.MetadataRoutingConfig")
	proto.RegisterMapType((map[string]string)(nil), "v2.MetadataRoutingConfig.LabelsEntry")
	proto.RegisterType((*RuleRoutingConfig)(nil), "v2.RuleRoutingConfig")
	proto.RegisterType((*Source)(nil), "v2.Source")
	proto.RegisterType((*Destination)(nil), "v2.Destination")
	proto.RegisterMapType((map[string]*MatchString)(nil), "v2.Destination.LabelsEntry")
	proto.RegisterType((*SourceMatch)(nil), "v2.SourceMatch")
	proto.RegisterEnum("v2.RoutingPolicy", RoutingPolicy_name, RoutingPolicy_value)
	proto.RegisterEnum("v2.MetadataFailover_FailoverRange", MetadataFailover_FailoverRange_name, MetadataFailover_FailoverRange_value)
	proto.RegisterEnum("v2.SourceMatch_Type", SourceMatch_Type_name, SourceMatch_Type_value)
}

func init() { proto.RegisterFile("routing_v2.proto", fileDescriptor_routing_v2_eb1407801dce78dd) }

var fileDescriptor_routing_v2_eb1407801dce78dd = []byte{
	// 855 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xdd, 0x8e, 0x1b, 0x35,
	0x14, 0xee, 0x4c, 0xfe, 0x4f, 0x9a, 0x61, 0x6a, 0x2d, 0x68, 0x08, 0x5c, 0x44, 0xa3, 0xad, 0x88,
	0x90, 0x48, 0x51, 0x16, 0xa4, 0x16, 0xc1, 0x45, 0xb4, 0x9b, 0x6a, 0x4b, 0xb3, 0xec, 0xe2, 0xa4,
	0x88, 0x5e, 0xad, 0xbc, 0x13, 0x67, 0x6a, 0x75, 0xe2, 0x89, 0x3c, 0x4e, 0xaa, 0xbc, 0x01, 0x8f,
	0xc3, 0x13, 0x20, 0x9e, 0x82, 0x97, 0xe1, 0x06, 0xd9, 0x1e, 0x27, 0x33, 0x21, 0xb4, 0xd2, 0x72,
	0x35, 0x3e, 0xdf, 0x1c, 0xfb, 0xf8, 0x7c, 0xe7, 0xf3, 0x07, 0xbe, 0x48, 0xd7, 0x92, 0xf1, 0xf8,
	0x76, 0x33, 0x1c, 0xac, 0x44, 0x2a, 0x53, 0xe4, 0x6e, 0x86, 0x5d, 0x6f, 0x99, 0xce, 0x69, 0xb2,
	0xc3, 0xba, 0x9f, 0xc6, 0x69, 0x1a, 0x27, 0xf4, 0x89, 0x8e, 0xee, 0xd6, 0x8b, 0x27, 0x84, 0x6f,
	0xcd, 0xaf, 0xf0, 0x57, 0x78, 0xf8, 0x3c, 0x49, 0xdf, 0x4d, 0x25, 0x61, 0x9c, 0xf1, 0x18, 0x79,
	0xe0, 0xb2, 0x79, 0x00, 0x3d, 0xa7, 0xdf, 0xc2, 0x2e, 0x9b, 0x23, 0x04, 0x55, 0x4e, 0x96, 0x34,
	0x70, 0x35, 0xa2, 0xd7, 0xe8, 0x14, 0x6a, 0x62, 0x9d, 0xd0, 0x2c, 0xa8, 0xf4, 0x2a, 0xfd, 0xf6,
	0xd0, 0x1b, 0x6c, 0x86, 0x03, 0x7d, 0x00, 0xc5, 0xeb, 0x84, 0x62, 0xf3, 0x33, 0xfc, 0xc3, 0x01,
	0xd8, 0xa3, 0xe8, 0x2b, 0x68, 0x11, 0x11, 0xaf, 0x97, 0x94, 0xcb, 0x2c, 0x70, 0xf4, 0xc6, 0x8f,
	0xf4, 0xc6, 0x74, 0x2d, 0x22, 0x7a, 0x45, 0x64, 0xf4, 0x06, 0xef, 0x33, 0xd0, 0x17, 0x50, 0x4f,
	0xc8, 0x1d, 0x4d, 0xb2, 0xc0, 0x2d, 0xe4, 0xea, 0xe3, 0x26, 0x0a, 0xc7, 0xf9, 0x6f, 0xd4, 0x85,
	0xe6, 0x4a, 0xb0, 0x54, 0x30, 0xb9, 0x0d, 0x2a, 0x3d, 0xa7, 0xdf, 0xc1, 0xbb, 0x18, 0x7d, 0x02,
	0x75, 0xca, 0xc9, 0x5d, 0x42, 0x83, 0x6a, 0xcf, 0xe9, 0x37, 0x71, 0x1e, 0xa1, 0x53, 0xe8, 0x64,
	0xfa, 0xa8, 0x1b, 0x2a, 0x22, 0xca, 0x65, 0x50, 0xd3, 0x1b, 0xcb, 0x60, 0xf8, 0x2d, 0xb4, 0x0b,
	0x05, 0x91, 0x0f, 0x95, 0x97, 0x74, 0x1b, 0x38, 0x9a, 0x08, 0xb5, 0x44, 0x27, 0x50, 0xfb, 0x85,
	0x24, 0x6b, 0x4b, 0x8e, 0x09, 0xc2, 0xbf, 0x5c, 0x68, 0x60, 0x33, 0x95, 0x9c, 0x4d, 0xe7, 0xbd,
	0x6c, 0x7e, 0x0e, 0x2d, 0xf5, 0xcd, 0x56, 0x24, 0xa2, 0xba, 0x83, 0x16, 0xde, 0x03, 0xff, 0xd9,
	0xc2, 0x33, 0xf0, 0xec, 0xe8, 0x57, 0x69, 0xc2, 0xa2, 0xad, 0xee, 0xc1, 0x1b, 0x3e, 0x52, 0x3c,
	0xe5, 0xe5, 0x6f, 0xf4, 0x0f, 0x7c, 0x90, 0x88, 0xbe, 0xdf, 0x6f, 0x8d, 0x52, 0xbe, 0x60, 0x71,
	0x50, 0xef, 0x39, 0xfd, 0xf6, 0xf0, 0x64, 0x60, 0x64, 0x32, 0xb0, 0x32, 0x19, 0x8c, 0x78, 0x61,
	0xb7, 0xc9, 0x55, 0x7c, 0x0b, 0xba, 0x61, 0x19, 0x4b, 0x79, 0xd0, 0xd0, 0xb7, 0xdd, 0xc5, 0x8a,
	0x90, 0x48, 0xb2, 0x25, 0x0d, 0x9a, 0x86, 0x10, 0x1d, 0x28, 0x74, 0xa9, 0xd1, 0x96, 0x41, 0x97,
	0x16, 0xa5, 0x1a, 0x35, 0x5a, 0x33, 0x41, 0x69, 0x9a, 0xed, 0xf2, 0x34, 0xc3, 0xdf, 0x5c, 0xf0,
	0xaf, 0xa8, 0x24, 0x73, 0x22, 0xc9, 0x73, 0xc2, 0x92, 0x74, 0x43, 0x05, 0x7a, 0x01, 0xde, 0x22,
	0x5f, 0xdf, 0x0a, 0xc2, 0x63, 0xaa, 0xd9, 0xf6, 0x86, 0xa1, 0xe2, 0xe1, 0x30, 0x7b, 0x60, 0x17,
	0x58, 0x65, 0xe2, 0xce, 0xa2, 0x18, 0xa2, 0xa7, 0x07, 0x92, 0xeb, 0x1d, 0x3d, 0x42, 0x8b, 0x21,
	0x1b, 0x73, 0x29, 0xb6, 0x56, 0x83, 0xdd, 0x67, 0xd0, 0x2e, 0xc0, 0x4a, 0x29, 0x6f, 0xf7, 0x4a,
	0x79, 0x6b, 0x94, 0xb2, 0x29, 0x2a, 0x45, 0x07, 0xdf, 0xb9, 0x4f, 0x9d, 0xf0, 0x1b, 0xe8, 0x94,
	0x2e, 0x85, 0x1a, 0x50, 0x19, 0x4d, 0x26, 0xfe, 0x03, 0x04, 0x50, 0xbf, 0x9e, 0x5d, 0x8e, 0xf1,
	0xd4, 0x77, 0x90, 0x07, 0xa0, 0xd7, 0xb7, 0x2f, 0xc7, 0xaf, 0xa7, 0xbe, 0x1b, 0xfe, 0xed, 0xc0,
	0xc7, 0xf6, 0x66, 0xf9, 0xb0, 0xcf, 0xcd, 0x78, 0x02, 0x68, 0x64, 0x54, 0x6c, 0x58, 0x44, 0xf3,
	0xfa, 0x36, 0x2c, 0xeb, 0xcc, 0x3d, 0xd4, 0xd9, 0x0f, 0xbb, 0xe6, 0xcd, 0xa3, 0x7e, 0x5c, 0x6c,
	0xbe, 0x54, 0xe2, 0x18, 0x03, 0xe8, 0x6b, 0x68, 0x5a, 0x32, 0xb5, 0x50, 0x95, 0x9a, 0x8e, 0xb0,
	0x87, 0x77, 0x59, 0xff, 0x87, 0x33, 0x0e, 0x8f, 0xb4, 0xd1, 0x94, 0x1a, 0x3f, 0x85, 0x46, 0xa6,
	0xad, 0xc4, 0xba, 0x0b, 0xec, 0xdd, 0x05, 0xdb, 0x5f, 0xe8, 0x0c, 0x1e, 0xce, 0x69, 0x26, 0x19,
	0x27, 0x92, 0xa5, 0xbc, 0x64, 0x2e, 0x17, 0x7b, 0x1c, 0x97, 0x92, 0xc2, 0x14, 0xea, 0xe6, 0x9c,
	0x7b, 0xb3, 0x5b, 0x32, 0xbf, 0xca, 0x87, 0xcc, 0x2f, 0xfc, 0xdd, 0x85, 0x76, 0xe1, 0x3a, 0xf7,
	0x2e, 0x7b, 0x76, 0x30, 0xd4, 0xcf, 0x0e, 0xfa, 0x3c, 0x3a, 0xca, 0xe2, 0x13, 0xac, 0xfe, 0xdb,
	0x50, 0xdf, 0x51, 0x16, 0xbf, 0xb1, 0x8e, 0x99, 0x47, 0x6a, 0x8f, 0x14, 0x84, 0x67, 0x0b, 0x2a,
	0xb4, 0x99, 0xb4, 0xf0, 0x2e, 0x56, 0x97, 0x67, 0x59, 0x9a, 0x10, 0x49, 0xb5, 0x5f, 0x34, 0xb1,
	0x0d, 0xbb, 0x3f, 0x7e, 0x48, 0x02, 0x8f, 0x8b, 0x12, 0xc8, 0x29, 0xd3, 0x64, 0x4d, 0xa5, 0x60,
	0x3c, 0x2e, 0x6a, 0xe2, 0x4f, 0x07, 0xda, 0x05, 0x36, 0x51, 0x1f, 0xaa, 0x72, 0xbb, 0xb2, 0x6e,
	0x70, 0x72, 0x40, 0xf6, 0x60, 0xb6, 0x5d, 0x51, 0xac, 0x33, 0x6c, 0x59, 0xf7, 0x48, 0xd9, 0xca,
	0xfb, 0xca, 0x86, 0x3f, 0x41, 0x55, 0x1d, 0xa3, 0x1e, 0xea, 0xf9, 0xab, 0xe9, 0xec, 0xfa, 0xca,
	0x3c, 0xda, 0xab, 0xf1, 0xec, 0xf2, 0xfa, 0xc2, 0x77, 0xd4, 0xfa, 0x72, 0x3c, 0xba, 0x18, 0x63,
	0xdf, 0x45, 0x2d, 0xa8, 0xfd, 0xfc, 0x6a, 0x8c, 0x5f, 0xfb, 0x15, 0xd4, 0x81, 0xd6, 0xf9, 0x68,
	0x32, 0x19, 0xe3, 0xdb, 0x17, 0x37, 0x7e, 0x15, 0x35, 0xa1, 0x7a, 0x33, 0x9a, 0x5d, 0xfa, 0xb5,
	0x2f, 0xcf, 0xa0, 0x53, 0x32, 0x6e, 0xf5, 0xea, 0x95, 0xce, 0x4d, 0xe4, 0x3f, 0x40, 0x08, 0x3c,
	0xfb, 0xa0, 0x72, 0xcc, 0xb9, 0xab, 0x6b, 0xb3, 0x3e, 0xfb, 0x27, 0x00, 0x00, 0xff, 0xff, 0x85,
	0x1e, 0x7c, 0xfb, 0x09, 0x08, 0x00, 0x00,
}
