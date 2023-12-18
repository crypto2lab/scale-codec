package scale_codec

import (
	"fmt"
	"strings"
	"testing"
)

func TestEnumParser(t *testing.T) {
	const input = `enum MyEnum {
        Test(uint64)
        Other(bool)
        Another(bool)
    }`

	lexer := newLexer("", strings.NewReader(input))

	yyParse(lexer)
	fmt.Println(lexer.err)
}
