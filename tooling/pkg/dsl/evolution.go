// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

const (
	// Annotations referenced in serialization codegen
	VersionAnnotationKey = "version"

	// Annotations used only for validation model evolution (local to this file)
	schemaAnnotationKey = "schema"
)

type ChangeTable map[string]DefinitionChange

func ValidateEvolution(env *Environment, predecessors []*Environment, versionLabels []string) (*Environment, error) {

	initializeChangeAnnotations(env)

	changeTable := make(ChangeTable)

	for i, predecessor := range predecessors {
		log.Info().Msgf("Resolving changes from predecessor %s", versionLabels[i])
		annotatePredecessorSchemas(predecessor)

		if err := annotateAllChanges(env, predecessor, changeTable, versionLabels[i]); err != nil {
			return nil, err
		}

		if err := validateChanges(env, changeTable); err != nil {
			return nil, err
		}

		saveChangedDefinitions(env, changeTable, versionLabels[i])
	}

	return env, nil
}

func validateChanges(env *Environment, changeTable ChangeTable) error {
	// Emit User Warnings and aggregate Errors
	errorSink := &validation.ErrorSink{}
	for _, ns := range env.Namespaces {

		validateTypeDefinitionChanges(ns.TypeDefinitions, changeTable, errorSink)
		validateProtocolChanges(ns.Protocols, changeTable, errorSink)
	}
	return errorSink.AsError()
}

func validateTypeDefinitionChanges(typeDefs []TypeDefinition, changeTable ChangeTable, errorSink *validation.ErrorSink) {
	for _, td := range typeDefs {
		defChange, ok := changeTable[td.GetDefinitionMeta().GetQualifiedName()]
		if !ok || defChange == nil {
			continue
		}

		switch defChange := defChange.(type) {

		case *DefinitionChangeIncompatible:
			errorSink.Add(validationError(td, "changing '%s' is not backward compatible", td.GetDefinitionMeta().Name))

		case *RecordChange:
			oldRec := defChange.PreviousDefinition().(*RecordDefinition)
			newRec := defChange.LatestDefinition().(*RecordDefinition)

			for _, added := range defChange.FieldsAdded {
				if !TypeHasNullOption(added.Type) {
					log.Warn().Msgf("Adding non-Optional field '%s' may result in undefined behavior with previous versions", added.Name)
				}
			}

			for i, field := range oldRec.Fields {
				if defChange.FieldRemoved[i] {
					if !TypeHasNullOption(oldRec.Fields[i].Type) {
						log.Warn().Msgf("Removing non-Optional field '%s' may result in undefined behavior with previous versions", field.Name)
					}
					continue
				}

				if tc := defChange.FieldChanges[i]; tc != nil {
					if typeChangeIsError(tc) {
						newField := newRec.Fields[defChange.NewFieldIndex[i]]
						errorSink.Add(validationError(newField, "changing field '%s' from %s", newField.Name, typeChangeToError(tc)))
					}

					if warn := typeChangeToWarning(tc); warn != "" {
						log.Warn().Msgf("Changing field '%s' from %s", field.Name, warn)
					}
				}
			}

		case *NamedTypeChange:
			if tc := defChange.TypeChange; tc != nil {
				if typeChangeIsError(tc) {
					errorSink.Add(validationError(td, "changing type '%s' from %s", td.GetDefinitionMeta().Name, typeChangeToWarning(tc)))
				}
				if warn := typeChangeToWarning(tc); warn != "" {
					log.Warn().Msgf("Changing type '%s' from %s", td.GetDefinitionMeta().Name, warn)
				}
			}

		case *EnumChange:
			if tc := defChange.BaseTypeChange; tc != nil {
				errorSink.Add(validationError(td, "changing base type of '%s' is not backward compatible", td.GetDefinitionMeta().Name))
			}

		default:
			panic("Shouldn't get here")
		}
	}
}

