package scale_codec

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/exp/constraints"
)

const (
	SingleByteModeMask uint8 = 0b00000000
	TwoByteModeMask    uint8 = 0b01000000
	FourByteModeMask   uint8 = 0b10000000
	BigIntegerModeMask uint8 = 0b11000000
)

type CompactMode uint8

const (
	SingleByteMode CompactMode = iota
	TwoByteMode
	FourByteMode
	BigIntegerMode
)

type CompactValue interface {
	isCompactValue()
}

var _ CompactValue = (*BigIntCompact)(nil)

type CompactInteger[T constraints.Integer] struct {
	Value *Integer[T]
}

func (CompactInteger[T]) isCompactValue() {}

type BigIntCompact struct {
	Value *big.Int
}

func (BigIntCompact) isCompactValue() {}

type Compact struct {
	value CompactValue
}

func (c Compact) Value() CompactValue {
	return c.value
}

func (c Compact) MarshalSCALE() ([]byte, error) {
	panic("not implemented yet")
}

func checkCompactMode(b uint8) CompactMode {
	switch b << 6 {
	case SingleByteModeMask:
		return SingleByteMode
	case TwoByteModeMask:
		return TwoByteMode
	case FourByteModeMask:
		return FourByteMode
	case BigIntegerModeMask:
		return BigIntegerMode
	default:
		panic(fmt.Sprintf("cannot define mode: %d", b))
	}
}

func (c *Compact) UnmarshalSCALE(reader io.Reader) error {
	fstByte := make([]byte, 1)
	_, err := reader.Read(fstByte)
	if err != nil {
		return fmt.Errorf("reading first compact byte: %w", err)
	}

	mode := checkCompactMode(fstByte[0])
	switch mode {
	case SingleByteMode:
		integer := &Integer[uint8]{}
		err := integer.UnmarshalSCALE(bytes.NewReader([]byte{fstByte[0] >> 2}))
		if err != nil {
			return err
		}
		c.value = &CompactInteger[uint8]{integer}
	case TwoByteMode:
		integer := &Integer[uint16]{}
		nextByte := make([]byte, 1)
		_, err := reader.Read(nextByte)
		if err != nil {
			return err
		}

		numberToDecode := bytes.Join(
			[][]byte{
				{(fstByte[0] >> 2) | (nextByte[0] << 6)},
				{nextByte[0] >> 2},
			},
			nil)
		err = integer.UnmarshalSCALE(bytes.NewReader(numberToDecode))
		if err != nil {
			return err
		}

		c.value = &CompactInteger[uint16]{integer}
	case FourByteMode:
		integer := &Integer[uint32]{}
		nextBytes := make([]byte, 3)
		_, err := reader.Read(nextBytes)
		if err != nil {
			return err
		}

		numberToDecode := bytes.Join(
			[][]byte{
				{(fstByte[0] >> 2) | (nextBytes[0] << 6)},
				{(nextBytes[0] >> 2) | (nextBytes[1] << 6)},
				{(nextBytes[1] >> 2) | (nextBytes[2] << 6)},
				{nextBytes[2] >> 2},
			},
			nil)
		err = integer.UnmarshalSCALE(bytes.NewReader(numberToDecode))
		if err != nil {
			return err
		}

		c.value = &CompactInteger[uint32]{integer}
	case BigIntegerMode:
		amountOfNextBytes := (fstByte[0] >> 2) + 4
		nextBytes := make([]byte, amountOfNextBytes)
		_, err := reader.Read(nextBytes)
		if err != nil {
			return err
		}

		switch {
		case amountOfNextBytes < 8:
			integer := &Integer[uint32]{}
			err = integer.UnmarshalSCALE(bytes.NewReader(nextBytes))
			if err != nil {
				return err
			}
			c.value = &CompactInteger[uint32]{integer}
			return nil
		case amountOfNextBytes == 8:
			integer := &Integer[uint64]{}
			err = integer.UnmarshalSCALE(bytes.NewReader(nextBytes))
			if err != nil {
				return err
			}
			c.value = &CompactInteger[uint64]{integer}
			return nil
		default:
			u128Value := new(U128)
			err = u128Value.UnmarshalSCALE(bytes.NewReader(nextBytes))
			if err != nil {
				return err
			}

			c.value = &BigIntCompact{u128Value.ToBigInt()}
			return nil
		}
	default:
		panic(fmt.Sprintf("mode not supported: %v", mode))
	}

	return nil
}
