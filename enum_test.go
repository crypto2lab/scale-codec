package scale_codec

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEnumParser(t *testing.T) {
	// to extract better info from the goyacc parser set:
	//yyDebug = 50
	const input = `
	enum Nested {
		Number(int32)
	}

	enum MyEnum {
		Single
        Int(uint64)
        Bool(bool)
		A(Option<bool>)
		B(Result<uint64, uint64>)
		C(Option<Nested>)
		D(Result<Nested, uint64>)
		E(Result<Nested, Nested>)
		F(Result<uint64, Nested>)
		G((uint64, bool))
		H(Option<(uint64, bool)>)
		J(Result<(uint64, bool), bool>)
		K((Option<bool>, Result<bool, bool>))
		L(Result<Option<(uint64, bool)>, uint64>)
		M(Option<Nested>)
		N(Result<Nested, bool>)
		O(Result<bool, Nested>)
		P(Result<Nested, Error>)
		Q((Nested, uint64, Error))
		R((Result<uint64, bool>, Option<uint64>, Error))
    }`

	expectedEnum := []Enum{
		{
			Name: "Nested",
			Variants: []EnumField{
				{
					Name:            "Number",
					Type:            "*scale_codec.Integer[int32]",
					TypeConstructor: "new(scale_codec.Integer[int32])",
				},
			},
		},
		{
			Name: "MyEnum",
			Variants: []EnumField{
				{
					Name:            "Single",
					Type:            "*scale_codec.SimpleVariant",
					TypeConstructor: "new(scale_codec.SimpleVariant)",
				},
				{
					Name:            "Int",
					Type:            "*scale_codec.Integer[uint64]",
					TypeConstructor: "new(scale_codec.Integer[uint64])",
				},
				{
					Name:            "Bool",
					Type:            "*scale_codec.Bool",
					TypeConstructor: "new(scale_codec.Bool)",
				},
				{
					Name:            "A",
					Type:            "*scale_codec.OptionG[*scale_codec.Bool]",
					TypeConstructor: "new(scale_codec.OptionG[*scale_codec.Bool])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,scale_codec.BoolFromRawBytes)",
				},
				{
					Name:            "B",
					Type:            "*scale_codec.ResultG[*scale_codec.Integer[uint64],*scale_codec.Integer[uint64]]",
					TypeConstructor: "new(scale_codec.ResultG[*scale_codec.Integer[uint64],*scale_codec.Integer[uint64]])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, scale_codec.IntegerFromRawBytes[uint64], scale_codec.IntegerFromRawBytes[uint64])",
				},
				{
					Name:            "C",
					Type:            "*scale_codec.OptionG[Nested]",
					TypeConstructor: "new(scale_codec.OptionG[Nested])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested)",
				},
				{
					Name:            "D",
					Type:            "*scale_codec.ResultG[Nested,*scale_codec.Integer[uint64]]",
					TypeConstructor: "new(scale_codec.ResultG[Nested,*scale_codec.Integer[uint64]])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested, scale_codec.IntegerFromRawBytes[uint64])",
				},
				{
					Name:            "E",
					Type:            "*scale_codec.ResultG[Nested,Nested]",
					TypeConstructor: "new(scale_codec.ResultG[Nested,Nested])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested, UnmarshalNested)",
				},
				{
					Name:            "F",
					Type:            "*scale_codec.ResultG[*scale_codec.Integer[uint64],Nested]",
					TypeConstructor: "new(scale_codec.ResultG[*scale_codec.Integer[uint64],Nested])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, scale_codec.IntegerFromRawBytes[uint64], UnmarshalNested)",
				},
				{
					Name:            "G",
					Type:            "*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]",
					TypeConstructor: "new(T2[*scale_codec.Integer[uint64],*scale_codec.Bool])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,scale_codec.IntegerFromRawBytes[uint64],scale_codec.BoolFromRawBytes)",
				},
				{
					Name:            "H",
					Type:            "*scale_codec.OptionG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]]",
					TypeConstructor: "new(scale_codec.OptionG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,UnmarshalT2FromRawBytes[*scale_codec.Integer[uint64],*scale_codec.Bool](scale_codec.IntegerFromRawBytes[uint64],scale_codec.BoolFromRawBytes))",
				},
				{
					Name:            "J",
					Type:            "*scale_codec.ResultG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool],*scale_codec.Bool]",
					TypeConstructor: "new(scale_codec.ResultG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool],*scale_codec.Bool])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalT2FromRawBytes[*scale_codec.Integer[uint64],*scale_codec.Bool](scale_codec.IntegerFromRawBytes[uint64],scale_codec.BoolFromRawBytes), scale_codec.BoolFromRawBytes)",
				},
				{
					Name:            "K",
					Type:            "*T2[*scale_codec.OptionG[*scale_codec.Bool],*scale_codec.ResultG[*scale_codec.Bool,*scale_codec.Bool]]",
					TypeConstructor: "new(T2[*scale_codec.OptionG[*scale_codec.Bool],*scale_codec.ResultG[*scale_codec.Bool,*scale_codec.Bool]])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,scale_codec.UnmarshalOptionFromRawBytes[*scale_codec.Bool](scale_codec.BoolFromRawBytes),scale_codec.UnmarshalResultFromRawBytes[*scale_codec.Bool,*scale_codec.Bool](scale_codec.BoolFromRawBytes,scale_codec.BoolFromRawBytes))",
				},
				{
					Name:            "L",
					Type:            "*scale_codec.ResultG[*scale_codec.OptionG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]],*scale_codec.Integer[uint64]]",
					TypeConstructor: "new(scale_codec.ResultG[*scale_codec.OptionG[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]],*scale_codec.Integer[uint64]])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, scale_codec.UnmarshalOptionFromRawBytes[*T2[*scale_codec.Integer[uint64],*scale_codec.Bool]](UnmarshalT2FromRawBytes[*scale_codec.Integer[uint64],*scale_codec.Bool](scale_codec.IntegerFromRawBytes[uint64],scale_codec.BoolFromRawBytes)), scale_codec.IntegerFromRawBytes[uint64])",
				},
				{
					Name:            "M",
					Type:            "*scale_codec.OptionG[Nested]",
					TypeConstructor: "new(scale_codec.OptionG[Nested])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested)",
				},
				{
					Name:            "N",
					Type:            "*scale_codec.ResultG[Nested,*scale_codec.Bool]",
					TypeConstructor: "new(scale_codec.ResultG[Nested,*scale_codec.Bool])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested, scale_codec.BoolFromRawBytes)",
				},
				{
					Name:            "O",
					Type:            "*scale_codec.ResultG[*scale_codec.Bool,Nested]",
					TypeConstructor: "new(scale_codec.ResultG[*scale_codec.Bool,Nested])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, scale_codec.BoolFromRawBytes, UnmarshalNested)",
				},
				{
					Name:            "P",
					Type:            "*scale_codec.ResultG[Nested,Error]",
					TypeConstructor: "new(scale_codec.ResultG[Nested,Error])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader, UnmarshalNested, UnmarshalError)",
				},
				{
					Name:            "Q",
					Type:            "*T3[Nested,*scale_codec.Integer[uint64],Error]",
					TypeConstructor: "new(T3[Nested,*scale_codec.Integer[uint64],Error])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,UnmarshalNested,scale_codec.IntegerFromRawBytes[uint64],UnmarshalError)",
				},
				{
					Name:            "R",
					Type:            "*T3[*scale_codec.ResultG[*scale_codec.Integer[uint64],*scale_codec.Bool],*scale_codec.OptionG[*scale_codec.Integer[uint64]],Error]",
					TypeConstructor: "new(T3[*scale_codec.ResultG[*scale_codec.Integer[uint64],*scale_codec.Bool],*scale_codec.OptionG[*scale_codec.Integer[uint64]],Error])",
					UnmarshalScale:  "return i.Inner.UnmarshalSCALE(reader,scale_codec.UnmarshalResultFromRawBytes[*scale_codec.Integer[uint64],*scale_codec.Bool](scale_codec.IntegerFromRawBytes[uint64],scale_codec.BoolFromRawBytes),scale_codec.UnmarshalOptionFromRawBytes[*scale_codec.Integer[uint64]](scale_codec.IntegerFromRawBytes[uint64]),UnmarshalError)",
				},
			},
		},
	}

	result := ParseEnum("", strings.NewReader(input))
	if result != 0 {
		t.Fatalf("error to parse enum")
	}

	for i, expected := range expectedEnum {
		actual := Enums[i]
		if expected.Name != actual.Name {
			t.Fatalf("\nexpected: %v\ngot: %v", expected.Name, actual.Name)
		}

		for j, variant := range expected.Variants {
			actualVariant := actual.Variants[j]
			if !reflect.DeepEqual(variant, actualVariant) {
				t.Fatalf("\nexpected: %v\ngot: %v", variant, actualVariant)
			}
		}
	}

	fmt.Println(GenericTuple)
}
