// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: proto/master_worker.proto

package backuptask

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

type BackupTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileSize int64  `protobuf:"varint,2,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
}

func (x *BackupTaskRequest) Reset() {
	*x = BackupTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_master_worker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackupTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackupTaskRequest) ProtoMessage() {}

func (x *BackupTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_master_worker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackupTaskRequest.ProtoReflect.Descriptor instead.
func (*BackupTaskRequest) Descriptor() ([]byte, []int) {
	return file_proto_master_worker_proto_rawDescGZIP(), []int{0}
}

func (x *BackupTaskRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *BackupTaskRequest) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

type WorkerAssignmentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkerIp     string `protobuf:"bytes,1,opt,name=worker_ip,json=workerIp,proto3" json:"worker_ip,omitempty"`
	WorkerPort   int32  `protobuf:"varint,2,opt,name=worker_port,json=workerPort,proto3" json:"worker_port,omitempty"`
	SftpUsername string `protobuf:"bytes,3,opt,name=sftp_username,json=sftpUsername,proto3" json:"sftp_username,omitempty"`
	SftpPassword string `protobuf:"bytes,4,opt,name=sftp_password,json=sftpPassword,proto3" json:"sftp_password,omitempty"`
}

func (x *WorkerAssignmentResponse) Reset() {
	*x = WorkerAssignmentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_master_worker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerAssignmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerAssignmentResponse) ProtoMessage() {}

func (x *WorkerAssignmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_master_worker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerAssignmentResponse.ProtoReflect.Descriptor instead.
func (*WorkerAssignmentResponse) Descriptor() ([]byte, []int) {
	return file_proto_master_worker_proto_rawDescGZIP(), []int{1}
}

func (x *WorkerAssignmentResponse) GetWorkerIp() string {
	if x != nil {
		return x.WorkerIp
	}
	return ""
}

func (x *WorkerAssignmentResponse) GetWorkerPort() int32 {
	if x != nil {
		return x.WorkerPort
	}
	return 0
}

func (x *WorkerAssignmentResponse) GetSftpUsername() string {
	if x != nil {
		return x.SftpUsername
	}
	return ""
}

func (x *WorkerAssignmentResponse) GetSftpPassword() string {
	if x != nil {
		return x.SftpPassword
	}
	return ""
}

type WorkerStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkerId    string `protobuf:"bytes,1,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	IsAvailable bool   `protobuf:"varint,2,opt,name=is_available,json=isAvailable,proto3" json:"is_available,omitempty"`
}

func (x *WorkerStatusRequest) Reset() {
	*x = WorkerStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_master_worker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerStatusRequest) ProtoMessage() {}

func (x *WorkerStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_master_worker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerStatusRequest.ProtoReflect.Descriptor instead.
func (*WorkerStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_master_worker_proto_rawDescGZIP(), []int{2}
}

func (x *WorkerStatusRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *WorkerStatusRequest) GetIsAvailable() bool {
	if x != nil {
		return x.IsAvailable
	}
	return false
}

type WorkerStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *WorkerStatusResponse) Reset() {
	*x = WorkerStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_master_worker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerStatusResponse) ProtoMessage() {}

func (x *WorkerStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_master_worker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerStatusResponse.ProtoReflect.Descriptor instead.
func (*WorkerStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_master_worker_proto_rawDescGZIP(), []int{3}
}

func (x *WorkerStatusResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_master_worker_proto protoreflect.FileDescriptor

var file_proto_master_worker_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x11, 0x42,
	0x61, 0x63, 0x6b, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x18, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x5f, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x49, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x66, 0x74, 0x70, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x66,
	0x74, 0x70, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x66,
	0x74, 0x70, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x73, 0x66, 0x74, 0x70, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x55, 0x0a, 0x13, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x41, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2e, 0x0a, 0x14, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x92, 0x01, 0x0a, 0x0d, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x42, 0x61, 0x63, 0x6b,
	0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14,
	0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x62,
	0x61, 0x63, 0x6b, 0x75, 0x70, 0x74, 0x61, 0x73, 0x6b, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_master_worker_proto_rawDescOnce sync.Once
	file_proto_master_worker_proto_rawDescData = file_proto_master_worker_proto_rawDesc
)

func file_proto_master_worker_proto_rawDescGZIP() []byte {
	file_proto_master_worker_proto_rawDescOnce.Do(func() {
		file_proto_master_worker_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_master_worker_proto_rawDescData)
	})
	return file_proto_master_worker_proto_rawDescData
}

var file_proto_master_worker_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_master_worker_proto_goTypes = []any{
	(*BackupTaskRequest)(nil),        // 0: BackupTaskRequest
	(*WorkerAssignmentResponse)(nil), // 1: WorkerAssignmentResponse
	(*WorkerStatusRequest)(nil),      // 2: WorkerStatusRequest
	(*WorkerStatusResponse)(nil),     // 3: WorkerStatusResponse
}
var file_proto_master_worker_proto_depIdxs = []int32{
	0, // 0: MasterService.RequestWorker:input_type -> BackupTaskRequest
	2, // 1: MasterService.ReportWorkerStatus:input_type -> WorkerStatusRequest
	1, // 2: MasterService.RequestWorker:output_type -> WorkerAssignmentResponse
	3, // 3: MasterService.ReportWorkerStatus:output_type -> WorkerStatusResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_master_worker_proto_init() }
func file_proto_master_worker_proto_init() {
	if File_proto_master_worker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_master_worker_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*BackupTaskRequest); i {
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
		file_proto_master_worker_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*WorkerAssignmentResponse); i {
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
		file_proto_master_worker_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*WorkerStatusRequest); i {
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
		file_proto_master_worker_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*WorkerStatusResponse); i {
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
			RawDescriptor: file_proto_master_worker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_master_worker_proto_goTypes,
		DependencyIndexes: file_proto_master_worker_proto_depIdxs,
		MessageInfos:      file_proto_master_worker_proto_msgTypes,
	}.Build()
	File_proto_master_worker_proto = out.File
	file_proto_master_worker_proto_rawDesc = nil
	file_proto_master_worker_proto_goTypes = nil
	file_proto_master_worker_proto_depIdxs = nil
}
