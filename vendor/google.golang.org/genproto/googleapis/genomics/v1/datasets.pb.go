// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/genomics/v1/datasets.proto

package genomics // import "google.golang.org/genproto/googleapis/genomics/v1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import v1 "google.golang.org/genproto/googleapis/iam/v1"
import field_mask "google.golang.org/genproto/protobuf/field_mask"

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

// A Dataset is a collection of genomic data.
//
// For more genomics resource definitions, see [Fundamentals of Google
// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
type Dataset struct {
	// The server-generated dataset ID, unique across all datasets.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The Google Cloud project ID that this dataset belongs to.
	ProjectId string `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// The dataset name.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// The time this dataset was created, in seconds from the epoch.
	CreateTime           *timestamp.Timestamp `protobuf:"bytes,4,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Dataset) Reset()         { *m = Dataset{} }
func (m *Dataset) String() string { return proto.CompactTextString(m) }
func (*Dataset) ProtoMessage()    {}
func (*Dataset) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{0}
}
func (m *Dataset) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dataset.Unmarshal(m, b)
}
func (m *Dataset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dataset.Marshal(b, m, deterministic)
}
func (m *Dataset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dataset.Merge(m, src)
}
func (m *Dataset) XXX_Size() int {
	return xxx_messageInfo_Dataset.Size(m)
}
func (m *Dataset) XXX_DiscardUnknown() {
	xxx_messageInfo_Dataset.DiscardUnknown(m)
}

var xxx_messageInfo_Dataset proto.InternalMessageInfo

func (m *Dataset) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Dataset) GetProjectId() string {
	if m != nil {
		return m.ProjectId
	}
	return ""
}

func (m *Dataset) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Dataset) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

