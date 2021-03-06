// Code generated by protoc-gen-go. DO NOT EDIT.
// source: channel.proto

package channel

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

// Channel is a publisher network with subscribers
type Channel struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	OwnerId              string   `protobuf:"bytes,4,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	CreatedTime          string   `protobuf:"bytes,5,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	Subscribers          int32    `protobuf:"varint,6,opt,name=subscribers,proto3" json:"subscribers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Channel) Reset()         { *m = Channel{} }
func (m *Channel) String() string { return proto.CompactTextString(m) }
func (*Channel) ProtoMessage()    {}
func (*Channel) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{0}
}

func (m *Channel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Channel.Unmarshal(m, b)
}
func (m *Channel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Channel.Marshal(b, m, deterministic)
}
func (m *Channel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Channel.Merge(m, src)
}
func (m *Channel) XXX_Size() int {
	return xxx_messageInfo_Channel.Size(m)
}
func (m *Channel) XXX_DiscardUnknown() {
	xxx_messageInfo_Channel.DiscardUnknown(m)
}

var xxx_messageInfo_Channel proto.InternalMessageInfo

func (m *Channel) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *Channel) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Channel) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Channel) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *Channel) GetCreatedTime() string {
	if m != nil {
		return m.CreatedTime
	}
	return ""
}

func (m *Channel) GetSubscribers() int32 {
	if m != nil {
		return m.Subscribers
	}
	return 0
}

// CreateChannelRequest is request to create a new channel
type CreateChannelRequest struct {
	Channel              *Channel `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChannelRequest) Reset()         { *m = CreateChannelRequest{} }
func (m *CreateChannelRequest) String() string { return proto.CompactTextString(m) }
func (*CreateChannelRequest) ProtoMessage()    {}
func (*CreateChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{1}
}

func (m *CreateChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChannelRequest.Unmarshal(m, b)
}
func (m *CreateChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChannelRequest.Marshal(b, m, deterministic)
}
func (m *CreateChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChannelRequest.Merge(m, src)
}
func (m *CreateChannelRequest) XXX_Size() int {
	return xxx_messageInfo_CreateChannelRequest.Size(m)
}
func (m *CreateChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChannelRequest proto.InternalMessageInfo

func (m *CreateChannelRequest) GetChannel() *Channel {
	if m != nil {
		return m.Channel
	}
	return nil
}

// CreateChannelResponse is response after creating a channel
type CreateChannelResponse struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChannelResponse) Reset()         { *m = CreateChannelResponse{} }
func (m *CreateChannelResponse) String() string { return proto.CompactTextString(m) }
func (*CreateChannelResponse) ProtoMessage()    {}
func (*CreateChannelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{2}
}

func (m *CreateChannelResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChannelResponse.Unmarshal(m, b)
}
func (m *CreateChannelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChannelResponse.Marshal(b, m, deterministic)
}
func (m *CreateChannelResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChannelResponse.Merge(m, src)
}
func (m *CreateChannelResponse) XXX_Size() int {
	return xxx_messageInfo_CreateChannelResponse.Size(m)
}
func (m *CreateChannelResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChannelResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChannelResponse proto.InternalMessageInfo

func (m *CreateChannelResponse) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

// UpdateChannelRequest request to update a channel resource
type UpdateChannelRequest struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Channel              *Channel `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateChannelRequest) Reset()         { *m = UpdateChannelRequest{} }
func (m *UpdateChannelRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateChannelRequest) ProtoMessage()    {}
func (*UpdateChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{3}
}

func (m *UpdateChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateChannelRequest.Unmarshal(m, b)
}
func (m *UpdateChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateChannelRequest.Marshal(b, m, deterministic)
}
func (m *UpdateChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateChannelRequest.Merge(m, src)
}
func (m *UpdateChannelRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateChannelRequest.Size(m)
}
func (m *UpdateChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateChannelRequest proto.InternalMessageInfo

func (m *UpdateChannelRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *UpdateChannelRequest) GetChannel() *Channel {
	if m != nil {
		return m.Channel
	}
	return nil
}

// DeleteChannelRequest is request to delete a channel resource
type DeleteChannelRequest struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteChannelRequest) Reset()         { *m = DeleteChannelRequest{} }
func (m *DeleteChannelRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteChannelRequest) ProtoMessage()    {}
func (*DeleteChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{4}
}

