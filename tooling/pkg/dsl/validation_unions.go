// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func assignUnionCaseLabels(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		if t, ok := node.(*GeneralizedType); ok && t.Cases.IsUnion() {
			// assign labels to union cases
			for _, typeCase := range t.Cases {
				if !typeCase.IsNullType() {
					typeCase.Label = typeLabel(typeCase.Type, true)
				}
			}

			duplicates := make(map[string][]int)
			for i, typeCase := range t.Cases {
				if !typeCase.IsNullType() {
					duplicates[typeCase.Label] = append(duplicates[typeCase.Label], i)
				}
			}

			for _, v := range duplicates {
				if len(v) > 1 {
					for _, i := range v {
						t.Cases[i].Label = typeLabel(t.Cases[i].Type, false)
					}
				}
			}

			labels := make(map[string]any)
			for _, item := range t.Cases {
				if item != nil {
					if _, found := labels[item.Label]; found {
						errorSink.Add(validationError(node, "internal error: union cases must have distinct labels within the union"))
					}
				}
			}
		}

		self.VisitChildren(node)
	})

	return env
}

func validateUnionCases(env *Environment, errorSink *validation.ErrorSink) *Environment {
	if len(errorSink.Errors) > 0 {
		// Only perform this if all types are resolved
		return env
	}

	VisitWithContext(env, false, func(self VisitorWithContext[bool], node Node, visitingReference bool) {
		switch t := node.(type) {
		case *GeneralizedType:
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

				for i, item := range t.Cases {
					for j := i + 1; j < len(t.Cases); j++ {
						otherItem := t.Cases[j]

						if TypesEqual(item.Type, otherItem.Type) {
							additionalExplanation := ""
							if !item.IsNullType() {
								// determine if this is because size and uint64 were used, which are equivalent but not aliases
								if itemPrimitive, ok := GetPrimitiveType(item.Type); ok {
									if otherItemPrimitive, ok := GetPrimitiveType(otherItem.Type); ok {
										if itemPrimitive != otherItemPrimitive &&
											(itemPrimitive == PrimitiveUint64 && otherItemPrimitive == PrimitiveSize ||
												itemPrimitive == PrimitiveSize && otherItemPrimitive == PrimitiveUint64) {
											additionalExplanation = " (uint64 and size are equivalent)"
										}
									}
								}

								// Determine if the types are defined at a different location than the cases
								// This indicates that the cause of the duplicate is a type argument.

								itemNodeMeta := item.GetNodeMeta()
								itemTypeNodeMeta := item.Type.GetNodeMeta()

								otherItemNodeMeta := otherItem.GetNodeMeta()
								otherItemTypeNodeMeta := otherItem.Type.GetNodeMeta()

								itemDefinedElsewhere := !itemNodeMeta.Equals(itemTypeNodeMeta)
								otherItemDefinedElsewhere := !otherItemNodeMeta.Equals(otherItemTypeNodeMeta)

								if itemDefinedElsewhere || otherItemDefinedElsewhere {
									if itemDefinedElsewhere && otherItemDefinedElsewhere {
										// both are type arguments
										errorSink.Add(validationError(item, "redundant union type cases resulting from the type arguments given at %s and %s%s", itemTypeNodeMeta, otherItemTypeNodeMeta, additionalExplanation))
										continue
									}

									// only one is a type argument
									var typeParameterNode Node
									var redundantNode Node
									if itemDefinedElsewhere {
										typeParameterNode = itemTypeNodeMeta
										redundantNode = otherItemNodeMeta
									} else {
										typeParameterNode = otherItemTypeNodeMeta
										redundantNode = itemNodeMeta
									}

									errorSink.Add(validationError(redundantNode, "redundant union type cases resulting from the type argument given at %s%s", typeParameterNode, additionalExplanation))
									continue
								}
							}
							// No contributions from type arguments.
							// To avoid reporting the same error multiple times, we only report the error
							// if we we visiting the type directly, i.e. not through a reference.
							if !visitingReference {
								errorSink.Add(validationError(item, "redundant union type cases%s", additionalExplanation))
							}
						}
					}
				}
			}

			self.VisitChildren(node, visitingReference)

		case *SimpleType:
			if len(t.ResolvedDefinition.GetDefinitionMeta().TypeArguments) > 0 {
				// Check the referenced type with the type arguments provided
				self.Visit(t.ResolvedDefinition, true)
			}
		default:
			self.VisitChildren(node, visitingReference)
		}
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
