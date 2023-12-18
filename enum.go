package scale_codec

import (
	"errors"
	"io"
	"text/scanner"
)

//go:generate goyacc -l -o enum_parser.go enum_parser.y

var ErrWrongEnumTag = errors.New("wrong enum tag")

type Enumerable interface {
	Encodable
	Index() byte
}

type yySymType struct {
	sval       string
	enum       Enum
	enumField  EnumField
	enumFields []EnumField
	yys        int
}

func ParseEnum(filename string, src io.Reader) int {
	lexer := newLexer(filename, src)
	return yyParse(lexer)
}

type lexer struct {
	s   scanner.Scanner
	err error
}

func newLexer(filename string, src io.Reader) *lexer {
	var s scanner.Scanner
	s.Init(src)
	s.Filename = filename

	return &lexer{s: s}
}

func (l *lexer) Error(msg string) {
	l.err = errors.New(msg)
}

func (l *lexer) Lex(lval *yySymType) int {
	token := l.s.Scan()
	if token == scanner.EOF {
		return -1
	}

	lexeme := l.s.TokenText()

	switch lexeme {
	case "enum":
		return ENUM
	case "{", "}", "(", ")":
		return int(rune(lexeme[0]))
	case "uint64", "bool":
		lval.sval = lexeme
		return TYPE
	default:
		lval.sval = lexeme
		return IDENTIFIER
	}
}
