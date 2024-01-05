package scale_codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"unsafe"

	"golang.org/x/exp/constraints"
)

var ErrUnexpectedInteger = errors.New("unexpected integer")

func IntegerFromRawBytes[T constraints.Integer](reader io.Reader) (*Integer[T], error) {
	scaleInteger := new(Integer[T])
	if err := scaleInteger.UnmarshalSCALE(reader); err != nil {
		return nil, err
	}
	return scaleInteger, nil
}

type Integer[T constraints.Integer] struct {
	Value T
}

func (in Integer[T]) MarshalSCALE() ([]byte, error) {
	sizeof := unsafe.Sizeof(T(0))
	enc := make([]byte, sizeof)

	for i := 0; i < int(sizeof); i++ {
		enc[i] = byte(in.Value >> (8 * i))
	}

	return enc, nil
}

func (i *Integer[T]) UnmarshalSCALE(reader io.Reader) error {
	sizeof := unsafe.Sizeof(T(0))
	enc := make([]byte, sizeof)

	n, err := reader.Read(enc)
	if err != nil {
		return err
	}

	if n != int(sizeof) {
		return fmt.Errorf("%w: want: %v, got: %v", ErrUnexpectedReadBytes, sizeof, n)
	}

	acc := T(enc[0])
	for i := 1; i < int(sizeof); i++ {
		acc |= T(enc[i]) << (i * 8)
	}

	i.Value = acc
	return nil
}

var MaxU128 = U128{
	lower: ^uint64(0),
	upper: ^uint64(0),
}

type U128 struct {
	upper uint64
	lower uint64
}

func U128FromBigInt(b *big.Int) U128 {
	words := b.Bits()
	switch len(words) {
	case 0:
		return U128{}
	case 1:
		return U128{lower: uint64(words[0])}
	case 2:
		return U128{
			lower: uint64(words[0]),
			upper: uint64(words[1]),
		}
	default:
		return MaxU128
	}
}

func (u U128) MarshalSCALE() ([]byte, error) {
	encoded := make([]byte, 16)
	binary.LittleEndian.PutUint64(encoded[:8], u.lower)
	binary.LittleEndian.PutUint64(encoded[8:], u.upper)
	return encoded, nil
}

func (u *U128) UnmarshalSCALE(reader io.Reader) error {
	encoded := make([]byte, 16)
	_, err := reader.Read(encoded)
	if err != nil {
		return err
	}

	u.lower = binary.LittleEndian.Uint64(encoded[:8])
	u.upper = binary.LittleEndian.Uint64(encoded[8:])
	return nil
}

func (u *U128) ToBigInt() *big.Int {
	bigint := new(big.Int)
	bigint.SetBits([]big.Word{big.Word(u.lower), big.Word(u.upper)})
	return bigint
}
