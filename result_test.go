package scale_codec_test

import (
	"bytes"
	"fmt"
	"io"
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

type TestResultNested interface {
	scale_codec.Encodable
	isTestResultEnum()
}

func UnmarshalTestResultNested(reader io.Reader) (TestResultNested, error) {
	enumTag := make([]byte, 1)
	n, err := reader.Read(enumTag)
	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, fmt.Errorf("%w: got %v", scale_codec.ErrExpectedOneByteRead, n)
	}

	switch enumTag[0] {
	case 0:
		numberx := &SimpleN{Inner: new(scale_codec.Integer[uint32])}
		err := numberx.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}

		return numberx, err
	default:
		return nil, fmt.Errorf("failure to decode bla bla bla")
	}
}

type SimpleN struct {
	Inner *scale_codec.Integer[uint32]
}

func (SimpleN) isTestResultEnum() {}

func (n SimpleN) MarshalSCALE() ([]byte, error) {
	innerEncode, err := n.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{[]byte{0x00}, innerEncode}, nil), nil
}

func (n *SimpleN) UnmarshalSCALE(reader io.Reader) error {
	return n.Inner.UnmarshalSCALE(reader)
}

func TestResultMarshalerGeneric(t *testing.T) {
	marshaler := scale_codec.OkG[TestResultNested, *scale_codec.Bool](
		&SimpleN{Inner: &scale_codec.Integer[uint32]{Value: 78}})

	output, err := marshaler.MarshalSCALE()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(output)
}

func TestResultUnmarshalerGeneric(t *testing.T) {
	inputBytes := []byte{0, 0, 78, 0, 0, 0}
	unmarshaler := new(scale_codec.ResultG[TestResultNested, *scale_codec.Bool])

	err := unmarshaler.UnmarshalSCALE(bytes.NewReader(inputBytes), UnmarshalTestResultNested, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := scale_codec.OkG[TestResultNested, *scale_codec.Bool](
		&SimpleN{Inner: &scale_codec.Integer[uint32]{Value: 78}})

	if !reflect.DeepEqual(expected, unmarshaler) {
		t.Fatalf("\nexpected: %v\nactual: %v", expected, unmarshaler)
	}

	inputBytes = []byte{1, 0}
	unmarshaler = new(scale_codec.ResultG[TestResultNested, *scale_codec.Bool])

	err = unmarshaler.UnmarshalSCALE(bytes.NewReader(inputBytes), UnmarshalTestResultNested, scale_codec.BoolFromRawBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected = scale_codec.ErrG[TestResultNested, *scale_codec.Bool](
		&scale_codec.Bool{Value: false},
	)

	if !reflect.DeepEqual(expected, unmarshaler) {
		t.Fatalf("\nexpected: %v\nactual: %v", expected, unmarshaler)
	}
}
