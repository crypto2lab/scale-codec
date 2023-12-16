package scale_codec_test

import (
	"bytes"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestTupleMarshaler(t *testing.T) {
	cases := []struct {
		marshaler     *scale_codec.Tuple
		expectedBytes []byte
	}{
		{
			marshaler: scale_codec.NewTuple(
				scale_codec.Some(&scale_codec.Integer[uint64]{79}),
				scale_codec.Ok(&scale_codec.Bool{true}),
			),
			expectedBytes: []byte{1, 79, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			marshaler: scale_codec.NewTuple(
				scale_codec.None(),
				scale_codec.Err(&scale_codec.Integer[uint64]{44}),
			),
			expectedBytes: []byte{0, 1, 44, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range cases {
		actual, err := tt.marshaler.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, actual) {
			t.Fatalf("\nexpected: %v\nactual: %v", tt.expectedBytes, actual)
		}
	}
}

func TestTupleUnmarshaler(t *testing.T) {
	cases := []struct {
		unmarshaler *scale_codec.Tuple
		expected    *scale_codec.Tuple
		inputBytes  []byte
	}{
		{
			unmarshaler: scale_codec.NewTuple(
				scale_codec.NewOption(new(scale_codec.Integer[uint64])),
				scale_codec.NewResult(
					new(scale_codec.Bool),
					new(scale_codec.Integer[uint64])),
			),
			expected: scale_codec.NewTuple(
				scale_codec.None(),
				scale_codec.Err(&scale_codec.Integer[uint64]{44}),
			),
			inputBytes: []byte{0, 1, 44, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			unmarshaler: scale_codec.NewTuple(
				scale_codec.NewOption(new(scale_codec.Integer[uint64])),
				scale_codec.NewResult(
					new(scale_codec.Bool),
					new(scale_codec.Integer[uint64])),
			),
			expected: scale_codec.NewTuple(
				scale_codec.Some(&scale_codec.Integer[uint64]{79}),
				scale_codec.Ok(&scale_codec.Bool{true}),
			),
			inputBytes: []byte{1, 79, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
	}

	for _, tt := range cases {
		reader := bytes.NewReader(tt.inputBytes)
		err := tt.unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if reader.Len() != 0 {
			t.Fatalf("expected empty reader, missing %v bytes to read", reader.Len())
		}

		if !reflect.DeepEqual(tt.expected, tt.unmarshaler) {
			t.Fatalf("\nexpected: %v\nactual: %v", tt.expected, tt.unmarshaler)
		}
	}
}
