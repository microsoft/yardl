// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Parses a set of models and returns the latest version, previous versions, and corresponding version names
func parseVersions(t *testing.T, models []string) (*Environment, []*Environment, []string) {
	var versions []*Environment
	var labels []string

	for i, model := range models {
		env, err := parseAndValidate(t, model)
		assert.Nil(t, err)
		versions = append(versions, env)
		labels = append(labels, fmt.Sprintf("v%d", i))
	}

	return versions[len(versions)-1], versions[:len(versions)-1], labels[:len(labels)-1]
}

func TestAddProtocolSteps(t *testing.T) {
	oldModel := `
P: !protocol
  sequence:
    x: int
`
	newModel := `
P: !protocol
  sequence:
    x: int
    y: %s
`
	tests := []string{
		"bool",
		"int",
		"uint",
		"float",
		"double",
		"string",
		"complexfloat",
		"complexdouble",
		"date",
		"time",
		"datetime",
		"int[]",
		"float[,]",
		"double[4, 5]",
		"[int, float, string]",
	}

	for _, ts := range tests {
		latest, previous, labels := parseVersions(t, []string{oldModel, fmt.Sprintf(newModel, ts)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "type: %s", ts)
	}
}

func TestRemoveProtocolSteps(t *testing.T) {
	models := []string{`
P: !protocol
  sequence:
    x: int
    y: int
`, `
P: !protocol
  sequence:
    x: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestProtocolReorderSteps(t *testing.T) {
	models := []string{`
P: !protocol
  sequence:
    x: int
    y: int
`, `
P: !protocol
  sequence:
    y: int
    x: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestEnumChanges(t *testing.T) {
	base := `
X: !enum
  base: int
  values: [a, b, c]
`

	versions := []string{`
# Enum base type changed
X: !enum
  base: uint64
  values: [a, b, c]
`, `
# Enum value removed
X: !enum
  base: int
  values:
    - a
    - b
`, `
# Enum values changed due to "placement" of new value
X: !enum
  base: int
  values:
    - x
    - a
    - b
    - c
`, `
# Enum value explicitly changed
X: !enum
  base: int
  values:
    a: 3
    b: 1
    c: 2
`}

	for _, version := range versions {
		latest, previous, labels := parseVersions(t, []string{base, version})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err)
	}
}

func TestInvalidProtocolStepDefinitionChanges(t *testing.T) {
	model := `
AS: string
AO: string?
AU: [string, int]

E1: !enum
  values:
    - a
    - b
    - c

E2: !enum
  values:
    - a
    - b
    - c

F1: !flags
  values:
    - a
    - b
    - c

F2: !flags
  values:
    - a
    - b
    - c

R1: !record
  fields:
    x: string

R2: !record
  fields:
    x: string

P: !protocol
  sequence:
    step: %s
`

	tests := []struct {
		typeA string
		typeB string
	}{
		{"AS", "E1"},
		{"AS", "F1"},
		{"AS", "R1"},

		{"AO", "E1"},
		{"AO", "F1"},
		{"AO", "R1"},

		{"AU", "E1"},
		{"AU", "F1"},
		{"AU", "R1"},

		{"E1", "E2"},
		{"E1", "F1"},
		{"E1", "R1"},

		{"F1", "F2"},
		{"F1", "E1"},
		{"F1", "R1"},

		{"R1", "R2"},
		{"R1", "E1"},
		{"R1", "F1"},
	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(model, tt.typeA), fmt.Sprintf(model, tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(model, tt.typeB), fmt.Sprintf(model, tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

func TestInvalidTypeChanges(t *testing.T) {
	model := `
R: !record
  fields:
    x: %s

P: !protocol
  sequence:
    r: R
`

	tests := []struct {
		typeA string
		typeB string
	}{
		{"bool", "complexfloat"},
		{"int", "complexfloat"},
		{"uint", "complexfloat"},
		{"float", "complexfloat"},
		{"double", "complexfloat"},
		{"string", "complexfloat"},

		{"bool", "complexdouble"},
		{"int", "complexdouble"},
		{"uint", "complexdouble"},
		{"float", "complexdouble"},
		{"double", "complexdouble"},
		{"string", "complexdouble"},

		{"bool", "date"},
		{"int", "date"},
		{"uint", "date"},
		{"float", "date"},
		{"double", "date"},
		{"complexfloat", "date"},
		{"complexdouble", "date"},
		{"string", "date"},

		{"bool", "time"},
		{"int", "time"},
		{"uint", "time"},
		{"float", "time"},
		{"double", "time"},
		{"complexfloat", "time"},
		{"complexdouble", "time"},
		{"string", "time"},

		{"bool", "datetime"},
		{"int", "datetime"},
		{"uint", "datetime"},
		{"float", "datetime"},
		{"double", "datetime"},
		{"complexfloat", "datetime"},
		{"complexdouble", "datetime"},
		{"string", "datetime"},

		{"int->int", "float->int"},
		{"int->int", "int->float"},

		{"int*3", "int*4"},

		{"int[]", "float[]"},
		{"int[,]", "float[,]"},
		{"int[3]", "float[3]"},
		{"int[3]", "int[4]"},

		{"int->int", "int*"},
		{"int->int", "int[]"},
		{"int->int", "int?"},
		{"int->int", "[int, float]"},
		{"int*", "int[]"},
		{"int*", "int?"},
		{"int*", "[int, float]"},
		{"int[]", "int?"},
		{"int[]", "[int, float]"},

		{"bool?", "int?"},

		{"int?", "[int, float]"},
	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(model, tt.typeA), fmt.Sprintf(model, tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(model, tt.typeB), fmt.Sprintf(model, tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

func TestUnionChanges(t *testing.T) {
	model := `
R: !record
  fields:
    x: %s

Header: !record
  fields:
    name: string
    age: Number

Number: int

P: !protocol
  sequence:
    x: %s
    r: R
`
	formatModel := func(yardlType string) string {
		return fmt.Sprintf(model, yardlType, yardlType)
	}

	valid := []struct {
		typeA string
		typeB string
	}{
		// Optional to Union with null alternative
		{"int?", "[null, int, float]"},
		{"Header?", "[null, Header]"},

		// Union with type reordered
		{"[int, float]", "[float, int]"},
		{"[Header, Number]", "[Number, Header]"},

		// Union with types added
		{"[int, float]", "[float, int, string]"},
		{"[Header, string]", "[string, Number, Header]"},

		// Union with types removed
		{"[int, float, string]", "[float, int]"},
		{"[Header, string, Number]", "[Number, string]"},
	}

	for _, tt := range valid {
		latest, previous, labels := parseVersions(t, []string{formatModel(tt.typeA), formatModel(tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.Nil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{formatModel(tt.typeB), formatModel(tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.Nil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}

	invalid := []struct {
		typeA string
		typeB string
	}{
		// Optional to Union without null alternative
		{"int?", "[int, float]"},
		{"Header?", "[Header, Number]"},

		// Union of completely different types
		{"[bool, int]", "[float, string]"},
		{"[Header, Number]", "[string, bool]"},
	}

	for _, tt := range invalid {
		latest, previous, labels := parseVersions(t, []string{formatModel(tt.typeA), formatModel(tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{formatModel(tt.typeB), formatModel(tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

func TestNamedTypeChangesInvalidMigrationPath(t *testing.T) {
	// The change to AliasedType is valid, but there isn't a path from A to B, so ProtocolStep x's TypeChange is not valid
	models := []string{`
P: !protocol
  sequence:
    x: AliasedType
AliasedType: A
A: !record
  fields:
    i: int
`, `
P: !protocol
  sequence:
    x: AliasedType
AliasedType: [B, string]
B: !record
  fields:
    i: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)

	slices.Reverse(models)
	latest, previous, labels = parseVersions(t, models)
	_, _, err = ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestNamedTypeChangesNotSemanticallyEqual(t *testing.T) {
	// Can't use NamedTypes to add/remove a "duplicate" (structurally equal but not semantically equal) type definition
	models := []string{`
P: !protocol
  sequence:
    x: A
    y: C
A: B
B: !record
  fields:
    i: int
C: !record
  fields:
    i: int
`, `
P: !protocol
  sequence:
    x: A
    y: A
A: B
B: !record
  fields:
    i: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)

	slices.Reverse(models)
	latest, previous, labels = parseVersions(t, models)
	_, _, err = ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestNamedTypeChangesNoMigrationPath(t *testing.T) {
	// Can't use the exact same definitions but with all different names (i.e. there has to be a migration "path" between the old and new model definitions)
	//	i.e. doesn't matter if they are structurally equal! There must be a named path between versions!
	models := []string{`
P: !protocol
  sequence:
    x: A
A: B
B: C
C: !record
  fields:
    i: int
`, `
P: !protocol
  sequence:
    x: X
X: Y
Y: Z
Z: !record
  fields:
    i: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)

	slices.Reverse(models)
	latest, previous, labels = parseVersions(t, models)
	_, _, err = ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestNamedTypeInvalidTypeChanges(t *testing.T) {
	// Can't change NamedTypes if underlying type change is invalid
	// NOTE: Could add/remove levels of indirection in these tests too
	model := `
P: !protocol
  sequence:
    x: AliasedType

AliasedType: %s

R: !record
  fields:
    i: int

E: !enum
  values:
    - a

Z: !record
  fields:
    i: int

U: [Z, string]

V: float*

M: string->R
`
	typeNames := []string{
		"R",
		"E",
		"int",
		"U",
		"V",
		"M",
	}

	for _, tOld := range typeNames {
		for _, tNew := range typeNames {
			if tOld != tNew {
				latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(model, tOld), fmt.Sprintf(model, tNew)})
				_, _, err := ValidateEvolution(latest, previous, labels)
				assert.NotNil(t, err, "typeOld: %s, typeNew: %s", tOld, tNew)

				latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(model, tNew), fmt.Sprintf(model, tOld)})
				_, _, err = ValidateEvolution(latest, previous, labels)
				assert.NotNil(t, err, "typeOld: %s, typeNew: %s", tNew, tOld)
			}
		}
	}
}

// This is also tested in the evolution integration tests, but keeping it here for now as well
// The point is to test that we can add/remove levels of indirection (NamedTypes) to a GeneralizedType (e.g. Union, Optional, etc.)
// Most (if not all) of the other tests check with only SimpleTypes (e.g. Primitives, Records)
func TestGeneralizedNamedTypeChanges(t *testing.T) {
	models := []string{`
P: !protocol
  sequence:
    u: AliasedUnion
    uc: [int, string, float]
    o: string?
    oc: AliasedOptionalWithChange

AliasedUnion: [int, string]
AliasedOptionalWithChange: string?
`,
		`
P: !protocol
  sequence:
    u: [int, string]
    uc: AliasedUnionWithChange
    o: AliasedOptional
    oc: int?

AliasedOptional: string?
AliasedUnionWithChange: [float, string]
`,
	}

	latest, previous, labels := parseVersions(t, models)
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.Nil(t, err)

	slices.Reverse(models)
	latest, previous, labels = parseVersions(t, models)
	_, _, err = ValidateEvolution(latest, previous, labels)
	assert.Nil(t, err)
}

// yardl is currently strict with Generic changes.
// The following are currently invalid:
// 1. Changing the number of TypeParameters on a Generic TypeDefinition
// 1. Renaming the TypeParameters on a Generic TypeDefinition
func TestGenericDefinitionParameterChanges(t *testing.T) {
	pairs := []struct {
		previous string
		latest   string
	}{
		// Rename Generic Record TypeParameters
		{`
P: !protocol
  sequence:
    x: GenericRecord<int, string>
GenericRecord<T1, T2>: !record
  fields:
    x: T1
    y: T2
`, `
P: !protocol
  sequence:
    x: GenericRecord<int, string>
GenericRecord<A, B>: !record
  fields:
    x: A
    y: B
`},

		// Add TypeParameter to Generic Record
		{`
P: !protocol
  sequence:
    x: GenericRecord<int, string>
GenericRecord<T1, T2>: !record
  fields:
    x: T1
    y: T2
    z: bool
`, `
P: !protocol
  sequence:
    x: GenericRecord<int, string, bool>
GenericRecord<T1, T2, T3>: !record
  fields:
    x: T1
    y: T2
    z: T3
`},

		// Add/Remove TypeParameter from Generic Record
		{`
P: !protocol
  sequence:
    x: GenericRecord<int, string>
GenericRecord<T1, T2>: !record
  fields:
    x: T1
    y: T2
`, `
P: !protocol
  sequence:
    x: GenericRecord<int>
GenericRecord<T1>: !record
  fields:
    x: T1
    y: string
`},

		// Rename Generic Union TypeParameters
		{`
P: !protocol
  sequence:
    x: GenericUnion<int, string>
GenericUnion<T1, T2>: [T1, T2]
`, `
P: !protocol
  sequence:
    x: GenericUnion<int, string>
GenericUnion<A, B>: [A, B]
`},

		// Add TypeParameter to Generic Union
		{`
P: !protocol
  sequence:
    x: GenericUnion<int, string>
GenericUnion<T1, T2>: [T1, T2, bool]
`, `
P: !protocol
  sequence:
    x: GenericUnion<int, string, bool>
GenericUnion<T1, T2, T3>: [T1, T2, T3]
`},

		// Remove TypeParameter from Generic Union
		{`
P: !protocol
  sequence:
    x: GenericUnion<int, string>
GenericUnion<T1, T2>: [T1, T2]
`, `
P: !protocol
  sequence:
    x: GenericUnion<int>
GenericUnion<T1>: [T1, string]
`},
	}

	for _, pair := range pairs {
		latest, previous, labels := parseVersions(t, []string{pair.previous, pair.latest})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "previous: %s, latest: %s", pair.previous, pair.latest)

		latest, previous, labels = parseVersions(t, []string{pair.latest, pair.previous})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "previous: %s, latest: %s", pair.latest, pair.previous)
	}
}

func TestValidGenericArgChanges(t *testing.T) {
	modelA := `
P: !protocol
  sequence:
    x: %s

AliasedType<T>: T

GU2<T1, T2>: [T1, T2]
`

	modelB := `
P: !protocol
  sequence:
    x: %s

AliasedType<T>: T

GU2<T1, T2>: [T1, T2]
AU2<A, B>: GU2<A, B>

TU1<X>: [X, string]
GU1<T>: GU2<T, string>
AU1<U>: AU2<U, string>

UIS: [int, string]
GUIS: GU2<int, string>
AUIS: AU2<int, string>
TUIS: TU1<int>
GUI: GU1<int>
AUI: AU1<int>


AliasedOptional<T>: T?
MaybeString: AliasedOptional<string>
AliasedUnionWithMaybeString<T>: [T, MaybeString]
`

	tests := []struct {
		typeA string
		typeB string
	}{
		// Basic Type -> Aliased Generic Type
		{"bool", "AliasedType<bool>"},
		{"float", "AliasedType<float>"},
		{"string", "AliasedType<string>"},
		{"complexdouble", "AliasedType<complexdouble>"},
		{"string->int", "AliasedType<string->int>"},
		{"int*", "AliasedType<int*>"},
		{"int[]", "AliasedType<int[]>"},

		{"int?", "AliasedType<int?>"},
		{"int?", "AliasedType<AliasedOptional<int>>"},

		{"[int, float]", "AliasedType<AU2<int, float>>"},

		// Basic Type -> Aliased Generic Type CHANGE
		{"int", "AliasedType<float>"},
		{"string", "AliasedType<int>"},
		{"string", "AliasedType<string>"},
		{"int?", "AliasedType<float?>"},
		{"int?", "AliasedType<AliasedOptional<float>>"},
		{"[int, float]", "AliasedType<AU2<float, int>>"},
		{"[int, string]", "AliasedType<AU2<string, float>>"},

		{"[int, string]", "AliasedType<AliasedUnionWithMaybeString<int>>"},

		// Generalized Type -> Aliased Generic Named Type
		{"[int, string]", "GU2<int, string>"},
		{"[int, string]", "AU2<int, string>"},
		{"[int, string]", "TU1<int>"},
		{"[int, string]", "GU1<int>"},
		{"[int, string]", "AU1<int>"},

		{"[int, string]", "UIS"},
		{"[int, string]", "GUIS"},
		{"[int, string]", "AUIS"},
		{"[int, string]", "TUIS"},
		{"[int, string]", "GUI"},
		{"[int, string]", "AUI"},

		// Generic Generalized Type -> Aliased Generic Named Type
		{"GU2<int, string>", "GU2<int, string>"},
		{"GU2<int, string>", "AU2<int, string>"},
		{"GU2<int, string>", "TU1<int>"},
		{"GU2<int, string>", "GU1<int>"},
		{"GU2<int, string>", "AU1<int>"},
		{"GU2<int, string>", "UIS"},
		{"GU2<int, string>", "GUIS"},
		{"GU2<int, string>", "AUIS"},
		{"GU2<int, string>", "TUIS"},
		{"GU2<int, string>", "GUI"},
		{"GU2<int, string>", "AUI"},
	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(modelA, tt.typeA), fmt.Sprintf(modelB, tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.Nil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(modelB, tt.typeB), fmt.Sprintf(modelA, tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.Nil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

func TestInvalidGenericArgChanges(t *testing.T) {
	modelA := `
P: !protocol
  sequence:
    x: %s

AliasedType<T>: T

UnchangedGenericRecord<T1, T2>: !record
  fields:
    x: T1
    y: T2

ChangedGenericRecord<A, B>: !record
  fields:
    a: A
    b: B
    x: bool
`

	modelB := `
P: !protocol
  sequence:
    x: %s

AliasedType<T>: T
AliasedOptional<T>: T?

UnchangedGenericRecord<T1, T2>: !record
  fields:
    x: T1
    y: T2

ChangedGenericRecord<A, B>: !record
  fields:
    y: datetime
    b: B
    a: A
`

	tests := []struct {
		typeA string
		typeB string
	}{
		// Basic Type -> Aliased Generic Type CHANGE invalid
		{"bool", "AliasedType<datetime>"},
		{"float", "AliasedType<float*>"},
		{"string", "AliasedType<string[]>"},
		{"string->int", "AliasedType<string->float>"},

		{"int?", "AliasedType<complexdouble?>"},
		{"int?", "AliasedType<AliasedOptional<datetime>>"},

		{"UnchangedGenericRecord<int, string>", "ChangedGenericRecord<int, string>"},

		{"UnchangedGenericRecord<int, string>", "UnchangedGenericRecord<float, string>"},
		{"ChangedGenericRecord<int, string>", "ChangedGenericRecord<float, string>"},
		{"UnchangedGenericRecord<int, string>", "UnchangedGenericRecord<int, float>"},
		{"ChangedGenericRecord<int, string>", "ChangedGenericRecord<int, float>"},
	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(modelA, tt.typeA), fmt.Sprintf(modelB, tt.typeB)})
		_, _, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(modelB, tt.typeB), fmt.Sprintf(modelA, tt.typeA)})
		_, _, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

func TestUnchangedGenericAliases(t *testing.T) {
	model := `
P: !protocol
  sequence:
    x: AliasedClosedGeneric


AliasedClosedGeneric: AliasedOpenGeneric<int>
AliasedOpenGeneric<T>: GenericRecord<T>
GenericRecord<T>: !record
  fields:
    x: T
`

	latest, previous, labels := parseVersions(t, []string{model, model})
	_, _, err := ValidateEvolution(latest, previous, labels)
	assert.Nil(t, err)
}
