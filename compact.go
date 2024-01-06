package scale_codec

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"math/bits"

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

var _ CompactValue = (*CompactBigInt)(nil)

type CompactInteger[T constraints.Integer] struct {
	Value T
}

func (CompactInteger[T]) isCompactValue() {}

type CompactBigInt struct {
	Value *big.Int
}

func (b *CompactBigInt) LessOrEqual(v int64) bool {
	cmp := b.Value.Cmp(big.NewInt(v))
	return cmp == -1 || cmp == 0
}

func (b *CompactBigInt) Less(v int64) bool {
	cmp := b.Value.Cmp(big.NewInt(v))
	return cmp == -1
}

func (b *CompactBigInt) GreaterOrEqual(v int64) bool {
	cmp := b.Value.Cmp(big.NewInt(v))
	return cmp == 0 || cmp == 1
}

func (b *CompactBigInt) Greater(v int64) bool {
	cmp := b.Value.Cmp(big.NewInt(v))
	return cmp == 1
}

func (*CompactBigInt) isCompactValue() {}

type Compact struct {
	Value CompactValue
}

func (c Compact) MarshalSCALE() ([]byte, error) {
	switch compactValue := c.Value.(type) {
	case *CompactInteger[uint8]:
		switch {
		case compactValue.Value <= 0b0011_1111:
			toMarshal := Integer[uint8]{Value: uint8(compactValue.Value) << 2}
			return toMarshal.MarshalSCALE()
		default:
			toMarshal := Integer[uint16]{Value: (uint16(compactValue.Value) << 2) | 0b00000001}
			return toMarshal.MarshalSCALE()
		}
	case *CompactInteger[uint16]:
		switch {
		case compactValue.Value <= 0b0011_1111:
			toMarshal := Integer[uint8]{Value: uint8(compactValue.Value) << 2}
			return toMarshal.MarshalSCALE()
		case 0b0011_1111 < compactValue.Value && compactValue.Value <= 0b0011_1111_1111_1111:
			toMarshal := Integer[uint16]{Value: (compactValue.Value << 2) | 0b00000001}
			return toMarshal.MarshalSCALE()
		default:
			toMarshal := Integer[uint32]{Value: (uint32(compactValue.Value) << 2) | 0b00000010}
			return toMarshal.MarshalSCALE()
		}
	case *CompactInteger[uint32]:
		switch {
		case compactValue.Value <= 0b0011_1111:
			toMarshal := Integer[uint8]{Value: uint8(compactValue.Value) << 2}
			return toMarshal.MarshalSCALE()
		case 0b0011_1111 < compactValue.Value && compactValue.Value <= 0b0011_1111_1111_1111:
			toMarshal := Integer[uint16]{Value: (uint16(compactValue.Value) << 2) | 0b00000001}
			return toMarshal.MarshalSCALE()
		case 0b0011_1111_1111_1111 < compactValue.Value && compactValue.Value <= 0b0011_1111_1111_1111_1111_1111_1111_1111:
			toMarshal := Integer[uint32]{Value: (uint32(compactValue.Value) << 2) | 0b00000010}
			return toMarshal.MarshalSCALE()
		default:
			toMarshal := Integer[uint32]{Value: compactValue.Value}
			encoded, err := toMarshal.MarshalSCALE()
			if err != nil {
				return nil, err
			}
			return bytes.Join([][]byte{{0b11}, encoded}, nil), nil
		}
	case *CompactInteger[uint64]:
		switch {
		case compactValue.Value <= 0b0011_1111:
			toMarshal := Integer[uint8]{Value: uint8(compactValue.Value) << 2}
			return toMarshal.MarshalSCALE()
		case 0b0011_1111 < compactValue.Value && compactValue.Value <= 0b0011_1111_1111_1111:
			toMarshal := Integer[uint16]{Value: (uint16(compactValue.Value) << 2) | 0b00000001}
			return toMarshal.MarshalSCALE()
		case 0b0011_1111_1111_1111 < compactValue.Value && compactValue.Value <= 0b0011_1111_1111_1111_1111_1111_1111_1111:
			toMarshal := Integer[uint32]{Value: (uint32(compactValue.Value) << 2) | 0b00000010}
			return toMarshal.MarshalSCALE()
		default:
			bytesNeeded := 8 - bits.LeadingZeros64(compactValue.Value)/8
			if bytesNeeded < 4 {
				panic("previous match arm matches anyting less than 2^30; qed")
			}

			output := make([]byte, bytesNeeded+1)
			output[0] = 0b11 + uint8((bytesNeeded-4)<<2)

			v := compactValue.Value
			for idx := 0; idx < bytesNeeded; idx++ {
				output[idx+1] = uint8(v)
				v >>= 8
			}

			if v != 0 {
				panic("shifted sufficient bits right to lead only leading zeros; qed")
			}

			return output, nil
		}
	case *CompactBigInt:
		switch {
		case compactValue.LessOrEqual(int64(0b0011_1111)):
			toMarshal := Integer[uint8]{Value: uint8(compactValue.Value.Int64()) << 2}
			return toMarshal.MarshalSCALE()
		case compactValue.Greater(int64(0b0011_1111)) && compactValue.LessOrEqual(int64(0b0011_1111_1111_1111)):
			toMarshal := Integer[uint16]{Value: (uint16(compactValue.Value.Int64()) << 2) | 0b00000001}
			return toMarshal.MarshalSCALE()
		case compactValue.Greater(int64(0b0011_1111_1111_1111)) && compactValue.LessOrEqual(int64(0b0011_1111_1111_1111_1111_1111_1111_1111)):
			toMarshal := Integer[uint32]{Value: (uint32(compactValue.Value.Int64()) << 2) | 0b00000010}
			return toMarshal.MarshalSCALE()
		default:
			u128Value := U128FromBigInt(compactValue.Value)
			bytesNeeded := 16 - u128Value.LeadingZeros()/8
			if bytesNeeded < 4 {
				panic("previous match arm matches anyting less than 2^30; qed")
			}

			output := make([]byte, bytesNeeded+1)
			output[0] = 0b11 + uint8((bytesNeeded-4)<<2)

			for idx := 0; idx < bytesNeeded; idx++ {
				output[idx+1] = u128Value.AsUint8()
				u128Value.Rsh(8)
			}

			if !u128Value.IsZero() {
				panic("shifted sufficient bits right to lead only leading zeros; qed")
			}

			return output, nil
		}
	default:
		panic("not implemented yet")
	}
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
		c.Value = &CompactInteger[uint8]{integer.Value}
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

		c.Value = &CompactInteger[uint16]{integer.Value}
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

		c.Value = &CompactInteger[uint32]{integer.Value}
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
			c.Value = &CompactInteger[uint32]{integer.Value}
			return nil
		case amountOfNextBytes == 8:
			integer := &Integer[uint64]{}
			err = integer.UnmarshalSCALE(bytes.NewReader(nextBytes))
			if err != nil {
				return err
			}
			c.Value = &CompactInteger[uint64]{integer.Value}
			return nil
		default:
			u128Value := new(U128)
			err = u128Value.UnmarshalSCALE(bytes.NewReader(nextBytes))
			if err != nil {
				return err
			}

			c.Value = &CompactBigInt{u128Value.ToBigInt()}
			return nil
		}
	default:
		panic(fmt.Sprintf("mode not supported: %v", mode))
	}

	return nil
}
