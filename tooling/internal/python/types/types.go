// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package types

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteTypes(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)

	common.WriteComment(w, "pyright: reportUnusedImport=false")
	common.WriteComment(w, "pyright: reportUnknownArgumentType=false")
	common.WriteComment(w, "pyright: reportUnknownMemberType=false")
	common.WriteComment(w, "pyright: reportUnknownVariableType=false")

	relativePath := ".."
	if ns.IsTopLevel {
		relativePath = "."
	}

	fmt.Fprintf(w, `
import datetime
import enum
import types
import typing

import numpy as np
import numpy.typing as npt

from %s import yardl_types as yardl
from %s import _dtypes

`, relativePath, relativePath)

	for _, ref := range ns.GetAllChildReferences() {
		fmt.Fprintf(w, "from %s import %s\n", relativePath, common.NamespaceIdentifierName(ref.Name))
	}
	w.WriteStringln("")

	writeTypes(w, st, ns)

	writeGetDTypeFunc(w, ns)

	definitionsPath := path.Join(packageDir, "types.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeTypes(w *formatting.IndentedWriter, st dsl.SymbolTable, ns *dsl.Namespace) {
	writeTypeVars(w, ns)

	unions := make(map[string]any)

	for _, td := range ns.TypeDefinitions {
		writeUnionClasses(w, td, unions)
		switch td := td.(type) {
		case *dsl.EnumDefinition:
			writeEnum(w, td)
		case *dsl.RecordDefinition:
			writeRecord(w, td, st)
		case *dsl.NamedType:
			if _, found := unions[td.Name]; !found {
				writeNamedType(w, td)
			}
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}
	}

	for _, p := range ns.Protocols {
		writeUnionClasses(w, p, unions)
	}
}

func writeTypeVars(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	typeVars := make(map[string]any)
	for _, td := range ns.TypeDefinitions {
		for _, tp := range td.GetDefinitionMeta().TypeParameters {
			identifier := common.TypeIdentifierName(tp.Name)
			if _, ok := typeVars[identifier]; !ok {
				typeVars[identifier] = nil
				fmt.Fprintf(w, "%s = typing.TypeVar(\"%s\")\n", identifier, identifier)
				numpyTypeVarName := common.NumpyTypeParameterSyntax(tp)
				fmt.Fprintf(w, "%s = typing.TypeVar(\"%s\", bound=np.generic)\n", numpyTypeVarName, numpyTypeVarName)
			}
		}
	}
	if len(typeVars) > 0 {
		w.WriteStringln("\n")
	}
}

func writeUnionClasses(w *formatting.IndentedWriter, td dsl.TypeDefinition, unions map[string]any) {
	dsl.Visit(td, func(self dsl.Visitor, node dsl.Node) {
		switch node := node.(type) {
		case *dsl.GeneralizedType:
			if node.Cases.IsUnion() {
				unionClassName, typeParameters := common.UnionClassName(node)
				if _, ok := unions[unionClassName]; !ok {
					if _, isNamedType := td.(*dsl.NamedType); isNamedType {
						// This is a named type defining a union, so we will use the named type's name instead
						unionClassName = td.GetDefinitionMeta().Name
					}
					if len(unions) == 0 {
						w.WriteStringln("_T = typing.TypeVar('_T')\n")
					}
					writeUnionClass(w, unionClassName, typeParameters, node, td.GetDefinitionMeta().Namespace)
					unions[unionClassName] = nil
				}
			}
		}
		self.VisitChildren(node)
	})
}

func writeUnionClass(w *formatting.IndentedWriter, className string, typeParameters string, generalizedType *dsl.GeneralizedType, contextNamespace string) {
	typeCases := generalizedType.Cases
	var baseClassSpec string
	var genericSubscript string
	if len(typeParameters) > 0 {
		baseClassSpec = fmt.Sprintf("(typing.Generic[%s])", typeParameters)
		genericSubscript = fmt.Sprintf("[%s]", typeParameters)
	}

	unionCaseType := fmt.Sprintf("%sUnionCase", className)
	fmt.Fprintf(w, "class %s%s:\n", className, baseClassSpec)
	w.Indented(func() {
		for _, tc := range typeCases {
			if tc.Type == nil {
				continue
			}

			if len(typeParameters) > 0 {
				fmt.Fprintf(w, "%s: type[\"%s[%s, %s]\"]\n", formatting.ToPascalCase(tc.Tag), unionCaseType, typeParameters, common.TypeSyntax(tc.Type, contextNamespace))
			} else {
				fmt.Fprintf(w, "%s: typing.ClassVar[type[\"%s[%s]\"]]\n", formatting.ToPascalCase(tc.Tag), unionCaseType, common.TypeSyntax(tc.Type, contextNamespace))
			}
		}
	})
	w.WriteStringln("")

	fmt.Fprintf(w, "class %s(%s%s, yardl.UnionCase[_T]):\n", unionCaseType, className, genericSubscript)
	w.Indented(func() {
		w.WriteStringln("pass")
	})
	w.WriteStringln("")
	i := 0
	for _, tc := range typeCases {
		if tc.Type == nil {
			continue
		}
		pascalTag := formatting.ToPascalCase(tc.Tag)
		fmt.Fprintf(w, "%s.%s = type(\"%s.%s\", (%s,), {\"index\": %d, \"tag\": \"%s\"})\n", className, pascalTag, className, pascalTag, unionCaseType, i, tc.Tag)
		i++
	}
	fmt.Fprintf(w, "del %s\n", unionCaseType)
	w.WriteStringln("")
}

