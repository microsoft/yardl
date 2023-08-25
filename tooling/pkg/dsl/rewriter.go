// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import "fmt"

// Rewrites a Node by visiting it from the given root.
// If the given node in an *Environment, the symbol table will be updated.
func Rewrite(node Node, rewriterFunc RewriterFunc) Node {
	wrapper := &Rewriter{rewriterFunc: rewriterFunc}
	wrapper.rewriterWithContext = RewriterWithContext[any]{
		rewriterFunc: func(node Node, context any, self *RewriterWithContext[any]) Node {
			return rewriterFunc(wrapper, node)
		}}

	if env, ok := node.(*Environment); ok {
		wrapper.rewriterWithContext.symbolTable = &env.SymbolTable
	}

	rewritten := wrapper.Rewrite(node)
	if len(wrapper.rewriterWithContext.updatedSymbols) == 0 {
		return rewritten
	}

	if newEnv, ok := rewritten.(*Environment); ok {
		if newEnv == node {
			clone := *newEnv
			newEnv = &clone
		}

		newEnv.SymbolTable = *wrapper.rewriterWithContext.symbolTable
		return newEnv
	}

	return rewritten
}

type RewriterFunc func(self *Rewriter, node Node) Node

// Rewrites a Node by visiting it from the given root, threading a context parameter throughout.
// If the given node in an *Environment, the symbol table will be updated.
func RewriteWithContext[T any](node Node, context T, rewriter RewriterWithContextFunc[T]) Node {
	wrapper := RewriterWithContext[T]{rewriterFunc: rewriter}
	if env, ok := node.(*Environment); ok {
		wrapper.symbolTable = &env.SymbolTable
	}
	rewritten := wrapper.Rewrite(node, context)
	if len(wrapper.updatedSymbols) == 0 {
		return rewritten
	}

	if newEnv, ok := rewritten.(*Environment); ok {
		if newEnv == node {
			clone := *newEnv
			newEnv = &clone
		}

		newEnv.SymbolTable = *wrapper.symbolTable
		return newEnv
	}

	return rewritten
}

type RewriterWithContextFunc[T any] func(node Node, context T, self *RewriterWithContext[T]) Node

type Rewriter struct {
	rewriterFunc        RewriterFunc
	rewriterWithContext RewriterWithContext[any]
}

func (rewriter *Rewriter) Rewrite(node Node) Node {
	rewritten := rewriter.rewriterFunc(rewriter, node)
	updateSymbolTable(node, rewritten, &rewriter.rewriterWithContext)
	return rewritten
}

func (rewriter *Rewriter) DefaultRewrite(node Node) Node {
	return rewriter.rewriterWithContext.DefaultRewrite(node, nil)
}

type RewriterWithContext[T any] struct {
	rewriterFunc   RewriterWithContextFunc[T]
	symbolTable    *SymbolTable
	updatedSymbols map[string]TypeDefinition
}

func (rewriter *RewriterWithContext[T]) Rewrite(node Node, context T) Node {
	rewritten := rewriter.rewriterFunc(node, context, rewriter)
	updateSymbolTable(node, rewritten, rewriter)
	return rewritten
}

func (rewriter *RewriterWithContext[T]) DefaultRewrite(node Node, context T) Node {
	rewritten := defaultRewriteImpl(rewriter, node, context)
	updateSymbolTable(node, rewritten, rewriter)
	return rewritten
}

