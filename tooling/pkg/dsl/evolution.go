// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

/*
	TODO: It might be better to use an interface here for a few reasons:

1. Primitive type changes only need to store the PrimitiveDefinition (not the entire Type)

2. There are various kinds of Type Changes:
  - Fully compatible (e.g. casting numbers)
  - Partially compatible (e.g. converting between scalars and unions)
  - Incompatible (e.g. converting between primitives and records)

3. Convenient to GetInverseChange (e.g. Num->String vs String->Num are both needed for Read/Write)
*/
type TypeChange struct {
	Old  Type
	New  Type
	Kind TypeChangeKind
}

type TypeChangeKind int

const (
	TypeChangeNoChange TypeChangeKind = iota
	TypeChangeNumberToNumber
	TypeChangeNumberToString
	TypeChangeStringToNumber
	TypeChangeScalarToOptional
	TypeChangeOptionalToScalar
	TypeChangeScalarToUnion
	TypeChangeUnionToScalar
	// UnionToUnion
	TypeChangeDefinitionChanged
	TypeChangeIncompatible
)

type ProtocolChange struct {
	PreviousDefinition *ProtocolDefinition
	Added              []*ProtocolStep
	Removed            []*ProtocolStep
}

type RecordChange struct {
	PreviousDefinition *RecordDefinition
	Added              []*Field
	Removed            []*Field
}

const (
	ChangeAnnotationKey             = "changed"
	AllVersionChangesAnnotationKey  = "all-changes"
	SchemaAnnotationKey             = "schema"
	AllVersionSchemasAnnotationKey  = "all-schemas"
	FieldOrStepRemovedAnnotationKey = "removed"
	VersionAnnotationKey            = "version"
)

func ValidateEvolution(env *Environment, predecessors []*Environment) (*Environment, error) {

	initializeChangeAnnotations(env)

	for versionId, predecessor := range predecessors {
		// log.Info().Msgf("Resolving changes from predecessor '%s'", predecessor.Label)
		initializePredecessorAnnotations(predecessor)

		if err := annotateAllChanges(env, predecessor); err != nil {
			return nil, err
		}

		if err := validateChanges(env); err != nil {
			return nil, err
		}

		saveChangeAnnotations(env, versionId)
	}

	return env, nil
}

func validateChanges(env *Environment) error {
	// Emit User Warnings and aggregate Errors
	errorSink := &validation.ErrorSink{}
	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {

			ch, ok := td.GetDefinitionMeta().Annotations[ChangeAnnotationKey]
			if !ok || ch == nil {
				continue
			}
			prevDef := ch.(TypeDefinition)

			if prevDef != nil {
				switch prevDef := prevDef.(type) {
				case *RecordDefinition:
					for _, field := range prevDef.Fields {
						if tc, ok := field.Annotations[ChangeAnnotationKey].(*TypeChange); ok {
							if typeChangeIsError(tc) {
								errorSink.Add(validationError(td, "Changing field '%s' from %s", field.Name, typeChangeToError(tc)))
							}

							if warn := typeChangeToWarning(tc); warn != "" {
								log.Warn().Msgf("Changing field '%s' from %s", field.Name, warn)
							}
						}
					}

				case *NamedType:
					log.Debug().Msgf("NamedType '%s' changed", td.GetDefinitionMeta().Name)
					if tc, ok := prevDef.Annotations[ChangeAnnotationKey].(*TypeChange); ok {
						if typeChangeIsError(tc) {
							errorSink.Add(validationError(td, "Changing type '%s' from %s", td.GetDefinitionMeta().Name, typeChangeToWarning(tc)))
						}
						if warn := typeChangeToWarning(tc); warn != "" {
							log.Warn().Msgf("Changing type '%s' from %s", td.GetDefinitionMeta().Name, warn)
						}
					}

				case *EnumDefinition:
					// TODO

				default:
					panic("Shouldn't get here")
				}

			}
		}

		for _, pd := range ns.Protocols {
			ch, ok := pd.GetDefinitionMeta().Annotations[ChangeAnnotationKey].(*ProtocolDefinition)
			if ok && ch != nil {
				// prevDef := ch
				for _, step := range pd.Sequence {
					tc := step.Annotations[ChangeAnnotationKey].(*TypeChange)
					if tc != nil {
						if typeChangeIsError(tc) {
							errorSink.Add(validationError(step, "Changing step '%s' from %s", step.Name, typeChangeToError(tc)))
						}
						if warn := typeChangeToWarning(tc); warn != "" {
							log.Warn().Msgf("Changing step '%s' from %s", step.Name, warn)
						}
					}
				}
			}
		}
	}

	return errorSink.AsError()
}