func validateProtocolChanges(protocols []*ProtocolDefinition, changeTable ChangeTable, errorSink *validation.ErrorSink) {
	for _, pd := range protocols {
		protChange, ok := changeTable[pd.GetQualifiedName()].(*ProtocolChange)
		if !ok || protChange == nil {
			continue
		}

		for _, reordered := range protChange.StepsReordered {
			errorSink.Add(validationError(reordered, "reordering step '%s' is not backward compatible", reordered.Name))
		}

		for _, removed := range protChange.StepsRemoved {
			errorSink.Add(validationError(pd, "removing step '%s' is not backward compatible", removed.Name))
		}

		for i, step := range pd.Sequence {
			if tc := protChange.StepChanges[i]; tc != nil {
				switch tc := tc.(type) {
				case *TypeChangeStepAdded:
					// A Step can be added to a Protocol if its Type can have an "empty" state
					typeCanBeEmpty := false
					switch t := GetUnderlyingType(step.Type).(type) {
					case *GeneralizedType:
						if t.Cases.HasNullOption() {
							typeCanBeEmpty = true
						} else if t.Dimensionality != nil {
							switch t.Dimensionality.(type) {
							case *Stream, *Vector, *Map:
								typeCanBeEmpty = true
							}
						}
					}
					if !typeCanBeEmpty {
						errorSink.Add(validationError(step, "adding step '%s' is not backward compatible", step.Name))
					}
				default:
					if typeChangeIsError(tc) {
						errorSink.Add(validationError(step.Type, "changing step '%s' from %s", step.Name, typeChangeToError(tc)))
					}

					if warn := typeChangeToWarning(tc); warn != "" {
						log.Warn().Msgf("Changing step '%s' from %s", step.Name, warn)
					}
				}
			}
		}
	}
}

// Annotate the previous model with Protocol Schema strings for later
func annotatePredecessorSchemas(predecessor *Environment) {
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations[schemaAnnotationKey] = GetProtocolSchemaString(node, predecessor.SymbolTable)

		default:
			self.VisitChildren(node)
		}
	})
}

// Prepare Annotations on the new model
func initializeChangeAnnotations(env *Environment) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *Namespace:
			node.TypeDefChanges = make(map[string][]DefinitionChange)
			self.VisitChildren(node)

		case *ProtocolDefinition:
			node.Versions = make(map[string]*ProtocolChange)

		default:
			self.VisitChildren(node)
		}
	})
}

func saveChangedDefinitions(env *Environment, changeTable ChangeTable, versionLabel string) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *Namespace:
			node.Versions = append(node.Versions, versionLabel)

			for _, change := range node.TypeDefChanges[versionLabel] {
				if change.PreviousDefinition().GetDefinitionMeta().Annotations == nil {
					change.PreviousDefinition().GetDefinitionMeta().Annotations = make(map[string]any)
				}
				change.PreviousDefinition().GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionLabel
			}

			self.VisitChildren(node)

		case *ProtocolDefinition:
			var changed *ProtocolChange
			if ch, ok := changeTable[node.GetQualifiedName()].(*ProtocolChange); ok {
				changed = ch
			}
			node.Versions[versionLabel] = changed

		default:
			self.VisitChildren(node)
		}
	})
}

func annotateAllChanges(newNode, oldNode *Environment, changeTable ChangeTable, versionLabel string) error {
	oldNamespaces := make(map[string]*Namespace)
	for _, oldNs := range oldNode.Namespaces {
		oldNamespaces[oldNs.Name] = oldNs
	}

	for _, newNs := range newNode.Namespaces {
		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
			annotateNamespaceChanges(newNs, oldNs, changeTable, versionLabel)
		} else {
			return fmt.Errorf("Namespace '%s' does not exist in previous version", newNs.Name)
		}
	}

	return nil
}

