// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func parsePair(t *testing.T, v0Src, v1Src string) (*Environment, *Environment) {
	v0, err := parseAndValidate(t, v0Src)
	assert.Nil(t, err)

	v1, err := parseAndValidate(t, v1Src)
	assert.Nil(t, err)

	return v0, v1
}

// func TestStrictCompareNamespaces(t *testing.T) {
// 	v0, err := parseAndValidate(t, `
// A: !record
//   fields:
//     a: int
// B: !record
//   fields:
//     b: string
// C: !protocol
//   sequence:
//     c: int
// D: !protocol
//   sequence:
//     d: string
// }`)
// 	assert.Nil(t, err)

// 	t.Run("remove record", func(t *testing.T) {
// 		v1, err := parseAndValidate(t, `

// }

func TestStrictCompareProtocols(t *testing.T) {
	v0, err := parseAndValidate(t, `
# This is a comment on a protocol
X: !protocol
  sequence:
    # comment on step
    a: int
    b: float
    c: string`)
	assert.Nil(t, err)

	t.Run("identical", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float
    c: string`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.Nil(t, err)
	})

	t.Run("removed comments", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float
    c: string`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.Nil(t, err)
	})

	t.Run("removed step", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.NotNil(t, err)
	})

	t.Run("added step", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float
    c: string
    d: bool`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.NotNil(t, err)
	})

	t.Run("renamed step", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float
    s: string`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.NotNil(t, err)
	})

	t.Run("changed step type", func(t *testing.T) {
		v1, err := parseAndValidate(t, `
X: !protocol
  sequence:
    a: int
    b: float
    c: bool`)
		assert.Nil(t, err)

		_, err = ensureNoChanges(v1, v0, nil)
		assert.NotNil(t, err)
	})
}
