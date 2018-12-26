// Code generated by protoc-gen-go. DO NOT EDIT.
// source: log_entry.proto

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

//
//Represents an entry which is stored in the log.
type LogEntry struct {
	// the index of the entry inside the leader's log
	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// current term of this entry
	// which is the term in which the client requested this command
	Term uint64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	// The content of the entry
	Content              []byte   `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (*LogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_1af096ea28a6a163, []int{0}
}

func (m *LogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntry.Unmarshal(m, b)
}
func (m *LogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntry.Marshal(b, m, deterministic)
}
func (m *LogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntry.Merge(m, src)
}
func (m *LogEntry) XXX_Size() int {
	return xxx_messageInfo_LogEntry.Size(m)
}
func (m *LogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntry proto.InternalMessageInfo

func (m *LogEntry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *LogEntry) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *LogEntry) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "protobuf.LogEntry")
}

func init() { proto.RegisterFile("log_entry.proto", fileDescriptor_1af096ea28a6a163) }

var fileDescriptor_1af096ea28a6a163 = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xc9, 0x4f, 0x8f,
	0x4f, 0xcd, 0x2b, 0x29, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49,
	0xa5, 0x69, 0x4a, 0x7e, 0x5c, 0x1c, 0x3e, 0xf9, 0xe9, 0xae, 0x20, 0x39, 0x21, 0x11, 0x2e, 0xd6,
	0xcc, 0xbc, 0x94, 0xd4, 0x0a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x08, 0x47, 0x48, 0x88,
	0x8b, 0xa5, 0x24, 0xb5, 0x28, 0x57, 0x82, 0x09, 0x2c, 0x08, 0x66, 0x0b, 0x49, 0x70, 0xb1, 0x27,
	0xe7, 0xe7, 0x95, 0xa4, 0xe6, 0x95, 0x48, 0x30, 0x2b, 0x30, 0x6a, 0xf0, 0x04, 0xc1, 0xb8, 0x49,
	0x6c, 0x60, 0x93, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7e, 0x85, 0xb1, 0x7f, 0x73, 0x00,
	0x00, 0x00,
}
