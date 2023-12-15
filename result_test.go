package scale_codec_test

import (
	"bytes"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestResultMarshaler(t *testing.T) {
	cases := []struct {
		result        *scale_codec.Result
		expectedBytes []byte
	}{
		{
			result:        scale_codec.Ok(&scale_codec.Integer[uint64]{Value: 332290}),
			expectedBytes: []byte{0, 2, 18, 5, 0, 0, 0, 0, 0},
		},
		{
			result:        scale_codec.Err(&scale_codec.Bool{Value: false}),
			expectedBytes: []byte{1, 0},
		},
		{
			result: scale_codec.Ok(
				scale_codec.Err(&scale_codec.Integer[uint8]{Value: 10})),
			expectedBytes: []byte{0, 1, 10},
		},
	}

	for _, tt := range cases {
		encoded, err := tt.result.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, encoded) {
			t.Fatalf("\nexpected: %v\nactual: %v", tt.expectedBytes, encoded)
		}
	}
}

func TestResultUnmarshal(t *testing.T) {
	cases := []struct {
		inputBytes  []byte
		result      *scale_codec.Result
		expected    scale_codec.Encodable
		expectedErr bool
	}{
		{
			inputBytes: []byte{0, 49, 0, 0, 0, 0, 0, 0, 0},
			result:     scale_codec.NewResult(new(scale_codec.Integer[uint64]), new(scale_codec.Bool)),
			expected:   &scale_codec.Integer[uint64]{Value: 49},
		},
		{
			inputBytes:  []byte{1, 0},
			result:      scale_codec.NewResult(new(scale_codec.Integer[uint64]), new(scale_codec.Bool)),
			expected:    &scale_codec.Bool{Value: false},
			expectedErr: true,
		},
		{
			inputBytes: []byte{0, 1, 10},
			result: scale_codec.NewResult(
				scale_codec.NewResult(
					new(scale_codec.Integer[uint64]),
					new(scale_codec.Integer[uint8])),
				new(scale_codec.Bool)),
			expected: scale_codec.Err(&scale_codec.Integer[uint8]{Value: 10}),
		},
	}

	for _, tt := range cases {
		err := tt.result.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var actualValue scale_codec.Encodable
		if tt.expectedErr {
			if !tt.result.IsErr() {
				t.Fatalf("exepected result error")
			}
			actualValue = tt.result.Err()
		} else {
			if tt.result.IsErr() {
				t.Fatalf("exepected result ok")
			}
			actualValue = tt.result.Unwrap()
		}

		if !reflect.DeepEqual(tt.expected, actualValue) {
			t.Fatalf("\nexpected: %v\nactual: %v", tt.expected, actualValue)
		}
	}
}
