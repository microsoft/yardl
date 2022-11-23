// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumBaseTypeValid(t *testing.T) {
	src := `
X: !enum
  base: int
  values:
    - a
`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestEnumSymbolNameInvalid(t *testing.T) {
	src := `
X: !enum
  base: int
  values:
    - FOOBAR
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "in enum 'X', the symbol name 'FOOBAR' must be camelCased matching the format")
}

func TestValidEnumBaseTypes(t *testing.T) {
	testCases := []string{
		`X: !enum {base: int8}`,
		`X: !enum {base: uint8}`,
		`X: !enum {base: int16}`,
		`X: !enum {base: uint16}`,
		`X: !enum {base: int32}`,
		`X: !enum {base: uint32}`,
		`X: !enum {base: int64}`,
		`X: !enum {base: uint64}`,
		`{MyInt: int, X: !enum {base: MyInt}}`,
		`{MyGeneric<T>: T, X: !enum {base: MyGeneric<int>}}`,
	}
	for _, src := range testCases {
		t.Run(src, func(t *testing.T) {
			_, err := parseAndValidate(t, src)
			assert.Nil(t, err)
		})
	}
}

func TestInvalidEnumBaseTypes(t *testing.T) {
	testCases := []string{
		`X: !enum {base: float32}`,
		`X: !enum {base: float64}`,
		`X: !enum {base: "int?"}`,
		`X: !enum {base: !vector {items: int}}`,
		`X: !enum {base: [int, string]}`,
		`{MyString: string, X: !enum {base: MyString}}`,
	}
	for _, src := range testCases {
		t.Run(src, func(t *testing.T) {
			_, err := parseAndValidate(t, src)
			assert.ErrorContains(t, err, "in enum 'X', the base type must be an integer type")
		})
	}
}

func TestUnrecognizedEnumBaseTypes(t *testing.T) {
	testCases := []string{
		`X: !enum {base: missing}`,
	}
	for _, src := range testCases {
		t.Run(src, func(t *testing.T) {
			_, err := parseAndValidate(t, src)
			assert.ErrorContains(t, err, "the type 'missing' is not recognized")
		})
	}
}

func TestEnumValuesWithinRange(t *testing.T) {
	testCases := []string{
		`X: !enum {base: int8, values: {a: -0x80, b: 0x7f}}`,
		`X: !enum {base: uint8, values: {a: 0, b: 0xff}}`,
		`X: !enum {base: int16, values: {a: -0x8000, b: 0x7fff}}`,
		`X: !enum {base: uint16, values: {a: 0, b: 0xffff}}`,
		`X: !enum {base: int32, values: {a: -0x80000000, b: 0x7fffffff}}`,
		`X: !enum {base: uint32, values: {a: 0, b: 0xffffffff}}`,
		`X: !enum {base: int64, values: {a: -0x8000000000000000, b: 0x7fffffffffffffff}}`,
		`X: !enum {base: uint64, values: {a: 0, b: 0xffffffffffffffff}}`,
	}
	for _, src := range testCases {
		t.Run(src, func(t *testing.T) {
			_, err := parseAndValidate(t, src)
			assert.Nil(t, err)
		})
	}
}

func TestEnumValuesTooLarge(t *testing.T) {
	testCases := []string{
		`X: !enum {base: int8, values: {a: -0x81}}`,
		`X: !enum {base: int8, values: {b: 0x80}}`,
		`X: !enum {base: uint8, values: {b: 0x100}}`,
		`X: !enum {base: int16, values: {a: -0x8001}}`,
		`X: !enum {base: int16, values: {b: 0x8000}}`,
		`X: !enum {base: uint16, values: {b: 0x10000}}`,
		`X: !enum {base: int32, values: {a: -0x80000001}}`,
		`X: !enum {base: int32, values: {b: 0x80000000}}`,
		`X: !enum {base: uint32, values: {b: 0x100000000}}`,
		`X: !enum {base: int64, values: {a: -0x8000000000000001}}`,
		`X: !enum {base: int64, values: {b: 0x8000000000000000}}`,
		`X: !enum {base: uint64, values: {b: 0x10000000000000000}}`,
	}
	for _, src := range testCases {
		t.Run(src, func(t *testing.T) {
			_, err := parseAndValidate(t, src)
			assert.ErrorContains(t, err, "out of range for the base type")
		})
	}
}
