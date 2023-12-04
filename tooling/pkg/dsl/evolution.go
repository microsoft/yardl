// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

type EvolutionPass func(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) (*Environment, error)

func ValidateEvolution(env *Environment, predecessor *Environment) (*Environment, error) {

	errorSink := validation.ErrorSink{}

	passes := []EvolutionPass{
		// ensureNoChanges,
		ensureBackwardCompatible,
	}

	for _, pass := range passes {
		var err error
		env, err = pass(env, predecessor, &errorSink)
		if err != nil {
			return env, err
		}
	}

	return env, errorSink.AsError()
}

// Allows reordering type definitions and record fields, but no other changes
func ensureBackwardCompatible(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) (*Environment, error) {
	visitor := func(self CompareVisitor[string], newRoot, oldRoot Node, context string) error {
		switch oldNode := oldRoot.(type) {
		case *Namespace:
			newNode, ok := newRoot.(*Namespace)
			if !ok {
				panic(fmt.Sprintf("expected a %T", oldNode))
			}

			if newNode.Name != oldNode.Name {
				return validationError(newRoot, "cannot rename namespace '%s' to '%s'", oldNode.Name, newNode.Name)
			}

			// TypeDefinitions may be reordered
			oldTds := make(map[string]TypeDefinition)
			for _, oldTd := range oldNode.TypeDefinitions {
				oldTds[oldTd.GetDefinitionMeta().Name] = oldTd
			}
			newTds := make(map[string]TypeDefinition)
			for _, newTd := range newNode.TypeDefinitions {
				newTds[newTd.GetDefinitionMeta().Name] = newTd
			}
			// TypeDefinitions may be added but not removed
			for oldName, oldTd := range oldTds {
				newTd, ok := newTds[oldName]
				if !ok {
					return validationError(newRoot, "missing type definition for '%s'", oldName)
				}
				if err := self.Compare(newTd, oldTd, context); err != nil {
					return err
				}
			}

			// Protocols may be reordered
			oldProts := make(map[string]*ProtocolDefinition)
			for _, oldProt := range oldNode.Protocols {
				oldProts[oldProt.Name] = oldProt
			}
			newProts := make(map[string]*ProtocolDefinition)
			for _, newProt := range newNode.Protocols {
				newProts[newProt.Name] = newProt
			}
			// Protocols may be added but not removed
			for oldName, oldProt := range oldProts {
				newProt, ok := newProts[oldName]
				if !ok {
					return validationError(newRoot, "missing protocol definition for '%s'", oldName)
				}
				if err := self.Compare(newProt, oldProt, context); err != nil {
					return err
				}
			}
			return nil

		case *RecordDefinition:
			newNode, ok := newRoot.(*RecordDefinition)
			if !ok {
				return validationError(newRoot, "expected a record definition")
			}

			if newNode.Name != oldNode.Name {
				return validationError(newRoot, "cannot rename record '%s' to '%s'", oldNode.Name, newNode.Name)
			}

			// Fields may be reordered
			oldFields := make(map[string]*Field)
			for _, f := range oldNode.Fields {
				oldFields[f.Name] = f
			}

			newFields := make(map[string]*Field)
			for _, f := range newNode.Fields {
				if _, ok := oldFields[f.Name]; !ok {
					if !TypeHasNullOption(f.Type) {
						return validationError(f, "cannot add new field '%s'", f.Name)
					}
				}
				newFields[f.Name] = f
			}

			for fname, oldField := range oldFields {
				newField, ok := newFields[fname]
				if !ok {
					return validationError(newRoot, "cannot remove field '%s'", fname)
				}

				if err := self.Compare(newField, oldField, context); err != nil {
					return err
				}
			}

			return nil

		default:
			return self.StrictCompare(newRoot, oldRoot, context)
		}
	}

	if err := Compare[string](env, predecessor, "", visitor); err != nil {
		return env, err
	}
	return env, nil
}
