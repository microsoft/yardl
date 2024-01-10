// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

type TypeChange interface {
	OldType() Type
	NewType() Type
	Inverse() TypeChange
}

type WrappedTypeChange interface {
	TypeChange
	Inner() TypeChange
}

type TypePair struct {
	Old Type
	New Type
}

func (tc TypePair) Swap() TypePair {
	return TypePair{tc.New, tc.Old}
}

func (tc *TypePair) OldType() Type {
	return tc.Old
}
func (tc *TypePair) NewType() Type {
	return tc.New
}

type TypeChangeNumberToNumber struct{ TypePair }

func (tc *TypeChangeNumberToNumber) Inverse() TypeChange {
	return &TypeChangeNumberToNumber{tc.Swap()}
}

type TypeChangeNumberToString struct{ TypePair }

func (tc *TypeChangeNumberToString) Inverse() TypeChange {
	return &TypeChangeStringToNumber{tc.Swap()}
}

type TypeChangeStringToNumber struct{ TypePair }

func (tc *TypeChangeStringToNumber) Inverse() TypeChange {
	return &TypeChangeNumberToString{tc.Swap()}
}

type TypeChangeScalarToOptional struct{ TypePair }

func (tc *TypeChangeScalarToOptional) Inverse() TypeChange {
	return &TypeChangeOptionalToScalar{tc.Swap()}
}

type TypeChangeOptionalToScalar struct{ TypePair }

func (tc *TypeChangeOptionalToScalar) Inverse() TypeChange {
	return &TypeChangeScalarToOptional{tc.Swap()}
}

type TypeChangeScalarToUnion struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeScalarToUnion) Inverse() TypeChange {
	return &TypeChangeUnionToScalar{tc.Swap(), tc.TypeIndex}
}

type TypeChangeUnionToScalar struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeUnionToScalar) Inverse() TypeChange {
	return &TypeChangeScalarToUnion{tc.Swap(), tc.TypeIndex}
}

type TypeChangeOptionalTypeChanged struct {
	TypePair
	InnerChange TypeChange
}

func (tc *TypeChangeOptionalTypeChanged) Inverse() TypeChange {
	return &TypeChangeOptionalTypeChanged{tc.Swap(), tc.InnerChange.Inverse()}
}

func (tc *TypeChangeOptionalTypeChanged) Inner() TypeChange {
	return tc.InnerChange
}

type TypeChangeUnionToOptional struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeUnionToOptional) Inverse() TypeChange {
	return &TypeChangeOptionalToUnion{tc.Swap(), tc.TypeIndex}
}

type TypeChangeOptionalToUnion struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeOptionalToUnion) Inverse() TypeChange {
	return &TypeChangeUnionToOptional{tc.Swap(), tc.TypeIndex}
}

type TypeChangeUnionTypesetChanged struct {
	TypePair
	OldMatches []bool
	NewMatches []bool
}

func (tc *TypeChangeUnionTypesetChanged) Inverse() TypeChange {
	return &TypeChangeUnionTypesetChanged{TypePair: tc.Swap(), OldMatches: tc.NewMatches, NewMatches: tc.OldMatches}
}

type TypeChangeStreamTypeChanged struct {
	TypePair
	InnerChange TypeChange
}

func (tc *TypeChangeStreamTypeChanged) Inverse() TypeChange {
	return &TypeChangeStreamTypeChanged{tc.Swap(), tc.InnerChange.Inverse()}
}

func (tc *TypeChangeStreamTypeChanged) Inner() TypeChange {
	return tc.InnerChange
}

type TypeChangeVectorTypeChanged struct {
	TypePair
	InnerChange TypeChange
}

func (tc *TypeChangeVectorTypeChanged) Inverse() TypeChange {
	return &TypeChangeVectorTypeChanged{tc.Swap(), tc.InnerChange.Inverse()}
}

func (tc *TypeChangeVectorTypeChanged) Inner() TypeChange {
	return tc.InnerChange
}

type TypeChangeDefinitionChanged struct{ TypePair }

func (tc *TypeChangeDefinitionChanged) Inverse() TypeChange {
	return &TypeChangeDefinitionChanged{tc.Swap()}
}

type TypeChangeIncompatible struct{ TypePair }

