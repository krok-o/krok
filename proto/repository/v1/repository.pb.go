// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: repository/v1/repository.proto

package repov1

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Url       string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Vcs       int32  `protobuf:"varint,4,opt,name=vcs,proto3" json:"vcs,omitempty"`
	UniqueUrl string `protobuf:"bytes,5,opt,name=unique_url,json=uniqueUrl,proto3" json:"unique_url,omitempty"`
}

func (x *Repository) Reset() {
	*x = Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repository) ProtoMessage() {}

func (x *Repository) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repository.ProtoReflect.Descriptor instead.
func (*Repository) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{0}
}

func (x *Repository) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Repository) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Repository) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Repository) GetVcs() int32 {
	if x != nil {
		return x.Vcs
	}
	return 0
}

func (x *Repository) GetUniqueUrl() string {
	if x != nil {
		return x.UniqueUrl
	}
	return ""
}

type Repositories struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Repository `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *Repositories) Reset() {
	*x = Repositories{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repositories) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repositories) ProtoMessage() {}

func (x *Repositories) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repositories.ProtoReflect.Descriptor instead.
func (*Repositories) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{1}
}

func (x *Repositories) GetItems() []*Repository {
	if x != nil {
		return x.Items
	}
	return nil
}

type CreateRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url  string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Vcs  int32  `protobuf:"varint,3,opt,name=vcs,proto3" json:"vcs,omitempty"`
}

func (x *CreateRepositoryRequest) Reset() {
	*x = CreateRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRepositoryRequest) ProtoMessage() {}

func (x *CreateRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRepositoryRequest.ProtoReflect.Descriptor instead.
func (*CreateRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRepositoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRepositoryRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateRepositoryRequest) GetVcs() int32 {
	if x != nil {
		return x.Vcs
	}
	return 0
}

type GetRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *wrappers.Int32Value `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRepositoryRequest) Reset() {
	*x = GetRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRepositoryRequest) ProtoMessage() {}

func (x *GetRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRepositoryRequest.ProtoReflect.Descriptor instead.
func (*GetRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{3}
}

func (x *GetRepositoryRequest) GetId() *wrappers.Int32Value {
	if x != nil {
		return x.Id
	}
	return nil
}

type ListRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Vcs  int32  `protobuf:"varint,2,opt,name=vcs,proto3" json:"vcs,omitempty"`
}

func (x *ListRepositoryRequest) Reset() {
	*x = ListRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRepositoryRequest) ProtoMessage() {}

func (x *ListRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRepositoryRequest.ProtoReflect.Descriptor instead.
func (*ListRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{4}
}

func (x *ListRepositoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListRepositoryRequest) GetVcs() int32 {
	if x != nil {
		return x.Vcs
	}
	return 0
}

type UpdateRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   *wrappers.Int32Value `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UpdateRepositoryRequest) Reset() {
	*x = UpdateRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRepositoryRequest) ProtoMessage() {}

func (x *UpdateRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRepositoryRequest.ProtoReflect.Descriptor instead.
func (*UpdateRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateRepositoryRequest) GetId() *wrappers.Int32Value {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateRepositoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DeleteRepositoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *wrappers.Int32Value `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteRepositoryRequest) Reset() {
	*x = DeleteRepositoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_v1_repository_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRepositoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRepositoryRequest) ProtoMessage() {}

func (x *DeleteRepositoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_v1_repository_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRepositoryRequest.ProtoReflect.Descriptor instead.
func (*DeleteRepositoryRequest) Descriptor() ([]byte, []int) {
	return file_repository_v1_repository_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteRepositoryRequest) GetId() *wrappers.Int32Value {
	if x != nil {
		return x.Id
	}
	return nil
}

var File_repository_v1_repository_proto protoreflect.FileDescriptor

var file_repository_v1_repository_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f,
	0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x0a,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x12, 0x10, 0x0a, 0x03, 0x76, 0x63, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76,
	0x63, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x55, 0x72,
	0x6c, 0x22, 0x3f, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x22, 0x51, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x63, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x76, 0x63, 0x73, 0x22, 0x43, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33,
	0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3d, 0x0a, 0x15, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x63, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76, 0x63, 0x73, 0x22, 0x5a, 0x0a, 0x17, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x46, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49,
	0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x69, 0x64, 0x32, 0xc7, 0x03,
	0x0a, 0x11, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x26, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x10,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x26, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x23, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x72, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x24, 0x2e, 0x72,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x22,
	0x00, 0x12, 0x54, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x26, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x6f, 0x6b, 0x2d, 0x6f, 0x2f, 0x6b, 0x72, 0x6f,
	0x6b, 0x2d, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x65, 0x70, 0x6f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_repository_v1_repository_proto_rawDescOnce sync.Once
	file_repository_v1_repository_proto_rawDescData = file_repository_v1_repository_proto_rawDesc
)

