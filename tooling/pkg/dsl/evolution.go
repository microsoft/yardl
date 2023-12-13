// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
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

func ValidateEvolution(env *Environment, predecessor *Environment, versionId int) (*Environment, error) {
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

	oldNamespaces := make(map[string]*Namespace)
	for _, oldNs := range predecessor.Namespaces {
		oldNamespaces[oldNs.Name] = oldNs
	}
	for _, newNs := range env.Namespaces {
		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
			// TODO: This will detect changes between Imported Namespaces first
			// 		Ensure we use ValidationError messages to capture the source/location of the change error
			annotateChangedNamespace(newNs, oldNs, versionId)
		}
	}

	return env, nil
}

func annotateChangedNamespace(newNode, oldNode *Namespace, version_index int) {
	if newNode.Name != oldNode.Name {
		log.Warn().Msgf("Changing namespaces between versions is not yet supported")
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

	for _, oldTd := range oldNode.TypeDefinitions {
		newTd, ok := newTds[oldTd.GetDefinitionMeta().Name]
		if !ok {
			// CHANGE: Removed TypeDefinition
			continue
		}

		changedTypeDef := annotateChangedTypeDefinition(newTd, oldTd, version_index)
		if changedTypeDef != nil {
			changedTypeDef.GetDefinitionMeta().Annotations["version"] = version_index
		}

		// Mark the "new" TypeDefinition as having changed from previous version.
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

	for _, oldProt := range oldNode.Protocols {
		newProt, ok := newProts[oldProt.GetDefinitionMeta().Name]
		if !ok {
			// CHANGE: Removed ProtocolDefinition
			continue
		}

		changedProtocolDef := annotateChangedProtocolDefinition(newProt, oldProt, version_index)
		if changedProtocolDef != nil {
			oldSchema, ok := oldProt.GetDefinitionMeta().Annotations["schema"]
			if !ok {
				panic("Expected annotation containing old protocol schema string")
			}
			newProt.GetDefinitionMeta().Annotations["schemas"] = append(newProt.GetDefinitionMeta().Annotations["schemas"].([]string), oldSchema.(string))
		}

		// Mark the "new" TypeDefinition as having changed from previous version.
		newProt.GetDefinitionMeta().Annotations["changes"] = append(newProt.GetDefinitionMeta().Annotations["changes"].([]*ProtocolDefinition), changedProtocolDef)
	}
}

// Compares two TypeDefinitions with matching names
func annotateChangedTypeDefinition(newNode, oldNode TypeDefinition, version_index int) TypeDefinition {
	switch newNode := newNode.(type) {
	case *RecordDefinition:
		oldNode, ok := oldNode.(*RecordDefinition)
		if !ok {
			log.Warn().Msgf("Changing '%s' to a Record is not backward compatible", newNode.Name)
			return oldNode
		}
		if res := annotateChangedRecordDefinition(newNode, oldNode, version_index); res != nil {
			return res
		}
		return nil

	case *NamedType:
		oldNode, ok := oldNode.(*NamedType)
		if !ok {
			log.Warn().Msgf("Changing '%s' to a named type is not backward compatible", newNode.Name)
			return oldNode
		}

		if ch := annotateChangedTypes(newNode.Type, oldNode.Type, version_index); ch != nil {
			// CHANGE: Changed NamedType type
			if ch.Kind == TypeChangeIncompatible {
				log.Warn().Msgf("Changing '%s' from '%s' to '%s' is not backward compatible", newNode.Name, TypeToShortSyntax(oldNode.Type, true), TypeToShortSyntax(newNode.Type, true))
			}
			return oldNode
		}
		return nil

	case *EnumDefinition:
		oldNode, ok := oldNode.(*EnumDefinition)
		if !ok {
			log.Warn().Msgf("Changing '%s' to an Enum is not backward compatible", newNode.Name)
			return oldNode
		}
		if res := annotateChangedEnumDefinitions(newNode, oldNode, version_index); res != nil {
			return res
		}
		return nil

	// case PrimitiveDefinition:
	// 	oldNode := oldNode.(PrimitiveDefinition)
	// 	if newNode != oldNode {
	// 		// CHANGE: Changed Primitive type
	// 		return oldNode
	// 	}
	// 	return nil

	default:
		panic("Expected a TypeDefinition")
	}
}

// Compares two ProtocolDefinitions with matching names
func annotateChangedProtocolDefinition(newNode, oldNode *ProtocolDefinition, version_index int) *ProtocolDefinition {
	changed := false

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

	for _, oldStep := range oldNode.Sequence {
		newStep, ok := newSequence[oldStep.Name]
		if !ok {
			log.Warn().Msgf("Removing a step from a Protocol is not backward compatible")
			changed = true
			continue
		}

		typeChange := annotateChangedTypes(newStep.Type, oldStep.Type, version_index)
		if typeChange != nil {
			changed = true

			log.Debug().Msgf("Protocol %s step %s changed from %s to %s", newNode.Name, newStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			if typeChange.Kind == TypeChangeIncompatible {
				log.Warn().Msgf("Changing step '%s' from '%s' to '%s' is not backward compatible", oldStep.Name, TypeToShortSyntax(oldStep.Type, true), TypeToShortSyntax(newStep.Type, true))
			}
		}

		// Annotate the change to ProtocolStep so we can handle compatibility later in Protocol Reader/Writer
		newStep.Annotations["changes"] = append(newStep.Annotations["changes"].([]*TypeChange), typeChange)
	}

	if changed {
		return oldNode
	}
	return nil
}

// Compares two RecordDefinitions with matching names
func annotateChangedRecordDefinition(newRecord, oldRecord *RecordDefinition, version_index int) *RecordDefinition {
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
				log.Warn().Msgf("Adding a non-Optional record field is not backward compatible")
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
				log.Warn().Msgf("Removing a non-Optional record field is not backward compatible")
			}
			// CHANGE: Removed field
			oldRecord.Fields[i].Annotations["removed"] = true
			changed = true
			continue
		}

		// log.Debug().Msgf("Comparing fields %s and %s", newField.Name, oldField.Name)
		if typeChange := annotateChangedTypes(newField.Type, oldField.Type, version_index); typeChange != nil {
			// CHANGE: Changed field type
			changed = true
			oldRecord.Fields[i].Annotations["changed"] = typeChange
			if typeChange.Kind == TypeChangeIncompatible {
				log.Warn().Msgf("Changing field '%s' from '%s' to '%s' is not backward compatible", oldField.Name, TypeToShortSyntax(oldField.Type, true), TypeToShortSyntax(newField.Type, true))
			}
			continue
		}
	}

	if changed {
		// log.Debug().Msgf("Record '%s' changed", newRecord.Name)
		return oldRecord
	}
	// log.Debug().Msgf("Record '%s' did NOT change", newRecord.Name)
	return nil
}

