// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func assignUnionCaseTags(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		if t, ok := node.(*GeneralizedType); ok && (t.Cases.IsUnion() || (len(t.Cases) > 0 && t.Cases[0].Tag != "")) {
			// assign tags to union cases
			for _, typeCase := range t.Cases {
				if typeCase.Tag == "" {
					typeCase.Tag = TypeToShortSyntax(typeCase.Type, false)
				}
			}
		}

		self.VisitChildren(node)
	})

	return env
}

func containsOpenGeneric(node Node) bool {
	res := false
	Visit(node, func(self Visitor, node Node) {
		switch t := node.(type) {
		case *GenericTypeParameter:
			res = true
			return
		case *SimpleType:
			if _, ok := t.ResolvedDefinition.(*GenericTypeParameter); ok {
				res = true
				return
			}
		}
		self.VisitChildren(node)
	})

	return res
}

func validateUnionCases(env *Environment, errorSink *validation.ErrorSink) *Environment {
	if len(errorSink.Errors) > 0 {
		// Only perform this if all types are resolved
		return env
	}

	tagTypeMap := make(map[string]Type)

	VisitWithContext(env, false, func(self VisitorWithContext[bool], node Node, visitingReference bool) {
		switch t := node.(type) {
		case *GeneralizedType:
			errorCountSnapshot := len(errorSink.Errors)
			cases := t.Cases
			if len(cases) == 0 {
				errorSink.Add(validationError(node, "a union type must have at least one option"))
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
							// if we are visiting the type directly, i.e. not through a reference.
							if !visitingReference {
								errorSink.Add(validationError(item, "redundant union type cases%s", additionalExplanation))
							}
						}
					}
				}
			}

			// validate tags
			if (t.Cases.IsUnion() || len(t.Cases) > 0 && t.Cases[0].ExplicitTag) &&
				len(errorSink.Errors) == errorCountSnapshot && !visitingReference {
				for _, typeCase := range t.Cases {
					if typeCase.ExplicitTag {
						if !memberNameRegex.MatchString(typeCase.Tag) {
							errorSink.Add(validationError(typeCase, "union tag '%s' must be camelCased matching the format %s", typeCase.Tag, memberNameRegex.String()))
						}
					} else if !memberNameRegex.MatchString(strings.ToLower(typeCase.Tag)) {
						explicitExample := fmt.Sprintf("!union { myTag: \"%s\", ... }", typeCase.Tag)
						if containsOpenGeneric(t) {
							errorSink.Add(
								validationError(
									typeCase, "the type '%s' cannot be used as a tag for the union case. An explicit tag can be given using the `!union` syntax (e.g. `%s`)",
									typeCase.Tag, explicitExample))
						} else {
							aliasExample := fmt.Sprintf("MyTypeAlias = %s\nMyUnion = [..., MyTypeAlias, ...]", typeCase.Tag)
							errorSink.Add(
								validationError(
									typeCase, "the type '%s' cannot be used as a tag for the union case. Explicit tags can be given using the `!union` syntax (e.g. `%s`) or the type can be aliased for the type case (e.g. `%s`)",
									typeCase.Tag, explicitExample, aliasExample))
						}
					}
				}

				tags := make(map[string]any)
				areCustomTags := false
				for _, item := range t.Cases {
					if item != nil {
						areCustomTags = areCustomTags || item.ExplicitTag
						if _, found := tags[item.Tag]; found {
							errorSink.Add(validationError(node, "all union cases must have distinct tags"))
						} else {
							tags[item.Tag] = nil
						}
					}
				}

				if areCustomTags {
					keys := make([]string, 0, len(tags))
					for k := range tags {
						keys = append(keys, k)
					}
					sort.Slice(keys, func(i, j int) bool {
						return keys[i] < keys[j]
					})
					tagsString := strings.Join(keys, ", ")

					if existing, found := tagTypeMap[tagsString]; found {
						if !TypesEqual(existing, t.ToScalar()) {
							existingNodeMeta := existing.GetNodeMeta()
							errorSink.Add(validationError(node, "the combination of tags used by the union are already in use with different types in file '%s' line '%d'", existingNodeMeta.File, existingNodeMeta.Line))
						}
					} else {
						tagTypeMap[tagsString] = t.ToScalar()
					}
				}
			}

			self.VisitChildren(node, visitingReference)

		case *SimpleType:
			if len(t.ResolvedDefinition.GetDefinitionMeta().TypeArguments) > 0 {
				if !t.IsRecursive {
					// Check the referenced type with the type arguments provided
					self.Visit(t.ResolvedDefinition, true)
				}
			}
		default:
			self.VisitChildren(node, visitingReference)
		}
	})

	return env
}
