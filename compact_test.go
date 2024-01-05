package scale_codec_test

import (
	"bytes"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestCompactUnmarshal(t *testing.T) {
	cases := []struct {
		inputBytes []byte
		expected   scale_codec.CompactValue
	}{
		{
			inputBytes: []byte{16},
			expected: &scale_codec.CompactInteger[uint8]{
				Value: &scale_codec.Integer[uint8]{4},
			},
		},
		{
			inputBytes: []byte{1, 1},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: &scale_codec.Integer[uint16]{64},
			},
		},
		{
			inputBytes: []byte{17, 14},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: &scale_codec.Integer[uint16]{900},
			},
		},
		{
			inputBytes: []byte{1, 128},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: &scale_codec.Integer[uint16]{8192},
			},
		},
		{
			inputBytes: []byte{2, 0, 1, 0},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: &scale_codec.Integer[uint32]{16384},
			},
		},
		{
			inputBytes: []byte{66, 126, 5, 0},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: &scale_codec.Integer[uint32]{90000},
			},
		},
		{
			inputBytes: []byte{2, 0, 0, 128},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: &scale_codec.Integer[uint32]{536870912},
			},
		},
		{
			inputBytes: []byte{3, 0, 0, 0, 64},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: &scale_codec.Integer[uint32]{1073741824},
			},
		},
		{
			inputBytes: []byte{19, 255, 255, 255, 255, 255, 255, 255, 255},
			expected: &scale_codec.CompactInteger[uint64]{
				Value: &scale_codec.Integer[uint64]{^uint64(0)},
			},
		},
		{
			inputBytes: []byte{51, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expected: &scale_codec.BigIntCompact{
				Value: scale_codec.MaxU128.ToBigInt(),
			},
		},
	}

	for _, tt := range cases {
		compact := &scale_codec.Compact{}
		err := compact.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tt.expected, compact.Value()) {
			t.Fatalf("\nexpected: %+v\nactual: %+v\nfor input: %v\n", tt.expected, compact.Value(), tt.inputBytes)
		}
	}
}
