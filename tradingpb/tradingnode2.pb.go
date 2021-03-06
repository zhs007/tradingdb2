// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.11.2
// source: tradingnode2.proto

package tradingpb

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

// RequestServerInfo - request server info
type RequestServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BasicRequest *BasicRequestData `protobuf:"bytes,1,opt,name=basicRequest,proto3" json:"basicRequest,omitempty"`
}

func (x *RequestServerInfo) Reset() {
	*x = RequestServerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingnode2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestServerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestServerInfo) ProtoMessage() {}

func (x *RequestServerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_tradingnode2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestServerInfo.ProtoReflect.Descriptor instead.
func (*RequestServerInfo) Descriptor() ([]byte, []int) {
	return file_tradingnode2_proto_rawDescGZIP(), []int{0}
}

func (x *RequestServerInfo) GetBasicRequest() *BasicRequestData {
	if x != nil {
		return x.BasicRequest
	}
	return nil
}

type ReplyServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeInfo *TradingNodeInfo `protobuf:"bytes,1,opt,name=nodeInfo,proto3" json:"nodeInfo,omitempty"`
}

func (x *ReplyServerInfo) Reset() {
	*x = ReplyServerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingnode2_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplyServerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyServerInfo) ProtoMessage() {}

func (x *ReplyServerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_tradingnode2_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyServerInfo.ProtoReflect.Descriptor instead.
func (*ReplyServerInfo) Descriptor() ([]byte, []int) {
	return file_tradingnode2_proto_rawDescGZIP(), []int{1}
}

func (x *ReplyServerInfo) GetNodeInfo() *TradingNodeInfo {
	if x != nil {
		return x.NodeInfo
	}
	return nil
}

// RequestCalcPNL - request calcPNL
type RequestCalcPNL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BasicRequest *BasicRequestData `protobuf:"bytes,1,opt,name=basicRequest,proto3" json:"basicRequest,omitempty"`
	Params       *SimTradingParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params,omitempty"`
}

func (x *RequestCalcPNL) Reset() {
	*x = RequestCalcPNL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingnode2_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestCalcPNL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestCalcPNL) ProtoMessage() {}

func (x *RequestCalcPNL) ProtoReflect() protoreflect.Message {
	mi := &file_tradingnode2_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestCalcPNL.ProtoReflect.Descriptor instead.
func (*RequestCalcPNL) Descriptor() ([]byte, []int) {
	return file_tradingnode2_proto_rawDescGZIP(), []int{2}
}

func (x *RequestCalcPNL) GetBasicRequest() *BasicRequestData {
	if x != nil {
		return x.BasicRequest
	}
	return nil
}

func (x *RequestCalcPNL) GetParams() *SimTradingParams {
	if x != nil {
		return x.Params
	}
	return nil
}

// ReplyCalcPNL - reply calcPNL
type ReplyCalcPNL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeInfo   *TradingNodeInfo `protobuf:"bytes,1,opt,name=nodeInfo,proto3" json:"nodeInfo,omitempty"`
	Pnl        []*PNLData       `protobuf:"bytes,2,rep,name=pnl,proto3" json:"pnl,omitempty"`
	RunSeconds int64            `protobuf:"varint,3,opt,name=runSeconds,proto3" json:"runSeconds,omitempty"`
}

func (x *ReplyCalcPNL) Reset() {
	*x = ReplyCalcPNL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingnode2_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplyCalcPNL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyCalcPNL) ProtoMessage() {}

