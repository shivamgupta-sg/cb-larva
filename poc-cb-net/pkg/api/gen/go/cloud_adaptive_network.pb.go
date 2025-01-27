// cloud_adaptive_network.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: cloud_adaptive_network.proto

package cb_larva

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// A specification of Cloud Adpative Network
type CLADNetSpecification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// AWS - Create Virtual Private Cloud (VPC): IPv4 CIDR block
	// MS Azure - Create Virtual Network (vNet): IPv4 address space
	// GCP - Create VPC Network: IPv4 address range
	// Alibaba Cloud - Create VPC: IPv4 CIDR Block
	Ipv4AddressSpace string `protobuf:"bytes,3,opt,name=ipv4_address_space,json=ipv4AddressSpace,proto3" json:"ipv4_address_space,omitempty"`
	Description      string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *CLADNetSpecification) Reset() {
	*x = CLADNetSpecification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloud_adaptive_network_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CLADNetSpecification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CLADNetSpecification) ProtoMessage() {}

func (x *CLADNetSpecification) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_adaptive_network_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CLADNetSpecification.ProtoReflect.Descriptor instead.
func (*CLADNetSpecification) Descriptor() ([]byte, []int) {
	return file_cloud_adaptive_network_proto_rawDescGZIP(), []int{0}
}

func (x *CLADNetSpecification) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CLADNetSpecification) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CLADNetSpecification) GetIpv4AddressSpace() string {
	if x != nil {
		return x.Ipv4AddressSpace
	}
	return ""
}

func (x *CLADNetSpecification) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// An ID of Cloud Adpative Network
type CLADNetID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *CLADNetID) Reset() {
	*x = CLADNetID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloud_adaptive_network_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CLADNetID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CLADNetID) ProtoMessage() {}

func (x *CLADNetID) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_adaptive_network_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CLADNetID.ProtoReflect.Descriptor instead.
func (*CLADNetID) Descriptor() ([]byte, []int) {
	return file_cloud_adaptive_network_proto_rawDescGZIP(), []int{1}
}

func (x *CLADNetID) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Request message for CreateCLADNet method
type CreateCLADNetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CladnetSpecification *CLADNetSpecification `protobuf:"bytes,1,opt,name=cladnet_specification,json=cladnetSpecification,proto3" json:"cladnet_specification,omitempty"`
}

func (x *CreateCLADNetRequest) Reset() {
	*x = CreateCLADNetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloud_adaptive_network_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCLADNetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCLADNetRequest) ProtoMessage() {}

