// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: messages.proto

package qitmeer_p2p_v1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	github_com_prysmaticlabs_go_bitfield "github.com/prysmaticlabs/go-bitfield"
	io "io"
	math "math"
	math_bits "math/bits"
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

type ErrorResponse struct {
	Message              []byte   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty" ssz-max:"256"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorResponse) Reset()         { *m = ErrorResponse{} }
func (m *ErrorResponse) String() string { return proto.CompactTextString(m) }
func (*ErrorResponse) ProtoMessage()    {}
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{0}
}
func (m *ErrorResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ErrorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ErrorResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ErrorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorResponse.Merge(m, src)
}
func (m *ErrorResponse) XXX_Size() int {
	return m.Size()
}
func (m *ErrorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorResponse proto.InternalMessageInfo

func (m *ErrorResponse) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

type MetaData struct {
	SeqNumber            uint64                                           `protobuf:"varint,1,opt,name=seq_number,json=seqNumber,proto3" json:"seq_number,omitempty"`
	Subnets              github_com_prysmaticlabs_go_bitfield.Bitvector64 `protobuf:"bytes,2,opt,name=subnets,proto3,casttype=github.com/prysmaticlabs/go-bitfield.Bitvector64" json:"subnets,omitempty" ssz-size:"8"`
	XXX_NoUnkeyedLiteral struct{}                                         `json:"-"`
	XXX_unrecognized     []byte                                           `json:"-"`
	XXX_sizecache        int32                                            `json:"-"`
}

func (m *MetaData) Reset()         { *m = MetaData{} }
func (m *MetaData) String() string { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()    {}
func (*MetaData) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{1}
}
func (m *MetaData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetaData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetaData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetaData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaData.Merge(m, src)
}
func (m *MetaData) XXX_Size() int {
	return m.Size()
}
func (m *MetaData) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaData.DiscardUnknown(m)
}

var xxx_messageInfo_MetaData proto.InternalMessageInfo

func (m *MetaData) GetSeqNumber() uint64 {
	if m != nil {
		return m.SeqNumber
	}
	return 0
}

func (m *MetaData) GetSubnets() github_com_prysmaticlabs_go_bitfield.Bitvector64 {
	if m != nil {
		return m.Subnets
	}
	return nil
}

type ChainState struct {
	GenesisHash          []byte   `protobuf:"bytes,1,opt,name=genesisHash,proto3" json:"genesisHash,omitempty" ssz-size:"32"`
	ProtocolVersion      uint32   `protobuf:"varint,2,opt,name=protocolVersion,proto3" json:"protocolVersion,omitempty"`
	Timestamp            uint64   `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Services             uint64   `protobuf:"varint,4,opt,name=services,proto3" json:"services,omitempty"`
	GraphState           uint32   `protobuf:"varint,5,opt,name=graphState,proto3" json:"graphState,omitempty"`
	UserAgent            []byte   `protobuf:"bytes,6,opt,name=userAgent,proto3" json:"userAgent,omitempty" ssz-max:"256"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChainState) Reset()         { *m = ChainState{} }
func (m *ChainState) String() string { return proto.CompactTextString(m) }
func (*ChainState) ProtoMessage()    {}
func (*ChainState) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{2}
}
func (m *ChainState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainState.Merge(m, src)
}
func (m *ChainState) XXX_Size() int {
	return m.Size()
}
func (m *ChainState) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainState.DiscardUnknown(m)
}

var xxx_messageInfo_ChainState proto.InternalMessageInfo

func (m *ChainState) GetGenesisHash() []byte {
	if m != nil {
		return m.GenesisHash
	}
	return nil
}

func (m *ChainState) GetProtocolVersion() uint32 {
	if m != nil {
		return m.ProtocolVersion
	}
	return 0
}

func (m *ChainState) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ChainState) GetServices() uint64 {
	if m != nil {
		return m.Services
	}
	return 0
}

func (m *ChainState) GetGraphState() uint32 {
	if m != nil {
		return m.GraphState
	}
	return 0
}

func (m *ChainState) GetUserAgent() []byte {
	if m != nil {
		return m.UserAgent
	}
	return nil
}

