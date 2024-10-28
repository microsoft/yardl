// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

func topologicalSortTypes(env *Environment, errorSink *validation.ErrorSink) *Environment {
	var rootSentinel Node = &RecordDefinition{}
	for _, ns := range env.Namespaces {
		sortedTypes := []TypeDefinition{}
		predecessors := make(map[Node]Node)
		visitingRecursiveType := false

		VisitWithContext(ns, rootSentinel, func(self VisitorWithContext[Node], node Node, parent Node) {
			switch t := node.(type) {
			case *ProtocolDefinition:
				// Cannot be referenced from other types
			case TypeDefinition:
				pred, found := predecessors[t]
				if found {
					if !visitingRecursiveType && pred != nil {
						nodeNameFunc := func(n Node) string {
							switch nt := n.(type) {
							case *RecordDefinition:
								return fmt.Sprintf("Record '%s'", nt.Name)
							case *EnumDefinition:
								return fmt.Sprintf("Enum '%s'", nt.Name)
							case *NamedType:
								return fmt.Sprintf("Alias '%s'", nt.Name)
							case *Field:
								return fmt.Sprintf("Field '%s'", nt.Name)
							default:
								panic("unexpected node type")
							}
						}

						path := []string{nodeNameFunc(t), nodeNameFunc(parent)}
						for pred = predecessors[parent]; pred != rootSentinel; pred = predecessors[pred] {
							path = append(path, nodeNameFunc(pred))
							if pred == t {
								break
							}
						}

						for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
							path[i], path[j] = path[j], path[i]
						}

						errorSink.Add(validationError(parent, "there is a reference cycle, which is not supported, within namespace '%s': %s", ns.Name, strings.Join(path, " -> ")))
					}
					return
				}

				predecessors[t] = parent
				self.VisitChildren(node, node)
				predecessors[t] = nil
				sortedTypes = append(sortedTypes, t)

			case *Field:
				predecessors[t] = parent
				self.VisitChildren(node, node)
				delete(predecessors, t)

			case *SimpleType:
				wasVisitingRecursiveType := visitingRecursiveType

				if t.ResolvedDefinition != nil {
					if t.IsRecursive {
						visitingRecursiveType = true
					}
					definitionMeta := t.ResolvedDefinition.GetDefinitionMeta()
					if definitionMeta.Namespace == ns.Name {
						self.Visit(env.SymbolTable.GetGenericTypeDefinition(t.ResolvedDefinition), parent)
					}
					for _, typeArg := range t.ResolvedDefinition.GetDefinitionMeta().TypeParameters {
						if typeArg.GetDefinitionMeta().Namespace == ns.Name {
							self.Visit(env.SymbolTable.GetGenericTypeDefinition(typeArg), parent)
						}
					}
				}

				self.VisitChildren(node, parent)
				visitingRecursiveType = wasVisitingRecursiveType

			default:
				self.VisitChildren(node, parent)
			}
		})

		ns.TypeDefinitions = sortedTypes
	}

	return env
}
