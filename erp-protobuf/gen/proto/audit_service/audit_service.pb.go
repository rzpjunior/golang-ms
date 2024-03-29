// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: proto/audit_service/audit_service.proto

package audit_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId         int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserIdGp       string                 `protobuf:"bytes,3,opt,name=user_id_gp,json=userIdGp,proto3" json:"user_id_gp,omitempty"`
	ReferenceId    string                 `protobuf:"bytes,4,opt,name=reference_id,json=referenceId,proto3" json:"reference_id,omitempty"`
	Type           string                 `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Function       string                 `protobuf:"bytes,6,opt,name=function,proto3" json:"function,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Note           string                 `protobuf:"bytes,8,opt,name=note,proto3" json:"note,omitempty"`
	SupportiveData string                 `protobuf:"bytes,9,opt,name=supportiveData,proto3" json:"supportiveData,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_audit_service_audit_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_proto_audit_service_audit_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_proto_audit_service_audit_service_proto_rawDescGZIP(), []int{0}
}

func (x *Log) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Log) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Log) GetUserIdGp() string {
	if x != nil {
		return x.UserIdGp
	}
	return ""
}

func (x *Log) GetReferenceId() string {
	if x != nil {
		return x.ReferenceId
	}
	return ""
}

func (x *Log) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Log) GetFunction() string {
	if x != nil {
		return x.Function
	}
	return ""
}

func (x *Log) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Log) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *Log) GetSupportiveData() string {
	if x != nil {
		return x.SupportiveData
	}
	return ""
}

type CreateLogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Log *Log `protobuf:"bytes,1,opt,name=log,proto3" json:"log,omitempty"`
}

func (x *CreateLogRequest) Reset() {
	*x = CreateLogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_audit_service_audit_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLogRequest) ProtoMessage() {}

func (x *CreateLogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_audit_service_audit_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLogRequest.ProtoReflect.Descriptor instead.
func (*CreateLogRequest) Descriptor() ([]byte, []int) {
	return file_proto_audit_service_audit_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateLogRequest) GetLog() *Log {
	if x != nil {
		return x.Log
	}
	return nil
}

type CreateLogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *Log   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *CreateLogResponse) Reset() {
	*x = CreateLogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_audit_service_audit_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLogResponse) ProtoMessage() {}

func (x *CreateLogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_audit_service_audit_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLogResponse.ProtoReflect.Descriptor instead.
func (*CreateLogResponse) Descriptor() ([]byte, []int) {
	return file_proto_audit_service_audit_service_proto_rawDescGZIP(), []int{2}
}

func (x *CreateLogResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CreateLogResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CreateLogResponse) GetData() *Log {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_audit_service_audit_service_proto protoreflect.FileDescriptor

var file_proto_audit_service_audit_service_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x96, 0x02, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1c, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x5f, 0x67, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x47, 0x70, 0x12, 0x21,
	0x0a, 0x0c, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x6f, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65,
	0x12, 0x26, 0x0a, 0x0e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x76, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x3e, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x03,
	0x6c, 0x6f, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4c, 0x6f, 0x67, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x22, 0x6f, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4c, 0x6f, 0x67, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x6c, 0x0a, 0x0c, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x09, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x5a, 0x5a, 0x52, 0x67, 0x69, 0x74, 0x2e, 0x65,
	0x64, 0x65, 0x6e, 0x66, 0x61, 0x72, 0x6d, 0x2e, 0x69, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x2d, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x33, 0x2f, 0x65, 0x72, 0x70, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65, 0x72, 0x70, 0x2d, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0xa2, 0x02, 0x03,
	0x45, 0x4f, 0x50, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_audit_service_audit_service_proto_rawDescOnce sync.Once
	file_proto_audit_service_audit_service_proto_rawDescData = file_proto_audit_service_audit_service_proto_rawDesc
)

func file_proto_audit_service_audit_service_proto_rawDescGZIP() []byte {
	file_proto_audit_service_audit_service_proto_rawDescOnce.Do(func() {
		file_proto_audit_service_audit_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_audit_service_audit_service_proto_rawDescData)
	})
	return file_proto_audit_service_audit_service_proto_rawDescData
}

var file_proto_audit_service_audit_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_audit_service_audit_service_proto_goTypes = []interface{}{
	(*Log)(nil),                   // 0: proto.audit_service.Log
	(*CreateLogRequest)(nil),      // 1: proto.audit_service.CreateLogRequest
	(*CreateLogResponse)(nil),     // 2: proto.audit_service.CreateLogResponse
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_proto_audit_service_audit_service_proto_depIdxs = []int32{
	3, // 0: proto.audit_service.Log.created_at:type_name -> google.protobuf.Timestamp
	0, // 1: proto.audit_service.CreateLogRequest.log:type_name -> proto.audit_service.Log
	0, // 2: proto.audit_service.CreateLogResponse.data:type_name -> proto.audit_service.Log
	1, // 3: proto.audit_service.AuditService.CreateLog:input_type -> proto.audit_service.CreateLogRequest
	2, // 4: proto.audit_service.AuditService.CreateLog:output_type -> proto.audit_service.CreateLogResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_audit_service_audit_service_proto_init() }
func file_proto_audit_service_audit_service_proto_init() {
	if File_proto_audit_service_audit_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_audit_service_audit_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
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
		file_proto_audit_service_audit_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLogRequest); i {
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
		file_proto_audit_service_audit_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLogResponse); i {
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
			RawDescriptor: file_proto_audit_service_audit_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_audit_service_audit_service_proto_goTypes,
		DependencyIndexes: file_proto_audit_service_audit_service_proto_depIdxs,
		MessageInfos:      file_proto_audit_service_audit_service_proto_msgTypes,
	}.Build()
	File_proto_audit_service_audit_service_proto = out.File
	file_proto_audit_service_audit_service_proto_rawDesc = nil
	file_proto_audit_service_audit_service_proto_goTypes = nil
	file_proto_audit_service_audit_service_proto_depIdxs = nil
}