func init() {
	proto.RegisterType((*ErrorResponse)(nil), "qitmeer.p2p.v1.ErrorResponse")
	proto.RegisterType((*MetaData)(nil), "qitmeer.p2p.v1.MetaData")
	proto.RegisterType((*ChainState)(nil), "qitmeer.p2p.v1.ChainState")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_4dc296cbfe5ffcd5) }

var fileDescriptor_4dc296cbfe5ffcd5 = []byte{
	// 390 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0xc1, 0xae, 0x12, 0x31,
	0x14, 0x75, 0xf4, 0xf9, 0xde, 0xe3, 0x0a, 0xa8, 0x5d, 0x4d, 0x88, 0x0e, 0xa6, 0x2b, 0x12, 0xc3,
	0x8c, 0x82, 0x12, 0x43, 0xdc, 0x88, 0x98, 0xb8, 0xd1, 0xc5, 0x98, 0xb8, 0xd4, 0x74, 0xc6, 0x4b,
	0xa7, 0x09, 0x9d, 0x0e, 0xbd, 0x1d, 0xa2, 0xfc, 0x81, 0x7f, 0xe0, 0x27, 0xb9, 0xf4, 0x0b, 0x88,
	0xc1, 0x3f, 0x60, 0xe1, 0xc2, 0x95, 0xb1, 0x08, 0x12, 0x12, 0x77, 0x3d, 0xa7, 0xe7, 0x9e, 0x7b,
	0x4e, 0x0b, 0x6d, 0x8d, 0x44, 0x42, 0x22, 0xc5, 0x95, 0x35, 0xce, 0xb0, 0xf6, 0x42, 0x39, 0x8d,
	0x68, 0xe3, 0x6a, 0x50, 0xc5, 0xcb, 0x87, 0x9d, 0xbe, 0x54, 0xae, 0xa8, 0xb3, 0x38, 0x37, 0x3a,
	0x91, 0x46, 0x9a, 0xc4, 0xcb, 0xb2, 0x7a, 0xe6, 0x91, 0x07, 0xfe, 0xb4, 0x1b, 0xe7, 0x4f, 0xa1,
	0xf5, 0xc2, 0x5a, 0x63, 0x53, 0xa4, 0xca, 0x94, 0x84, 0xec, 0x3e, 0x5c, 0xfc, 0xdd, 0x10, 0x06,
	0xf7, 0x82, 0x5e, 0x73, 0x72, 0x7b, 0xbb, 0xee, 0xb6, 0x88, 0x56, 0x7d, 0x2d, 0x3e, 0x8e, 0xf9,
	0xe0, 0xf1, 0x88, 0xa7, 0x7b, 0x05, 0xff, 0x1c, 0xc0, 0xe5, 0x2b, 0x74, 0x62, 0x2a, 0x9c, 0x60,
	0x77, 0x01, 0x08, 0x17, 0xef, 0xcb, 0x5a, 0x67, 0x68, 0xfd, 0xf0, 0x59, 0xda, 0x20, 0x5c, 0xbc,
	0xf6, 0x04, 0x7b, 0x07, 0x17, 0x54, 0x67, 0x25, 0x3a, 0x0a, 0xaf, 0x7a, 0xe3, 0xe9, 0x76, 0xdd,
	0x6d, 0xfe, 0x31, 0x26, 0xb5, 0xc2, 0x31, 0x7f, 0xc2, 0x7f, 0xad, 0xbb, 0x0f, 0x8e, 0xd2, 0x57,
	0xf6, 0x13, 0x69, 0xe1, 0x54, 0x3e, 0x17, 0x19, 0x25, 0xd2, 0xf4, 0x33, 0xe5, 0x66, 0x0a, 0xe7,
	0x1f, 0xe2, 0x89, 0x72, 0x4b, 0xcc, 0x9d, 0xb1, 0xa3, 0x47, 0xe9, 0xde, 0x94, 0xff, 0x0c, 0x00,
	0x9e, 0x17, 0x42, 0x95, 0x6f, 0x9c, 0x70, 0xc8, 0x86, 0x70, 0x43, 0x62, 0x89, 0xa4, 0xe8, 0xa5,
	0xa0, 0xe2, 0xb4, 0xcb, 0x6e, 0xe5, 0x70, 0xc0, 0xd3, 0x63, 0x15, 0xeb, 0xc1, 0x4d, 0xff, 0x2c,
	0xb9, 0x99, 0xbf, 0x45, 0x4b, 0xca, 0x94, 0x3e, 0x6b, 0x2b, 0x3d, 0xa5, 0xd9, 0x1d, 0x68, 0x38,
	0xa5, 0x91, 0x9c, 0xd0, 0x55, 0x78, 0x6d, 0xd7, 0xf5, 0x40, 0xb0, 0x0e, 0x5c, 0x12, 0xda, 0xa5,
	0xca, 0x91, 0xc2, 0x33, 0x7f, 0x79, 0xc0, 0x2c, 0x02, 0x90, 0x56, 0x54, 0x85, 0x8f, 0x19, 0x5e,
	0xf7, 0xf6, 0x47, 0x0c, 0x4b, 0xa0, 0x51, 0x13, 0xda, 0x67, 0x12, 0x4b, 0x17, 0x9e, 0xff, 0xef,
	0x0b, 0xfe, 0x69, 0x26, 0xb7, 0xbe, 0x6e, 0xa2, 0xe0, 0xdb, 0x26, 0x0a, 0xbe, 0x6f, 0xa2, 0xe0,
	0xcb, 0x8f, 0xe8, 0x4a, 0x76, 0xee, 0xd3, 0x0e, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xd7, 0x02,
	0x07, 0x8f, 0x2c, 0x02, 0x00, 0x00,
}

