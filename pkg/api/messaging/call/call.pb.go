// Code generated by protoc-gen-go. DO NOT EDIT.
// source: call.proto

package call

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type CallPayload struct {
	DestinationPhones    []string `protobuf:"bytes,2,rep,name=destination_phones,json=destinationPhones,proto3" json:"destination_phones,omitempty"`
	Keyword              string   `protobuf:"bytes,1,opt,name=keyword,proto3" json:"keyword,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallPayload) Reset()         { *m = CallPayload{} }
func (m *CallPayload) String() string { return proto.CompactTextString(m) }
func (*CallPayload) ProtoMessage()    {}
func (*CallPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_caa5955d5eab2d2d, []int{0}
}

func (m *CallPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallPayload.Unmarshal(m, b)
}
func (m *CallPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallPayload.Marshal(b, m, deterministic)
}
func (m *CallPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallPayload.Merge(m, src)
}
func (m *CallPayload) XXX_Size() int {
	return xxx_messageInfo_CallPayload.Size(m)
}
func (m *CallPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_CallPayload.DiscardUnknown(m)
}

var xxx_messageInfo_CallPayload proto.InternalMessageInfo

func (m *CallPayload) GetDestinationPhones() []string {
	if m != nil {
		return m.DestinationPhones
	}
	return nil
}

func (m *CallPayload) GetKeyword() string {
	if m != nil {
		return m.Keyword
	}
	return ""
}

func (m *CallPayload) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*CallPayload)(nil), "umrs.messaging.call.CallPayload")
}

func init() { proto.RegisterFile("call.proto", fileDescriptor_caa5955d5eab2d2d) }

var fileDescriptor_caa5955d5eab2d2d = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xcf, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0x99, 0x1b, 0x1b, 0x8b, 0x27, 0x23, 0x48, 0xa9, 0x1e, 0x8a, 0xa7, 0x1d, 0xec, 0x2b,
	0xe8, 0x3f, 0xe0, 0x0f, 0x3c, 0x88, 0x97, 0xb2, 0xa3, 0x17, 0x79, 0x5b, 0x9f, 0xb1, 0x98, 0xe6,
	0x85, 0x26, 0x63, 0xf4, 0xbf, 0x97, 0x24, 0x1d, 0x14, 0xf4, 0x94, 0xbc, 0x7c, 0x3e, 0x21, 0xf9,
	0x7e, 0x85, 0xd8, 0xa3, 0xd6, 0x60, 0x7b, 0xf6, 0x2c, 0x2f, 0x0f, 0x5d, 0xef, 0xa0, 0x23, 0xe7,
	0x50, 0xb5, 0x46, 0x41, 0x40, 0xf9, 0xb5, 0x62, 0x56, 0x9a, 0xaa, 0xa8, 0xec, 0x0e, 0x5f, 0x15,
	0x75, 0xd6, 0x0f, 0xe9, 0x46, 0x7e, 0x33, 0x42, 0xb4, 0x6d, 0x85, 0xc6, 0xb0, 0x47, 0xdf, 0xb2,
	0x71, 0x23, 0xbd, 0x8b, 0xcb, 0xbe, 0x54, 0x64, 0x4a, 0x77, 0x44, 0xa5, 0xa8, 0xaf, 0xd8, 0x46,
	0xe3, 0xaf, 0x7d, 0x6b, 0xc5, 0xf9, 0x0b, 0x6a, 0x5d, 0xe3, 0xa0, 0x19, 0x1b, 0x59, 0x0a, 0xd9,
	0x90, 0xf3, 0xad, 0x89, 0xd2, 0xa7, 0xfd, 0x66, 0x43, 0x2e, 0x3b, 0x2b, 0xe6, 0x9b, 0xf5, 0xf6,
	0x62, 0x42, 0xea, 0x08, 0x64, 0x26, 0x56, 0x3f, 0x34, 0x1c, 0xb9, 0x6f, 0xb2, 0x59, 0x31, 0xdb,
	0xac, 0xb7, 0xa7, 0x31, 0x90, 0x14, 0x89, 0xb2, 0x79, 0x22, 0xe3, 0x78, 0xff, 0x2e, 0x56, 0xe1,
	0xc5, 0xa7, 0xfa, 0x4d, 0x3e, 0x8a, 0x45, 0xd8, 0xca, 0x02, 0xfe, 0xe9, 0x00, 0x26, 0xff, 0xca,
	0xaf, 0x20, 0x65, 0x86, 0x53, 0x21, 0xf0, 0x1a, 0x0a, 0x79, 0x5e, 0x7e, 0x2c, 0x82, 0xbb, 0x5b,
	0xc6, 0xf3, 0x87, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x19, 0xfa, 0x86, 0x4b, 0x59, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CallAPIClient is the client API for CallAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CallAPIClient interface {
	Call(ctx context.Context, in *CallPayload, opts ...grpc.CallOption) (*empty.Empty, error)
}

type callAPIClient struct {
	cc *grpc.ClientConn
}

func NewCallAPIClient(cc *grpc.ClientConn) CallAPIClient {
	return &callAPIClient{cc}
}

func (c *callAPIClient) Call(ctx context.Context, in *CallPayload, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.messaging.call.CallAPI/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CallAPIServer is the server API for CallAPI service.
type CallAPIServer interface {
	Call(context.Context, *CallPayload) (*empty.Empty, error)
}

func RegisterCallAPIServer(s *grpc.Server, srv CallAPIServer) {
	s.RegisterService(&_CallAPI_serviceDesc, srv)
}

func _CallAPI_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CallAPIServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.messaging.call.CallAPI/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CallAPIServer).Call(ctx, req.(*CallPayload))
	}
	return interceptor(ctx, in, info, handler)
}

var _CallAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "umrs.messaging.call.CallAPI",
	HandlerType: (*CallAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _CallAPI_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "call.proto",
}