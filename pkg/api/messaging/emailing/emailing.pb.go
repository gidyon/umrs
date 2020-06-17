// Code generated by protoc-gen-go. DO NOT EDIT.
// source: emailing.proto

package emailing

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

type Email struct {
	Destinations         []string `protobuf:"bytes,2,rep,name=destinations,proto3" json:"destinations,omitempty"`
	From                 string   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	Subject              string   `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	BodyContentType      string   `protobuf:"bytes,4,opt,name=body_content_type,json=bodyContentType,proto3" json:"body_content_type,omitempty"`
	Body                 string   `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Email) Reset()         { *m = Email{} }
func (m *Email) String() string { return proto.CompactTextString(m) }
func (*Email) ProtoMessage()    {}
func (*Email) Descriptor() ([]byte, []int) {
	return fileDescriptor_260a6c41a6a28e68, []int{0}
}

func (m *Email) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Email.Unmarshal(m, b)
}
func (m *Email) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Email.Marshal(b, m, deterministic)
}
func (m *Email) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Email.Merge(m, src)
}
func (m *Email) XXX_Size() int {
	return xxx_messageInfo_Email.Size(m)
}
func (m *Email) XXX_DiscardUnknown() {
	xxx_messageInfo_Email.DiscardUnknown(m)
}

var xxx_messageInfo_Email proto.InternalMessageInfo

func (m *Email) GetDestinations() []string {
	if m != nil {
		return m.Destinations
	}
	return nil
}

func (m *Email) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Email) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Email) GetBodyContentType() string {
	if m != nil {
		return m.BodyContentType
	}
	return ""
}

func (m *Email) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*Email)(nil), "umrs.messaging.emailing.Email")
}

func init() { proto.RegisterFile("emailing.proto", fileDescriptor_260a6c41a6a28e68) }

var fileDescriptor_260a6c41a6a28e68 = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0x99, 0xdb, 0x74, 0x0d, 0xa2, 0x98, 0x83, 0x86, 0x2a, 0x32, 0x76, 0x1a, 0xe2, 0x52,
	0xd0, 0x6f, 0xe0, 0xe8, 0x59, 0x98, 0x9e, 0xbc, 0x8c, 0xb4, 0x7d, 0x0b, 0x91, 0x25, 0x2f, 0x34,
	0x29, 0xd2, 0x8f, 0xe2, 0xb7, 0x95, 0x24, 0xed, 0x41, 0xc4, 0x53, 0x5f, 0xff, 0xbf, 0x7f, 0x5e,
	0xc2, 0x8f, 0x5c, 0x80, 0x16, 0xea, 0xa8, 0x8c, 0xe4, 0xb6, 0x45, 0x8f, 0xf4, 0xa6, 0xd3, 0xad,
	0xe3, 0x1a, 0x9c, 0x13, 0x32, 0xa4, 0x23, 0xce, 0x6f, 0x25, 0xa2, 0x3c, 0x42, 0x11, 0x6b, 0x55,
	0x77, 0x28, 0x40, 0x5b, 0xdf, 0xa7, 0x53, 0xf9, 0xdd, 0x00, 0x85, 0x55, 0x85, 0x30, 0x06, 0xbd,
	0xf0, 0x0a, 0x8d, 0x1b, 0xe8, 0x63, 0xfc, 0xd4, 0x1b, 0x09, 0x66, 0xe3, 0xbe, 0x84, 0x94, 0xd0,
	0x16, 0x68, 0x63, 0xe3, 0x6f, 0x7b, 0xf5, 0x3d, 0x21, 0xf3, 0x32, 0xdc, 0x4a, 0x57, 0xe4, 0xbc,
	0x01, 0xe7, 0x95, 0x49, 0x9c, 0x9d, 0x2c, 0xa7, 0xeb, 0x6c, 0xf7, 0x2b, 0xa3, 0x94, 0xcc, 0x0e,
	0x2d, 0x6a, 0x36, 0x59, 0x4e, 0xd6, 0xd9, 0x2e, 0xce, 0x94, 0x91, 0x33, 0xd7, 0x55, 0x9f, 0x50,
	0x7b, 0x36, 0x8d, 0xf1, 0xf8, 0x4b, 0x1f, 0xc8, 0x55, 0x85, 0x4d, 0xbf, 0xaf, 0xd1, 0x78, 0x30,
	0x7e, 0xef, 0x7b, 0x0b, 0x6c, 0x16, 0x3b, 0x97, 0x01, 0x6c, 0x53, 0xfe, 0xde, 0x5b, 0x08, 0x9b,
	0x43, 0xc4, 0xe6, 0x69, 0x73, 0x98, 0x9f, 0x5e, 0xc9, 0xa2, 0x1c, 0x84, 0xd0, 0x2d, 0xc9, 0xde,
	0xc0, 0x34, 0xe9, 0xa9, 0xf7, 0xfc, 0x1f, 0x6f, 0x3c, 0xf2, 0xfc, 0x9a, 0x27, 0x43, 0x7c, 0xd4,
	0xc7, 0xcb, 0xa0, 0xef, 0x85, 0x7c, 0x2c, 0xc6, 0x66, 0x75, 0x1a, 0xd9, 0xf3, 0x4f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x94, 0xf7, 0xab, 0xfd, 0x93, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EmailingClient is the client API for Emailing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EmailingClient interface {
	SendEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*empty.Empty, error)
}

type emailingClient struct {
	cc *grpc.ClientConn
}

func NewEmailingClient(cc *grpc.ClientConn) EmailingClient {
	return &emailingClient{cc}
}

func (c *emailingClient) SendEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.messaging.emailing.Emailing/SendEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailingServer is the server API for Emailing service.
type EmailingServer interface {
	SendEmail(context.Context, *Email) (*empty.Empty, error)
}

func RegisterEmailingServer(s *grpc.Server, srv EmailingServer) {
	s.RegisterService(&_Emailing_serviceDesc, srv)
}

func _Emailing_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailingServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.messaging.emailing.Emailing/SendEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailingServer).SendEmail(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

var _Emailing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "umrs.messaging.emailing.Emailing",
	HandlerType: (*EmailingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _Emailing_SendEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emailing.proto",
}