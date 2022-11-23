// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicTologigicalSort(t *testing.T) {
	src := `
R: !record
  fields:
    f: E
E: !enum
`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)
	require.Equal(t, "E", env.Namespaces[0].TypeDefinitions[0].GetDefinitionMeta().Name)
	require.Equal(t, "R", env.Namespaces[0].TypeDefinitions[1].GetDefinitionMeta().Name)
}

func TestBasicCycleDection(t *testing.T) {
	src := `
R1: !record
  fields:
    f1: R2
R2: !record
  fields:
    f2: R1
`
	_, err := parseAndValidate(t, src)
	require.ErrorContains(t, err, "there is a reference cycle, which is not supported, within namespace 'test': Record 'R1' -> Field 'f1' -> Record 'R2' -> Field 'f2' -> Record 'R1'")
}

func TestCycleDectionThroughGenericTypeParams(t *testing.T) {
	src := `
R1: !record
  fields:
    f1: R2<R1>
R2<T>: !record
`
	_, err := parseAndValidate(t, src)
	require.ErrorContains(t, err, "there is a reference cycle, which is not supported, within namespace 'test': Record 'R1' -> Field 'f1' -> Record 'R1'")
}

func TestDirectCycleDectionThroughGenericTypeParamsOnNamedRecords(t *testing.T) {
	src := `
R1<T>: !record
  fields:
    f1: R1<T>
---
`
	_, err := parseAndValidate(t, src)
	require.ErrorContains(t, err, "there is a reference cycle, which is not supported, within namespace 'test': Record 'R1' -> Field 'f1' -> Record 'R1'")
}

func TestDirectCycleDectionThroughGenericTypeParamsOnNamedArray(t *testing.T) {
	src := `
Image<T>: !array
  items: Image<T>
---
`
	_, err := parseAndValidate(t, src)
	require.ErrorContains(t, err, "there is a reference cycle, which is not supported, within namespace 'test': Array 'Image' -> Array 'Image'")
}
