package scale_codec_test

import (
	"bytes"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestBoolMarshaler(t *testing.T) {
	cases := []struct {
		enc           scale_codec.Marshaler
		expectedBytes []byte
	}{
		{
			enc:           scale_codec.Bool{true},
			expectedBytes: []byte{0x01},
		},
		{
			enc:           scale_codec.Bool{false},
			expectedBytes: []byte{0x00},
		},
		{
			enc:           scale_codec.OptionBool{},
			expectedBytes: []byte{0x00},
		},
		{
			enc:           scale_codec.OptionBool{&scale_codec.Bool{true}},
			expectedBytes: []byte{0x01},
		},
		{
			enc:           scale_codec.OptionBool{&scale_codec.Bool{false}},
			expectedBytes: []byte{0x02},
		},
	}

	for _, tt := range cases {
		actual, err := tt.enc.MarshalSCALE()
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if !bytes.EqualFold(tt.expectedBytes, actual) {
			t.Errorf("expected %v, got %v", tt.expectedBytes, actual)
		}
	}
}

func TestBoolUnmarshaler(t *testing.T) {
	cases := []struct {
		unmarshaler    scale_codec.Unmarshaler
		inputBytes     []byte
		expectedOutput any
	}{
		{
			unmarshaler:    new(scale_codec.Bool),
			inputBytes:     []byte{0x01},
			expectedOutput: &scale_codec.Bool{true},
		},
		{
			unmarshaler:    new(scale_codec.Bool),
			inputBytes:     []byte{0x00},
			expectedOutput: &scale_codec.Bool{false},
		},
		{
			unmarshaler:    new(scale_codec.OptionBool),
			inputBytes:     []byte{0x00},
			expectedOutput: &scale_codec.OptionBool{nil},
		},
		{
			unmarshaler:    new(scale_codec.OptionBool),
			inputBytes:     []byte{0x01},
			expectedOutput: &scale_codec.OptionBool{&scale_codec.Bool{true}},
		},
		{
			unmarshaler:    new(scale_codec.OptionBool),
			inputBytes:     []byte{0x02},
			expectedOutput: &scale_codec.OptionBool{&scale_codec.Bool{false}},
		},
	}

	for _, tt := range cases {
		err := tt.unmarshaler.UnmarshalSCALE(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if !reflect.DeepEqual(tt.unmarshaler, tt.expectedOutput) {
			t.Errorf("\nexpected: %v\ngot: %v", tt.expectedOutput, tt.unmarshaler)
		}
	}
}
