package main

import (
	"bytes"
	"fmt"
	"io"

	scale_codec "github.com/crypto2lab/scale-codec"
)

type MyScaleEncodedEnum interface {
	scale_codec.Encodable
	IsMyScaleEncodedEnum()
}

func UnmarhalMyScaleEncodedEnum(reader io.Reader) (MyScaleEncodedEnum, error) {
	enumTag := make([]byte, 1)
	n, err := reader.Read(enumTag)
	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, fmt.Errorf("%w: got %v", scale_codec.ErrExpectedOneByteRead, n)
	}

	switch enumTag[0] {
	
	case IntIndex:
		unmarshaler := NewInt()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err
	
	case BoolIndex:
		unmarshaler := NewBool()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err
	
	default:
		return nil, fmt.Errorf("unexpected enum tag: %v", enumTag[0])
	}
}


var IntIndex byte = 0

var _ MyScaleEncodedEnum = (*Int)(nil)

type Int struct {
	Inner *scale_codec.Integer[uint64]
}

func NewInt() *Int {
	return &Int{
		Inner: new(scale_codec.Integer[uint64]),
	}
}

func (Int) IsMyScaleEncodedEnum() {}

func (i Int) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := IntIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *Int) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}
var BoolIndex byte = 1

var _ MyScaleEncodedEnum = (*Bool)(nil)

type Bool struct {
	Inner *scale_codec.Bool
}

func NewBool() *Bool {
	return &Bool{
		Inner: new(scale_codec.Bool),
	}
}

func (Bool) IsMyScaleEncodedEnum() {}

func (i Bool) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := BoolIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *Bool) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}
