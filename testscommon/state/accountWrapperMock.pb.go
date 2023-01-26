// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: accountWrapperMock.proto

package state

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type AccountWrapMockData struct {
	MockValue int64 `protobuf:"varint,1,opt,name=MockValue,proto3" json:"MockValue"`
}

func (m *AccountWrapMockData) Reset()      { *m = AccountWrapMockData{} }
func (*AccountWrapMockData) ProtoMessage() {}
func (*AccountWrapMockData) Descriptor() ([]byte, []int) {
	return fileDescriptor_8fc2a076cf238279, []int{0}
}
func (m *AccountWrapMockData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountWrapMockData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *AccountWrapMockData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountWrapMockData.Merge(m, src)
}
func (m *AccountWrapMockData) XXX_Size() int {
	return m.Size()
}
func (m *AccountWrapMockData) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountWrapMockData.DiscardUnknown(m)
}

var xxx_messageInfo_AccountWrapMockData proto.InternalMessageInfo

func (m *AccountWrapMockData) GetMockValue() int64 {
	if m != nil {
		return m.MockValue
	}
	return 0
}

func init() {
	proto.RegisterType((*AccountWrapMockData)(nil), "proto.AccountWrapMockData")
}

func init() { proto.RegisterFile("accountWrapperMock.proto", fileDescriptor_8fc2a076cf238279) }

var fileDescriptor_8fc2a076cf238279 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x48, 0x4c, 0x4e, 0xce,
	0x2f, 0xcd, 0x2b, 0x09, 0x2f, 0x4a, 0x2c, 0x28, 0x48, 0x2d, 0xf2, 0xcd, 0x4f, 0xce, 0xd6, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xba, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49,
	0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xe9, 0xf9, 0xe9, 0xf9, 0xfa, 0x60, 0xe1, 0xa4, 0xd2, 0x34, 0x30,
	0x0f, 0xcc, 0x01, 0xb3, 0x20, 0xba, 0x94, 0x9c, 0xb8, 0x84, 0x1d, 0x11, 0x26, 0x82, 0x8c, 0x73,
	0x49, 0x2c, 0x49, 0x14, 0xd2, 0xe6, 0xe2, 0x04, 0xb1, 0xc3, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x98, 0x9d, 0x78, 0x5f, 0xdd, 0x93, 0x47, 0x08, 0x06, 0x21, 0x98, 0x4e, 0xf6,
	0x17, 0x1e, 0xca, 0x31, 0xdc, 0x78, 0x28, 0xc7, 0xf0, 0xe1, 0xa1, 0x1c, 0x63, 0xc3, 0x23, 0x39,
	0xc6, 0x15, 0x8f, 0xe4, 0x18, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc6, 0x23,
	0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x5f, 0x3c, 0x92, 0x63, 0xf8, 0xf0, 0x48, 0x8e, 0x71, 0xc2,
	0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0x62, 0x2d, 0x2e, 0x49,
	0x2c, 0x49, 0x4d, 0x62, 0x03, 0xbb, 0xc5, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe4, 0xb5, 0x13,
	0x00, 0xdd, 0x00, 0x00, 0x00,
}

func (this *AccountWrapMockData) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountWrapMockData)
	if !ok {
		that2, ok := that.(AccountWrapMockData)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MockValue != that1.MockValue {
		return false
	}
	return true
}
func (this *AccountWrapMockData) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&state.AccountWrapMockData{")
	s = append(s, "MockValue: "+fmt.Sprintf("%#v", this.MockValue)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringAccountWrapperMock(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *AccountWrapMockData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountWrapMockData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountWrapMockData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MockValue != 0 {
		i = encodeVarintAccountWrapperMock(dAtA, i, uint64(m.MockValue))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintAccountWrapperMock(dAtA []byte, offset int, v uint64) int {
	offset -= sovAccountWrapperMock(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AccountWrapMockData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MockValue != 0 {
		n += 1 + sovAccountWrapperMock(uint64(m.MockValue))
	}
	return n
}

func sovAccountWrapperMock(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAccountWrapperMock(x uint64) (n int) {
	return sovAccountWrapperMock(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *AccountWrapMockData) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountWrapMockData{`,
		`MockValue:` + fmt.Sprintf("%v", this.MockValue) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAccountWrapperMock(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *AccountWrapMockData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountWrapperMock
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
			return fmt.Errorf("proto: AccountWrapMockData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountWrapMockData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MockValue", wireType)
			}
			m.MockValue = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountWrapperMock
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MockValue |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAccountWrapperMock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAccountWrapperMock
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAccountWrapperMock
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAccountWrapperMock(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAccountWrapperMock
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
					return 0, ErrIntOverflowAccountWrapperMock
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
					return 0, ErrIntOverflowAccountWrapperMock
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
				return 0, ErrInvalidLengthAccountWrapperMock
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAccountWrapperMock
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAccountWrapperMock
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAccountWrapperMock        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAccountWrapperMock          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAccountWrapperMock = fmt.Errorf("proto: unexpected end of group")
)
