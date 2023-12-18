package scale_codec

import (
	"errors"
	"io"
	"text/scanner"
)

//go:generate goyacc -l -o enum_parser.go enum_parser.y

type Enumerable interface {
	Encodable
	Index() uint
}

type yySymType struct {
	sval string
	yys  int
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
