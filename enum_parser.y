%{
package scale_codec

import (
    "fmt"
)

type Enum struct {
    Name string
    Variants []EnumField
}

type EnumField struct {
    Name            string
    Type            string
    TypeConstructor string
	UnmarshalScale  string
}

var Enums []Enum

%}

%token ENUM
%token IDENTIFIER
%token TYPE

%%

Enums: /* empty */ {
    Enums = make([]Enum, 0)
} | Enums Enum

Enum: ENUM IDENTIFIER "{" EnumFields "}" {
    $$.enum = Enum{Name: $2.sval, Variants: $4.enumFields}
    Enums = append(Enums, $$.enum)
    fmt.Printf("Parsed enum: %s with fields: %+v\n", $2.sval, Enums)
};

EnumFields: /* empty */ {
    $$.enumFields = nil // Initialize as an empty slice
} | EnumFields EnumField {
    $$.enumFields = append($1.enumFields, $2.enumField)
};

EnumField: IDENTIFIER {
    $$.enumField = EnumField{
        Name: $1.sval, 
        Type: "*scale_codec.SimpleVariant", 
        TypeConstructor: "new(scale_codec.SimpleVariant)",
        UnmarshalScale: "",
    }
} | IDENTIFIER "(" ComplexType ")" {
    $$.enumField = EnumField{
        Name: $1.sval, 
        Type: $3.ttype, 
        TypeConstructor: $3.sval,
        UnmarshalScale: $3.unmarshalScale,
    }
};


ComplexType: TYPE {
    $$.sval = "new(" + $1.sval + ")"
    $$.ttype = "*" + $1.sval
    $$.fromRawBytesFunc = $1.fromRawBytesFunc
} | Tuple | Option | Result ;


Tuple: "(" TypeList ")" {
    $$.sval = "scale_codec.NewTuple(" + $2.sval + ")"
    $$.ttype = "*scale_codec.Tuple"
};

TypeList: ComplexType {
    $$.sval = $1.sval
} | TypeList "," ComplexType {
    $$.sval += "," + $3.sval
} ;

Option: TYPE "<" ComplexType ">" {
    $$.sval = "scale_codec.NewOption(" + $3.sval + ")"
    $$.ttype = "*scale_codec.Option"
} | TYPE "<" IDENTIFIER ">" {
    $$.sval = "new(scale_codec.OptionG[" + $3.sval + "])"
    $$.ttype = "*scale_codec.OptionG[" + $3.sval + "]"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ")"
} ;

Result: TYPE "<" ComplexType "," ComplexType ">" {
    $$.sval = "scale_codec.NewResult(" + $3.sval + "," + $5.sval + ")"
    $$.ttype = "*scale_codec.Result"
} | TYPE "<" IDENTIFIER "," ComplexType ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.sval + "," + $5.ttype + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.sval + "," + $5.ttype +"]"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ", " + $5.fromRawBytesFunc + ")"
} | TYPE "<" ComplexType "," IDENTIFIER ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.ttype + "," + $5.sval + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.ttype + "," + $5.sval +"]"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, " + $3.fromRawBytesFunc + ", Unmarshal" + $5.sval + ")"
} | TYPE "<" IDENTIFIER "," IDENTIFIER ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.sval + "," + $5.sval + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.sval + "," + $5.sval +"]"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ", Unmarshal" + $5.sval + ")"
} ;

%%
