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

// TODO: Unit test this
func (tc *TypeChangeNumberToNumber) IsDemotion() bool {
	newPrim := tc.Old.(*SimpleType).ResolvedDefinition.(PrimitiveDefinition)
	oldPrim := tc.New.(*SimpleType).ResolvedDefinition.(PrimitiveDefinition)

	if newPrim == oldPrim {
		return false
	}

	width := func(prim PrimitiveDefinition) int {
		switch prim {
		case PrimitiveInt8, PrimitiveUint8:
			return 1
		case PrimitiveInt16, PrimitiveUint16:
			return 2
		case PrimitiveInt32, PrimitiveUint32, PrimitiveFloat32:
			return 4
		case PrimitiveInt64, PrimitiveUint64, PrimitiveFloat64:
			return 8
		default:
			panic("Shouldn't get here")
		}
	}

	isPromotion := func(a, b PrimitiveDefinition) bool {
		if GetPrimitiveKind(a) == PrimitiveKindInteger && GetPrimitiveKind(b) == PrimitiveKindFloatingPoint {
			return true
		} else if width(b) > width(a) {
			return true
		}
		return false
	}
	return !isPromotion(oldPrim, newPrim)
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

type TypeChangeUnionShrink struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeUnionShrink) Inverse() TypeChange {
	return &TypeChangeUnionGrow{tc.Swap(), tc.TypeIndex}
}

type TypeChangeUnionGrow struct {
	TypePair
	TypeIndex int
}

func (tc *TypeChangeUnionGrow) Inverse() TypeChange {
	return &TypeChangeUnionShrink{tc.Swap(), tc.TypeIndex}
}

type TypeChangeOptionalTypeChanged struct {
	TypePair
	Change TypeChange
}

func (tc *TypeChangeOptionalTypeChanged) Inverse() TypeChange {
	return &TypeChangeOptionalTypeChanged{tc.Swap(), tc.Change.Inverse()}
}

func (tc *TypeChangeOptionalTypeChanged) Inner() TypeChange {
	return tc.Change
}

type TypeChangeStreamTypeChanged struct {
	TypePair
	Change TypeChange
}

func (tc *TypeChangeStreamTypeChanged) Inverse() TypeChange {
	return &TypeChangeStreamTypeChanged{tc.Swap(), tc.Change.Inverse()}
}

func (tc *TypeChangeStreamTypeChanged) Inner() TypeChange {
	return tc.Change
}

type TypeChangeVectorTypeChanged struct {
	TypePair
	Change TypeChange
}

func (tc *TypeChangeVectorTypeChanged) Inverse() TypeChange {
	return &TypeChangeVectorTypeChanged{tc.Swap(), tc.Change.Inverse()}
}

