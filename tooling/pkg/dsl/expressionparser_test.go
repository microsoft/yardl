// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input         string
		errorContains string
		parsed        string
	}{
		{`1`, "", "1"},
		{`+1`, "", "1"},
		{`+ 1`, "", "1"},
		{`+(1 + 2)`, "", "(+ 1 2)"},
		{`-1`, "", "-1"},
		{`-(1 + 2)`, "", "(- (+ 1 2))"},
		{`0x10`, "", "16"},
		{`0XAAb`, "", "2731"},
		{`-0x10`, "", "-16"},
		{`0x10478349823749782398327643786348762387326`, "", "1487018271446938533842446339972202262126988391206"},
		{`-0x10478349823749782398327643786348762387326`, "", "-1487018271446938533842446339972202262126988391206"},
		{`1.2`, "", "1.2"},
		{`.1`, "", ".1"},
		{`1.`, "", "1."},
		{`1.2e-3`, "", "1.2e-3"},
		{`1.2e3`, "", "1.2e3"},
		{`1e-3`, "", "1e-3"},
		{`1.2 as float32`, "", "(as 1.2 float32)"},
		{`(1 + 2) as float32`, "", "(as (+ 1 2) float32)"},
		{`(1 + 2) as (float32)`, "", "(as (+ 1 2) float32)"},
		{`1 + 2 as float32`, "", "(+ 1 (as 2 float32))"},
		{`1 * 2 as float32`, "", "(* 1 (as 2 float32))"},
		{`1 ** 2 as float32`, "", "(** 1 (as 2 float32))"},
		{`a.b as float32`, "", "(as (. a b) float32)"},
		{`a.b[] as float32`, "", "(as (subscript (. a b)) float32)"},
		{`(float32)`, "", "float32"},
		{`"s"`, "", "«s»"},
		{`"s\"s"`, "", `«s"s»`},
		{`"s\\s"`, "", `«s\s»`},
		{"'s'", "", `«s»`},
		{`'s\'s'`, "", `«s's»`},
		{`'s\\s'`, "", `«s\s»`},
		{`'s"s'`, "", `«s"s»`},
		{`"s's"`, "", `«s's»`},
		{" foo ( 1 , 2 ) ", "", "(call foo 1 2)"},
		{"foo(1,) ", `unexpected token ")" (expected an expression)`, ""},
		{"foo(,1) ", `unexpected token "," (expected an expression)`, ""},
		{"foo()", "", "(call foo)"},
		{"foo(bar())", "", "(call foo (call bar))"},
		{"foo(bar(),1)", "", "(call foo (call bar) 1)"},
		{"foo(bar(),)", `unexpected token ")" (expected an expression)`, ""},
		{"foo", "", "foo"},
		{"foo.bar", "", "(. foo bar)"},
		{"foo.bar.baz", "", "(. (. foo bar) baz)"},
		{"foo.(1+2)", `The right-hand side of a '.' operator must be an identifier`, ""},
		{"foo[]", "", "(subscript foo)"},
		{"foo[1]", "", "(subscript foo 1)"},
		{"foo[1, 2]", "", "(subscript foo 1 2)"},
		{"foo[x:1]", "", "(subscript foo x:1)"},
		{"foo[x:1, y:2]", "", "(subscript foo x:1 y:2)"},
		{"foo[1.1 as int]", "", "(subscript foo (as 1.1 int))"},
		{"foo[x: 1.1 as int, y:2]", "", "(subscript foo x:(as 1.1 int) y:2)"},
		{"foo[1,]", `unexpected token "]" (expected a subscript argument)`, ""},
		{"foo[,1]", `unexpected token "," (expected a subscript argument)`, ""},
		{"foo[1][2]", "", "(subscript (subscript foo 1) 2)"},
		{"foo[1].bar[2]", "", "(subscript (. (subscript foo 1) bar) 2)"},
		{"foo[bar()]", "", "(subscript foo (call bar))"},
		{"[1]", `unexpected token "["`, ""},
		{"1 + 1", "", "(+ 1 1)"},
		{"1 - 1", "", "(- 1 1)"},
		{"1 + 2 - 3", "", "(- (+ 1 2) 3)"},
		{"1 * 2 / 3", "", "(/ (* 1 2) 3)"},
		{"1 * 2 / 3 + 4", "", "(+ (/ (* 1 2) 3) 4)"},
		{"1 + 2 * 3", "", "(+ 1 (* 2 3))"},
		{"(1 + 2) * 3", "", "(* (+ 1 2) 3)"},
		{"2 ** 2 ** 3", "", "(** 2 (** 2 3))"},
		{"1 * 2 ** 3 + 2", "", "(+ (* 1 (** 2 3)) 2)"},

		{"Foo(1 + 2, a[3 * 4])", "", "(call Foo (+ 1 2) (subscript a (* 3 4)))"},
		{"1 ** x() + 2", "", "(+ (** 1 (call x)) 2)"},

		{`"a`, `unterminated string`, ""},
		{`%^`, `unexpected token "%"`, ``},
		{`99bottles`, `unexpected token "bottles"`, ``},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			exp, err := ParseExpression(tt.input, 1, 1)
			if tt.errorContains == "" {
				if err != nil {
					assert.NoError(t, err)
				} else {
					assert.Equal(t, tt.parsed, expressionToString(exp))
				}
			} else {
				if err == nil {
					require.Error(t, err, "got parse result when an error was expected: %v", expressionToString(exp))
				}

				require.IsType(t, validation.ValidationError{}, err, "expected error to be a ValidationError")
				require.ErrorContains(t, err, tt.errorContains)
			}
		})
	}
}

