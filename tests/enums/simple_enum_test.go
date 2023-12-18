package main

import (
	"bytes"
	"reflect"
	"testing"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func TestSimpleEnumMarshaler(t *testing.T) {

}

func TestSimpleEnumUnmarshaler(t *testing.T) {
	cases := []struct {
		inputBytes      []byte
		expectedVariant MyScaleEncodedEnum
	}{
		{
			inputBytes: []byte{0, 32, 0, 0, 0, 0, 0, 0, 0},
			expectedVariant: &Int{
				Inner: &scale_codec.Integer[uint64]{Value: 32},
			},
		},
		{
			inputBytes: []byte{1, 1},
			expectedVariant: &Bool{
				Inner: &scale_codec.Bool{Value: true},
			},
		},
	}

	for _, tt := range cases {
		variant, err := UnmarhalMyScaleEncodedEnum(bytes.NewReader(tt.inputBytes))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tt.expectedVariant, variant) {
			t.Fatalf("\nexpected: %v (%T)\ngot: %v (%T)", tt.expectedVariant, tt.expectedVariant,
				variant, variant)
		}
	}
}