func defaultRewriteImpl[T any](rewriter *RewriterWithContext[T], node Node, context T) Node {
	switch t := node.(type) {
	case *Environment:
		rewrittenNamespaces := rewriteSlice(t.Namespaces, context, rewriter)
		if rewrittenNamespaces == nil {
			return t
		}

		rewrittenEnv := *t
		rewrittenEnv.Namespaces = rewrittenNamespaces
		return &rewrittenEnv
	case *Namespace:
		rewrittenTypes := rewriteInterfaceSlice(t.TypeDefinitions, context, rewriter)
		rewrittenProtocols := rewriteSlice(t.Protocols, context, rewriter)

		if rewrittenTypes == nil && rewrittenProtocols == nil {
			return t
		}

		rewrittenNamespace := *t
		if rewrittenTypes != nil {
			rewrittenNamespace.TypeDefinitions = rewrittenTypes
		}
		if rewrittenProtocols != nil {
			rewrittenNamespace.Protocols = rewrittenProtocols
		}

		return &rewrittenNamespace
	case *DefinitionMeta:
		return t
	case *RecordDefinition:
		rewrittenDimensionMeta := rewriter.Rewrite(t.DefinitionMeta, context)
		rewrittenFields := rewriteSlice(t.Fields, context, rewriter)
		rewrittenComputedFields := rewriteSlice(t.ComputedFields, context, rewriter)

		if rewrittenDimensionMeta == t.DefinitionMeta && rewrittenFields == nil && rewrittenComputedFields == nil {
			return t
		}

		rewrittenRecord := *t
		rewrittenRecord.DefinitionMeta = rewrittenDimensionMeta.(*DefinitionMeta)
		if rewrittenFields != nil {
			rewrittenRecord.Fields = rewrittenFields
		}
		if rewrittenComputedFields != nil {
			rewrittenRecord.ComputedFields = rewrittenComputedFields
		}
		return &rewrittenRecord
	case *Vector:
		return t
	case *Array:
		if t.Dimensions == nil {
			return t
		}

		rewrittenDimensions := rewriteSlice(*t.Dimensions, context, rewriter)
		if rewrittenDimensions == nil {
			return t
		}

		rewrittenArray := *t
		rewrittenArray.Dimensions = &rewrittenDimensions
		return &rewrittenArray
	case *ArrayDimension:
		return t
	case *Stream:
		return t
	case *Map:
		rewrittenKeyType := rewriter.Rewrite(t.KeyType, context)
		if rewrittenKeyType == t.KeyType {
			return t
		}

		rewrittenMap := *t
		rewrittenMap.KeyType = rewrittenKeyType.(Type)
		return &rewrittenMap
	case *EnumDefinition:
		rewrittenDimensionMeta := rewriter.Rewrite(t.DefinitionMeta, context)
		rewrittenBaseType := t.BaseType
		if t.BaseType != nil {
			rewriter.Rewrite(t.BaseType, context)
		}
		rewrittenValues := rewriteSlice(t.Values, context, rewriter)

		if t.DefinitionMeta == rewrittenDimensionMeta && t.BaseType == rewrittenBaseType && rewrittenValues == nil {
			return t
		}

		rewrittenEnum := *t
		rewrittenEnum.DefinitionMeta = rewrittenDimensionMeta.(*DefinitionMeta)
		rewrittenEnum.BaseType = rewrittenBaseType
		if rewrittenValues != nil {
			rewrittenEnum.Values = rewrittenValues
		}
		return &rewrittenEnum
	case *EnumValue:
		return t
	case *NamedType:
		rewrittenDimensionMeta := rewriter.Rewrite(t.DefinitionMeta, context)
		rewrittenType := rewriter.Rewrite(t.Type, context)
		if rewrittenDimensionMeta == t.DefinitionMeta && rewrittenType == t.Type {
			return t
		}
		rewrittenNamedType := *t
		rewrittenNamedType.DefinitionMeta = rewrittenDimensionMeta.(*DefinitionMeta)
		rewrittenNamedType.Type = rewrittenType.(Type)
		return &rewrittenNamedType
	case *ProtocolDefinition:
		rewrittenDimensionMeta := rewriter.Rewrite(t.DefinitionMeta, context)
		rewrittenSteps := rewriteSlice(t.Sequence, context, rewriter)

		if t.DefinitionMeta == rewrittenDimensionMeta && rewrittenSteps == nil {
			return t
		}

		rewrittenProtocol := *t
		rewrittenProtocol.DefinitionMeta = rewrittenDimensionMeta.(*DefinitionMeta)
		if rewrittenSteps != nil {
			rewrittenProtocol.Sequence = rewrittenSteps
		}
		return &rewrittenProtocol
	case *Field:
		rewrittenType := rewriter.Rewrite(t.Type, context)
		if rewrittenType == t.Type {
			return t
		}
		rewrittenField := *t
		rewrittenField.Type = rewrittenType.(Type)
		return &rewrittenField
	case *ProtocolStep:
		rewrittenType := rewriter.Rewrite(t.Type, context)
		if rewrittenType == t.Type {
			return t
		}
		rewrittenStep := *t
		rewrittenStep.Type = rewrittenType.(Type)
		return &rewrittenStep
	case *GenericTypeParameter:
		return t
	case *SimpleType:
		rewrittenTypeArguments := rewriteInterfaceSlice(t.TypeArguments, context, rewriter)
		if rewrittenTypeArguments == nil {
			return updateTypeRefence(rewriter, t)
		}

		rewrittenSimpleType := *t
		rewrittenSimpleType.TypeArguments = rewrittenTypeArguments
		return updateTypeRefence(rewriter, &rewrittenSimpleType)
	case *GeneralizedType:
		rewrittenTypeCases := rewriteInterfaceSlice(t.Cases, context, rewriter)

		var rewrittenDimensionality Dimensionality
		if t.Dimensionality != nil {
			rewrittenDimensionality = rewriter.Rewrite(t.Dimensionality, context).(Dimensionality)
		}

		if rewrittenTypeCases == nil && rewrittenDimensionality == t.Dimensionality {
			return t
		}

		rewrittenType := *t
		rewrittenType.Dimensionality = rewrittenDimensionality
		if rewrittenTypeCases != nil {
			rewrittenType.Cases = rewrittenTypeCases
		}
		return &rewrittenType
	case *TypeCase:
		if t.IsNullType() {
			return t
		}

		rewrittenType := rewriter.Rewrite(t.Type, context)
		if rewrittenType == t.Type {
			return t
		}
		rewrittenTypeCase := *t
		rewrittenTypeCase.Type = rewrittenType.(Type)
		return &rewrittenTypeCase
	case Dimensionality:
		return t
	case PrimitiveDefinition:
		return t
	case *ComputedField:
		rewrittenExpression := rewriter.Rewrite(t.Expression, context)
		if rewrittenExpression == t.Expression {
			return t
		}

		rewrittenField := *t
		rewrittenField.Expression = rewrittenExpression.(Expression)
		return &rewrittenField
	case *IntegerLiteralExpression:
		return t
	case *StringLiteralExpression:
		return t
	case *MemberAccessExpression:
		if t.Target == nil {
			return t
		}

		rewrittenTarget := rewriter.Rewrite(t.Target, context)
		if rewrittenTarget == t.Target {
			return t
		}

		rewrittenExpression := *t
		rewrittenExpression.Target = rewrittenTarget.(Expression)
		return &rewrittenExpression
	case *IndexExpression:
		var rewrittenTarget Expression
		if t.Target != nil {
			rewrittenTarget = rewriter.Rewrite(t.Target, context).(Expression)
		}

		rewrittenArguments := rewriteInterfaceSlice(t.Arguments, context, rewriter)

		if rewrittenTarget == t.Target && rewrittenArguments == nil {
			return t
		}

		rewrittenExpression := *t
		rewrittenExpression.Target = rewrittenTarget
		if rewrittenArguments != nil {
			rewrittenExpression.Arguments = rewrittenArguments
		}
		return &rewrittenExpression
	case *IndexArgument:
		rewrittenValue := rewriter.Rewrite(t.Value, context)
		if rewrittenValue == t.Value {
			return t
		}

		rewrittenArgument := *t
		rewrittenArgument.Value = rewrittenValue.(Expression)
		return &rewrittenArgument

	case *FunctionCallExpression:
		rewrittenArguments := rewriteInterfaceSlice(t.Arguments, context, rewriter)

		if rewrittenArguments == nil {
			return t
		}

		rewrittenExpression := *t
		rewrittenExpression.Arguments = rewrittenArguments
		return &rewrittenExpression

	case *TypeConversionExpression:
		rewrittenTarget := rewriter.Rewrite(t.Expression, context)
		if rewrittenTarget == t.Expression {
			return t
		}

		rewrittenExpression := *t
		rewrittenExpression.Expression = rewrittenTarget.(Expression)
		return &rewrittenExpression
	case *SwitchExpression:
		rewrittenTarget := rewriter.Rewrite(t.Target, context)
		rewrittenCases := rewriteSlice(t.Cases, context, rewriter)

		if rewrittenTarget == t.Target && rewrittenCases == nil {
			return t
		}

		rewrittenSwitch := *t
		rewrittenSwitch.Target = rewrittenTarget.(Expression)
		if rewrittenCases != nil {
			rewrittenSwitch.Cases = rewrittenCases
		}
		return &rewrittenSwitch
	case *SwitchCase:
		rewrittenPattern := rewriter.Rewrite(t.Pattern, context)
		rewrittenExpression := rewriter.Rewrite(t.Expression, context)

		if rewrittenPattern == t.Pattern && rewrittenExpression == t.Expression {
			return t
		}

		rewrittenCase := *t
		rewrittenCase.Pattern = rewrittenPattern.(Pattern)
		rewrittenCase.Expression = rewrittenExpression.(Expression)
		return &rewrittenCase
	case *TypePattern:
		if t.Type == nil {
			return t
		}
		rewrittenType := rewriter.Rewrite(t.Type, context)
		if rewrittenType == t.Type {
			return t
		}

		rewrittenPattern := *t
		rewrittenPattern.Type = rewrittenType.(Type)
		return &rewrittenPattern
	case *DeclarationPattern:
		rewrittenType := rewriter.Rewrite(&t.TypePattern, context)
		if rewrittenType == &t.TypePattern {
			return t
		}

		rewrittenPattern := *t
		rewrittenPattern.TypePattern = *rewrittenType.(*TypePattern)
		return &rewrittenPattern
	case *DiscardPattern:
		return t
	default:
		panic(fmt.Sprintf("unhandled type %T", t))
	}
}

