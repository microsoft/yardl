// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
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
		_, err := ValidateEvolution(latest, previous, labels)
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
	_, err := ValidateEvolution(latest, previous, labels)
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
	_, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestRecordChanges(t *testing.T) {
	// All RecordDefinition changes are "valid" but some may produce Warnings
	// TOOD: Mechanism for testing warnings (i.e. return them and log at top-level instead of logging them within evolution.go)
}

func TestEnumChanges(t *testing.T) {
	models := []string{`
X: !enum
  base: int
  values:
    - a
`, `
X: !enum
  base: uint64
  values:
    - a

P: !protocol
  sequence:
    x: X
`}

	latest, previous, labels := parseVersions(t, models)

	assert.NotNil(t, latest)
	assert.Len(t, previous, 1)
	assert.Len(t, labels, 1)

	_, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestStreamTypeChanges(t *testing.T) {
	models := []string{`
P: !protocol
  sequence:
    s: !stream
      items: int
`, `
P: !protocol
  sequence:
    s: !stream
      items: string
`}

	latest, previous, labels := parseVersions(t, models)
	_, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
}

func TestInvalidDefinitionChanges(t *testing.T) {
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
		_, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(model, tt.typeB), fmt.Sprintf(model, tt.typeA)})
		_, err = ValidateEvolution(latest, previous, labels)
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

		{"int*", "float*"},
		{"int*3", "float*3"},
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
		_, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{fmt.Sprintf(model, tt.typeB), fmt.Sprintf(model, tt.typeA)})
		_, err = ValidateEvolution(latest, previous, labels)
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
		_, err := ValidateEvolution(latest, previous, labels)
		assert.Nil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{formatModel(tt.typeB), formatModel(tt.typeA)})
		_, err = ValidateEvolution(latest, previous, labels)
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
		_, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeA, tt.typeB)

		latest, previous, labels = parseVersions(t, []string{formatModel(tt.typeB), formatModel(tt.typeA)})
		_, err = ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, "typeA: %s, typeB: %s", tt.typeB, tt.typeA)
	}
}

// NamedType tests:
//
// NamedType primitive changes
// Add level of indirection
// Remove level of indirection
// NamedType to underlying DIFFERENT records

func TestNamedTypeChanges(t *testing.T) {
	//
	// TODO: This is a "successful" evolution, so move it to the integration test
	//
	// I think also I should add the example from evolution/aliases (commented out)
	// where the only migration path is through an alias whose type changes from Record to Union[Record, ...]
	// 		I think this should demonstrate that the "migration path" concept works, including top-level changes to NamedTypes (oooh spooky!)
	//
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
AliasedType: B
B: !record
  fields:
    i: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, err := ValidateEvolution(latest, previous, labels)
	assert.Nil(t, err)

	// TODO: Test the following things that can't be tested by a "successful" evolution integration test
	//
	// 1. Can't change a step type from a "duplicate" record (i.e. R1 = R2 in old model) to a "different" record (i.e. R1 != R2 in new model, accounting for named type resolution)
	// 		This example is in the evolution/aliases demo in my folder
	// 1. And vice versa?
	// 1. Can't use the exact same definitions but with all different names (i.e. there has to be a migration "path" between the old and new model definitions)
	// 		This example is also in the evolution/aliases demo in my folder, but commented out
	// 1. Can't add levels of indirection to different definition types (foreach pair of {record, enum, primitive, union[different_types], vector, etc.)
	// 		I think any reasonable combination of these should cover it for now
	//		Weellll I already have a function that tests combinations of invalid type changes
	// 		I could add to that, or just make a new one here to simultaneously test NamedType resolution...

}