func writeNamedType(w *formatting.IndentedWriter, td *dsl.NamedType) {
	// Does this NamedType resolve to a RecordDefinition?
	resolvesToRecord := false
	if t, ok := dsl.GetUnderlyingType(td.Type).(*dsl.SimpleType); ok {
		if _, ok := t.ResolvedDefinition.(*dsl.RecordDefinition); ok {
			resolvesToRecord = true
		}
	}

	// // If the NamedType is Generic and resolves to a RecordDefinition, we can drop the type parameters in the alias declaration
	if dsl.IsGeneric(td) && resolvesToRecord {
		fmt.Fprintf(w, "%s = %s\n", common.TypeIdentifierName(td.Name), common.TypeSyntaxWithoutTypeParameters(td.Type, td.Namespace))
	} else {
		fmt.Fprintf(w, "%s = %s\n", common.TypeIdentifierName(td.Name), common.TypeSyntax(td.Type, td.Namespace))
	}
	common.WriteDocstring(w, td.Comment)
	w.Indent().WriteStringln("")
}

func writeRecord(w *formatting.IndentedWriter, rec *dsl.RecordDefinition, st dsl.SymbolTable) {
	fmt.Fprintf(w, "class %s%s:\n", common.TypeSyntaxWithoutTypeParameters(rec, rec.Namespace), GetGenericBase(rec))
	w.Indented(func() {
		common.WriteDocstring(w, rec.Comment)
		for _, field := range rec.Fields {
			fmt.Fprintf(w, "%s: %s\n", common.FieldIdentifierName(field.Name), common.TypeSyntax(field.Type, rec.Namespace))

			common.WriteDocstring(w, field.Comment)
		}
		w.WriteStringln("")

		if len(rec.Fields) > 0 {
			w.WriteStringln("def __init__(self, *,")
			w.Indented(func() {
				for _, f := range rec.Fields {
					fieldName := common.FieldIdentifierName(f.Name)
					fieldTypeSyntax := common.TypeSyntax(f.Type, rec.Namespace)
					fmt.Fprintf(w, "%s: ", fieldName)

					defaultExpression, defaultExpressionKind := typeDefault(f.Type, rec.Namespace, "", st)
					switch defaultExpressionKind {
					case defaultValueKindNone:
						w.WriteString(fieldTypeSyntax)
					case defaultValueKindImmutable:
						fmt.Fprintf(w, "%s = %s", fieldTypeSyntax, defaultExpression)
					case defaultValueKindMutable:
						fmt.Fprintf(w, "typing.Optional[%s] = None", fieldTypeSyntax)
					}
					w.WriteStringln(",")
				}
			})

			w.WriteStringln("):")
			w.Indented(func() {
				for _, f := range rec.Fields {
					fieldName := common.FieldIdentifierName(f.Name)
					defaultExpression, defaultExpressionKind := typeDefault(f.Type, rec.Namespace, "", st)
					switch defaultExpressionKind {
					case defaultValueKindNone, defaultValueKindImmutable:
						fmt.Fprintf(w, "self.%s = %s\n", fieldName, fieldName)
					case defaultValueKindMutable:
						fmt.Fprintf(w, "self.%s = %s if %s is not None else %s\n", fieldName, fieldName, fieldName, defaultExpression)
					}
				}
			})
			w.WriteStringln("")
		}

		for _, computedField := range rec.ComputedFields {
			expressionTypeSyntax := common.TypeSyntax(computedField.Expression.GetResolvedType(), rec.Namespace)
			fieldName := common.ComputedFieldIdentifierName(computedField.Name)
			fmt.Fprintf(w, "def %s(self) -> %s:\n", fieldName, expressionTypeSyntax)
			w.Indented(func() {
				common.WriteDocstring(w, computedField.Comment)
				writeComputedFieldExpression(w, computedField.Expression, rec.Namespace)
				w.WriteStringln("")
			})
		}

		writeEqMethod(w, rec)

		w.WriteStringln("def __str__(self) -> str:")
		w.Indented(func() {
			fmt.Fprintf(w, "return f\"%s(", common.TypeSyntaxWithoutTypeParameters(rec, rec.Namespace))
			formatting.Delimited(w, ", ", rec.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
				fmt.Fprintf(w, "%s={self.%s}", f.Name, common.FieldIdentifierName(f.Name))
			})
			w.WriteString(")\"\n")
		})
		w.WriteStringln("")

		w.WriteStringln("def __repr__(self) -> str:")
		w.Indented(func() {
			fmt.Fprintf(w, "return f\"%s(", common.TypeSyntaxWithoutTypeParameters(rec, rec.Namespace))
			formatting.Delimited(w, ", ", rec.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
				fmt.Fprintf(w, "%s={repr(self.%s)}", f.Name, common.FieldIdentifierName(f.Name))
			})
			w.WriteString(")\"\n")
		})
		w.WriteStringln("")
	})
	w.WriteStringln("")
}

func writeEqMethod(w *formatting.IndentedWriter, rec *dsl.RecordDefinition) {
	w.WriteStringln("def __eq__(self, other: object) -> bool:")
	w.Indented(func() {
		{
			w.WriteStringln("return (")
			w.Indented(func() {
				fmt.Fprintf(w, "isinstance(other, %s)\n", common.TypeSyntaxWithoutTypeParameters(rec, rec.Namespace))
				for _, field := range rec.Fields {
					w.WriteString("and ")
					fieldIdentifier := common.FieldIdentifierName(field.Name)
					w.WriteStringln(typeEqualityExpression(field.Type, "self."+fieldIdentifier, "other."+fieldIdentifier))
				}
			})
			w.WriteStringln(")")
		}
	})
	w.WriteStringln("")
}

