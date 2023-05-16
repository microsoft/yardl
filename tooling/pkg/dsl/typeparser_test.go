// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

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

		// {input: "Foo[]", expected: ``},
		// {input: "Foo[]?", expected: ``},
		// {input: "Foo?[]", expected: ``},
		// {input: "Foo?[]?", expected: ``},
		// {input: "Foo[*]", expected: ``},
		// {input: "Foo[]", expected: ``},
		// {input: "Foo[,]", expected: ``},
		// {input: "Foo[x,y]", expected: ``},
		// {input: "Foo[2,3]", expected: ``},
		// {input: "Foo[x:2,y:3]", expected: ``},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ts, err := parseSimpleType2(tc.input)
			assert.Nil(t, err)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, ts.String())
		})
	}
}

// func TestTypeParsing_Invalid(t *testing.T) {
// 	testCases := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{input: "Foo<", expected: `missing '>'`},
// 		{input: "Foo?<int>", expected: `'?' at position 4 must appear after generic type arguments`},
// 		{input: "Foo??", expected: `unexpected trailing '?' in type string`},
// 		{input: "Foo<int>>", expected: `unexpected trailing '>' in type string`},
// 		{input: "Foo<int>,bar>", expected: `unexpected trailing ',bar>' in type string`},
// 		{input: "<int>", expected: `the type name cannot be empty`},
// 		{input: "Foo<>", expected: `the type parameter name cannot be empty at position 5`},
// 		{input: "Foo<int,>", expected: `the type parameter name cannot be empty at position 9`},
// 		{input: "string->", expected: `missing type name after '->'`},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.input, func(t *testing.T) {
// 			_, err := parseSimpleTypeString(tc.input)
// 			assert.ErrorContains(t, err, tc.expected)
// 		})
// 	}
// }

func TestTypeParsing2(t *testing.T) {
	parseSimpleType2("Foo<A>")
}