func typeChangeIsError(tc *TypeChange) bool {
	return tc.Kind == TypeChangeIncompatible
}

func typeChangeToError(tc *TypeChange) string {
	return fmt.Sprintf("'%s' to '%s' is not backward compatible", TypeToShortSyntax(tc.Old, true), TypeToShortSyntax(tc.New, true))
}

func typeChangeToWarning(tc *TypeChange) string {
	message := fmt.Sprintf("'%s' to '%s' ", TypeToShortSyntax(tc.Old, true), TypeToShortSyntax(tc.New, true))
	switch tc.Kind {
	case TypeChangeNumberToNumber:
		// TODO: Warn only for numeric demotion (not promotion)
		return message + "may result in loss of precision"
	case TypeChangeNumberToString, TypeChangeStringToNumber:
		return message + "may result in loss of precision"
	case TypeChangeScalarToOptional, TypeChangeOptionalToScalar:
		return message + "may result in undefined behavior"
	case TypeChangeScalarToUnion, TypeChangeUnionToScalar:
		return message + "may result in undefined behavior"
	// case UnionToUnion
	default:
		return ""
	}
}

// Prepare Annotations on the old model
func initializePredecessorAnnotations(predecessor *Environment) {
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			// node.Annotations[ChangeAnnotationKey] = nil
			node.GetDefinitionMeta().Annotations[SchemaAnnotationKey] = GetProtocolSchemaString(node, predecessor.SymbolTable)
			self.VisitChildren(node)

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			// node.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = nil
			self.VisitChildren(node)

		case *Field:
			if node.Annotations == nil {
				node.Annotations = make(map[string]any)
			}
			node.Annotations[ChangeAnnotationKey] = nil

		case *ProtocolStep:
			if node.Annotations == nil {
				node.Annotations = make(map[string]any)
			}
			node.Annotations[ChangeAnnotationKey] = nil

		default:
			self.VisitChildren(node)
		}
	})
}

// Prepare Annotations on the new model
func initializeChangeAnnotations(env *Environment) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] == nil {
				node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = make([]*ProtocolDefinition, 0)
			}
			if node.GetDefinitionMeta().Annotations[AllVersionSchemasAnnotationKey] == nil {
				node.GetDefinitionMeta().Annotations[AllVersionSchemasAnnotationKey] = make([]string, 0)
			}
			node.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = nil
			self.VisitChildren(node)

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] == nil {
				node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = make([]TypeDefinition, 0)
			}
			node.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = nil

		case *ProtocolStep:
			if node.Annotations == nil {
				node.Annotations = make(map[string]any)
			}
			if node.Annotations[AllVersionChangesAnnotationKey] == nil {
				node.Annotations[AllVersionChangesAnnotationKey] = make([]*TypeChange, 0)
			}
			node.Annotations[ChangeAnnotationKey] = nil

		default:
			self.VisitChildren(node)
		}
	})
}

