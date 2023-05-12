// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeParsing_Valid(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "Foo", expected: `{"name":"Foo"}`},
		{input: "Foo ", expected: `{"name":"Foo"}`},
		{input: " Foo ", expected: `{"name":"Foo","positionOffset":1}`},
		{input: "Foo?", expected: `{"name":"Foo","optional":true}`},
		{input: "Foo ? ", expected: `{"name":"Foo","optional":true}`},
		{input: "Foo<int>", expected: `{"name":"Foo","args":[{"name":"int","positionOffset":4}]}`},
		{input: "Foo< int >", expected: `{"name":"Foo","args":[{"name":"int","positionOffset":5}]}`},
		{input: "Foo<int?>", expected: `{"name":"Foo","args":[{"name":"int","optional":true,"positionOffset":4}]}`},
		{input: "Foo<int>?", expected: `{"name":"Foo","args":[{"name":"int","positionOffset":4}],"optional":true}`},
		{input: "Foo<int,float>", expected: `{"name":"Foo","args":[{"name":"int","positionOffset":4},{"name":"float","positionOffset":8}]}`},
		{input: " Foo < int , float > ", expected: `{"name":"Foo","args":[{"name":"int","positionOffset":7},{"name":"float","positionOffset":13}],"positionOffset":1}`},
		{input: "Foo<Bar<int>>", expected: `{"name":"Foo","args":[{"name":"Bar","args":[{"name":"int","positionOffset":8}],"positionOffset":4}]}`},
		{input: "Foo<Bar<int>,Baz<long>>", expected: `{"name":"Foo","args":[{"name":"Bar","args":[{"name":"int","positionOffset":8}],"positionOffset":4},{"name":"Baz","args":[{"name":"long","positionOffset":17}],"positionOffset":13}]}`},
		{input: "string->int", expected: `{"name":"int","mapKey":{"name":"string"},"positionOffset":8}`},
		{input: "string -> int", expected: `{"name":"int","mapKey":{"name":"string"},"positionOffset":10}`},
		{input: "string->Foo<int>?", expected: `{"name":"Foo","mapKey":{"name":"string"},"args":[{"name":"int","positionOffset":12}],"optional":true,"positionOffset":8}`},
		{input: "Foo<string>->Bar<int>?", expected: `{"name":"Bar","mapKey":{"name":"Foo","args":[{"name":"string","positionOffset":4}]},"args":[{"name":"int","positionOffset":17}],"optional":true,"positionOffset":13}`},
		{input: "Foo?->Bar<int>", expected: `{"name":"Bar","mapKey":{"name":"Foo","optional":true},"args":[{"name":"int","positionOffset":10}],"positionOffset":6}`},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ts, err := parseSimpleTypeString(tc.input)
			assert.Nil(t, err)
			bytes, err := json.Marshal(ts)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, string(bytes))
		})
	}
}

func TestTypeParsing_Invalid(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "Foo<", expected: `missing '>'`},
		{input: "Foo?<int>", expected: `'?' at position 4 must appear after generic type arguments`},
		{input: "Foo??", expected: `unexpected trailing '?' in type string`},
		{input: "Foo<int>>", expected: `unexpected trailing '>' in type string`},
		{input: "Foo<int>,bar>", expected: `unexpected trailing ',bar>' in type string`},
		{input: "<int>", expected: `the type name cannot be empty`},
		{input: "Foo<>", expected: `the type parameter name cannot be empty at position 5`},
		{input: "Foo<int,>", expected: `the type parameter name cannot be empty at position 9`},
		{input: "string->", expected: `missing type name after '->'`},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			_, err := parseSimpleTypeString(tc.input)
			assert.ErrorContains(t, err, tc.expected)
		})
	}
}