func typeEqualityExpression(t dsl.Type, a, b string) string {
	if hasSimpleEquality(t) {
		return fmt.Sprintf("%s == %s", a, b)
	}

	switch t := t.(type) {
	case *dsl.SimpleType:
		return typeDefinitionEqualityExpression(t.ResolvedDefinition, a, b)
	case *dsl.GeneralizedType:
		switch t.Dimensionality.(type) {
		case nil:
			if t.Cases.IsSingle() {
				return typeEqualityExpression(t.Cases[0].Type, a, b)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("(%s is None if %s is None else (%s is not None and %s))", b, a, b, typeEqualityExpression(t.Cases[1].Type, a, b))
			}
		case *dsl.Vector:
			return fmt.Sprintf("len(%s) == len(%s) and all(%s for %s, %s in zip(%s, %s))", a, b, typeEqualityExpression(t.ToScalar(), "a", "b"), "a", "b", a, b)
		}
		return fmt.Sprintf("yardl.structural_equal(%s, %s)", a, b)
	default:
		panic(fmt.Sprintf("unsupported type: %T", t))
	}
}

func typeDefinitionEqualityExpression(t dsl.TypeDefinition, a, b string) string {
	switch t := t.(type) {
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("yardl.structural_equal(%s, %s)", a, b)
	case *dsl.NamedType:
		return typeEqualityExpression(t.Type, a, b)
	}

	return fmt.Sprintf("%s == %s", a, b)
}

func hasSimpleEquality(t dsl.Node) bool {
	res := true
	dsl.Visit(t, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.SimpleType:
			self.Visit(t.ResolvedDefinition)
		case *dsl.Array, *dsl.GenericTypeParameter:
			res = false
			return
		}

		self.VisitChildren(node)
	})
	return res
}

type tailHandler func(next func())

type tailWrapper struct {
	outer   *tailWrapper
	handler tailHandler
}

func (t tailWrapper) Append(handler tailHandler) tailWrapper {
	return tailWrapper{
		outer:   &t,
		handler: handler,
	}
}

func (t tailWrapper) Run(body func()) {
	t.composeFunc(body)()
}

func (t tailWrapper) composeFunc(next func()) func() {
	if t.handler == nil {
		return next
	}
	this := func() { t.handler(next) }
	if t.outer == nil {
		return this
	}

	return t.outer.composeFunc(this)
}