func annotateNamespaceChanges(newNs, oldNs *Namespace, changeTable ChangeTable, versionLabel string) {
	// TypeDefinitions may be reordered, added, or removed, so we compare them by name
	oldTds := make(map[string]TypeDefinition)
	for _, oldTd := range oldNs.TypeDefinitions {
		oldTds[oldTd.GetDefinitionMeta().GetQualifiedName()] = oldTd
	}

	isUserTypeDef := make(map[string]bool)
	for _, newTd := range newNs.TypeDefinitions {
		isUserTypeDef[newTd.GetDefinitionMeta().GetQualifiedName()] = true
	}

	// Collect only pre-existing TypeDefinitions that are used within a Protocol
	// Keeping them in definition order!
	typeDefCollected := make(map[string]bool)
	var newUsedTypeDefs []TypeDefinition
	for _, protocol := range newNs.Protocols {
		Visit(protocol, func(self Visitor, node Node) {
			switch node := node.(type) {
			case TypeDefinition:
				self.VisitChildren(node)
				name := node.GetDefinitionMeta().GetQualifiedName()
				if isUserTypeDef[name] && !typeDefCollected[name] {
					typeDefCollected[name] = true
					newUsedTypeDefs = append(newUsedTypeDefs, node)
				}
			case *SimpleType:
				self.VisitChildren(node)
				self.Visit(node.ResolvedDefinition)
			default:
				self.VisitChildren(node)
			}
		})
	}

	typeDefChanges := make([]DefinitionChange, 0)
	alreadyCompared := make(map[string]bool)
	for _, newTd := range newUsedTypeDefs {
		oldTd, ok := oldTds[newTd.GetDefinitionMeta().GetQualifiedName()]
		if !ok {
			// Skip new TypeDefinition
			continue
		}

		type NamedTypeUnwinder = func(TypeDefinition) TypeDefinition
		removedAliases := make([]DefinitionChange, 0)
		var unwindOldAlias, unwindNewAlias NamedTypeUnwinder

		unwindOldAlias = func(oldTd TypeDefinition) TypeDefinition {
			switch old := oldTd.(type) {
			case *NamedType:
				if _, isNamedType := newTd.(*NamedType); !isNamedType {
					// Alias removed and we need to generate its compatibility serializers.
					if oldType, ok := old.Type.(*SimpleType); ok {
						compat := &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
						removedAliases = append([]DefinitionChange{compat}, removedAliases...)
						oldTd = oldType.ResolvedDefinition
						return unwindOldAlias(oldTd)
					}
				}
			}
			return oldTd
		}

		unwindNewAlias = func(newTd TypeDefinition) TypeDefinition {
			switch new := newTd.(type) {
			case *NamedType:
				if _, isNamedType := oldTd.(*NamedType); !isNamedType {
					// Alias added and we can safely ignore it.
					if newType, ok := new.Type.(*SimpleType); ok {
						newTd = newType.ResolvedDefinition
						return unwindNewAlias(newTd)
					}
				}
			}
			return newTd
		}

		// TODO: The unwind helpers might modify oldTd/newTd out-of-order. Make them NOT capture locals

		// "Unwind" any NamedTypes so we only compare underlying TypeDefinitions
		oldTd = unwindOldAlias(oldTd)
		newTd = unwindNewAlias(newTd)

		if alreadyCompared[newTd.GetDefinitionMeta().GetQualifiedName()] {
			// TODO: Remove this check if not needed, once integration tests are "complete"
			panic(fmt.Sprintf("Already Compared %s", newTd.GetDefinitionMeta().GetQualifiedName()))
			continue
		}

		defChange := detectTypeDefinitionChanges(newTd, oldTd, changeTable)
		if defChange != nil {
			typeDefChanges = append(typeDefChanges, defChange)
			typeDefChanges = append(typeDefChanges, removedAliases...)

			// Save this DefinitionChange so that, later, detectSimpleTypeChanges can determine if an underlying TypeDefinition changed
			changeTable[newTd.GetDefinitionMeta().GetQualifiedName()] = defChange
		}

		alreadyCompared[newTd.GetDefinitionMeta().GetQualifiedName()] = true
	}

	// Save all TypeDefinition changes for generating of compatibility serializers
	newNs.TypeDefChanges[versionLabel] = typeDefChanges

	// Protocols may be reordered, added, or removed
	// We only care about pre-existing Protocols that CHANGED
	oldProts := make(map[string]*ProtocolDefinition)
	for _, oldProt := range oldNs.Protocols {
		oldProts[oldProt.GetQualifiedName()] = oldProt
	}

	for _, newProt := range newNs.Protocols {
		oldProt, ok := oldProts[newProt.GetDefinitionMeta().GetQualifiedName()]
		if !ok {
			// Skip new ProtocolDefinition
			continue
		}

		// Annotate this ProtocolDefinition with any changes from previous version.
		protocolChange := detectProtocolDefinitionChanges(newProt, oldProt, changeTable)
		changeTable[newProt.GetQualifiedName()] = protocolChange
	}
}