func (x *CreateCLADNetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_adaptive_network_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCLADNetRequest.ProtoReflect.Descriptor instead.
func (*CreateCLADNetRequest) Descriptor() ([]byte, []int) {
	return file_cloud_adaptive_network_proto_rawDescGZIP(), []int{2}
}

func (x *CreateCLADNetRequest) GetCladnetSpecification() *CLADNetSpecification {
	if x != nil {
		return x.CladnetSpecification
	}
	return nil
}

// Response message for common
type CLADNetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsSucceeded          bool                  `protobuf:"varint,1,opt,name=is_succeeded,json=isSucceeded,proto3" json:"is_succeeded,omitempty"`
	Message              string                `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	CladnetSpecification *CLADNetSpecification `protobuf:"bytes,3,opt,name=cladnet_specification,json=cladnetSpecification,proto3" json:"cladnet_specification,omitempty"`
}

func (x *CLADNetResponse) Reset() {
	*x = CLADNetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloud_adaptive_network_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CLADNetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CLADNetResponse) ProtoMessage() {}

func (x *CLADNetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_adaptive_network_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CLADNetResponse.ProtoReflect.Descriptor instead.
func (*CLADNetResponse) Descriptor() ([]byte, []int) {
	return file_cloud_adaptive_network_proto_rawDescGZIP(), []int{3}
}

func (x *CLADNetResponse) GetIsSucceeded() bool {
	if x != nil {
		return x.IsSucceeded
	}
	return false
}

func (x *CLADNetResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CLADNetResponse) GetCladnetSpecification() *CLADNetSpecification {
	if x != nil {
		return x.CladnetSpecification
	}
	return nil
}

var File_cloud_adaptive_network_proto protoreflect.FileDescriptor

var file_cloud_adaptive_network_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x61, 0x64, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65,
	0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x63, 0x62, 0x6e, 0x65, 0x74, 0x22, 0x8a, 0x01, 0x0a, 0x14, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65,
	0x74, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x69, 0x70, 0x76, 0x34, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x5f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x69, 0x70, 0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x21, 0x0a, 0x09, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x49, 0x44, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x68, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x50, 0x0a,
	0x15, 0x63, 0x6c, 0x61, 0x64, 0x6e, 0x65, 0x74, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63,
	0x62, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x53, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x14, 0x63, 0x6c, 0x61, 0x64, 0x6e,
	0x65, 0x74, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0xa0, 0x01, 0x0a, 0x0f, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65,
	0x64, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x53, 0x75, 0x63,
	0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x50, 0x0a, 0x15, 0x63, 0x6c, 0x61, 0x64, 0x6e, 0x65, 0x74, 0x5f, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x63, 0x62, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x53,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x14, 0x63, 0x6c,
	0x61, 0x64, 0x6e, 0x65, 0x74, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x32, 0x98, 0x01, 0x0a, 0x14, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x41, 0x64, 0x61, 0x70,
	0x74, 0x69, 0x76, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x46, 0x0a, 0x0d, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x63,
	0x62, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x4c, 0x41, 0x44, 0x4e,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x63, 0x62, 0x6e, 0x65,
	0x74, 0x2e, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0a, 0x67, 0x65, 0x74, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65,
	0x74, 0x12, 0x10, 0x2e, 0x63, 0x62, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x4c, 0x41, 0x44, 0x4e, 0x65,
	0x74, 0x49, 0x44, 0x1a, 0x16, 0x2e, 0x63, 0x62, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x4c, 0x41, 0x44,
	0x4e, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x23, 0x5a,
	0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x2d, 0x62, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x63, 0x62, 0x2d, 0x6c, 0x61, 0x72,
	0x76, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cloud_adaptive_network_proto_rawDescOnce sync.Once
	file_cloud_adaptive_network_proto_rawDescData = file_cloud_adaptive_network_proto_rawDesc
)

func file_cloud_adaptive_network_proto_rawDescGZIP() []byte {
	file_cloud_adaptive_network_proto_rawDescOnce.Do(func() {
		file_cloud_adaptive_network_proto_rawDescData = protoimpl.X.CompressGZIP(file_cloud_adaptive_network_proto_rawDescData)
	})
	return file_cloud_adaptive_network_proto_rawDescData
}

var file_cloud_adaptive_network_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_cloud_adaptive_network_proto_goTypes = []interface{}{
	(*CLADNetSpecification)(nil), // 0: cbnet.CLADNetSpecification
	(*CLADNetID)(nil),            // 1: cbnet.CLADNetID
	(*CreateCLADNetRequest)(nil), // 2: cbnet.CreateCLADNetRequest
	(*CLADNetResponse)(nil),      // 3: cbnet.CLADNetResponse
}
var file_cloud_adaptive_network_proto_depIdxs = []int32{
	0, // 0: cbnet.CreateCLADNetRequest.cladnet_specification:type_name -> cbnet.CLADNetSpecification
	0, // 1: cbnet.CLADNetResponse.cladnet_specification:type_name -> cbnet.CLADNetSpecification
	2, // 2: cbnet.CloudAdaptiveNetwork.createCLADNet:input_type -> cbnet.CreateCLADNetRequest
	1, // 3: cbnet.CloudAdaptiveNetwork.getCLADNet:input_type -> cbnet.CLADNetID
	3, // 4: cbnet.CloudAdaptiveNetwork.createCLADNet:output_type -> cbnet.CLADNetResponse
	3, // 5: cbnet.CloudAdaptiveNetwork.getCLADNet:output_type -> cbnet.CLADNetResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cloud_adaptive_network_proto_init() }
func file_cloud_adaptive_network_proto_init() {
	if File_cloud_adaptive_network_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cloud_adaptive_network_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CLADNetSpecification); i {
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
		file_cloud_adaptive_network_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CLADNetID); i {
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
		file_cloud_adaptive_network_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCLADNetRequest); i {
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
		file_cloud_adaptive_network_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CLADNetResponse); i {
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
			RawDescriptor: file_cloud_adaptive_network_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cloud_adaptive_network_proto_goTypes,
		DependencyIndexes: file_cloud_adaptive_network_proto_depIdxs,
		MessageInfos:      file_cloud_adaptive_network_proto_msgTypes,
	}.Build()
	File_cloud_adaptive_network_proto = out.File
	file_cloud_adaptive_network_proto_rawDesc = nil
	file_cloud_adaptive_network_proto_goTypes = nil
	file_cloud_adaptive_network_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CloudAdaptiveNetworkClient is the client API for CloudAdaptiveNetwork service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CloudAdaptiveNetworkClient interface {
	// Creates a new CLADNet
	CreateCLADNet(ctx context.Context, in *CreateCLADNetRequest, opts ...grpc.CallOption) (*CLADNetResponse, error)
	// Returns a specific CLADNet
	GetCLADNet(ctx context.Context, in *CLADNetID, opts ...grpc.CallOption) (*CLADNetResponse, error)
}

type cloudAdaptiveNetworkClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudAdaptiveNetworkClient(cc grpc.ClientConnInterface) CloudAdaptiveNetworkClient {
	return &cloudAdaptiveNetworkClient{cc}
}

func (c *cloudAdaptiveNetworkClient) CreateCLADNet(ctx context.Context, in *CreateCLADNetRequest, opts ...grpc.CallOption) (*CLADNetResponse, error) {
	out := new(CLADNetResponse)
	err := c.cc.Invoke(ctx, "/cbnet.CloudAdaptiveNetwork/createCLADNet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudAdaptiveNetworkClient) GetCLADNet(ctx context.Context, in *CLADNetID, opts ...grpc.CallOption) (*CLADNetResponse, error) {
	out := new(CLADNetResponse)
	err := c.cc.Invoke(ctx, "/cbnet.CloudAdaptiveNetwork/getCLADNet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudAdaptiveNetworkServer is the server API for CloudAdaptiveNetwork service.
type CloudAdaptiveNetworkServer interface {
	// Creates a new CLADNet
	CreateCLADNet(context.Context, *CreateCLADNetRequest) (*CLADNetResponse, error)
	// Returns a specific CLADNet
	GetCLADNet(context.Context, *CLADNetID) (*CLADNetResponse, error)
}

// UnimplementedCloudAdaptiveNetworkServer can be embedded to have forward compatible implementations.
type UnimplementedCloudAdaptiveNetworkServer struct {
}

func (*UnimplementedCloudAdaptiveNetworkServer) CreateCLADNet(context.Context, *CreateCLADNetRequest) (*CLADNetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCLADNet not implemented")
}
func (*UnimplementedCloudAdaptiveNetworkServer) GetCLADNet(context.Context, *CLADNetID) (*CLADNetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCLADNet not implemented")
}

func RegisterCloudAdaptiveNetworkServer(s *grpc.Server, srv CloudAdaptiveNetworkServer) {
	s.RegisterService(&_CloudAdaptiveNetwork_serviceDesc, srv)
}

func _CloudAdaptiveNetwork_CreateCLADNet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCLADNetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudAdaptiveNetworkServer).CreateCLADNet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cbnet.CloudAdaptiveNetwork/CreateCLADNet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudAdaptiveNetworkServer).CreateCLADNet(ctx, req.(*CreateCLADNetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudAdaptiveNetwork_GetCLADNet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CLADNetID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudAdaptiveNetworkServer).GetCLADNet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cbnet.CloudAdaptiveNetwork/GetCLADNet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudAdaptiveNetworkServer).GetCLADNet(ctx, req.(*CLADNetID))
	}
	return interceptor(ctx, in, info, handler)
}

var _CloudAdaptiveNetwork_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cbnet.CloudAdaptiveNetwork",
	HandlerType: (*CloudAdaptiveNetworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createCLADNet",
			Handler:    _CloudAdaptiveNetwork_CreateCLADNet_Handler,
		},
		{
			MethodName: "getCLADNet",
			Handler:    _CloudAdaptiveNetwork_GetCLADNet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cloud_adaptive_network.proto",
}
