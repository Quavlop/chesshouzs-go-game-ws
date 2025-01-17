// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: match_service_model.proto

package pb

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

type ValidateMoveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OldState string `protobuf:"bytes,1,opt,name=oldState,proto3" json:"oldState,omitempty"`
	NewState string `protobuf:"bytes,2,opt,name=newState,proto3" json:"newState,omitempty"`
}

func (x *ValidateMoveReq) Reset() {
	*x = ValidateMoveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_match_service_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateMoveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateMoveReq) ProtoMessage() {}

func (x *ValidateMoveReq) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateMoveReq.ProtoReflect.Descriptor instead.
func (*ValidateMoveReq) Descriptor() ([]byte, []int) {
	return file_match_service_model_proto_rawDescGZIP(), []int{0}
}

func (x *ValidateMoveReq) GetOldState() string {
	if x != nil {
		return x.OldState
	}
	return ""
}

func (x *ValidateMoveReq) GetNewState() string {
	if x != nil {
		return x.NewState
	}
	return ""
}

type ValidateMoveResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid bool `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
}

func (x *ValidateMoveResp) Reset() {
	*x = ValidateMoveResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_match_service_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateMoveResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateMoveResp) ProtoMessage() {}

func (x *ValidateMoveResp) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateMoveResp.ProtoReflect.Descriptor instead.
func (*ValidateMoveResp) Descriptor() ([]byte, []int) {
	return file_match_service_model_proto_rawDescGZIP(), []int{1}
}

func (x *ValidateMoveResp) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

var File_match_service_model_proto protoreflect.FileDescriptor

var file_match_service_model_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22,
	0x49, 0x0a, 0x0f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x6e, 0x65, 0x77, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6e, 0x65, 0x77, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x28, 0x0a, 0x10, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_match_service_model_proto_rawDescOnce sync.Once
	file_match_service_model_proto_rawDescData = file_match_service_model_proto_rawDesc
)

func file_match_service_model_proto_rawDescGZIP() []byte {
	file_match_service_model_proto_rawDescOnce.Do(func() {
		file_match_service_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_match_service_model_proto_rawDescData)
	})
	return file_match_service_model_proto_rawDescData
}

var file_match_service_model_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_match_service_model_proto_goTypes = []any{
	(*ValidateMoveReq)(nil),  // 0: pb.ValidateMoveReq
	(*ValidateMoveResp)(nil), // 1: pb.ValidateMoveResp
}
var file_match_service_model_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_match_service_model_proto_init() }
func file_match_service_model_proto_init() {
	if File_match_service_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_match_service_model_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ValidateMoveReq); i {
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
		file_match_service_model_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ValidateMoveResp); i {
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
			RawDescriptor: file_match_service_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_match_service_model_proto_goTypes,
		DependencyIndexes: file_match_service_model_proto_depIdxs,
		MessageInfos:      file_match_service_model_proto_msgTypes,
	}.Build()
	File_match_service_model_proto = out.File
	file_match_service_model_proto_rawDesc = nil
	file_match_service_model_proto_goTypes = nil
	file_match_service_model_proto_depIdxs = nil
}
