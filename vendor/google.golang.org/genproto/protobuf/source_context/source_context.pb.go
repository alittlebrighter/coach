// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/protobuf/source_context.proto

package source_context // import "google.golang.org/genproto/protobuf/source_context"

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

// `SourceContext` represents information about the source of a
// protobuf element, like the file in which it is defined.
type SourceContext struct {
	// The path-qualified name of the .proto file that contained the associated
	// protobuf element.  For example: `"google/protobuf/source_context.proto"`.
	FileName             string   `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SourceContext) Reset()         { *m = SourceContext{} }
func (m *SourceContext) String() string { return proto.CompactTextString(m) }
func (*SourceContext) ProtoMessage()    {}
func (*SourceContext) Descriptor() ([]byte, []int) {
	return fileDescriptor_b686cdb126d509db, []int{0}
}
func (m *SourceContext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SourceContext.Unmarshal(m, b)
}
func (m *SourceContext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SourceContext.Marshal(b, m, deterministic)
}
func (m *SourceContext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SourceContext.Merge(m, src)
}
func (m *SourceContext) XXX_Size() int {
	return xxx_messageInfo_SourceContext.Size(m)
}
func (m *SourceContext) XXX_DiscardUnknown() {
	xxx_messageInfo_SourceContext.DiscardUnknown(m)
}

var xxx_messageInfo_SourceContext proto.InternalMessageInfo

func (m *SourceContext) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func init() {
	proto.RegisterType((*SourceContext)(nil), "google.protobuf.SourceContext")
}

func init() {
	proto.RegisterFile("google/protobuf/source_context.proto", fileDescriptor_b686cdb126d509db)
}

var fileDescriptor_b686cdb126d509db = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x2f, 0xce, 0x2f, 0x2d,
	0x4a, 0x4e, 0x8d, 0x4f, 0xce, 0xcf, 0x2b, 0x49, 0xad, 0x28, 0xd1, 0x03, 0x8b, 0x0b, 0xf1, 0x43,
	0x54, 0xe9, 0xc1, 0x54, 0x29, 0xe9, 0x70, 0xf1, 0x06, 0x83, 0x15, 0x3a, 0x43, 0xd4, 0x09, 0x49,
	0x73, 0x71, 0xa6, 0x65, 0xe6, 0xa4, 0xc6, 0xe7, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a,
	0x70, 0x06, 0x71, 0x80, 0x04, 0xfc, 0x12, 0x73, 0x53, 0x9d, 0xa6, 0x32, 0x72, 0x09, 0x27, 0xe7,
	0xe7, 0xea, 0xa1, 0x99, 0xe2, 0x24, 0x84, 0x62, 0x46, 0x00, 0x48, 0x38, 0x80, 0x31, 0xca, 0x11,
	0xaa, 0x2c, 0x3d, 0x3f, 0x27, 0x31, 0x2f, 0x5d, 0x2f, 0xbf, 0x28, 0x5d, 0x3f, 0x3d, 0x35, 0x0f,
	0xac, 0x09, 0x97, 0x33, 0xad, 0x51, 0xb9, 0x8b, 0x98, 0x98, 0xdd, 0x03, 0x9c, 0x56, 0x31, 0xc9,
	0xb9, 0x43, 0x4c, 0x0a, 0x80, 0xea, 0xd2, 0x0b, 0x4f, 0xcd, 0xc9, 0xf1, 0xce, 0xcb, 0x2f, 0xcf,
	0x0b, 0xa9, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03, 0x1b, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff,
	0x5c, 0xbd, 0xa4, 0x22, 0x05, 0x01, 0x00, 0x00,
}
