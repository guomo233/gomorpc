// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: gomorpcpb/mheader.proto

package __

import (
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

type MHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Method  string `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Id      uint32 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Error   uint32 `protobuf:"varint,4,opt,name=error,proto3" json:"error,omitempty"`
	Timeout uint64 `protobuf:"varint,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *MHeader) Reset() {
	*x = MHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gomorpcpb_mheader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MHeader) ProtoMessage() {}

func (x *MHeader) ProtoReflect() protoreflect.Message {
	mi := &file_gomorpcpb_mheader_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MHeader.ProtoReflect.Descriptor instead.
func (*MHeader) Descriptor() ([]byte, []int) {
	return file_gomorpcpb_mheader_proto_rawDescGZIP(), []int{0}
}

func (x *MHeader) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *MHeader) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *MHeader) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MHeader) GetError() uint32 {
	if x != nil {
		return x.Error
	}
	return 0
}

func (x *MHeader) GetTimeout() uint64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

var File_gomorpcpb_mheader_proto protoreflect.FileDescriptor

var file_gomorpcpb_mheader_proto_rawDesc = []byte{
	0x0a, 0x17, 0x67, 0x6f, 0x6d, 0x6f, 0x72, 0x70, 0x63, 0x70, 0x62, 0x2f, 0x6d, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x6f, 0x6d, 0x6f, 0x72,
	0x70, 0x63, 0x70, 0x62, 0x22, 0x7b, 0x0a, 0x07, 0x4d, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gomorpcpb_mheader_proto_rawDescOnce sync.Once
	file_gomorpcpb_mheader_proto_rawDescData = file_gomorpcpb_mheader_proto_rawDesc
)

func file_gomorpcpb_mheader_proto_rawDescGZIP() []byte {
	file_gomorpcpb_mheader_proto_rawDescOnce.Do(func() {
		file_gomorpcpb_mheader_proto_rawDescData = protoimpl.X.CompressGZIP(file_gomorpcpb_mheader_proto_rawDescData)
	})
	return file_gomorpcpb_mheader_proto_rawDescData
}

var file_gomorpcpb_mheader_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_gomorpcpb_mheader_proto_goTypes = []interface{}{
	(*MHeader)(nil), // 0: gomorpcpb.MHeader
}
var file_gomorpcpb_mheader_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gomorpcpb_mheader_proto_init() }
func file_gomorpcpb_mheader_proto_init() {
	if File_gomorpcpb_mheader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gomorpcpb_mheader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MHeader); i {
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
			RawDescriptor: file_gomorpcpb_mheader_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gomorpcpb_mheader_proto_goTypes,
		DependencyIndexes: file_gomorpcpb_mheader_proto_depIdxs,
		MessageInfos:      file_gomorpcpb_mheader_proto_msgTypes,
	}.Build()
	File_gomorpcpb_mheader_proto = out.File
	file_gomorpcpb_mheader_proto_rawDesc = nil
	file_gomorpcpb_mheader_proto_goTypes = nil
	file_gomorpcpb_mheader_proto_depIdxs = nil
}