func updateTypeRefence[T any](rewriter *RewriterWithContext[T], t *SimpleType) *SimpleType {
	if rewriter.updatedSymbols == nil {
		return t
	}
	updatedDefinition := rewriter.updatedSymbols[t.Name]
	if updatedDefinition == nil {
		return t
	}

	newType := *t
	if len(t.TypeArguments) > 0 {
		genericType, err := MakeGenericType(updatedDefinition, t.TypeArguments, false)
		if err != nil {
			panic(fmt.Sprintf("failed to make generic rewritten type: %s", err))
		}
		newType.ResolvedDefinition = genericType
	} else {
		newType.ResolvedDefinition = updatedDefinition
	}

	return &newType
}

func updateSymbolTable[T any](original, rewritten Node, rewriterWithContext *RewriterWithContext[T]) {
	if rewriterWithContext.symbolTable == nil {
		return
	}

	if rewritten != original {
		if td, ok := rewritten.(TypeDefinition); ok {
			meta := td.GetDefinitionMeta()
			if len(meta.TypeArguments) > 0 {
				// This is a generic type with arguments (not a generic type definition),
				// so it does not belong in the symbol table
				return
			}
			name := meta.GetQualifiedName()
			if _, wasFound := (*rewriterWithContext.symbolTable)[name]; wasFound {
				if len(rewriterWithContext.updatedSymbols) == 0 {
					rewriterWithContext.updatedSymbols = make(map[string]TypeDefinition)
					newSymbolTable := rewriterWithContext.symbolTable.Clone()
					rewriterWithContext.symbolTable = &newSymbolTable
				}

				(*rewriterWithContext.symbolTable)[name] = td
				rewriterWithContext.updatedSymbols[name] = td
			}
		}
	}
}