func (m *ErrorResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ErrorResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ErrorResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MetaData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetaData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetaData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Subnets) > 0 {
		i -= len(m.Subnets)
		copy(dAtA[i:], m.Subnets)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Subnets)))
		i--
		dAtA[i] = 0x12
	}
	if m.SeqNumber != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.SeqNumber))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ChainState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.UserAgent) > 0 {
		i -= len(m.UserAgent)
		copy(dAtA[i:], m.UserAgent)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.UserAgent)))
		i--
		dAtA[i] = 0x32
	}
	if m.GraphState != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.GraphState))
		i--
		dAtA[i] = 0x28
	}
	if m.Services != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.Services))
		i--
		dAtA[i] = 0x20
	}
	if m.Timestamp != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x18
	}
	if m.ProtocolVersion != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.ProtocolVersion))
		i--
		dAtA[i] = 0x10
	}
	if len(m.GenesisHash) > 0 {
		i -= len(m.GenesisHash)
		copy(dAtA[i:], m.GenesisHash)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.GenesisHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessages(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessages(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ErrorResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *MetaData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SeqNumber != 0 {
		n += 1 + sovMessages(uint64(m.SeqNumber))
	}
	l = len(m.Subnets)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ChainState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GenesisHash)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.ProtocolVersion != 0 {
		n += 1 + sovMessages(uint64(m.ProtocolVersion))
	}
	if m.Timestamp != 0 {
		n += 1 + sovMessages(uint64(m.Timestamp))
	}
	if m.Services != 0 {
		n += 1 + sovMessages(uint64(m.Services))
	}
	if m.GraphState != 0 {
		n += 1 + sovMessages(uint64(m.GraphState))
	}
	l = len(m.UserAgent)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessages(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessages(x uint64) (n int) {
	return sovMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ErrorResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ErrorResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ErrorResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = append(m.Message[:0], dAtA[iNdEx:postIndex]...)
			if m.Message == nil {
				m.Message = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MetaData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNumber", wireType)
			}
			m.SeqNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SeqNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subnets", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subnets = append(m.Subnets[:0], dAtA[iNdEx:postIndex]...)
			if m.Subnets == nil {
				m.Subnets = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ChainState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ChainState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenesisHash = append(m.GenesisHash[:0], dAtA[iNdEx:postIndex]...)
			if m.GenesisHash == nil {
				m.GenesisHash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolVersion", wireType)
			}
			m.ProtocolVersion = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProtocolVersion |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Services", wireType)
			}
			m.Services = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Services |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GraphState", wireType)
			}
			m.GraphState = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GraphState |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserAgent", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserAgent = append(m.UserAgent[:0], dAtA[iNdEx:postIndex]...)
			if m.UserAgent == nil {
				m.UserAgent = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessages
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessages
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessages
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessages        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessages          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessages = fmt.Errorf("proto: unexpected end of group")
)