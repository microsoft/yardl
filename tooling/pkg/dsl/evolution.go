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
	// TypeChangeIncompatible
)

func TypeChangeIncompatibleError(oldType, newType Type) error {
	return fmt.Errorf("changing '%s' to '%s' is not backward compatible", TypeToShortSyntax(oldType, true), TypeToShortSyntax(newType, true))
}

func ValidateEvolution(env *Environment, predecessor *Environment, versionId int) (*Environment, error) {

	preprocessPredecessor(predecessor)
	preprocessCurrent(env)

	oldNamespaces := make(map[string]*Namespace)
	for _, oldNs := range predecessor.Namespaces {
		oldNamespaces[oldNs.Name] = oldNs
	}

	for _, newNs := range env.Namespaces {
		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
			if err := annotateNamespaceChanges(newNs, oldNs, versionId); err != nil {
				return env, err
			}
		}
	}

	// Emit User Warnings and aggregate Errors
	errorSink := &validation.ErrorSink{}
	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {
			changes := td.GetDefinitionMeta().Annotations["changes"].([]TypeDefinition)
			if versionId >= len(changes) {
				continue
			}

			prevDef := changes[versionId]
			if prevDef != nil {
				switch prevDef := prevDef.(type) {
				case *RecordDefinition:
					for _, field := range prevDef.Fields {
						if tc, ok := field.Annotations["changed"].(*TypeChange); ok {
							if typeChangeIsError(tc) {
								errorSink.Add(validationError(td, "Changing field '%s' from %s", field.Name, typeChangeToWarning(tc)))
							}

							if warn := typeChangeToWarning(tc); warn != "" {
								log.Warn().Msgf("Changing field '%s' from %s", field.Name, warn)
							}
						}
					}

				case *NamedType:
					if tc, ok := prevDef.Annotations["changed"].(*TypeChange); ok {
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
			changes := pd.GetDefinitionMeta().Annotations["changes"].([]*ProtocolDefinition)
			if versionId >= len(changes) {
				continue
			}

			prevDef := changes[versionId]
			if prevDef != nil {
				for _, step := range pd.Sequence {
					tc := step.Annotations["changes"].([]*TypeChange)[versionId]
					if tc != nil {
						if typeChangeIsError(tc) {
							errorSink.Add(validationError(step, "Changing step '%s' from %s", step.Name, typeChangeToWarning(tc)))
						}
						if warn := typeChangeToWarning(tc); warn != "" {
							log.Warn().Msgf("Changing step '%s' from %s", step.Name, warn)
						}
					}
				}
			}
		}
	}

	return env, errorSink.AsError()
}

func typeChangeIsError(tc *TypeChange) bool {
	return false
}

func typeChangeToWarning(tc *TypeChange) string {
	message := fmt.Sprintf("'%s' to '%s' ", TypeToShortSyntax(tc.Old, true), TypeToShortSyntax(tc.New, true))
	switch tc.Kind {
	case TypeChangeNumberToNumber:
		// TODO: Warning only for numeric demotion (not promotion)
		return message + "may result in loss of precision"
	case TypeChangeNumberToString, TypeChangeStringToNumber:
		return message + "may result in loss of precision"
	case TypeChangeScalarToOptional, TypeChangeOptionalToScalar:
		return message + "may result in undefined behavior"
	case TypeChangeScalarToUnion, TypeChangeUnionToScalar:
		return message + "may result in undefined behavior"
	// case UnionToUnion
	// case TypeChangeIncompatible
	default:
		return ""
	}
}

// Prepare Annotations on the old model
func preprocessPredecessor(predecessor *Environment) {
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations["schema"] = GetProtocolSchemaString(node, predecessor.SymbolTable)
			return

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			self.VisitChildren(node)
			return

		case *Field:
			if node.Annotations == nil {
				node.Annotations = make(map[string]any)
			}
			return

		default:
			self.VisitChildren(node)
		}
	})
}

// Prepare Annotations on the new model
func preprocessCurrent(env *Environment) {
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
			return

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
				node.Annotations["changes"] = make([]*TypeChange, 0)
			}
			return

		default:
			self.VisitChildren(node)
		}
	})
}

