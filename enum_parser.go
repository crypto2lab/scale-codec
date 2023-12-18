// Code generated by goyacc -l -o enum_parser.go enum_parser.y. DO NOT EDIT.
package scale_codec

import __yyfmt__ "fmt"

import (
	"fmt"
)

type Enum struct {
	Name     string
	Variants []EnumField
}

type EnumField struct {
	Name string
	Type string
}

var Enums []Enum

const ENUM = 57346
const IDENTIFIER = 57347
const TYPE = 57348

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ENUM",
	"IDENTIFIER",
	"TYPE",
	"\"{\"",
	"\"}\"",
	"\"(\"",
	"\")\"",
	"\",\"",
	"\"<\"",
	"\">\"",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 43

var yyAct = [...]int8{
	11, 28, 26, 27, 25, 37, 36, 35, 34, 23,
	24, 17, 18, 12, 33, 12, 16, 20, 16, 21,
	31, 12, 22, 12, 16, 29, 16, 30, 10, 32,
	9, 5, 4, 7, 3, 19, 15, 14, 13, 8,
	6, 2, 1,
}

var yyPact = [...]int16{
	-1000, 30, -1000, 27, 24, -1000, 25, -1000, -1000, 19,
	7, 1, 0, -1000, -1000, -1000, 7, -1000, 17, -1,
	-1000, -9, -10, -1000, 7, -1000, 15, -1000, 9, -1000,
	-5, -6, -7, -8, -1000, -1000, -1000, -1000,
}

var yyPgo = [...]int8{
	0, 42, 41, 40, 39, 0, 38, 37, 36, 35,
}

var yyR1 = [...]int8{
	0, 1, 1, 2, 3, 3, 4, 5, 5, 5,
	5, 6, 9, 9, 7, 7, 8, 8, 8, 8,
}

var yyR2 = [...]int8{
	0, 0, 2, 5, 0, 2, 4, 1, 1, 1,
	1, 3, 1, 3, 4, 4, 6, 6, 6, 6,
}

var yyChk = [...]int16{
	-1000, -1, -2, 4, 5, 7, -3, 8, -4, 5,
	9, -5, 6, -6, -7, -8, 9, 10, 12, -9,
	-5, -5, 5, 10, 11, 13, 11, 13, 11, -5,
	-5, 5, -5, 5, 13, 13, 13, 13,
}

var yyDef = [...]int8{
	1, -2, 2, 0, 0, 4, 0, 3, 5, 0,
	0, 0, 7, 8, 9, 10, 0, 6, 0, 0,
	12, 0, 0, 11, 0, 14, 0, 15, 0, 13,
	0, 0, 0, 0, 16, 18, 17, 19,
}

var yyTok1 = [...]int8{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	9, 10, 3, 3, 11, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	12, 3, 13, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 8,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
		{
			Enums = make([]Enum, 0)
		}
	case 3:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.enum = Enum{Name: yyDollar[2].sval, Variants: yyDollar[4].enumFields}
			Enums = append(Enums, yyVAL.enum)
			fmt.Printf("Parsed enum: %s with fields: %+v\n", yyDollar[2].sval, Enums)
		}
	case 4:
		yyDollar = yyS[yypt-0 : yypt+1]
		{
			yyVAL.enumFields = nil // Initialize as an empty slice
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.enumFields = append(yyDollar[1].enumFields, yyDollar[2].enumField)
		}
	case 6:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.enumField = EnumField{Name: yyDollar[1].sval, Type: yyDollar[3].sval}
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.sval = "Tuple<" + yyDollar[2].sval + ">"
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.sval += "," + yyDollar[3].sval
		}
	case 14:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + ">"
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + ">"
		}
	case 16:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + "," + yyDollar[5].sval + ">"
		}
	case 17:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + "," + yyDollar[5].sval + ">"
		}
	case 18:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + "," + yyDollar[5].sval + ">"
		}
	case 19:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.sval = yyDollar[1].sval + "<" + yyDollar[3].sval + "," + yyDollar[5].sval + ">"
		}
	}
	goto yystack /* stack new state and value */
}
