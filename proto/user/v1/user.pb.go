// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: user/v1/user.proto

package userv1

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type APIKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	UserId    int32                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	KeyId     string               `protobuf:"bytes,4,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	KeySecret string               `protobuf:"bytes,5,opt,name=key_secret,json=keySecret,proto3" json:"key_secret,omitempty"`
	Ttl       *timestamp.Timestamp `protobuf:"bytes,6,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *APIKey) Reset() {
	*x = APIKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APIKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APIKey) ProtoMessage() {}

func (x *APIKey) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APIKey.ProtoReflect.Descriptor instead.
func (*APIKey) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *APIKey) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *APIKey) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *APIKey) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *APIKey) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *APIKey) GetKeySecret() string {
	if x != nil {
		return x.KeySecret
	}
	return ""
}

func (x *APIKey) GetTtl() *timestamp.Timestamp {
	if x != nil {
		return x.Ttl
	}
	return nil
}

type APIKeys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*APIKey `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *APIKeys) Reset() {
	*x = APIKeys{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APIKeys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APIKeys) ProtoMessage() {}

func (x *APIKeys) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APIKeys.ProtoReflect.Descriptor instead.
func (*APIKeys) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{1}
}

func (x *APIKeys) GetItems() []*APIKey {
	if x != nil {
		return x.Items
	}
	return nil
}

type CreateAPIKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId *wrappers.Int32Value `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name   string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *CreateAPIKeyRequest) Reset() {
	*x = CreateAPIKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAPIKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAPIKeyRequest) ProtoMessage() {}

