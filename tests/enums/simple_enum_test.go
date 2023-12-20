package main

import (
	"bytes"
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
	}

	for _, tt := range cases {
		variant, err := UnmarhalMyScaleEncodedEnum(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tt.expectedVariant, variant) {
			t.Fatalf("\nexpected: %q (%T)\ngot: %q (%T)", tt.expectedVariant, tt.expectedVariant,
				variant, variant)
		}
	}
}
