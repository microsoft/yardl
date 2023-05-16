// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Expression interface {
	fmt.Stringer
	_expression()
}

type IntegerLiteral struct {
	Pos   lexer.Position
	Value string `parser:"@Int"`
}

func (IntegerLiteral) _expression() {}

func (n IntegerLiteral) String() string {
	return n.Value
}

type StringLiteral struct {
	Pos lexer.Position

	Value string `parser:"@String"`
}

func (StringLiteral) _expression() {}

func (n StringLiteral) String() string {
	return fmt.Sprintf("«%s»", n.Value)
}

type FunctionCall struct {
	Pos lexer.Position

	FunctionName string       `parser:"@Ident"`
	Arguments    []Expression `parser:"'(' (@@ (',' @@)*)? ')'"`
}

func (n FunctionCall) String() string {
	var args []string
	for _, arg := range n.Arguments {
		args = append(args, arg.String())
	}

	argString := strings.Join(args, ", ")
	if argString != "" {
		argString = " " + argString
	}

	return fmt.Sprintf("(call %s%s)", n.FunctionName, argString)
}

func (FunctionCall) _expression() {}

type PathExpr struct {
	Pos lexer.Position

	Parts []PathPart `parser:"@@ ( '.' @@ )*"`
}

func (PathExpr) _expression() {}

func (n PathExpr) String() string {
	var parts []string
	for _, part := range n.Parts {
		parts = append(parts, part.String())
	}
	return strings.Join(parts, ".")
}

type PathPart struct {
	Pos lexer.Position

	Name    string  `parser:"@Ident"`
	Indexes []Index `parser:"('[' @@ ']')*"`
}

func (p PathPart) String() string {
	var indexes []string
	for _, index := range p.Indexes {
		indexes = append(indexes, index.String())
	}
	return fmt.Sprintf("(part %s %s)", p.Name, strings.Join(indexes, ""))
}

type Index struct {
	Pos lexer.Position

	IndexArgs []IndexArg `parser:"(@@ (',' @@)*)?"`
}

func (i Index) String() string {
	var args []string
	for _, arg := range i.IndexArgs {
		args = append(args, arg.String())
	}
	argString := strings.Join(args, ", ")
	if argString != "" {
		argString = " " + argString
	}

	return fmt.Sprintf("(index%s)", argString)
}

type IndexArg struct {
	Pos lexer.Position

	Label *Label     `parser:"(@@ ':')?"`
	Value Expression `parser:"@@"`
}

func (a IndexArg) String() string {
	if a.Label != nil {
		return fmt.Sprintf("%s:%s", a.Label, a.Value)
	}
	return a.Value.String()
}

type Label struct {
	Pos lexer.Position

	Name string `parser:"@Ident"`
}

func (l Label) String() string {
	return l.Name
}

var (
	expressionLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "String", Pattern: `("([^\\"]|\\.)*")|('(\\'|[^'])*')`},
		{Name: "Int", Pattern: `[-+]?((0[xX][0-9A-Fa-f]+)|\d+)`},
		{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{Name: "whitespace", Pattern: `[ \t]+`},
	})

	expressionParser = participle.MustBuild[Expression](
		participle.Lexer(expressionLexer),
		participle.Unquote("String"),
		participle.Union[Expression](
			IntegerLiteral{},
			StringLiteral{},
			FunctionCall{},
			PathExpr{}))
)

func ParseExpression(input string) (*Expression, error) {
	return expressionParser.ParseString("", input)
}

type ParseError participle.Error