// Rewrites a slice of pointers to types that implement the Node interface, e.g, []*Field
// Returns nil if no changes were made and the original slice should be used.
func rewriteSlice[TContext any, TElement any, T interface {
	*TElement
	Node
}, TSlice interface {
	~[]T
}](slice TSlice, context TContext, rewriter *RewriterWithContext[TContext]) TSlice {
	var rewrittenSlice []T
	for i, element := range slice {
		visited := rewriter.Rewrite(T(element), context)
		if visited != element && rewrittenSlice == nil {
			rewrittenSlice = make([]T, len(slice))
			copy(rewrittenSlice, slice)
		}
		if rewrittenSlice != nil {
			rewrittenSlice[i] = visited.(T)
		}
	}

	return rewrittenSlice
}

// Rewrites a slice of an interface that implements the Node interface, e.g, []Expression
// Returns nil if no changes were made and the original slice should be used.
func rewriteInterfaceSlice[TContext any, T Node](slice []T, context TContext, rewriter *RewriterWithContext[TContext]) []T {
	var rewrittenSlice []T
	for i, element := range slice {
		visited := rewriter.Rewrite(T(element), context)
		if visited != Node(element) && rewrittenSlice == nil {
			rewrittenSlice = make([]T, len(slice))
			copy(rewrittenSlice, slice)
		}
		if rewrittenSlice != nil {
			rewrittenSlice[i] = visited.(T)
		}
	}

	return rewrittenSlice
}
