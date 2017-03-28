// Code generated by protoc-gen-go.
// source: controlMsgBody.proto
// DO NOT EDIT!

/*
Package overmind is a generated protocol buffer package.

It is generated from these files:
	controlMsgBody.proto

It has these top-level messages:
	ControlMsgBody
*/
package overmind

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ControlMsgBody struct {
	ControlType      *uint32  `protobuf:"varint,1,opt,name=controlType" json:"controlType,omitempty"`
	Sign             []byte   `protobuf:"bytes,2,opt,name=sign" json:"sign,omitempty"`
	Condition        []byte   `protobuf:"bytes,3,opt,name=condition" json:"condition,omitempty"`
	Data             []byte   `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
	Product          []byte   `protobuf:"bytes,5,opt,name=product" json:"product,omitempty"`
	MultiSigns       []string `protobuf:"bytes,6,rep,name=multiSigns" json:"multiSigns,omitempty"`
	Source           *string  `protobuf:"bytes,7,opt,name=source" json:"source,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ControlMsgBody) Reset()                    { *m = ControlMsgBody{} }
func (m *ControlMsgBody) String() string            { return proto.CompactTextString(m) }
func (*ControlMsgBody) ProtoMessage()               {}
func (*ControlMsgBody) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *ControlMsgBody) GetControlType() uint32 {
	if m != nil && m.ControlType != nil {
		return *m.ControlType
	}
	return 0
}

func (m *ControlMsgBody) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

func (m *ControlMsgBody) GetCondition() []byte {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (m *ControlMsgBody) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ControlMsgBody) GetProduct() []byte {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *ControlMsgBody) GetMultiSigns() []string {
	if m != nil {
		return m.MultiSigns
	}
	return nil
}

func (m *ControlMsgBody) GetSource() string {
	if m != nil && m.Source != nil {
		return *m.Source
	}
	return ""
}

func init() {
	proto.RegisterType((*ControlMsgBody)(nil), "overmind.ControlMsgBody")
}

func init() { proto.RegisterFile("controlMsgBody.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 165 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x49, 0xce, 0xcf, 0x2b,
	0x29, 0xca, 0xcf, 0xf1, 0x2d, 0x4e, 0x77, 0xca, 0x4f, 0xa9, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0xc8, 0x2f, 0x4b, 0x2d, 0xca, 0xcd, 0xcc, 0x4b, 0x51, 0xea, 0x64, 0xe4, 0xe2, 0x73,
	0x46, 0x51, 0x22, 0x24, 0xcc, 0xc5, 0x0d, 0xd5, 0x14, 0x52, 0x59, 0x90, 0x2a, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0x2b, 0xc4, 0xc3, 0xc5, 0x52, 0x9c, 0x99, 0x9e, 0x27, 0xc1, 0x04, 0xe4, 0xf1, 0x08,
	0x09, 0x72, 0x71, 0x02, 0x95, 0xa4, 0x64, 0x96, 0x64, 0xe6, 0xe7, 0x49, 0x30, 0x83, 0x85, 0x80,
	0x0a, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x58, 0xc0, 0x3c, 0x7e, 0x2e, 0x76, 0xa0, 0x4d, 0x29, 0xa5,
	0xc9, 0x25, 0x12, 0xac, 0x60, 0x01, 0x21, 0x2e, 0xae, 0xdc, 0xd2, 0x9c, 0x92, 0xcc, 0x60, 0xa0,
	0x21, 0xc5, 0x12, 0x6c, 0x0a, 0xcc, 0x1a, 0x9c, 0x42, 0x7c, 0x5c, 0x6c, 0xc5, 0xf9, 0xa5, 0x45,
	0xc9, 0xa9, 0x12, 0xec, 0x40, 0x35, 0x9c, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x1d, 0xe4,
	0x63, 0xac, 0x00, 0x00, 0x00,
}