// The dataset list request.
type ListDatasetsRequest struct {
	// Required. The Google Cloud project ID to list datasets for.
	ProjectId string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// The maximum number of results to return in a single page. If unspecified,
	// defaults to 50. The maximum value is 1024.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The continuation token, which is used to page through large result sets.
	// To get the next page of results, set this parameter to the value of
	// `nextPageToken` from the previous response.
	PageToken            string   `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListDatasetsRequest) Reset()         { *m = ListDatasetsRequest{} }
func (m *ListDatasetsRequest) String() string { return proto.CompactTextString(m) }
func (*ListDatasetsRequest) ProtoMessage()    {}
func (*ListDatasetsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{1}
}
func (m *ListDatasetsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDatasetsRequest.Unmarshal(m, b)
}
func (m *ListDatasetsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDatasetsRequest.Marshal(b, m, deterministic)
}
func (m *ListDatasetsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDatasetsRequest.Merge(m, src)
}
func (m *ListDatasetsRequest) XXX_Size() int {
	return xxx_messageInfo_ListDatasetsRequest.Size(m)
}
func (m *ListDatasetsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDatasetsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListDatasetsRequest proto.InternalMessageInfo

func (m *ListDatasetsRequest) GetProjectId() string {
	if m != nil {
		return m.ProjectId
	}
	return ""
}

func (m *ListDatasetsRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListDatasetsRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

// The dataset list response.
type ListDatasetsResponse struct {
	// The list of matching Datasets.
	Datasets []*Dataset `protobuf:"bytes,1,rep,name=datasets,proto3" json:"datasets,omitempty"`
	// The continuation token, which is used to page through large result sets.
	// Provide this value in a subsequent request to return the next page of
	// results. This field will be empty if there aren't any additional results.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListDatasetsResponse) Reset()         { *m = ListDatasetsResponse{} }
func (m *ListDatasetsResponse) String() string { return proto.CompactTextString(m) }
func (*ListDatasetsResponse) ProtoMessage()    {}
func (*ListDatasetsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{2}
}
func (m *ListDatasetsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDatasetsResponse.Unmarshal(m, b)
}
func (m *ListDatasetsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDatasetsResponse.Marshal(b, m, deterministic)
}
func (m *ListDatasetsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDatasetsResponse.Merge(m, src)
}
func (m *ListDatasetsResponse) XXX_Size() int {
	return xxx_messageInfo_ListDatasetsResponse.Size(m)
}
func (m *ListDatasetsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDatasetsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListDatasetsResponse proto.InternalMessageInfo

func (m *ListDatasetsResponse) GetDatasets() []*Dataset {
	if m != nil {
		return m.Datasets
	}
	return nil
}

func (m *ListDatasetsResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type CreateDatasetRequest struct {
	// The dataset to be created. Must contain projectId and name.
	Dataset              *Dataset `protobuf:"bytes,1,opt,name=dataset,proto3" json:"dataset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateDatasetRequest) Reset()         { *m = CreateDatasetRequest{} }
func (m *CreateDatasetRequest) String() string { return proto.CompactTextString(m) }
func (*CreateDatasetRequest) ProtoMessage()    {}
func (*CreateDatasetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{3}
}
func (m *CreateDatasetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateDatasetRequest.Unmarshal(m, b)
}
func (m *CreateDatasetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateDatasetRequest.Marshal(b, m, deterministic)
}
func (m *CreateDatasetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateDatasetRequest.Merge(m, src)
}
func (m *CreateDatasetRequest) XXX_Size() int {
	return xxx_messageInfo_CreateDatasetRequest.Size(m)
}
func (m *CreateDatasetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateDatasetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateDatasetRequest proto.InternalMessageInfo

func (m *CreateDatasetRequest) GetDataset() *Dataset {
	if m != nil {
		return m.Dataset
	}
	return nil
}

type UpdateDatasetRequest struct {
	// The ID of the dataset to be updated.
	DatasetId string `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	// The new dataset data.
	Dataset *Dataset `protobuf:"bytes,2,opt,name=dataset,proto3" json:"dataset,omitempty"`
	// An optional mask specifying which fields to update. At this time, the only
	// mutable field is [name][google.genomics.v1.Dataset.name]. The only
	// acceptable value is "name". If unspecified, all mutable fields will be
	// updated.
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,3,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateDatasetRequest) Reset()         { *m = UpdateDatasetRequest{} }
func (m *UpdateDatasetRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateDatasetRequest) ProtoMessage()    {}
func (*UpdateDatasetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{4}
}
func (m *UpdateDatasetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateDatasetRequest.Unmarshal(m, b)
}
func (m *UpdateDatasetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateDatasetRequest.Marshal(b, m, deterministic)
}
func (m *UpdateDatasetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateDatasetRequest.Merge(m, src)
}
func (m *UpdateDatasetRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateDatasetRequest.Size(m)
}
func (m *UpdateDatasetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateDatasetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateDatasetRequest proto.InternalMessageInfo

func (m *UpdateDatasetRequest) GetDatasetId() string {
	if m != nil {
		return m.DatasetId
	}
	return ""
}

func (m *UpdateDatasetRequest) GetDataset() *Dataset {
	if m != nil {
		return m.Dataset
	}
	return nil
}

func (m *UpdateDatasetRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type DeleteDatasetRequest struct {
	// The ID of the dataset to be deleted.
	DatasetId            string   `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteDatasetRequest) Reset()         { *m = DeleteDatasetRequest{} }
func (m *DeleteDatasetRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteDatasetRequest) ProtoMessage()    {}
func (*DeleteDatasetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{5}
}
func (m *DeleteDatasetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteDatasetRequest.Unmarshal(m, b)
}
func (m *DeleteDatasetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteDatasetRequest.Marshal(b, m, deterministic)
}
func (m *DeleteDatasetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteDatasetRequest.Merge(m, src)
}
func (m *DeleteDatasetRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteDatasetRequest.Size(m)
}
func (m *DeleteDatasetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteDatasetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteDatasetRequest proto.InternalMessageInfo

func (m *DeleteDatasetRequest) GetDatasetId() string {
	if m != nil {
		return m.DatasetId
	}
	return ""
}

type UndeleteDatasetRequest struct {
	// The ID of the dataset to be undeleted.
	DatasetId            string   `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UndeleteDatasetRequest) Reset()         { *m = UndeleteDatasetRequest{} }
func (m *UndeleteDatasetRequest) String() string { return proto.CompactTextString(m) }
func (*UndeleteDatasetRequest) ProtoMessage()    {}
func (*UndeleteDatasetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{6}
}
func (m *UndeleteDatasetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UndeleteDatasetRequest.Unmarshal(m, b)
}
func (m *UndeleteDatasetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UndeleteDatasetRequest.Marshal(b, m, deterministic)
}
func (m *UndeleteDatasetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UndeleteDatasetRequest.Merge(m, src)
}
func (m *UndeleteDatasetRequest) XXX_Size() int {
	return xxx_messageInfo_UndeleteDatasetRequest.Size(m)
}
func (m *UndeleteDatasetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UndeleteDatasetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UndeleteDatasetRequest proto.InternalMessageInfo

func (m *UndeleteDatasetRequest) GetDatasetId() string {
	if m != nil {
		return m.DatasetId
	}
	return ""
}

type GetDatasetRequest struct {
	// The ID of the dataset.
	DatasetId            string   `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDatasetRequest) Reset()         { *m = GetDatasetRequest{} }
func (m *GetDatasetRequest) String() string { return proto.CompactTextString(m) }
func (*GetDatasetRequest) ProtoMessage()    {}
func (*GetDatasetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddd0efa223187e29, []int{7}
}
func (m *GetDatasetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDatasetRequest.Unmarshal(m, b)
}
func (m *GetDatasetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDatasetRequest.Marshal(b, m, deterministic)
}
func (m *GetDatasetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDatasetRequest.Merge(m, src)
}
func (m *GetDatasetRequest) XXX_Size() int {
	return xxx_messageInfo_GetDatasetRequest.Size(m)
}
func (m *GetDatasetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDatasetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDatasetRequest proto.InternalMessageInfo

func (m *GetDatasetRequest) GetDatasetId() string {
	if m != nil {
		return m.DatasetId
	}
	return ""
}

func init() {
	proto.RegisterType((*Dataset)(nil), "google.genomics.v1.Dataset")
	proto.RegisterType((*ListDatasetsRequest)(nil), "google.genomics.v1.ListDatasetsRequest")
	proto.RegisterType((*ListDatasetsResponse)(nil), "google.genomics.v1.ListDatasetsResponse")
	proto.RegisterType((*CreateDatasetRequest)(nil), "google.genomics.v1.CreateDatasetRequest")
	proto.RegisterType((*UpdateDatasetRequest)(nil), "google.genomics.v1.UpdateDatasetRequest")
	proto.RegisterType((*DeleteDatasetRequest)(nil), "google.genomics.v1.DeleteDatasetRequest")
	proto.RegisterType((*UndeleteDatasetRequest)(nil), "google.genomics.v1.UndeleteDatasetRequest")
	proto.RegisterType((*GetDatasetRequest)(nil), "google.genomics.v1.GetDatasetRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DatasetServiceV1Client is the client API for DatasetServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DatasetServiceV1Client interface {
	// Lists datasets within a project.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	ListDatasets(ctx context.Context, in *ListDatasetsRequest, opts ...grpc.CallOption) (*ListDatasetsResponse, error)
	// Creates a new dataset.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	CreateDataset(ctx context.Context, in *CreateDatasetRequest, opts ...grpc.CallOption) (*Dataset, error)
	// Gets a dataset by ID.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	GetDataset(ctx context.Context, in *GetDatasetRequest, opts ...grpc.CallOption) (*Dataset, error)
	// Updates a dataset.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// This method supports patch semantics.
	UpdateDataset(ctx context.Context, in *UpdateDatasetRequest, opts ...grpc.CallOption) (*Dataset, error)
	// Deletes a dataset and all of its contents (all read group sets,
	// reference sets, variant sets, call sets, annotation sets, etc.)
	// This is reversible (up to one week after the deletion) via
	// the
	// [datasets.undelete][google.genomics.v1.DatasetServiceV1.UndeleteDataset]
	// operation.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	DeleteDataset(ctx context.Context, in *DeleteDatasetRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Undeletes a dataset by restoring a dataset which was deleted via this API.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// This operation is only possible for a week after the deletion occurred.
	UndeleteDataset(ctx context.Context, in *UndeleteDatasetRequest, opts ...grpc.CallOption) (*Dataset, error)
	// Sets the access control policy on the specified dataset. Replaces any
	// existing policy.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// See <a href="/iam/docs/managing-policies#setting_a_policy">Setting a
	// Policy</a> for more information.
	SetIamPolicy(ctx context.Context, in *v1.SetIamPolicyRequest, opts ...grpc.CallOption) (*v1.Policy, error)
	// Gets the access control policy for the dataset. This is empty if the
	// policy or resource does not exist.
	//
	// See <a href="/iam/docs/managing-policies#getting_a_policy">Getting a
	// Policy</a> for more information.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	GetIamPolicy(ctx context.Context, in *v1.GetIamPolicyRequest, opts ...grpc.CallOption) (*v1.Policy, error)
	// Returns permissions that a caller has on the specified resource.
	// See <a href="/iam/docs/managing-policies#testing_permissions">Testing
	// Permissions</a> for more information.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	TestIamPermissions(ctx context.Context, in *v1.TestIamPermissionsRequest, opts ...grpc.CallOption) (*v1.TestIamPermissionsResponse, error)
}

type datasetServiceV1Client struct {
	cc *grpc.ClientConn
}

func NewDatasetServiceV1Client(cc *grpc.ClientConn) DatasetServiceV1Client {
	return &datasetServiceV1Client{cc}
}

func (c *datasetServiceV1Client) ListDatasets(ctx context.Context, in *ListDatasetsRequest, opts ...grpc.CallOption) (*ListDatasetsResponse, error) {
	out := new(ListDatasetsResponse)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/ListDatasets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) CreateDataset(ctx context.Context, in *CreateDatasetRequest, opts ...grpc.CallOption) (*Dataset, error) {
	out := new(Dataset)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/CreateDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) GetDataset(ctx context.Context, in *GetDatasetRequest, opts ...grpc.CallOption) (*Dataset, error) {
	out := new(Dataset)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/GetDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) UpdateDataset(ctx context.Context, in *UpdateDatasetRequest, opts ...grpc.CallOption) (*Dataset, error) {
	out := new(Dataset)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/UpdateDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) DeleteDataset(ctx context.Context, in *DeleteDatasetRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/DeleteDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) UndeleteDataset(ctx context.Context, in *UndeleteDatasetRequest, opts ...grpc.CallOption) (*Dataset, error) {
	out := new(Dataset)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/UndeleteDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) SetIamPolicy(ctx context.Context, in *v1.SetIamPolicyRequest, opts ...grpc.CallOption) (*v1.Policy, error) {
	out := new(v1.Policy)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/SetIamPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) GetIamPolicy(ctx context.Context, in *v1.GetIamPolicyRequest, opts ...grpc.CallOption) (*v1.Policy, error) {
	out := new(v1.Policy)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/GetIamPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceV1Client) TestIamPermissions(ctx context.Context, in *v1.TestIamPermissionsRequest, opts ...grpc.CallOption) (*v1.TestIamPermissionsResponse, error) {
	out := new(v1.TestIamPermissionsResponse)
	err := c.cc.Invoke(ctx, "/google.genomics.v1.DatasetServiceV1/TestIamPermissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatasetServiceV1Server is the server API for DatasetServiceV1 service.
type DatasetServiceV1Server interface {
	// Lists datasets within a project.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	ListDatasets(context.Context, *ListDatasetsRequest) (*ListDatasetsResponse, error)
	// Creates a new dataset.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	CreateDataset(context.Context, *CreateDatasetRequest) (*Dataset, error)
	// Gets a dataset by ID.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	GetDataset(context.Context, *GetDatasetRequest) (*Dataset, error)
	// Updates a dataset.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// This method supports patch semantics.
	UpdateDataset(context.Context, *UpdateDatasetRequest) (*Dataset, error)
	// Deletes a dataset and all of its contents (all read group sets,
	// reference sets, variant sets, call sets, annotation sets, etc.)
	// This is reversible (up to one week after the deletion) via
	// the
	// [datasets.undelete][google.genomics.v1.DatasetServiceV1.UndeleteDataset]
	// operation.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	DeleteDataset(context.Context, *DeleteDatasetRequest) (*empty.Empty, error)
	// Undeletes a dataset by restoring a dataset which was deleted via this API.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// This operation is only possible for a week after the deletion occurred.
	UndeleteDataset(context.Context, *UndeleteDatasetRequest) (*Dataset, error)
	// Sets the access control policy on the specified dataset. Replaces any
	// existing policy.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	//
	// See <a href="/iam/docs/managing-policies#setting_a_policy">Setting a
	// Policy</a> for more information.
	SetIamPolicy(context.Context, *v1.SetIamPolicyRequest) (*v1.Policy, error)
	// Gets the access control policy for the dataset. This is empty if the
	// policy or resource does not exist.
	//
	// See <a href="/iam/docs/managing-policies#getting_a_policy">Getting a
	// Policy</a> for more information.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	GetIamPolicy(context.Context, *v1.GetIamPolicyRequest) (*v1.Policy, error)
	// Returns permissions that a caller has on the specified resource.
	// See <a href="/iam/docs/managing-policies#testing_permissions">Testing
	// Permissions</a> for more information.
	//
	// For the definitions of datasets and other genomics resources, see
	// [Fundamentals of Google
	// Genomics](https://cloud.google.com/genomics/fundamentals-of-google-genomics)
	TestIamPermissions(context.Context, *v1.TestIamPermissionsRequest) (*v1.TestIamPermissionsResponse, error)
}

func RegisterDatasetServiceV1Server(s *grpc.Server, srv DatasetServiceV1Server) {
	s.RegisterService(&_DatasetServiceV1_serviceDesc, srv)
}

func _DatasetServiceV1_ListDatasets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDatasetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).ListDatasets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/ListDatasets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).ListDatasets(ctx, req.(*ListDatasetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_CreateDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).CreateDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/CreateDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).CreateDataset(ctx, req.(*CreateDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_GetDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).GetDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/GetDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).GetDataset(ctx, req.(*GetDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_UpdateDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).UpdateDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/UpdateDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).UpdateDataset(ctx, req.(*UpdateDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_DeleteDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).DeleteDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/DeleteDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).DeleteDataset(ctx, req.(*DeleteDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_UndeleteDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UndeleteDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).UndeleteDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/UndeleteDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).UndeleteDataset(ctx, req.(*UndeleteDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_SetIamPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.SetIamPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).SetIamPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/SetIamPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).SetIamPolicy(ctx, req.(*v1.SetIamPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_GetIamPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.GetIamPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).GetIamPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/GetIamPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).GetIamPolicy(ctx, req.(*v1.GetIamPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetServiceV1_TestIamPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.TestIamPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceV1Server).TestIamPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.genomics.v1.DatasetServiceV1/TestIamPermissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceV1Server).TestIamPermissions(ctx, req.(*v1.TestIamPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DatasetServiceV1_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.genomics.v1.DatasetServiceV1",
	HandlerType: (*DatasetServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDatasets",
			Handler:    _DatasetServiceV1_ListDatasets_Handler,
		},
		{
			MethodName: "CreateDataset",
			Handler:    _DatasetServiceV1_CreateDataset_Handler,
		},
		{
			MethodName: "GetDataset",
			Handler:    _DatasetServiceV1_GetDataset_Handler,
		},
		{
			MethodName: "UpdateDataset",
			Handler:    _DatasetServiceV1_UpdateDataset_Handler,
		},
		{
			MethodName: "DeleteDataset",
			Handler:    _DatasetServiceV1_DeleteDataset_Handler,
		},
		{
			MethodName: "UndeleteDataset",
			Handler:    _DatasetServiceV1_UndeleteDataset_Handler,
		},
		{
			MethodName: "SetIamPolicy",
			Handler:    _DatasetServiceV1_SetIamPolicy_Handler,
		},
		{
			MethodName: "GetIamPolicy",
			Handler:    _DatasetServiceV1_GetIamPolicy_Handler,
		},
		{
			MethodName: "TestIamPermissions",
			Handler:    _DatasetServiceV1_TestIamPermissions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/genomics/v1/datasets.proto",
}

func init() { proto.RegisterFile("google/genomics/v1/datasets.proto", fileDescriptor_ddd0efa223187e29) }

var fileDescriptor_ddd0efa223187e29 = []byte{
	// 786 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xd1, 0x4e, 0x13, 0x4d,
	0x14, 0xce, 0x16, 0xfe, 0x1f, 0x7a, 0xa0, 0xa0, 0x63, 0xc5, 0xda, 0x8a, 0x96, 0x8d, 0x42, 0xad,
	0xba, 0x4d, 0x6b, 0x08, 0x49, 0x89, 0x37, 0x88, 0x12, 0x12, 0x49, 0x9a, 0x02, 0x5e, 0x78, 0xd3,
	0x0c, 0xdd, 0xa1, 0x8e, 0x74, 0x77, 0xd6, 0x9d, 0x29, 0x28, 0xc8, 0x0d, 0x77, 0x5c, 0xfb, 0x00,
	0x26, 0xde, 0xf9, 0x3c, 0xbe, 0x82, 0x0f, 0xe1, 0xa5, 0x99, 0xd9, 0xd9, 0x76, 0xdb, 0x2e, 0x05,
	0x8c, 0x77, 0xdb, 0x73, 0xbe, 0x73, 0xbe, 0xef, 0xcc, 0xf9, 0x76, 0xba, 0xb0, 0xd0, 0x62, 0xac,
	0xd5, 0x26, 0xa5, 0x16, 0x71, 0x99, 0x43, 0x9b, 0xbc, 0x74, 0x58, 0x2e, 0xd9, 0x58, 0x60, 0x4e,
	0x04, 0xb7, 0x3c, 0x9f, 0x09, 0x86, 0x50, 0x00, 0xb1, 0x42, 0x88, 0x75, 0x58, 0xce, 0xde, 0xd3,
	0x65, 0xd8, 0xa3, 0x25, 0xec, 0xba, 0x4c, 0x60, 0x41, 0x99, 0xab, 0x2b, 0xb2, 0xf7, 0x75, 0x96,
	0x62, 0x47, 0xf6, 0xa3, 0xd8, 0x69, 0x78, 0xac, 0x4d, 0x9b, 0x9f, 0x75, 0x3e, 0xdb, 0x9f, 0xef,
	0xcb, 0xe5, 0x74, 0x4e, 0xfd, 0xda, 0xeb, 0xec, 0x97, 0x88, 0xe3, 0x89, 0x30, 0x99, 0x1f, 0x4c,
	0xee, 0x53, 0xd2, 0xb6, 0x1b, 0x0e, 0xe6, 0x07, 0x1a, 0xf1, 0x60, 0x10, 0x21, 0xa8, 0x43, 0xb8,
	0xc0, 0x8e, 0x17, 0x00, 0xcc, 0x73, 0x03, 0x26, 0xd6, 0x83, 0x01, 0xd1, 0x0c, 0x24, 0xa8, 0x9d,
	0x31, 0xf2, 0x46, 0x21, 0x59, 0x4f, 0x50, 0x1b, 0xcd, 0x03, 0x78, 0x3e, 0xfb, 0x40, 0x9a, 0xa2,
	0x41, 0xed, 0x4c, 0x42, 0xc5, 0x93, 0x3a, 0xb2, 0x69, 0x23, 0x04, 0xe3, 0x2e, 0x76, 0x48, 0x66,
	0x4c, 0x25, 0xd4, 0x33, 0x5a, 0x85, 0xa9, 0xa6, 0x4f, 0xb0, 0x20, 0x0d, 0x49, 0x94, 0x19, 0xcf,
	0x1b, 0x85, 0xa9, 0x4a, 0xd6, 0xd2, 0x47, 0x16, 0xaa, 0xb0, 0x76, 0x42, 0x15, 0x75, 0x08, 0xe0,
	0x32, 0x60, 0x7a, 0x70, 0xeb, 0x0d, 0xe5, 0x42, 0xcb, 0xe1, 0x75, 0xf2, 0xb1, 0x43, 0xb8, 0x18,
	0x90, 0x61, 0x0c, 0xca, 0xc8, 0x41, 0xd2, 0xc3, 0x2d, 0xd2, 0xe0, 0xf4, 0x98, 0x28, 0x91, 0xff,
	0xd5, 0x27, 0x65, 0x60, 0x9b, 0x1e, 0x13, 0x55, 0x2b, 0x93, 0x82, 0x1d, 0x10, 0x57, 0x2b, 0x55,
	0xf0, 0x1d, 0x19, 0x30, 0x8f, 0x20, 0xdd, 0xcf, 0xc8, 0x3d, 0xe6, 0x72, 0x82, 0x56, 0x60, 0x32,
	0xdc, 0x7a, 0xc6, 0xc8, 0x8f, 0x15, 0xa6, 0x2a, 0x39, 0x6b, 0x78, 0xed, 0x96, 0xae, 0xab, 0x77,
	0xc1, 0x68, 0x11, 0x66, 0x5d, 0xf2, 0x49, 0x34, 0x22, 0xa4, 0xc1, 0xb9, 0xa5, 0x64, 0xb8, 0xd6,
	0x25, 0xde, 0x82, 0xf4, 0x4b, 0x35, 0x78, 0xd8, 0x42, 0xcf, 0xba, 0x0c, 0x13, 0xba, 0x97, 0x1a,
	0xf4, 0x12, 0xde, 0x10, 0x6b, 0xfe, 0x30, 0x20, 0xbd, 0xeb, 0xd9, 0xc3, 0xfd, 0xe6, 0x01, 0x34,
	0x26, 0x72, 0x76, 0x3a, 0xb2, 0x69, 0x47, 0xe9, 0x12, 0x57, 0xa7, 0x93, 0x5b, 0xee, 0x28, 0x36,
	0x65, 0x35, 0x75, 0xac, 0x71, 0x5b, 0x7e, 0x2d, 0xdd, 0xb8, 0x85, 0xf9, 0x41, 0x1d, 0x02, 0xb8,
	0x7c, 0x36, 0x97, 0x21, 0xbd, 0x4e, 0xda, 0xe4, 0x9a, 0x52, 0xcd, 0x15, 0x98, 0xdb, 0x75, 0xed,
	0xbf, 0x28, 0xac, 0xc0, 0xcd, 0x0d, 0x22, 0xae, 0x55, 0x53, 0xf9, 0x96, 0x84, 0x1b, 0xba, 0x62,
	0x9b, 0xf8, 0x87, 0xb4, 0x49, 0xde, 0x96, 0xd1, 0x11, 0x4c, 0x47, 0xcd, 0x82, 0x96, 0xe2, 0xce,
	0x2a, 0xc6, 0xc0, 0xd9, 0xc2, 0xe5, 0xc0, 0xc0, 0x77, 0x66, 0xfa, 0xec, 0xe7, 0xaf, 0xaf, 0x89,
	0x19, 0x34, 0x1d, 0xbd, 0x77, 0x50, 0x07, 0x52, 0x7d, 0x66, 0x41, 0xb1, 0x0d, 0xe3, 0xfc, 0x94,
	0x1d, 0xb5, 0x4f, 0x73, 0x5e, 0xb1, 0xdd, 0x31, 0xfb, 0xd8, 0xaa, 0xdd, 0x2d, 0x73, 0x80, 0xde,
	0xc1, 0xa1, 0x47, 0x71, 0x9d, 0x86, 0x0e, 0x76, 0x34, 0xe1, 0x82, 0x22, 0xcc, 0xa1, 0xbb, 0x51,
	0xc2, 0xd2, 0x49, 0x6f, 0x13, 0xa7, 0xe8, 0xcc, 0x80, 0x54, 0x9f, 0x93, 0xe3, 0x87, 0x8d, 0x33,
	0xfb, 0x68, 0xee, 0xa2, 0xe2, 0x7e, 0x58, 0xb9, 0x98, 0xbb, 0x37, 0xb9, 0x80, 0x54, 0x9f, 0x45,
	0xe3, 0x35, 0xc4, 0xb9, 0x38, 0x3b, 0x37, 0xf4, 0x16, 0xbc, 0x92, 0x17, 0x76, 0x38, 0x7a, 0x71,
	0xc4, 0xe8, 0xe7, 0x06, 0xcc, 0x0e, 0x58, 0x1c, 0x15, 0x63, 0x87, 0x8f, 0x7d, 0x0f, 0x46, 0x8f,
	0xff, 0x4c, 0xf1, 0x2f, 0x99, 0xe6, 0xc5, 0xe3, 0x77, 0x74, 0xdb, 0xaa, 0x51, 0x44, 0x5f, 0x60,
	0x7a, 0x9b, 0x88, 0x4d, 0xec, 0xd4, 0xd4, 0x9f, 0x11, 0x32, 0xc3, 0xde, 0x14, 0x3b, 0xb2, 0x6d,
	0x34, 0x19, 0xf2, 0xdf, 0x1e, 0xc0, 0x04, 0x59, 0xb3, 0xac, 0x98, 0x9f, 0x98, 0x8b, 0x92, 0xf9,
	0xc4, 0x27, 0x9c, 0x75, 0xfc, 0x26, 0x79, 0xd1, 0xd5, 0x50, 0x3c, 0xad, 0xf2, 0x48, 0x37, 0xcd,
	0xbe, 0x31, 0x8a, 0x7d, 0xe3, 0x9f, 0xb2, 0xb7, 0x06, 0xd8, 0xbf, 0x1b, 0x80, 0x76, 0x08, 0x57,
	0x41, 0xe2, 0x3b, 0x94, 0x73, 0xf9, 0x5f, 0xde, 0xf3, 0x80, 0x26, 0x18, 0x86, 0x84, 0x52, 0x1e,
	0x5f, 0x01, 0xa9, 0x5f, 0xf8, 0x15, 0x25, 0xaf, 0x6c, 0x3e, 0xbd, 0x58, 0x9e, 0x18, 0xaa, 0xae,
	0x1a, 0xc5, 0xb5, 0xf7, 0x30, 0xd7, 0x64, 0x4e, 0xcc, 0xc6, 0xd7, 0x52, 0xe1, 0xad, 0x52, 0x93,
	0x0e, 0xac, 0x19, 0xef, 0xaa, 0x21, 0x88, 0xb5, 0xb1, 0xdb, 0xb2, 0x98, 0xdf, 0x92, 0x9f, 0x37,
	0xca, 0x9f, 0xa5, 0x20, 0x85, 0x3d, 0xca, 0xa3, 0x9f, 0x3c, 0xab, 0xe1, 0xf3, 0x6f, 0xc3, 0xd8,
	0xfb, 0x5f, 0x21, 0x9f, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x87, 0x48, 0x07, 0xbb, 0x1b, 0x09,
	0x00, 0x00,
}
