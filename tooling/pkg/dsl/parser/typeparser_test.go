// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeParsing_Valid(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "Foo", expected: `'Foo'`},
		{input: "Foo ", expected: `'Foo'`},
		{input: " Foo ", expected: `'Foo'`},

		{input: "Foo?", expected: `(Optional 'Foo')`},
		{input: "Foo ? ", expected: `(Optional 'Foo')`},

		{input: "Foo<int>", expected: `(Generic 'Foo' 'int')`},
		{input: "Foo < int > ", expected: `(Generic 'Foo' 'int')`},
		{input: "Foo<int?>", expected: `(Generic 'Foo' (Optional 'int'))`},
		{input: "Foo<int>?", expected: `(Optional (Generic 'Foo' 'int'))`},
		{input: "Foo<int,float>", expected: `(Generic 'Foo' 'int' 'float')`},
		{input: "Foo < int , float > ", expected: `(Generic 'Foo' 'int' 'float')`},
		{input: "Foo<Bar<int>>", expected: `(Generic 'Foo' (Generic 'Bar' 'int'))`},
		{input: "Foo<Bar<int>, Baz<long>>", expected: `(Generic 'Foo' (Generic 'Bar' 'int') (Generic 'Baz' 'long'))`},

		{input: "string->int", expected: `(Map 'string' 'int')`},
		{input: "string -> int", expected: `(Map 'string' 'int')`},
		{input: "string->int?", expected: `(Map 'string' (Optional 'int'))`},
		{input: "Foo<string>->Bar<int>?", expected: `(Map (Generic 'Foo' 'string') (Optional (Generic 'Bar' 'int')))`},
		{input: "Foo<string>->(Bar<int>)?", expected: `(Map (Generic 'Foo' 'string') (Optional (Generic 'Bar' 'int')))`},
		{input: "(Foo<string>->Bar<int>)?", expected: `(Optional (Map (Generic 'Foo' 'string') (Generic 'Bar' 'int')))`},
		{input: "Foo?->Bar<int>", expected: `(Map (Optional 'Foo') (Generic 'Bar' 'int'))`},

		{input: "Foo*", expected: `(Vector 'Foo')`},
		{input: "Foo * ", expected: `(Vector 'Foo')`},
		{input: "Foo*3", expected: `(Vector[3] 'Foo')`},
		{input: "Foo?*", expected: `(Vector (Optional 'Foo'))`},
		{input: "Foo*?", expected: `(Optional (Vector 'Foo'))`},
		{input: "Foo?*3?", expected: `(Optional (Vector[3] (Optional 'Foo')))`},

		{input: "Foo[]", expected: `(Array 'Foo')`},
		{input: "Foo[()]", expected: `(Array[[]] 'Foo')`},
		{input: "Foo[()]", expected: `(Array[[]] 'Foo')`},
		{input: "Foo[,]", expected: `(Array[[][]] 'Foo')`},
		{input: "Foo[, ,]", expected: `(Array[[][][]] 'Foo')`},
		{input: "Foo[x,y]", expected: `(Array[[x][y]] 'Foo')`},
		{input: "Foo[2,3]", expected: `(Array[[2][3]] 'Foo')`},
		{input: "Foo[x:2,y:3]", expected: `(Array[[x 2][y 3]] 'Foo')`},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ts, err := ParseType(tc.input)
			assert.Nil(t, err)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, ts.String())
		})
	}
}

func TestTypeParsing_Invalid(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// TODO: Inspect this carefully before finalizing the "Recursive Types" feature
		// {input: "Foo<", expected: `unexpected token "<EOF>" (expected Type ("," Type)* ">")`},
		{input: "Foo<", expected: `unexpected token "<EOF>" (expected <ident> ("." <ident>)* ("<" Type ("," Type)* ">")?)`},

		{input: "Foo?<int>", expected: `unexpected token "<"`},
		{input: "Foo<int>>", expected: `unexpected token ">"`},
		{input: "Foo<int>,bar>", expected: `unexpected token ","`},
		{input: "<int>", expected: `unexpected token "<"`},
		{input: "Foo<>", expected: `unexpected token ">"`},
		{input: "Foo<int,>", expected: `unexpected token "," (expected ">")`},

		// TODO: Inspect this carefully before finalizing the "Recursive Types" feature
		// {input: "string->", expected: `unexpected token "<EOF>" (expected Type)`},
		{input: "string->", expected: `unexpected token "<EOF>" (expected <ident> ("." <ident>)* ("<" Type ("," Type)* ">")?)`},

		{input: "int[", expected: `unexpected token "<EOF>" (expected "]")`},
		{input: "int[/]", expected: `unexpected token "/" (expected "]")`},
		{input: "int[x:]", expected: `unexpected token "]" (expected integer)`},
		{input: "int[(x:2]", expected: `unexpected token "]" (expected ")")`},
		{input: "int[((x:2)]", expected: `unexpected token "]" (expected ")")`},
		{input: "int[x:4987439128739182743918274]", expected: `integer out of range`},
		{input: "int[4987439128739182743918274]", expected: `integer out of range`},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ParseType(tc.input)
			t.Log(err)
			assert.ErrorContains(t, err, tc.expected)
		})
	}
}

func TestTypeParsing2(t *testing.T) {
	ParseType("Foo<A>")
}