func writeComputedFieldExpression(w *formatting.IndentedWriter, expression dsl.Expression, contextNamespace string) {
	helperFunctionLookup := make(map[any]string)
	dsl.Visit(expression, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.FunctionCallExpression:
			if t.FunctionName == dsl.FunctionDimensionIndex {
				arrType := (t.Arguments[0].GetResolvedType())
				if _, ok := helperFunctionLookup[arrType]; !ok {
					funcName := fmt.Sprintf("_helper_%d", len(helperFunctionLookup))
					helperFunctionLookup[arrType] = funcName
					fmt.Fprintf(w, "def %s(dim_name: str) -> int:\n", funcName)
					w.Indented(func() {
						dims := dsl.ToGeneralizedType(arrType).Dimensionality.(*dsl.Array).Dimensions
						for i, d := range *dims {
							fmt.Fprintf(w, "if dim_name == \"%s\":\n", *d.Name)
							w.Indented(func() {
								fmt.Fprintf(w, "return %d\n", i)
							})
						}
						fmt.Fprintf(w, "raise KeyError(f\"Unknown dimension name: '{dim_name}'\")\n")
						w.WriteStringln("")
					})
				}
			}
		}
		self.VisitChildren(node)
	})

	varCounter := 0
	newVarName := func() string {
		varName := fmt.Sprintf("_var%d", varCounter)
		varCounter++
		return varName
	}

	tail := tailWrapper{}.Append(func(next func()) {
		w.WriteString("return ")
		next()
		w.WriteStringln("")
	})

	dsl.VisitWithContext(expression, tail, func(self dsl.VisitorWithContext[tailWrapper], node dsl.Node, tail tailWrapper) {
		switch t := node.(type) {
		case *dsl.UnaryExpression:
			tail.Run(func() {
				if t.Operator != dsl.UnaryOpNegate {
					panic(fmt.Sprintf("unexpected unary operator %d", t.Operator))
				}
				w.WriteString("-(")
				self.Visit(t.Expression, tailWrapper{})
				w.WriteString(")")
			})
		case *dsl.BinaryExpression:
			tail.Run(func() {
				requiresParentheses := false
				if l, ok := t.Left.(*dsl.BinaryExpression); ok && l.Operator.Precedence() < t.Operator.Precedence() {
					requiresParentheses = true
				}

				if requiresParentheses {
					w.WriteString("(")
				}
				self.Visit(t.Left, tailWrapper{})
				if requiresParentheses {
					w.WriteString(")")
				}

				w.WriteString(" ")

				switch t.Operator {
				case dsl.BinaryOpAdd:
					w.WriteString("+")
				case dsl.BinaryOpSub:
					w.WriteString("-")
				case dsl.BinaryOpMul:
					w.WriteString("*")
				case dsl.BinaryOpDiv:
					w.WriteString("//")
				case dsl.BinaryOpPow:
					w.WriteString("**")
				default:
					panic(fmt.Sprintf("unexpected binary operator %d", t.Operator))
				}

				w.WriteString(" ")

				requiresParentheses = false
				if r, ok := t.Right.(*dsl.BinaryExpression); ok && r.Operator.Precedence() < t.Operator.Precedence() {
					requiresParentheses = true
				}

				if requiresParentheses {
					w.WriteString("(")
				}
				self.Visit(t.Right, tailWrapper{})
				if requiresParentheses {
					w.WriteString(")")
				}
			})
		case *dsl.IntegerLiteralExpression:
			tail.Run(func() {
				fmt.Fprintf(w, "%d", &t.Value)
			})
		case *dsl.FloatingPointLiteralExpression:
			tail.Run(func() {
				w.WriteString(t.Value)
			})
		case *dsl.StringLiteralExpression:
			tail.Run(func() {
				fmt.Fprintf(w, "%q", t.Value)
			})
		case *dsl.MemberAccessExpression:
			tail.Run(func() {
				if t.Target == nil {
					if t.Kind == dsl.MemberAccessVariable {
						w.WriteString(common.FieldIdentifierName(t.Member))
						return
					}

					w.WriteString("self")
				} else {
					self.Visit(t.Target, tailWrapper{})
				}
				w.WriteString(".")
				if t.Kind == dsl.MemberAccessComputedField {
					fmt.Fprintf(w, "%s()", common.ComputedFieldIdentifierName(t.Member))
				} else {
					w.WriteString(common.FieldIdentifierName(t.Member))
				}
			})
		case *dsl.SubscriptExpression:
			tail.Run(func() {
				isTargetArray := false
				if t.Target != nil {
					if gt, ok := t.Target.GetResolvedType().(*dsl.GeneralizedType); ok {
						if _, ok := gt.Dimensionality.(*dsl.Array); ok {
							isTargetArray = true
						}
					}
				}
				if isTargetArray {
					// a cast is needed for numpy subscripting
					fmt.Fprintf(w, "typing.cast(%s, ", common.TypeSyntax(t.GetResolvedType(), contextNamespace))
				}

				self.Visit(t.Target, tailWrapper{})
				w.WriteString("[")
				formatting.Delimited(w, ", ", t.Arguments, func(w *formatting.IndentedWriter, i int, a *dsl.SubscriptArgument) {
					self.Visit(a.Value, tailWrapper{})
				})
				w.WriteString("]")

				if isTargetArray {
					w.WriteString(")")
				}
			})
		case *dsl.FunctionCallExpression:
			tail.Run(func() {
				switch t.FunctionName {
				case dsl.FunctionSize:
					switch dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(type) {
					case *dsl.Vector, *dsl.Map:
						fmt.Fprintf(w, "len(")
						self.Visit(t.Arguments[0], tailWrapper{})
						fmt.Fprintf(w, ")")
					case *dsl.Array:
						self.Visit(t.Arguments[0], tailWrapper{})

						if len(t.Arguments) == 1 {
							fmt.Fprintf(w, ".size")
						} else {
							fmt.Fprintf(w, ".shape[")
							remainingArgs := t.Arguments[1:]
							formatting.Delimited(w, ", ", remainingArgs, func(w *formatting.IndentedWriter, i int, arg dsl.Expression) {
								self.Visit(arg, tailWrapper{})
							})
							fmt.Fprintf(w, "]")
						}
					}
				case dsl.FunctionDimensionIndex:
					helperFuncName := helperFunctionLookup[t.Arguments[0].GetResolvedType()]
					fmt.Fprintf(w, "%s(", helperFuncName)
					self.Visit(t.Arguments[1], tailWrapper{})
					w.WriteString(")")

				case dsl.FunctionDimensionCount:
					self.Visit(t.Arguments[0], tailWrapper{})
					fmt.Fprintf(w, ".ndim")
				default:
					panic(fmt.Sprintf("Unknown function '%s'", t.FunctionName))
				}
			})
		case *dsl.SwitchExpression:
			targetType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Target.GetResolvedType()))

			unionVariableName := newVarName()
			fmt.Fprintf(w, "%s = ", unionVariableName)
			self.Visit(t.Target, tailWrapper{})
			w.WriteStringln("")

			if targetType.Cases.IsOptional() {
				for i, switchCase := range t.Cases {
					writeSwitchCaseOverOptional(w, switchCase, unionVariableName, i == len(targetType.Cases)-1, self, tail)
				}
				return
			}

			if targetType.Cases.IsUnion() {

				// Special handling for SwitchExpression over a Union from an imported namespace
				targetTypeNamespace := ""
				dsl.Visit(t.Target, func(self dsl.Visitor, node dsl.Node) {
					switch node := node.(type) {
					case *dsl.SimpleType:
						self.Visit(node.ResolvedDefinition)
					case *dsl.RecordDefinition:
						for _, field := range node.Fields {
							u := dsl.GetUnderlyingType(field.Type)
							if u == t.Target.GetResolvedType() {
								if targetTypeNamespace == "" {
									meta := node.GetDefinitionMeta()
									targetTypeNamespace = meta.Namespace
								}
								return
							}
						}
					case dsl.Expression:
						t := node.GetResolvedType()
						if t != nil {
							self.Visit(t)
						}
					}
					self.VisitChildren(node)
				})

				unionClassName, _ := common.UnionClassName(targetType)
				if targetTypeNamespace != "" && targetTypeNamespace != contextNamespace {
					unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(targetTypeNamespace), unionClassName)
				}

				for _, switchCase := range t.Cases {
					writeSwitchCaseOverUnion(w, targetType, unionClassName, switchCase, unionVariableName, self, tail)
				}

				fmt.Fprintf(w, "raise RuntimeError(\"Unexpected union case\")\n")
				return
			}

			// this is over a single type
			if len(t.Cases) != 1 {
				panic("switch expression over a single type expected to have exactly one case")
			}
			switchCase := t.Cases[0]
			switch pattern := switchCase.Pattern.(type) {
			case *dsl.DeclarationPattern:
				fmt.Fprintf(w, "%s = %s\n", common.FieldIdentifierName(pattern.Identifier), unionVariableName)
				self.Visit(switchCase.Expression, tail)
			case *dsl.TypePattern, *dsl.DiscardPattern:
				self.Visit(switchCase.Expression, tail)
			default:
				panic(fmt.Sprintf("Unexpected pattern type %T", t.Cases[0].Pattern))
			}

		case *dsl.TypeConversionExpression:
			tail = tail.Append(func(next func()) {
				fmt.Fprintf(w, "%s(", typeConversionCallable(t.Type))
				next()
				w.WriteString(")")
			})

			self.Visit(t.Expression, tail)
		default:
			panic(fmt.Sprintf("Unknown expression type '%T'", t))
		}
	})
}