func expressionToString(exp Expression) string {
	switch exp := exp.(type) {
	case *IntegerLiteralExpression:
		return exp.Value.String()
	case *FloatingPointLiteralExpression:
		return exp.Value
	case *StringLiteralExpression:
		return fmt.Sprintf("«%s»", exp.Value)
	case *BinaryExpression:
		var op string
		switch exp.Operator {
		case BinaryOpAdd:
			op = "+"
		case BinaryOpSub:
			op = "-"
		case BinaryOpMul:
			op = "*"
		case BinaryOpDiv:
			op = "/"
		case BinaryOpPow:
			op = "**"
		default:
			panic(fmt.Sprintf("unexpected binary operator %d", exp.Operator))
		}
		return fmt.Sprintf("(%s %s %s)", op, expressionToString(exp.Left), expressionToString(exp.Right))
	case *UnaryExpression:
		if exp.Operator != UnaryOpNegate {
			panic(fmt.Sprintf("unexpected unary operator %d", exp.Operator))
		}
		return fmt.Sprintf("(- %s)", expressionToString(exp.Expression))
	case *TypeConversionExpression:
		return fmt.Sprintf("(as %s %s)", expressionToString(exp.Expression), TypeToShortSyntax(exp.Type, true))
	case *MemberAccessExpression:
		if exp.Target == nil {
			return exp.Member
		}
		return fmt.Sprintf("(. %s %s)", expressionToString(exp.Target), exp.Member)
	case *FunctionCallExpression:
		if len(exp.Arguments) == 0 {
			return fmt.Sprintf("(call %s)", exp.FunctionName)
		}

		args := make([]string, len(exp.Arguments))
		for i, arg := range exp.Arguments {
			args[i] = expressionToString(arg)
		}
		return fmt.Sprintf("(call %s %s)", exp.FunctionName, strings.Join(args, " "))
	case *SubscriptExpression:
		if len(exp.Arguments) == 0 {
			return fmt.Sprintf("(subscript %s)", expressionToString(exp.Target))
		}

		args := make([]string, len(exp.Arguments))
		for i, arg := range exp.Arguments {
			if arg.Label != "" {
				args[i] = fmt.Sprintf("%s:%s", arg.Label, expressionToString(arg.Value))
			} else {
				args[i] = expressionToString(arg.Value)
			}
		}

		return fmt.Sprintf("(subscript %s %s)", expressionToString(exp.Target), strings.Join(args, " "))

	default:
		return fmt.Sprintf("UNRECOGNIZED %T", exp)
	}
}
