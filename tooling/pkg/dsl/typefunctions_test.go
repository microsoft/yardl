// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) {
	base := `G<T, U>: !record
  fields:
    t: T
    u: U

S: string
X: `

	testCases := []struct {
		input    string
		expected string
	}{
		{input: "int", expected: "int32"},
		{input: "S"},
		{input: "G<int32, G<date, time>>"},
		{input: "!generic {name: G, args: [int32, [int32, float32]]}", expected: "G<int32, int32 | float32>"},
		{input: "int32?"},
		{input: "[null, int32, float32]", expected: "null | int32 | float32"},
		{input: "!union {anInt: int, aFloat: float}", expected: "anInt: int32 | aFloat: float32"},
		{input: "int32->string"},
		{input: "int32->(string?)", expected: "int32->string?"},
		{input: "(int32->string)?", expected: "(int32->string)?"},
		{input: "!map { keys: int32, values: [string, int32] }", expected: "int32->(string | int32)"},
		{input: "int32*"},
		{input: "int32*2"},
		{input: "(int32?)*2?", expected: "int32?*2?"},
		{input: "!vector {items: [int32, string]}", expected: "(int32 | string)*"},
		{input: "int32[]"},
		{input: "int32[]?"},
		{input: "int32[]?*?"},
		{input: "!array {items: [int32, string]}", expected: "(int32 | string)[]"},
		{input: "int32[,]"},
		{input: "int32[,,]"},
		{input: "int32[2, 3]"},
		{input: "int32[x:2, y:3]"},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			env, err := parseAndValidate(t, base+tC.input)
			require.NoError(t, err)
			x := env.SymbolTable["test.X"].(*NamedType)
			str := TypeToShortSyntax(x.Type, false)
			if tC.expected == "" {
				assert.Equal(t, tC.input, str)
			} else {
				assert.Equal(t, tC.expected, str)
			}
		})
	}

}