func saveChangeAnnotations(env *Environment, versionId int) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			var changed *ProtocolDefinition
			var schema string
			if ch, ok := node.GetDefinitionMeta().Annotations[ChangeAnnotationKey].(*ProtocolDefinition); ok {
				changed = ch

				if s, ok := changed.GetDefinitionMeta().Annotations[SchemaAnnotationKey].(string); ok {
					schema = s
				}
			}

			node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = append(node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey].([]*ProtocolDefinition), changed)

			node.GetDefinitionMeta().Annotations[AllVersionSchemasAnnotationKey] = append(node.GetDefinitionMeta().Annotations[AllVersionSchemasAnnotationKey].([]string), schema)

			node.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = nil

			for _, step := range node.Sequence {
				var changed *TypeChange
				if ch, ok := step.Annotations[ChangeAnnotationKey].(*TypeChange); ok {
					changed = ch
				}
				step.Annotations[AllVersionChangesAnnotationKey] = append(step.Annotations[AllVersionChangesAnnotationKey].([]*TypeChange), changed)
				step.Annotations[ChangeAnnotationKey] = nil
			}

		case TypeDefinition:
			var changed TypeDefinition
			if ch, ok := node.GetDefinitionMeta().Annotations[ChangeAnnotationKey].(TypeDefinition); ok {
				changed = ch

				changed.GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionId
			}

			node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = append(node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey].([]TypeDefinition), changed)

			node.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = nil

		default:
			self.VisitChildren(node)
		}
	})
}

func annotateAllChanges(newNode, oldNode *Environment) error {
	oldNamespaces := make(map[string]*Namespace)
	for _, oldNs := range oldNode.Namespaces {
		oldNamespaces[oldNs.Name] = oldNs
	}

	for _, newNs := range newNode.Namespaces {
		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
			if err := annotateNamespaceChanges(newNs, oldNs); err != nil {
				return err
			}
		}
	}

	return nil
}

func annotateNamespaceChanges(newNode, oldNode *Namespace) error {
	if newNode.Name != oldNode.Name {
		return validationError(newNode, "changing namespaces between versions is not yet supported")
	}

	// TypeDefinitions may be reordered, added, or removed
	// We only care about pre-existing TypeDefinitions that CHANGED
	oldTds := make(map[string]TypeDefinition)
	for _, oldTd := range oldNode.TypeDefinitions {
		oldTds[oldTd.GetDefinitionMeta().Name] = oldTd
	}
	newTds := make(map[string]TypeDefinition)
	for _, newTd := range newNode.TypeDefinitions {
		newTds[newTd.GetDefinitionMeta().Name] = newTd
	}

	for _, newTd := range newNode.TypeDefinitions {
		oldTd, ok := oldTds[newTd.GetDefinitionMeta().Name]
		if !ok {
			// Skip new TypeDefinition
			continue
		}
		changedTypeDef, err := annotateChangedTypeDefinition(newTd, oldTd)
		if err != nil {
			return err
		}

		// Annotate this TypeDefinition as having changed from previous version.
		newTd.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = changedTypeDef
	}

	// Protocols may be reordered, added, or removed
	// We only care about pre-existing Protocols that CHANGED
	oldProts := make(map[string]*ProtocolDefinition)
	for _, oldProt := range oldNode.Protocols {
		oldProts[oldProt.Name] = oldProt
	}
	newProts := make(map[string]*ProtocolDefinition)
	for _, newProt := range newNode.Protocols {
		newProts[newProt.GetDefinitionMeta().Name] = newProt
	}

	for _, newProt := range newNode.Protocols {
		oldProt, ok := oldProts[newProt.GetDefinitionMeta().Name]
		if !ok {
			// Skip new ProtocolDefinition
			continue
		}

		changedProtocolDef, err := annotateChangedProtocolDefinition(newProt, oldProt)
		if err != nil {
			return err
		}

		// Annotate this ProtocolDefinition as having changed from previous version.
		newProt.GetDefinitionMeta().Annotations[ChangeAnnotationKey] = changedProtocolDef
	}

	return nil
}

