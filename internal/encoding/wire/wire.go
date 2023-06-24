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
	VarientType  Type = 0
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