func (tc *TypeChangeIncompatible) Inverse() TypeChange {
	return &TypeChangeIncompatible{tc.Swap()}
}

var (
	_ TypeChange = (*TypeChangeNumberToNumber)(nil)
	_ TypeChange = (*TypeChangeNumberToString)(nil)
	_ TypeChange = (*TypeChangeStringToNumber)(nil)
	_ TypeChange = (*TypeChangeScalarToOptional)(nil)
	_ TypeChange = (*TypeChangeOptionalToScalar)(nil)
	_ TypeChange = (*TypeChangeScalarToUnion)(nil)
	_ TypeChange = (*TypeChangeUnionToScalar)(nil)
	_ TypeChange = (*TypeChangeOptionalTypeChanged)(nil)
	_ TypeChange = (*TypeChangeOptionalToUnion)(nil)
	_ TypeChange = (*TypeChangeUnionToOptional)(nil)
	_ TypeChange = (*TypeChangeUnionTypesetChanged)(nil)
	_ TypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ TypeChange = (*TypeChangeVectorTypeChanged)(nil)
	_ TypeChange = (*TypeChangeDefinitionChanged)(nil)
	_ TypeChange = (*TypeChangeIncompatible)(nil)

	_ WrappedTypeChange = (*TypeChangeOptionalTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeVectorTypeChanged)(nil)

	_ DefinitionChange = (*DefinitionChangeIncompatible)(nil)
	_ DefinitionChange = (*NamedTypeChange)(nil)
	_ DefinitionChange = (*RecordChange)(nil)
	_ DefinitionChange = (*EnumChange)(nil)
	_ DefinitionChange = (*ProtocolChange)(nil)
)

type DefinitionChange interface {
	PreviousDefinition() TypeDefinition
	LatestDefinition() TypeDefinition
}

type DefinitionPair struct {
	Old TypeDefinition
	New TypeDefinition
}

func (tc *DefinitionPair) PreviousDefinition() TypeDefinition {
	return tc.Old
}
func (tc *DefinitionPair) LatestDefinition() TypeDefinition {
	return tc.New
}

type DefinitionChangeIncompatible struct {
	DefinitionPair
}

type NamedTypeChange struct {
	DefinitionPair
	TypeChange TypeChange
}

type ProtocolChange struct {
	DefinitionPair
	PreviousSchema string
	StepsAdded     []*ProtocolStep
	StepRemoved    []bool
	StepChanges    []TypeChange
	StepsReordered bool
}

type RecordChange struct {
	DefinitionPair
	FieldsAdded     []*Field
	FieldRemoved    []bool
	FieldChanges    []TypeChange
	FieldsReordered bool
}

type EnumChange struct {
	DefinitionPair
	BaseTypeChange TypeChange
}

const (
	// Annotations referenced in serialization codegen
	VersionAnnotationKey = "version"

	// Annotations used only for validation model evolution (local to this file)
	changeAnnotationKey = "changed"
	schemaAnnotationKey = "schema"
)

func ValidateEvolution(env *Environment, predecessors []*Environment, versionLabels []string) (*Environment, error) {

	initializeChangeAnnotations(env)

	for i, predecessor := range predecessors {
		log.Info().Msgf("Resolving changes from predecessor %s", versionLabels[i])
		annotatePredecessorSchemas(predecessor)

		if err := annotateAllChanges(env, predecessor, versionLabels[i]); err != nil {
			return nil, err
		}

		if err := validateChanges(env); err != nil {
			return nil, err
		}

		saveChangeAnnotations(env, versionLabels[i])
	}

	return env, nil
}

