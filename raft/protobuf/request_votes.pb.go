// Code generated by protoc-gen-go. DO NOT EDIT.
// source: request_votes.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type RequestVoteCall struct {
	// current term of the candidate
	Term uint64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	// id of the candidate
	CandidateId string `protobuf:"bytes,2,opt,name=CandidateId,proto3" json:"CandidateId,omitempty"`
	// index of the last log entry
	LastLogIndex uint64 `protobuf:"varint,3,opt,name=LastLogIndex,proto3" json:"LastLogIndex,omitempty"`
	// term of the last log entry
	LastLogTerm          uint64   `protobuf:"varint,4,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteCall) Reset()         { *m = RequestVoteCall{} }
func (m *RequestVoteCall) String() string { return proto.CompactTextString(m) }
func (*RequestVoteCall) ProtoMessage()    {}
func (*RequestVoteCall) Descriptor() ([]byte, []int) {
	return fileDescriptor_5afbe2cbd34751c5, []int{0}
}

func (m *RequestVoteCall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteCall.Unmarshal(m, b)
}
func (m *RequestVoteCall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteCall.Marshal(b, m, deterministic)
}
func (m *RequestVoteCall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteCall.Merge(m, src)
}
func (m *RequestVoteCall) XXX_Size() int {
	return xxx_messageInfo_RequestVoteCall.Size(m)
}
func (m *RequestVoteCall) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteCall.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteCall proto.InternalMessageInfo

func (m *RequestVoteCall) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteCall) GetCandidateId() string {
	if m != nil {
		return m.CandidateId
	}
	return ""
}

func (m *RequestVoteCall) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteCall) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

type RequestVoteRespond struct {
	// current term of the node sending this respond
	Term uint32 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	// true if the node grants the vote for the candidate
	VoteGranted          bool     `protobuf:"varint,2,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteRespond) Reset()         { *m = RequestVoteRespond{} }
func (m *RequestVoteRespond) String() string { return proto.CompactTextString(m) }
func (*RequestVoteRespond) ProtoMessage()    {}
func (*RequestVoteRespond) Descriptor() ([]byte, []int) {
	return fileDescriptor_5afbe2cbd34751c5, []int{1}
}

func (m *RequestVoteRespond) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteRespond.Unmarshal(m, b)
}
func (m *RequestVoteRespond) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteRespond.Marshal(b, m, deterministic)
}
func (m *RequestVoteRespond) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteRespond.Merge(m, src)
}
func (m *RequestVoteRespond) XXX_Size() int {
	return xxx_messageInfo_RequestVoteRespond.Size(m)
}
func (m *RequestVoteRespond) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteRespond.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteRespond proto.InternalMessageInfo

func (m *RequestVoteRespond) GetTerm() uint32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteRespond) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

func init() {
	proto.RegisterType((*RequestVoteCall)(nil), "protobuf.RequestVoteCall")
	proto.RegisterType((*RequestVoteRespond)(nil), "protobuf.RequestVoteRespond")
}

func init() { proto.RegisterFile("request_votes.proto", fileDescriptor_5afbe2cbd34751c5) }

var fileDescriptor_5afbe2cbd34751c5 = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x4a, 0x2d, 0x2c,
	0x4d, 0x2d, 0x2e, 0x89, 0x2f, 0xcb, 0x2f, 0x49, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a, 0xbd, 0x8c, 0x5c, 0xfc, 0x41, 0x10, 0x15, 0x61, 0xf9,
	0x25, 0xa9, 0xce, 0x89, 0x39, 0x39, 0x42, 0x42, 0x5c, 0x2c, 0x21, 0xa9, 0x45, 0xb9, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x60, 0xb6, 0x90, 0x02, 0x17, 0xb7, 0x73, 0x62, 0x5e, 0x4a, 0x66,
	0x4a, 0x62, 0x49, 0xaa, 0x67, 0x8a, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0xb2, 0x90, 0x90,
	0x12, 0x17, 0x8f, 0x4f, 0x62, 0x71, 0x89, 0x4f, 0x7e, 0xba, 0x67, 0x5e, 0x4a, 0x6a, 0x85, 0x04,
	0x33, 0x58, 0x37, 0x8a, 0x18, 0xc8, 0x14, 0x28, 0x1f, 0x6c, 0x01, 0x0b, 0x58, 0x09, 0xb2, 0x90,
	0x92, 0x17, 0x97, 0x10, 0x92, 0x73, 0x82, 0x52, 0x8b, 0x0b, 0xf2, 0xf3, 0x52, 0x50, 0x5c, 0xc4,
	0x8b, 0x70, 0x11, 0x48, 0x89, 0x7b, 0x51, 0x62, 0x5e, 0x49, 0x2a, 0xc4, 0x45, 0x1c, 0x41, 0xc8,
	0x42, 0x49, 0x6c, 0x60, 0x5f, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa7, 0xeb, 0xdc, 0x5a,
	0x03, 0x01, 0x00, 0x00,
}