func writeSwitchCaseOverOptional(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, isLastCase bool, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeCore := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if declarationIdentifier != "" {
			fmt.Fprintf(w, "%s = %s\n", declarationIdentifier, variableName)
		}
		visitor.Visit(switchCase.Expression, tail)
	}

	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if isLastCase {
			writeCore(typePattern, declarationIdentifier)
			return
		}

		if typePattern.Type == nil {
			fmt.Fprintf(w, "if %s is None:\n", variableName)
		} else {
			fmt.Fprintf(w, "if %s is not None:\n", variableName)
		}

		w.Indented(func() {
			writeCore(typePattern, declarationIdentifier)
		})
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		writeCore(nil, "")
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func writeSwitchCaseOverUnion(w *formatting.IndentedWriter, unionType *dsl.GeneralizedType, unionClassName string, switchCase *dsl.SwitchCase, variableName string, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		for _, typeCase := range unionType.Cases {
			if dsl.TypesEqual(typePattern.Type, typeCase.Type) {
				if typePattern.Type == nil {
					fmt.Fprintf(w, "if %s is None:\n", variableName)
					w.Indented(func() {
						visitor.Visit(switchCase.Expression, tail)
					})
				} else {
					fmt.Fprintf(w, "if isinstance(%s, %s.%s):\n", variableName, unionClassName, formatting.ToPascalCase(typeCase.Tag))
					w.Indented(func() {
						if declarationIdentifier != "" {
							fmt.Fprintf(w, "%s = %s.value\n", declarationIdentifier, variableName)
						}
						visitor.Visit(switchCase.Expression, tail)
					})
				}
				return
			}
		}
		panic(fmt.Sprintf("Did not find pattern type  '%s'", dsl.TypeToShortSyntax(typePattern.Type, false)))
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		visitor.Visit(switchCase.Expression, tail)
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func typeConversionCallable(t dsl.Type) string {
	switch t := t.(type) {
	case *dsl.SimpleType:
		switch t := t.ResolvedDefinition.(type) {
		case dsl.PrimitiveDefinition:
			switch t {
			case dsl.Bool:
				return "bool"
			case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Size:
				return "int"
			case dsl.Float32, dsl.Float64:
				return "float"
			case dsl.ComplexFloat32, dsl.ComplexFloat64:
				return "complex"
			case dsl.String:
				return "str"
			case dsl.Date:
				return "datetime.date"
			case dsl.Time:
				return "datetime.time"
			case dsl.DateTime:
				return "datetime.datetime"
			}
		}
	}
	panic(fmt.Sprintf("Unsupported type '%s'", t))
}

func GetGenericBase(t dsl.TypeDefinition) string {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) == 0 {
		return ""
	}

	var typeParams []string
	for _, tp := range meta.TypeParameters {
		use := tp.Annotations[common.TypeParameterUseAnnotationKey].(common.TypeParameterUse)
		if use&common.TypeParameterUseScalar != 0 {
			typeParams = append(typeParams, common.TypeIdentifierName(tp.Name))
		}
		if use&common.TypeParameterUseArray != 0 {
			typeParams = append(typeParams, common.NumpyTypeParameterSyntax(tp))
		}
	}

	if len(typeParams) == 0 {
		return ""
	}

	return fmt.Sprintf("(typing.Generic[%s])", strings.Join(typeParams, ", "))
}

func writeEnum(w *formatting.IndentedWriter, enum *dsl.EnumDefinition) {
	var base string
	if enum.IsFlags {
		base = "enum.IntFlag"
	} else {
		base = "yardl.OutOfRangeEnum"
	}

	enumTypeSyntax := common.TypeSyntax(enum, enum.Namespace)
	fmt.Fprintf(w, "class %s(%s):\n", enumTypeSyntax, base)

	w.Indented(func() {
		common.WriteDocstring(w, enum.Comment)
		for _, value := range enum.Values {
			fmt.Fprintf(w, "%s = %d\n", common.EnumValueIdentifierName(value.Symbol), &value.IntegerValue)
			common.WriteDocstring(w, value.Comment)
		}

		if enum.IsFlags {
			w.WriteStringln("")
			w.WriteStringln("def __eq__(self, other: object) -> bool:")
			w.Indented(func() {
				fmt.Fprintf(w, "return isinstance(other, %s) and self.value == other.value\n", enumTypeSyntax)
			})
			w.WriteStringln("")

			w.WriteStringln("def __hash__(self) -> int:")
			w.Indented(func() {
				w.WriteStringln("return hash(self.value)")
			})
			w.WriteStringln("")

			w.WriteStringln("__str__ = enum.Flag.__str__ # type: ignore")

		}
	})
	w.WriteStringln("")
}