func validateChanges(env *Environment) error {
	// Emit User Warnings and aggregate Errors
	errorSink := &validation.ErrorSink{}
	for _, ns := range env.Namespaces {

		for _, td := range ns.TypeDefinitions {
			defChange, ok := td.GetDefinitionMeta().Annotations[changeAnnotationKey].(DefinitionChange)
			if !ok || defChange == nil {
				continue
			}

			switch defChange := defChange.(type) {

			case *DefinitionChangeIncompatible:
				errorSink.Add(validationError(td, "changing '%s' is not backward compatible", td.GetDefinitionMeta().Name))

			case *RecordChange:
				oldRec := defChange.PreviousDefinition().(*RecordDefinition)

				for _, added := range defChange.FieldsAdded {
					if !TypeHasNullOption(added.Type) {
						log.Warn().Msgf("Adding a non-Optional record field may result in undefined behavior")
					}
				}

				for i, field := range oldRec.Fields {
					if defChange.FieldRemoved[i] {
						if !TypeHasNullOption(oldRec.Fields[i].Type) {
							log.Warn().Msgf("Removing non-Optional field '%s' may result in undefined behavior", field.Name)
						}
						continue
					}

					if tc := defChange.FieldChanges[i]; tc != nil {
						if typeChangeIsError(tc) {
							errorSink.Add(validationError(tc.NewType(), "changing field '%s' from %s", field.Name, typeChangeToError(tc)))
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

		for _, pd := range ns.Protocols {
			protChange, ok := pd.GetDefinitionMeta().Annotations[changeAnnotationKey].(*ProtocolChange)
			if !ok || protChange == nil {
				continue
			}

			if protChange.StepsReordered {
				errorSink.Add(validationError(pd, "reordering steps in a Protocol is not backward compatible"))
			}

			oldProt := protChange.PreviousDefinition().(*ProtocolDefinition)
			for _, added := range protChange.StepsAdded {
				errorSink.Add(validationError(added, "adding steps to a Protocol is not backward compatible"))
			}

			for i, step := range oldProt.Sequence {
				if protChange.StepRemoved[i] {
					errorSink.Add(validationError(pd, "removing step '%s' from a Protocol is not backward compatible", step.Name))
					continue
				}

				if tc := protChange.StepChanges[i]; tc != nil {
					if typeChangeIsError(tc) {
						errorSink.Add(validationError(tc.NewType(), "changing step '%s' from %s", step.Name, typeChangeToError(tc)))
					}

					if warn := typeChangeToWarning(tc); warn != "" {
						log.Warn().Msgf("Changing step '%s' from %s", step.Name, warn)
					}
				}
			}
		}
	}

	return errorSink.AsError()
}

func typeChangeIsError(tc TypeChange) bool {
	switch tc := tc.(type) {
	case WrappedTypeChange:
		switch tc := tc.(type) {
		case *TypeChangeStreamTypeChanged:
			// A Stream's Type can only change if it is a changed TypeDefinition
			if _, ok := tc.Inner().(*TypeChangeDefinitionChanged); !ok {
				return true
			}
		case *TypeChangeVectorTypeChanged:
			// A Vector's Type can only change if it is a changed TypeDefinition
			if _, ok := tc.Inner().(*TypeChangeDefinitionChanged); !ok {
				return true
			}
		}

		// Otherwise, is the inner type change an error?
		return typeChangeIsError(tc.Inner())

	case *TypeChangeIncompatible:
		return true
	}
	return false
}

func typeChangeToError(tc TypeChange) string {
	return fmt.Sprintf("'%s' to '%s' is not backward compatible", TypeToShortSyntax(tc.OldType(), true), TypeToShortSyntax(tc.NewType(), true))
}

func typeChangeWarningReason(tc TypeChange) string {
	switch tc := tc.(type) {
	case WrappedTypeChange:
		return typeChangeWarningReason(tc.Inner())

	case *TypeChangeNumberToNumber:
		return "may result in numeric overflow or loss of precision"
	case *TypeChangeNumberToString, *TypeChangeStringToNumber:
		return "may result in loss of precision"
	case *TypeChangeScalarToOptional, *TypeChangeOptionalToScalar:
		return "may result in undefined behavior"
	case *TypeChangeScalarToUnion, *TypeChangeUnionToScalar:
		return "may result in undefined behavior"
	case *TypeChangeUnionTypesetChanged:
		return "may produce runtime errors"
	}
	return ""
}

func typeChangeToWarning(tc TypeChange) string {
	message := fmt.Sprintf("'%s' to '%s' ", TypeToShortSyntax(tc.OldType(), true), TypeToShortSyntax(tc.NewType(), true))
	if reason := typeChangeWarningReason(tc); reason != "" {
		return message + reason
	}
	return ""
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
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

		default:
			self.VisitChildren(node)
		}
	})
}

func saveChangeAnnotations(env *Environment, versionLabel string) {
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
			if ch, ok := node.GetDefinitionMeta().Annotations[changeAnnotationKey].(*ProtocolChange); ok {
				changed = ch
			}
			node.Versions[versionLabel] = changed
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

		default:
			self.VisitChildren(node)
		}
	})
}

