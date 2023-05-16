// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input         string
		errorContains string
		parsed        string
	}{
		{`1`, "", "1"},
		{`-1`, "", "-1"},
		{`0x10`, "", "0x10"},
		{`0XAAb`, "", "0XAAb"},
		{`-0x10`, "", "-0x10"},
		{`0x10478349823749782398327643786348762387326`, "", "0x10478349823749782398327643786348762387326"},
		{`-0x10478349823749782398327643786348762387326`, "", "-0x10478349823749782398327643786348762387326"},
		{`"s"`, "", "«s»"},
		{`"s\"s"`, "", `«s"s»`},
		{`"s\\s"`, "", `«s\s»`},
		{"'s'", "", `«s»`},
		{`'s\'s'`, "", `«s's»`},
		{`'s\\s'`, "", `«s\s»`},
		{`'s"s'`, "", `«s"s»`},
		{`"s's"`, "", `«s's»`},
		{" foo ( 1 , 2 ) ", "", "(call foo 1, 2)"},
		{"foo()", "", "(call foo)"},
		{"foo(bar())", "", "(call foo (call bar))"},
		{"foo(bar(),1)", "", "(call foo (call bar), 1)"},
		{"foo(bar(),)", `unexpected token "," (expected ")")`, ""},
		{"foo[]", "", "(part foo (index))"},
		{"foo[1]", "", "(part foo (index 1))"},
		{"foo[1, 2]", "", "(part foo (index 1, 2))"},
		{"foo[x:1, y:2]", "", "(part foo (index x:1, y:2))"},
		{"foo[1,]", `unexpected token "," (expected "]")`, ""},
		{"foo[,1]", `unexpected token "," (expected Expression)`, ""},
		{"foo[1][2]", "", "(part foo (index 1)(index 2))"},
		{"foo[1].bar[2]", "", "(part foo (index 1)).(part bar (index 2))"},
		{"foo[bar()]", "", "(part foo (index (call bar)))"},
		{"[1]", `unexpected token "["`, ""},

		// This is not a great error message, since
		// the real issue is that the string is not terminated.
		// Looks like we would need to write a custom lexer to
		// handle this better
		{`"a`, `unexpected token "`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			exp, err := ParseExpression(tt.input)
			if tt.errorContains == "" {
				if err != nil {
					assert.Nil(t, err, "unexpected error: %v", err)
				} else {
					assert.Equal(t, tt.parsed, (*exp).String())
				}
			} else {
				err, ok := err.(ParseError)
				assert.True(t, ok, "expected error to be a ParseError")
				if ok {
					assert.Contains(t, err.Message(), tt.errorContains)
				}
			}
		})
	}
}