// Compares two TypeDefinitions with matching names
func detectTypeDefinitionChanges(newTd, oldTd TypeDefinition, changeTable ChangeTable) DefinitionChange {
	switch newNode := newTd.(type) {
	case *RecordDefinition:
		switch oldTd := oldTd.(type) {
		case *RecordDefinition:
			if ch := detectRecordDefinitionChanges(newNode, oldTd, changeTable); ch != nil {
				return ch
			}
			return nil

		default:
			// Changing a non-Record to a Record is not backward compatible
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	case *NamedType:
		switch oldTd := oldTd.(type) {
		case *NamedType:
			if typeChange := detectTypeChanges(newNode.Type, oldTd.Type, changeTable); typeChange != nil {
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, typeChange}
			}
			return nil

		default:
			panic("Shouldn't get here")
		}

	case *EnumDefinition:
		oldTd, ok := oldTd.(*EnumDefinition)
		if !ok {
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}
		if ch := detectEnumDefinitionChanges(newNode, oldTd, changeTable); ch != nil {
			return ch
		}
		return nil

	default:
		// log.Debug().Msgf("What is this? %s was %s", newNode, oldTd)
		panic("Expected a TypeDefinition")
	}
}

// Compares two ProtocolDefinitions with matching names
func detectProtocolDefinitionChanges(newProtocol, oldProtocol *ProtocolDefinition, changeTable ChangeTable) *ProtocolChange {
	change := &ProtocolChange{
		DefinitionPair: DefinitionPair{oldProtocol, newProtocol},
		PreviousSchema: oldProtocol.GetDefinitionMeta().Annotations[schemaAnnotationKey].(string),
		StepChanges:    make([]TypeChange, len(newProtocol.Sequence)),
	}

	newSteps := make(map[string]*ProtocolStep)
	for _, newStep := range newProtocol.Sequence {
		newSteps[newStep.Name] = newStep
	}

	oldSteps := make(map[string]*ProtocolStep)
	oldStepIndices := make(map[string]int)
	for i, oldStep := range oldProtocol.Sequence {
		oldSteps[oldStep.Name] = oldStep
		oldStepIndices[oldStep.Name] = i

		if _, ok := newSteps[oldStep.Name]; !ok {
			// CHANGE: Removed this ProtocolStep
			change.StepsRemoved = append(change.StepsRemoved, oldStep)
		}
	}

	expectedIndex := 0
	for i, newStep := range newProtocol.Sequence {
		oldStep, ok := oldSteps[newStep.Name]
		if !ok {
			// CHANGE: Added this ProtocolStep
			change.StepChanges[i] = &TypeChangeStepAdded{TypePair{nil, newStep.Type}}
			continue
		}

		if oldStepIndices[newStep.Name] != expectedIndex {
			// CHANGE: Reordered this ProtocolStep
			change.StepsReordered = append(change.StepsReordered, newStep)
		}
		expectedIndex++

		if typeChange := detectTypeChanges(newStep.Type, oldStep.Type, changeTable); typeChange != nil {
			// CHANGE: ProtocolStep type changed
			change.StepChanges[i] = typeChange
		}
	}

	if len(change.StepsReordered) > 0 || len(change.StepsRemoved) > 0 {
		return change
	}
	for _, ch := range change.StepChanges {
		if ch != nil {
			return change
		}
	}

	return nil
}