func annotateChangedEnumDefinitions(newNode, oldNode *EnumDefinition, version_index int) *EnumDefinition {
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
		if ch := annotateChangedTypes(newNode.BaseType, oldNode.BaseType, version_index); ch != nil {
			// CHANGE: Changed Enum base type
			log.Warn().Msgf("Changing '%s' base type is not backward compatible", newNode.Name)
			changed = true
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
		return oldNode
	}

	return nil
}

// Compares two Types to determine what changed
// NOTE: We can't just use the `TypesEqual` function because we need to know *how* a Type changed.
func annotateChangedTypes(newType, oldType Type, version_index int) *TypeChange {
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

func detectSimpleTypeChanges(newType, oldType *SimpleType, version_index int) *TypeChange {
	// TODO: Compare TypeArguments
	// This comparison depends on whether the ResolvedDefinition changed!
	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
		// CHANGE: Changed number of TypeArguments

	} else {
		for i := range newType.TypeArguments {
			if ch := annotateChangedTypes(newType.TypeArguments[i], oldType.TypeArguments[i], version_index); ch != nil {
				// CHANGE: Changed TypeArgument
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
			return primitiveToPrimitiveTypeChange(newType, oldType, version_index)
		}
		log.Warn().Msgf("Converting non-primitive to primitive type is not backward compatible")
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	if _, ok := oldDef.(PrimitiveDefinition); ok {
		log.Warn().Msgf("Converting primitive to non-primitive type is not backward compatible")
		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	// At this point, both Types should be TypeDefinitions
	if newDef.GetDefinitionMeta().Name != oldDef.GetDefinitionMeta().Name {
		// CHANGE: Type changed to a different TypeDefinition
		return &TypeChange{oldType, newType, TypeChangeDefinitionChanged}
	}

	// log.Debug().Msgf("Comparing TypeDefinitions %s and %s", newDef.GetDefinitionMeta().Name, oldDef.GetDefinitionMeta().Name)

	// At this point, only the underlying TypeDefinition with matching name could have changed
	// And it would have been annotated earlier when comparing Namespace TypeDefinitions
	changes := newDef.GetDefinitionMeta().Annotations["changes"].([]TypeDefinition)
	if ch := changes[version_index]; ch != nil {
		// log.Debug().Msgf("SimpleType '%s' changed", newType.Name)
		return &TypeChange{oldType, newType, TypeChangeDefinitionChanged}
	}

	// log.Debug().Msgf("SimpleType '%s' did NOT change", newType.Name)
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
func primitiveToPrimitiveTypeChange(newType, oldType *SimpleType, version_index int) *TypeChange {
	newPrimitive := newType.ResolvedDefinition.(PrimitiveDefinition)
	oldPrimitive := oldType.ResolvedDefinition.(PrimitiveDefinition)

	if newPrimitive == oldPrimitive {
		return nil
	}

	// CHANGE: Changed Primitive type
	switch oldPrimitive {

	case PrimitiveString:
		switch newPrimitive {
		case PrimitiveInt8, PrimitiveInt16, PrimitiveInt32, PrimitiveInt64, PrimitiveFloat32, PrimitiveFloat64:
			return &TypeChange{oldType, newType, TypeChangeStringToNumber}
		}

	case PrimitiveInt8, PrimitiveInt16, PrimitiveInt32, PrimitiveInt64, PrimitiveFloat32, PrimitiveFloat64:
		switch newPrimitive {
		case PrimitiveString:
			return &TypeChange{oldType, newType, TypeChangeNumberToString}
		case PrimitiveInt8, PrimitiveInt16, PrimitiveInt32, PrimitiveInt64, PrimitiveFloat32, PrimitiveFloat64:
			return &TypeChange{oldType, newType, TypeChangeNumberToNumber}
		}
	}
	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType, version_index int) *TypeChange {
	changeKind := TypeChangeNoChange

	// TODO: Compare Cases
	// A GeneralizedType can change in many ways...

	if len(newType.Cases) != len(oldType.Cases) {
		// TODO: Handle adding types to a Union/Optional
		// TODO: Handle removing types from a Union/Optional

		return &TypeChange{oldType, newType, TypeChangeIncompatible}
	}

	// Compare the Type of each TypeCase
	for i, newCase := range newType.Cases {
		oldCase := oldType.Cases[i]

		if oldCase.Type == nil {
			if newCase.Type != nil {
				// TODO: Improve useless warning
				log.Warn().Msg("Adding a type case type is not backward compatible")
				changeKind = TypeChangeIncompatible
				continue
			}
			// Both types are nil, so no change for this TypeCase
			continue
		}

		if oldCase.Type != nil {
			if newCase.Type == nil {
				// TODO: Improve useless warning
				log.Warn().Msg("Removing a type case type is not backward compatible")
				changeKind = TypeChangeIncompatible
				continue
			}
		}

		if ch := annotateChangedTypes(newCase.Type, oldCase.Type, version_index); ch != nil {
			changeKind = ch.Kind
			break
		}
	}

	// Compare Dimensonality
	if ch := annotateChangedDimensionality(newType.Dimensionality, oldType.Dimensionality, version_index); ch != nil {
		changeKind = TypeChangeIncompatible
	}

	if changeKind != TypeChangeNoChange {
		return &TypeChange{oldType, newType, changeKind}
	}
	return nil

}

func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType, version_index int) *TypeChange {
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
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeUnionToScalar}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType, version_index int) *TypeChange {
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
		// TODO: Probably need to capture index of the matching type in the new Union's Cases
		if compatible {
			return &TypeChange{oldType, newType, TypeChangeScalarToUnion}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChange{oldType, newType, TypeChangeIncompatible}
}

func annotateChangedArrayDimensions(newNode, oldNode *ArrayDimension, version_index int) Node {
	changed := false
	// TODO: Do we care about named ArrayDimension changes?
	{
		if oldNode.Name == nil {
			if newNode.Name != nil {
				// CHANGE: Added array dimension name
				changed = true
			}
		}
		if newNode.Name == nil {
			// CHANGE: Removed array dimension name
			changed = true
		}
		if *newNode.Name != *oldNode.Name {
			// CHANGE: Renamed array dimension
			changed = true
		}
	}

	if oldNode.Length == nil {
		if newNode.Length != nil {
			// CHANGE: Added array dimension length
			changed = true
		}
	}
	if newNode.Length == nil {
		// CHANGE: Removed array dimension length
		changed = true
	}
	if *newNode.Length != *oldNode.Length {
		// CHANGE: Changed array dimension length
		changed = true
	}

	if changed {
		return oldNode
	}

	return nil
}

func annotateChangedDimensionality(newNode, oldNode Dimensionality, version_index int) Node {
	// TODO: Handle Dimensionality changes

	switch newDim := newNode.(type) {
	case nil, *Stream:
		return nil
	case *Vector:
		oldDim, ok := oldNode.(*Vector)
		if !ok {
			log.Warn().Msgf("expected a vector")
			return oldNode
		}
		changed := false
		if oldDim.Length == nil {
			if newDim.Length != nil {
				// CHANGE: Added vector length
				changed = true
			}
		} else {
			if newDim.Length == nil {
				// CHANGE: Removed vector length
				changed = true
			}
			if *newDim.Length != *oldDim.Length {
				// CHANGE: Changed vector length
				changed = true
			}
		}
		if changed {
			return oldNode
		}
		return nil

	case *Array:
		oldDim, ok := oldNode.(*Array)
		if !ok {
			log.Warn().Msgf("expected an array")
			return oldNode
		}
		changed := false
		if oldDim.Dimensions == nil {
			if newDim.Dimensions != nil {
				// CHANGE: Added array dimensions
				changed = true
			}
		} else {
			if newDim.Dimensions == nil {
				// CHANGE: Removed array dimensions
				changed = true
			}

			if len(*newDim.Dimensions) != len(*oldDim.Dimensions) {
				// CHANGE: Mismatch in number of array dimensions
				changed = true
			}

			for i := range *newDim.Dimensions {
				if ch := annotateChangedArrayDimensions((*newDim.Dimensions)[i], (*oldDim.Dimensions)[i], version_index); ch != nil {
					changed = true
				}
			}
		}

		if changed {
			return oldNode
		}
		return nil

	case *Map:
		oldDim, ok := oldNode.(*Map)
		if !ok {
			log.Warn().Msgf("expected a map")
			return oldNode
		}
		if ch := annotateChangedTypes(newDim.KeyType, oldDim.KeyType, version_index); ch != nil {
			return oldNode
		}
		return nil

	default:
		log.Panic().Msgf("unhandled type %T", newNode)
	}

	return nil
}