type defaultValueKind int

const (
	defaultValueKindNone defaultValueKind = iota
	defaultValueKindImmutable
	defaultValueKindMutable
)

func typeDefault(t dsl.Type, contextNamespace string, namedType string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case nil:
		return "None", defaultValueKindImmutable
	case *dsl.SimpleType:
		return typeDefinitionDefault(t.ResolvedDefinition, contextNamespace, st)
	case *dsl.GeneralizedType:
		switch td := t.Dimensionality.(type) {
		case nil:
			defaultExpression, defaultKind := typeDefault(t.Cases[0].Type, contextNamespace, "", st)
			if t.Cases.IsSingle() || t.Cases.HasNullOption() {
				return defaultExpression, defaultKind
			}

			var unionClassName string
			if namedType != "" {
				unionClassName = namedType
			} else {
				unionClassName, _ = common.UnionClassName(t)
			}

			unionCaseConstructor := fmt.Sprintf("%s.%s", unionClassName, formatting.ToPascalCase(t.Cases[0].Tag))

			switch defaultKind {
			case defaultValueKindNone:
				return "", defaultKind
			default:
				return fmt.Sprintf(`%s(%s)`, unionCaseConstructor, defaultExpression), defaultKind
			}

		case *dsl.Vector:
			if td.Length == nil {
				return "[]", defaultValueKindMutable
			}

			scalarDefault, scalarDefaultKind := typeDefault(t.Cases[0].Type, contextNamespace, "", st)

			switch scalarDefaultKind {
			case defaultValueKindNone:
				return "", defaultValueKindNone
			case defaultValueKindImmutable:
				return fmt.Sprintf("[%s] * %d", scalarDefault, *td.Length), defaultValueKindMutable
			case defaultValueKindMutable:
				return fmt.Sprintf("[%s for _ in range(%d)]", scalarDefault, *td.Length), defaultValueKindMutable
			}

		case *dsl.Array:
			context := dTypeExpressionContext{
				namespace: contextNamespace,
				root:      false,
			}

			scalar := t.ToScalar()
			if dsl.TypeContainsGenericTypeParameter(scalar) {
				return "", defaultValueKindNone
			}

			dtype := typeDTypeExpression(scalar, context)

			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("np.zeros((%s,), dtype=%s)", strings.Join(dims, ", "), dtype), defaultValueKindMutable
			}

			if td.HasKnownNumberOfDimensions() {
				shape := fmt.Sprintf("(%s)", strings.Repeat("0, ", len(*td.Dimensions))[0:len(*td.Dimensions)*3-2])
				return fmt.Sprintf("np.zeros(%s, dtype=%s)", shape, dtype), defaultValueKindMutable
			}

			return fmt.Sprintf("np.zeros((), dtype=%s)", dtype), defaultValueKindMutable

		case *dsl.Map:
			return "{}", defaultValueKindMutable
		}
	}

	return "", defaultValueKindNone
}

func typeDefinitionDefault(t dsl.TypeDefinition, contextNamespace string, st dsl.SymbolTable) (string, defaultValueKind) {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Bool:
			return "False", defaultValueKindImmutable
		case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Size:
			return "0", defaultValueKindImmutable
		case dsl.Float32, dsl.Float64:
			return "0.0", defaultValueKindImmutable
		case dsl.ComplexFloat32, dsl.ComplexFloat64:
			return "0j", defaultValueKindImmutable
		case dsl.String:
			return `""`, defaultValueKindImmutable
		case dsl.Date:
			return "datetime.date(1970, 1, 1)", defaultValueKindImmutable
		case dsl.Time:
			return "yardl.Time()", defaultValueKindImmutable
		case dsl.DateTime:
			return "yardl.DateTime()", defaultValueKindImmutable
		}
	case *dsl.EnumDefinition:
		zeroValue := t.GetZeroValue()
		if t.IsFlags {
			if zeroValue == nil {
				return fmt.Sprintf("%s(0)", common.TypeSyntax(t, contextNamespace)), defaultValueKindImmutable
			} else {
				return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindImmutable
			}
		}

		if zeroValue == nil {
			return "", defaultValueKindNone
		}

		return fmt.Sprintf("%s.%s", common.TypeSyntax(t, contextNamespace), common.EnumValueIdentifierName(zeroValue.Symbol)), defaultValueKindImmutable
	case *dsl.NamedType:
		return typeDefault(t.Type, contextNamespace, common.TypeSyntax(t, contextNamespace), st)

	case *dsl.RecordDefinition:
		if len(t.TypeArguments) == 0 && len(t.TypeParameters) > 0 {
			// *Open* Generic Record type
			// Should never get here - typeDefault is only called on Fields, which must be closed if generic
			panic(fmt.Sprintf("No typeDefault for open generic record %s", t.Name))
		}

		// t is a *closed* generic record type
		// genericDef is its original generic type definition
		genericDef := st[t.GetQualifiedName()].(*dsl.RecordDefinition)
		args := make([]string, 0)
		for i, f := range t.Fields {
			fieldDefaultExpr, fieldDefaultKind := typeDefault(f.Type, contextNamespace, "", st)
			if fieldDefaultKind == defaultValueKindNone {
				return "", defaultValueKindNone
			}

			// Only write a constructor argument if it is needed, e.g. the record definition's field is generic and doesn't have a default value
			_, genDefaultKind := typeDefault(genericDef.Fields[i].Type, contextNamespace, "", st)
			if genDefaultKind == defaultValueKindNone {
				args = append(args, fmt.Sprintf("%s=%s", common.FieldIdentifierName(f.Name), fieldDefaultExpr))
			}
		}

		return fmt.Sprintf("%s(%s)", common.TypeSyntaxWithoutTypeParameters(t, contextNamespace), strings.Join(args, ", ")), defaultValueKindMutable
	}

	return "", defaultValueKindNone
}

