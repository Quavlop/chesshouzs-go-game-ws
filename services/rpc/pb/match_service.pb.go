// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: match_service.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_match_service_proto protoreflect.FileDescriptor

var file_match_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x19, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0x49, 0x0a, 0x0c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x0c, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x6f, 0x76, 0x65, 0x12, 0x13, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x42,
	0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_match_service_proto_goTypes = []any{
	(*ValidateMoveReq)(nil),  // 0: pb.ValidateMoveReq
	(*ValidateMoveResp)(nil), // 1: pb.ValidateMoveResp
}
var file_match_service_proto_depIdxs = []int32{
	0, // 0: pb.MatchService.ValidateMove:input_type -> pb.ValidateMoveReq
	1, // 1: pb.MatchService.ValidateMove:output_type -> pb.ValidateMoveResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_match_service_proto_init() }
func file_match_service_proto_init() {
	if File_match_service_proto != nil {
		return
	}
	file_match_service_model_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_match_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_match_service_proto_goTypes,
		DependencyIndexes: file_match_service_proto_depIdxs,
	}.Build()
	File_match_service_proto = out.File
	file_match_service_proto_rawDesc = nil
	file_match_service_proto_goTypes = nil
	file_match_service_proto_depIdxs = nil
}
