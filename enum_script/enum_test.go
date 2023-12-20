package main

import (
	"bytes"
	"fmt"
	"io"

	scale_codec "github.com/crypto2lab/scale-codec"
)

type Nested interface {
	scale_codec.Encodable
	IsNested()
}

func UnmarhalNested(reader io.Reader) (Nested, error) {
	enumTag := make([]byte, 1)
	n, err := reader.Read(enumTag)
	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, fmt.Errorf("%w: got %v", scale_codec.ErrExpectedOneByteRead, n)
	}

	switch enumTag[0] {

	case NumberIndex:
		unmarshaler := NewNumber()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	default:
		return nil, fmt.Errorf("unexpected enum tag: %v", enumTag[0])
	}
}

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

	case SingleIndex:
		unmarshaler := NewSingle()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

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

	case AIndex:
		unmarshaler := NewA()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case BIndex:
		unmarshaler := NewB()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case CIndex:
		unmarshaler := NewC()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case DIndex:
		unmarshaler := NewD()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case EIndex:
		unmarshaler := NewE()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case FIndex:
		unmarshaler := NewF()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case GIndex:
		unmarshaler := NewG()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case HIndex:
		unmarshaler := NewH()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case JIndex:
		unmarshaler := NewJ()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case KIndex:
		unmarshaler := NewK()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	case LIndex:
		unmarshaler := NewL()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err

	default:
		return nil, fmt.Errorf("unexpected enum tag: %v", enumTag[0])
	}
}

var NumberIndex byte = 0

var _ Nested = (*Number)(nil)

type Number struct {
	Inner *scale_codec.Integer[int32]
}

func NewNumber() *Number {
	return &Number{
		Inner: new(scale_codec.Integer[int32]),
	}
}

func (Number) IsNested() {}

func (i Number) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := NumberIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *Number) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var SingleIndex byte = 0

var _ MyScaleEncodedEnum = (*Single)(nil)

type Single struct {
	Inner *scale_codec.SimpleVariant
}

func NewSingle() *Single {
	return &Single{
		Inner: new(scale_codec.SimpleVariant),
	}
}

func (Single) IsMyScaleEncodedEnum() {}

func (i Single) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := SingleIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *Single) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var IntIndex byte = 1

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

var BoolIndex byte = 2

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

var AIndex byte = 3

var _ MyScaleEncodedEnum = (*A)(nil)

type A struct {
	Inner *scale_codec.Option
}

func NewA() *A {
	return &A{
		Inner: scale_codec.NewOption(new(scale_codec.Bool)),
	}
}

func (A) IsMyScaleEncodedEnum() {}

func (i A) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := AIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *A) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var BIndex byte = 4

var _ MyScaleEncodedEnum = (*B)(nil)

type B struct {
	Inner *scale_codec.Result
}

func NewB() *B {
	return &B{
		Inner: scale_codec.NewResult(new(scale_codec.Integer[uint64]), new(scale_codec.Integer[uint64])),
	}
}

func (B) IsMyScaleEncodedEnum() {}

func (i B) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := BIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *B) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var CIndex byte = 5

var _ MyScaleEncodedEnum = (*C)(nil)

type C struct {
	Inner *scale_codec.OptionG[Nested]
}

func NewC() *C {
	return &C{
		Inner: scale_codec.NewOptionG[Nested](nil),
	}
}

func (C) IsMyScaleEncodedEnum() {}

func (i C) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := CIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *C) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var DIndex byte = 6

var _ MyScaleEncodedEnum = (*D)(nil)

type D struct {
	Inner *scale_codec.Result
}

func NewD() *D {
	return &D{
		Inner: scale_codec.NewResult(Nested, new(scale_codec.Integer[uint64])),
	}
}

func (D) IsMyScaleEncodedEnum() {}

func (i D) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := DIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *D) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var EIndex byte = 7

var _ MyScaleEncodedEnum = (*E)(nil)

type E struct {
	Inner *scale_codec.Result
}

func NewE() *E {
	return &E{
		Inner: scale_codec.NewResult(Nested, Nested),
	}
}

func (E) IsMyScaleEncodedEnum() {}

func (i E) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := EIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *E) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var FIndex byte = 8

var _ MyScaleEncodedEnum = (*F)(nil)

type F struct {
	Inner *scale_codec.Result
}

func NewF() *F {
	return &F{
		Inner: scale_codec.NewResult(new(scale_codec.Integer[uint64]), Nested),
	}
}

func (F) IsMyScaleEncodedEnum() {}

func (i F) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := FIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *F) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var GIndex byte = 9

var _ MyScaleEncodedEnum = (*G)(nil)

type G struct {
	Inner *scale_codec.Tuple
}

func NewG() *G {
	return &G{
		Inner: scale_codec.NewTuple(new(scale_codec.Integer[uint64]), new(scale_codec.Bool)),
	}
}

func (G) IsMyScaleEncodedEnum() {}

func (i G) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := GIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *G) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var HIndex byte = 10

var _ MyScaleEncodedEnum = (*H)(nil)

type H struct {
	Inner *scale_codec.Option
}

func NewH() *H {
	return &H{
		Inner: scale_codec.NewOption(scale_codec.NewTuple(new(scale_codec.Integer[uint64]), new(scale_codec.Bool))),
	}
}

func (H) IsMyScaleEncodedEnum() {}

func (i H) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := HIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *H) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var JIndex byte = 11

var _ MyScaleEncodedEnum = (*J)(nil)

type J struct {
	Inner *scale_codec.Result
}

func NewJ() *J {
	return &J{
		Inner: scale_codec.NewResult(scale_codec.NewTuple(new(scale_codec.Integer[uint64]), new(scale_codec.Bool)), new(scale_codec.Bool)),
	}
}

func (J) IsMyScaleEncodedEnum() {}

func (i J) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := JIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *J) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var KIndex byte = 12

var _ MyScaleEncodedEnum = (*K)(nil)

type K struct {
	Inner *scale_codec.Tuple
}

func NewK() *K {
	return &K{
		Inner: scale_codec.NewTuple(scale_codec.NewOption(new(scale_codec.Bool)), scale_codec.NewResult(new(scale_codec.Bool), new(scale_codec.Bool))),
	}
}

func (K) IsMyScaleEncodedEnum() {}

func (i K) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := KIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *K) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}

var LIndex byte = 13

var _ MyScaleEncodedEnum = (*L)(nil)

type L struct {
	Inner *scale_codec.Result
}

func NewL() *L {
	return &L{
		Inner: scale_codec.NewResult(scale_codec.NewOption(scale_codec.NewTuple(new(scale_codec.Integer[uint64]), new(scale_codec.Bool))), new(scale_codec.Integer[uint64])),
	}
}

func (L) IsMyScaleEncodedEnum() {}

func (i L) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := LIndex
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *L) UnmarshalSCALE(reader io.Reader) error {
	return i.Inner.UnmarshalSCALE(reader)
}