func writeGetDTypeFunc(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	w.WriteStringln("def _mk_get_dtype():")
	w.Indented(func() {
		w.WriteStringln("dtype_map: dict[typing.Union[type, types.GenericAlias], typing.Union[np.dtype[typing.Any], typing.Callable[[tuple[type, ...]], np.dtype[typing.Any]]]] = {}")
		w.WriteStringln("get_dtype = _dtypes.make_get_dtype_func(dtype_map)\n")

		context := dTypeExpressionContext{
			namespace: ns.Name,
			root:      true,
		}

		writeUnionCaseDtypes := func(gt *dsl.GeneralizedType, unionClassName string) {
			for _, tc := range gt.Cases {
				if tc.Type != nil && !dsl.TypeContainsGenericTypeParameter(tc.Type) {
					tag := formatting.ToPascalCase(tc.Tag)
					fmt.Fprintf(w, "dtype_map.setdefault(%s.%s, %s)\n", unionClassName, tag, typeDTypeExpression(tc.Type, context))
				}
			}
		}

		writeUnionDtypeIfNeeded := func(td dsl.TypeDefinition, unions map[string]bool, callingNamespace string) {
			dsl.VisitWithContext(td, "", func(self dsl.VisitorWithContext[string], node dsl.Node, currentNamespace string) {
				switch node := node.(type) {
				case dsl.TypeDefinition:
					currentNamespace = node.GetDefinitionMeta().Namespace

				case *dsl.GeneralizedType:
					if node.Cases.IsUnion() {
						unionClassName, _ := common.UnionClassName(node)
						nt, isNamedType := td.(*dsl.NamedType)
						if isNamedType {
							// This is a named type defining a union, so we will use the named type's name instead
							unionClassName = td.GetDefinitionMeta().Name
						}
						if currentNamespace != callingNamespace {
							unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(currentNamespace), unionClassName)
						}

						if !unions[unionClassName] {
							unions[unionClassName] = true
							if isNamedType {
								fmt.Fprintf(w, "dtype_map.setdefault(%s, %s)\n", unionClassName, typeDefinitionDTypeExpression(nt, context))
							} else {
								fmt.Fprintf(w, "dtype_map.setdefault(%s, %s)\n", unionClassName, typeDTypeExpression(node, context))
							}
							writeUnionCaseDtypes(node, unionClassName)
						}
					}
				}

				self.VisitChildren(node, currentNamespace)
			})
		}

		writeDefaults := func(ns *dsl.Namespace, contextNamespace string) {
			unions := make(map[string]bool)

			for _, td := range ns.TypeDefinitions {
				isUnion := false

				if td, ok := td.(*dsl.NamedType); ok {
					if gt, ok := dsl.GetUnderlyingType(td.Type).(*dsl.GeneralizedType); ok {
						if gt.Cases.IsUnion() && !gt.Cases.HasNullOption() {
							isUnion = true
						}
						switch gt.Dimensionality.(type) {
						// Skip Named arrays, vectors, and maps because they resolve to Python
						// 	built-in types having dtype=np.object_
						case *dsl.Array, *dsl.Vector, *dsl.Map:
							continue
						}
					}
				}

				fmt.Fprintf(w, "dtype_map.setdefault(%s, %s)\n", common.TypeSyntaxWithoutTypeParameters(td, contextNamespace), typeDefinitionDTypeExpression(td, context))

				if !isUnion {
					writeUnionDtypeIfNeeded(td, unions, contextNamespace)
				}
			}

			for _, p := range ns.Protocols {
				writeUnionDtypeIfNeeded(p, unions, contextNamespace)
			}
		}

		for _, refNs := range ns.GetAllChildReferences() {
			writeDefaults(refNs, ns.Name)
		}

		writeDefaults(ns, ns.Name)

		w.WriteStringln("\nreturn get_dtype")
	})
	w.WriteStringln("")

	w.WriteStringln("get_dtype = _mk_get_dtype()\n")
}

type dTypeExpressionContext struct {
	namespace            string
	root                 bool
	typeParameterIndexes map[*dsl.GenericTypeParameter]int
}

