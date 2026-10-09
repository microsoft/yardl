// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"encoding/json"
	"sort"
)

type ProtocolSchema struct {
	Protocol *ProtocolDefinition `json:"protocol"`
	Types    []TypeDefinition    `json:"types"`
}

func GetProtocolSchema(protocol *ProtocolDefinition, symbolTable SymbolTable) *ProtocolSchema {
	baseNamespace := protocol.Namespace
	schema := &ProtocolSchema{Protocol: removeComments(protocol)}
	visitedTypeDefinitions := make(map[TypeDefinition]any)
	Visit(protocol, func(self Visitor, node Node) {
		switch t := node.(type) {
		case *ProtocolDefinition:
			break
		case PrimitiveDefinition:
			break
		case *GenericTypeParameter:
			break
		case TypeDefinition:
			if _, visited := visitedTypeDefinitions[t]; visited {
				return
			}

			visitedTypeDefinitions[t] = nil

			// We don't want to include computed fields in the schema json
			// since they are not used for (de)serialization.
			if rec, ok := t.(*RecordDefinition); ok {
				clone := *rec
				clone.ComputedFields = nil
				t = &clone
			}

			schema.Types = append(schema.Types, qualifyImportedName(removeComments(t), baseNamespace))

		case *SimpleType:
			self.Visit(symbolTable.GetGenericTypeDefinition(t.ResolvedDefinition))
			for _, typeArg := range t.ResolvedDefinition.GetDefinitionMeta().TypeParameters {
				self.Visit(symbolTable.GetGenericTypeDefinition(typeArg))
			}
		}

		self.VisitChildren(node)
	})

	sort.Slice(schema.Types, func(i, j int) bool {
		return schema.Types[i].GetDefinitionMeta().GetQualifiedName() < schema.Types[j].GetDefinitionMeta().GetQualifiedName()
	})

	return schema
}

// qualifyImportedName ensures the definition's serialized name uniquely
// identifies it in the embedded schema. Definitions from the protocol's own
// namespace keep their unqualified name (preserving existing schema strings
// for models that don't cross namespaces). Definitions from other namespaces
// are rewritten to use their fully-qualified name so that consumers which
// resolve types by name can disambiguate cross-namespace collisions.
func qualifyImportedName(t TypeDefinition, baseNamespace string) TypeDefinition {
	meta := t.GetDefinitionMeta()
	if meta.Namespace == "" || meta.Namespace == baseNamespace {
		return t
	}
	qualified := meta.GetQualifiedName()
	if meta.Name == qualified {
		return t
	}
	newMeta := *meta
	newMeta.Name = qualified
	switch v := t.(type) {
	case *RecordDefinition:
		clone := *v
		clone.DefinitionMeta = &newMeta
		return &clone
	case *NamedType:
		clone := *v
		clone.DefinitionMeta = &newMeta
		return &clone
	case *EnumDefinition:
		clone := *v
		clone.DefinitionMeta = &newMeta
		return &clone
	default:
		return t
	}
}

func removeComments[T Node](typeDefinition T) T {
	return Rewrite(typeDefinition, func(self *Rewriter, node Node) Node {
		switch t := node.(type) {
		case *DefinitionMeta:
			if t.Comment == "" {
				return t
			}

			clone := *t
			clone.Comment = ""
			return &clone

		case *Field:
			if t.Comment == "" {
				return self.DefaultRewrite(t)
			}

			clone := *t
			clone.Comment = ""
			return self.DefaultRewrite(&clone)
		case *ProtocolStep:
			if t.Comment == "" {
				return self.DefaultRewrite(t)
			}

			clone := *t
			clone.Comment = ""
			return self.DefaultRewrite(&clone)
		case *ArrayDimension:
			if t.Comment == "" {
				return t
			}
			clone := *t
			clone.Comment = ""
			return &clone

		case *EnumValue:
			if t.Comment == "" {
				return t
			}
			clone := *t
			clone.Comment = ""
			return &clone

		default:
			return self.DefaultRewrite(t)
		}
	}).(T)
}

func GetProtocolSchemaString(protocol *ProtocolDefinition, symbolTable SymbolTable) string {
	schema := GetProtocolSchema(protocol, symbolTable)
	bytes, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
