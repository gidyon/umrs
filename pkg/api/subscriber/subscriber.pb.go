// Code generated by protoc-gen-go. DO NOT EDIT.
// source: subscriber.proto

package subscriber

import (
	context "context"
	fmt "fmt"
	"github.com/gidyon/umrs/pkg/api/notification"
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

// Subscriber represent an entity that can subscribes to channels
type Subscriber struct {
	AccountId            string                  `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Email                string                  `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Phone                string                  `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	SendMethod           notification.SendMethod `protobuf:"varint,4,opt,name=send_method,json=sendMethod,proto3,enum=umrs.notification.SendMethod" json:"send_method,omitempty"`
	Channels             []string                `protobuf:"bytes,5,rep,name=channels,proto3" json:"channels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Subscriber) Reset()         { *m = Subscriber{} }
func (m *Subscriber) String() string { return proto.CompactTextString(m) }
func (*Subscriber) ProtoMessage()    {}
func (*Subscriber) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{0}
}

func (m *Subscriber) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subscriber.Unmarshal(m, b)
}
func (m *Subscriber) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subscriber.Marshal(b, m, deterministic)
}
func (m *Subscriber) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subscriber.Merge(m, src)
}
func (m *Subscriber) XXX_Size() int {
	return xxx_messageInfo_Subscriber.Size(m)
}
func (m *Subscriber) XXX_DiscardUnknown() {
	xxx_messageInfo_Subscriber.DiscardUnknown(m)
}

var xxx_messageInfo_Subscriber proto.InternalMessageInfo

func (m *Subscriber) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *Subscriber) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Subscriber) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Subscriber) GetSendMethod() notification.SendMethod {
	if m != nil {
		return m.SendMethod
	}
	return notification.SendMethod_EMAIL
}

func (m *Subscriber) GetChannels() []string {
	if m != nil {
		return m.Channels
	}
	return nil
}

// SubscriberRequest request to subscribes a user to a channel
type SubscriberRequest struct {
	AccountId            string   `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscriberRequest) Reset()         { *m = SubscriberRequest{} }
func (m *SubscriberRequest) String() string { return proto.CompactTextString(m) }
func (*SubscriberRequest) ProtoMessage()    {}
func (*SubscriberRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{1}
}

func (m *SubscriberRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscriberRequest.Unmarshal(m, b)
}
func (m *SubscriberRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscriberRequest.Marshal(b, m, deterministic)
}
func (m *SubscriberRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscriberRequest.Merge(m, src)
}
func (m *SubscriberRequest) XXX_Size() int {
	return xxx_messageInfo_SubscriberRequest.Size(m)
}
func (m *SubscriberRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscriberRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscriberRequest proto.InternalMessageInfo

func (m *SubscriberRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *SubscriberRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

// ListSubscribersRequest is request to list subscribers for a particular channel
type ListSubscribersRequest struct {
	Channel              string   `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	PageToken            int32    `protobuf:"varint,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	PageSize             int32    `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListSubscribersRequest) Reset()         { *m = ListSubscribersRequest{} }
func (m *ListSubscribersRequest) String() string { return proto.CompactTextString(m) }
func (*ListSubscribersRequest) ProtoMessage()    {}
func (*ListSubscribersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{2}
}

