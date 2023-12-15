package scale_codec_test

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestIntegersMarshaler(t *testing.T) {
	cases := []struct {
		marshaler     scale_codec.Marshaler
		expectedBytes []byte
	}{
		{
			marshaler: scale_codec.Integer[uint8]{
				Value: 255,
			},
			expectedBytes: []byte{255},
		},
		{
			marshaler: scale_codec.Integer[int8]{
				Value: -10,
			},
			expectedBytes: []byte{246},
		},
		{
			marshaler: scale_codec.Integer[uint16]{
				Value: ^uint16(0),
			},
			expectedBytes: []byte{255, 255},
		},
		{
			marshaler: scale_codec.Integer[int16]{
				Value: -32768,
			},
			expectedBytes: []byte{0, 128},
		},
		{
			marshaler: scale_codec.Integer[uint32]{
				Value: ^uint32(0),
			},
			expectedBytes: []byte{255, 255, 255, 255},
		},
		{
			marshaler: scale_codec.Integer[int32]{
				Value: -2147483648,
			},
			expectedBytes: []byte{0, 0, 0, 128},
		},
		{
			marshaler: scale_codec.Integer[uint64]{
				Value: ^uint64(0),
			},
			expectedBytes: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			marshaler: scale_codec.Integer[int64]{
				Value: -9223372036854775808,
			},
			expectedBytes: []byte{0, 0, 0, 0, 0, 0, 0, 128},
		},
	}

	for _, tt := range cases {
		actual, err := tt.marshaler.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, actual) {
			t.Fatalf("\nexpected: %v\ngot: %v\n", tt.expectedBytes, actual)
		}
	}
}

func TestIntegersUnmarshaler(t *testing.T) {
	cases := []struct {
		unmarshaler         scale_codec.Unmarshaler
		expectedUnmarshaler scale_codec.Unmarshaler
		inputBytes          []byte
	}{
		{
			unmarshaler: &scale_codec.Integer[uint8]{},
			expectedUnmarshaler: &scale_codec.Integer[uint8]{
				Value: 255,
			},
			inputBytes: []byte{255},
		},
		{
			unmarshaler: &scale_codec.Integer[int8]{},
			expectedUnmarshaler: &scale_codec.Integer[int8]{
				Value: -10,
			},
			inputBytes: []byte{246},
		},
		{
			unmarshaler: &scale_codec.Integer[uint16]{},
			expectedUnmarshaler: &scale_codec.Integer[uint16]{
				Value: ^uint16(0),
			},
			inputBytes: []byte{255, 255},
		},
		{
			unmarshaler: &scale_codec.Integer[int16]{},
			expectedUnmarshaler: &scale_codec.Integer[int16]{
				Value: -32768,
			},
			inputBytes: []byte{0, 128},
		},
		{
			unmarshaler: &scale_codec.Integer[uint32]{},
			expectedUnmarshaler: &scale_codec.Integer[uint32]{
				Value: ^uint32(0),
			},
			inputBytes: []byte{255, 255, 255, 255},
		},
		{
			unmarshaler: &scale_codec.Integer[int32]{},
			expectedUnmarshaler: &scale_codec.Integer[int32]{
				Value: -2147483648,
			},
			inputBytes: []byte{0, 0, 0, 128},
		},
		{
			unmarshaler: &scale_codec.Integer[uint64]{},
			expectedUnmarshaler: &scale_codec.Integer[uint64]{
				Value: ^uint64(0),
			},
			inputBytes: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			unmarshaler: &scale_codec.Integer[int64]{},
			expectedUnmarshaler: &scale_codec.Integer[int64]{
				Value: -9223372036854775808,
			},
			inputBytes: []byte{0, 0, 0, 0, 0, 0, 0, 128},
		},
	}

	for _, tt := range cases {
		err := tt.unmarshaler.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tt.expectedUnmarshaler, tt.unmarshaler) {
			t.Fatalf("\nexpected: %v\ngot: %v\n", tt.expectedUnmarshaler, tt.unmarshaler)
		}
	}
}

func TestU128Marshaler(t *testing.T) {
	cases := []struct {
		bignumber      string
		expectedOutput []byte
	}{
		{
			bignumber:      "340282366920938463463374607431768211455",
			expectedOutput: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			bignumber:      "340282366920938463463374607431768211454",
			expectedOutput: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			bignumber:      "18446744073709551615",
			expectedOutput: []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			bignumber:      "0",
			expectedOutput: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range cases {
		value, ok := new(big.Int).SetString(tt.bignumber, 10)
		if !ok {
			t.Fatalf("failed to convert %s to big int", tt.bignumber)
		}

		u128Value := scale_codec.U128FromBigInt(value)
		output, err := u128Value.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedOutput, output) {
			t.Fatalf("\nexpected: %v\ngot: %v\n", tt.expectedOutput, output)
		}
	}
}

func TestU128Unmarshaler(t *testing.T) {
	cases := []struct {
		expectedBigNumber string
		inputBytes        []byte
	}{
		{
			inputBytes:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedBigNumber: "0",
		},
		{
			inputBytes:        []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedBigNumber: "18446744073709551615",
		},
		{
			inputBytes:        []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expectedBigNumber: "340282366920938463463374607431768211455",
		},
		{
			inputBytes:        []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expectedBigNumber: "340282366920938463463374607431768211454",
		},
	}

	for _, tt := range cases {
		expected, ok := new(big.Int).SetString(tt.expectedBigNumber, 10)
		if !ok {
			t.Fatalf("failed to convert %s to big int", tt.expectedBigNumber)
		}

		actual := new(scale_codec.U128)
		err := actual.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reflect.DeepEqual(expected, actual.ToBigInt())
	}
}
