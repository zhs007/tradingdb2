// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.2
// source: tradingdb2.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Candle - candle
type Candle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts           int64 `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
	Open         int64 `protobuf:"varint,2,opt,name=open,proto3" json:"open,omitempty"`
	Close        int64 `protobuf:"varint,3,opt,name=close,proto3" json:"close,omitempty"`
	High         int64 `protobuf:"varint,4,opt,name=high,proto3" json:"high,omitempty"`
	Low          int64 `protobuf:"varint,5,opt,name=low,proto3" json:"low,omitempty"`
	Volume       int64 `protobuf:"varint,6,opt,name=volume,proto3" json:"volume,omitempty"`
	OpenInterest int64 `protobuf:"varint,7,opt,name=openInterest,proto3" json:"openInterest,omitempty"`
}

func (x *Candle) Reset() {
	*x = Candle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingdb2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Candle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Candle) ProtoMessage() {}

func (x *Candle) ProtoReflect() protoreflect.Message {
	mi := &file_tradingdb2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Candle.ProtoReflect.Descriptor instead.
func (*Candle) Descriptor() ([]byte, []int) {
	return file_tradingdb2_proto_rawDescGZIP(), []int{0}
}

func (x *Candle) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *Candle) GetOpen() int64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *Candle) GetClose() int64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *Candle) GetHigh() int64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *Candle) GetLow() int64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *Candle) GetVolume() int64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *Candle) GetOpenInterest() int64 {
	if x != nil {
		return x.OpenInterest
	}
	return 0
}

// Candles - candles
type Candles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Market  string    `protobuf:"bytes,1,opt,name=market,proto3" json:"market,omitempty"`
	Symbol  string    `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Tag     string    `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	Candles []*Candle `protobuf:"bytes,4,rep,name=candles,proto3" json:"candles,omitempty"`
}

func (x *Candles) Reset() {
	*x = Candles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingdb2_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Candles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Candles) ProtoMessage() {}

func (x *Candles) ProtoReflect() protoreflect.Message {
	mi := &file_tradingdb2_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Candles.ProtoReflect.Descriptor instead.
func (*Candles) Descriptor() ([]byte, []int) {
	return file_tradingdb2_proto_rawDescGZIP(), []int{1}
}

func (x *Candles) GetMarket() string {
	if x != nil {
		return x.Market
	}
	return ""
}

func (x *Candles) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Candles) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Candles) GetCandles() []*Candle {
	if x != nil {
		return x.Candles
	}
	return nil
}

// ReplyUpdCandle - reply UpdCandle
type ReplyUpdCandle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error    string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	LengthOK int32  `protobuf:"varint,2,opt,name=lengthOK,proto3" json:"lengthOK,omitempty"`
}

func (x *ReplyUpdCandle) Reset() {
	*x = ReplyUpdCandle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradingdb2_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplyUpdCandle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyUpdCandle) ProtoMessage() {}

func (x *ReplyUpdCandle) ProtoReflect() protoreflect.Message {
	mi := &file_tradingdb2_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyUpdCandle.ProtoReflect.Descriptor instead.
func (*ReplyUpdCandle) Descriptor() ([]byte, []int) {
	return file_tradingdb2_proto_rawDescGZIP(), []int{2}
}

func (x *ReplyUpdCandle) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ReplyUpdCandle) GetLengthOK() int32 {
	if x != nil {
		return x.LengthOK
	}
	return 0
}

var File_tradingdb2_proto protoreflect.FileDescriptor

var file_tradingdb2_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32, 0x70, 0x62,
	0x22, 0xa4, 0x01, 0x0a, 0x06, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6f,
	0x70, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x76,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x76, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x65, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x70, 0x65, 0x6e, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x22, 0x7b, 0x0a, 0x07, 0x43, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x74, 0x61, 0x67, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x64,
	0x62, 0x32, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x07, 0x63, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x73, 0x22, 0x42, 0x0a, 0x0e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x55, 0x70, 0x64,
	0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1a, 0x0a, 0x08,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x4f, 0x4b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x4f, 0x4b, 0x32, 0x56, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x44, 0x42, 0x32, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a,
	0x09, 0x75, 0x70, 0x64, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x14, 0x2e, 0x74, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x1a, 0x1c, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x55, 0x70, 0x64, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x22, 0x00,
	0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a,
	0x68, 0x73, 0x30, 0x30, 0x37, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x64, 0x62, 0x32,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tradingdb2_proto_rawDescOnce sync.Once
	file_tradingdb2_proto_rawDescData = file_tradingdb2_proto_rawDesc
)

func file_tradingdb2_proto_rawDescGZIP() []byte {
	file_tradingdb2_proto_rawDescOnce.Do(func() {
		file_tradingdb2_proto_rawDescData = protoimpl.X.CompressGZIP(file_tradingdb2_proto_rawDescData)
	})
	return file_tradingdb2_proto_rawDescData
}

var file_tradingdb2_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tradingdb2_proto_goTypes = []interface{}{
	(*Candle)(nil),         // 0: tradingdb2pb.Candle
	(*Candles)(nil),        // 1: tradingdb2pb.Candles
	(*ReplyUpdCandle)(nil), // 2: tradingdb2pb.ReplyUpdCandle
}
var file_tradingdb2_proto_depIdxs = []int32{
	0, // 0: tradingdb2pb.Candles.candles:type_name -> tradingdb2pb.Candle
	0, // 1: tradingdb2pb.TradingDB2Service.updCandle:input_type -> tradingdb2pb.Candle
	2, // 2: tradingdb2pb.TradingDB2Service.updCandle:output_type -> tradingdb2pb.ReplyUpdCandle
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tradingdb2_proto_init() }
func file_tradingdb2_proto_init() {
	if File_tradingdb2_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tradingdb2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Candle); i {
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
		file_tradingdb2_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Candles); i {
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
		file_tradingdb2_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplyUpdCandle); i {
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
			RawDescriptor: file_tradingdb2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tradingdb2_proto_goTypes,
		DependencyIndexes: file_tradingdb2_proto_depIdxs,
		MessageInfos:      file_tradingdb2_proto_msgTypes,
	}.Build()
	File_tradingdb2_proto = out.File
	file_tradingdb2_proto_rawDesc = nil
	file_tradingdb2_proto_goTypes = nil
	file_tradingdb2_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TradingDB2ServiceClient is the client API for TradingDB2Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TradingDB2ServiceClient interface {
	// updCandle - update candle
	UpdCandle(ctx context.Context, in *Candle, opts ...grpc.CallOption) (*ReplyUpdCandle, error)
}

type tradingDB2ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTradingDB2ServiceClient(cc grpc.ClientConnInterface) TradingDB2ServiceClient {
	return &tradingDB2ServiceClient{cc}
}

func (c *tradingDB2ServiceClient) UpdCandle(ctx context.Context, in *Candle, opts ...grpc.CallOption) (*ReplyUpdCandle, error) {
	out := new(ReplyUpdCandle)
	err := c.cc.Invoke(ctx, "/tradingdb2pb.TradingDB2Service/updCandle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TradingDB2ServiceServer is the server API for TradingDB2Service service.
type TradingDB2ServiceServer interface {
	// updCandle - update candle
	UpdCandle(context.Context, *Candle) (*ReplyUpdCandle, error)
}

// UnimplementedTradingDB2ServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTradingDB2ServiceServer struct {
}

func (*UnimplementedTradingDB2ServiceServer) UpdCandle(context.Context, *Candle) (*ReplyUpdCandle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdCandle not implemented")
}

func RegisterTradingDB2ServiceServer(s *grpc.Server, srv TradingDB2ServiceServer) {
	s.RegisterService(&_TradingDB2Service_serviceDesc, srv)
}

func _TradingDB2Service_UpdCandle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Candle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradingDB2ServiceServer).UpdCandle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tradingdb2pb.TradingDB2Service/UpdCandle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradingDB2ServiceServer).UpdCandle(ctx, req.(*Candle))
	}
	return interceptor(ctx, in, info, handler)
}

var _TradingDB2Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tradingdb2pb.TradingDB2Service",
	HandlerType: (*TradingDB2ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "updCandle",
			Handler:    _TradingDB2Service_UpdCandle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tradingdb2.proto",
}
