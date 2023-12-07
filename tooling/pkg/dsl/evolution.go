// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

func CollectChanges(env *Environment, predecessor *Environment, versionId int) (*Environment, error) {
	fmt.Println("Collecting changes")

	// Pre-process the predecessor Protocols to annotate them with their protocol string
	// Instead of trying to stuff it into a context parameter in recursive Comparison functions
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations["schema"] = GetProtocolSchemaString(node, predecessor.SymbolTable)
			return
		default:
			self.VisitChildren(node)
		}
	})

	// Pre-process the new Model to prepare Annotation slices
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations["changes"] == nil {
				node.GetDefinitionMeta().Annotations["changes"] = make([]*ProtocolDefinition, 0)
			}
			if node.GetDefinitionMeta().Annotations["schemas"] == nil {
				node.GetDefinitionMeta().Annotations["schemas"] = make([]string, 0)
			}
			self.VisitChildren(node)

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations["changes"] == nil {
				node.GetDefinitionMeta().Annotations["changes"] = make([]TypeDefinition, 0)
			}
			return

		case *ProtocolStep:
			if node.Annotations == nil {
				node.Annotations = make(map[string]any)
			}
			if node.Annotations["changes"] == nil {
				node.Annotations["changes"] = make([]*ProtocolStep, 0)
			}
			return

		default:
			self.VisitChildren(node)
		}
	})

	detectChanges(env, predecessor, versionId)

	return env, nil
}

// func GetChange(node Node, version_index int) Node {
// 	changes, ok := node.GetDefinitionMeta().Annotations["changes"]
// }

