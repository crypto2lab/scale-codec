package scale_codec_test

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestCompactMarshal(t *testing.T) {
	cases := []struct {
		toMarshal     *scale_codec.Compact
		expectedBytes []byte
	}{
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint8]{Value: 4},
			},
			expectedBytes: []byte{16},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint8]{Value: 64},
			},
			expectedBytes: []byte{1, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint16]{Value: 4},
			},
			expectedBytes: []byte{16},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint16]{Value: 64},
			},
			expectedBytes: []byte{1, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint16]{Value: ^uint16(0)},
			},
			expectedBytes: []byte{254, 255, 3, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint32]{Value: 4},
			},
			expectedBytes: []byte{16},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint32]{Value: 64},
			},
			expectedBytes: []byte{1, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint32]{Value: uint32(^uint16(0))},
			},
			expectedBytes: []byte{254, 255, 3, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint32]{Value: uint32(^uint16(0)) + 10},
			},
			expectedBytes: []byte{38, 0, 4, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint32]{Value: ^uint32(0)},
			},
			expectedBytes: []byte{3, 255, 255, 255, 255},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: 4},
			},
			expectedBytes: []byte{16},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: 64},
			},
			expectedBytes: []byte{1, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: uint64(^uint16(0))},
			},
			expectedBytes: []byte{254, 255, 3, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: uint64(^uint16(0)) + 10},
			},
			expectedBytes: []byte{38, 0, 4, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: uint64(^uint32(0))},
			},
			expectedBytes: []byte{3, 255, 255, 255, 255},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: uint64(^uint32(0)) + 1},
			},
			expectedBytes: []byte{7, 0, 0, 0, 0, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactInteger[uint64]{Value: ^uint64(0)},
			},
			expectedBytes: []byte{19, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(4)},
			},
			expectedBytes: []byte{16},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(64)},
			},
			expectedBytes: []byte{1, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(int64(uint64(^uint16(0))))},
			},
			expectedBytes: []byte{254, 255, 3, 0},
		},

		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(int64(uint64(^uint16(0)) + 10))},
			},
			expectedBytes: []byte{38, 0, 4, 0},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(int64(uint64(^uint32(0))))},
			},
			expectedBytes: []byte{3, 255, 255, 255, 255},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: big.NewInt(int64(uint64(^uint32(0)) + 1))},
			},
			expectedBytes: []byte{7, 0, 0, 0, 0, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{Value: scale_codec.
					U128FromUpperLower(0, ^uint64(0)).
					ToBigInt(),
				},
			},
			expectedBytes: []byte{19, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{
					Value: new(big.Int).Add(
						scale_codec.
							U128FromUpperLower(0, ^uint64(0)).
							ToBigInt(),
						big.NewInt(1),
					),
				},
			},
			expectedBytes: []byte{23, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{
					Value: new(big.Int).Mul(
						scale_codec.
							U128FromUpperLower(0, ^uint64(0)).
							ToBigInt(),
						big.NewInt(int64(^uint32(0))),
					),
				},
			},
			expectedBytes: []byte{35, 1, 0, 0, 0, 255, 255, 255, 255, 254, 255, 255, 255},
		},

		{
			toMarshal: &scale_codec.Compact{
				Value: &scale_codec.CompactBigInt{
					Value: scale_codec.MaxU128.ToBigInt(),
				},
			},
			expectedBytes: []byte{51, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	for _, tt := range cases {
		result, err := tt.toMarshal.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, result) {
			t.Fatalf("\nexpected: %v\nactual: %v\n", tt.expectedBytes, result)
		}
	}
}

func TestCompactUnmarshal(t *testing.T) {
	cases := []struct {
		inputBytes []byte
		expected   scale_codec.CompactValue
	}{
		{
			inputBytes: []byte{16},
			expected: &scale_codec.CompactInteger[uint8]{
				Value: 4,
			},
		},
		{
			inputBytes: []byte{1, 1},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: 64,
			},
		},
		{
			inputBytes: []byte{17, 14},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: 900,
			},
		},
		{
			inputBytes: []byte{1, 128},
			expected: &scale_codec.CompactInteger[uint16]{
				Value: 8192,
			},
		},
		{
			inputBytes: []byte{2, 0, 1, 0},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: 16384,
			},
		},
		{
			inputBytes: []byte{66, 126, 5, 0},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: 90000,
			},
		},
		{
			inputBytes: []byte{2, 0, 0, 128},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: 536870912,
			},
		},
		{
			inputBytes: []byte{3, 0, 0, 0, 64},
			expected: &scale_codec.CompactInteger[uint32]{
				Value: 1073741824,
			},
		},
		{
			inputBytes: []byte{19, 255, 255, 255, 255, 255, 255, 255, 255},
			expected: &scale_codec.CompactInteger[uint64]{
				Value: ^uint64(0),
			},
		},
		{
			inputBytes: []byte{51, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expected: &scale_codec.CompactBigInt{
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

		if !reflect.DeepEqual(tt.expected, compact.Value) {
			t.Fatalf("\nexpected: %+v\nactual: %+v\nfor input: %v\n", tt.expected, compact.Value, tt.inputBytes)
		}
	}
}