func (tc *TypeChangeVectorTypeChanged) Inner() TypeChange {
	return tc.Change
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
	_ TypeChange = (*TypeChangeUnionShrink)(nil)
	_ TypeChange = (*TypeChangeUnionGrow)(nil)
	_ TypeChange = (*TypeChangeOptionalTypeChanged)(nil)
	_ TypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ TypeChange = (*TypeChangeVectorTypeChanged)(nil)
	_ TypeChange = (*TypeChangeDefinitionChanged)(nil)
	_ TypeChange = (*TypeChangeIncompatible)(nil)

	_ WrappedTypeChange = (*TypeChangeOptionalTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeStreamTypeChanged)(nil)
	_ WrappedTypeChange = (*TypeChangeVectorTypeChanged)(nil)
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
		log.Info().Msgf("Resolving changes from predecessor %d", versionId)
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
						if tc, ok := field.Annotations[ChangeAnnotationKey].(TypeChange); ok {
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
					if tc, ok := prevDef.Annotations[ChangeAnnotationKey].(TypeChange); ok {
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
			for _, step := range pd.Sequence {
				if tc, ok := step.Annotations[ChangeAnnotationKey].(TypeChange); ok {
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
		if tc.IsDemotion() {
			return "may result in loss of precision"
		}
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

// Prepare Annotations on the old model
func initializePredecessorAnnotations(predecessor *Environment) {
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations[SchemaAnnotationKey] = GetProtocolSchemaString(node, predecessor.SymbolTable)
			self.VisitChildren(node)

		case TypeDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
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
				node.Annotations[AllVersionChangesAnnotationKey] = make([]TypeChange, 0)
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
				var changed TypeChange
				if ch, ok := step.Annotations[ChangeAnnotationKey].(TypeChange); ok {
					changed = ch
				}
				step.Annotations[AllVersionChangesAnnotationKey] = append(step.Annotations[AllVersionChangesAnnotationKey].([]TypeChange), changed)
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

func annotateNamespaceChanges(newNs, oldNs *Namespace) error {
	if newNs.Name != oldNs.Name {
		return validationError(newNs, "changing namespaces between versions is not yet supported")
	}

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
	for _, oldProt := range oldNs.Protocols {
		oldProts[oldProt.Name] = oldProt
	}

	for _, newProt := range newNs.Protocols {
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
func annotateChangedTypeDefinition(newTd, oldTd TypeDefinition) (TypeDefinition, error) {
	switch newNode := newTd.(type) {
	case *RecordDefinition:
		oldTd, ok := oldTd.(*RecordDefinition)
		if !ok {
			return oldTd, fmt.Errorf("changing '%s' to a Record is not backward compatible", newNode.Name)
		}
		res, err := annotateChangedRecordDefinition(newNode, oldTd)
		if err != nil {
			return res, err
		}
		if res != nil {
			return res, nil
		}
		return nil, nil

	case *NamedType:
		oldTd, ok := oldTd.(*NamedType)
		if !ok {
			return oldTd, fmt.Errorf("changing '%s' to a named type is not backward compatible", newNode.Name)
		}
		typeChange := detectChangedTypes(newNode.Type, oldTd.Type)
		oldTd.Annotations[ChangeAnnotationKey] = typeChange
		if typeChange != nil {
			return oldTd, nil
		}
		return nil, nil

	case *EnumDefinition:
		oldTd, ok := oldTd.(*EnumDefinition)
		if !ok {
			return oldTd, fmt.Errorf("changing '%s' to an Enum is not backward compatible", newNode.Name)
		}
		res, err := annotateChangedEnumDefinitions(newNode, oldTd)
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
		}

		if i > len(oldProtocol.Sequence) {
			// CHANGE: Reordered ProtocolSteps
			return oldProtocol, fmt.Errorf("reordering Protocol steps is not backward compatible")
		}
		if newStep.Name != oldProtocol.Sequence[i].Name {
			// CHANGE: Reordered/Renamed ProtocolSteps
			return oldProtocol, fmt.Errorf("reordering or renaming Protocol steps is not backward compatible")
		}
	}

	for _, oldStep := range oldProtocol.Sequence {
		newStep, ok := newSequence[oldStep.Name]
		if !ok {
			oldStep.Annotations[FieldOrStepRemovedAnnotationKey] = true
			return oldProtocol, fmt.Errorf("removing Protocol steps is not backward compatible")
		}

		typeChange := detectChangedTypes(newStep.Type, oldStep.Type)
		if typeChange != nil {
			changed = true
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

		typeChange := detectChangedTypes(newField.Type, oldField.Type)
		if typeChange != nil {
			changed = true
		}
		oldRecord.Fields[i].Annotations[ChangeAnnotationKey] = typeChange
	}

	if changed {
		return oldRecord, nil
	}
	return nil, nil
}

func annotateChangedEnumDefinitions(newNode, oldEnum *EnumDefinition) (*EnumDefinition, error) {
	changed := false

	if newNode.Name != oldEnum.Name {
		// CHANGE: Renamed Enum
		changed = true
	}
	if newNode.IsFlags != oldEnum.IsFlags {
		// CHANGE: Changed Enum to Flags or vice versa
		changed = true
	}

	if oldEnum.BaseType != nil {
		if newNode.BaseType == nil {
			// CHANGE: Removed enum base type?
			changed = true
		}
		if ch := detectChangedTypes(newNode.BaseType, oldEnum.BaseType); ch != nil {
			// CHANGE: Changed Enum base type
			return oldEnum, fmt.Errorf("changing '%s' base type is not backward compatible", newNode.Name)
		}
	} else {
		if newNode.BaseType != nil {
			// CHANGE: Added an enum base type?
			changed = true
		}
	}

	for i, newEnumValue := range newNode.Values {
		oldEnumValue := oldEnum.Values[i]

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
		return oldEnum, nil
	}

	return nil, nil
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
	if ch, ok := newDef.GetDefinitionMeta().Annotations[ChangeAnnotationKey]; ok && ch != nil {
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
		// CHANGE: Changed a non-Optional to an Optional
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
		// CHANGE: Changed a non-Union to a Union
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if len(newType.Cases) == len(oldType.Cases) {
		// TODO: Determine if newType and oldType Union types are an equal set
		// TODO: Account for reordering the Union types
		// return &TypeChangeUnionTypesChange{TypePair{oldType, newType}, mappings}
	} else if len(newType.Cases) > len(oldType.Cases) {
		// TODO: Determine if the oldType Union types are a subset of the newType Union types
		// TODO: Account for reordering the Union type
		// return &TypeChangeUnionGrow{TypePair{oldType, newType}, indexes}
	} else if len(newType.Cases) < len(oldType.Cases) {
		// TODO: Determine if the newType Union types are a subset of the oldType Union types
		// TODO: Account for reordering the Union type
		// return &TypeChangeUnionShrink{TypePair{oldType, newType}, indexes}
	}
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
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

			if newDimension.Length != nil && newDimension.Length != oldDimension.Length {
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
