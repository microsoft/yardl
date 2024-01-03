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

func TestProtocolAddSteps(t *testing.T) {
	models := []string{`
P: !protocol
  sequence:
    x: int
`, `
P: !protocol
  sequence:
    x: int
    y: int
`}

	latest, previous, labels := parseVersions(t, models)
	_, err := ValidateEvolution(latest, previous, labels)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "adding steps to a Protocol")
}

func TestProtocolRemoveSteps(t *testing.T) {
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
	assert.ErrorContains(t, err, "removing steps from a Protocol")
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
	assert.ErrorContains(t, err, "reordering steps in a Protocol")
}

func TestRecordChanges(t *testing.T) {
	// All RecordDefinition changes are "valid" but some may produce Warnings
	// TOOD: Mechanism for capturing warnings (i.e. return them and log at top-level instead of logging them within evolution.go)
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
	assert.ErrorContains(t, err, "base type")
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
		// errorContains string
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

		{"[bool, int]", "[int, bool]"},
		// {"[bool, int]", "[bool, float]"},
		// {"[bool, int]", "[bool, int, float]"},

	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(model, tt.typeA), fmt.Sprintf(model, tt.typeB)})
		_, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, fmt.Sprintf("typeA: %s, typeB: %s", tt.typeA, tt.typeB))
	}

	for _, tt := range tests {
		latest, previous, labels := parseVersions(t, []string{fmt.Sprintf(model, tt.typeB), fmt.Sprintf(model, tt.typeA)})
		_, err := ValidateEvolution(latest, previous, labels)
		assert.NotNil(t, err, fmt.Sprintf("typeA: %s, typeB: %s", tt.typeA, tt.typeB))
	}

}
