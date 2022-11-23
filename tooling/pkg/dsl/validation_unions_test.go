// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullCannotBeSingleUnionElement(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "null cannot be the only option in a union type")
}

func TestFieldTypeCannotBeNull(t *testing.T) {
	src := `
X: !record
  fields:
    f: null`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "a field or protocol step cannot be null")
}

func TestFieldTypeCannotBeEmptyUnion(t *testing.T) {
	src := `
X: !record
  fields:
    f: []`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "a field cannot be a union type with no options")
}

func TestNullMustBeFirstUnionElement(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, null]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "if null is specified in a union type, it must be the first option")
}

func TestUnionElementsMustBeDistict(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, int]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestUnionElementsMustBeDistictWithGenerics(t *testing.T) {
	src := `
X: !record
  fields:
    f: [GenericRecord<int>, GenericRecord<int>]
GenericRecord<T>: !record`

	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestUnionElementsAreDistictWithGenerics(t *testing.T) {
	src := `
X: !record
  fields:
    f: [GenericRecord<float>, GenericRecord<double>]
GenericRecord<T>: !record`

	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestUnionElementsMustBeDistict_SameUnrecognizedType(t *testing.T) {
	src := `
X: !record
  fields:
    f: [Bar, Bar]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestUnionElementsMustBeDistict_DifferentUnrecognizedType(t *testing.T) {
	src := `
X: !record
  fields:
    f: [Foo, Bar]`
	_, err := parseAndValidate(t, src)
	assert.NotNil(t, err)
	assert.NotContains(t, err.Error(), "all type cases in a union must be distinct")
}

func TestUnionElementsMustBeDistict_MultipleNulls(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null, null, null]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestUnionElementsMustBeDistict_Complex(t *testing.T) {
	src := `
X: !record
  fields:
    f:
      - !vector
        items: [int, float]
        length: 10
      - !vector
        items: [int, float]
        length: 10`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestUnionElementsMustBeDistict_Nested(t *testing.T) {
	src := `
X: !record
  fields:
    f:
      - int
      - [ float, float ]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all type cases in a union must be distinct")
}

func TestVectorLengthCannotBeNegative(t *testing.T) {
	src := `
X: !record
  fields:
    f: !vector
      items: int
      length: -1`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "length cannot be negative")
}

func TestUnionLabels_Simple(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null, int, string]`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)

	f := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[0]
	require.Equal(t, "int32", f.Type.(*GeneralizedType).Cases[1].Label)
	require.Equal(t, "string", f.Type.(*GeneralizedType).Cases[2].Label)
}

func TestUnionLabels_SimpleVectorsAndArrays(t *testing.T) {
	src := `
X: !record
  fields:
    f: [!vector {items: int}, !vector {items: string}]
    g: [!array {items: int}, !array {items: string}]`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)

	f := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[0]
	require.Equal(t, "int32Vector", f.Type.(*GeneralizedType).Cases[0].Label)
	require.Equal(t, "stringVector", f.Type.(*GeneralizedType).Cases[1].Label)
	g := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[1]
	require.Equal(t, "int32Array", g.Type.(*GeneralizedType).Cases[0].Label)
	require.Equal(t, "stringArray", g.Type.(*GeneralizedType).Cases[1].Label)
}

func TestUnionLabels_VectorsAndArraysWithDisambiguation(t *testing.T) {
	src := `
X: !record
  fields:
    f: [!vector {items: string}, !vector {items: string, length: 10}]
    g: [!array {items: string}, !array {items: string, dimensions: 2}]
    h: [!array {items: string, dimensions: [2,3]}, !array {items: string, dimensions: 2}]`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)

	f := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[0]
	require.Equal(t, "stringVector", f.Type.(*GeneralizedType).Cases[0].Label)
	require.Equal(t, "stringVector[10]", f.Type.(*GeneralizedType).Cases[1].Label)
	g := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[1]
	require.Equal(t, "stringArray", g.Type.(*GeneralizedType).Cases[0].Label)
	require.Equal(t, "stringArray[,,]", g.Type.(*GeneralizedType).Cases[1].Label)
	h := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[2]
	require.Equal(t, "stringArray[2,3]", h.Type.(*GeneralizedType).Cases[0].Label)
	require.Equal(t, "stringArray[,,]", h.Type.(*GeneralizedType).Cases[1].Label)
}

func TestUnionsCannotContainUnions(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, [null, int]]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "unions may not immediately contain other unions")
}

func parseAndValidate(t *testing.T, src string) (*Environment, error) {
	d := t.TempDir()
	os.WriteFile(path.Join(d, "t.yaml"), []byte(src), 0644)
	ns, err := ParseYamlInDir(d, "test")
	if err != nil {
		return nil, err
	}

	return Validate([]*Namespace{ns})
}
