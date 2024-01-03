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
	InnerChange TypeChange
	TypeIndex   int
}

func (tc *TypeChangeUnionToOptional) Inverse() TypeChange {
	return &TypeChangeOptionalToUnion{tc.Swap(), tc.InnerChange.Inverse(), tc.TypeIndex}
}

func (tc *TypeChangeUnionToOptional) Inner() TypeChange {
	return tc.InnerChange
}

type TypeChangeOptionalToUnion struct {
	TypePair
	InnerChange TypeChange
	TypeIndex   int
}

func (tc *TypeChangeOptionalToUnion) Inverse() TypeChange {
	return &TypeChangeUnionToOptional{tc.Swap(), tc.InnerChange.Inverse(), tc.TypeIndex}
}

func (tc *TypeChangeOptionalToUnion) Inner() TypeChange {
	return tc.InnerChange
}

type TypeChangeUnionTypesChange struct {
	TypePair
	InnerChanges []TypeChange
}

func (tc *TypeChangeUnionTypesChange) Inverse() TypeChange {
	innerChanges := make([]TypeChange, len(tc.InnerChanges))
	for i, ch := range tc.InnerChanges {
		innerChanges[i] = ch.Inverse()
	}
	return &TypeChangeUnionTypesChange{tc.Swap(), innerChanges}
}

type TypeChangeUnionShrink struct {
	TypePair
	InnerChanges []TypeChange
}

func (tc *TypeChangeUnionShrink) Inverse() TypeChange {
	innerChanges := make([]TypeChange, len(tc.InnerChanges))
	for i, ch := range tc.InnerChanges {
		innerChanges[i] = ch.Inverse()
	}
	return &TypeChangeUnionGrow{tc.Swap(), innerChanges}
}

type TypeChangeUnionGrow struct {
	TypePair
	InnerChanges []TypeChange
}

func (tc *TypeChangeUnionGrow) Inverse() TypeChange {
	innerChanges := make([]TypeChange, len(tc.InnerChanges))
	for i, ch := range tc.InnerChanges {
		innerChanges[i] = ch.Inverse()
	}
	return &TypeChangeUnionShrink{tc.Swap(), innerChanges}
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
	_ TypeChange = (*TypeChangeUnionTypesChange)(nil)
	_ TypeChange = (*TypeChangeUnionShrink)(nil)
	_ TypeChange = (*TypeChangeUnionGrow)(nil)
	_ TypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ TypeChange = (*TypeChangeVectorTypeChanged)(nil)
	_ TypeChange = (*TypeChangeDefinitionChanged)(nil)
	_ TypeChange = (*TypeChangeIncompatible)(nil)

	_ WrappedTypeChange = (*TypeChangeOptionalTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeVectorTypeChanged)(nil)
)

type DefinitionChange interface {
	PreviousVersion() TypeDefinition
}

type DefinitionChangeIncompatible struct {
	PreviousDefinition TypeDefinition
}

func (dc *DefinitionChangeIncompatible) PreviousVersion() TypeDefinition {
	return dc.PreviousDefinition
}

type NamedTypeChange struct {
	PreviousDefinition *NamedType
	TypeChange         TypeChange
}

func (ntc *NamedTypeChange) PreviousVersion() TypeDefinition {
	return ntc.PreviousDefinition
}

type ProtocolChange struct {
	PreviousDefinition *ProtocolDefinition
	PreviousSchema     string
	StepsAdded         []*ProtocolStep
	StepRemoved        []bool
	StepChanges        []TypeChange
	StepsReordered     bool
}

func (pc *ProtocolChange) PreviousVersion() TypeDefinition {
	return pc.PreviousDefinition
}

type RecordChange struct {
	PreviousDefinition *RecordDefinition
	FieldsAdded        []*Field
	FieldRemoved       []bool
	FieldChanges       []TypeChange
	FieldsReordered    bool
}

func (rc *RecordChange) PreviousVersion() TypeDefinition {
	return rc.PreviousDefinition
}

type EnumChange struct {
	PreviousDefinition *EnumDefinition
	BaseTypeChange     TypeChange
}