func annotateAllChanges(newNode, oldNode *Environment, versionLabel string) error {
	oldNamespaces := make(map[string]*Namespace)
	for _, oldNs := range oldNode.Namespaces {
		oldNamespaces[oldNs.Name] = oldNs
	}

	for _, newNs := range newNode.Namespaces {
		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
			annotateNamespaceChanges(newNs, oldNs, versionLabel)
		} else {
			return fmt.Errorf("Namespace '%s' does not exist in previous version", newNs.Name)
		}
	}

	return nil
}

func annotateNamespaceChanges(newNs, oldNs *Namespace, versionLabel string) {
	// TypeDefinitions may be reordered, added, or removed, so we compare them by name
	oldTds := make(map[string]TypeDefinition)
	for _, oldTd := range oldNs.TypeDefinitions {
		oldTds[oldTd.GetDefinitionMeta().Name] = oldTd
	}

	isUserTypeDef := make(map[string]bool)
	for _, newTd := range newNs.TypeDefinitions {
		isUserTypeDef[newTd.GetDefinitionMeta().Name] = true
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
				name := node.GetDefinitionMeta().Name
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
		oldTd, ok := oldTds[newTd.GetDefinitionMeta().Name]
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

		// "Unwind" any NamedTypes so we only compare underlying TypeDefinitions
		oldTd = unwindOldAlias(oldTd)
		newTd = unwindNewAlias(newTd)

		if alreadyCompared[newTd.GetDefinitionMeta().Name] {
			// TODO: Remove this check if not needed, once integration tests are "complete"
			panic(fmt.Sprintf("Already Compared %s", newTd.GetDefinitionMeta().Name))
			continue
		}

		defChange := detectTypeDefinitionChanges(newTd, oldTd)
		if defChange != nil {
			typeDefChanges = append(typeDefChanges, defChange)
			typeDefChanges = append(typeDefChanges, removedAliases...)

			// Annotate this TypeDefinition if it changed from previous version
			// Later, detectSimpleTypeChanges will look for this to see if an underlying TypeDefinition changed
			newTd.GetDefinitionMeta().Annotations[changeAnnotationKey] = defChange
		}

		alreadyCompared[newTd.GetDefinitionMeta().Name] = true
	}

	// Save all TypeDefinition changes for generating of compatibility serializers
	newNs.TypeDefChanges[versionLabel] = typeDefChanges

	// Protocols may be reordered, added, or removed
	// We only care about pre-existing Protocols that CHANGED
	oldProts := make(map[string]*ProtocolDefinition)
	for _, oldProt := range oldNs.Protocols {
		oldProts[oldProt.Name] = oldProt
	}

	for _, newProt := range newNs.Protocols {
		oldProt, ok := oldProts[newProt.GetDefinitionMeta().Name]
		if !ok {
			// Skip new ProtocolDefinition
			continue
		}

		// Annotate this ProtocolDefinition with any changes from previous version.
		protocolChange := detectProtocolDefinitionChanges(newProt, oldProt)
		newProt.GetDefinitionMeta().Annotations[changeAnnotationKey] = protocolChange
	}
}

