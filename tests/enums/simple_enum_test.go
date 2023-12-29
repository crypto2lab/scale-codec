package main

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestSimpleEnumMarshaler(t *testing.T) {
	cases := []struct {
		marshaler     MyScaleEncodedEnum
		expectedBytes []byte
	}{
		{
			expectedBytes: []byte{0},
			marshaler:     NewSingle(),
		},
		{
			marshaler: &Int{
				Inner: &scale_codec.Integer[uint64]{Value: 32},
			},
			expectedBytes: []byte{1, 32, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			marshaler: &Bool{
				Inner: &scale_codec.Bool{Value: true},
			},
			expectedBytes: []byte{2, 1},
		},
		{
			marshaler: &A{
				Inner: scale_codec.Some(&scale_codec.Bool{Value: true}),
			},
			expectedBytes: []byte{3, 1, 1},
		},
		{
			marshaler: &A{
				Inner: scale_codec.None(),
			},
			expectedBytes: []byte{3, 0},
		},
		{
			marshaler: &B{
				Inner: scale_codec.Ok(&scale_codec.Integer[uint64]{Value: 108}),
			},
			expectedBytes: []byte{4, 0, 108, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			expectedBytes: []byte{4, 1, 90, 0, 0, 0, 0, 0, 0, 0},
			marshaler: &B{
				Inner: scale_codec.Err(&scale_codec.Integer[uint64]{Value: 90}),
			},
		},
		{
			expectedBytes: []byte{5, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			marshaler: &G{
				Inner: scale_codec.NewTuple(
					&scale_codec.Integer[uint64]{60},
					&scale_codec.Bool{false},
				),
			},
		},
		{
			expectedBytes: []byte{6, 1, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			marshaler: &H{
				Inner: scale_codec.Some(
					scale_codec.NewTuple(
						&scale_codec.Integer[uint64]{60},
						&scale_codec.Bool{false},
					),
				),
			},
		},
		{
			expectedBytes: []byte{7, 0, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			marshaler: &J{
				Inner: scale_codec.Ok(
					scale_codec.NewTuple(
						&scale_codec.Integer[uint64]{60},
						&scale_codec.Bool{false},
					),
				),
			},
		},
		{
			expectedBytes: []byte{8, 1, 1, 0, 0},
			marshaler: &K{
				Inner: scale_codec.NewTuple(
					scale_codec.Some(&scale_codec.Bool{Value: true}),
					scale_codec.Ok(&scale_codec.Bool{Value: false}),
				),
			},
		},
		{
			expectedBytes: []byte{10, 1, 0, 10, 0, 0, 0},
			marshaler: &M{
				Inner: scale_codec.SomeG[Nested](
					&Number{Inner: &scale_codec.Integer[uint32]{10}},
				),
			},
		},
		{
			expectedBytes: []byte{10, 0},
			marshaler: &M{
				Inner: scale_codec.NoneG[Nested](),
			},
		},
		{
			expectedBytes: []byte{11, 0, 0, 78, 0, 0, 0},
			marshaler: &N{
				Inner: scale_codec.OkG[Nested, *scale_codec.Bool](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: 78}}),
			},
		},
		{
			expectedBytes: []byte{11, 1, 1},
			marshaler: &N{
				Inner: scale_codec.ErrG[Nested, *scale_codec.Bool](
					&scale_codec.Bool{Value: true}),
			},
		},
		{
			expectedBytes: []byte{12, 0, 1},
			marshaler: &O{
				Inner: scale_codec.OkG[*scale_codec.Bool, Nested](
					&scale_codec.Bool{Value: true}),
			},
		},
		{
			expectedBytes: []byte{12, 1, 0, 76, 0, 0, 0},
			marshaler: &O{
				Inner: scale_codec.ErrG[*scale_codec.Bool, Nested](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: 76}}),
			},
		},
		{
			expectedBytes: []byte{13, 0, 0, 255, 255, 255, 255},
			marshaler: &P{
				Inner: scale_codec.OkG[Nested, Error](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: math.MaxUint32}}),
			},
		},
		{
			expectedBytes: []byte{13, 1, 0},
			marshaler: &P{
				Inner: scale_codec.ErrG[Nested, Error](NewFailureX()),
			},
		},
		{
			expectedBytes: []byte{14, 0, 77, 0, 0, 0, 89, 0, 0, 0, 0, 0, 0, 0, 0},
			marshaler: &Q{
				Inner: &T3[Nested, *scale_codec.Integer[uint64], Error]{
					F0: &Number{Inner: &scale_codec.Integer[uint32]{77}},
					F1: &scale_codec.Integer[uint64]{89},
					F2: NewFailureX(),
				},
			},
		},
	}

	for _, tt := range cases {
		output, err := tt.marshaler.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, output) {
			t.Fatalf("\nexpected: %v\nactual: %v",
				tt.expectedBytes, output)
		}
	}
}

