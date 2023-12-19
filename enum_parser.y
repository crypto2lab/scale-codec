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
    $$.enumField = EnumField{Name: $1.sval, Type: ""}
} | IDENTIFIER "(" ComplexType ")" {
    $$.enumField = EnumField{Name: $1.sval, Type: $3.sval}
};


ComplexType: TYPE
    | Tuple
    | Option
    | Result ;


Tuple: "(" TypeList ")" {
    $$.sval = "Tuple<" + $2.sval + ">"
};

TypeList: ComplexType {
    $$.sval = $1.sval
} | TypeList "," ComplexType {
    $$.sval += "," + $3.sval
} ;

Option: TYPE "<" ComplexType ">" {
    $$.sval = $1.sval + "<" + $3.sval + ">"
} | TYPE "<" IDENTIFIER ">" {
    $$.sval = $1.sval + "<" + $3.sval + ">"
} ;

Result: TYPE "<" ComplexType "," ComplexType ">" {
        $$.sval = $1.sval + "<" + $3.sval + "," + $5.sval + ">"
    } 
    | TYPE "<" IDENTIFIER "," ComplexType ">" {
        $$.sval = $1.sval + "<" + $3.sval + "," + $5.sval + ">"
    } 
    | TYPE "<" ComplexType "," IDENTIFIER ">" {
        $$.sval = $1.sval + "<" + $3.sval + "," + $5.sval + ">"
    } 
    | TYPE "<" IDENTIFIER "," IDENTIFIER ">" {
        $$.sval = $1.sval + "<" + $3.sval + "," + $5.sval + ">"
    } ;

%%
