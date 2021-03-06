// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

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

type ScriptsQuery struct {
	Query                string   `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScriptsQuery) Reset()         { *m = ScriptsQuery{} }
func (m *ScriptsQuery) String() string { return proto.CompactTextString(m) }
func (*ScriptsQuery) ProtoMessage()    {}
func (*ScriptsQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *ScriptsQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScriptsQuery.Unmarshal(m, b)
}
func (m *ScriptsQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScriptsQuery.Marshal(b, m, deterministic)
}
func (m *ScriptsQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScriptsQuery.Merge(m, src)
}
func (m *ScriptsQuery) XXX_Size() int {
	return xxx_messageInfo_ScriptsQuery.Size(m)
}
func (m *ScriptsQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ScriptsQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ScriptsQuery proto.InternalMessageInfo

func (m *ScriptsQuery) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type GetScriptsResponse struct {
	Scripts              []*DocumentedScript `protobuf:"bytes,1,rep,name=scripts,proto3" json:"scripts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetScriptsResponse) Reset()         { *m = GetScriptsResponse{} }
func (m *GetScriptsResponse) String() string { return proto.CompactTextString(m) }
func (*GetScriptsResponse) ProtoMessage()    {}
func (*GetScriptsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *GetScriptsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetScriptsResponse.Unmarshal(m, b)
}
func (m *GetScriptsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetScriptsResponse.Marshal(b, m, deterministic)
}
func (m *GetScriptsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetScriptsResponse.Merge(m, src)
}
func (m *GetScriptsResponse) XXX_Size() int {
	return xxx_messageInfo_GetScriptsResponse.Size(m)
}
func (m *GetScriptsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetScriptsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetScriptsResponse proto.InternalMessageInfo

func (m *GetScriptsResponse) GetScripts() []*DocumentedScript {
	if m != nil {
		return m.Scripts
	}
	return nil
}

// RunEventIn has multiple uses, the first event should contain the tag of the script
// to run in the input.  Subsequent events should contain responses to requests for input
// and finally it should end with an EOF.
type RunEventIn struct {
	Input string    `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	Env   []*EnvVar `protobuf:"bytes,2,rep,name=env,proto3" json:"env,omitempty"`
	// specify the batch size in bytes you would like the response output to come in
	ResponseSize         uint32   `protobuf:"varint,3,opt,name=responseSize,proto3" json:"responseSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunEventIn) Reset()         { *m = RunEventIn{} }
func (m *RunEventIn) String() string { return proto.CompactTextString(m) }
func (*RunEventIn) ProtoMessage()    {}
func (*RunEventIn) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *RunEventIn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunEventIn.Unmarshal(m, b)
}
func (m *RunEventIn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunEventIn.Marshal(b, m, deterministic)
}
func (m *RunEventIn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunEventIn.Merge(m, src)
}
func (m *RunEventIn) XXX_Size() int {
	return xxx_messageInfo_RunEventIn.Size(m)
}
func (m *RunEventIn) XXX_DiscardUnknown() {
	xxx_messageInfo_RunEventIn.DiscardUnknown(m)
}

var xxx_messageInfo_RunEventIn proto.InternalMessageInfo

func (m *RunEventIn) GetInput() string {
	if m != nil {
		return m.Input
	}
	return ""
}

func (m *RunEventIn) GetEnv() []*EnvVar {
	if m != nil {
		return m.Env
	}
	return nil
}

func (m *RunEventIn) GetResponseSize() uint32 {
	if m != nil {
		return m.ResponseSize
	}
	return 0
}

type RunEventOut struct {
	Output string `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	// exitCode means nothing until stream finishes
	ExitCode             int32    `protobuf:"varint,3,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunEventOut) Reset()         { *m = RunEventOut{} }
func (m *RunEventOut) String() string { return proto.CompactTextString(m) }
func (*RunEventOut) ProtoMessage()    {}
func (*RunEventOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{3}
}

func (m *RunEventOut) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunEventOut.Unmarshal(m, b)
}
func (m *RunEventOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunEventOut.Marshal(b, m, deterministic)
}
func (m *RunEventOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunEventOut.Merge(m, src)
}
func (m *RunEventOut) XXX_Size() int {
	return xxx_messageInfo_RunEventOut.Size(m)
}
func (m *RunEventOut) XXX_DiscardUnknown() {
	xxx_messageInfo_RunEventOut.DiscardUnknown(m)
}

var xxx_messageInfo_RunEventOut proto.InternalMessageInfo

func (m *RunEventOut) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func (m *RunEventOut) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *RunEventOut) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

type EnvVar struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnvVar) Reset()         { *m = EnvVar{} }
func (m *EnvVar) String() string { return proto.CompactTextString(m) }
func (*EnvVar) ProtoMessage()    {}
func (*EnvVar) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{4}
}