func (x *CreateAPIKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAPIKeyRequest.ProtoReflect.Descriptor instead.
func (*CreateAPIKeyRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAPIKeyRequest) GetUserId() *wrappers.Int32Value {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *CreateAPIKeyRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetAPIKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     *wrappers.Int32Value `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId *wrappers.Int32Value `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetAPIKeyRequest) Reset() {
	*x = GetAPIKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAPIKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAPIKeyRequest) ProtoMessage() {}

func (x *GetAPIKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAPIKeyRequest.ProtoReflect.Descriptor instead.
func (*GetAPIKeyRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetAPIKeyRequest) GetId() *wrappers.Int32Value {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *GetAPIKeyRequest) GetUserId() *wrappers.Int32Value {
	if x != nil {
		return x.UserId
	}
	return nil
}

type ListAPIKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId *wrappers.Int32Value `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ListAPIKeyRequest) Reset() {
	*x = ListAPIKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAPIKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAPIKeyRequest) ProtoMessage() {}

func (x *ListAPIKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAPIKeyRequest.ProtoReflect.Descriptor instead.
func (*ListAPIKeyRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{4}
}

func (x *ListAPIKeyRequest) GetUserId() *wrappers.Int32Value {
	if x != nil {
		return x.UserId
	}
	return nil
}

type DeleteAPIKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     *wrappers.Int32Value `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId *wrappers.Int32Value `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DeleteAPIKeyRequest) Reset() {
	*x = DeleteAPIKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAPIKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAPIKeyRequest) ProtoMessage() {}

func (x *DeleteAPIKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAPIKeyRequest.ProtoReflect.Descriptor instead.
func (*DeleteAPIKeyRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteAPIKeyRequest) GetId() *wrappers.Int32Value {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DeleteAPIKeyRequest) GetUserId() *wrappers.Int32Value {
	if x != nil {
		return x.UserId
	}
	return nil
}

var File_user_v1_user_proto protoreflect.FileDescriptor

var file_user_v1_user_proto_rawDesc = []byte{
	0x0a, 0x12, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x01, 0x0a, 0x06, 0x41, 0x50,
	0x49, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6b, 0x65, 0x79, 0x5f,
	0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6b, 0x65,
	0x79, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x03, 0x74, 0x74, 0x6c, 0x22, 0x30, 0x0a, 0x07, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x73,
	0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x5f, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x75, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41,
	0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x69, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74,
	0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x49, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x78, 0x0a, 0x13, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x02, 0x69, 0x64, 0x12, 0x34,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x32, 0x92, 0x02, 0x0a, 0x0d, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x50, 0x49, 0x4b, 0x65, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x39, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0b, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x1a, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x73, 0x22, 0x00, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x6f, 0x6b, 0x2d, 0x6f, 0x2f, 0x6b,
	0x72, 0x6f, 0x6b, 0x2d, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x76, 0x31, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_user_v1_user_proto_rawDescOnce sync.Once
	file_user_v1_user_proto_rawDescData = file_user_v1_user_proto_rawDesc
)

func file_user_v1_user_proto_rawDescGZIP() []byte {
	file_user_v1_user_proto_rawDescOnce.Do(func() {
		file_user_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_v1_user_proto_rawDescData)
	})
	return file_user_v1_user_proto_rawDescData
}

var file_user_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_user_v1_user_proto_goTypes = []interface{}{
	(*APIKey)(nil),              // 0: user.v1.APIKey
	(*APIKeys)(nil),             // 1: user.v1.APIKeys
	(*CreateAPIKeyRequest)(nil), // 2: user.v1.CreateAPIKeyRequest
	(*GetAPIKeyRequest)(nil),    // 3: user.v1.GetAPIKeyRequest
	(*ListAPIKeyRequest)(nil),   // 4: user.v1.ListAPIKeyRequest
	(*DeleteAPIKeyRequest)(nil), // 5: user.v1.DeleteAPIKeyRequest
	(*timestamp.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*wrappers.Int32Value)(nil), // 7: google.protobuf.Int32Value
	(*empty.Empty)(nil),         // 8: google.protobuf.Empty
}
var file_user_v1_user_proto_depIdxs = []int32{
	6,  // 0: user.v1.APIKey.ttl:type_name -> google.protobuf.Timestamp
	0,  // 1: user.v1.APIKeys.items:type_name -> user.v1.APIKey
	7,  // 2: user.v1.CreateAPIKeyRequest.user_id:type_name -> google.protobuf.Int32Value
	7,  // 3: user.v1.GetAPIKeyRequest.id:type_name -> google.protobuf.Int32Value
	7,  // 4: user.v1.GetAPIKeyRequest.user_id:type_name -> google.protobuf.Int32Value
	7,  // 5: user.v1.ListAPIKeyRequest.user_id:type_name -> google.protobuf.Int32Value
	7,  // 6: user.v1.DeleteAPIKeyRequest.id:type_name -> google.protobuf.Int32Value
	7,  // 7: user.v1.DeleteAPIKeyRequest.user_id:type_name -> google.protobuf.Int32Value
	2,  // 8: user.v1.APIKeyService.CreateAPIKey:input_type -> user.v1.CreateAPIKeyRequest
	5,  // 9: user.v1.APIKeyService.DeleteAPIKey:input_type -> user.v1.DeleteAPIKeyRequest
	3,  // 10: user.v1.APIKeyService.GetAPIKey:input_type -> user.v1.GetAPIKeyRequest
	4,  // 11: user.v1.APIKeyService.ListAPIKeys:input_type -> user.v1.ListAPIKeyRequest
	0,  // 12: user.v1.APIKeyService.CreateAPIKey:output_type -> user.v1.APIKey
	8,  // 13: user.v1.APIKeyService.DeleteAPIKey:output_type -> google.protobuf.Empty
	0,  // 14: user.v1.APIKeyService.GetAPIKey:output_type -> user.v1.APIKey
	1,  // 15: user.v1.APIKeyService.ListAPIKeys:output_type -> user.v1.APIKeys
	12, // [12:16] is the sub-list for method output_type
	8,  // [8:12] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_user_v1_user_proto_init() }
func file_user_v1_user_proto_init() {
	if File_user_v1_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_v1_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APIKey); i {
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
		file_user_v1_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APIKeys); i {
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
		file_user_v1_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAPIKeyRequest); i {
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
		file_user_v1_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAPIKeyRequest); i {
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
		file_user_v1_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAPIKeyRequest); i {
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
		file_user_v1_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAPIKeyRequest); i {
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
			RawDescriptor: file_user_v1_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_v1_user_proto_goTypes,
		DependencyIndexes: file_user_v1_user_proto_depIdxs,
		MessageInfos:      file_user_v1_user_proto_msgTypes,
	}.Build()
	File_user_v1_user_proto = out.File
	file_user_v1_user_proto_rawDesc = nil
	file_user_v1_user_proto_goTypes = nil
	file_user_v1_user_proto_depIdxs = nil
}
