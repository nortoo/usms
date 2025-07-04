// This file defines some common data models,
// which are useful in a lot of scenarios.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: pkg/proto/common/v1/types/shared.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TimeModel is the shared model for time fields.
type TimeModel struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CreatedAt     int64                  `protobuf:"varint,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimeModel) Reset() {
	*x = TimeModel{}
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeModel) ProtoMessage() {}

func (x *TimeModel) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeModel.ProtoReflect.Descriptor instead.
func (*TimeModel) Descriptor() ([]byte, []int) {
	return file_pkg_proto_common_v1_types_shared_proto_rawDescGZIP(), []int{0}
}

func (x *TimeModel) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *TimeModel) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// Pagination is the essential query parameter.
type Pagination struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_pkg_proto_common_v1_types_shared_proto_rawDescGZIP(), []int{1}
}

func (x *Pagination) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *Pagination) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

// PaginationResp contains the common fields when returning a list.
type PaginationResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Total         int64                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaginationResp) Reset() {
	*x = PaginationResp{}
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaginationResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationResp) ProtoMessage() {}

func (x *PaginationResp) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_common_v1_types_shared_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationResp.ProtoReflect.Descriptor instead.
func (*PaginationResp) Descriptor() ([]byte, []int) {
	return file_pkg_proto_common_v1_types_shared_proto_rawDescGZIP(), []int{2}
}

func (x *PaginationResp) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *PaginationResp) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_pkg_proto_common_v1_types_shared_proto protoreflect.FileDescriptor

const file_pkg_proto_common_v1_types_shared_proto_rawDesc = "" +
	"\n" +
	"&pkg/proto/common/v1/types/shared.proto\x12\x1bnortoo.usms.common.v1.types\"I\n" +
	"\tTimeModel\x12\x1d\n" +
	"\n" +
	"created_at\x18\x01 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x02 \x01(\x03R\tupdatedAt\"=\n" +
	"\n" +
	"Pagination\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x05R\bpageSize\":\n" +
	"\x0ePaginationResp\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x03R\x05totalB8Z6github.com/nortoo/usms/pkg/proto/common/v1/types;typesb\x06proto3"

var (
	file_pkg_proto_common_v1_types_shared_proto_rawDescOnce sync.Once
	file_pkg_proto_common_v1_types_shared_proto_rawDescData []byte
)

func file_pkg_proto_common_v1_types_shared_proto_rawDescGZIP() []byte {
	file_pkg_proto_common_v1_types_shared_proto_rawDescOnce.Do(func() {
		file_pkg_proto_common_v1_types_shared_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_common_v1_types_shared_proto_rawDesc), len(file_pkg_proto_common_v1_types_shared_proto_rawDesc)))
	})
	return file_pkg_proto_common_v1_types_shared_proto_rawDescData
}

var file_pkg_proto_common_v1_types_shared_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_proto_common_v1_types_shared_proto_goTypes = []any{
	(*TimeModel)(nil),      // 0: nortoo.usms.common.v1.types.TimeModel
	(*Pagination)(nil),     // 1: nortoo.usms.common.v1.types.Pagination
	(*PaginationResp)(nil), // 2: nortoo.usms.common.v1.types.PaginationResp
}
var file_pkg_proto_common_v1_types_shared_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_proto_common_v1_types_shared_proto_init() }
func file_pkg_proto_common_v1_types_shared_proto_init() {
	if File_pkg_proto_common_v1_types_shared_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_common_v1_types_shared_proto_rawDesc), len(file_pkg_proto_common_v1_types_shared_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_common_v1_types_shared_proto_goTypes,
		DependencyIndexes: file_pkg_proto_common_v1_types_shared_proto_depIdxs,
		MessageInfos:      file_pkg_proto_common_v1_types_shared_proto_msgTypes,
	}.Build()
	File_pkg_proto_common_v1_types_shared_proto = out.File
	file_pkg_proto_common_v1_types_shared_proto_goTypes = nil
	file_pkg_proto_common_v1_types_shared_proto_depIdxs = nil
}