func TestSimpleEnumUnmarshaler(t *testing.T) {
	cases := []struct {
		inputBytes      []byte
		expectedVariant MyScaleEncodedEnum
	}{
		{
			inputBytes:      []byte{0},
			expectedVariant: NewSingle(),
		},
		{
			inputBytes: []byte{1, 32, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &Int{
				Inner: &scale_codec.Integer[uint64]{Value: 32},
			},
		},
		{
			inputBytes: []byte{2, 1},
			expectedVariant: &Bool{
				Inner: &scale_codec.Bool{Value: true},
			},
		},
		{
			inputBytes: []byte{3, 1, 1},
			expectedVariant: &A{
				Inner: scale_codec.Some(&scale_codec.Bool{Value: true}),
			},
		},
		{
			inputBytes: []byte{3, 0},
			expectedVariant: &A{
				Inner: scale_codec.None(),
			},
		},
		{
			inputBytes: []byte{4, 0, 108, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &B{
				Inner: scale_codec.Ok(&scale_codec.Integer[uint64]{Value: 108}),
			},
		},
		{
			inputBytes: []byte{4, 1, 90, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &B{
				Inner: scale_codec.Err(&scale_codec.Integer[uint64]{Value: 90}),
			},
		},
		{
			inputBytes: []byte{5, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &G{
				Inner: scale_codec.NewTuple(
					&scale_codec.Integer[uint64]{60},
					&scale_codec.Bool{false},
				),
			},
		},
		{
			inputBytes: []byte{6, 1, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &H{
				Inner: scale_codec.Some(
					scale_codec.NewTuple(
						&scale_codec.Integer[uint64]{60},
						&scale_codec.Bool{false},
					),
				),
			},
		},
		{
			inputBytes: []byte{7, 0, 60, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &J{
				Inner: scale_codec.Ok(
					scale_codec.NewTuple(
						&scale_codec.Integer[uint64]{60},
						&scale_codec.Bool{false},
					),
				),
			},
		},
		{
			inputBytes: []byte{8, 1, 1, 0, 0},
			expectedVariant: &K{
				Inner: scale_codec.NewTuple(
					scale_codec.Some(&scale_codec.Bool{Value: true}),
					scale_codec.Ok(&scale_codec.Bool{Value: false}),
				),
			},
		},
		{
			inputBytes: []byte{10, 1, 0, 10, 0, 0, 0},
			expectedVariant: &M{
				Inner: scale_codec.SomeG[Nested](
					&Number{
						Inner: &scale_codec.Integer[uint32]{10},
					},
				),
			},
		},
		{
			inputBytes: []byte{10, 0},
			expectedVariant: &M{
				Inner: scale_codec.NoneG[Nested](),
			},
		},
		{
			inputBytes: []byte{11, 0, 0, 78, 0, 0, 0},
			expectedVariant: &N{
				Inner: scale_codec.OkG[Nested, *scale_codec.Bool](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: 78}}),
			},
		},
		{
			inputBytes: []byte{11, 1, 1},
			expectedVariant: &N{
				Inner: scale_codec.ErrG[Nested, *scale_codec.Bool](
					&scale_codec.Bool{Value: true}),
			},
		},
		{
			inputBytes: []byte{12, 0, 1},
			expectedVariant: &O{
				Inner: scale_codec.OkG[*scale_codec.Bool, Nested](
					&scale_codec.Bool{Value: true}),
			},
		},
		{
			inputBytes: []byte{12, 1, 0, 76, 0, 0, 0},
			expectedVariant: &O{
				Inner: scale_codec.ErrG[*scale_codec.Bool, Nested](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: 76}}),
			},
		},
		{
			inputBytes: []byte{13, 0, 0, 255, 255, 255, 255},
			expectedVariant: &P{
				Inner: scale_codec.OkG[Nested, Error](
					&Number{Inner: &scale_codec.Integer[uint32]{Value: math.MaxUint32}}),
			},
		},
		{
			inputBytes: []byte{13, 1, 0},
			expectedVariant: &P{
				Inner: scale_codec.ErrG[Nested, Error](NewFailureX()),
			},
		},
		{
			inputBytes: []byte{14, 0, 77, 0, 0, 0, 89, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &Q{
				Inner: &T3[Nested, *scale_codec.Integer[uint64], Error]{
					F0: &Number{Inner: &scale_codec.Integer[uint32]{77}},
					F1: &scale_codec.Integer[uint64]{89},
					F2: NewFailureX(),
				},
			},
		},
	}

	for _, tt := range cases {
		variant, err := UnmarshalMyScaleEncodedEnum(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tt.expectedVariant, variant) {
			t.Fatalf("\nexpected: %q (%T)\ngot: %q (%T)", tt.expectedVariant, tt.expectedVariant,
				variant, variant)
		}
	}
}
