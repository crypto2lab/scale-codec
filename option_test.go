package scale_codec_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestOptionMarshal(t *testing.T) {
	cases := []struct {
		marshaler     *scale_codec.Option
		expectedBytes []byte
	}{
		{
			marshaler:     scale_codec.Some(scale_codec.Ok(&scale_codec.Bool{true})),
			expectedBytes: []byte{1, 0, 1},
		},
		{
			marshaler:     scale_codec.None(),
			expectedBytes: []byte{0x00},
		},
		{
			marshaler:     scale_codec.Some(scale_codec.None()),
			expectedBytes: []byte{1, 0},
		},
	}

	for _, tt := range cases {
		output, err := tt.marshaler.MarshalSCALE()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(tt.expectedBytes, output) {
			t.Fatalf("\nexpected: %v\nactual: %v", tt.expectedBytes, output)
		}
	}
}

func TestOptionUnmarshaler(t *testing.T) {
	cases := []struct {
		unmarshaler  *scale_codec.Option
		expected     scale_codec.Encodable
		expectedNone bool
		inputBytes   []byte
	}{
		{
			inputBytes: []byte{1, 0, 1},
			expected:   scale_codec.Some(scale_codec.Ok(&scale_codec.Bool{true})),
			unmarshaler: scale_codec.NewOption(
				scale_codec.NewResult(new(scale_codec.Bool), nil),
			),
		},
		{
			inputBytes:   []byte{0x00},
			expectedNone: true,
			unmarshaler:  scale_codec.NewOption(new(scale_codec.Bool)),
		},
	}

	for _, tt := range cases {
		err := tt.unmarshaler.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tt.expectedNone {
			if !tt.unmarshaler.IsNone() {
				t.Fatalf("expected option to be none")
			}
		} else {
			if !reflect.DeepEqual(tt.unmarshaler, tt.expected) {
				t.Fatalf("\nexpected: %v\nactual: %v", tt.expected, tt.unmarshaler)
			}
		}
	}
}

type TestEnum interface {
	scale_codec.Encodable
	isTestEnum()
}

func UnmarshalTestEnum(reader io.Reader) (TestEnum, error) {
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
		number := &Number{Inner: new(scale_codec.Integer[uint32])}
		err := number.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}

		return number, err
	default:
		return nil, fmt.Errorf("failure to decode bla bla bla")
	}
}

type Number struct {
	Inner *scale_codec.Integer[uint32]
}

func (Number) isTestEnum() {}
func (n *Number) MarshalSCALE() ([]byte, error) {
	innerEncode, err := n.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{[]byte{0}, innerEncode}, nil), nil
}

func (n *Number) UnmarshalSCALE(reader io.Reader) error {
	return n.Inner.UnmarshalSCALE(reader)
}

func TestMarshalOptionGeneric(t *testing.T) {
	marshaler := scale_codec.SomeG[TestEnum](
		&Number{Inner: &scale_codec.Integer[uint32]{10}})
	output, err := marshaler.MarshalSCALE()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte{1, 0, 10, 0, 0, 0}

	if !bytes.Equal(expected, output) {
		t.Fatalf("\nexpected: %v\nactual: %v\n", expected, output)
	}
}

func TestUnmarshalOptionGeneric(t *testing.T) {
	inputBytes := []byte{1, 0, 10, 0, 0, 0}

	unmarshaler := new(scale_codec.OptionG[TestEnum])
	err := unmarshaler.UnmarshalSCALE(bytes.NewReader(inputBytes), UnmarshalTestEnum)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := scale_codec.SomeG[TestEnum](
		&Number{Inner: &scale_codec.Integer[uint32]{10}})
	if !reflect.DeepEqual(expected, unmarshaler) {
		t.Fatalf("\nexpected: %v\nactual: %v\n", expected, unmarshaler)
	}
}
