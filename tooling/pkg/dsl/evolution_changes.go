// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import "fmt"

type TypeChange interface {
	OldType() Type
	NewType() Type
	Inverse() TypeChange
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

func (tc *TypeChangeUnionTypesetChanged) HasTypesAdded() bool {
	for _, match := range tc.NewMatches {
		if !match {
			return true
		}
	}
	return false
}

func (tc *TypeChangeUnionTypesetChanged) HasTypesRemoved() bool {
	for _, match := range tc.OldMatches {
		if !match {
			return true
		}
	}
	return false
}

type TypeChangeStreamTypeChanged struct {
	TypePair
	InnerChange TypeChange
}

func (tc *TypeChangeStreamTypeChanged) Inverse() TypeChange {
	return &TypeChangeStreamTypeChanged{tc.Swap(), tc.InnerChange.Inverse()}
}

type TypeChangeVectorTypeChanged struct {
	TypePair
	InnerChange TypeChange
}

func (tc *TypeChangeVectorTypeChanged) Inverse() TypeChange {
	return &TypeChangeVectorTypeChanged{tc.Swap(), tc.InnerChange.Inverse()}
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
	StepChanges    []TypeChange
	StepMapping    []int
}

func (pc *ProtocolChange) HasReorderedSteps() bool {
	for i, index := range pc.StepMapping {
		if index >= 0 && index != i {
			return true
		}
	}
	return false
}

type RecordChange struct {
	DefinitionPair
	FieldsAdded  []*Field
	FieldRemoved []bool
	FieldChanges []TypeChange
	FieldMapping []int
}

func (rc *RecordChange) HasReorderedFields() bool {
	for i, index := range rc.FieldMapping {
		if index != i {
			return true
		}
	}
	return false
}

type EnumChange struct {
	DefinitionPair
	BaseTypeChange TypeChange
}

func typeChangeIsError(tc TypeChange) bool {
	switch tc := tc.(type) {
	case *TypeChangeStreamTypeChanged:
		// A Stream's Type can only change if it is a changed TypeDefinition
		if _, ok := tc.InnerChange.(*TypeChangeDefinitionChanged); !ok {
			return true
		}
	case *TypeChangeVectorTypeChanged:
		// A Vector's Type can only change if it is a changed TypeDefinition
		if _, ok := tc.InnerChange.(*TypeChangeDefinitionChanged); !ok {
			return true
		}
	case *TypeChangeOptionalTypeChanged:
		return typeChangeIsError(tc.InnerChange)

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
	case *TypeChangeNumberToNumber:
		return "may result in numeric overflow or loss of precision"
	case *TypeChangeNumberToString, *TypeChangeScalarToOptional, *TypeChangeScalarToUnion, *TypeChangeOptionalToUnion:
		return "may produce runtime errors when serializing previous versions"
	case *TypeChangeStringToNumber, *TypeChangeOptionalToScalar, *TypeChangeUnionToScalar, *TypeChangeUnionToOptional:
		return "may produce runtime errors when deserializing previous versions"
	case *TypeChangeUnionTypesetChanged:
		if tc.HasTypesRemoved() {
			return "may produce runtime errors when deserializing previous versions"
		}
		if tc.HasTypesAdded() {
			return "may produce runtime errors when serializing previous versions"
		}
	case *TypeChangeStreamTypeChanged:
		return typeChangeWarningReason(tc.InnerChange)
	case *TypeChangeVectorTypeChanged:
		return typeChangeWarningReason(tc.InnerChange)
	case *TypeChangeOptionalTypeChanged:
		return typeChangeWarningReason(tc.InnerChange)
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
