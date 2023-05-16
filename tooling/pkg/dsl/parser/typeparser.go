// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package parser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Type struct {
	Pos   lexer.Position
	Named *TypeName  `parser:"(@@"`
	Sub   *Type      `parser:" | '(' @@ ')')"`
	Tails []TypeTail `parser:"@@*"`
}

func (t Type) String() string {
	var val string
	if t.Named != nil {
		val = t.Named.String()
	} else {
		val = t.Sub.String()
	}
	for _, tail := range t.Tails {
		val = tail.String(val)
	}

	return val
}

type TypeName struct {
	Pos      lexer.Position
	Name     string  `parser:"@Ident"`
	TypeArgs []*Type `parser:"('<' @@ (',' @@)* '>')?"`
}

func (n TypeName) String() string {
	if len(n.TypeArgs) == 0 {
		return fmt.Sprintf("'%s'", n.Name)
	}

	var args []string
	for _, arg := range n.TypeArgs {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("(Generic '%s' %s)", n.Name, strings.Join(args, " "))
}

type TypeTail struct {
	Optional bool    `parser:"  @'?'"`
	MapValue *Type   `parser:"| '-' '>' @@"`
	Vector   *Vector `parser:"| '*' @@"`
	Array    *Array  `parser:"| '[' @@"`
}

func (t TypeTail) String(target string) string {
	if t.Optional {
		return fmt.Sprintf("(Optional %s)", target)
	}
	if t.MapValue != nil {
		return fmt.Sprintf("(Map %s %s)", target, t.MapValue)
	}
	if t.Vector != nil {
		if t.Vector.Length != nil {
			return fmt.Sprintf("(Vector[%d] %s)", *t.Vector.Length, target)
		}
		return fmt.Sprintf("(Vector %s)", target)
	}
	if t.Array != nil {
		if len(t.Array.Dimensions) == 0 {
			return fmt.Sprintf("(Array %s)", target)
		}
		dims := make([]string, 0)
		for _, dim := range t.Array.Dimensions {
			args := make([]string, 0)
			if dim.Name != nil {
				args = append(args, *dim.Name)
			}
			if dim.Length != nil {
				args = append(args, fmt.Sprintf("%d", *dim.Length))
			}
			dims = append(dims, fmt.Sprintf("[%s]", strings.Join(args, " ")))
		}

		return fmt.Sprintf("(Array[%s] %s)", strings.Join(dims, ""), target)
	}

	panic("unreachable")
}

type Vector struct {
	Length *uint64 `parser:"@Int?"`
}

type Array struct {
	Dimensions []ArrayDimension `parser:"']' | ((@@ (',' @@)*) ']')"`
}

type ArrayDimension struct {
	Name   *string
	Length *uint64
}

func (a *ArrayDimension) Parse(lex *lexer.PeekingLexer) error {
	parenCount := 0
	for lex.Peek().Type == '(' {
		lex.Next()
		parenCount++
	}

	if lex.Peek().Type == scanner.Ident {
		name := lex.Next().Value
		a.Name = &name
		if lex.Peek().Type == ':' {
			lex.Next()

			l, err := parseUnit64(lex)
			if err != nil {
				return err
			}
			a.Length = &l
		}
	} else if lex.Peek().Type == scanner.Int {
		l, err := parseUnit64(lex)
		if err != nil {
			return err
		}
		a.Length = &l
	}

	for i := 0; i < parenCount; i++ {
		if lex.Peek().Type != ')' {
			return &participle.UnexpectedTokenError{
				Unexpected: *lex.Peek(),
				Expect:     `")"`,
			}
		}

		lex.Next()
	}

	return nil
}

type Pattern struct {
	Discard  bool    `parser:"  @'_'"`
	Type     *Type   `parser:"| (@@"`
	Variable *string `parser:"   @Ident?)"`
}

var (
	typeParser    = participle.MustBuild[Type]()
	patternParser = participle.MustBuild[Pattern]()
)

func ParseType(input string) (*Type, error) {
	return typeParser.ParseString("", input)
}

func ParsePattern(input string) (*Pattern, error) {
	return patternParser.ParseString("", input)
}

func parseUnit64(lex *lexer.PeekingLexer) (uint64, error) {
	if lex.Peek().Type != scanner.Int {
		return 0, &participle.UnexpectedTokenError{
			Unexpected: *lex.Peek(),
			Expect:     "integer",
		}
	}
	tok := lex.Next()
	if l, err := strconv.ParseUint(tok.Value, 10, 64); err == nil {
		return l, nil
	}
	return 0, &participle.ParseError{
		Msg: "integer out of range",
		Pos: tok.Pos,
	}

}
