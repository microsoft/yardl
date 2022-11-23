// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayDimensionLengthCannotBeNegative(t *testing.T) {
	src := `
X: !record
  fields:
    f: !array
      items: int
      dimensions: [-1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "length cannot be negative")
}

func TestArrayDimensionsLengthsPartiallySpecified(t *testing.T) {
	src := `
X: !record
  fields:
    f: !array
      items: int
      dimensions:
        a:
        b: 3`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "lengths must either be specified on all dimensions or none of them")
}

func TestArrayDimensionNameInvalid(t *testing.T) {
	src := `
X: !record
  fields:
    f: !array
      items: int
      dimensions:
        "9": 3`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "dimension name '9' must match the format")
}

func TestArrayDimensionNameDuplicated(t *testing.T) {
	src := `
X: !record
  fields:
    f: !array
      items: int
      dimensions:
        x: 3
        x: 3`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "a dimension with the name 'x' is already defined on the array")
}
