// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/monitored_resource.proto

package monitoredres // import "google.golang.org/genproto/googleapis/api/monitoredres"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _struct "github.com/golang/protobuf/ptypes/struct"
import label "google.golang.org/genproto/googleapis/api/label"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// An object that describes the schema of a [MonitoredResource][google.api.MonitoredResource] object using a
// type name and a set of labels.  For example, the monitored resource
// descriptor for Google Compute Engine VM instances has a type of
// `"gce_instance"` and specifies the use of the labels `"instance_id"` and
// `"zone"` to identify particular VM instances.
//
// Different APIs can support different monitored resource types. APIs generally
// provide a `list` method that returns the monitored resource descriptors used
// by the API.
type MonitoredResourceDescriptor struct {
	// Optional. The resource name of the monitored resource descriptor:
	// `"projects/{project_id}/monitoredResourceDescriptors/{type}"` where
	// {type} is the value of the `type` field in this object and
	// {project_id} is a project ID that provides API-specific context for
	// accessing the type.  APIs that do not use project information can use the
	// resource name format `"monitoredResourceDescriptors/{type}"`.
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// Required. The monitored resource type. For example, the type
	// `"cloudsql_database"` represents databases in Google Cloud SQL.
	// The maximum length of this value is 256 characters.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Optional. A concise name for the monitored resource type that might be
	// displayed in user interfaces. It should be a Title Cased Noun Phrase,
	// without any article or other determiners. For example,
	// `"Google Cloud SQL Database"`.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Optional. A detailed description of the monitored resource type that might
	// be used in documentation.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Required. A set of labels used to describe instances of this monitored
	// resource type. For example, an individual Google Cloud SQL database is
	// identified by values for the labels `"database_id"` and `"zone"`.
	Labels               []*label.LabelDescriptor `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *MonitoredResourceDescriptor) Reset()         { *m = MonitoredResourceDescriptor{} }
func (m *MonitoredResourceDescriptor) String() string { return proto.CompactTextString(m) }
func (*MonitoredResourceDescriptor) ProtoMessage()    {}
func (*MonitoredResourceDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cd8bd738b08f2bf, []int{0}
}
func (m *MonitoredResourceDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitoredResourceDescriptor.Unmarshal(m, b)
}
func (m *MonitoredResourceDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitoredResourceDescriptor.Marshal(b, m, deterministic)
}
func (m *MonitoredResourceDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitoredResourceDescriptor.Merge(m, src)
}
func (m *MonitoredResourceDescriptor) XXX_Size() int {
	return xxx_messageInfo_MonitoredResourceDescriptor.Size(m)
}
func (m *MonitoredResourceDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitoredResourceDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_MonitoredResourceDescriptor proto.InternalMessageInfo

func (m *MonitoredResourceDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *MonitoredResourceDescriptor) GetLabels() []*label.LabelDescriptor {
	if m != nil {
		return m.Labels
	}
	return nil
}

// An object representing a resource that can be used for monitoring, logging,
// billing, or other purposes. Examples include virtual machine instances,
// databases, and storage devices such as disks. The `type` field identifies a
// [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] object that describes the resource's
// schema. Information in the `labels` field identifies the actual resource and
// its attributes according to the schema. For example, a particular Compute
// Engine VM instance could be represented by the following object, because the
// [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] for `"gce_instance"` has labels
// `"instance_id"` and `"zone"`:
//
//     { "type": "gce_instance",
//       "labels": { "instance_id": "12345678901234",
//                   "zone": "us-central1-a" }}
type MonitoredResource struct {
	// Required. The monitored resource type. This field must match
	// the `type` field of a [MonitoredResourceDescriptor][google.api.MonitoredResourceDescriptor] object. For
	// example, the type of a Compute Engine VM instance is `gce_instance`.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Required. Values for all of the labels listed in the associated monitored
	// resource descriptor. For example, Compute Engine VM instances use the
	// labels `"project_id"`, `"instance_id"`, and `"zone"`.
	Labels               map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MonitoredResource) Reset()         { *m = MonitoredResource{} }
func (m *MonitoredResource) String() string { return proto.CompactTextString(m) }
func (*MonitoredResource) ProtoMessage()    {}
func (*MonitoredResource) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cd8bd738b08f2bf, []int{1}
}
func (m *MonitoredResource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitoredResource.Unmarshal(m, b)
}
func (m *MonitoredResource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitoredResource.Marshal(b, m, deterministic)
}
func (m *MonitoredResource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitoredResource.Merge(m, src)
}
func (m *MonitoredResource) XXX_Size() int {
	return xxx_messageInfo_MonitoredResource.Size(m)
}
func (m *MonitoredResource) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitoredResource.DiscardUnknown(m)
}

var xxx_messageInfo_MonitoredResource proto.InternalMessageInfo

func (m *MonitoredResource) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MonitoredResource) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

// Auxiliary metadata for a [MonitoredResource][google.api.MonitoredResource] object.
// [MonitoredResource][google.api.MonitoredResource] objects contain the minimum set of information to
// uniquely identify a monitored resource instance. There is some other useful
// auxiliary metadata. Google Stackdriver Monitoring & Logging uses an ingestion
// pipeline to extract metadata for cloud resources of all types , and stores
// the metadata in this message.
type MonitoredResourceMetadata struct {
	// Output only. Values for predefined system metadata labels.
	// System labels are a kind of metadata extracted by Google Stackdriver.
	// Stackdriver determines what system labels are useful and how to obtain
	// their values. Some examples: "machine_image", "vpc", "subnet_id",
	// "security_group", "name", etc.
	// System label values can be only strings, Boolean values, or a list of
	// strings. For example:
	//
	//     { "name": "my-test-instance",
	//       "security_group": ["a", "b", "c"],
	//       "spot_instance": false }
	SystemLabels *_struct.Struct `protobuf:"bytes,1,opt,name=system_labels,json=systemLabels,proto3" json:"system_labels,omitempty"`
	// Output only. A map of user-defined metadata labels.
	UserLabels           map[string]string `protobuf:"bytes,2,rep,name=user_labels,json=userLabels,proto3" json:"user_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MonitoredResourceMetadata) Reset()         { *m = MonitoredResourceMetadata{} }
func (m *MonitoredResourceMetadata) String() string { return proto.CompactTextString(m) }
func (*MonitoredResourceMetadata) ProtoMessage()    {}
func (*MonitoredResourceMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cd8bd738b08f2bf, []int{2}
}
func (m *MonitoredResourceMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitoredResourceMetadata.Unmarshal(m, b)
}
func (m *MonitoredResourceMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitoredResourceMetadata.Marshal(b, m, deterministic)
}
func (m *MonitoredResourceMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitoredResourceMetadata.Merge(m, src)
}
func (m *MonitoredResourceMetadata) XXX_Size() int {
	return xxx_messageInfo_MonitoredResourceMetadata.Size(m)
}
func (m *MonitoredResourceMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitoredResourceMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_MonitoredResourceMetadata proto.InternalMessageInfo

func (m *MonitoredResourceMetadata) GetSystemLabels() *_struct.Struct {
	if m != nil {
		return m.SystemLabels
	}
	return nil
}

func (m *MonitoredResourceMetadata) GetUserLabels() map[string]string {
	if m != nil {
		return m.UserLabels
	}
	return nil
}

func init() {
	proto.RegisterType((*MonitoredResourceDescriptor)(nil), "google.api.MonitoredResourceDescriptor")
	proto.RegisterType((*MonitoredResource)(nil), "google.api.MonitoredResource")
	proto.RegisterMapType((map[string]string)(nil), "google.api.MonitoredResource.LabelsEntry")
	proto.RegisterType((*MonitoredResourceMetadata)(nil), "google.api.MonitoredResourceMetadata")
	proto.RegisterMapType((map[string]string)(nil), "google.api.MonitoredResourceMetadata.UserLabelsEntry")
}

func init() {
	proto.RegisterFile("google/api/monitored_resource.proto", fileDescriptor_6cd8bd738b08f2bf)
}

var fileDescriptor_6cd8bd738b08f2bf = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x4d, 0xab, 0xd3, 0x40,
	0x14, 0x65, 0xd2, 0x0f, 0xf0, 0xa6, 0x7e, 0x0d, 0x52, 0x63, 0xea, 0xa2, 0xd6, 0x4d, 0xdd, 0x24,
	0xd0, 0x22, 0xf8, 0xb9, 0x68, 0x55, 0x44, 0xb0, 0x52, 0x22, 0xba, 0x70, 0x13, 0xa6, 0xc9, 0x18,
	0x82, 0x49, 0x26, 0xcc, 0x4c, 0x84, 0xfc, 0x1d, 0xc1, 0xdf, 0xe1, 0x5f, 0x72, 0xe9, 0x52, 0x32,
	0x33, 0x69, 0xd3, 0x97, 0xc7, 0x83, 0xb7, 0xbb, 0xf7, 0xdc, 0x73, 0xcf, 0x3d, 0x27, 0x43, 0xe0,
	0x71, 0xc2, 0x58, 0x92, 0x51, 0x9f, 0x94, 0xa9, 0x9f, 0xb3, 0x22, 0x95, 0x8c, 0xd3, 0x38, 0xe4,
	0x54, 0xb0, 0x8a, 0x47, 0xd4, 0x2b, 0x39, 0x93, 0x0c, 0x83, 0x26, 0x79, 0xa4, 0x4c, 0xdd, 0x69,
	0x67, 0x21, 0x23, 0x07, 0x9a, 0x69, 0x8e, 0xfb, 0xd0, 0xe0, 0xaa, 0x3b, 0x54, 0xdf, 0x7d, 0x21,
	0x79, 0x15, 0x49, 0x3d, 0x5d, 0xfc, 0x41, 0x30, 0xdb, 0xb5, 0xf2, 0x81, 0x51, 0x7f, 0x4b, 0x45,
	0xc4, 0xd3, 0x52, 0x32, 0x8e, 0x31, 0x0c, 0x0b, 0x92, 0x53, 0x67, 0x34, 0x47, 0xcb, 0x1b, 0x81,
	0xaa, 0x1b, 0x4c, 0xd6, 0x25, 0x75, 0x90, 0xc6, 0x9a, 0x1a, 0x3f, 0x82, 0x49, 0x9c, 0x8a, 0x32,
	0x23, 0x75, 0xa8, 0xf8, 0x96, 0x9a, 0xd9, 0x06, 0xfb, 0xd4, 0xac, 0xcd, 0xc1, 0x8e, 0x8d, 0x70,
	0xca, 0x0a, 0x67, 0x60, 0x18, 0x27, 0x08, 0xaf, 0x61, 0xac, 0x9c, 0x0b, 0x67, 0x38, 0x1f, 0x2c,
	0xed, 0xd5, 0xcc, 0x3b, 0xe5, 0xf3, 0x3e, 0x36, 0x93, 0x93, 0xb3, 0xc0, 0x50, 0x17, 0xbf, 0x11,
	0xdc, 0xed, 0x25, 0xb8, 0xd4, 0xe3, 0xe6, 0x28, 0x6f, 0x29, 0xf9, 0x27, 0x5d, 0xf9, 0x9e, 0x84,
	0x3e, 0x28, 0xde, 0x15, 0x92, 0xd7, 0xed, 0x31, 0xf7, 0x39, 0xd8, 0x1d, 0x18, 0xdf, 0x81, 0xc1,
	0x0f, 0x5a, 0x9b, 0x23, 0x4d, 0x89, 0xef, 0xc1, 0xe8, 0x27, 0xc9, 0xaa, 0xf6, 0x03, 0xe8, 0xe6,
	0x85, 0xf5, 0x0c, 0x2d, 0xfe, 0x22, 0x78, 0xd0, 0x3b, 0xb2, 0xa3, 0x92, 0xc4, 0x44, 0x12, 0xfc,
	0x0a, 0x6e, 0x8a, 0x5a, 0x48, 0x9a, 0x87, 0xc6, 0x62, 0xa3, 0x69, 0xaf, 0xee, 0xb7, 0x16, 0xdb,
	0xd7, 0xf3, 0x3e, 0xab, 0xd7, 0x0b, 0x26, 0x9a, 0xad, 0xcd, 0xe0, 0xaf, 0x60, 0x57, 0x82, 0xf2,
	0xf0, 0x2c, 0xde, 0xd3, 0x2b, 0xe3, 0xb5, 0x97, 0xbd, 0x2f, 0x82, 0xf2, 0x6e, 0x54, 0xa8, 0x8e,
	0x80, 0xfb, 0x1a, 0x6e, 0x5f, 0x18, 0x5f, 0x27, 0xf2, 0xb6, 0x86, 0x5b, 0x11, 0xcb, 0x3b, 0x36,
	0xb6, 0xd3, 0x9e, 0x8f, 0x7d, 0x13, 0x6c, 0x8f, 0xbe, 0xbd, 0x31, 0xac, 0x84, 0x65, 0xa4, 0x48,
	0x3c, 0xc6, 0x13, 0x3f, 0xa1, 0x85, 0x8a, 0xed, 0xeb, 0x11, 0x29, 0x53, 0x71, 0xfe, 0x3b, 0x70,
	0x2a, 0x5e, 0x76, 0x9b, 0x7f, 0x08, 0xfd, 0xb2, 0x86, 0xef, 0x37, 0xfb, 0x0f, 0x87, 0xb1, 0xda,
	0x5c, 0xff, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x10, 0x16, 0x7c, 0xe9, 0x47, 0x03, 0x00, 0x00,
}