// TODO: Rename to "annotateChanges"
func detectChanges(newRoot Node, oldRoot Node, version_index int) Node {
	switch newNode := newRoot.(type) {
	case *Environment:
		/* 	Ensure each Namespace matches */
		oldNode := oldRoot.(*Environment)

		if len(newNode.Namespaces) != len(oldNode.Namespaces) {
			panic("mismatch in number of namespaces")
		}
		for i := range newNode.Namespaces {
			// changes := detectChanges(newNode.Namespaces[i], oldNode.Namespaces[i])
			// if changes != nil {
			// 	newNode.Namespaces[i].VersionChanges = append(newNode.Namespaces[i].VersionChanges, changes.(*Namespace))
			// }

			detectChanges(newNode.Namespaces[i], oldNode.Namespaces[i], version_index)
		}

		return nil

	case *Namespace:
		/* 	For each matching TypeDefinition, detect changes:
			If change is detected, append the Change to output Namespace.TypeDefinitions
			Otherwise, skip it.

		For each matching ProtocolDefinition, detect changes:
			If change is detected, append the Change to output Namespace.Protocols
		*/
		oldNode := oldRoot.(*Namespace)

		if newNode.Name != oldNode.Name {
			// return validationError(newRoot, "cannot rename namespace '%s' to '%s'", oldNode.Name, newNode.Name)
			panic("cannot rename namespace")
		}

		result := *oldNode
		// We only need to save the Definitions that changed between schemas
		result.TypeDefinitions = nil
		result.Protocols = nil
		changed := false

		oldTds := make(map[string]TypeDefinition)
		for _, oldTd := range oldNode.TypeDefinitions {
			oldTds[oldTd.GetDefinitionMeta().Name] = oldTd
		}
		newTds := make(map[string]TypeDefinition)
		for i, newTd := range newNode.TypeDefinitions {
			newTds[newTd.GetDefinitionMeta().Name] = newTd
			_, ok := oldTds[newTd.GetDefinitionMeta().Name]
			if !ok {
				// CHANGE: New TypeDefinition
				changed = true
				continue
			}
			if i > len(oldNode.TypeDefinitions) {
				// CHANGE: Reordered TypeDefinitions
				changed = true
				continue
			}
			if newTd.GetDefinitionMeta().Name != oldNode.TypeDefinitions[i].GetDefinitionMeta().Name {
				// CHANGE: Reordered TypeDefinitions
				changed = true
				continue
			}
		}

		for _, oldTd := range oldNode.TypeDefinitions {
			newTd, ok := newTds[oldTd.GetDefinitionMeta().Name]
			if !ok {
				// return validationError(newRoot, "missing type definition for '%s'", oldName)
				log.Warn().Msgf("missing type definition")
				changed = true
				continue
			}

			var changedTypeDef TypeDefinition
			if ch := detectChanges(newTd, oldTd, version_index); ch != nil {
				changedTypeDef = ch.(TypeDefinition)
				result.TypeDefinitions = append(result.TypeDefinitions, changedTypeDef)
				changed = true
			}

			// Mark the "new" TypeDefinition as having changed from previous version.
			// This annotation is used when recursively comparing Types for changes
			newTd.GetDefinitionMeta().Annotations["changes"] = append(newTd.GetDefinitionMeta().Annotations["changes"].([]TypeDefinition), changedTypeDef)
		}

		// Protocols may be reordered
		// Don't care about Protocols that were added or removed
		oldProts := make(map[string]*ProtocolDefinition)
		for _, oldProt := range oldNode.Protocols {
			oldProts[oldProt.Name] = oldProt
		}

		newProts := make(map[string]*ProtocolDefinition)
		for _, newTd := range newNode.Protocols {
			newProts[newTd.GetDefinitionMeta().Name] = newTd
		}
		for i, newProt := range newNode.Protocols {
			_, ok := oldProts[newProt.Name]
			if !ok {
				// CHANGE: New Protocol
				changed = true
				continue
			}
			if i > len(oldNode.Protocols) {
				// CHANGE: Reordered Protocols
				changed = true
				continue
			}
			if newProt.GetDefinitionMeta().Name != oldNode.TypeDefinitions[i].GetDefinitionMeta().Name {
				// CHANGE: Reordered Protocols
				changed = true
				continue
			}
		}

		for _, oldProt := range oldNode.Protocols {
			newProt, ok := newProts[oldProt.GetDefinitionMeta().Name]
			if !ok {
				// return validationError(newRoot, "missing type definition for '%s'", oldName)
				log.Warn().Msgf("missing protocol definition")
				changed = true
				continue
			}

			var changedProtocolDef *ProtocolDefinition
			if ch := detectChanges(newProt, oldProt, version_index); ch != nil {
				changedProtocolDef = ch.(*ProtocolDefinition)
				oldSchema, ok := oldProt.GetDefinitionMeta().Annotations["schema"]
				if !ok {
					panic("Expected annotation containing old protocol schema string")
				}
				newProt.GetDefinitionMeta().Annotations["schemas"] = append(newProt.GetDefinitionMeta().Annotations["schemas"].([]string), oldSchema.(string))

				result.Protocols = append(result.Protocols, changedProtocolDef)
				changed = true
			}

			// Mark the "new" TypeDefinition as having changed from previous version.
			// This annotation is used when recursively comparing Types for changes
			newProt.GetDefinitionMeta().Annotations["changes"] = append(newProt.GetDefinitionMeta().Annotations["changes"].([]*ProtocolDefinition), changedProtocolDef)
		}

		if changed {
			oldNode.TypeDefinitions = result.TypeDefinitions
			oldNode.Protocols = result.Protocols
			return &result
		}
		return nil

	case *ProtocolDefinition:
		oldNode := oldRoot.(*ProtocolDefinition)
		log.Debug().Msgf("Comparing Protocols %s and %s", newNode.GetDefinitionMeta().Name, oldNode.GetDefinitionMeta().Name)

		result := *oldNode
		changed := false

		if len(newNode.Sequence) != len(oldNode.Sequence) {
			// Error if protocol steps are added/removed
			log.Warn().Msgf("mismatch in number of protocol steps")
			changed = true
		}

		oldSequence := make(map[string]*ProtocolStep)
		for _, f := range oldNode.Sequence {
			oldSequence[f.Name] = f
		}
		newSequence := make(map[string]*ProtocolStep)
		for i, newStep := range newNode.Sequence {
			newSequence[newStep.Name] = newStep

			if _, ok := oldSequence[newStep.Name]; !ok {
				// CHANGE: New ProtocolStep
				log.Warn().Msg("Adding new Protocol steps is not backward compatible")
				changed = true
				continue
			}

			if i > len(oldNode.Sequence) {
				// CHANGE: Reordered ProtocolSteps
				log.Warn().Msg("Reordering Protocol steps is not backward compatible")
				changed = true
				continue
			}
			if newStep.Name != oldNode.Sequence[i].Name {
				// CHANGE: Reordered/Renamed ProtocolSteps
				log.Warn().Msg("Renaming Protocol steps is not backward compatible")
				changed = true
				continue
			}
		}

		for i, oldStep := range oldNode.Sequence {
			newStep, ok := newSequence[oldStep.Name]
			if !ok {
				log.Warn().Msgf("Removing a step from a Protocol is not backward compatible")
				changed = true
				continue
			}

			var changedProtocolStep *ProtocolStep
			if ch := detectChanges(newStep, oldStep, version_index); ch != nil {
				changedProtocolStep = ch.(*ProtocolStep)
				result.Sequence[i] = changedProtocolStep
				changed = true
			}

			// Annotate the change to ProtocolStep so we can handle compatibility later in Protocol Reader/Writer
			newStep.Annotations["changes"] = append(newStep.Annotations["changes"].([]*ProtocolStep), changedProtocolStep)
		}

		if changed {
			return &result
		}
		return nil

	case *ProtocolStep:
		oldNode := oldRoot.(*ProtocolStep)
		if newNode.Name != oldNode.Name {
			// return validationError(newRoot, "cannot rename protocol step '%s' to '%s'", oldNode.Name, newNode.Name)
			return oldNode
		}

		if ch := detectChanges(newNode.Type, oldNode.Type, version_index); ch != nil {
			return oldNode
		}
		return nil

	case *RecordDefinition:
		oldNode := oldRoot.(*RecordDefinition)

		log.Debug().Msgf("Comparing RecordDefinitions %s and %s", newNode.Name, oldNode.Name)

		result := oldNode
		changed := false

		if newNode.Name != oldNode.Name {
			// CHANGE: Renamed Record
			changed = true
		}

		// Fields may be reordered
		// If they are, we want result to represent the old Record, for Serialization compatibility
		oldFieldsSlice := make([]*Field, len(oldNode.Fields))
		copy(oldFieldsSlice, oldNode.Fields)

		oldFields := make(map[string]*Field)
		for _, f := range oldNode.Fields {
			oldFields[f.Name] = f
		}
		newFields := make(map[string]*Field)
		for i, newField := range newNode.Fields {
			newFields[newField.Name] = newField

			if _, ok := oldFields[newField.Name]; !ok {
				// if !TypeHasNullOption(f.Type) {
				// 	// return validationError(f, "cannot add new field '%s'", f.Name)
				// 	panic("cannot add new field")
				// }

				// CHANGE: New field
				changed = true
				continue
			}

			if i > len(oldNode.Fields) {
				// CHANGE: Reordered fields
				changed = true
				continue
			}
			if newField.Name != oldNode.Fields[i].Name {
				// CHANGE: Reordered/Renamed fields
				changed = true
				continue
			}
		}

		for i, oldField := range oldNode.Fields {
			newField, ok := newFields[oldField.Name]
			if !ok {
				// return validationError(newRoot, "cannot add new field '%s'", f.Name)
				log.Warn().Msgf("cannot remove a field")
				changed = true
				continue
			}

			log.Debug().Msgf("Comparing fields %s and %s", newField.Name, oldField.Name)
			if ch := detectChanges(newField, oldField, version_index); ch != nil {
				result.Fields[i] = ch.(*Field)
				changed = true
				continue
			}
		}

		if changed {
			log.Warn().Msgf("Record '%s' changed", newNode.Name)
			return result
		}
		return nil

	case *Field:
		oldNode := oldRoot.(*Field)
		return compareFields(newNode, oldNode, version_index)

	case *NamedType:
		oldNode := oldRoot.(*NamedType)

		log.Debug().Msgf("Comparing NamedTypes %s and %s", newNode.Name, oldNode.Name)
		if newNode.Name != oldNode.Name {
			// CHANGE: Renamed NamedType
			return oldNode
		}

		if ch := detectChanges(newNode.Type, oldNode.Type, version_index); ch != nil {
			return oldNode
		}
		return nil

	// case *EnumDefinition:
	// TODO

	case PrimitiveDefinition:
		oldNode := oldRoot.(PrimitiveDefinition)
		if newNode != oldNode {
			// CHANGE: Changed Primitive type
			return oldNode
		}
		return nil

	case *SimpleType:
		oldNode := oldRoot.(*SimpleType)
		// TODO: Handle generalization change -> Simple to Generalized

		if newNode.Name != oldNode.Name {
			// CHANGE: Renamed SimpleType
			return oldNode
		}

		// TODO: Compare TypeArguments
		// ...

		newDef := newNode.ResolvedDefinition
		ch, ok := newDef.GetDefinitionMeta().Annotations["changes"]
		if !ok {
			if _, ok := newDef.(PrimitiveDefinition); !ok {
				panic("Expected annotation containing TypeDefinition changes")
			}
			return nil
		}

		changes := ch.([]TypeDefinition)
		if ch := changes[version_index]; ch != nil {
			log.Warn().Msgf("SimpleType '%s' changed", newNode.Name)
			return ch
		}
		return nil

	case *GeneralizedType:
		oldNode := oldRoot.(*GeneralizedType)

		changed := false

		// TODO: Compare Cases

		for i, newCase := range newNode.Cases {
			oldCase := oldNode.Cases[i]
			if ch := detectChanges(newCase, oldCase, version_index); ch != nil {
				log.Warn().Msg("GeneralizedType changed")
				changed = true
				continue
			}
		}

		// TODO: Dimensionality
		// switch oldDim := oldNode.Dimensionality.(type) {
		// }

		if changed {
			return oldNode
		}
		return nil

	case *TypeCase:
		oldNode := oldRoot.(*TypeCase)

		changed := false

		if oldNode.Type == nil {
			if newNode.Type != nil {
				log.Warn().Msg("cannot add type case type")
			}
			changed = true
		}
		if newNode.Type == nil {
			log.Warn().Msg("cannot remove type case type")
			changed = true
		}

		if newNode.Tag != oldNode.Tag {
			// TODO: ??
			changed = true
		}

		if ch := detectChanges(newNode.Type, oldNode.Type, version_index); ch != nil {
			changed = true
		}

		if changed {
			return oldNode
		}
		return nil

	default:
		return nil
	}
}

func compareFields(newField *Field, oldField *Field, version_index int) Node {
	if newField.Name != oldField.Name {
		// CHANGE: Renamed Field
		return oldField
	}

	if ch := detectChanges(newField.Type, oldField.Type, version_index); ch != nil {
		return oldField
	}
	return nil
}

type EvolutionPass func(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) (*Environment, error)

func ValidateEvolution(env *Environment, predecessor *Environment, versionId int) (*Environment, error) {
	fmt.Println("Validating evolution")

	return CollectChanges(env, predecessor, versionId)

	// errorSink := validation.ErrorSink{}

	// passes := []EvolutionPass{
	// 	// ensureNoChanges,
	// 	ensureBackwardCompatible,
	// }

	// for _, pass := range passes {
	// 	var err error
	// 	env, err = pass(env, predecessor, &errorSink)
	// 	if err != nil {
	// 		return env, err
	// 	}
	// }

	// return env, errorSink.AsError()
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
				log.Debug().Msgf("Comparing TypeDefinitions %s and %s", newTd.GetDefinitionMeta().Name, oldTd.GetDefinitionMeta().Name)
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

			log.Debug().Msgf("Comparing records %s and %s", newNode.Name, oldNode.Name)
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