// Compares two TypeDefinitions with matching names
func annotateChangedTypeDefinition(newNode, oldNode TypeDefinition) (TypeDefinition, error) {
	switch newNode := newNode.(type) {
	case *RecordDefinition:
		oldNode, ok := oldNode.(*RecordDefinition)
		if !ok {
			return oldNode, fmt.Errorf("changing '%s' to a Record is not backward compatible", newNode.Name)
		}
		res, err := annotateChangedRecordDefinition(newNode, oldNode)
		if err != nil {
			return res, err
		}
		if res != nil {
			return res, nil
		}
		return nil, nil

	case *NamedType:
		oldNode, ok := oldNode.(*NamedType)
		if !ok {
			return oldNode, fmt.Errorf("changing '%s' to a named type is not backward compatible", newNode.Name)
		}

		if ch := detectChangedTypes(newNode.Type, oldNode.Type); ch != nil {
			// CHANGE: Changed NamedType type
			// if ch.Kind == TypeChangeIncompatible {
			// 	log.Warn().Msgf("Changing '%s' from '%s' to '%s' is not backward compatible", newNode.Name, TypeToShortSyntax(oldNode.Type, true), TypeToShortSyntax(newNode.Type, true))
			// }
			oldNode.Annotations[ChangeAnnotationKey] = ch
			return oldNode, nil
		}
		return nil, nil

	case *EnumDefinition:
		oldNode, ok := oldNode.(*EnumDefinition)
		if !ok {
			return oldNode, fmt.Errorf("changing '%s' to an Enum is not backward compatible", newNode.Name)
		}
		res, err := annotateChangedEnumDefinitions(newNode, oldNode)
		if err != nil {
			return res, err
		}
		if res != nil {
			return res, nil
		}
		return nil, nil

	default:
		panic("Expected a TypeDefinition")
	}
}

// Compares two ProtocolDefinitions with matching names
func annotateChangedProtocolDefinition(newProtocol, oldProtocol *ProtocolDefinition) (*ProtocolDefinition, error) {
	changed := false

	oldSequence := make(map[string]*ProtocolStep)
	for _, f := range oldProtocol.Sequence {
		oldSequence[f.Name] = f
	}
	newSequence := make(map[string]*ProtocolStep)
	for i, newStep := range newProtocol.Sequence {
		newSequence[newStep.Name] = newStep

		if _, ok := oldSequence[newStep.Name]; !ok {
			// CHANGE: New ProtocolStep
			return oldProtocol, fmt.Errorf("adding new Protocol steps is not backward compatible")
			// log.Warn().Msg("Adding new Protocol steps is not backward compatible")
			// changed = true
			// continue
		}

		if i > len(oldProtocol.Sequence) {
			// CHANGE: Reordered ProtocolSteps
			return oldProtocol, fmt.Errorf("reordering Protocol steps is not backward compatible")
			// log.Warn().Msg("Reordering Protocol steps is not backward compatible")
			// changed = true
			// continue
		}
		if newStep.Name != oldProtocol.Sequence[i].Name {
			// CHANGE: Reordered/Renamed ProtocolSteps
			return oldProtocol, fmt.Errorf("reordering or renaming Protocol steps is not backward compatible")
			// log.Warn().Msg("Renaming Protocol steps is not backward compatible")
			// changed = true
			// continue
		}
	}

	for _, oldStep := range oldProtocol.Sequence {
		newStep, ok := newSequence[oldStep.Name]
		if !ok {
			oldStep.Annotations[FieldOrStepRemovedAnnotationKey] = true
			return oldProtocol, fmt.Errorf("removing Protocol steps is not backward compatible")
			// log.Warn().Msgf("Removing a step from a Protocol is not backward compatible")
			// changed = true
			// continue
		}

		typeChange := detectChangedTypes(newStep.Type, oldStep.Type)
		if typeChange != nil {
			changed = true

			// log.Debug().Msgf("Protocol %s step %s changed from %s to %s", newProtocol.Name, newStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			// if typeChange.Kind == TypeChangeIncompatible {
			// 	log.Warn().Msgf("Changing step '%s' from '%s' to '%s' is not backward compatible", oldStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			// }
		}

		// Annotate the change to ProtocolStep so we can handle compatibility later in Protocol Reader/Writer
		newStep.Annotations[ChangeAnnotationKey] = typeChange
	}

	if changed {
		return oldProtocol, nil
	}
	return nil, nil
}