func (m *ListSubscribersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListSubscribersRequest.Unmarshal(m, b)
}
func (m *ListSubscribersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListSubscribersRequest.Marshal(b, m, deterministic)
}
func (m *ListSubscribersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListSubscribersRequest.Merge(m, src)
}
func (m *ListSubscribersRequest) XXX_Size() int {
	return xxx_messageInfo_ListSubscribersRequest.Size(m)
}
func (m *ListSubscribersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListSubscribersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListSubscribersRequest proto.InternalMessageInfo

func (m *ListSubscribersRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *ListSubscribersRequest) GetPageToken() int32 {
	if m != nil {
		return m.PageToken
	}
	return 0
}

func (m *ListSubscribersRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

// ListSubscribersResponse is response containing collection of subscribers
type ListSubscribersResponse struct {
	Subscribers          []*Subscriber `protobuf:"bytes,1,rep,name=subscribers,proto3" json:"subscribers,omitempty"`
	NextPageToken        int32         `protobuf:"varint,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ListSubscribersResponse) Reset()         { *m = ListSubscribersResponse{} }
func (m *ListSubscribersResponse) String() string { return proto.CompactTextString(m) }
func (*ListSubscribersResponse) ProtoMessage()    {}
func (*ListSubscribersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{3}
}

func (m *ListSubscribersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListSubscribersResponse.Unmarshal(m, b)
}
func (m *ListSubscribersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListSubscribersResponse.Marshal(b, m, deterministic)
}
func (m *ListSubscribersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListSubscribersResponse.Merge(m, src)
}
func (m *ListSubscribersResponse) XXX_Size() int {
	return xxx_messageInfo_ListSubscribersResponse.Size(m)
}
func (m *ListSubscribersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListSubscribersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListSubscribersResponse proto.InternalMessageInfo

func (m *ListSubscribersResponse) GetSubscribers() []*Subscriber {
	if m != nil {
		return m.Subscribers
	}
	return nil
}

func (m *ListSubscribersResponse) GetNextPageToken() int32 {
	if m != nil {
		return m.NextPageToken
	}
	return 0
}

// GetSubscriberRequest is request to retrieve a subscriber
type GetSubscriberRequest struct {
	AccountId            string   `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSubscriberRequest) Reset()         { *m = GetSubscriberRequest{} }
func (m *GetSubscriberRequest) String() string { return proto.CompactTextString(m) }
func (*GetSubscriberRequest) ProtoMessage()    {}
func (*GetSubscriberRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{4}
}

func (m *GetSubscriberRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSubscriberRequest.Unmarshal(m, b)
}
func (m *GetSubscriberRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSubscriberRequest.Marshal(b, m, deterministic)
}
func (m *GetSubscriberRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSubscriberRequest.Merge(m, src)
}
func (m *GetSubscriberRequest) XXX_Size() int {
	return xxx_messageInfo_GetSubscriberRequest.Size(m)
}
func (m *GetSubscriberRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSubscriberRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSubscriberRequest proto.InternalMessageInfo

func (m *GetSubscriberRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

// GetSendMethodRequest is request to get the send method of a user
type GetSendMethodRequest struct {
	AccountId            string   `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSendMethodRequest) Reset()         { *m = GetSendMethodRequest{} }
func (m *GetSendMethodRequest) String() string { return proto.CompactTextString(m) }
func (*GetSendMethodRequest) ProtoMessage()    {}
func (*GetSendMethodRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{5}
}

func (m *GetSendMethodRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSendMethodRequest.Unmarshal(m, b)
}
func (m *GetSendMethodRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSendMethodRequest.Marshal(b, m, deterministic)
}
func (m *GetSendMethodRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSendMethodRequest.Merge(m, src)
}
func (m *GetSendMethodRequest) XXX_Size() int {
	return xxx_messageInfo_GetSendMethodRequest.Size(m)
}
func (m *GetSendMethodRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSendMethodRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSendMethodRequest proto.InternalMessageInfo

func (m *GetSendMethodRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

// GetSendMethodResponse is response containing the send method
type GetSendMethodResponse struct {
	SendMethod           notification.SendMethod `protobuf:"varint,1,opt,name=send_method,json=sendMethod,proto3,enum=umrs.notification.SendMethod" json:"send_method,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *GetSendMethodResponse) Reset()         { *m = GetSendMethodResponse{} }
func (m *GetSendMethodResponse) String() string { return proto.CompactTextString(m) }
func (*GetSendMethodResponse) ProtoMessage()    {}
func (*GetSendMethodResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a004d881b890f1e0, []int{6}
}

func (m *GetSendMethodResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSendMethodResponse.Unmarshal(m, b)
}
func (m *GetSendMethodResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSendMethodResponse.Marshal(b, m, deterministic)
}
func (m *GetSendMethodResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSendMethodResponse.Merge(m, src)
}
func (m *GetSendMethodResponse) XXX_Size() int {
	return xxx_messageInfo_GetSendMethodResponse.Size(m)
}
func (m *GetSendMethodResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSendMethodResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSendMethodResponse proto.InternalMessageInfo

func (m *GetSendMethodResponse) GetSendMethod() notification.SendMethod {
	if m != nil {
		return m.SendMethod
	}
	return notification.SendMethod_EMAIL
}

func init() {
	proto.RegisterType((*Subscriber)(nil), "umrs.notification.subscriber.Subscriber")
	proto.RegisterType((*SubscriberRequest)(nil), "umrs.notification.subscriber.SubscriberRequest")
	proto.RegisterType((*ListSubscribersRequest)(nil), "umrs.notification.subscriber.ListSubscribersRequest")
	proto.RegisterType((*ListSubscribersResponse)(nil), "umrs.notification.subscriber.ListSubscribersResponse")
	proto.RegisterType((*GetSubscriberRequest)(nil), "umrs.notification.subscriber.GetSubscriberRequest")
	proto.RegisterType((*GetSendMethodRequest)(nil), "umrs.notification.subscriber.GetSendMethodRequest")
	proto.RegisterType((*GetSendMethodResponse)(nil), "umrs.notification.subscriber.GetSendMethodResponse")
}

func init() { proto.RegisterFile("subscriber.proto", fileDescriptor_a004d881b890f1e0) }

var fileDescriptor_a004d881b890f1e0 = []byte{
	// 599 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x41, 0x6b, 0x13, 0x41,
	0x14, 0xc7, 0x99, 0xc6, 0xd4, 0xe6, 0xc5, 0x58, 0x1d, 0x6a, 0x5d, 0xb6, 0x2d, 0xc6, 0x41, 0x4a,
	0x08, 0xed, 0x2e, 0x24, 0xe6, 0x60, 0x05, 0x41, 0x41, 0xa4, 0x52, 0xa1, 0x6c, 0x14, 0xc1, 0x4b,
	0xd8, 0x6c, 0x26, 0x9b, 0xc5, 0x64, 0x66, 0xcd, 0xcc, 0xa2, 0xb6, 0x14, 0x44, 0xbc, 0x7a, 0xf2,
	0xe4, 0x47, 0xf0, 0xdc, 0xa3, 0x1f, 0xc3, 0xaf, 0xe0, 0x07, 0x91, 0x9d, 0xdd, 0xcd, 0x6e, 0x92,
	0x25, 0x4d, 0xc0, 0x53, 0xe6, 0xbd, 0x79, 0xff, 0xd9, 0xdf, 0xbc, 0x79, 0xff, 0xc0, 0x2d, 0x11,
	0x74, 0x85, 0x33, 0xf6, 0xba, 0x74, 0x6c, 0xf8, 0x63, 0x2e, 0x39, 0xde, 0x1d, 0x39, 0x7d, 0xdf,
	0x60, 0x5c, 0x7a, 0x7d, 0xcf, 0xb1, 0xa5, 0xc7, 0x99, 0x91, 0xd6, 0xe8, 0x3b, 0x2e, 0xe7, 0xee,
	0x90, 0x9a, 0xaa, 0xb6, 0x1b, 0xf4, 0x4d, 0x3a, 0xf2, 0xe5, 0xe7, 0x48, 0xaa, 0xef, 0xc6, 0x9b,
	0xb6, 0xef, 0x99, 0x36, 0x63, 0x5c, 0x2a, 0xb9, 0x88, 0x77, 0x0f, 0xd4, 0x8f, 0x73, 0xe8, 0x52,
	0x76, 0x28, 0x3e, 0xda, 0xae, 0x4b, 0xc7, 0x26, 0xf7, 0x55, 0x45, 0x4e, 0x35, 0x9e, 0x22, 0x50,
	0x39, 0x72, 0x89, 0x00, 0xda, 0x13, 0x16, 0xbc, 0x07, 0x60, 0x3b, 0x0e, 0x0f, 0x98, 0xec, 0x78,
	0x3d, 0x0d, 0x55, 0x51, 0xad, 0x64, 0x95, 0xe2, 0xcc, 0x71, 0x0f, 0x6f, 0x41, 0x91, 0x8e, 0x6c,
	0x6f, 0xa8, 0xad, 0xa9, 0x9d, 0x28, 0x08, 0xb3, 0xfe, 0x80, 0x33, 0xaa, 0x15, 0xa2, 0xac, 0x0a,
	0xf0, 0x13, 0x28, 0x0b, 0xca, 0x7a, 0x9d, 0x11, 0x95, 0x03, 0xde, 0xd3, 0xae, 0x55, 0x51, 0xed,
	0x66, 0x63, 0xcf, 0x98, 0x6f, 0x45, 0x9b, 0xb2, 0xde, 0x2b, 0x55, 0x64, 0x81, 0x98, 0xac, 0xb1,
	0x0e, 0x1b, 0xce, 0xc0, 0x66, 0x8c, 0x0e, 0x85, 0x56, 0xac, 0x16, 0x6a, 0x25, 0x6b, 0x12, 0x93,
	0x13, 0xb8, 0x9d, 0x42, 0x5b, 0xf4, 0x43, 0x40, 0x85, 0xbc, 0x8a, 0x5d, 0x83, 0xeb, 0xb1, 0x3e,
	0xa6, 0x4f, 0x42, 0xc2, 0x60, 0xfb, 0xc4, 0x13, 0x32, 0x3d, 0x51, 0x24, 0x47, 0x66, 0x34, 0x68,
	0x4a, 0x13, 0x7e, 0xcc, 0xb7, 0x5d, 0xda, 0x91, 0xfc, 0x3d, 0x65, 0xea, 0xc0, 0xa2, 0x55, 0x0a,
	0x33, 0xaf, 0xc3, 0x04, 0xde, 0x01, 0x15, 0x74, 0x84, 0x77, 0x16, 0xb5, 0xa5, 0x68, 0x6d, 0x84,
	0x89, 0xb6, 0x77, 0x46, 0xc9, 0x77, 0x04, 0x77, 0xe7, 0x3e, 0x28, 0x7c, 0xce, 0x04, 0xc5, 0x2f,
	0xa1, 0x9c, 0x8e, 0x86, 0xd0, 0x50, 0xb5, 0x50, 0x2b, 0x37, 0x6a, 0xc6, 0xa2, 0x01, 0x32, 0x32,
	0xad, 0xc8, 0x8a, 0xf1, 0x3e, 0x6c, 0x32, 0xfa, 0x49, 0x76, 0xe6, 0x40, 0x2b, 0x61, 0xfa, 0x34,
	0x81, 0x25, 0x2d, 0xd8, 0x7a, 0x41, 0xe5, 0xaa, 0x0d, 0x4d, 0x64, 0xe9, 0xeb, 0x2d, 0x27, 0x7b,
	0x0b, 0x77, 0x66, 0x64, 0xf1, 0xd5, 0x67, 0x06, 0x06, 0xad, 0x38, 0x30, 0x8d, 0x5f, 0xeb, 0x50,
	0x49, 0x2f, 0xf1, 0xf4, 0xf4, 0x18, 0x7f, 0x41, 0x50, 0x9a, 0x64, 0xb0, 0xb9, 0x74, 0x17, 0xa3,
	0x8b, 0xe8, 0xdb, 0x46, 0x64, 0x3e, 0x23, 0x71, 0xa6, 0xf1, 0x3c, 0x74, 0x26, 0xa9, 0x7f, 0xfd,
	0xf3, 0xf7, 0xc7, 0xda, 0x03, 0x72, 0x4f, 0xb9, 0x32, 0x3c, 0xd4, 0xcc, 0xb4, 0x3d, 0x5d, 0x1f,
	0xa1, 0x3a, 0xfe, 0x86, 0xa0, 0xfc, 0x86, 0x89, 0xff, 0x0f, 0x71, 0xa0, 0x20, 0xf6, 0xc9, 0xfd,
	0x7c, 0x88, 0x80, 0x4d, 0x61, 0xfc, 0x46, 0xb0, 0x39, 0x33, 0x72, 0xf8, 0xe1, 0x62, 0x94, 0x7c,
	0x4b, 0xe8, 0xad, 0x15, 0x55, 0xd1, 0xe3, 0x92, 0xc7, 0x0a, 0xb7, 0x85, 0x9b, 0xf9, 0xb8, 0xb1,
	0xad, 0xcc, 0xf3, 0x78, 0x71, 0x91, 0xdd, 0xc5, 0x3f, 0x11, 0x54, 0xa6, 0x26, 0x14, 0x37, 0x16,
	0x53, 0xe4, 0x8d, 0xb3, 0xbe, 0xb4, 0x8b, 0x92, 0x07, 0xc6, 0x24, 0x1f, 0xf6, 0x3c, 0x1d, 0xef,
	0x0b, 0x7c, 0x19, 0xb3, 0xa5, 0x7f, 0x5c, 0x4b, 0xb0, 0xcd, 0x7a, 0x46, 0x6f, 0xae, 0xa4, 0x89,
	0x7b, 0xfa, 0x48, 0x61, 0x36, 0x89, 0x71, 0x35, 0xa6, 0x99, 0x71, 0xd6, 0x11, 0xaa, 0x3f, 0xbb,
	0xf1, 0x0e, 0xd2, 0xd2, 0xee, 0xba, 0x9a, 0xad, 0xe6, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01,
	0x37, 0xe2, 0xb2, 0xba, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SubscriberAPIClient is the client API for SubscriberAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SubscriberAPIClient interface {
	// Subscribes a user to a channel
	Subscribe(ctx context.Context, in *SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Unsubscribes a user from a channel
	Unsubscribe(ctx context.Context, in *SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Lists subscribers for a channel
	ListSubscribers(ctx context.Context, in *ListSubscribersRequest, opts ...grpc.CallOption) (*ListSubscribersResponse, error)
	// GetSubscriber retrieves information about a single subscriber
	GetSubscriber(ctx context.Context, in *GetSubscriberRequest, opts ...grpc.CallOption) (*Subscriber, error)
	// Retrieve send method for an account
	GetSendMethod(ctx context.Context, in *GetSendMethodRequest, opts ...grpc.CallOption) (*GetSendMethodResponse, error)
}

type subscriberAPIClient struct {
	cc *grpc.ClientConn
}

func NewSubscriberAPIClient(cc *grpc.ClientConn) SubscriberAPIClient {
	return &subscriberAPIClient{cc}
}

func (c *subscriberAPIClient) Subscribe(ctx context.Context, in *SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.notification.subscriber.SubscriberAPI/Subscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriberAPIClient) Unsubscribe(ctx context.Context, in *SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.notification.subscriber.SubscriberAPI/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriberAPIClient) ListSubscribers(ctx context.Context, in *ListSubscribersRequest, opts ...grpc.CallOption) (*ListSubscribersResponse, error) {
	out := new(ListSubscribersResponse)
	err := c.cc.Invoke(ctx, "/umrs.notification.subscriber.SubscriberAPI/ListSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriberAPIClient) GetSubscriber(ctx context.Context, in *GetSubscriberRequest, opts ...grpc.CallOption) (*Subscriber, error) {
	out := new(Subscriber)
	err := c.cc.Invoke(ctx, "/umrs.notification.subscriber.SubscriberAPI/GetSubscriber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriberAPIClient) GetSendMethod(ctx context.Context, in *GetSendMethodRequest, opts ...grpc.CallOption) (*GetSendMethodResponse, error) {
	out := new(GetSendMethodResponse)
	err := c.cc.Invoke(ctx, "/umrs.notification.subscriber.SubscriberAPI/GetSendMethod", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubscriberAPIServer is the server API for SubscriberAPI service.
type SubscriberAPIServer interface {
	// Subscribes a user to a channel
	Subscribe(context.Context, *SubscriberRequest) (*empty.Empty, error)
	// Unsubscribes a user from a channel
	Unsubscribe(context.Context, *SubscriberRequest) (*empty.Empty, error)
	// Lists subscribers for a channel
	ListSubscribers(context.Context, *ListSubscribersRequest) (*ListSubscribersResponse, error)
	// GetSubscriber retrieves information about a single subscriber
	GetSubscriber(context.Context, *GetSubscriberRequest) (*Subscriber, error)
	// Retrieve send method for an account
	GetSendMethod(context.Context, *GetSendMethodRequest) (*GetSendMethodResponse, error)
}

func RegisterSubscriberAPIServer(s *grpc.Server, srv SubscriberAPIServer) {
	s.RegisterService(&_SubscriberAPI_serviceDesc, srv)
}

func _SubscriberAPI_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscriberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberAPIServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.subscriber.SubscriberAPI/Subscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberAPIServer).Subscribe(ctx, req.(*SubscriberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscriberAPI_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscriberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberAPIServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.subscriber.SubscriberAPI/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberAPIServer).Unsubscribe(ctx, req.(*SubscriberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscriberAPI_ListSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSubscribersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberAPIServer).ListSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.subscriber.SubscriberAPI/ListSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberAPIServer).ListSubscribers(ctx, req.(*ListSubscribersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscriberAPI_GetSubscriber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubscriberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberAPIServer).GetSubscriber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.subscriber.SubscriberAPI/GetSubscriber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberAPIServer).GetSubscriber(ctx, req.(*GetSubscriberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscriberAPI_GetSendMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSendMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberAPIServer).GetSendMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.subscriber.SubscriberAPI/GetSendMethod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberAPIServer).GetSendMethod(ctx, req.(*GetSendMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SubscriberAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "umrs.notification.subscriber.SubscriberAPI",
	HandlerType: (*SubscriberAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Subscribe",
			Handler:    _SubscriberAPI_Subscribe_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _SubscriberAPI_Unsubscribe_Handler,
		},
		{
			MethodName: "ListSubscribers",
			Handler:    _SubscriberAPI_ListSubscribers_Handler,
		},
		{
			MethodName: "GetSubscriber",
			Handler:    _SubscriberAPI_GetSubscriber_Handler,
		},
		{
			MethodName: "GetSendMethod",
			Handler:    _SubscriberAPI_GetSendMethod_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "subscriber.proto",
}
