// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"reflect"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

// Very strict comparison - No changes allowed
// Ignores comments and annotations
func ensureNoChanges(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) (*Environment, error) {
	visitor := func(self CompareVisitor[string], newRoot, oldRoot Node, context string) error {
		return self.StrictCompare(newRoot, oldRoot, context)
	}

	if err := Compare[string](env, predecessor, "", visitor); err != nil {
		return env, err
	}
	return env, nil
}

// Returns an error if there is *any* change between oldRoot and newRoot
// This serve[sd] as the base for developing the model evolution approach
func (cv CompareVisitor[T]) StrictCompare(newRoot, oldRoot Node, context T) error {
	// log.Info().Msgf("Comparing %s and %s", newRoot.GetNodeMeta(), oldRoot.GetNodeMeta())

	/*
		**TODO**:

		- Determine what are backward compatible changes and add support for them
		- Write unit tests
		- Could wrap errors (e.g. ComparisonError) to add context for informing user of changes
	*/

	// We're comparing new, working changes to a previous, established schema.
	// Old is the reference from the user's perspective.
	switch oldNode := oldRoot.(type) {
	case *Environment:
		newNode, ok := newRoot.(*Environment)
		if !ok {
			panic(fmt.Sprintf("expected a %T", oldNode))
		}

		return StrictCompareSlices(newNode.Namespaces, oldNode.Namespaces, context, cv)

	case *Namespace:
		newNode, ok := newRoot.(*Namespace)
		if !ok {
			panic(fmt.Sprintf("expected a %T", oldNode))
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename namespace %s to %s", oldNode.Name, newNode.Name)
		}

		if err := StrictCompareSlices(newNode.TypeDefinitions, oldNode.TypeDefinitions, context, cv); err != nil {
			return err
		}

		return StrictCompareSlices(newNode.Protocols, oldNode.Protocols, context, cv)

	case *DefinitionMeta:
		newNode, ok := newRoot.(*DefinitionMeta)
		if !ok {
			return validationError(newRoot, "expected a definition meta")
		}

		// These should be redundant
		if newNode.Namespace != oldNode.Namespace {
			return validationError(newRoot, "cannot change definition namespace")
		}
		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename definition %s to %s", oldNode.Name, newNode.Name)
		}

		if err := StrictCompareSlices(newNode.TypeParameters, oldNode.TypeParameters, context, cv); err != nil {
			return err
		}
		return StrictCompareSlices(newNode.TypeArguments, oldNode.TypeArguments, context, cv)

	case *RecordDefinition:
		newNode, ok := newRoot.(*RecordDefinition)
		if !ok {
			return validationError(newRoot, "expected a record definition")
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename record %s to %s", oldNode.Name, newNode.Name)
		}

		if err := cv.Compare(newNode.DefinitionMeta, oldNode.DefinitionMeta, context); err != nil {
			return err
		}

		if err := StrictCompareSlices(newNode.Fields, oldNode.Fields, context, cv); err != nil {
			return err
		}

		return StrictCompareSlices(newNode.ComputedFields, oldNode.ComputedFields, context, cv)

	case *Field:
		newNode, ok := newRoot.(*Field)
		if !ok {
			return validationError(newRoot, "expected a field")
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename field %s to %s", oldNode.Name, newNode.Name)
		}

		return cv.Compare(newNode.Type, oldNode.Type, context)

	case *ProtocolDefinition:
		newNode, ok := newRoot.(*ProtocolDefinition)
		if !ok {
			return validationError(newRoot, "expected a protocol definition")
		}
		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename protocol %s to %s", oldNode.Name, newNode.Name)
		}

		if err := cv.Compare(newNode.DefinitionMeta, oldNode.DefinitionMeta, context); err != nil {
			return err
		}

		return StrictCompareSlices(newNode.Sequence, oldNode.Sequence, context, cv)

	case *ProtocolStep:
		newNode, ok := newRoot.(*ProtocolStep)
		if !ok {
			return validationError(newRoot, "expected a protocol step")
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename protocol step %s to %s", oldNode.Name, newNode.Name)
		}

		return cv.Compare(newNode.Type, oldNode.Type, context)

	case *EnumDefinition:
		newNode, ok := newRoot.(*EnumDefinition)
		if !ok {
			return validationError(newRoot, "expected an enum definition")
		}
		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename enum %s to %s", oldNode.Name, newNode.Name)
		}
		// if err := cv.Compare(newNode.DefinitionMeta, oldNode.DefinitionMeta, context); err != nil {
		// 	return err
		// }
		if newNode.IsFlags != oldNode.IsFlags {
			return validationError(newRoot, "cannot change enum to flags or vice versa")
		}

		if oldNode.BaseType != nil {
			if newNode.BaseType == nil {
				return validationError(newRoot, "cannot change enum base type")
			}
			if err := cv.Compare(newNode.BaseType, oldNode.BaseType, context); err != nil {
				return err
			}
		} else {
			if newNode.BaseType != nil {
				return validationError(newRoot, "cannot change enum base type")
			}
		}

		return StrictCompareSlices(newNode.Values, oldNode.Values, context, cv)

	case PrimitiveDefinition:
		newNode, ok := newRoot.(PrimitiveDefinition)
		if !ok {
			return validationError(newRoot, "expected a primitive definition")
		}
		if newNode != oldNode {
			return validationError(newRoot, "cannot change primitive type %s to %s", oldNode, newNode)
		}
		return nil

	case *NamedType:
		newNode, ok := newRoot.(*NamedType)
		if !ok {
			return validationError(newRoot, "expected a named type")
		}
		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename type %s to %s", oldNode.Name, newNode.Name)
		}
		if err := cv.Compare(newNode.DefinitionMeta, oldNode.DefinitionMeta, context); err != nil {
			return err
		}
		return cv.Compare(newNode.Type, oldNode.Type, context)

	case *SimpleType:
		newNode, ok := newRoot.(*SimpleType)
		if !ok {
			return validationError(newRoot, "expected a simple type")
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot change type %s to %s", oldNode.Name, newNode.Name)
		}

		// Don't compare ResolveDefinition. This is handled by comparing Namespace.TypeDefinitions
		// return cv.Compare(newNode.ResolvedDefinition, oldNode.ResolvedDefinition, context)
		return nil

	case *GeneralizedType:
		newNode, ok := newRoot.(*GeneralizedType)
		if !ok {
			// TODO: Clarify in error msg what kind of generalized type (and dimensionality)
			return validationError(newRoot, "expected a generalized type")
		}

		if err := StrictCompareSlices(newNode.Cases, oldNode.Cases, context, cv); err != nil {
			return err
		}

		switch oldDim := oldNode.Dimensionality.(type) {
		case nil, *Stream:
			return nil
		case *Vector:
			newDim, ok := newNode.Dimensionality.(*Vector)
			if !ok {
				return validationError(newRoot, "expected a vector")
			}
			if oldDim.Length == nil {
				if newDim.Length != nil {
					return validationError(newRoot, "cannot add vector length")
				}
				return nil
			}
			if newDim.Length == nil {
				return validationError(newRoot, "cannot remove vector length")
			}
			if *newDim.Length != *oldDim.Length {
				return validationError(newRoot, "cannot change vector length")
			}
			return nil

		case *Array:
			newDim, ok := newNode.Dimensionality.(*Array)
			if !ok {
				return validationError(newRoot, "expected an array")
			}
			if oldDim.Dimensions == nil {
				if newDim.Dimensions != nil {
					return validationError(newRoot, "cannot add dimensions to array")
				}
				return nil
			}
			if newDim.Dimensions == nil {
				return validationError(newRoot, "cannot remove dimensions from array")
			}

			if len(*newDim.Dimensions) != len(*oldDim.Dimensions) {
				return validationError(newRoot, "mismatch in number of array dimensions")
			}

			for i := range *oldDim.Dimensions {
				if err := cv.Compare((*newDim.Dimensions)[i], (*oldDim.Dimensions)[i], context); err != nil {
					return err
				}
			}

			return StrictCompareSlices(*newDim.Dimensions, *oldDim.Dimensions, context, cv)

		case *Map:
			dimA, ok := newNode.Dimensionality.(*Map)
			if !ok {
				return validationError(newRoot, "expected a map")
			}
			return cv.Compare(dimA.KeyType, oldDim.KeyType, context)
		default:
			log.Panic().Msgf("unhandled type %v", reflect.TypeOf(newRoot))
		}

	case *TypeCase:
		newNode, ok := newRoot.(*TypeCase)
		if !ok {
			return validationError(newRoot, "expected a type case")
		}

		if oldNode.Type == nil {
			if newNode.Type != nil {
				return validationError(newNode.Type, "cannot change type case type")
			}
			return nil
		}
		if newNode.Type == nil {
			return validationError(newNode.Type, "cannot remove type case type")
		}

		if newNode.Tag != oldNode.Tag {
			return validationError(newRoot, "cannot change type case tag")
		}

		return cv.Compare(newNode.Type, oldNode.Type, context)

	case *ArrayDimension:
		newNode, ok := newRoot.(*ArrayDimension)
		if !ok {
			return validationError(newRoot, "expected an array dimension")
		}

		if oldNode.Name == nil {
			if newNode.Name != nil {
				return validationError(newRoot, "cannot add array dimension name")
			}
			return nil
		}
		if newNode.Name == nil {
			return validationError(newRoot, "cannot remove array dimension name")
		}
		if *newNode.Name != *oldNode.Name {
			return validationError(newRoot, "cannot rename array dimension %s to %s", *oldNode.Name, *newNode.Name)
		}

		if oldNode.Length == nil {
			if newNode.Length != nil {
				return validationError(newRoot, "cannot add array dimension length")
			}
			return nil
		}
		if newNode.Length == nil {
			return validationError(newRoot, "cannot remove array dimension length")
		}
		if *newNode.Length != *oldNode.Length {
			return validationError(newRoot, "cannot change array dimension length")
		}

		return nil

	case *EnumValue:
		newNode, ok := newRoot.(*EnumValue)
		if !ok {
			return validationError(newRoot, "expected an enum value")
		}
		if newNode.Symbol != oldNode.Symbol {
			return validationError(newRoot, "cannot change enum value symbol")
		}
		if newNode.IntegerValue.Cmp(&oldNode.IntegerValue) != 0 {
			return validationError(newRoot, "cannot change enum value integer value")
		}
		return nil

	case *GenericTypeParameter:
		newNode, ok := newRoot.(*GenericTypeParameter)
		if !ok {
			return validationError(newRoot, "expected a generic type parameter")
		}
		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename generic type parameter %s to %s", oldNode.Name, newNode.Name)
		}
		return nil

	case *ComputedField:
		newNode, ok := newRoot.(*ComputedField)
		if !ok {
			return validationError(newRoot, "expected a computed field")
		}

		if newNode.Name != oldNode.Name {
			return validationError(newRoot, "cannot rename computed field %s to %s", oldNode.Name, newNode.Name)
		}

		return cv.Compare(newNode.Expression, oldNode.Expression, context)

	case *UnaryExpression:
		newNode, ok := newRoot.(*UnaryExpression)
		if !ok {
			return validationError(newRoot, "expected a unary expression")
		}
		return cv.Compare(newNode.Expression, oldNode.Expression, context)

	case *BinaryExpression:
		newNode, ok := newRoot.(*BinaryExpression)
		if !ok {
			return validationError(newRoot, "expected a binary expression")
		}
		if err := cv.Compare(newNode.Left, oldNode.Left, context); err != nil {
			return err
		}
		if err := cv.Compare(newNode.Right, oldNode.Right, context); err != nil {
			return err
		}

		return nil

	case *IntegerLiteralExpression:
		newNode, ok := newRoot.(*IntegerLiteralExpression)
		if !ok {
			return validationError(newRoot, "expected an integer literal")
		}
		if newNode.Value.Cmp(&oldNode.Value) != 0 {
			return validationError(newNode, "cannot change integer literal value")
		}
		return nil

	case *FloatingPointLiteralExpression:
		newNode, ok := newRoot.(*FloatingPointLiteralExpression)
		if !ok {
			return validationError(newRoot, "expected a floating point literal")
		}
		if newNode.Value != oldNode.Value {
			return validationError(newNode, "cannot change integer literal value")
		}
		return nil

	case *StringLiteralExpression:
		newNode, ok := newRoot.(*StringLiteralExpression)
		if !ok {
			return validationError(newRoot, "expected a string literal")
		}
		if newNode.Value != oldNode.Value {
			return validationError(newNode, "cannot change integer literal value")
		}
		return nil

	case *MemberAccessExpression:
		newNode, ok := newRoot.(*MemberAccessExpression)
		if !ok {
			return validationError(newRoot, "expected a member access expression")
		}
		if newNode.Kind != oldNode.Kind {
			return validationError(newRoot, "cannot change kind of member access")
		}
		if newNode.Member != oldNode.Member {
			return validationError(newRoot, "cannot change accessed member %s to %s", oldNode.Member, newNode.Member)
		}

		if oldNode.Target == nil {
			if newNode.Target != nil {
				return validationError(newRoot, "cannot change member access target")
			}
			return nil
		}
		return cv.Compare(newNode.Target, oldNode.Target, context)

	case *SubscriptExpression:
		newNode, ok := newRoot.(*SubscriptExpression)
		if !ok {
			return validationError(newRoot, "expected a subscript expression")
		}

		if len(newNode.Arguments) != len(oldNode.Arguments) {
			return validationError(newNode.Arguments[0], "mismatch in number of subscript arguments")
		}
		for i := range oldNode.Arguments {
			if err := cv.Compare(newNode.Arguments[i], oldNode.Arguments[i], context); err != nil {
				return err
			}
		}

		if err := cv.Compare(newNode.Target, oldNode.Target, context); err != nil {
			return err
		}
		return nil

	case *SubscriptArgument:
		newNode, ok := newRoot.(*SubscriptArgument)
		if !ok {
			return validationError(newRoot, "expected a subscript argument")
		}

		if newNode.Label != oldNode.Label {
			return validationError(newRoot, "cannot change subscript argument kind")
		}

		if err := cv.Compare(newNode.Value, oldNode.Value, context); err != nil {
			return err
		}
		return nil

	case *FunctionCallExpression:
		newNode, ok := newRoot.(*FunctionCallExpression)
		if !ok {
			return validationError(newRoot, "expected a function call expression")
		}
		if newNode.FunctionName != oldNode.FunctionName {
			return validationError(newRoot, "cannot change function name")
		}
		if len(newNode.Arguments) != len(oldNode.Arguments) {
			return validationError(newNode.Arguments[0], "mismatch in number of function call arguments")
		}
		for i := range oldNode.Arguments {
			if err := cv.Compare(newNode.Arguments[i], oldNode.Arguments[i], context); err != nil {
				return err
			}
		}
		return nil

	case *TypeConversionExpression:
		newNode, ok := newRoot.(*TypeConversionExpression)
		if !ok {
			return validationError(newRoot, "expected a type conversion expression")
		}
		if err := cv.Compare(newNode.Type, oldNode.Type, context); err != nil {
			return err
		}
		if err := cv.Compare(newNode.Expression, oldNode.Expression, context); err != nil {
			return err
		}
		return nil

	case *SwitchExpression:
		newNode, ok := newRoot.(*SwitchExpression)
		if !ok {
			return validationError(newRoot, "expected a switch expression")
		}
		if err := cv.Compare(newNode.Target, oldNode.Target, context); err != nil {
			return err
		}
		if len(newNode.Cases) != len(oldNode.Cases) {
			return validationError(newNode.Cases[0], "mismatch in number of switch cases")
		}
		for i := range oldNode.Cases {
			if err := cv.Compare(newNode.Cases[i], oldNode.Cases[i], context); err != nil {
				return err
			}
		}
		return nil

	case *SwitchCase:
		newNode, ok := newRoot.(*SwitchCase)
		if !ok {
			return validationError(newRoot, "expected a switch case")
		}
		if err := cv.Compare(newNode.Pattern, oldNode.Pattern, context); err != nil {
			return err
		}
		if err := cv.Compare(newNode.Expression, oldNode.Expression, context); err != nil {
			return err
		}
		return nil

	case *TypePattern:
		newNode, ok := newRoot.(*TypePattern)
		if !ok {
			return validationError(newRoot, "expected a type pattern")
		}
		if oldNode.Type == nil {
			if newNode.Type != nil {
				return validationError(newRoot, "cannot change type pattern type")
			}
			return nil
		}
		if err := cv.Compare(newNode.Type, oldNode.Type, context); err != nil {
			return err
		}
		return nil

	case *DeclarationPattern:
		newNode, ok := newRoot.(*DeclarationPattern)
		if !ok {
			return validationError(newRoot, "expected a declaration pattern")
		}
		if newNode.Identifier != oldNode.Identifier {
			return validationError(newRoot, "cannot change declaration pattern identifier")
		}
		if oldNode.Type == nil {
			if newNode.Type != nil {
				return validationError(newRoot, "cannot change declaration pattern type")
			}
			return nil
		}
		if err := cv.Compare(newNode.Type, oldNode.Type, context); err != nil {
			return err
		}
		return nil

	case *DiscardPattern:
		_, ok := newRoot.(*DiscardPattern)
		if !ok {
			return validationError(newRoot, "expected a discard pattern")
		}
		return nil

	default:
		log.Panic().Msgf("unhandled type %v", reflect.TypeOf(newRoot))
	}

	log.Panic().Msg("should not reach this")
	return fmt.Errorf("should not reach this")
}

func StrictCompareSlices[T any, N Node](newNodes, oldNodes []N, context T, cv CompareVisitor[T]) error {
	if newNodes == nil && oldNodes == nil {
		return nil
	}

	if newNodes == nil && oldNodes != nil || newNodes != nil && oldNodes == nil || len(oldNodes) != len(newNodes) {
		return fmt.Errorf("mismatch in number of nodes")
	}

	for i := range newNodes {
		if err := cv.Compare(newNodes[i], oldNodes[i], context); err != nil {
			return err
		}
	}
	return nil
}
