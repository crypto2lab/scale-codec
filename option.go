package scale_codec

import (
	"bytes"
	"fmt"
	"io"
)

var NoneEncoded = []byte{0x00}

type OptionG[T Marshaler] struct {
	inner  T
	isNone bool
}

func UnmarshalOptionFromRawBytes[T Marshaler](
	f func(io.Reader) (T, error)) func(reader io.Reader) (*OptionG[T], error) {
	return func(reader io.Reader) (*OptionG[T], error) {
		option := &OptionG[T]{}
		err := option.UnmarshalSCALE(reader, f)
		if err != nil {
			return nil, err
		}
		return option, nil
	}
}

func (o *OptionG[T]) MarshalSCALE() ([]byte, error) {
	if o.isNone {
		return NoneEncoded, nil
	}

	encodedOptionTag := []byte{0x01}
	encodedInner, err := o.inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{encodedOptionTag, encodedInner}, nil), nil
}

func (o *OptionG[T]) UnmarshalSCALE(reader io.Reader, f func(io.Reader) (T, error)) error {
	encodedOptionTag := make([]byte, 1)
	n, err := reader.Read(encodedOptionTag)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: want: 1, got: %v", ErrUnexpectedReadBytes, n)
	}

	switch encodedOptionTag[0] {
	case 0x00:
		o.isNone = true
		return nil
	case 0x01:
		innerValue, err := f(reader)
		if err != nil {
			return err
		}

		o.inner = innerValue
		o.isNone = false
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrUnexpectedOptionTag, encodedOptionTag[0])
	}
}

func SomeG[T Marshaler](inner T) *OptionG[T] {
	return &OptionG[T]{inner, false}
}

func NoneG[T Marshaler]() *OptionG[T] {
	return &OptionG[T]{
		isNone: true,
	}
}

type Option struct {
	inner  Encodable
	isNone bool
}

func NewOption(encodable Encodable) *Option {
	return &Option{inner: encodable}
}

func None() *Option {
	return &Option{
		isNone: true,
	}
}
func Some(v Encodable) *Option {
	return &Option{
		inner:  v,
		isNone: false,
	}
}

func (o Option) MarshalSCALE() ([]byte, error) {
	if o.isNone {
		return NoneEncoded, nil
	}

	encodedOptionTag := []byte{0x01}
	encodedInner, err := o.inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{encodedOptionTag, encodedInner}, nil), nil
}

func (o *Option) UnmarshalSCALE(reader io.Reader) error {
	encodedOptionTag := make([]byte, 1)
	n, err := reader.Read(encodedOptionTag)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: want: 1, got: %v", ErrUnexpectedReadBytes, n)
	}

	switch encodedOptionTag[0] {
	case 0x00:
		o.inner = nil
		o.isNone = true
		return nil
	case 0x01:
		o.isNone = false
		return o.inner.UnmarshalSCALE(reader)
	default:
		return fmt.Errorf("%w: %v", ErrUnexpectedOptionTag, encodedOptionTag[0])
	}
}

func (o *Option) IsNone() bool {
	return o.isNone
}
