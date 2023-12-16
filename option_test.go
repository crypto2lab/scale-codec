package scale_codec_test

import (
	"bytes"
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