func (m *EnvVar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnvVar.Unmarshal(m, b)
}
func (m *EnvVar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnvVar.Marshal(b, m, deterministic)
}
func (m *EnvVar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnvVar.Merge(m, src)
}
func (m *EnvVar) XXX_Size() int {
	return xxx_messageInfo_EnvVar.Size(m)
}
func (m *EnvVar) XXX_DiscardUnknown() {
	xxx_messageInfo_EnvVar.DiscardUnknown(m)
}

var xxx_messageInfo_EnvVar proto.InternalMessageInfo

func (m *EnvVar) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EnvVar) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SaveScriptRequest struct {
	Script               *DocumentedScript `protobuf:"bytes,1,opt,name=script,proto3" json:"script,omitempty"`
	Overwrite            bool              `protobuf:"varint,2,opt,name=overwrite,proto3" json:"overwrite,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SaveScriptRequest) Reset()         { *m = SaveScriptRequest{} }
func (m *SaveScriptRequest) String() string { return proto.CompactTextString(m) }
func (*SaveScriptRequest) ProtoMessage()    {}
func (*SaveScriptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{5}
}

func (m *SaveScriptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveScriptRequest.Unmarshal(m, b)
}
func (m *SaveScriptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveScriptRequest.Marshal(b, m, deterministic)
}
func (m *SaveScriptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveScriptRequest.Merge(m, src)
}
func (m *SaveScriptRequest) XXX_Size() int {
	return xxx_messageInfo_SaveScriptRequest.Size(m)
}
func (m *SaveScriptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveScriptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveScriptRequest proto.InternalMessageInfo

func (m *SaveScriptRequest) GetScript() *DocumentedScript {
	if m != nil {
		return m.Script
	}
	return nil
}

func (m *SaveScriptRequest) GetOverwrite() bool {
	if m != nil {
		return m.Overwrite
	}
	return false
}

type Response struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{6}
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

func (m *Response) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*ScriptsQuery)(nil), "io.alittlebrighter.coach.ScriptsQuery")
	proto.RegisterType((*GetScriptsResponse)(nil), "io.alittlebrighter.coach.GetScriptsResponse")
	proto.RegisterType((*RunEventIn)(nil), "io.alittlebrighter.coach.RunEventIn")
	proto.RegisterType((*RunEventOut)(nil), "io.alittlebrighter.coach.RunEventOut")
	proto.RegisterType((*EnvVar)(nil), "io.alittlebrighter.coach.EnvVar")
	proto.RegisterType((*SaveScriptRequest)(nil), "io.alittlebrighter.coach.SaveScriptRequest")
	proto.RegisterType((*Response)(nil), "io.alittlebrighter.coach.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CoachRPCClient is the client API for CoachRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CoachRPCClient interface {
	// query shoud be a comma delimited list of tags
	QueryScripts(ctx context.Context, in *ScriptsQuery, opts ...grpc.CallOption) (*GetScriptsResponse, error)
	// query should be a script alias
	GetScript(ctx context.Context, in *ScriptsQuery, opts ...grpc.CallOption) (*DocumentedScript, error)
	// the first RunEventIn event should contain the alias for the script
	// you want to run, subsequent events are treated as stdin passed to
	// the script while running
	RunScript(ctx context.Context, opts ...grpc.CallOption) (CoachRPC_RunScriptClient, error)
	SaveScript(ctx context.Context, in *SaveScriptRequest, opts ...grpc.CallOption) (*Response, error)
}

type coachRPCClient struct {
	cc *grpc.ClientConn
}

func NewCoachRPCClient(cc *grpc.ClientConn) CoachRPCClient {
	return &coachRPCClient{cc}
}

func (c *coachRPCClient) QueryScripts(ctx context.Context, in *ScriptsQuery, opts ...grpc.CallOption) (*GetScriptsResponse, error) {
	out := new(GetScriptsResponse)
	err := c.cc.Invoke(ctx, "/io.alittlebrighter.coach.CoachRPC/QueryScripts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coachRPCClient) GetScript(ctx context.Context, in *ScriptsQuery, opts ...grpc.CallOption) (*DocumentedScript, error) {
	out := new(DocumentedScript)
	err := c.cc.Invoke(ctx, "/io.alittlebrighter.coach.CoachRPC/GetScript", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coachRPCClient) RunScript(ctx context.Context, opts ...grpc.CallOption) (CoachRPC_RunScriptClient, error) {
	stream, err := c.cc.NewStream(ctx, &_CoachRPC_serviceDesc.Streams[0], "/io.alittlebrighter.coach.CoachRPC/RunScript", opts...)
	if err != nil {
		return nil, err
	}
	x := &coachRPCRunScriptClient{stream}
	return x, nil
}

type CoachRPC_RunScriptClient interface {
	Send(*RunEventIn) error
	Recv() (*RunEventOut, error)
	grpc.ClientStream
}

type coachRPCRunScriptClient struct {
	grpc.ClientStream
}

func (x *coachRPCRunScriptClient) Send(m *RunEventIn) error {
	return x.ClientStream.SendMsg(m)
}

func (x *coachRPCRunScriptClient) Recv() (*RunEventOut, error) {
	m := new(RunEventOut)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *coachRPCClient) SaveScript(ctx context.Context, in *SaveScriptRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/io.alittlebrighter.coach.CoachRPC/SaveScript", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoachRPCServer is the server API for CoachRPC service.
type CoachRPCServer interface {
	// query shoud be a comma delimited list of tags
	QueryScripts(context.Context, *ScriptsQuery) (*GetScriptsResponse, error)
	// query should be a script alias
	GetScript(context.Context, *ScriptsQuery) (*DocumentedScript, error)
	// the first RunEventIn event should contain the alias for the script
	// you want to run, subsequent events are treated as stdin passed to
	// the script while running
	RunScript(CoachRPC_RunScriptServer) error
	SaveScript(context.Context, *SaveScriptRequest) (*Response, error)
}

func RegisterCoachRPCServer(s *grpc.Server, srv CoachRPCServer) {
	s.RegisterService(&_CoachRPC_serviceDesc, srv)
}

func _CoachRPC_QueryScripts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScriptsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoachRPCServer).QueryScripts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/io.alittlebrighter.coach.CoachRPC/QueryScripts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoachRPCServer).QueryScripts(ctx, req.(*ScriptsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoachRPC_GetScript_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScriptsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoachRPCServer).GetScript(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/io.alittlebrighter.coach.CoachRPC/GetScript",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoachRPCServer).GetScript(ctx, req.(*ScriptsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoachRPC_RunScript_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CoachRPCServer).RunScript(&coachRPCRunScriptServer{stream})
}

type CoachRPC_RunScriptServer interface {
	Send(*RunEventOut) error
	Recv() (*RunEventIn, error)
	grpc.ServerStream
}

type coachRPCRunScriptServer struct {
	grpc.ServerStream
}

func (x *coachRPCRunScriptServer) Send(m *RunEventOut) error {
	return x.ServerStream.SendMsg(m)
}

func (x *coachRPCRunScriptServer) Recv() (*RunEventIn, error) {
	m := new(RunEventIn)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CoachRPC_SaveScript_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveScriptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoachRPCServer).SaveScript(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/io.alittlebrighter.coach.CoachRPC/SaveScript",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoachRPCServer).SaveScript(ctx, req.(*SaveScriptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CoachRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "io.alittlebrighter.coach.CoachRPC",
	HandlerType: (*CoachRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryScripts",
			Handler:    _CoachRPC_QueryScripts_Handler,
		},
		{
			MethodName: "GetScript",
			Handler:    _CoachRPC_GetScript_Handler,
		},
		{
			MethodName: "SaveScript",
			Handler:    _CoachRPC_SaveScript_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RunScript",
			Handler:       _CoachRPC_RunScript_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x5f, 0x6b, 0x13, 0x41,
	0x14, 0xc5, 0xd9, 0xc6, 0x24, 0x9b, 0x9b, 0xf4, 0xc1, 0x8b, 0xc8, 0xb2, 0xf8, 0x10, 0x86, 0x2a,
	0x41, 0x65, 0x91, 0xf8, 0xe6, 0x63, 0xd3, 0x22, 0x3e, 0x55, 0x27, 0xa0, 0x50, 0x14, 0xd9, 0x6e,
	0x2e, 0x76, 0x20, 0x9d, 0x49, 0xe7, 0xcf, 0xaa, 0xfd, 0x10, 0x7e, 0x66, 0xd9, 0x99, 0xd9, 0x44,
	0xa5, 0x5b, 0x9b, 0xa7, 0x9d, 0x73, 0xb9, 0x73, 0x7f, 0x97, 0x33, 0x67, 0xe1, 0xd0, 0x90, 0xae,
	0x45, 0x45, 0xc5, 0x46, 0x2b, 0xab, 0x30, 0x13, 0xaa, 0x28, 0xd7, 0xc2, 0xda, 0x35, 0x5d, 0x68,
	0xf1, 0xed, 0xd2, 0x92, 0x2e, 0x2a, 0x55, 0x56, 0x97, 0xf9, 0xd8, 0x7f, 0x42, 0x1b, 0x3b, 0x82,
	0xc9, 0xb2, 0xd2, 0x62, 0x63, 0xcd, 0x07, 0x47, 0xfa, 0x27, 0x3e, 0x82, 0xfe, 0x75, 0x73, 0xc8,
	0x92, 0x69, 0x32, 0x1b, 0xf1, 0x20, 0xd8, 0x39, 0xe0, 0x5b, 0xb2, 0xb1, 0x91, 0x93, 0xd9, 0x28,
	0x69, 0x08, 0x4f, 0x60, 0x68, 0x42, 0x29, 0x4b, 0xa6, 0xbd, 0xd9, 0x78, 0xfe, 0xbc, 0xe8, 0x82,
	0x16, 0x27, 0xaa, 0x72, 0x57, 0x24, 0x2d, 0xad, 0xc2, 0x14, 0xde, 0x5e, 0x65, 0x37, 0x00, 0xdc,
	0xc9, 0xd3, 0x9a, 0xa4, 0x7d, 0x27, 0x1b, 0xbe, 0x90, 0x1b, 0x67, 0x5b, 0xbe, 0x17, 0x38, 0x87,
	0x1e, 0xc9, 0x3a, 0x3b, 0xf0, 0x94, 0x69, 0x37, 0xe5, 0x54, 0xd6, 0x1f, 0x4b, 0xcd, 0x9b, 0x66,
	0x64, 0x30, 0xd1, 0x71, 0xd3, 0xa5, 0xb8, 0xa1, 0xac, 0x37, 0x4d, 0x66, 0x87, 0xfc, 0xaf, 0x1a,
	0xfb, 0x04, 0xe3, 0x96, 0x7d, 0xe6, 0x2c, 0x3e, 0x86, 0x81, 0x72, 0x76, 0x47, 0x8f, 0xaa, 0x59,
	0x8a, 0xb4, 0x56, 0x3a, 0x3b, 0x08, 0x4b, 0x79, 0x81, 0x39, 0xa4, 0xf4, 0x43, 0xd8, 0x85, 0x5a,
	0x85, 0xe1, 0x7d, 0xbe, 0xd5, 0x6c, 0x0e, 0x83, 0xb0, 0x0b, 0x22, 0x3c, 0x90, 0xe5, 0x15, 0xc5,
	0x89, 0xfe, 0xdc, 0xcc, 0xab, 0xcb, 0xb5, 0xa3, 0x76, 0x9e, 0x17, 0xcc, 0xc1, 0xc3, 0x65, 0x59,
	0x53, 0xf4, 0x87, 0xae, 0x1d, 0x19, 0x8b, 0xc7, 0x30, 0x08, 0x46, 0xf9, 0x01, 0xfb, 0x59, 0x1c,
	0x6f, 0xe2, 0x13, 0x18, 0xa9, 0x9a, 0xf4, 0x77, 0x2d, 0x6c, 0x40, 0xa6, 0x7c, 0x57, 0x60, 0x6f,
	0x20, 0xdd, 0xbe, 0x68, 0x06, 0x43, 0xe3, 0xaa, 0x8a, 0x8c, 0xf1, 0xb8, 0x94, 0xb7, 0xf2, 0x76,
	0x0b, 0xe6, 0xbf, 0x7a, 0x90, 0x2e, 0x1a, 0x38, 0x7f, 0xbf, 0xc0, 0x15, 0x4c, 0x7c, 0x86, 0x62,
	0x4c, 0xf0, 0x59, 0xf7, 0xaa, 0x7f, 0x46, 0x2e, 0x7f, 0xd9, 0xdd, 0x77, 0x4b, 0xe8, 0xbe, 0xc2,
	0x68, 0x5b, 0xbd, 0x37, 0x62, 0x0f, 0xd7, 0xf0, 0x33, 0x8c, 0xb8, 0x93, 0x51, 0x1c, 0x75, 0x5f,
	0xdc, 0x85, 0x36, 0x7f, 0xfa, 0xff, 0xae, 0x33, 0x67, 0x67, 0xc9, 0xab, 0x04, 0xbf, 0x00, 0xec,
	0x1e, 0x19, 0x5f, 0xdc, 0xb1, 0xff, 0xbf, 0x51, 0xc8, 0xd9, 0x1d, 0x94, 0xe8, 0xce, 0xf1, 0xf0,
	0xbc, 0xef, 0xff, 0xeb, 0x8b, 0x81, 0xff, 0xbc, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xeb, 0x2a,
	0x89, 0x63, 0x16, 0x04, 0x00, 0x00,
}