func (x *ReplyCalcPNL) ProtoReflect() protoreflect.Message {
	mi := &file_tradingnode2_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyCalcPNL.ProtoReflect.Descriptor instead.
func (*ReplyCalcPNL) Descriptor() ([]byte, []int) {
	return file_tradingnode2_proto_rawDescGZIP(), []int{3}
}

func (x *ReplyCalcPNL) GetNodeInfo() *TradingNodeInfo {
	if x != nil {
		return x.NodeInfo
	}
	return nil
}

func (x *ReplyCalcPNL) GetPnl() []*PNLData {
	if x != nil {
		return x.Pnl
	}
	return nil
}

func (x *ReplyCalcPNL) GetRunSeconds() int64 {
	if x != nil {
		return x.RunSeconds
	}
	return 0
}

var File_tradingnode2_proto protoreflect.FileDescriptor

var file_tradingnode2_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x6e, 0x6f, 0x64, 0x65, 0x32, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x1a,
	0x0e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x54, 0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3f, 0x0a, 0x0c, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x49, 0x0a, 0x0f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x36, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x22, 0x86, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x63,
	0x50, 0x4e, 0x4c, 0x12, 0x3f, 0x0a, 0x0c, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x72, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62,
	0x2e, 0x53, 0x69, 0x6d, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x8c, 0x01, 0x0a, 0x0c, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x43, 0x61, 0x6c, 0x63, 0x50, 0x4e, 0x4c, 0x12, 0x36, 0x0a, 0x08, 0x6e, 0x6f,
	0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74,
	0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x24, 0x0a, 0x03, 0x70, 0x6e, 0x6c, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x50, 0x4e, 0x4c, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x03, 0x70, 0x6e, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x75, 0x6e, 0x53,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x75,
	0x6e, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x32, 0x9c, 0x01, 0x0a, 0x0c, 0x54, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x64, 0x65, 0x32, 0x12, 0x4b, 0x0a, 0x0d, 0x67, 0x65, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x2e, 0x74, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x1a, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x07, 0x63, 0x61, 0x6c, 0x63, 0x50, 0x4e,
	0x4c, 0x12, 0x19, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x63, 0x50, 0x4e, 0x4c, 0x1a, 0x17, 0x2e, 0x74,
	0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x43, 0x61,
	0x6c, 0x63, 0x50, 0x4e, 0x4c, 0x22, 0x00, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x68, 0x73, 0x30, 0x30, 0x37, 0x2f, 0x74, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tradingnode2_proto_rawDescOnce sync.Once
	file_tradingnode2_proto_rawDescData = file_tradingnode2_proto_rawDesc
)

func file_tradingnode2_proto_rawDescGZIP() []byte {
	file_tradingnode2_proto_rawDescOnce.Do(func() {
		file_tradingnode2_proto_rawDescData = protoimpl.X.CompressGZIP(file_tradingnode2_proto_rawDescData)
	})
	return file_tradingnode2_proto_rawDescData
}

var file_tradingnode2_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_tradingnode2_proto_goTypes = []interface{}{
	(*RequestServerInfo)(nil), // 0: tradingpb.RequestServerInfo
	(*ReplyServerInfo)(nil),   // 1: tradingpb.ReplyServerInfo
	(*RequestCalcPNL)(nil),    // 2: tradingpb.RequestCalcPNL
	(*ReplyCalcPNL)(nil),      // 3: tradingpb.ReplyCalcPNL
	(*BasicRequestData)(nil),  // 4: tradingpb.BasicRequestData
	(*TradingNodeInfo)(nil),   // 5: tradingpb.TradingNodeInfo
	(*SimTradingParams)(nil),  // 6: tradingpb.SimTradingParams
	(*PNLData)(nil),           // 7: tradingpb.PNLData
}
var file_tradingnode2_proto_depIdxs = []int32{
	4, // 0: tradingpb.RequestServerInfo.basicRequest:type_name -> tradingpb.BasicRequestData
	5, // 1: tradingpb.ReplyServerInfo.nodeInfo:type_name -> tradingpb.TradingNodeInfo
	4, // 2: tradingpb.RequestCalcPNL.basicRequest:type_name -> tradingpb.BasicRequestData
	6, // 3: tradingpb.RequestCalcPNL.params:type_name -> tradingpb.SimTradingParams
	5, // 4: tradingpb.ReplyCalcPNL.nodeInfo:type_name -> tradingpb.TradingNodeInfo
	7, // 5: tradingpb.ReplyCalcPNL.pnl:type_name -> tradingpb.PNLData
	0, // 6: tradingpb.TradingNode2.getServerInfo:input_type -> tradingpb.RequestServerInfo
	2, // 7: tradingpb.TradingNode2.calcPNL:input_type -> tradingpb.RequestCalcPNL
	1, // 8: tradingpb.TradingNode2.getServerInfo:output_type -> tradingpb.ReplyServerInfo
	3, // 9: tradingpb.TradingNode2.calcPNL:output_type -> tradingpb.ReplyCalcPNL
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_tradingnode2_proto_init() }
func file_tradingnode2_proto_init() {
	if File_tradingnode2_proto != nil {
		return
	}
	file_trading2_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tradingnode2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestServerInfo); i {
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
		file_tradingnode2_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplyServerInfo); i {
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
		file_tradingnode2_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestCalcPNL); i {
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
		file_tradingnode2_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplyCalcPNL); i {
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
			RawDescriptor: file_tradingnode2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tradingnode2_proto_goTypes,
		DependencyIndexes: file_tradingnode2_proto_depIdxs,
		MessageInfos:      file_tradingnode2_proto_msgTypes,
	}.Build()
	File_tradingnode2_proto = out.File
	file_tradingnode2_proto_rawDesc = nil
	file_tradingnode2_proto_goTypes = nil
	file_tradingnode2_proto_depIdxs = nil
}