func (ec *EnumChange) PreviousVersion() TypeDefinition {
	return ec.PreviousDefinition
}

const (
	// Annotations referenced in serialization codegen
	AllVersionChangesAnnotationKey = "all-changes"
	VersionAnnotationKey           = "version"

	// Annotations used only for validation model evolution (local to this file)
	changeAnnotationKey = "changed"
	schemaAnnotationKey = "schema"
)

func ValidateEvolution(env *Environment, predecessors []*Environment, versionLabels []string) (*Environment, error) {

	initializeChangeAnnotations(env)

	for i, predecessor := range predecessors {
		log.Info().Msgf("Resolving changes from predecessor %s", versionLabels[i])
		annotatePredecessorSchemas(predecessor)

		if err := annotateAllChanges(env, predecessor); err != nil {
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
				oldRec := defChange.PreviousDefinition

				for _, added := range defChange.FieldsAdded {
					if !TypeHasNullOption(added.Type) {
						log.Warn().Msgf("Adding a non-Optional record field may result in undefined behavior")
					}
				}

				for i, field := range oldRec.Fields {
					if defChange.FieldRemoved[i] {
						if !TypeHasNullOption(oldRec.Fields[i].Type) {
							log.Warn().Msgf("Removing a non-Optional record field may result in undefined behavior")
						}
						continue
					}

					if tc := defChange.FieldChanges[i]; tc != nil {
						if typeChangeIsError(tc) {
							errorSink.Add(validationError(td, "changing field '%s' from %s", field.Name, typeChangeToError(tc)))
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

			oldProt := protChange.PreviousDefinition
			for _, added := range protChange.StepsAdded {
				errorSink.Add(validationError(added, "adding steps to a Protocol is not backward compatible"))
			}

			for i, step := range oldProt.Sequence {
				if protChange.StepRemoved[i] {
					errorSink.Add(validationError(step, "removing steps from a Protocol is not backward compatible"))
					continue
				}

				if tc := protChange.StepChanges[i]; tc != nil {
					if typeChangeIsError(tc) {
						errorSink.Add(validationError(pd, "changing step '%s' from %s", step.Name, typeChangeToError(tc)))
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
	case *TypeChangeUnionGrow, *TypeChangeUnionShrink:
		return "may result in undefined behavior"
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
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] == nil {
				node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = make(map[string]*ProtocolChange)
			}
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			if node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] == nil {
				node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey] = make(map[string]DefinitionChange)
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
			self.VisitChildren(node)

		case *ProtocolDefinition:
			var changed *ProtocolChange
			if ch, ok := node.GetDefinitionMeta().Annotations[changeAnnotationKey].(*ProtocolChange); ok {
				changed = ch
			}
			node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey].(map[string]*ProtocolChange)[versionLabel] = changed
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

			// Annotate *each* ProtocolStep with any changes from previous versions.
			// NOTE: This assumes ProtocolStep are NOT REORDERED between versions
			for i, step := range node.Sequence {
				if step.Annotations == nil {
					step.Annotations = make(map[string]any)
				}
				if step.Annotations[AllVersionChangesAnnotationKey] == nil {
					step.Annotations[AllVersionChangesAnnotationKey] = make(map[string]TypeChange)
				}

				var stepChange TypeChange
				if changed != nil && i < len(changed.StepChanges) {
					stepChange = changed.StepChanges[i]
				}

				step.Annotations[AllVersionChangesAnnotationKey].(map[string]TypeChange)[versionLabel] = stepChange
			}

		case TypeDefinition:
			var change DefinitionChange
			if ch, ok := node.GetDefinitionMeta().Annotations[changeAnnotationKey].(DefinitionChange); ok {
				change = ch

				// Annotate the previous Definition with VersionLabel so we can later use to write compatibility serializers
				if change.PreviousVersion().GetDefinitionMeta().Annotations == nil {
					change.PreviousVersion().GetDefinitionMeta().Annotations = make(map[string]any)
				}
				change.PreviousVersion().GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionLabel
			}
			node.GetDefinitionMeta().Annotations[AllVersionChangesAnnotationKey].(map[string]DefinitionChange)[versionLabel] = change
			node.GetDefinitionMeta().Annotations[changeAnnotationKey] = nil

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
			annotateNamespaceChanges(newNs, oldNs)
		} else {
			return fmt.Errorf("Namespace '%s' does not exist in previous version", newNs.Name)
		}
	}

	return nil
}

func annotateNamespaceChanges(newNs, oldNs *Namespace) {
	// TypeDefinitions may be reordered, added, or removed, so we compare them by name
	oldTds := make(map[string]TypeDefinition)
	for _, oldTd := range oldNs.TypeDefinitions {
		oldTds[oldTd.GetDefinitionMeta().Name] = oldTd
	}

	newTds := make(map[string]TypeDefinition)
	for _, newTd := range newNs.TypeDefinitions {
		newTds[newTd.GetDefinitionMeta().Name] = newTd
	}

	// Collect only pre-existing TypeDefinitions that are used within a Protocol
	// Keeping them in definition order!
	var newUsedTds []TypeDefinition
	for _, protocol := range newNs.Protocols {
		Visit(protocol, func(self Visitor, node Node) {
			switch node := node.(type) {
			case *ProtocolDefinition:
				self.VisitChildren(node)
			case TypeDefinition:
				self.VisitChildren(node)
				if td, ok := newTds[node.GetDefinitionMeta().Name]; ok {
					newUsedTds = append(newUsedTds, td)
				}
			case *SimpleType:
				self.VisitChildren(node)
				self.Visit(node.ResolvedDefinition)
			default:
				self.VisitChildren(node)
			}
		})
	}

	for _, newTd := range newUsedTds {
		oldTd, ok := oldTds[newTd.GetDefinitionMeta().Name]
		if !ok {
			// Skip new TypeDefinition
			continue
		}

		// Annotate this TypeDefinition with any changes from previous version.
		defChange := detectTypeDefinitionChanges(newTd, oldTd)
		newTd.GetDefinitionMeta().Annotations[changeAnnotationKey] = defChange
	}

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
		oldTd, ok := oldTd.(*RecordDefinition)
		if !ok {
			// Changing a non-Record to a Record is not backward compatible
			return &DefinitionChangeIncompatible{oldTd}
		}
		if ch := detectRecordDefinitionChanges(newNode, oldTd); ch != nil {
			return ch
		}
		return nil

	case *NamedType:
		oldTd, ok := oldTd.(*NamedType)
		if !ok {
			// Changing a non-NamedType to a NamedType is not backward compatible
			return &DefinitionChangeIncompatible{oldTd}
		}
		if typeChange := detectChangedTypes(newNode.Type, oldTd.Type); typeChange != nil {
			return &NamedTypeChange{oldTd, typeChange}
		}
		return nil

	case *EnumDefinition:
		oldTd, ok := oldTd.(*EnumDefinition)
		if !ok {
			return &DefinitionChangeIncompatible{oldTd}
		}
		if ch := detectEnumDefinitionChanges(newNode, oldTd); ch != nil {
			return ch
		}
		return nil

	default:
		panic("Expected a TypeDefinition")
	}
}

// Compares two ProtocolDefinitions with matching names
func detectProtocolDefinitionChanges(newProtocol, oldProtocol *ProtocolDefinition) *ProtocolChange {
	change := &ProtocolChange{
		PreviousDefinition: oldProtocol,
		PreviousSchema:     oldProtocol.GetDefinitionMeta().Annotations[schemaAnnotationKey].(string),
		StepRemoved:        make([]bool, len(oldProtocol.Sequence)),
		StepChanges:        make([]TypeChange, len(oldProtocol.Sequence)),
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

		if typeChange := detectChangedTypes(newStep.Type, oldStep.Type); typeChange != nil {
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
		PreviousDefinition: oldRecord,
		FieldRemoved:       make([]bool, len(oldRecord.Fields)),
		FieldChanges:       make([]TypeChange, len(oldRecord.Fields)),
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

		if typeChange := detectChangedTypes(newField.Type, oldField.Type); typeChange != nil {
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
		return &DefinitionChangeIncompatible{oldEnum}
	}

	oldBaseType := oldEnum.BaseType
	if oldBaseType == nil {
		oldBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	newBaseType := newNode.BaseType
	if newBaseType == nil {
		newBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	if ch := detectChangedTypes(newBaseType, oldBaseType); ch != nil {
		// CHANGE: Changed Enum base type
		return &EnumChange{PreviousDefinition: oldEnum, BaseTypeChange: ch}
	}

	return nil
}

// Compares two Types to determine if and how they changed
func detectChangedTypes(newType, oldType Type) TypeChange {
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
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if _, ok := oldDef.(PrimitiveDefinition); ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// At this point, both Types should be TypeDefinitions
	if newDef.GetDefinitionMeta().Name != oldDef.GetDefinitionMeta().Name {
		// CHANGE: Type changed to a different TypeDefinition
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// At this point, only the underlying TypeDefinition with matching name could have changed
	// And it would have been annotated earlier when comparing Namespace TypeDefinitions
	if ch, ok := newDef.GetDefinitionMeta().Annotations[changeAnnotationKey]; ok && ch != nil {
		return &TypeChangeDefinitionChanged{TypePair{oldType, newType}}
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
		if TypesEqual(newType, oldType.Cases[1].Type) {
			return &TypeChangeOptionalToScalar{TypePair{oldType, newType}}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		for i, tc := range oldType.Cases {
			if TypesEqual(newType, tc.Type) {
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
		if TypesEqual(newType.Cases[1].Type, oldType) {
			return &TypeChangeScalarToOptional{TypePair{oldType, newType}}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		for i, tc := range newType.Cases {
			if TypesEqual(tc.Type, oldType) {
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

// TODO: Handle Union to Optional
func detectOptionalChanges(newType, oldType *GeneralizedType) TypeChange {
	if !oldType.Cases.IsOptional() {
		if oldType.Cases.IsUnion() && oldType.Cases.HasNullOption() {
			// We want to find a matching type in the old Union
			// Should first look for identical type, but what if the user
			// changed the "matching" inner type in a compatible fashion
			// e.g. Union<string, Header_v1> to Optional<Header_v2>??

			for i, c := range oldType.Cases[1:] {
				// if c.IsNullType() {
				// 	continue
				// }
				if TypesEqual(newType.Cases[1].Type, c.Type) {
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, nil, i}
				}
			}

			for i, c := range oldType.Cases {
				ch := detectChangedTypes(newType.Cases[1].Type, c.Type)
				if ch == nil {
					// I think we shouldn't get here because we already checked above if TypesEqual...
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, nil, i}
				}

				switch ch.(type) {
				case *TypeChangeIncompatible:
					// Not a match
				default:
					// Match, but the inner type changed in a compatible way
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, ch, i}
				}
			}
		}

		// CHANGE: Changed a non-Optional/Union to an Optional
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectChangedTypes(newType.Cases[1].Type, oldType.Cases[1].Type); ch != nil {
		// CHANGE: Changed Optional type
		return &TypeChangeOptionalTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

// TODO: Handle Optional to Union
func detectUnionChanges(newType, oldType *GeneralizedType) TypeChange {
	if !oldType.Cases.IsUnion() {
		if oldType.Cases.IsOptional() && newType.Cases.HasNullOption() {
			for i, c := range newType.Cases[1:] {
				// if c.IsNullType() {
				// 	continue
				// }
				if TypesEqual(c.Type, oldType.Cases[1].Type) {
					return &TypeChangeOptionalToUnion{TypePair{oldType, newType}, nil, i}
				}
			}

			for i, c := range newType.Cases {
				ch := detectChangedTypes(c.Type, oldType.Cases[1].Type)
				if ch == nil {
					// I think we shouldn't get here because we already checked above if TypesEqual...
					return &TypeChangeOptionalToUnion{TypePair{oldType, newType}, nil, i}
				}

				switch ch.(type) {
				case *TypeChangeIncompatible:
					// Not a match
				default:
					// Match, but the inner type changed in a compatible way
					return &TypeChangeOptionalToUnion{TypePair{oldType, newType}, ch, i}
				}
			}

		}
		// CHANGE: Changed a non-Union/Optional to a Union
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// NOTE: Reordering Union types is currently NOT SUPPORTED
	if len(newType.Cases) == len(oldType.Cases) {
		innerChanges := make([]TypeChange, len(newType.Cases))
		for i, newCase := range newType.Cases {
			// Determine if newType and oldType Union types are an equal set
			oldCase := oldType.Cases[i]
			if ch := detectChangedTypes(newCase.Type, oldCase.Type); ch != nil {
				// CHANGE: Changed Union type
				if typeChangeIsError(ch) {
					return &TypeChangeIncompatible{TypePair{oldType, newType}}
				}
				innerChanges[i] = ch
			}
		}

		for _, ch := range innerChanges {
			if ch != nil {
				return &TypeChangeUnionTypesChange{TypePair{oldType, newType}, innerChanges}
			}
		}

	} else if len(newType.Cases) > len(oldType.Cases) {
		// Determine if the oldType Union types are a subset of the newType Union types
		// NOTE: Because type reordering is not yet supported, users can only add a type to the "end" of the Union
		innerChanges := make([]TypeChange, len(newType.Cases))
		for i, oldCase := range oldType.Cases {
			// Determine if newType and oldType Union types are an equal set
			newCase := newType.Cases[i]
			if ch := detectChangedTypes(newCase.Type, oldCase.Type); ch != nil {
				// CHANGE: Changed Union type
				if typeChangeIsError(ch) {
					return &TypeChangeIncompatible{TypePair{oldType, newType}}
				}
				innerChanges[i] = ch
			}
		}

		for _, ch := range innerChanges {
			if ch != nil {
				return &TypeChangeUnionGrow{TypePair{oldType, newType}, innerChanges}
			}
		}

	} else if len(newType.Cases) < len(oldType.Cases) {
		// Determine if the newType Union types are a subset of the oldType Union types
		innerChanges := make([]TypeChange, len(newType.Cases))
		for i, newCase := range newType.Cases {
			// Determine if newType and oldType Union types are an equal set
			oldCase := oldType.Cases[i]
			if ch := detectChangedTypes(newCase.Type, oldCase.Type); ch != nil {
				// CHANGE: Changed Union type
				if typeChangeIsError(ch) {
					return &TypeChangeIncompatible{TypePair{oldType, newType}}
				}
				innerChanges[i] = ch
			}
		}

		for _, ch := range innerChanges {
			if ch != nil {
				return &TypeChangeUnionShrink{TypePair{oldType, newType}, innerChanges}
			}
		}
	}

	return nil
}

func detectStreamChanges(newType, oldType *GeneralizedType) TypeChange {
	if _, ok := oldType.Dimensionality.(*Stream); !ok {
		// CHANGE: Changed a non-Stream to a Stream
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectChangedTypes(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
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

	if ch := detectChangedTypes(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
		// CHANGE: Changed Vector type
		return &TypeChangeVectorTypeChanged{TypePair{oldType, newType}, ch}
	}

	if (oldDim.Length == nil) != (newDim.Length == nil) {
		// CHANGE: Changed from a fixed-length Vector to a variable-length Vector or vice versa
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if newDim.Length != nil && *newDim.Length != *oldDim.Length {
		// CHANGE: Changed vector length
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}

func detectArrayChanges(newType, oldType *GeneralizedType) TypeChange {
	newDim := newType.Dimensionality.(*Array)
	oldDim, ok := oldType.Dimensionality.(*Array)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := detectChangedTypes(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
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

	if ch := detectChangedTypes(newDim.KeyType, oldDim.KeyType); ch != nil {
		// CHANGE: Changed Map key type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if ch := detectChangedTypes(newType.Cases[0].Type, oldType.Cases[0].Type); ch != nil {
		// CHANGE: Changed Map value type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}
