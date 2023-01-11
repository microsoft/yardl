// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeEquality(t *testing.T) {
	testCases := []struct {
		spec           string
		expectedResult bool
	}{
		{`
X: [int, string]
Y: X`, true,
		},
		{`
X: !vector { items: int}
Y: X`, true,
		},
		{`
X: !vector { items: int}
Y: !vector { items: float}`, false,
		},
		{`
X: [int, string]
Y: Z
Z: X`, true,
		},
		{`
X: int
Y: Z<int>
Z<T>: T`, true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.spec, func(t *testing.T) {
			env, err := parseAndValidate(t, tC.spec)
			require.Nil(t, err)
			x := typeDefinitionByName(env, "X").(*NamedType)
			y := typeDefinitionByName(env, "Y").(*NamedType)
			require.Equal(t, tC.expectedResult, TypesEqual(GetUnderlyingType(y.Type), x.Type))
		})
	}
}

func typeDefinitionByName(env *Environment, name string) TypeDefinition {
	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {
			if td.GetDefinitionMeta().Name == name {
				return td
			}
		}
	}

	return nil
}
