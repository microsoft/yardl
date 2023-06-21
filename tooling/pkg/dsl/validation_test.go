// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordFieldNameInvalid(t *testing.T) {
	src := `
Rec: !record
  fields:
    __: int
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "field name '__' must be camelCased matching the format")
}

func TestRecordComputedFieldNameInvalid(t *testing.T) {
	src := `
Rec: !record
  computedFields:
    A: 1
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "field name 'A' must be camelCased matching the format")
}

func TestRecordProtocolStepNameInvalid(t *testing.T) {
	src := `
Proto: !protocol
  sequence:
    _b: int
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "protocol step name '_b' must be camelCased matching the format")
}

func TestInvalidProtocolName(t *testing.T) {
	src := `
abc: !protocol
  sequence:
    a: int
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "type name 'abc' must be PascalCased matching the format")
}

func TestInvalidGenericTypeParameterName(t *testing.T) {
	src := `
Abc<T_A>: !record
  fields:
    a: t
`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "generic type parameter name 'T_A' must be PascalCased matching the format")
}
