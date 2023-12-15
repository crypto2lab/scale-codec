package scale_codec

import (
	"errors"
	"fmt"
	"io"
)

var ErrExpectedOneByteRead = errors.New("expected one byte read")

type Bool struct {
	Value bool
}

func (Bool) New() Encodable {
	return &Bool{}
}

func (b Bool) MarshalSCALE() ([]byte, error) {
	var value byte = 0x00
	if b.Value {
		value = 0x01
	}

	return []byte{value}, nil
}

func (b *Bool) UnmarshalSCALE(byteReader io.Reader) error {
	bValue := make([]byte, 1)
	n, err := byteReader.Read(bValue)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: %v", ErrExpectedOneByteRead, n)
	}

	switch bValue[0] {
	case 0x01:
		b.Value = true
	case 0x00:
		b.Value = false
	default:
		return fmt.Errorf("unknown byte to decode bool: %v", bValue)
	}

	return nil
}

type OptionBool struct {
	*Bool
}

func (o OptionBool) MarshalSCALE() ([]byte, error) {
	if o.Bool == nil {
		return []byte{0x00}, nil
	}

	if o.Bool.Value {
		return []byte{0x01}, nil
	}

	return []byte{0x02}, nil
}

func (o *OptionBool) UnmarshalSCALE(r io.Reader) error {
	bValue := make([]byte, 1)
	n, err := r.Read(bValue)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: %v", ErrExpectedOneByteRead, n)
	}

	switch bValue[0] {
	case 0x00:
		o.Bool = nil
	case 0x01:
		o.Bool = &Bool{true}
	case 0x02:
		o.Bool = &Bool{false}
	default:
		return fmt.Errorf("unknown byte to decode bool: %v", bValue)
	}

	return nil
}
