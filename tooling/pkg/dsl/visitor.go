// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"log"
	"reflect"
)

// Visits a Node tree from the given root.
func Visit(root Node, visitorFunc VisitorFunc) {
	wrapper := Visitor{visitorFunc: visitorFunc}
	wrapper.visitorWithContext = VisitorWithContext[any]{
		func(self VisitorWithContext[any], node Node, context any) {
			visitorFunc(wrapper, node)
		}}

	wrapper.Visit(root)
}

type VisitorFunc func(self Visitor, node Node)

// Visits a Node tree from the given root, threading a context parameter throughout.
func VisitWithContext[T any](root Node, context T, visitor VisitorWithContextFunc[T]) {
	visitorWithContext := VisitorWithContext[T]{visitorWithContext: visitor}
	visitor(visitorWithContext, root, context)
}

type VisitorWithContextFunc[T any] func(self VisitorWithContext[T], node Node, context T)

type Visitor struct {
	visitorFunc        VisitorFunc
	visitorWithContext VisitorWithContext[any]
}

func (visitor Visitor) Visit(node Node) {
	visitor.visitorFunc(visitor, node)
}

func (visitor Visitor) VisitChildren(node Node) {
	visitor.visitorWithContext.VisitChildren(node, nil)
}

type VisitorWithContext[T any] struct {
	visitorWithContext VisitorWithContextFunc[T]
}

func (visitor VisitorWithContext[T]) Visit(node Node, context T) {
	visitor.visitorWithContext(visitor, node, context)
}

func (visitor VisitorWithContext[T]) VisitChildren(node Node, context T) {
	switch t := node.(type) {
	case *Environment:
		for _, ns := range t.Namespaces {
			visitor.Visit(ns, context)
		}
	case *Namespace:
		for _, t := range t.TypeDefinitions {
			visitor.Visit(t, context)
		}

		for _, p := range t.Protocols {
			visitor.Visit(p, context)
		}
	case *DefinitionMeta:
		break
	case *RecordDefinition:
		visitor.Visit(t.DefinitionMeta, context)
		for _, f := range t.Fields {
			visitor.Visit(f, context)
		}

		for _, v := range t.ComputedFields {
			visitor.Visit(v, context)
		}
	case *Vector:
		break
	case *Array:
		if t.HasKnownNumberOfDimensions() {
			for _, d := range *t.Dimensions {
				visitor.Visit(d, context)
			}
		}
	case *ArrayDimension:
		break
	case *Stream:
		break
	case *EnumDefinition:
		visitor.Visit(t.DefinitionMeta, context)
		if t.BaseType != nil {
			visitor.Visit(t.BaseType, context)
		}
		for _, v := range t.Values {
			visitor.Visit(v, context)
		}
	case *EnumValue:
		break
	case *NamedType:
		visitor.Visit(t.DefinitionMeta, context)
		visitor.Visit(t.Type, context)
	case *ProtocolDefinition:
		visitor.Visit(t.DefinitionMeta, context)
		for _, step := range t.Sequence {
			visitor.Visit(step, context)
		}
	case *Field:
		visitor.Visit(t.Type, context)
	case *ProtocolStep:
		visitor.Visit(t.Type, context)
	case *GenericTypeParameter:
		break
	case *SimpleType:
		for _, typeArg := range t.TypeArguments {
			visitor.Visit(typeArg, context)
		}
	case *GeneralizedType:
		for _, typeCase := range t.Cases {
			visitor.Visit(typeCase, context)
		}
		if t.Dimensionality != nil {
			visitor.Visit(t.Dimensionality, context)
		}
	case *TypeCase:
		if !t.IsNullType() {
			visitor.Visit(t.Type, context)
		}
	case Dimensionality:
		break
	case PrimitiveDefinition:
		break
	case *ComputedField:
		visitor.Visit(t.Expression, context)
	case *IntegerLiteralExpression:
		break
	case *StringLiteralExpression:
		break
	case *MemberAccessExpression:
		if t.Target != nil {
			visitor.Visit(t.Target, context)
		}
	case *IndexExpression:
		if t.Target != nil {
			visitor.Visit(t.Target, context)
		}
		for _, arg := range t.Arguments {
			visitor.Visit(arg.Value, context)
		}

	case *FunctionCallExpression:
		for _, arg := range t.Arguments {
			visitor.Visit(arg, context)
		}
	case *TypeConversionExpression:
		visitor.Visit(t.Expression, context)
	case *SwitchExpression:
		visitor.Visit(t.Target, context)
		for _, c := range t.Cases {
			visitor.Visit(c, context)
		}
	case *SwitchCase:
		visitor.Visit(t.Pattern, context)
		visitor.Visit(t.Expression, context)
	case *TypePattern:
		if t.Type != nil {
			visitor.Visit(t.Type, context)
		}
	case *DeclarationPattern:
		visitor.Visit(&t.TypePattern, context)
	case *DiscardPattern:
		break

	default:
		log.Panicf("unhandled type %v", reflect.TypeOf(node))
	}
}
