// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeNamesNotUnique(t *testing.T) {
	src := `
X: !record
  fields:
    unused: int
X: !enum
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the name 'X' is already defined in ")
}

func TestTypeNamesNotUniqueOneWithGenericParameters(t *testing.T) {
	src := `
X<T>: !record
  fields:
    unused: T
X<T1, T2>: !record
  fields:
    unused: [T1, T2]
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the name 'X' is already defined in ")
}

func TestCannotReferenceProtocol(t *testing.T) {
	src := `
MyProtocol: !protocol
  sequence:
    s: string

Rec: !record
  fields:
    f: MyProtocol`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "cannot reference a protocol")
}

func TestEnumsCannotBeGeneric(t *testing.T) {
	src := `
Abc<T>: !enum
  values: [A, B, C]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "'Abc' cannot have generic type parameters")
}

func TestProtocolsCannotBeGeneric(t *testing.T) {
	src := `
MyProtocol<T>: !protocol
  sequence:
    s: T`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "'MyProtocol' cannot have generic type parameters")
}

func TestNoTypeArgsGiven(t *testing.T) {
	src := `
Rec1<T>: !record
  fields:
    f: T

Rec2: !record
  fields:
    f: Rec1
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "'Rec1' was given 0 type argument(s) but has 1 type parameter(s)")
}

func TestTypeParamNotResolved(t *testing.T) {
	src := `
Rec1<T>: !record
  fields:
    f: T

Rec2: !record
  fields:
    f: Rec1<Foo>`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the type 'Foo' is not recognized")
}

func TestTypeParameterUnused(t *testing.T) {
	src := `
MyUnion<T, U>: [T, int]
`

	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "generic type parameter 'U' is not used")
}
