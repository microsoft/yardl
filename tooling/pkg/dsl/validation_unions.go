// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func validateUnionCases(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		if t, ok := node.(*GeneralizedType); ok {
			cases := t.Cases
			if len(cases) == 0 {
				errorSink.Add(validationError(node, "a field cannot be a union type with no options"))
			}

			if len(cases) == 1 && cases[0].IsNullType() {
				errorSink.Add(validationError(node, "null cannot be the only option in a union type"))
			}

			for i, typeCase := range cases {
				if typeCase.IsNullType() && i != 0 {
					errorSink.Add(validationError(node, "if null is specified in a union type, it must be the first option"))
				}
			}

			if len(cases) > 1 {
				for _, typeCase := range cases {
					if childType, ok := typeCase.Type.(*GeneralizedType); ok && len(childType.Cases) > 1 {
						errorSink.Add(validationError(typeCase, "unions may not immediately contain other unions"))
					}
				}
			}

			if cases.IsUnion() {
				// assign labels to union cases
				for _, typeCase := range cases {
					if !typeCase.IsNullType() {
						typeCase.Label = typeLabel(typeCase.Type, true)
					}
				}

				duplicates := make(map[string][]int)
				for i, typeCase := range cases {
					if !typeCase.IsNullType() {
						duplicates[typeCase.Label] = append(duplicates[typeCase.Label], i)
					}
				}

				for _, v := range duplicates {
					if len(v) > 1 {
						for _, i := range v {
							cases[i].Label = typeLabel(cases[i].Type, false)
						}
					}
				}

				for i, item := range cases {
					for j := i + 1; j < len(cases); j++ {
						if TypesEqual(item.Type, cases[j].Type) {
							errorSink.Add(validationError(item, "all type cases in a union must be distinct"))
						}
					}
				}

				labels := make(map[string]any)
				for _, item := range cases {
					if item != nil {
						if _, found := labels[item.Label]; found {
							errorSink.Add(validationError(node, "union cases must have distict labels within the union"))
						}
					}
				}
			}
		}

		self.VisitChildren(node)
	})

	return env
}

func typeLabel(t Type, simple bool) string {
	switch t := t.(type) {
	case nil:
		return "null"
	case *SimpleType:
		baseName := func() string {
			if simple && t.ResolvedDefinition != nil {
				return t.ResolvedDefinition.GetDefinitionMeta().Name
			}
			return t.Name
		}()

		if len(t.TypeArguments) == 0 {
			return baseName
		}
		typeArguments := make([]string, len(t.TypeArguments))
		for i, typeArg := range t.TypeArguments {
			typeArguments[i] = typeLabel(typeArg, simple)
		}

		return fmt.Sprintf("%s<%s>", baseName, strings.Join(typeArguments, ","))

	case *GeneralizedType:
		casesLabel := func() string {
			if len(t.Cases) == 1 {
				return typeLabel(t.Cases[0].Type, simple)
			}

			caseLabels := make([]string, len(t.Cases))
			for i, typeCase := range t.Cases {
				caseLabels[i] = typeLabel(typeCase.Type, simple)
			}

			return fmt.Sprintf("{%s}", strings.Join(caseLabels, ","))
		}()

		switch d := t.Dimensionality.(type) {
		case nil:
			return casesLabel
		case *Vector:
			simpleLabel := casesLabel + "Vector"
			if simple || d.Length == nil {
				return simpleLabel
			}

			return fmt.Sprintf("%s[%d]", simpleLabel, *d.Length)
		case *Array:
			simpleLabel := casesLabel + "Array"
			if simple || !d.HasKnownNumberOfDimensions() {
				return simpleLabel
			}

			if d.IsFixed() {
				dims := make([]string, len(*d.Dimensions))
				for i, dim := range *d.Dimensions {
					dims[i] = strconv.FormatUint(*dim.Length, 10)
				}

				return fmt.Sprintf("%s[%s]", simpleLabel, strings.Join(dims, ","))
			}

			return fmt.Sprintf("%s[%s]", simpleLabel, strings.Repeat(",", len(*d.Dimensions)))

		default:
			panic(fmt.Sprintf("unexpected type %T", d))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}