// Compares two RecordDefinitions with matching names
func detectRecordDefinitionChanges(newRecord, oldRecord *RecordDefinition, changeTable ChangeTable) *RecordChange {
	change := &RecordChange{
		DefinitionPair: DefinitionPair{oldRecord, newRecord},
		FieldRemoved:   make([]bool, len(oldRecord.Fields)),
		FieldChanges:   make([]TypeChange, len(oldRecord.Fields)),
		NewFieldIndex:  make([]int, len(oldRecord.Fields)),
	}

	// Fields may be reordered
	// If they are, we want result to represent the old Record, for Serialization compatibility
	oldFields := make(map[string]*Field)
	for _, f := range oldRecord.Fields {
		oldFields[f.Name] = f
	}

	newFields := make(map[string]*Field)
	newFieldIndices := make(map[string]int)
	for i, newField := range newRecord.Fields {
		newFields[newField.Name] = newField
		newFieldIndices[newField.Name] = i

		if _, ok := oldFields[newField.Name]; !ok {
			// CHANGE: New field
			change.FieldsAdded = append(change.FieldsAdded, newField)
		}
	}

	fieldsReordered := false
	for i, oldField := range oldRecord.Fields {
		newField, ok := newFields[oldField.Name]
		if !ok {
			// CHANGE: Removed field
			change.FieldRemoved[i] = true
			change.NewFieldIndex[i] = -1
			continue
		}

		change.NewFieldIndex[i] = newFieldIndices[oldField.Name]
		if change.NewFieldIndex[i] != i {
			fieldsReordered = true
		}

		if typeChange := detectTypeChanges(newField.Type, oldField.Type, changeTable); typeChange != nil {
			// CHANGE: Field type changed
			change.FieldChanges[i] = typeChange
		}
	}

	if fieldsReordered || len(change.FieldsAdded) > 0 {
		return change
	}
	for i := range oldRecord.Fields {
		if change.FieldRemoved[i] || change.FieldChanges[i] != nil {
			return change
		}
	}
	return nil
}

