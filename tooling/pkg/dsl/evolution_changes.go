// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strings"
)

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

type TypeChangeDefinitionChanged struct {
	TypePair
	DefinitionChange
}

func (tc *TypeChangeDefinitionChanged) Inverse() TypeChange {
	return &TypeChangeDefinitionChanged{tc.Swap(), tc.DefinitionChange}
}

type TypeChangeIncompatible struct{ TypePair }

func (tc *TypeChangeIncompatible) Inverse() TypeChange {
	return &TypeChangeIncompatible{tc.Swap()}
}

type TypeChangeStepAdded struct{ TypePair }

func (tc *TypeChangeStepAdded) Inverse() TypeChange {
	return nil
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
	_ TypeChange = (*TypeChangeStepAdded)(nil)

	_ DefinitionChange = (*DefinitionChangeIncompatible)(nil)
	_ DefinitionChange = (*NamedTypeChange)(nil)
	_ DefinitionChange = (*RecordChange)(nil)
	_ DefinitionChange = (*EnumChange)(nil)
	_ DefinitionChange = (*ProtocolChange)(nil)
	_ DefinitionChange = (*CompatibilityChange)(nil)

	_ DefinitionChange = (*PrimitiveChangeNumberToNumber)(nil)
	_ DefinitionChange = (*PrimitiveChangeNumberToString)(nil)
	_ DefinitionChange = (*PrimitiveChangeStringToNumber)(nil)
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
	Reason string
}

const (
	IncompatibleDefinitions     = "definitions are incompatible"
	IncompatibleTypeParameters  = "type parameters do not match"
	IncompatibleBaseDefinitions = "base definitions are incompatible"
)

type NamedTypeChange struct {
	DefinitionPair
	TypeChange TypeChange
}

type ProtocolRemoved struct {
	DefinitionPair
}

type ProtocolChange struct {
	DefinitionPair
	PreviousSchema string
	StepsRemoved   []*ProtocolStep
	StepChanges    []TypeChange
	StepsReordered []*ProtocolStep
}

type RecordChange struct {
	DefinitionPair
	FieldsAdded   []*Field
	FieldRemoved  []bool
	FieldChanges  []TypeChange
	NewFieldIndex []int
}

type EnumChange struct {
	DefinitionPair
	BaseTypeChange TypeChange
	ValuesAdded    []string
	ValuesRemoved  []string
	ValuesChanged  []string
}

type CompatibilityChange struct {
	DefinitionPair
}

type PrimitiveChangeNumberToNumber struct{ DefinitionPair }
type PrimitiveChangeNumberToString struct{ DefinitionPair }
type PrimitiveChangeStringToNumber struct{ DefinitionPair }

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
	case *TypeChangeNumberToString:
		return fmt.Sprintf("will result in a write error if its value cannot be converted to type '%s' at runtime", TypeToShortSyntax(tc.OldType(), true))
	case *TypeChangeStringToNumber:
		return fmt.Sprintf("will result in a read error if its value cannot be converted to type '%s' at runtime", TypeToShortSyntax(tc.NewType(), true))

	case *TypeChangeScalarToOptional:
		return fmt.Sprintf("will result in writing the default zero value for '%s' if it does not have a value at runtime", TypeToShortSyntax(tc.OldType(), true))
	case *TypeChangeOptionalToScalar:
		return fmt.Sprintf("will result in reading the default zero value for '%s' if it does not have a value at runtime", TypeToShortSyntax(tc.NewType(), true))

	case *TypeChangeScalarToUnion:
		return fmt.Sprintf("will result in a write error if its value is not of type '%s' at runtime", TypeToShortSyntax(tc.OldType(), true))
	case *TypeChangeUnionToScalar:
		return fmt.Sprintf("will result in a read error if its value is not of type '%s' at runtime", TypeToShortSyntax(tc.NewType(), true))

	case *TypeChangeOptionalToUnion:
		oldInnerType := tc.OldType().(*GeneralizedType).Cases[1].Type
		return fmt.Sprintf("will result in a write error if its value is not of type '%s' at runtime", TypeToShortSyntax(oldInnerType, true))
	case *TypeChangeUnionToOptional:
		newInnerType := tc.NewType().(*GeneralizedType).Cases[1].Type
		return fmt.Sprintf("will result in a read error if its value is not of type '%s' at runtime", TypeToShortSyntax(newInnerType, true))

	case *TypeChangeUnionTypesetChanged:
		var removed []string
		for i, match := range tc.OldMatches {
			if !match {
				syntax := TypeToShortSyntax(tc.OldType().(*GeneralizedType).Cases[i].Type, true)
				removed = append(removed, fmt.Sprintf("'%s'", syntax))
			}
		}
		if len(removed) > 0 {
			return fmt.Sprintf("may result in a read error if its value at runtime is of type %s", strings.Join(removed, " or "))
		}

		var added []string
		for i, match := range tc.NewMatches {
			if !match {
				syntax := TypeToShortSyntax(tc.NewType().(*GeneralizedType).Cases[i].Type, true)
				added = append(added, fmt.Sprintf("'%s'", syntax))
			}
		}
		if len(added) > 0 {
			return fmt.Sprintf("may result in a write error if its value at runtime is of type %s", strings.Join(added, " or "))
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
