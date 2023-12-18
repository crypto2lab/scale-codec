package scale_codec

import (
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
    }`

	expectedEnum := []Enum{
		{
			Name: "Nested",
			Variants: []EnumField{
				{
					Name: "Number",
					Type: "int32",
				},
			},
		},
		{
			Name: "MyEnum",
			Variants: []EnumField{
				{
					Name: "Int",
					Type: "uint64",
				},
				{
					Name: "Bool",
					Type: "bool",
				},
				{
					Name: "A",
					Type: "Option<bool>",
				},
				{
					Name: "B",
					Type: "Result<uint64,uint64>",
				},
				{
					Name: "C",
					Type: "Option<Nested>",
				},
				{
					Name: "D",
					Type: "Result<Nested,uint64>",
				},
				{
					Name: "E",
					Type: "Result<Nested,Nested>",
				},
				{
					Name: "F",
					Type: "Result<uint64,Nested>",
				},
				{
					Name: "G",
					Type: "Tuple<uint64,bool>",
				},
				{
					Name: "H",
					Type: "Option<Tuple<uint64,bool>>",
				},
				{
					Name: "J",
					Type: "Result<Tuple<uint64,bool>,bool>",
				},
				{
					Name: "K",
					Type: "Tuple<Option<bool>,Result<bool,bool>>",
				},
				{
					Name: "L",
					Type: "Result<Option<Tuple<uint64,bool>>,uint64>",
				},
			},
		},
	}

	result := ParseEnum("", strings.NewReader(input))
	if result != 0 {
		t.Fatalf("error to parse enum")
	}

	if !reflect.DeepEqual(expectedEnum, Enums) {
		t.Fatalf("\nexpected: %v\ngot: %v", expectedEnum, Enums)
	}
}