func file_repository_v1_repository_proto_rawDescGZIP() []byte {
	file_repository_v1_repository_proto_rawDescOnce.Do(func() {
		file_repository_v1_repository_proto_rawDescData = protoimpl.X.CompressGZIP(file_repository_v1_repository_proto_rawDescData)
	})
	return file_repository_v1_repository_proto_rawDescData
}

var file_repository_v1_repository_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_repository_v1_repository_proto_goTypes = []interface{}{
	(*Repository)(nil),              // 0: repository.v1.Repository
	(*Repositories)(nil),            // 1: repository.v1.Repositories
	(*CreateRepositoryRequest)(nil), // 2: repository.v1.CreateRepositoryRequest
	(*GetRepositoryRequest)(nil),    // 3: repository.v1.GetRepositoryRequest
	(*ListRepositoryRequest)(nil),   // 4: repository.v1.ListRepositoryRequest
	(*UpdateRepositoryRequest)(nil), // 5: repository.v1.UpdateRepositoryRequest
	(*DeleteRepositoryRequest)(nil), // 6: repository.v1.DeleteRepositoryRequest
	(*wrappers.Int32Value)(nil),     // 7: google.protobuf.Int32Value
	(*empty.Empty)(nil),             // 8: google.protobuf.Empty
}
var file_repository_v1_repository_proto_depIdxs = []int32{
	0, // 0: repository.v1.Repositories.items:type_name -> repository.v1.Repository
	7, // 1: repository.v1.GetRepositoryRequest.id:type_name -> google.protobuf.Int32Value
	7, // 2: repository.v1.UpdateRepositoryRequest.id:type_name -> google.protobuf.Int32Value
	7, // 3: repository.v1.DeleteRepositoryRequest.id:type_name -> google.protobuf.Int32Value
	2, // 4: repository.v1.RepositoryService.CreateRepository:input_type -> repository.v1.CreateRepositoryRequest
	5, // 5: repository.v1.RepositoryService.UpdateRepository:input_type -> repository.v1.UpdateRepositoryRequest
	3, // 6: repository.v1.RepositoryService.GetRepository:input_type -> repository.v1.GetRepositoryRequest
	4, // 7: repository.v1.RepositoryService.ListRepositories:input_type -> repository.v1.ListRepositoryRequest
	6, // 8: repository.v1.RepositoryService.DeleteRepository:input_type -> repository.v1.DeleteRepositoryRequest
	0, // 9: repository.v1.RepositoryService.CreateRepository:output_type -> repository.v1.Repository
	0, // 10: repository.v1.RepositoryService.UpdateRepository:output_type -> repository.v1.Repository
	0, // 11: repository.v1.RepositoryService.GetRepository:output_type -> repository.v1.Repository
	1, // 12: repository.v1.RepositoryService.ListRepositories:output_type -> repository.v1.Repositories
	8, // 13: repository.v1.RepositoryService.DeleteRepository:output_type -> google.protobuf.Empty
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_repository_v1_repository_proto_init() }
func file_repository_v1_repository_proto_init() {
	if File_repository_v1_repository_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_repository_v1_repository_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repository); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repositories); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRepositoryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRepositoryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRepositoryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRepositoryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_v1_repository_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRepositoryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_repository_v1_repository_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_repository_v1_repository_proto_goTypes,
		DependencyIndexes: file_repository_v1_repository_proto_depIdxs,
		MessageInfos:      file_repository_v1_repository_proto_msgTypes,
	}.Build()
	File_repository_v1_repository_proto = out.File
	file_repository_v1_repository_proto_rawDesc = nil
	file_repository_v1_repository_proto_goTypes = nil
	file_repository_v1_repository_proto_depIdxs = nil
}
