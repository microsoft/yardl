// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagsAutoAssigned(t *testing.T) {
	src := `
X: !flags
  values:
    - a
    - b
    - c
`
	env, err := parseAndValidate(t, src)
	assert.Nil(t, err)
	e := env.Namespaces[0].TypeDefinitions[0].(*EnumDefinition)
	assert.Equal(t, int64(1), e.Values[0].IntegerValue.Int64())
	assert.Equal(t, int64(2), e.Values[1].IntegerValue.Int64())
	assert.Equal(t, int64(4), e.Values[2].IntegerValue.Int64())
}

func TestFlagsSomeAssigned(t *testing.T) {
	src := `
X: !flags
  values:
    a: 0
    b:
    c:
    d: 99
    e:
    f: -8
    g: -9
`
	env, err := parseAndValidate(t, src)
	assert.Nil(t, err)
	e := env.Namespaces[0].TypeDefinitions[0].(*EnumDefinition)
	assert.Equal(t, int64(0), e.Values[0].IntegerValue.Int64())
	assert.Equal(t, int64(1), e.Values[1].IntegerValue.Int64())
	assert.Equal(t, int64(2), e.Values[2].IntegerValue.Int64())
	assert.Equal(t, int64(99), e.Values[3].IntegerValue.Int64())
	assert.Equal(t, int64(128), e.Values[4].IntegerValue.Int64())
	assert.Equal(t, int64(-8), e.Values[5].IntegerValue.Int64())
	assert.Equal(t, int64(-9), e.Values[6].IntegerValue.Int64())
}

func TestFlagsAutoAssignedConflict(t *testing.T) {
	src := `
X: !flags
  values:
    a: 0
    b:
    c: 1
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "in flags 'X', the symbols [b c] have the same value of 1")
}

func TestFlagsCannotAutoAssignValueAfterNegative(t *testing.T) {
	src := `
X: !flags
  values:
    a: -1
    b:
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "flag value following a negative value must be explicitly specified")
}
