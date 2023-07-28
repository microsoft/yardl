// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"os"
	"path"
	"strings"
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
	assert.ErrorContains(t, err, "a union type must have at least one option")
}

func TestNullMustBeFirstUnionElement(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, null]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "if null is specified in a union type, it must be the first option")
}

func TestUnionElementsMustBeDistinct(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, int]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinctWithGenerics(t *testing.T) {
	src := `
X: !record
  fields:
    f: !union
      g1: GenericRecord<int>
      g2: GenericRecord<int>
GenericRecord<T>: !record`

	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsAreDistinctWithGenerics(t *testing.T) {
	src := `
X: !record
  fields:
    f: !union
      g1: GenericRecord<float>
      g2: GenericRecord<double>

GenericRecord<T>: !record
  fields:
    t: T`

	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestUnionElementsMustBeDistinct_SameUnrecognizedType(t *testing.T) {
	src := `
Bar: int
X: !record
  fields:
    f: [Bar, Bar]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_DifferentUnrecognizedType(t *testing.T) {
	src := `
X: !record
  fields:
    f: [Foo, Bar]`
	_, err := parseAndValidate(t, src)
	assert.NotNil(t, err)
	assert.NotContains(t, err.Error(), "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_MultipleNulls(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null, null, null]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_Complex(t *testing.T) {
	src := `
X: !record
  fields:
    f: !union
      o1: !vector
        items: [int, float]
        length: 10
      o2: !vector
        items: [int, float]
        length: 10`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_Nested(t *testing.T) {
	src := `
X: !record
  fields:
    f:
      - int
      - [ float, float ]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_AliasedType(t *testing.T) {
	src := `
MyIntType: uint64
MyRecord: !record
  fields:
    one: [uint64, MyIntType]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_AliasedWithinVector(t *testing.T) {
	src := `
MyIntType: uint64
MyRecord: !record
  fields:
    f:
      - !vector
        items: MyIntType
      - !vector
        items: uint64`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "redundant union type cases")
}

func TestUnionElementsMustBeDistinct_DifferentGenericArgs(t *testing.T) {
	src := `
MyIntType: uint64
Image<T>: !array
  items: T
MyRecord: !record
  fields:
    f: !union
      di: Image<double>
      df: Image<float>`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestUnionElementsMustBeDistinct_GenericUnionAlias_AllUnique(t *testing.T) {
	src := `
MyUnionType<T, U>: [T, U]
MyRecord: !record
  fields:
    f: MyUnionType<int, float>`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestUnionElementsMustBeDistinct_GenericUnionAlias_NotUnique_MultipleTypeArgs(t *testing.T) {
	src := `
MyUnionType<T, U>: [T, U]
MyRecord: !record
  fields:
    f: MyUnionType<int, int>`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.Regexp(t, ".yaml:2:21: redundant union type cases resulting from the type arguments given at .*.yaml:5:20 and .*.yaml:5:25", err.Error())
}

func TestUnionElementsMustBeDistinct_GenericUnionAlias_NotUnique_SingleTypeArg(t *testing.T) {
	src := `
MyUnionType<T, U>: [T, int]
MyRecord: !record
  fields:
    f: MyUnionType<int, float>`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.Regexp(t, ".yaml:2:24: redundant union type cases resulting from the type argument given at .*.yaml:5:20$", err.Error())
}

func TestUnionElementsMustBeDistinct_GenericUnionAliasChain_SingleTypeArg(t *testing.T) {
	src := `
Rec<T>: !record
  fields:
    f: [T, int]
Alias1<T>: Rec<T>
Alias2: Alias1<int>`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.Regexp(t, ".yaml:4:12: redundant union type cases resulting from the type argument given at .*.yaml:6:16$", err.Error())
}

func TestUnionElementsMustBeDistinct_GenericUnionAliasChain_ErrorsNotDuplicated(t *testing.T) {
	src := `
Rec<T>: !record
  fields:
    f: [int, int]
Alias1<T>: Rec<T>
Alias2: Alias1<int>`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.Regexp(t, ".yaml:4:9: redundant union type cases$", err.Error())
	assert.Equal(t, 1, len(strings.Split(err.Error(), "\n")))
}

func TestUnionElementsMustBeDistinct_SizeAndUnit64Direct(t *testing.T) {
	src := `
T: [size, uint64]`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "redundant union type cases (uint64 and size are equivalent)")
}

func TestUnionElementsMustBeDistinct_SizeAndUnit64Aliased(t *testing.T) {
	src := `
MySize: size
MyUint64: uint64
T: [MySize, MyUint64]`
	_, err := parseAndValidate(t, src)
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "redundant union type cases (uint64 and size are equivalent)")
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

func TestUnionTags_Simple(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null, int, string]`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)

	f := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[0]
	require.Equal(t, "int32", f.Type.(*GeneralizedType).Cases[1].Tag)
	require.Equal(t, "string", f.Type.(*GeneralizedType).Cases[2].Tag)
}

func TestUnionTags_SimpleVectorsAndArrays(t *testing.T) {
	src := `
X: !record
  fields:
    f: !union
      intVector: !vector {items: int}
      stringVector: !vector {items: string}
    g: !union
      intArray: !array {items: int}
      stringArray: !array {items: string}
    h: !union
      stringToInt: string->int
      stringToString: string->string`
	env, err := parseAndValidate(t, src)
	require.Nil(t, err)

	f := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[0]
	require.Equal(t, "intVector", f.Type.(*GeneralizedType).Cases[0].Tag)
	require.Equal(t, "stringVector", f.Type.(*GeneralizedType).Cases[1].Tag)
	g := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[1]
	require.Equal(t, "intArray", g.Type.(*GeneralizedType).Cases[0].Tag)
	require.Equal(t, "stringArray", g.Type.(*GeneralizedType).Cases[1].Tag)
	h := env.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).Fields[2]
	require.Equal(t, "stringToInt", h.Type.(*GeneralizedType).Cases[0].Tag)
	require.Equal(t, "stringToString", h.Type.(*GeneralizedType).Cases[1].Tag)
}

func TestUnionsCannotContainUnions(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int, [null, int]]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "unions may not immediately contain other unions")
}

func TestUnionsWithValidTagSyntax(t *testing.T) {
	src := `
U: !union
  s: string
  i: int`
	_, err := parseAndValidate(t, src)
	assert.NoError(t, err)
}

func TestUnionsHaveDistinctTags(t *testing.T) {
	src := `
U: !union
  s: string
  s: int`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "all union cases must have distinct tags")
}

func TestUnionsTagsMustBeValidIdentifiers(t *testing.T) {
	src := `
U: [int, int*2]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the type 'int32*2' cannot be used as a tag for the union case. Explicit tags can be given using the `!union` syntax (e.g. `!union { myTag: \"int32*2\", ... }`) or the type can be aliased for the type case (e.g. `MyTypeAlias = int32*2\nMyUnion = [..., MyTypeAlias, ...]`)")
}

func TestUnionsTagsMustBeValidIdentifiersWithOpenGeneric(t *testing.T) {
	src := `
VecOrList<T>: [T*, "T[]"]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "the type 'T*' cannot be used as a tag for the union case. An explicit tag can be given using the `!union` syntax (e.g. `!union { myTag: \"T*\", ... }`)")
}

func TestUnionsWithComplexTypesRequireTags(t *testing.T) {
	src := `
U: !union
  s*: string`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "union tag 's*' must be camelCased matching the format ^[a-z][a-zA-Z0-9]{0,63}$")
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