func annotateNamespaceChanges(newNode, oldNode *Namespace, version_index int) error {
	if newNode.Name != oldNode.Name {
		return validationError(newNode, "changing namespaces between versions is not yet supported")
		// log.Warn().Msgf("Changing namespaces between versions is not yet supported")
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
		changedTypeDef, err := annotateChangedTypeDefinition(newTd, oldTd, version_index)
		if err != nil {
			return err
		}
		if changedTypeDef != nil {
			changedTypeDef.GetDefinitionMeta().Annotations["version"] = version_index
		}

		// Annotate this TypeDefinition as having changed from previous version.
		newTd.GetDefinitionMeta().Annotations["changes"] = append(newTd.GetDefinitionMeta().Annotations["changes"].([]TypeDefinition), changedTypeDef)
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

		changedProtocolDef, err := annotateChangedProtocolDefinition(newProt, oldProt, version_index)
		if err != nil {
			return err
		}
		if changedProtocolDef != nil {
			oldSchema := oldProt.GetDefinitionMeta().Annotations["schema"]
			newProt.GetDefinitionMeta().Annotations["schemas"] = append(newProt.GetDefinitionMeta().Annotations["schemas"].([]string), oldSchema.(string))
		}

		// Annotate this ProtocolDefinition as having changed from previous version.
		newProt.GetDefinitionMeta().Annotations["changes"] = append(newProt.GetDefinitionMeta().Annotations["changes"].([]*ProtocolDefinition), changedProtocolDef)
	}

	return nil
}

// Compares two TypeDefinitions with matching names
func annotateChangedTypeDefinition(newNode, oldNode TypeDefinition, version_index int) (TypeDefinition, error) {
	switch newNode := newNode.(type) {
	case *RecordDefinition:
		oldNode, ok := oldNode.(*RecordDefinition)
		if !ok {
			return oldNode, fmt.Errorf("changing '%s' to a Record is not backward compatible", newNode.Name)
			// log.Warn().Msgf("Changing '%s' to a Record is not backward compatible", newNode.Name)
			// return oldNode
		}
		res, err := annotateChangedRecordDefinition(newNode, oldNode, version_index)
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
			// log.Warn().Msgf("Changing '%s' to a named type is not backward compatible", newNode.Name)
			// return oldNode
		}

		ch, err := detectChangedTypes(newNode.Type, oldNode.Type, version_index)
		if err != nil {
			return oldNode, err
		}

		if ch != nil {
			// // CHANGE: Changed NamedType type
			// if ch.Kind == TypeChangeIncompatible {
			// 	log.Warn().Msgf("Changing '%s' from '%s' to '%s' is not backward compatible", newNode.Name, TypeToShortSyntax(oldNode.Type, true), TypeToShortSyntax(newNode.Type, true))
			// }
			oldNode.Annotations["changed"] = ch
			return oldNode, nil
		}
		return nil, nil

	case *EnumDefinition:
		oldNode, ok := oldNode.(*EnumDefinition)
		if !ok {
			return oldNode, fmt.Errorf("changing '%s' to an Enum is not backward compatible", newNode.Name)
			// log.Warn().Msgf("Changing '%s' to an Enum is not backward compatible", newNode.Name)
			// return oldNode
		}
		res, err := annotateChangedEnumDefinitions(newNode, oldNode, version_index)
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
func annotateChangedProtocolDefinition(newProtocol, oldProtocol *ProtocolDefinition, version_index int) (*ProtocolDefinition, error) {
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
			return oldProtocol, fmt.Errorf("removing Protocol steps is not backward compatible")
			// log.Warn().Msgf("Removing a step from a Protocol is not backward compatible")
			// changed = true
			// continue
		}

		typeChange, err := detectChangedTypes(newStep.Type, oldStep.Type, version_index)
		if err != nil {
			return oldProtocol, err
		}

		if typeChange != nil {
			changed = true

			// log.Debug().Msgf("Protocol %s step %s changed from %s to %s", newProtocol.Name, newStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			// if typeChange.Kind == TypeChangeIncompatible {
			// 	log.Warn().Msgf("Changing step '%s' from '%s' to '%s' is not backward compatible", oldStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			// }
		}

		// Annotate the change to ProtocolStep so we can handle compatibility later in Protocol Reader/Writer
		newStep.Annotations["changes"] = append(newStep.Annotations["changes"].([]*TypeChange), typeChange)
	}

	if changed {
		return oldProtocol, nil
	}
	return nil, nil
}

