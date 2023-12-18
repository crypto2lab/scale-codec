%{
package scale_codec

import (
    "fmt"
)

type EnumField struct {
    Name string
    Type string
}

var enums []EnumField

%}

%token ENUM
%token IDENTIFIER
%token TYPE

%%

Enum: ENUM IDENTIFIER "{" EnumFields "}" {
    fmt.Printf("Parsed enum: %s with fields: %+v\n", $2.sval, enums)
};

EnumFields: /* empty */ 
    | EnumFields EnumField ;

EnumField: IDENTIFIER "(" TYPE ")" {
    enums = append(enums, EnumField{Name: $1.sval, Type: $3.sval})
};

%%
