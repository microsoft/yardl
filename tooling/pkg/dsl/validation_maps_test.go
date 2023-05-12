// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapNonPrimitiveKey(t *testing.T) {
	src := `
X: !map
  keys: !vector
    items: string
  values: int`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "map key type must be a primitive scalar type")
}

func TestMapAliasedPrimitiveKey(t *testing.T) {
	src := `
K: string
X: !map
  keys: K
  values: int`
	_, err := parseAndValidate(t, src)
	assert.NoError(t, err)
}

func TestMapUnresolvedKeyType(t *testing.T) {
	src := `
X: !map
  keys: Missing
  values: int`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the type 'Missing' is not recognized")
}