func (m *DeleteChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteChannelRequest.Unmarshal(m, b)
}
func (m *DeleteChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteChannelRequest.Marshal(b, m, deterministic)
}
func (m *DeleteChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteChannelRequest.Merge(m, src)
}
func (m *DeleteChannelRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteChannelRequest.Size(m)
}
func (m *DeleteChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteChannelRequest proto.InternalMessageInfo

func (m *DeleteChannelRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

// ListChannelsRequest is request to retrive a collection of channels resource
type ListChannelsRequest struct {
	PageToken            int32    `protobuf:"varint,1,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	PageSize             int32    `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListChannelsRequest) Reset()         { *m = ListChannelsRequest{} }
func (m *ListChannelsRequest) String() string { return proto.CompactTextString(m) }
func (*ListChannelsRequest) ProtoMessage()    {}
func (*ListChannelsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{5}
}

func (m *ListChannelsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListChannelsRequest.Unmarshal(m, b)
}
func (m *ListChannelsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListChannelsRequest.Marshal(b, m, deterministic)
}
func (m *ListChannelsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListChannelsRequest.Merge(m, src)
}
func (m *ListChannelsRequest) XXX_Size() int {
	return xxx_messageInfo_ListChannelsRequest.Size(m)
}
func (m *ListChannelsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListChannelsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListChannelsRequest proto.InternalMessageInfo

func (m *ListChannelsRequest) GetPageToken() int32 {
	if m != nil {
		return m.PageToken
	}
	return 0
}

func (m *ListChannelsRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

// ListChannelsResponse is response containing collection of channel resource
type ListChannelsResponse struct {
	Channels             []*Channel `protobuf:"bytes,1,rep,name=channels,proto3" json:"channels,omitempty"`
	NextPageToken        int32      `protobuf:"varint,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ListChannelsResponse) Reset()         { *m = ListChannelsResponse{} }
func (m *ListChannelsResponse) String() string { return proto.CompactTextString(m) }
func (*ListChannelsResponse) ProtoMessage()    {}
func (*ListChannelsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{6}
}

func (m *ListChannelsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListChannelsResponse.Unmarshal(m, b)
}
func (m *ListChannelsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListChannelsResponse.Marshal(b, m, deterministic)
}
func (m *ListChannelsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListChannelsResponse.Merge(m, src)
}
func (m *ListChannelsResponse) XXX_Size() int {
	return xxx_messageInfo_ListChannelsResponse.Size(m)
}
func (m *ListChannelsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListChannelsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListChannelsResponse proto.InternalMessageInfo

func (m *ListChannelsResponse) GetChannels() []*Channel {
	if m != nil {
		return m.Channels
	}
	return nil
}

func (m *ListChannelsResponse) GetNextPageToken() int32 {
	if m != nil {
		return m.NextPageToken
	}
	return 0
}

// GetChannelRequest is request to retrieve a channel resource
type GetChannelRequest struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetChannelRequest) Reset()         { *m = GetChannelRequest{} }
func (m *GetChannelRequest) String() string { return proto.CompactTextString(m) }
func (*GetChannelRequest) ProtoMessage()    {}
func (*GetChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8f385724121f37b, []int{7}
}

func (m *GetChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetChannelRequest.Unmarshal(m, b)
}
func (m *GetChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetChannelRequest.Marshal(b, m, deterministic)
}
func (m *GetChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetChannelRequest.Merge(m, src)
}
func (m *GetChannelRequest) XXX_Size() int {
	return xxx_messageInfo_GetChannelRequest.Size(m)
}
func (m *GetChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetChannelRequest proto.InternalMessageInfo

func (m *GetChannelRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func init() {
	proto.RegisterType((*Channel)(nil), "umrs.notification.channel.Channel")
	proto.RegisterType((*CreateChannelRequest)(nil), "umrs.notification.channel.CreateChannelRequest")
	proto.RegisterType((*CreateChannelResponse)(nil), "umrs.notification.channel.CreateChannelResponse")
	proto.RegisterType((*UpdateChannelRequest)(nil), "umrs.notification.channel.UpdateChannelRequest")
	proto.RegisterType((*DeleteChannelRequest)(nil), "umrs.notification.channel.DeleteChannelRequest")
	proto.RegisterType((*ListChannelsRequest)(nil), "umrs.notification.channel.ListChannelsRequest")
	proto.RegisterType((*ListChannelsResponse)(nil), "umrs.notification.channel.ListChannelsResponse")
	proto.RegisterType((*GetChannelRequest)(nil), "umrs.notification.channel.GetChannelRequest")
}

func init() { proto.RegisterFile("channel.proto", fileDescriptor_c8f385724121f37b) }

var fileDescriptor_c8f385724121f37b = []byte{
	// 599 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x96, 0x53, 0xdc, 0x24, 0xd3, 0x46, 0x88, 0x25, 0x20, 0xc7, 0x85, 0x92, 0x5a, 0x55, 0xa9,
	0xa2, 0xd6, 0x46, 0x41, 0x70, 0x40, 0x08, 0x09, 0x0a, 0x42, 0x91, 0x38, 0x94, 0x10, 0x2e, 0x5c,
	0x22, 0xc7, 0x9e, 0x9a, 0x15, 0x89, 0xd7, 0x78, 0x37, 0x2a, 0x14, 0x81, 0x04, 0x42, 0xe2, 0x01,
	0x38, 0xc0, 0xe3, 0xf0, 0x0e, 0x48, 0x3c, 0x01, 0x0f, 0x82, 0xbc, 0x5e, 0x53, 0xa7, 0x71, 0xf3,
	0x73, 0xb2, 0x76, 0x66, 0xbe, 0x99, 0x6f, 0x3f, 0x7f, 0x3b, 0x50, 0xf3, 0x5e, 0xbb, 0x61, 0x88,
	0x43, 0x3b, 0x8a, 0x99, 0x60, 0xa4, 0x31, 0x1e, 0xc5, 0xdc, 0x0e, 0x99, 0xa0, 0x47, 0xd4, 0x73,
	0x05, 0x65, 0xa1, 0xad, 0x0a, 0xcc, 0x8d, 0x80, 0xb1, 0x60, 0x88, 0x8e, 0x2c, 0x1c, 0x8c, 0x8f,
	0x1c, 0x1c, 0x45, 0xe2, 0x7d, 0x8a, 0x33, 0xaf, 0xa9, 0xa4, 0x1b, 0x51, 0xc7, 0x0d, 0x43, 0x26,
	0x24, 0x96, 0xab, 0xec, 0x9e, 0xfc, 0x78, 0xfb, 0x01, 0x86, 0xfb, 0xfc, 0xd8, 0x0d, 0x02, 0x8c,
	0x1d, 0x16, 0xc9, 0x8a, 0xe9, 0x6a, 0xeb, 0x97, 0x06, 0xe5, 0x83, 0x74, 0x28, 0xb9, 0x0e, 0xa0,
	0xe6, 0xf7, 0xa9, 0x6f, 0x68, 0x4d, 0x6d, 0xb7, 0xda, 0xad, 0xaa, 0x48, 0xc7, 0x27, 0x75, 0xd0,
	0x05, 0x15, 0x43, 0x34, 0x4a, 0x32, 0x93, 0x1e, 0x48, 0x13, 0xd6, 0x7c, 0xe4, 0x5e, 0x4c, 0xe5,
	0x08, 0x63, 0x45, 0xe6, 0xf2, 0x21, 0xd2, 0x80, 0x0a, 0x3b, 0x0e, 0x31, 0x4e, 0x9a, 0x5e, 0x90,
	0xe9, 0xb2, 0x3c, 0x77, 0x7c, 0xb2, 0x05, 0xeb, 0x5e, 0x8c, 0xae, 0x40, 0xbf, 0x2f, 0xe8, 0x08,
	0x0d, 0x3d, 0x45, 0xab, 0x58, 0x8f, 0x8e, 0x64, 0x7f, 0x3e, 0x1e, 0x24, 0xdd, 0x06, 0x18, 0x73,
	0x63, 0xb5, 0xa9, 0xed, 0xea, 0xdd, 0x7c, 0xc8, 0xea, 0x41, 0xfd, 0x40, 0x02, 0xd4, 0x3d, 0xba,
	0xf8, 0x76, 0x8c, 0x5c, 0x90, 0xfb, 0x50, 0x56, 0xe4, 0xe5, 0x5d, 0xd6, 0xda, 0x96, 0x7d, 0xae,
	0xe0, 0x76, 0x86, 0xcd, 0x20, 0xd6, 0x5d, 0xb8, 0x72, 0xa6, 0x2b, 0x8f, 0x58, 0xc8, 0x71, 0x8e,
	0x4a, 0x16, 0x87, 0xfa, 0xcb, 0xc8, 0x9f, 0x66, 0x33, 0x47, 0xdc, 0x1c, 0xd9, 0xd2, 0xf2, 0x64,
	0xef, 0x40, 0xfd, 0x31, 0x0e, 0x71, 0xc9, 0xa1, 0xd6, 0x73, 0xb8, 0xfc, 0x8c, 0x72, 0xa1, 0x40,
	0x3c, 0x87, 0x8a, 0xdc, 0x00, 0xfb, 0x82, 0xbd, 0xc1, 0x50, 0xa2, 0xf4, 0x6e, 0x35, 0x89, 0xf4,
	0x92, 0x00, 0xd9, 0x00, 0x79, 0xe8, 0x73, 0x7a, 0x92, 0x7a, 0x41, 0xef, 0x56, 0x92, 0xc0, 0x0b,
	0x7a, 0x82, 0xd6, 0x27, 0xa8, 0x4f, 0xb6, 0x54, 0xaa, 0x3d, 0x80, 0x8a, 0x9a, 0xcb, 0x0d, 0xad,
	0xb9, 0xb2, 0xe0, 0x05, 0xff, 0x63, 0xc8, 0x0e, 0x5c, 0x0c, 0xf1, 0x9d, 0xe8, 0xe7, 0x88, 0xa5,
	0xa3, 0x6b, 0x49, 0xf8, 0x30, 0x23, 0x67, 0xb5, 0xe1, 0xd2, 0x53, 0x14, 0x4b, 0xc9, 0xd0, 0xfe,
	0xa3, 0x03, 0x28, 0xc4, 0xc3, 0xc3, 0x0e, 0xf9, 0xa9, 0x41, 0x6d, 0xe2, 0xd7, 0x13, 0x67, 0x16,
	0xd5, 0x02, 0xeb, 0x99, 0xb7, 0x16, 0x07, 0xa4, 0xfa, 0x58, 0xdb, 0x5f, 0x7e, 0xff, 0xfd, 0x5e,
	0xda, 0xb4, 0x1a, 0xf2, 0x55, 0x27, 0x68, 0x27, 0xbb, 0xbb, 0x93, 0x3e, 0x87, 0x7b, 0x5a, 0x8b,
	0x7c, 0xd5, 0xa0, 0x36, 0xe1, 0xae, 0x99, 0xd4, 0x8a, 0x7c, 0x68, 0x5e, 0xb5, 0xd3, 0xed, 0x61,
	0x67, 0xab, 0xc5, 0x7e, 0x92, 0xac, 0x16, 0xab, 0x25, 0x09, 0x6c, 0x9b, 0x37, 0x0a, 0x08, 0x7c,
	0x38, 0x55, 0xef, 0x63, 0x42, 0xe3, 0xb3, 0x06, 0xb5, 0x09, 0xbf, 0xcd, 0xa4, 0x51, 0xe4, 0xcc,
	0x73, 0x69, 0xdc, 0x94, 0x34, 0xb6, 0x5a, 0xf3, 0x68, 0x90, 0x1f, 0x1a, 0xac, 0xe7, 0x9d, 0x46,
	0xec, 0x19, 0x14, 0x0a, 0x5c, 0x6e, 0x3a, 0x0b, 0xd7, 0xab, 0x5f, 0xb4, 0x23, 0xa9, 0x35, 0xc9,
	0x66, 0x01, 0x35, 0xd7, 0x4b, 0xf0, 0xce, 0x90, 0x72, 0x41, 0xbe, 0x69, 0x00, 0xa7, 0x1e, 0x24,
	0x7b, 0x33, 0xe6, 0x4c, 0x59, 0xd5, 0x5c, 0xe0, 0x55, 0x64, 0x1a, 0x91, 0x79, 0x1a, 0x3d, 0xaa,
	0xbe, 0xca, 0x36, 0xc4, 0x60, 0x55, 0xea, 0x7c, 0xfb, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x22,
	0x0c, 0x84, 0x26, 0x83, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChannelAPIClient is the client API for ChannelAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChannelAPIClient interface {
	// Creates a new subscriber channel
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error)
	// Updates an existing channel resource
	UpdateChannel(ctx context.Context, in *UpdateChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Removes a subscribers channel
	DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Retrieves a collection of channels resource
	ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error)
	// Retrieves a single channel resource
	GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*Channel, error)
}

type channelAPIClient struct {
	cc *grpc.ClientConn
}

func NewChannelAPIClient(cc *grpc.ClientConn) ChannelAPIClient {
	return &channelAPIClient{cc}
}

func (c *channelAPIClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error) {
	out := new(CreateChannelResponse)
	err := c.cc.Invoke(ctx, "/umrs.notification.channel.ChannelAPI/CreateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) UpdateChannel(ctx context.Context, in *UpdateChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.notification.channel.ChannelAPI/UpdateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/umrs.notification.channel.ChannelAPI/DeleteChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error) {
	out := new(ListChannelsResponse)
	err := c.cc.Invoke(ctx, "/umrs.notification.channel.ChannelAPI/ListChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*Channel, error) {
	out := new(Channel)
	err := c.cc.Invoke(ctx, "/umrs.notification.channel.ChannelAPI/GetChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelAPIServer is the server API for ChannelAPI service.
type ChannelAPIServer interface {
	// Creates a new subscriber channel
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error)
	// Updates an existing channel resource
	UpdateChannel(context.Context, *UpdateChannelRequest) (*empty.Empty, error)
	// Removes a subscribers channel
	DeleteChannel(context.Context, *DeleteChannelRequest) (*empty.Empty, error)
	// Retrieves a collection of channels resource
	ListChannels(context.Context, *ListChannelsRequest) (*ListChannelsResponse, error)
	// Retrieves a single channel resource
	GetChannel(context.Context, *GetChannelRequest) (*Channel, error)
}

func RegisterChannelAPIServer(s *grpc.Server, srv ChannelAPIServer) {
	s.RegisterService(&_ChannelAPI_serviceDesc, srv)
}

func _ChannelAPI_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.channel.ChannelAPI/CreateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_UpdateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).UpdateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.channel.ChannelAPI/UpdateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).UpdateChannel(ctx, req.(*UpdateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_DeleteChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).DeleteChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.channel.ChannelAPI/DeleteChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).DeleteChannel(ctx, req.(*DeleteChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_ListChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).ListChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.channel.ChannelAPI/ListChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).ListChannels(ctx, req.(*ListChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_GetChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).GetChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/umrs.notification.channel.ChannelAPI/GetChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).GetChannel(ctx, req.(*GetChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChannelAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "umrs.notification.channel.ChannelAPI",
	HandlerType: (*ChannelAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChannel",
			Handler:    _ChannelAPI_CreateChannel_Handler,
		},
		{
			MethodName: "UpdateChannel",
			Handler:    _ChannelAPI_UpdateChannel_Handler,
		},
		{
			MethodName: "DeleteChannel",
			Handler:    _ChannelAPI_DeleteChannel_Handler,
		},
		{
			MethodName: "ListChannels",
			Handler:    _ChannelAPI_ListChannels_Handler,
		},
		{
			MethodName: "GetChannel",
			Handler:    _ChannelAPI_GetChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "channel.proto",
}