func typeDefinitionDTypeExpression(t dsl.TypeDefinition, context dTypeExpressionContext) string {
	if !context.root {
		var dtypeExpression string
		switch t := t.(type) {
		case dsl.PrimitiveDefinition:
			switch t {
			case dsl.Bool:
				return "np.dtype(np.bool_)"
			case dsl.Int8, dsl.Uint8, dsl.Int16, dsl.Uint16, dsl.Int32, dsl.Uint32, dsl.Int64, dsl.Uint64, dsl.Float32, dsl.Float64:
				return fmt.Sprintf("np.dtype(np.%s)", strings.ToLower(string(t)))
			case dsl.Size:
				return "np.dtype(np.uint64)"
			case dsl.ComplexFloat32:
				return "np.dtype(np.complex64)"
			case dsl.ComplexFloat64:
				return "np.dtype(np.complex128)"
			case dsl.Date:
				return "np.dtype(np.datetime64)"
			case dsl.Time:
				return "np.dtype(np.timedelta64)"
			case dsl.DateTime:
				return "np.dtype(np.datetime64)"
			case dsl.String:
				return "np.dtype(np.object_)"
			default:
				panic(fmt.Sprintf("Not implemented %s", t))
			}
		case *dsl.GenericTypeParameter:
			index, ok := context.typeParameterIndexes[t]
			if !ok {
				panic("type parameter not found")
			}
			return fmt.Sprintf("get_dtype(type_args[%d])", index)
		}

		if len(t.GetDefinitionMeta().TypeParameters) > 0 {
			typeArgs := make([]string, 0)
			for _, ta := range t.GetDefinitionMeta().TypeArguments {
				typeArgs = append(typeArgs, getTypeSyntaxWithGenricArgsReadFromTupleArgs(ta, context))
			}

			dtypeExpression = fmt.Sprintf("get_dtype(types.GenericAlias(%s, (%s,)))", common.TypeSyntaxWithoutTypeParameters(t, context.namespace), strings.Join(typeArgs, ", "))
		} else {
			dtypeExpression = fmt.Sprintf("get_dtype(%s)", common.TypeSyntaxWithoutTypeParameters(t, context.namespace))
		}

		return dtypeExpression
	}

	meta := t.GetDefinitionMeta()
	lambdaDeclaration := ""
	if len(meta.TypeParameters) > 0 {
		context.typeParameterIndexes = make(map[*dsl.GenericTypeParameter]int)
		for i, p := range meta.TypeParameters {
			context.typeParameterIndexes[p] = i
		}

		lambdaDeclaration = "lambda type_args: "
	}

	switch t := t.(type) {
	case *dsl.NamedType:
		return lambdaDeclaration + typeDTypeExpression(t.Type, context)
	case *dsl.EnumDefinition:
		base := t.BaseType
		if base == nil {
			base = dsl.Int32Type
		}

		return typeDTypeExpression(base, context)

	case *dsl.RecordDefinition:
		fields := make([]string, len(t.Fields))
		for i, f := range t.Fields {
			subarrayShape := ""
			underyingType := dsl.GetUnderlyingType(f.Type)
			if gt, ok := underyingType.(*dsl.GeneralizedType); ok {
				if vec, ok := gt.Dimensionality.(*dsl.Vector); ok && vec.Length != nil {
					subarrayShape = fmt.Sprintf("(%d,)", *vec.Length)
				} else if arr, ok := gt.Dimensionality.(*dsl.Array); ok && arr.IsFixed() {
					dims := make([]string, len(*arr.Dimensions))
					for i, d := range *arr.Dimensions {
						dims[i] = strconv.FormatUint(*d.Length, 10)
					}
					subarrayShape = fmt.Sprintf("(%s,)", strings.Join(dims, ", "))
				}
			}

			if subarrayShape != "" {
				subarrayShape = fmt.Sprintf(", %s", subarrayShape)
			}

			fields[i] = fmt.Sprintf("('%s', %s%s)", common.FieldIdentifierName(f.Name), typeDTypeExpression(f.Type, context), subarrayShape)
		}

		return fmt.Sprintf("%snp.dtype([%s], align=True)", lambdaDeclaration, strings.Join(fields, ", "))
	}

	return "np.dtype(np.object_)"
}

func typeDTypeExpression(t dsl.Type, context dTypeExpressionContext) string {
	switch t := dsl.GetUnderlyingType(t).(type) {
	case *dsl.SimpleType:
		context.root = false
		return typeDefinitionDTypeExpression(t.ResolvedDefinition, context)

	case *dsl.GeneralizedType:
		switch td := t.Dimensionality.(type) {
		case nil:
			if t.Cases.IsOptional() {
				return fmt.Sprintf("np.dtype([('has_value', np.dtype(np.bool_)), ('value', %s)], align=True)", typeDTypeExpression(t.Cases[1].Type, context))
			}
		case *dsl.Vector:
			if td.Length != nil {
				return typeDTypeExpression(t.ToScalar(), context)
			}

		case *dsl.Array:
			if td.IsFixed() {
				return typeDTypeExpression(t.ToScalar(), context)
			}
		}
	}

	return "np.dtype(np.object_)"
}

func getTypeSyntaxWithGenricArgsReadFromTupleArgs(t dsl.Type, context dTypeExpressionContext) string {
	var f dsl.TypeSyntaxWriter[string] = func(self dsl.TypeSyntaxWriter[string], typeOrTypeDef dsl.Node, _ string) string {
		switch t := typeOrTypeDef.(type) {
		case *dsl.GenericTypeParameter:
			return fmt.Sprintf("type_args[%d]", context.typeParameterIndexes[t])
		}

		return common.TypeSyntaxWriter(self, typeOrTypeDef, context.namespace)
	}

	return f.ToSyntax(t, context.namespace)
}