// Compares two RecordDefinitions with matching names
func annotateChangedRecordDefinition(newRecord, oldRecord *RecordDefinition) (*RecordDefinition, error) {
	if newRecord.Name != oldRecord.Name {
		panic("Records name should match at this point")
		// CHANGE: Renamed Record
	}

	changed := false

	// Fields may be reordered
	// If they are, we want result to represent the old Record, for Serialization compatibility
	oldFields := make(map[string]*Field)
	for _, f := range oldRecord.Fields {
		oldFields[f.Name] = f
	}
	newFields := make(map[string]*Field)
	for i, newField := range newRecord.Fields {
		newFields[newField.Name] = newField

		if _, ok := oldFields[newField.Name]; !ok {
			if !TypeHasNullOption(newField.Type) {
				log.Warn().Msgf("Adding a non-Optional record field may result in undefined behavior")
			}

			// CHANGE: New field
			changed = true
			continue
		}

		if i > len(oldRecord.Fields) {
			// CHANGE: Reordered fields
			changed = true
			continue
		}
		if newField.Name != oldRecord.Fields[i].Name {
			// CHANGE: Reordered/Renamed fields
			changed = true
			continue
		}
	}

	for i, oldField := range oldRecord.Fields {
		newField, ok := newFields[oldField.Name]
		if !ok {
			if !TypeHasNullOption(oldField.Type) {
				log.Warn().Msgf("Removing a non-Optional record field may result in undefined behavior")
			}
			// CHANGE: Removed field
			oldRecord.Fields[i].Annotations[FieldOrStepRemovedAnnotationKey] = true
			changed = true
			continue
		}

		// log.Debug().Msgf("Comparing fields %s and %s", newField.Name, oldField.Name)
		if typeChange := detectChangedTypes(newField.Type, oldField.Type); typeChange != nil {
			// CHANGE: Changed field type
			changed = true
			oldRecord.Fields[i].Annotations[ChangeAnnotationKey] = typeChange
			// if typeChange.Kind == TypeChangeIncompatible {
			// 	log.Warn().Msgf("Changing field '%s' from '%s' to '%s' is not backward compatible", oldField.Name, TypeToShortSyntax(oldField.Type, true), TypeToShortSyntax(newField.Type, true))
			// }
			continue
		}
	}

	if changed {
		// log.Debug().Msgf("Record '%s' changed", newRecord.Name)
		return oldRecord, nil
	}
	// log.Debug().Msgf("Record '%s' did NOT change", newRecord.Name)
	return nil, nil
}

func annotateChangedEnumDefinitions(newNode, oldNode *EnumDefinition) (*EnumDefinition, error) {
	changed := false

	if newNode.Name != oldNode.Name {
		// CHANGE: Renamed Enum
		changed = true
	}
	if newNode.IsFlags != oldNode.IsFlags {
		// CHANGE: Changed Enum to Flags or vice versa
		changed = true
	}

	if oldNode.BaseType != nil {
		if newNode.BaseType == nil {
			// CHANGE: Removed enum base type?
			changed = true
		}
		if ch := detectChangedTypes(newNode.BaseType, oldNode.BaseType); ch != nil {
			// CHANGE: Changed Enum base type
			// log.Warn().Msgf("Changing '%s' base type is not backward compatible", newNode.Name)
			// changed = true
			return oldNode, fmt.Errorf("changing '%s' base type is not backward compatible", newNode.Name)
		}
	} else {
		if newNode.BaseType != nil {
			// CHANGE: Added an enum base type?
			changed = true
		}
	}

	for i, newEnumValue := range newNode.Values {
		oldEnumValue := oldNode.Values[i]

		if newEnumValue.Symbol != oldEnumValue.Symbol {
			// CHANGE: Renamed enum value
			changed = true
		}
		if newEnumValue.IntegerValue.Cmp(&oldEnumValue.IntegerValue) != 0 {
			// CHANGE: Changed enum value integer value
			changed = true
		}
	}

	if changed {
		return oldNode, nil
	}

	return nil, nil
}

