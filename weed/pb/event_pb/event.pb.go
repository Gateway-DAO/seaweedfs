// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: event.proto

package event_pb

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

type MerkleTree struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digest string            `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
	Tree   map[string]string `protobuf:"bytes,2,rep,name=tree,proto3" json:"tree,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MerkleTree) Reset() {
	*x = MerkleTree{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerkleTree) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerkleTree) ProtoMessage() {}

func (x *MerkleTree) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerkleTree.ProtoReflect.Descriptor instead.
func (*MerkleTree) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

func (x *MerkleTree) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *MerkleTree) GetTree() map[string]string {
	if x != nil {
		return x.Tree
	}
	return nil
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tree       *MerkleTree `protobuf:"bytes,1,opt,name=tree,proto3,oneof" json:"tree,omitempty"`
	PublicUrl  string      `protobuf:"bytes,2,opt,name=publicUrl,proto3" json:"publicUrl,omitempty"`
	Rack       string      `protobuf:"bytes,3,opt,name=rack,proto3" json:"rack,omitempty"`
	DataCenter string      `protobuf:"bytes,4,opt,name=dataCenter,proto3" json:"dataCenter,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetTree() *MerkleTree {
	if x != nil {
		return x.Tree
	}
	return nil
}

func (x *Server) GetPublicUrl() string {
	if x != nil {
		return x.PublicUrl
	}
	return ""
}

func (x *Server) GetRack() string {
	if x != nil {
		return x.Rack
	}
	return ""
}

func (x *Server) GetDataCenter() string {
	if x != nil {
		return x.DataCenter
	}
	return ""
}

type ProofOfHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreviousHash *string `protobuf:"bytes,1,opt,name=previous_hash,json=previousHash,proto3,oneof" json:"previous_hash,omitempty"`
	Hash         string  `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Signature    string  `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ProofOfHistory) Reset() {
	*x = ProofOfHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProofOfHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProofOfHistory) ProtoMessage() {}

func (x *ProofOfHistory) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProofOfHistory.ProtoReflect.Descriptor instead.
func (*ProofOfHistory) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{2}
}

func (x *ProofOfHistory) GetPreviousHash() string {
	if x != nil && x.PreviousHash != nil {
		return *x.PreviousHash
	}
	return ""
}

func (x *ProofOfHistory) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *ProofOfHistory) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

var File_event_proto protoreflect.FileDescriptor

var file_event_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x62, 0x22, 0x91, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x72, 0x6b,
	0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x32,
	0x0a, 0x04, 0x74, 0x72, 0x65, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72,
	0x65, 0x65, 0x2e, 0x54, 0x72, 0x65, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x74, 0x72,
	0x65, 0x65, 0x1a, 0x37, 0x0a, 0x09, 0x54, 0x72, 0x65, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x92, 0x01, 0x0a, 0x06,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x72, 0x65, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x62, 0x2e,
	0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x48, 0x00, 0x52, 0x04, 0x74, 0x72,
	0x65, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x55,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x72, 0x61, 0x63, 0x6b, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x43,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74,
	0x61, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x72, 0x65, 0x65,
	0x22, 0x7e, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x4f, 0x66, 0x48, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x28, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x48, 0x61, 0x73, 0x68, 0x88, 0x01, 0x01, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x10,
	0x0a, 0x0e, 0x5f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2d, 0x64, 0x61, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x77, 0x65,
	0x65, 0x64, 0x66, 0x73, 0x2f, 0x77, 0x65, 0x65, 0x64, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_proto_rawDescOnce sync.Once
	file_event_proto_rawDescData = file_event_proto_rawDesc
)

func file_event_proto_rawDescGZIP() []byte {
	file_event_proto_rawDescOnce.Do(func() {
		file_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_proto_rawDescData)
	})
	return file_event_proto_rawDescData
}

var file_event_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_event_proto_goTypes = []interface{}{
	(*MerkleTree)(nil),     // 0: event_pb.MerkleTree
	(*Server)(nil),         // 1: event_pb.Server
	(*ProofOfHistory)(nil), // 2: event_pb.ProofOfHistory
	nil,                    // 3: event_pb.MerkleTree.TreeEntry
}
var file_event_proto_depIdxs = []int32{
	3, // 0: event_pb.MerkleTree.tree:type_name -> event_pb.MerkleTree.TreeEntry
	0, // 1: event_pb.Server.tree:type_name -> event_pb.MerkleTree
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_event_proto_init() }
func file_event_proto_init() {
	if File_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerkleTree); i {
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
		file_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
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
		file_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProofOfHistory); i {
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
	file_event_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_event_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_event_proto_goTypes,
		DependencyIndexes: file_event_proto_depIdxs,
		MessageInfos:      file_event_proto_msgTypes,
	}.Build()
	File_event_proto = out.File
	file_event_proto_rawDesc = nil
	file_event_proto_goTypes = nil
	file_event_proto_depIdxs = nil
}
