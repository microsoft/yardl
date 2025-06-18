// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"errors"
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func buildSymbolTable(env *Environment, errorSink *validation.ErrorSink) *Environment {
	VisitWithContext(env, "", func(self VisitorWithContext[string], node Node, namespace string) {
		switch t := node.(type) {
		case *Namespace:
			self.VisitChildren(node, t.Name)
		case TypeDefinition:
			meta := t.GetDefinitionMeta()

			if meta.Name == "" {
				errorSink.Add(validationError(node, "the name field must be provided and non-empty"))
				return
			}

			if _, found := primitiveTypes[meta.Name]; found {
				errorSink.Add(validationError(node, "the name '%s' is reserved", meta.Name))
				return
			}

			meta.Namespace = namespace
			fullName := fmt.Sprintf("%s.%s", namespace, meta.Name)

			if other, exists := (env.SymbolTable)[fullName]; exists {
				otherMeta := other.GetNodeMeta()
				errorSink.Add(validationError(node, "the name '%s' is already defined in file '%s' line '%d'", meta.Name, otherMeta.File, otherMeta.Line))
			} else {
				env.SymbolTable[fullName] = t
			}

		default:
			self.VisitChildren(node, namespace)
		}
	})

	return env
}

func resolveTypes(env *Environment, errorSink *validation.ErrorSink) *Environment {
	type visitorContext struct {
		currentNamespace string
		symbolTable      SymbolTable
	}

	VisitWithContext(env, &visitorContext{symbolTable: env.SymbolTable}, func(self VisitorWithContext[*visitorContext], node Node, context *visitorContext) {
		switch t := node.(type) {
		case *Namespace:
			self.VisitChildren(node, &visitorContext{currentNamespace: t.Name, symbolTable: env.SymbolTable})
			return
		case TypeDefinition:
			definitionMeta := t.GetDefinitionMeta()
			if len(definitionMeta.TypeParameters) > 0 {
				scopedSymbolTable := context.symbolTable.Clone()
				for _, genericTypeParameter := range definitionMeta.TypeParameters {
					scopedSymbolTable[genericTypeParameter.Name] = genericTypeParameter
				}
				self.VisitChildren(node, &visitorContext{symbolTable: scopedSymbolTable, currentNamespace: context.currentNamespace})
			} else {
				self.VisitChildren(node, context)
			}
			return
		case *SimpleType:
			self.VisitChildren(node, context)
			err := resolveType(t, context.currentNamespace, context.symbolTable, true)
			if err != nil {
				errorSink.Add(validationError(t, "%s", err.Error()))
				break
			}
		}

		self.VisitChildren(node, context)
	})

	return env
}

func convertGenericReferences(env *Environment, errorSink *validation.ErrorSink) *Environment {
	if len(errorSink.Errors) > 0 {
		return env
	}

	type visitorContext struct {
		symbolTable      SymbolTable
		currentNamespace string
	}

	VisitWithContext(env, visitorContext{symbolTable: env.SymbolTable}, func(self VisitorWithContext[visitorContext], node Node, context visitorContext) {
		switch t := node.(type) {
		case *Namespace:
			self.VisitChildren(node, visitorContext{context.symbolTable, t.Name})
			return
		case TypeDefinition:
			definitionMeta := t.GetDefinitionMeta()
			if len(definitionMeta.TypeParameters) > 0 {
				scopedSymbolTable := context.symbolTable.Clone()
				for _, genericTypeParameter := range definitionMeta.TypeParameters {
					scopedSymbolTable[genericTypeParameter.Name] = genericTypeParameter
				}
				self.VisitChildren(node, visitorContext{scopedSymbolTable, context.currentNamespace})
			} else {
				self.VisitChildren(node, context)
			}
			return
		case *SimpleType:
			self.VisitChildren(node, context)
			err := resolveType(t, context.currentNamespace, context.symbolTable, false)
			if err != nil {
				errorSink.Add(validationError(t, "%s", err.Error()))
				break
			}
		}

		self.VisitChildren(node, context)
	})

	return env
}

func resolveTypeByName(typeName string, currentNamespace string, symbolTable SymbolTable) (TypeDefinition, error) {
	if primitiveType, found := primitiveTypes[typeName]; found {
		return primitiveType, nil
	}

	resolvedType, found := symbolTable[typeName]
	if !found {
		qualifiedName := fmt.Sprintf("%s.%s", currentNamespace, typeName)
		resolvedType, found = symbolTable[qualifiedName]
		if !found {
			return nil, fmt.Errorf("the type '%s' is not recognized", typeName)
		}
	}

	if _, isProtocol := resolvedType.(*ProtocolDefinition); isProtocol {
		return nil, errors.New("cannot reference a protocol")
	}

	return resolvedType, nil
}

func resolveType(simpleType *SimpleType, currentNamespace string, symbolTable SymbolTable, shallow bool) error {
	resolvedTypeDefinition, err := resolveTypeByName(simpleType.Name, currentNamespace, symbolTable)
	if err != nil {
		return err
	}

	meta := resolvedTypeDefinition.GetDefinitionMeta()
	simpleType.Name = meta.GetQualifiedName()

	if len(meta.TypeParameters) != len(simpleType.TypeArguments) {
		return fmt.Errorf("'%s' was given %d type argument(s) but has %d type parameter(s)", meta.Name, len(simpleType.TypeArguments), len(meta.TypeParameters))
	}

	if len(meta.TypeParameters) == 0 {
		simpleType.ResolvedDefinition = resolvedTypeDefinition
		return nil
	}

	simpleType.ResolvedDefinition, err = MakeGenericType(resolvedTypeDefinition, simpleType.TypeArguments, shallow)
	return err
}

func validateGenericTypeDefinitions(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *RecordDefinition, *NamedType:
			return
		case TypeDefinition:
			meta := node.GetDefinitionMeta()
			if len(meta.TypeParameters) > 0 {
				errorSink.Add(validationError(node, "'%s' cannot have generic type parameters", meta.Name))
			}
		default:
			self.VisitChildren(node)
		}
	})

	return env
}

func validateGenericParametersUsed(env *Environment, errorSink *validation.ErrorSink) *Environment {
	if len(errorSink.Errors) > 0 {
		return env
	}

	VisitWithContext(env, nil, func(self VisitorWithContext[map[*GenericTypeParameter]any], node Node, context map[*GenericTypeParameter]any) {
		switch t := node.(type) {
		case TypeDefinition:
			meta := t.GetDefinitionMeta()
			if len(meta.TypeParameters) == 0 {
				break
			}

			usedTypeParameters := make(map[*GenericTypeParameter]any)
			for _, p := range meta.TypeParameters {
				usedTypeParameters[p] = nil
			}
			self.VisitChildren(node, usedTypeParameters)
			for p := range usedTypeParameters {
				errorSink.Add(validationError(p, "generic type parameter '%s' is not used", p.Name))
			}

			return

		case *SimpleType:
			if gen, ok := t.ResolvedDefinition.(*GenericTypeParameter); ok {
				delete(context, gen)
			}
		}

		self.VisitChildren(node, context)
	})

	return env
}
