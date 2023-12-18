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
}

var Enums []Enum

%}

%token ENUM
%token IDENTIFIER
%token TYPE

%%

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

EnumField: IDENTIFIER "(" TYPE ")" {
    $$.enumField = EnumField{Name: $1.sval, Type: $3.sval}
};

%%
