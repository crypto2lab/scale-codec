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
    Name string
    Type string
    TypeConstructor string
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
    $$.enumField = EnumField{Name: $1.sval, Type: "*scale_codec.SimpleVariant", TypeConstructor: "new(scale_codec.SimpleVariant)"}
} | IDENTIFIER "(" ComplexType ")" {
    $$.enumField = EnumField{Name: $1.sval, Type: $3.ttype, TypeConstructor: $3.sval}
};


ComplexType: TYPE {
    $$.sval = "new(" + $1.sval + ")"
    $$.ttype = "*" + $1.sval
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
    $$.sval = "scale_codec.NewOption(" + $3.sval + ")"
    $$.ttype = "*scale_codec.Option"
} ;

Result: TYPE "<" ComplexType "," ComplexType ">" {
    $$.sval = "scale_codec.NewResult(" + $3.sval + "," + $5.sval + ")"
    $$.ttype = "*scale_codec.Result"
} | TYPE "<" IDENTIFIER "," ComplexType ">" {
    $$.sval = "scale_codec.NewResult(" + $3.sval + "," + $5.sval + ")"
    $$.ttype = "*scale_codec.Result"
} | TYPE "<" ComplexType "," IDENTIFIER ">" {
    $$.sval = "scale_codec.NewResult(" + $3.sval + "," + $5.sval + ")"
    $$.ttype = "*scale_codec.Result"
} | TYPE "<" IDENTIFIER "," IDENTIFIER ">" {
    $$.sval = "scale_codec.NewResult(" + $3.sval + "," + $5.sval + ")"
    $$.ttype = "*scale_codec.Result"
} ;

%%
