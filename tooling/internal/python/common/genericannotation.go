// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package common

import "github.com/microsoft/yardl/tooling/pkg/dsl"

type TypeParameterUse int

const TypeParameterUseAnnotationKey = "typeParameterUse"

const (
	TypeParameterUseNone   TypeParameterUse = 0
	TypeParameterUseScalar TypeParameterUse = 1
	TypeParameterUseArray  TypeParameterUse = 2
	TypeParameterUseBoth   TypeParameterUse = TypeParameterUseScalar | TypeParameterUseArray
)

func AnnotateGenerics(env *dsl.Environment) {
	dsl.Visit(env, func(self dsl.Visitor, node dsl.Node) {
		switch node := node.(type) {
		case *dsl.ProtocolDefinition:
			return
		case dsl.TypeDefinition:
			definitionMeta := node.GetDefinitionMeta()
			if len(definitionMeta.TypeParameters) == 0 {
				return
			}

			for _, typeParameter := range definitionMeta.TypeParameters {
				use := GetTypeParameterUse(node, typeParameter)
				if typeParameter.Annotations == nil {
					typeParameter.Annotations = make(map[string]any)
				}
				typeParameter.Annotations["typeParameterUse"] = use
			}
		default:
			self.VisitChildren(node)
		}
	})
}

func GetTypeParameterUse(root dsl.Node, typeParameter *dsl.GenericTypeParameter) TypeParameterUse {
	use := TypeParameterUseNone
	dsl.VisitWithContext(root, false, func(self dsl.VisitorWithContext[bool], node dsl.Node, inArray bool) {
		switch node := node.(type) {
		case *dsl.GeneralizedType:
			switch node.Dimensionality.(type) {
			case *dsl.Array:
				if node.Cases.IsSingle() {
					if st, ok := node.Cases[0].Type.(*dsl.SimpleType); ok {
						switch st.ResolvedDefinition.(type) {
						case *dsl.GenericTypeParameter:
							self.VisitChildren(node, true)
						}
					}
				}
				return
			}
		case *dsl.SimpleType:
			if node.ResolvedDefinition == typeParameter {
				if inArray {
					use |= TypeParameterUseArray
				} else {
					use |= TypeParameterUseScalar
				}

				return
			}

			if !node.IsRecursive && len(node.ResolvedDefinition.GetDefinitionMeta().TypeParameters) > 0 {
				self.Visit(node.ResolvedDefinition, inArray)
				return
			}
		}

		self.VisitChildren(node, inArray)
	})

	return use
}