// Compares two TypeDefinitions with matching names
func detectTypeDefinitionChanges(newTd, oldTd TypeDefinition) DefinitionChange {
	switch newNode := newTd.(type) {
	case *RecordDefinition:
		switch oldTd := oldTd.(type) {
		case *RecordDefinition:
			if ch := detectRecordDefinitionChanges(newNode, oldTd); ch != nil {
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
			if typeChange := detectTypeChanges(newNode.Type, oldTd.Type); typeChange != nil {
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
		if ch := detectEnumDefinitionChanges(newNode, oldTd); ch != nil {
			return ch
		}
		return nil

	default:
		// log.Debug().Msgf("What is this? %s was %s", newNode, oldTd)
		panic("Expected a TypeDefinition")
	}
}

// Compares two ProtocolDefinitions with matching names
func detectProtocolDefinitionChanges(newProtocol, oldProtocol *ProtocolDefinition) *ProtocolChange {
	change := &ProtocolChange{
		DefinitionPair: DefinitionPair{oldProtocol, newProtocol},
		PreviousSchema: oldProtocol.GetDefinitionMeta().Annotations[schemaAnnotationKey].(string),
		StepRemoved:    make([]bool, len(oldProtocol.Sequence)),
		StepChanges:    make([]TypeChange, len(oldProtocol.Sequence)),
	}

	oldSequence := make(map[string]*ProtocolStep)
	for _, f := range oldProtocol.Sequence {
		oldSequence[f.Name] = f
	}
	newSequence := make(map[string]*ProtocolStep)
	for i, newStep := range newProtocol.Sequence {
		newSequence[newStep.Name] = newStep

		if i >= len(oldProtocol.Sequence) || newStep.Name != oldProtocol.Sequence[i].Name {
			// CHANGE: Reordered or renamed steps
			change.StepsReordered = true
		}

		if _, ok := oldSequence[newStep.Name]; !ok {
			// CHANGE: New ProtocolStep
			change.StepsAdded = append(change.StepsAdded, newStep)
		}
	}

	anyStepChanged := false
	for i, oldStep := range oldProtocol.Sequence {
		newStep, ok := newSequence[oldStep.Name]
		if !ok {
			// CHANGE: Removed ProtocolStep
			anyStepChanged = true
			change.StepRemoved[i] = true
			continue
		}

		if typeChange := detectTypeChanges(newStep.Type, oldStep.Type); typeChange != nil {
			// CHANGE: ProtocolStep type changed
			anyStepChanged = true
			change.StepChanges[i] = typeChange
		}
	}

	if anyStepChanged || change.StepsReordered || len(change.StepsAdded) > 0 {
		return change
	}
	return nil
}

// Compares two RecordDefinitions with matching names
func detectRecordDefinitionChanges(newRecord, oldRecord *RecordDefinition) *RecordChange {
	change := &RecordChange{
		DefinitionPair: DefinitionPair{oldRecord, newRecord},
		FieldRemoved:   make([]bool, len(oldRecord.Fields)),
		FieldChanges:   make([]TypeChange, len(oldRecord.Fields)),
	}

	// Fields may be reordered
	// If they are, we want result to represent the old Record, for Serialization compatibility
	oldFields := make(map[string]*Field)
	for _, f := range oldRecord.Fields {
		oldFields[f.Name] = f
	}

	newFields := make(map[string]*Field)
	for i, newField := range newRecord.Fields {
		newFields[newField.Name] = newField

		if i >= len(oldRecord.Fields) || newField.Name != oldRecord.Fields[i].Name {
			// CHANGE: Reordered or renamed fields
			change.FieldsReordered = true
		}

		if _, ok := oldFields[newField.Name]; !ok {
			// CHANGE: New field
			change.FieldsAdded = append(change.FieldsAdded, newField)
		}
	}

	anyFieldChanged := false
	for i, oldField := range oldRecord.Fields {
		newField, ok := newFields[oldField.Name]
		if !ok {
			// CHANGE: Removed field
			anyFieldChanged = true
			change.FieldRemoved[i] = true
			continue
		}

		if typeChange := detectTypeChanges(newField.Type, oldField.Type); typeChange != nil {
			// CHANGE: Field type changed
			anyFieldChanged = true
			change.FieldChanges[i] = typeChange
		}
	}

	if anyFieldChanged || change.FieldsReordered || len(change.FieldsAdded) > 0 {
		return change
	}
	return nil
}

func detectEnumDefinitionChanges(newNode, oldEnum *EnumDefinition) DefinitionChange {
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

	if ch := detectTypeChanges(newBaseType, oldBaseType); ch != nil {
		// CHANGE: Changed Enum base type
		return &EnumChange{DefinitionPair{oldEnum, newNode}, ch}
	}

	return nil
}

// Compares two Types to determine if and how they changed
func detectTypeChanges(newType, oldType Type) TypeChange {
	newType = GetUnderlyingType(newType)
	oldType = GetUnderlyingType(oldType)

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

func detectSimpleTypeChanges(newType, oldType *SimpleType) TypeChange {
	// TODO: Compare TypeArguments
	// This comparison depends on whether the ResolvedDefinition changed!
	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
		// CHANGE: Changed number of TypeArguments
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	} else {
		for i := range newType.TypeArguments {
			if ch := detectTypeChanges(newType.TypeArguments[i], oldType.TypeArguments[i]); ch != nil {
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

	if ch, ok := newDef.GetDefinitionMeta().Annotations[changeAnnotationKey].(DefinitionChange); ok {
		if ch != nil && ch.PreviousDefinition() == oldDef {
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}}
		}
	}

	if newDef.GetDefinitionMeta().Name != oldDef.GetDefinitionMeta().Name {
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

func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType) TypeChange {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		switch detectTypeChanges(newType, oldType.Cases[1].Type).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeOptionalToScalar{TypePair{oldType, newType}}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		for i, tc := range oldType.Cases {
			switch detectTypeChanges(newType, tc.Type).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeUnionToScalar{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType) TypeChange {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		switch detectTypeChanges(newType.Cases[1].Type, oldType).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeScalarToOptional{TypePair{oldType, newType}}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		for i, tc := range newType.Cases {
			switch detectTypeChanges(tc.Type, oldType).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeScalarToUnion{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType) TypeChange {
	// A GeneralizedType can change in many ways...
	if newType.Cases.IsOptional() {
		return detectOptionalChanges(newType, oldType)
	} else if newType.Cases.IsUnion() {
		return detectUnionChanges(newType, oldType)
	} else {
		switch newType.Dimensionality.(type) {
		case nil:
			// TODO: Not an Optional, Union, Stream, Vector, Array, Map...
		case *Stream:
			return detectStreamChanges(newType, oldType)
		case *Vector:
			return detectVectorChanges(newType, oldType)
		case *Array:
			return detectArrayChanges(newType, oldType)
		case *Map:
			return detectMapChanges(newType, oldType)
		default:
			panic("Shouldn't get here")
		}
	}

	return nil
}

func detectOptionalChanges(newType, oldType *GeneralizedType) TypeChange {
	if !oldType.Cases.IsOptional() {
		if oldType.Cases.IsUnion() && oldType.Cases.HasNullOption() {
			// An Optional<T> can become a Union<null, T, ...> ONLY if
			// 	1. type T does not change, or
			// 	2. type T's TypeDefinition changed

			// Look for a matching type in the old Union
			for i, c := range oldType.Cases[1:] {
				switch detectTypeChanges(newType.Cases[1].Type, c.Type).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, i + 1}
				}
			}
		}

		// CHANGE: Changed a non-Optional/Union to an Optional
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[1].Type, oldType.Cases[1].Type); ch != nil {
		// CHANGE: Changed Optional type
		return &TypeChangeOptionalTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectUnionChanges(newType, oldType *GeneralizedType) TypeChange {
	if !oldType.Cases.IsUnion() {
		if oldType.Cases.IsOptional() && newType.Cases.HasNullOption() {
			for i, c := range newType.Cases[1:] {
				switch detectTypeChanges(c.Type, oldType.Cases[1].Type).(type) {
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

			switch detectTypeChanges(newCase.Type, oldCase.Type).(type) {
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

func detectStreamChanges(newType, oldType *GeneralizedType) TypeChange {
	if _, ok := oldType.Dimensionality.(*Stream); !ok {
		// CHANGE: Changed a non-Stream to a Stream
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
		// CHANGE: Changed Stream type
		return &TypeChangeStreamTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectVectorChanges(newType, oldType *GeneralizedType) TypeChange {
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

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
		// CHANGE: Changed Vector type
		return &TypeChangeVectorTypeChanged{TypePair{oldType, newType}, ch}
	}

	return nil
}

func detectArrayChanges(newType, oldType *GeneralizedType) TypeChange {
	newDim := newType.Dimensionality.(*Array)
	oldDim, ok := oldType.Dimensionality.(*Array)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
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

func detectMapChanges(newType, oldType *GeneralizedType) TypeChange {
	newDim := newType.Dimensionality.(*Map)
	oldDim, ok := oldType.Dimensionality.(*Map)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectTypeChanges(newDim.KeyType, oldDim.KeyType); ch != nil {
		// CHANGE: Changed Map key type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
		// CHANGE: Changed Map value type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}
