// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateComputedFieldName(t *testing.T) {
	src := `
X: !record
  computedFields:
    f: 1
    f: 2`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `a field or computed field with the name 'f' is already defined on the record 'X'`)
}

func TestUnboundField(t *testing.T) {
	src := `
X: !record
  computedFields:
    f: missingField`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "record 'X' does not have a field or computed field named 'missingField'")
}

func TestRecursiveField(t *testing.T) {
	src := `
X: !record
  computedFields:
    f: f`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "cycle detected in computed fields: f -> f")
}

func TestCycle(t *testing.T) {
	src := `
X: !record
  computedFields:
    a: b
    b: c
    c: a`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "cycle detected in computed fields: a -> b -> c -> a")
}

func TestNoCycle(t *testing.T) {
	src := `
X: !record
  computedFields:
    a: b
    b: c
    c: 9`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestIntTooLarge(t *testing.T) {
	src := `
X: !record
  computedFields:
    a: 0xFFFFFFFFFFFFFFFF1`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "literal is too large")
}

func TestNegativeIntTooLarge(t *testing.T) {
	src := `
X: !record
  computedFields:
    a: -0x8000000000000001`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, "literal is too large")
}

func TestChainedRecords(t *testing.T) {
	src := `
X: !record
  fields:
    y: Y
  computedFields:
    yInt: y.anInt
Y: !record
  fields:
    anInt: int`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestChainedExpressionTargetNotResolved(t *testing.T) {
	src := `
X: !record
  computedFields:
    yInt: y.anInt`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `record 'X' does not have a field or computed field named 'y'`)
}

func TestChainedExpressionTargetNotARecord(t *testing.T) {
	src := `
X: !record
  fields:
    y: int
  computedFields:
    yInt: y.anInt`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `member access target must be a !record type`)
}

func TestIndexVector(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
  computedFields:
    y1: y[1]`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestIndexVectorWithString(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
  computedFields:
    y1: y["1"]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `index argument must be an integral type`)
}

func TestIndexVectorWithNoArgs(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
  computedFields:
    y1: y[]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `vector index must have exactly one argument`)
}

func TestIndexVectorWithTwoArgs(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
  computedFields:
    y1: y[0,1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `vector index must have exactly one argument`)
}

func TestIndexVectorWithUnresolvable(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
  computedFields:
    y1: y[sin(0)]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `unknown function 'sin'`)
}

func TestIndexVectorFixedVectorTooLargeLiteral(t *testing.T) {
	src := `
X: !record
  fields:
    y: !vector
      items: int
      length: 10
  computedFields:
    y1: y[10]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `index argument (10) is too large for the vector of length 10`)
}

func TestIndexOnUnresolved(t *testing.T) {
	src := `
R: !record
X: !record
  fields:
    r: R
  computedFields:
    r0: r[0]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `index target must be a vector or array`)
}

func TestIndexVectorOfVectors(t *testing.T) {
	src := `
X: !record
  fields:
    v: !vector { items: !vector { items: int } }
  computedFields:
    v01: v[0][1]`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestIndexArray(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    y1: y[1, 2]
    y2: y[x:1, y:2]`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestIndexArrayMixingLabels(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    y2: y[2, z:1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array index cannot mix labeled and unlabeled arguments`)
}

func TestIndexArrayDimensionSpecifiedTwice(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    y2: y[x:1, x:1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array index has multiple arguments for dimension 'x'`)
}

func TestIndexArrayDimensionsOutOfOrder(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    y2: y[y:1, x:1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array index has arguments must be specified in order`)
}

func TestIndexArrayInvalidDimensionName(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    y2: y[x:2, z:1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the array has no dimension named 'z'`)
}

func TestIndexArrayTooFew(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions: [1, 2]
  computedFields:
    y1: y[0]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array index must provide arguments for all 2 dimensions`)
}

func TestIndexArrayTooManyDimensionsGiven(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions: [1, 2]
  computedFields:
    y1: y[0,1,2]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array index has more arguments than dimensions`)
}

func TestIndexFixedArrayDimensionTooLarge(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x: 10
        y: 20
  computedFields:
    y1: y[10,1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `index argument (10) is too large for array dimension '0' of length 10`)
}

func TestIndexFixedArrayDimensionTooLargeNamedDimension(t *testing.T) {
	src := `
X: !record
  fields:
    y: !array
      items: int
      dimensions:
        x: 10
        y: 20
  computedFields:
    y1: y[x:10, y:1]`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `index argument (10) is too large for array dimension 'x' of length 10`)
}

func TestSizeNoArgs(t *testing.T) {
	src := `
X: !record
  fields:
  computedFields:
    vSize: size()`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `size() expects 1 or 2 arguments, but called with 0`)
}

func TestSizeSingleArgWrongType(t *testing.T) {
	src := `
X: !record
  fields:
  computedFields:
    vSize: size(1)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `size() must be called with a !vector or an !array as the first argument`)
}

func TestVectorSize(t *testing.T) {
	src := `
X: !record
  fields:
    v: !vector { items: int }
  computedFields:
    vSize: size(v)`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestVectorSizeTooManyArgs(t *testing.T) {
	src := `
X: !record
  fields:
    v: !vector { items: int }
  computedFields:
    vSize: size(v, 1)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `size() does not accept a second argument when called with a !vector`)
}