func detectEnumDefinitionChanges(newNode, oldEnum *EnumDefinition, changeTable ChangeTable) DefinitionChange {
	if newNode.IsFlags != oldEnum.IsFlags {
		// CHANGE: Changed Enum to Flags or vice versa
		return &DefinitionChangeIncompatible{DefinitionPair{oldEnum, newNode}}
	}

	oldBaseType := oldEnum.BaseType
	if oldBaseType == nil {
		oldBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	newBaseType := newNode.BaseType
	if newBaseType == nil {
		newBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	if ch := detectTypeChanges(newBaseType, oldBaseType, changeTable); ch != nil {
		// CHANGE: Changed Enum base type
		return &EnumChange{DefinitionPair{oldEnum, newNode}, ch}
	}

	return nil
}

// Compares two Types to determine if and how they changed
func detectTypeChanges(newType, oldType Type, changeTable ChangeTable) TypeChange {
	newType = GetUnderlyingType(newType)
	oldType = GetUnderlyingType(oldType)

	switch newType := newType.(type) {

	case *SimpleType:
		switch oldType := oldType.(type) {
		case *SimpleType:
			return detectSimpleTypeChanges(newType, oldType, changeTable)
		case *GeneralizedType:
			return detectGeneralizedToSimpleTypeChanges(newType, oldType, changeTable)
		default:
			panic("Shouldn't get here")
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return detectGeneralizedTypeChanges(newType, oldType, changeTable)
		case *SimpleType:
			return detectSimpleToGeneralizedTypeChanges(newType, oldType, changeTable)
		default:
			panic("Shouldn't get here")
		}

	default:
		panic("Expected a type")
	}
}

func detectSimpleTypeChanges(newType, oldType *SimpleType, changeTable ChangeTable) TypeChange {
	// TODO: Compare TypeArguments
	// This comparison depends on whether the ResolvedDefinition changed!
	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
		// CHANGE: Changed number of TypeArguments
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	} else {
		for i := range newType.TypeArguments {
			if ch := detectTypeChanges(newType.TypeArguments[i], oldType.TypeArguments[i], changeTable); ch != nil {
				// CHANGE: Changed TypeArgument
				// TODO: Returning early skips other possible changes to the Type
				return ch
			}
		}
	}

	// Both newType and oldType are SimpleTypes
	// Thus, the possible type changes here are:
	//  - Primitive to Primitive (possibly valid)
	//  - TypeDefinition to TypeDefinition (possibly valid)
	//  - Primitive to TypeDefinition (invalid)
	//  - TypeDefinition to Primitive (invalid)

	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	if _, ok := oldDef.(PrimitiveDefinition); ok {
		if _, ok := newDef.(PrimitiveDefinition); ok {
			return detectPrimitiveTypeChange(newType, oldType)
		}
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if _, ok := newDef.(PrimitiveDefinition); ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch, ok := changeTable[newDef.GetDefinitionMeta().GetQualifiedName()]; ok {
		if ch != nil && ch.PreviousDefinition() == oldDef {
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}}
		}
	}

	if newDef.GetDefinitionMeta().GetQualifiedName() != oldDef.GetDefinitionMeta().GetQualifiedName() {
		// CHANGE: Not the same underlying TypeDefinition
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}

func detectPrimitiveTypeChange(newType, oldType *SimpleType) TypeChange {
	newPrimitive := newType.ResolvedDefinition.(PrimitiveDefinition)
	oldPrimitive := oldType.ResolvedDefinition.(PrimitiveDefinition)

	if newPrimitive == oldPrimitive {
		return nil
	}

	// CHANGE: Changed Primitive type
	if oldPrimitive == PrimitiveString {
		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChangeStringToNumber{TypePair{oldType, newType}}
		}
	}

	if GetPrimitiveKind(oldPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(oldPrimitive) == PrimitiveKindFloatingPoint {
		if newPrimitive == PrimitiveString {
			return &TypeChangeNumberToString{TypePair{oldType, newType}}
		}

		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChangeNumberToNumber{TypePair{oldType, newType}}
		}
	}

	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		switch detectTypeChanges(newType, oldType.Cases[1].Type, changeTable).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeOptionalToScalar{TypePair{oldType, newType}}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		for i, tc := range oldType.Cases {
			switch detectTypeChanges(newType, tc.Type, changeTable).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeUnionToScalar{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType, changeTable ChangeTable) TypeChange {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		switch detectTypeChanges(newType.Cases[1].Type, oldType, changeTable).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeScalarToOptional{TypePair{oldType, newType}}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		for i, tc := range newType.Cases {
			switch detectTypeChanges(tc.Type, oldType, changeTable).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeScalarToUnion{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	// A GeneralizedType can change in many ways...
	if newType.Cases.IsOptional() {
		return detectOptionalChanges(newType, oldType, changeTable)
	} else if newType.Cases.IsUnion() {
		return detectUnionChanges(newType, oldType, changeTable)
	} else {
		switch newType.Dimensionality.(type) {
		case nil:
			// TODO: Not an Optional, Union, Stream, Vector, Array, Map...
		case *Stream:
			return detectStreamChanges(newType, oldType, changeTable)
		case *Vector:
			return detectVectorChanges(newType, oldType, changeTable)
		case *Array:
			return detectArrayChanges(newType, oldType, changeTable)
		case *Map:
			return detectMapChanges(newType, oldType, changeTable)
		default:
			panic("Shouldn't get here")
		}
	}

	return nil
}

func detectOptionalChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	if !oldType.Cases.IsOptional() {
		if oldType.Cases.IsUnion() && oldType.Cases.HasNullOption() {
			// An Optional<T> can become a Union<null, T, ...> ONLY if
			// 	1. type T does not change, or
			// 	2. type T's TypeDefinition changed

			// Look for a matching type in the old Union
			for i, c := range oldType.Cases[1:] {
				switch detectTypeChanges(newType.Cases[1].Type, c.Type, changeTable).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, i + 1}
				}
			}
		}

		// CHANGE: Changed a non-Optional/Union to an Optional
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[1].Type, oldType.Cases[1].Type, changeTable); ch != nil {
		// CHANGE: Changed Optional type
		return &TypeChangeOptionalTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectUnionChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	if !oldType.Cases.IsUnion() {
		if oldType.Cases.IsOptional() && newType.Cases.HasNullOption() {
			for i, c := range newType.Cases[1:] {
				switch detectTypeChanges(c.Type, oldType.Cases[1].Type, changeTable).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeOptionalToUnion{TypePair{oldType, newType}, i + 1}
				}
			}
		}
		// CHANGE: Changed a non-Union/Optional to a Union
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	oldMatches := make([]bool, len(oldType.Cases))
	newMatches := make([]bool, len(newType.Cases))

	innerTypeDefsChanged := false
	// Search for a match for each Type in the new Union
	for i, newCase := range newType.Cases {
		for j, oldCase := range oldType.Cases {
			if oldMatches[j] {
				continue
			}

			switch detectTypeChanges(newCase.Type, oldCase.Type, changeTable).(type) {
			case nil:
				// Found matching type
				newMatches[i] = true
				oldMatches[j] = true
			case *TypeChangeDefinitionChanged:
				// Found matching type with underlying definition changed
				newMatches[i] = true
				oldMatches[j] = true
				innerTypeDefsChanged = true
			}
		}
	}

	// If newMatches is all False, then this isn't a valid Union type change
	// If newMatches is not all True, then type(s) were added to the Union
	// If oldMatches is not all True, then type(s) were removed from the Union
	// If newMatches and oldMatches are all true, then the Union types are the same, but possibly reordered
	anyMatch := false
	allMatch := true
	for _, m := range newMatches {
		if !m {
			allMatch = false
		} else {
			anyMatch = true
		}
	}
	for _, m := range oldMatches {
		if !m {
			allMatch = false
		}
	}

	if !anyMatch {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if innerTypeDefsChanged || !allMatch {
		return &TypeChangeUnionTypesetChanged{TypePair{oldType, newType}, oldMatches, newMatches}
	}

	return nil
}

func detectStreamChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	if _, ok := oldType.Dimensionality.(*Stream); !ok {
		// CHANGE: Changed a non-Stream to a Stream
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
		// CHANGE: Changed Stream type
		return &TypeChangeStreamTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectVectorChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	newDim := newType.Dimensionality.(*Vector)
	oldDim, ok := oldType.Dimensionality.(*Vector)
	if !ok {
		// CHANGE: Changed a non-Vector to a Vector
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if (oldDim.Length == nil) != (newDim.Length == nil) {
		// CHANGE: Changed from a fixed-length Vector to a variable-length Vector or vice versa
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if newDim.Length != nil && *newDim.Length != *oldDim.Length {
		// CHANGE: Changed vector length
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
		// CHANGE: Changed Vector type
		return &TypeChangeVectorTypeChanged{TypePair{oldType, newType}, ch}
	}

	return nil
}

func detectArrayChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	newDim := newType.Dimensionality.(*Array)
	oldDim, ok := oldType.Dimensionality.(*Array)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
		// CHANGE: Changed Array type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if (newDim.Dimensions == nil) != (oldDim.Dimensions == nil) {
		// CHANGE: Added or removed array dimensions
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if newDim.Dimensions != nil {
		newDimensions := *newDim.Dimensions
		oldDimensions := *oldDim.Dimensions

		if len(newDimensions) != len(oldDimensions) {
			// CHANGE: Mismatch in number of array dimensions
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		for i, newDimension := range newDimensions {
			oldDimension := oldDimensions[i]

			if (newDimension.Length == nil) != (oldDimension.Length == nil) {
				// CHANGE: Added or removed array dimension length
				return &TypeChangeIncompatible{TypePair{oldType, newType}}
			}

			if newDimension.Length != nil && *newDimension.Length != *oldDimension.Length {
				// CHANGE: Changed array dimension length
				return &TypeChangeIncompatible{TypePair{oldType, newType}}
			}
		}
	}
	return nil
}

func detectMapChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
	newDim := newType.Dimensionality.(*Map)
	oldDim, ok := oldType.Dimensionality.(*Map)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newDim.KeyType, oldDim.KeyType, changeTable); ch != nil {
		// CHANGE: Changed Map key type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
		// CHANGE: Changed Map value type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}