// Compares two Types to determine what changed
// NOTE: We can't just use the `TypesEqual` function because we need to know *how* a Type changed.
func detectChangedTypes(newType, oldType Type) *TypeChange {
	// TODO: This is a good example of where it would be nice to bubble up Type Change User Warnings
	// so they are reported in the context of the `Field` or `ProtocolStep` that changed.
	//
	// UPDATE: The caller can just check if the TypeChange.Kind is TypeChangeIncompatible

	switch newType := newType.(type) {

	case *SimpleType:
		switch oldType := oldType.(type) {
		case *SimpleType:
			return detectSimpleTypeChanges(newType, oldType)
		case *GeneralizedType:
			return detectGeneralizedToSimpleTypeChanges(newType, oldType)
		default:
			panic("Shouldn't get here")
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return detectGeneralizedTypeChanges(newType, oldType)
		case *SimpleType:
			return detectSimpleToGeneralizedTypeChanges(newType, oldType)
		default:
			panic("Shouldn't get here")
		}

	default:
		panic("Expected a type")
	}
}

func detectSimpleTypeChanges(newType, oldType *SimpleType) *TypeChange {
	// TODO: Compare TypeArguments
	// This comparison depends on whether the ResolvedDefinition changed!
	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
		// CHANGE: Changed number of TypeArguments

	} else {
		for i := range newType.TypeArguments {
			if ch := detectChangedTypes(newType.TypeArguments[i], oldType.TypeArguments[i]); ch != nil {
				// CHANGE: Changed TypeArgument
				// TODO: Returning early skips other possible changes to the Type
				return ch
			}
		}
	}

	// Both newType and oldType are SimpleTypes
	// Thus, the possible type changes here are:
	//  - Primitive to Primitive
	//  - Primitive to TypeDefinition
	//  - TypeDefinition to Primitive
	//  - TypeDefinition to TypeDefinition

	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	if _, ok := newDef.(PrimitiveDefinition); ok {
		if _, ok := oldDef.(PrimitiveDefinition); ok {
			return detectPrimitiveTypeChange(newType, oldType)
		}
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	if _, ok := oldDef.(PrimitiveDefinition); ok {
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	// At this point, both Types should be TypeDefinitions
	if newDef.GetDefinitionMeta().Name != oldDef.GetDefinitionMeta().Name {
		// CHANGE: Type changed to a different TypeDefinition
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	// At this point, only the underlying TypeDefinition with matching name could have changed
	// And it would have been annotated earlier when comparing Namespace TypeDefinitions
	if ch, ok := newDef.GetDefinitionMeta().Annotations[ChangeAnnotationKey]; ok && ch != nil {
		return &TypeChange{oldType, newType, TypeChangeDefinitionChanged}
	}

	return nil
}

/*
	TODO: Leverage the type functions:

- func GetPrimitiveType(t Type) (primitive PrimitiveDefinition, ok bool)
- func GetPrimitiveKind(t PrimitiveDefinition) PrimitiveKind
- func GetKindIfPrimitive(t Type) (primitiveKind PrimitiveKind, ok bool)
- func IsIntegralPrimitive(prim PrimitiveDefinition)
- func IsIntegralType(t Type) bool
*/
func detectPrimitiveTypeChange(newType, oldType *SimpleType) *TypeChange {
	newPrimitive := newType.ResolvedDefinition.(PrimitiveDefinition)
	oldPrimitive := oldType.ResolvedDefinition.(PrimitiveDefinition)

	if newPrimitive == oldPrimitive {
		return nil
	}

	// CHANGE: Changed Primitive type
	if oldPrimitive == PrimitiveString {
		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChange{oldType, newType, TypeChangeStringToNumber}
		}
	}

	if GetPrimitiveKind(oldPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(oldPrimitive) == PrimitiveKindFloatingPoint {
		if newPrimitive == PrimitiveString {
			return &TypeChange{oldType, newType, TypeChangeNumberToString}
		}

		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChange{oldType, newType, TypeChangeNumberToNumber}
		}
	}

	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType) *TypeChange {
	// A GeneralizedType can change in many ways...
	changeKind := TypeChangeNoChange

	// Compare Dimensonality
	if ch := detectChangedDimensionality(newType.Dimensionality, oldType.Dimensionality); ch != nil {
		// CHANGE: Dimensionality changed
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	if len(newType.Cases) != len(oldType.Cases) {
		// TODO: Handle adding types to a Union/Optional
		// TODO: Handle removing types from a Union/Optional

		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	// Compare the Type of each TypeCase
	for i, newCase := range newType.Cases {
		oldCase := oldType.Cases[i]

		if (newCase.Type == nil) != (oldCase.Type == nil) {
			// CHANGE: Added or removed a type case
			return &TypeChange{oldType, newType, TypeChangeIncompatible}
		}

		if newCase.Type == nil {
			continue
		}

		if ch := detectChangedTypes(newCase.Type, oldCase.Type); ch != nil {
			// CHANGE: Changed a type case
			changeKind = ch.Kind
			if ch.Kind == TypeChangeIncompatible {
				break
			}
		}
	}

	if changeKind != TypeChangeNoChange {
		return &TypeChange{oldType, newType, changeKind}
	}
	return nil
}

func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType) *TypeChange {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		if TypesEqual(newType, oldType.Cases[1].Type) {
			return &TypeChange{oldType, newType, TypeChangeOptionalToScalar}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		compatible := false
		for _, tc := range oldType.Cases {
			if TypesEqual(newType, tc.Type) {
				compatible = true
			}
		}
		// TODO: Need to capture the index of the matching type in the old Union's Cases
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeUnionToScalar}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType) *TypeChange {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		if TypesEqual(newType.Cases[1].Type, oldType) {
			return &TypeChange{oldType, newType, TypeChangeScalarToOptional}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		compatible := false
		for _, tc := range newType.Cases {
			if TypesEqual(tc.Type, oldType) {
				compatible = true
			}
		}
		// TODO: Need to capture the index of the matching type in the new Union's Cases
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeScalarToUnion}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectChangedDimensionality(newNode, oldNode Dimensionality) error {
	switch newDim := newNode.(type) {
	case nil:
		if oldNode != nil {
			// CHANGE: Removed dimensionality
			return fmt.Errorf("removing dimensionality is not backward compatible")
		}
		return nil

	case *Stream:
		_, ok := oldNode.(*Stream)
		if !ok {
			return fmt.Errorf("expected a stream")
		}
		return nil

	case *Vector:
		oldDim, ok := oldNode.(*Vector)
		if !ok {
			return fmt.Errorf("expected a vector")
		}
		if (oldDim.Length == nil) != (newDim.Length == nil) {
			// CHANGE: Added or removed vector length
			return fmt.Errorf("changing vector length is not backward compatible")
		}
		if newDim.Length != nil && *newDim.Length != *oldDim.Length {
			// CHANGE: Changed vector length
			return fmt.Errorf("changing vector length is not backward compatible")
		}
		return nil

	case *Array:
		oldDim, ok := oldNode.(*Array)
		if !ok {
			return fmt.Errorf("expected an array")
		}
		if (newDim.Dimensions == nil) != (oldDim.Dimensions == nil) {
			// CHANGE: Added or removed array dimensions
			return fmt.Errorf("changing array dimensions is not backward compatible")
		}

		if newDim.Dimensions != nil {
			newDimensions := *newDim.Dimensions
			oldDimensions := *oldDim.Dimensions

			if len(newDimensions) != len(oldDimensions) {
				// CHANGE: Mismatch in number of array dimensions
				return fmt.Errorf("changing array dimensions is not backward compatible")
			}

			for i, newDimension := range newDimensions {
				oldDimension := oldDimensions[i]

				if (newDimension.Length == nil) != (oldDimension.Length == nil) {
					// CHANGE: Added or removed array dimension length
					return fmt.Errorf("changing array dimensions is not backward compatible")
				}

				if newDimension.Length != nil && newDimension.Length != oldDimension.Length {
					// CHANGE: Changed array dimension length
					return fmt.Errorf("changing array dimensions is not backward compatible")
				}
			}
		}

		return nil

	case *Map:
		oldDim, ok := oldNode.(*Map)
		if !ok {
			return fmt.Errorf("expected a map")
		}
		if ch := detectChangedTypes(newDim.KeyType, oldDim.KeyType); ch != nil {
			return fmt.Errorf("changing map key type is not backward compatible")
		}
		return nil

	default:
		log.Panic().Msgf("unhandled type %T", newNode)
	}

	return nil
}
