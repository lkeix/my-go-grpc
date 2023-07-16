package wire

import (
	"errors"
	"io"
)

type Number int32

const (
	MinValidNumber      Number = 1
	FirstReservedNumber Number = 19000
	LastReservedNumber  Number = 19999
	MaxValidNumber      Number = 1<<29 - 1
)

func (n Number) IsValid() bool {
	return MinValidNumber <= n && n < FirstReservedNumber ||
		LastReservedNumber < n && n <= MaxValidNumber
}

type Type int8

const (
	VarintType   Type = 0
	Fixed32Type  Type = 1
	Fixed64Type  Type = 2
	BytesType    Type = 3
	EndGroupType Type = 4
)

const (
	_ = -iota
	errCodeTruncated
	errCodeFieldNumber
	errCodeOverflow
	errCodeReserved
	errCodeEndGroup
)

var (
	errFieldNumber = errors.New("invalid field number")
	errOverflow    = errors.New("variable length integer overflow")
	errReserved    = errors.New("cannnot parse reserved wire type")
	errEndGroup    = errors.New("mismatching end group marker")
	errParse       = errors.New("parse error")
)

func parseError(n int) error {
	if n >= 0 {
		return nil
	}

	switch n {
	case errCodeTruncated:
		return io.ErrUnexpectedEOF
	case errCodeFieldNumber:
		return errFieldNumber
	case errCodeOverflow:
		return errOverflow
	case errCodeReserved:
		return errReserved
	case errCodeEndGroup:
		return errEndGroup
	default:
		return errParse
	}
}

func ConsumeField(b []byte) (Number, Type, int) {
	num, typ, n := ConsumeTag(b)
	if n < 0 {
		return 0, 0, n
	}

	m := ConsumeFieldValue(num, typ, b[n:])
	if m < 0 {
		return 0, 0, m
	}

	return num, typ, n + m
}

func ConsumeFieldValue(num Number, typ Type, b []byte) int {
	var n int
	switch typ {
	case VarintType:
		_, n = ConsumeVarint(b)
	case Fixed32Type:
		_, n = ConsumeFixed32(b)
	case Fixed64Type:
		_, n = ConsumeFixed64(b)
	case BytesType:
		_, n = ConsumeBytes(b)
	case EndGroupType:
		return errCodeEndGroup
	default:
		return errCodeReserved
	}
	return n
}

func ConsumeTag(b []byte) (Number, Type, int) {
	v, n := ConsumeVarint(b)
	if n < 0 {
		return 0, 0, n
	}

	num, typ := DecodeTag(v)
	if num < MinValidNumber {
		return 0, 0, errCodeFieldNumber
	}
	return num, typ, n
}

func ConsumeFixed32(b []byte) (uint32, int) {
	if len(b) < 4 {
		return 0, errCodeTruncated
	}

	v := uint32(b[0])<<0 | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return v, 4
}

func ConsumeFixed64(b []byte) (uint64, int) {
	if len(b) < 8 {
		return 0, errCodeTruncated
	}

	v := uint64(b[0])<<0 | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	return v, 8
}

func ConsumeBytes(b []byte) ([]byte, int) {
	m, n := ConsumeVarint(b)
	if n < 0 {
		return nil, n
	}

	if m < uint64(len(b[n:])) {
		return nil, errCodeTruncated
	}

	return b[n : n+int(m)], n + int(m)
}

func ConsumeVarint(b []byte) (uint64, int) {
	var y uint64
	if len(b) <= 0 {
		return 0, errCodeTruncated
	}
	v := uint64(b[0])
	if v < 0x80 {
		return v, 1
	}
	v -= 0x80

	if len(b) <= 1 {
		return 0, errCodeTruncated
	}
	y = uint64(b[1])
	v += y << 7
	if y < 0x80 {
		return v, 2
	}
	v -= 0x80 << 7

	if len(b) <= 2 {
		return 0, errCodeTruncated
	}
	y = uint64(b[2])
	v += y << 14
	if y < 0x80 {
		return v, 3
	}
	v -= 0x80 << 14

	if len(b) <= 3 {
		return 0, errCodeTruncated
	}
	y = uint64(b[3])
	v += y << 21
	if y < 0x80 {
		return v, 4
	}
	v -= 0x80 << 21

	if len(b) <= 4 {
		return 0, errCodeTruncated
	}
	y = uint64(b[4])
	v += y << 28
	if y < 0x80 {
		return v, 5
	}
	v -= 0x80 << 28

	if len(b) <= 5 {
		return 0, errCodeTruncated
	}
	y = uint64(b[5])
	v += y << 35
	if y < 0x80 {
		return v, 6
	}
	v -= 0x80 << 35

	if len(b) <= 6 {
		return 0, errCodeTruncated
	}
	y = uint64(b[6])
	v += y << 42
	if y < 0x80 {
		return v, 7
	}
	v -= 0x80 << 42

	if len(b) <= 7 {
		return 0, errCodeTruncated
	}
	y = uint64(b[7])
	v += y << 49
	if y < 0x80 {
		return v, 8
	}
	v -= 0x80 << 49

	if len(b) <= 8 {
		return 0, errCodeTruncated
	}
	y = uint64(b[8])
	v += y << 56
	if y < 0x80 {
		return v, 9
	}
	v -= 0x80 << 56

	if len(b) <= 9 {
		return 0, errCodeTruncated
	}
	y = uint64(b[9])
	v += y << 63
	if y < 2 {
		return v, 10
	}
	return 0, errCodeOverflow
}

func DecodeTag(x uint64) (Number, Type) {
	num := Number(x >> 3)
	if num > MaxValidNumber {
		num = -1
	}
	return num, Type(x & 7)
}