func TestArraySize(t *testing.T) {
	src := `
NamedArray: !array
  items: int
X: !record
  fields:
    v: !array { items: int }
    i: int
    i2: int64
    i3: size
    namedArray: NamedArray
  computedFields:
    namedArraySize: size(namedArray)
    aSize: size(v)
    a1Size: size(v, 0)
    aiSize: size ( v , i )
    ai2Size: size(v,i2)
    ai3Size: size(v, i3)`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestArraySizeUnresolvableArg(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array { items: int }
  computedFields:
    aSize: size(v, i)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `record 'X' does not have a field or computed field named 'i'`)
}

func TestArraySizeNamedDimension(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
    stringField: string
  computedFields:
    aSize: size(v, "x")
    aSize2: size(v, stringField)
    aSize3: size(v, returnsString)
    returnsString: "\"y\""`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestArraySizeInvalidIndex(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    aSize: size(v, 5)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array dimension index is out of bounds`)
}

func TestArraySizeNegativeIndex(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    aSize: size(v, -1)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `array dimension cannot be negative`)
}

func TestArraySizeInvalidNamedDimension(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    aSize: size(v, "z")`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `this array does not have a dimension named 'z'`)
}

func TestArrayDimensionIndex(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    dimXIndeX: dimensionIndex(v, "x")
    dimYIndeX: dimensionIndex(v, 'y')`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestArrayDimensionIndexInvalidArg1Type(t *testing.T) {
	src := `
X: !record
  computedFields:
    dimIndeX: dimensionIndex(1, "x")`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `dimensionIndex() must be called with an !array as the first argument`)
}

func TestArrayDimensionIndexInvalidArg2Type(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    dimIndeX: dimensionIndex(v, 1)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the second argument to dimensionIndex() must be a dimension name string`)
}

func TestArrayDimensionIndexInvalidArgCount(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    dimIndeX: dimensionIndex(v, "x", "y")`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `dimensionIndex() expects 2 arguments, but called with 3`)
}

func TestArrayDimensionIndexInvalidName(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    dimIndeX: dimensionIndex(v, "w")`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the array does not have a dimension named 'w'`)
}

func TestArrayDimensionIndexNonamedDimensions(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions: 2
  computedFields:
    dimIndeX: dimensionIndex(v, "x")`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `dimensionIndex() is only valid for arrays with named dimensions`)
}

func TestArrayDimensionCount(t *testing.T) {
	src := `
X: !record
  fields:
    v: !array
      items: int
      dimensions:
        x:
        y:
  computedFields:
    dimCount: dimensionCount(v)`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestArrayDimensionCountInvalidArgType(t *testing.T) {
	src := `
X: !record
  computedFields:
    dimCount: dimensionCount(1)`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `dimensionCount() must be called with an !array argument`)
}

func TestSwitchErrorInTarget(t *testing.T) {
	src := `
X: !record
  computedFields:
    c:
      !switch missing:
        int: 1`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `record 'X' does not have a field or computed field named 'missing'`)
}

func TestValidSwitchWithTypePattern(t *testing.T) {
	src := `
X: !record
  fields:
    i: int
    u: [int, float]
    nu: [null, int, float]
  computedFields:
    ci:
      !switch i:
        int: i
    cu:
      !switch u:
        int: 1
        float: 2
    cnu:
      !switch nu:
        null: 0
        int: 1
        float: 2`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestSwitchImpossibleCase(t *testing.T) {
	src := `
X: !record
  fields:
    f: [float, complexfloat]
  computedFields:
    c:
      !switch f:
        int: 0
        float: 1
        complexfloat: 2`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the type is not a valid case for this switch expression`)
}

func TestSwitchImpossibleNullCase(t *testing.T) {
	src := `
X: !record
  fields:
    f: [complexfloat, float]
  computedFields:
    c:
      !switch f:
        null: 0
        complexfloat: f
        float: f`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the type is not a valid case for this switch expression`)
}

func TestSwitchDiscardMustBeLastCase(t *testing.T) {
	src := `
X: !record
  fields:
    f: [complexfloat, float]
  computedFields:
    c:
      !switch f:
        _: 1
        complexfloat: 0`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the switch case is not reachable`)
}

func TestSwitchTwoDiscards(t *testing.T) {
	src := `
X: !record
  fields:
    f: [complexfloat, float]
  computedFields:
    c:
      !switch f:
        _: 1
        _: 0`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `switch expression has no remaining cases to discard`)
}

func TestSwitchCaseRedundant(t *testing.T) {
	src := `
X: !record
  fields:
    f: [complexfloat, float]
  computedFields:
    c:
      !switch f:
        complexfloat: 1
        complexfloat i: 0`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `the switch case is not reachable`)
}

func TestSwitchNotExhaustive(t *testing.T) {
	src := `
X: !record
  fields:
    f: [complexfloat, float]
  computedFields:
    c:
      !switch f:
        complexfloat: 0`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `switch expression is not exhaustive`)
}

func TestSwitchNullCaseCannotDeclareVariable(t *testing.T) {
	src := `
X: !record
  fields:
    f: [null, float]
  computedFields:
    c:
      !switch f:
        null n: 1
        _: 0`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `a declaration pattern cannot be used with the null type`)
}

func TestSwitchFindCompatibleType(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int16, uint16]
  computedFields:
    c:
      !switch f:
        uint16 v: v
        int16 v: v`
	_, err := parseAndValidate(t, src)
	assert.Nil(t, err)
}

func TestSwitchIncompatibleCaseTypes(t *testing.T) {
	src := `
X: !record
  fields:
    f: [int64, uint64]
  computedFields:
    c:
      !switch f:
        uint64 v: v
        int64 v: v`
	_, err := parseAndValidate(t, src)
	assert.ErrorContains(t, err, `no best type was found for the switch expression`)
}