// Compares two RecordDefinitions with matching names
func annotateChangedRecordDefinition(newRecord, oldRecord *RecordDefinition, version_index int) (*RecordDefinition, error) {
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
			oldRecord.Fields[i].Annotations["removed"] = true
			changed = true
			continue
		}

		// log.Debug().Msgf("Comparing fields %s and %s", newField.Name, oldField.Name)
		typeChange, err := detectChangedTypes(newField.Type, oldField.Type, version_index)
		if err != nil {
			return oldRecord, err
		}

		if typeChange != nil {
			// CHANGE: Changed field type
			changed = true
			oldRecord.Fields[i].Annotations["changed"] = typeChange
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

func annotateChangedEnumDefinitions(newNode, oldNode *EnumDefinition, version_index int) (*EnumDefinition, error) {
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
		ch, err := detectChangedTypes(newNode.BaseType, oldNode.BaseType, version_index)
		if err != nil {
			return oldNode, err
		}
		if ch != nil {
			// CHANGE: Changed Enum base type
			return oldNode, fmt.Errorf("changing '%s' base type is not backward compatible", newNode.Name)
			// log.Warn().Msgf("Changing '%s' base type is not backward compatible", newNode.Name)
			// changed = true
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
func detectChangedTypes(newType, oldType Type, version_index int) (*TypeChange, error) {
	// TODO: This is a good example of where it would be nice to bubble up Type Change User Warnings
	// so they are reported in the context of the `Field` or `ProtocolStep` that changed.
	//
	// UPDATE: The caller can just check if the TypeChange.Kind is TypeChangeIncompatible

	switch newType := newType.(type) {

	case *SimpleType:
		switch oldType := oldType.(type) {
		case *SimpleType:
			return detectSimpleTypeChanges(newType, oldType, version_index)
		case *GeneralizedType:
			return detectGeneralizedToSimpleTypeChanges(newType, oldType, version_index)
		default:
			panic("Shouldn't get here")
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return detectGeneralizedTypeChanges(newType, oldType, version_index)
		case *SimpleType:
			return detectSimpleToGeneralizedTypeChanges(newType, oldType, version_index)
		default:
			panic("Shouldn't get here")
		}

	default:
		panic("Expected a type")
	}
}

func detectSimpleTypeChanges(newType, oldType *SimpleType, version_index int) (*TypeChange, error) {
	// TODO: Compare TypeArguments
	// This comparison depends on whether the ResolvedDefinition changed!
	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
		// CHANGE: Changed number of TypeArguments

	} else {
		for i := range newType.TypeArguments {
			ch, err := detectChangedTypes(newType.TypeArguments[i], oldType.TypeArguments[i], version_index)
			if err != nil {
				return nil, err
			}
			if ch != nil {
				// CHANGE: Changed TypeArgument
				// TODO: Capture it
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
			return detectPrimitiveTypeChange(newType, oldType, version_index)
		}
		// log.Warn().Msgf("Converting non-primitive to primitive type is not backward compatible")
		// return &TypeChange{oldType, newType, TypeChangeIncompatible}
		return nil, TypeChangeIncompatibleError(oldType, newType)
	}

	if _, ok := oldDef.(PrimitiveDefinition); ok {
		// log.Warn().Msgf("Converting primitive to non-primitive type is not backward compatible")
		// return &TypeChange{oldType, newType, TypeChangeIncompatible}
		return nil, TypeChangeIncompatibleError(oldType, newType)
	}

	// At this point, both Types should be TypeDefinitions
	if newDef.GetDefinitionMeta().Name != oldDef.GetDefinitionMeta().Name {
		// CHANGE: Type changed to a different TypeDefinition
		// return &TypeChange{oldType, newType, TypeChangeDefinitionChanged}
		return nil, TypeChangeIncompatibleError(oldType, newType)
	}

	// log.Debug().Msgf("Comparing TypeDefinitions %s and %s", newDef.GetDefinitionMeta().Name, oldDef.GetDefinitionMeta().Name)

	// At this point, only the underlying TypeDefinition with matching name could have changed
	// And it would have been annotated earlier when comparing Namespace TypeDefinitions
	changes := newDef.GetDefinitionMeta().Annotations["changes"].([]TypeDefinition)
	if ch := changes[version_index]; ch != nil {
		// log.Debug().Msgf("SimpleType '%s' changed", newType.Name)
		return &TypeChange{oldType, newType, TypeChangeDefinitionChanged}, nil
	}

	// log.Debug().Msgf("SimpleType '%s' did NOT change", newType.Name)
	return nil, nil
}

/*
	TODO: Leverage the type functions:

- func GetPrimitiveType(t Type) (primitive PrimitiveDefinition, ok bool)
- func GetPrimitiveKind(t PrimitiveDefinition) PrimitiveKind
- func GetKindIfPrimitive(t Type) (primitiveKind PrimitiveKind, ok bool)
- func IsIntegralPrimitive(prim PrimitiveDefinition)
- func IsIntegralType(t Type) bool
*/
func detectPrimitiveTypeChange(newType, oldType *SimpleType, version_index int) (*TypeChange, error) {
	newPrimitive := newType.ResolvedDefinition.(PrimitiveDefinition)
	oldPrimitive := oldType.ResolvedDefinition.(PrimitiveDefinition)

	if newPrimitive == oldPrimitive {
		return nil, nil
	}

	// CHANGE: Changed Primitive type
	if oldPrimitive == PrimitiveString {
		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChange{oldType, newType, TypeChangeStringToNumber}, nil
		}
	}

	if GetPrimitiveKind(oldPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(oldPrimitive) == PrimitiveKindFloatingPoint {
		if newPrimitive == PrimitiveString {
			return &TypeChange{oldType, newType, TypeChangeNumberToString}, nil
		}

		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChange{oldType, newType, TypeChangeNumberToNumber}, nil
		}
	}

	return nil, TypeChangeIncompatibleError(oldType, newType)
}

func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType, version_index int) (*TypeChange, error) {
	// A GeneralizedType can change in many ways...
	changeKind := TypeChangeNoChange

	// Compare Dimensonality
	if ch := detectChangedDimensionality(newType.Dimensionality, oldType.Dimensionality, version_index); ch != nil {
		// CHANGE: Dimensionality changed
		return nil, TypeChangeIncompatibleError(oldType, newType)
	}

	if len(newType.Cases) != len(oldType.Cases) {
		// TODO: Handle adding types to a Union/Optional
		// TODO: Handle removing types from a Union/Optional

		return nil, TypeChangeIncompatibleError(oldType, newType)
	}

	// Compare the Type of each TypeCase
	for i, newCase := range newType.Cases {
		oldCase := oldType.Cases[i]

		if (newCase.Type == nil) != (oldCase.Type == nil) {
			// CHANGE: Added or removed a type case
			return nil, TypeChangeIncompatibleError(oldType, newType)
		}

		if newCase.Type == nil {
			continue
		}

		ch, err := detectChangedTypes(newCase.Type, oldCase.Type, version_index)
		if err != nil {
			return nil, err
		}
		if ch != nil {
			changeKind = ch.Kind
			break
		}
	}

	if changeKind != TypeChangeNoChange {
		return &TypeChange{oldType, newType, changeKind}, nil
	}
	return nil, nil
}

func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType, version_index int) (*TypeChange, error) {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		if TypesEqual(newType, oldType.Cases[1].Type) {
			return &TypeChange{oldType, newType, TypeChangeOptionalToScalar}, nil
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
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeUnionToScalar}, nil
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return nil, TypeChangeIncompatibleError(oldType, newType)
	// return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType, version_index int) (*TypeChange, error) {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		if TypesEqual(newType.Cases[1].Type, oldType) {
			return &TypeChange{oldType, newType, TypeChangeScalarToOptional}, nil
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
		// TODO: Probably need to capture index of the matching type in the new Union's Cases
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeScalarToUnion}, nil
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return nil, TypeChangeIncompatibleError(oldType, newType)
	// return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectChangedDimensionality(newNode, oldNode Dimensionality, version_index int) error {
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
		ch, err := detectChangedTypes(newDim.KeyType, oldDim.KeyType, version_index)
		if err != nil {
			return err
		}
		if ch != nil {
			return fmt.Errorf("changing map key type is not backward compatible")
		}
		return nil

	default:
		log.Panic().Msgf("unhandled type %T", newNode)
	}

	return nil
}
