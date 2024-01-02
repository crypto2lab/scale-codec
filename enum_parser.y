%{
package scale_codec

import (
    "fmt"
    "strings"
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

var (
    Enums []Enum

    // map of tuple name and tuple qty of values
    GenericTuple map[string]int = make(map[string]int)
)

%}

%token ENUM
%token IDENTIFIER
%token TYPE
%token RESULT
%token OPTION

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
    GenericTuple["T" + fmt.Sprint(len($2.tupleValuesTypes))] = len($2.tupleValuesTypes)
    genericTuple := "T" + fmt.Sprint(len($2.tupleValuesTypes)) + "[" + strings.Join($2.tupleValuesTypes, ",") + "]"
    $$.ttype = "*" + genericTuple
    $$.sval = "new(" + genericTuple + ")"
    $$.fromRawBytesFunc = "UnmarshalT"+ fmt.Sprint(len($2.tupleValuesTypes)) + 
        "FromRawBytes[" + strings.Join($2.tupleValuesTypes, ",") + 
        "](" + strings.Join($2.tupleValuesUnmarshalScale, ",") + ")"
    $$.unmarshalScale = "return i.Inner.UnmarshalSCALE(reader," + strings.Join($2.tupleValuesUnmarshalScale, ",") + ")"
};

TypeList: IDENTIFIER {
    $$.sval = $1.sval
    $$.tupleValuesTypes = append($$.tupleValuesTypes, $1.sval)
    $$.tupleValuesUnmarshalScale = append($$.tupleValuesUnmarshalScale, "Unmarshal"+$1.sval)
} | ComplexType {
    $$.sval = $1.sval
    $$.tupleValuesTypes = append($$.tupleValuesTypes, $1.ttype)
    $$.tupleValuesUnmarshalScale = append($$.tupleValuesUnmarshalScale, $1.fromRawBytesFunc)
} | TypeList "," IDENTIFIER {
    $$.sval += "," + $3.sval
    $$.tupleValuesTypes = append($$.tupleValuesTypes, $3.sval)
    $$.tupleValuesUnmarshalScale = append($$.tupleValuesUnmarshalScale, "Unmarshal"+$3.sval)
} | TypeList "," ComplexType {
    $$.sval += "," + $3.sval
    $$.tupleValuesTypes = append($$.tupleValuesTypes, $3.ttype)
    $$.tupleValuesUnmarshalScale = append($$.tupleValuesUnmarshalScale, $3.fromRawBytesFunc)
};

Option: OPTION "<" ComplexType ">" {
    $$.sval = "new(scale_codec.OptionG[" + $3.ttype + "])"
    $$.ttype = "*scale_codec.OptionG[" + $3.ttype + "]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalOptionFromRawBytes[" + $3.ttype + "](" + $3.fromRawBytesFunc + ")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader," + $3.fromRawBytesFunc + ")"
} | OPTION "<" IDENTIFIER ">" {
    $$.sval = "new(scale_codec.OptionG[" + $3.sval + "])"
    $$.ttype = "*scale_codec.OptionG[" + $3.sval + "]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalOptionFromRawBytes[" + $3.sval + "](Unmarshal"+ $3.sval +")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ")"
} ;

Result: RESULT "<" ComplexType "," ComplexType ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.ttype + "," + $5.ttype + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.ttype + "," + $5.ttype +"]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalResultFromRawBytes[" + 
        $3.ttype + ","+ $5.ttype +"]("+ $3.fromRawBytesFunc +","+ $5.fromRawBytesFunc +")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, " + $3.fromRawBytesFunc + ", " + $5.fromRawBytesFunc + ")"
} | RESULT "<" IDENTIFIER "," ComplexType ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.sval + "," + $5.ttype + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.sval + "," + $5.ttype +"]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalResultFromRawBytes[" + 
        $3.sval + ","+ $5.ttype +"](Unmarshal"+ $3.sval +","+ $5.fromRawBytesFunc +")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ", " + $5.fromRawBytesFunc + ")"
} | RESULT "<" ComplexType "," IDENTIFIER ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.ttype + "," + $5.sval + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.ttype + "," + $5.sval +"]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalResultFromRawBytes[" + 
        $3.ttype + ","+ $5.sval +"]("+ $3.fromRawBytesFunc +",Unmarshal"+ $5.sval +")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, " + $3.fromRawBytesFunc + ", Unmarshal" + $5.sval + ")"
} | RESULT "<" IDENTIFIER "," IDENTIFIER ">" {
    $$.sval = "new(scale_codec.ResultG[" + $3.sval + "," + $5.sval + "])"
    $$.ttype = "*scale_codec.ResultG[" + $3.sval + "," + $5.sval +"]"
    $$.fromRawBytesFunc = "scale_codec.UnmarshalResultFromRawBytes[" + 
        $3.sval + ","+ $5.sval +"](Unmarshal"+ $3.sval +",Unmarshal"+ $5.sval +")"
    $$.unmarshalScale  = "return i.Inner.UnmarshalSCALE(reader, Unmarshal" + $3.sval + ", Unmarshal" + $5.sval + ")"
} ;

%%

func hasGenericArgs(args []string) bool {
    for _, input := range args {
        if !strings.HasPrefix(input, "*"){
            return true
        }
    }

    return false
}

