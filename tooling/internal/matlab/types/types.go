// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package types

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteTypes(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	for _, td := range ns.TypeDefinitions {
		b := bytes.Buffer{}
		w := formatting.NewIndentedWriter(&b, "  ")
		common.WriteGeneratedFileHeader(w)

		switch td := td.(type) {
		case *dsl.NamedType:
			writeNamedType(w, td)
		case *dsl.EnumDefinition:
			writeEnum(w, td)
		case *dsl.RecordDefinition:
			writeRecord(w, td, st)
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}

		fname := fmt.Sprintf("%s.m", common.TypeSyntax(td, ns.Name))

		definitionsPath := path.Join(packageDir, fname)
		if err := iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func writeNamedType(w *formatting.IndentedWriter, td *dsl.NamedType) {
	common.WriteComment(w, td.Comment)
	fmt.Fprintf(w, "classdef %s < %s\n", common.TypeSyntax(td, td.Namespace), common.TypeSyntax(td.Type, td.Namespace))
	w.WriteStringln("end")
}

func writeEnum(w *formatting.IndentedWriter, enum *dsl.EnumDefinition) {
	var base string
	if enum.BaseType == nil {
		base = "uint64"
	} else {
		base = common.TypeSyntax(enum.BaseType, enum.Namespace)
	}

	common.WriteComment(w, enum.Comment)
	enumTypeSyntax := common.TypeSyntax(enum, enum.Namespace)
	fmt.Fprintf(w, "classdef %s < %s\n", enumTypeSyntax, base)
	common.WriteBlockBody(w, func() {
		fmt.Fprintf(w, "enumeration\n")
		common.WriteBlockBody(w, func() {
			for _, value := range enum.Values {
				common.WriteComment(w, value.Comment)
				fmt.Fprintf(w, "%s (%d)\n", common.EnumValueIdentifierName(value.Symbol), &value.IntegerValue)
			}
		})
	})
}

func writeRecord(w *formatting.IndentedWriter, rec *dsl.RecordDefinition, st dsl.SymbolTable) {
	common.WriteComment(w, rec.Comment)

	fmt.Fprintf(w, "classdef %s < handle\n", common.TypeSyntax(rec, rec.Namespace))
	common.WriteBlockBody(w, func() {

		w.WriteStringln("properties")
		var fieldNames []string
		common.WriteBlockBody(w, func() {
			for i, field := range rec.Fields {
				common.WriteComment(w, field.Comment)
				fieldNames = append(fieldNames, common.FieldIdentifierName(field.Name))
				fmt.Fprintf(w, "%s\n", common.FieldIdentifierName(field.Name))
				if i < len(rec.Fields)-1 {
					w.WriteStringln("")
				}
			}
		})
		w.WriteStringln("")

		w.WriteStringln("methods")
		common.WriteBlockBody(w, func() {

			// Record Constructor
			fmt.Fprintf(w, "function obj = %s(%s)\n", rec.Name, strings.Join(fieldNames, ", "))
			common.WriteBlockBody(w, func() {
				for _, field := range rec.Fields {
					fmt.Fprintf(w, "obj.%s = %s;\n", common.FieldIdentifierName(field.Name), common.FieldIdentifierName(field.Name))
				}
			})
			w.WriteStringln("")

			// Computed Fields
			for _, computedField := range rec.ComputedFields {
				fieldName := common.ComputedFieldIdentifierName(computedField.Name)

				common.WriteComment(w, computedField.Comment)
				fmt.Fprintf(w, "function res = %s(self)\n", fieldName)
				common.WriteBlockBody(w, func() {
					writeComputedFieldExpression(w, computedField.Expression, rec.Namespace)
				})
				w.WriteStringln("")
			}
			w.WriteStringln("")

			// eq method
			w.WriteStringln("function res = eq(obj, other)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("res = ...")
				w.Indented(func() {
					fmt.Fprintf(w, "isa(other, '%s')", common.TypeSyntax(rec, rec.Namespace))
					for _, field := range rec.Fields {
						w.WriteStringln(" && ...")
						fieldIdentifier := common.FieldIdentifierName(field.Name)
						w.WriteString(typeEqualityExpression(field.Type, "obj."+fieldIdentifier, "other."+fieldIdentifier))
					}
					w.WriteStringln(";")
				})
			})

			// neq method
		})

	})
}

func typeEqualityExpression(t dsl.Type, a, b string) string {
	// TODO: Figure out equality because in Matlab both 'a' and 'b' can be scalar or non-scalar...
	if hasSimpleEquality(t) {
		// return fmt.Sprintf("%s == %s", a, b)
		return fmt.Sprintf("all(%s == %s)", a, b)
		// return fmt.Sprintf("all([%s] == [%s])", a, b)
	}

	// TODO: Other forms
	// panic(fmt.Sprintf("How about type equality expression for %s", dsl.TypeToShortSyntax(t, false)))
	return fmt.Sprintf("all(%s == %s)", a, b)
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
					funcName := fmt.Sprintf("helper_%d_", len(helperFunctionLookup))
					helperFunctionLookup[arrType] = funcName
					fmt.Fprintf(w, "function dim = %s(dim_name)\n", funcName)
					common.WriteBlockBody(w, func() {
						dims := dsl.ToGeneralizedType(arrType).Dimensionality.(*dsl.Array).Dimensions
						for i, d := range *dims {
							fmt.Fprintf(w, "if dim_name == \"%s\"\n", *d.Name)
							w.Indented(func() {
								fmt.Fprintf(w, "dim = %d;\n", i)
							})
							w.WriteString("else")
						}
						w.WriteStringln("")
						common.WriteBlockBody(w, func() {
							w.WriteStringln(`throw(yardl.KeyError("Unknown dimension name: '%s'", dim_name));`)
						})
						w.WriteStringln("")
					})
				}
			}
		}
		self.VisitChildren(node)
	})

	varCounter := 0
	newVarName := func() string {
		varCounter++
		return fmt.Sprintf("var%d", varCounter)
	}

	tail := tailWrapper{}.Append(func(next func()) {
		w.WriteString("res = ")
		next()
		w.WriteStringln(";")
		w.WriteStringln("return")
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
					w.WriteString(".*")
				case dsl.BinaryOpDiv:
					w.WriteString("./")
				case dsl.BinaryOpPow:
					w.WriteString("^")
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
			// Collapse adjacent subscript expressions (you can't do `array[x][y]` in Matlab, only `array(x,y)`)
			target := t.Target
			arguments := t.Arguments
			dsl.Visit(t.Target, func(self dsl.Visitor, node dsl.Node) {
				if t, ok := node.(*dsl.SubscriptExpression); ok {
					target = t.Target
					arguments = append(t.Arguments, arguments...)
					self.VisitChildren(t.Target)
				}
			})

			tail.Run(func() {
				self.Visit(target, tailWrapper{})
				w.WriteString("(")
				formatting.Delimited(w, ", ", arguments, func(w *formatting.IndentedWriter, i int, a *dsl.SubscriptArgument) {
					self.Visit(a.Value, tailWrapper{})
				})
				w.WriteString(")")
			})

		case *dsl.FunctionCallExpression:
			tail.Run(func() {
				switch t.FunctionName {
				case dsl.FunctionSize:
					switch dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(type) {
					case *dsl.Vector, *dsl.Map:
						fmt.Fprintf(w, "length(")
						self.Visit(t.Arguments[0], tailWrapper{})
						fmt.Fprintf(w, ")")
					case *dsl.Array:
						w.WriteString("size(")
						self.Visit(t.Arguments[0], tailWrapper{})

						if len(t.Arguments) > 1 {
							w.WriteString("(")
							remainingArgs := t.Arguments[1:]
							formatting.Delimited(w, ", ", remainingArgs, func(w *formatting.IndentedWriter, i int, arg dsl.Expression) {
								self.Visit(arg, tailWrapper{})
							})
							fmt.Fprintf(w, ")")
						}
						w.WriteString(")")
					}

				case dsl.FunctionDimensionIndex:
					helperFuncName := helperFunctionLookup[t.Arguments[0].GetResolvedType()]
					fmt.Fprintf(w, "%s(", helperFuncName)
					self.Visit(t.Arguments[1], tailWrapper{})
					w.WriteString(")")

				case dsl.FunctionDimensionCount:
					w.WriteString("ndims(")
					self.Visit(t.Arguments[0], tailWrapper{})
					w.WriteString(")")

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
				// targetTypeNamespace := ""
				// dsl.Visit(t.Target, func(self dsl.Visitor, node dsl.Node) {
				// 	switch node := node.(type) {
				// 	case *dsl.SimpleType:
				// 		self.Visit(node.ResolvedDefinition)
				// 	case *dsl.RecordDefinition:
				// 		for _, field := range node.Fields {
				// 			u := dsl.GetUnderlyingType(field.Type)
				// 			if u == t.Target.GetResolvedType() {
				// 				if targetTypeNamespace == "" {
				// 					meta := node.GetDefinitionMeta()
				// 					targetTypeNamespace = meta.Namespace
				// 				}
				// 				return
				// 			}
				// 		}
				// 	case dsl.Expression:
				// 		t := node.GetResolvedType()
				// 		if t != nil {
				// 			self.Visit(t)
				// 		}
				// 	}
				// 	self.VisitChildren(node)
				// })

				// unionClassName, _ := common.UnionClassName(targetType)
				// if targetTypeNamespace != "" && targetTypeNamespace != contextNamespace {
				// 	unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(targetTypeNamespace), unionClassName)
				// }

				// for _, switchCase := range t.Cases {
				// 	writeSwitchCaseOverUnion(w, targetType, unionClassName, switchCase, unionVariableName, self, tail)
				// }

				// fmt.Fprintf(w, "raise RuntimeError(\"Unexpected union case\")\n")
				return
			}

			// this is over a single type
			if len(t.Cases) != 1 {
				panic("switch expression over a single type expected to have exactly one case")
			}
			switchCase := t.Cases[0]
			switch pattern := switchCase.Pattern.(type) {
			case *dsl.DeclarationPattern:
				fmt.Fprintf(w, "%s = %s;\n", common.FieldIdentifierName(pattern.Identifier), unionVariableName)
				self.Visit(switchCase.Expression, tail)
			case *dsl.TypePattern, *dsl.DiscardPattern:
				self.Visit(switchCase.Expression, tail)
			default:
				panic(fmt.Sprintf("Unexpected pattern type %T", t.Cases[0].Pattern))
			}

		case *dsl.TypeConversionExpression:
			tail = tail.Append(func(next func()) {
				writeTypeConversion(w, t.Type, next)
			})

			self.Visit(t.Expression, tail)

		default:
			panic(fmt.Sprintf("Unknown expression type '%T'", t))
		}
	})
}

func writeSwitchCaseOverOptional(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, isLastCase bool, visitor dsl.VisitorWithContext[tailWrapper], tail tailWrapper) {
	writeCore := func(declarationIdentifier string) {
		if declarationIdentifier != "" {
			fmt.Fprintf(w, "%s = %s;\n", declarationIdentifier, variableName)
		}
		visitor.Visit(switchCase.Expression, tail)
	}

	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if isLastCase {
			writeCore(declarationIdentifier)
			return
		}

		if typePattern.Type == nil {
			fmt.Fprintf(w, "if isa(%s, yardl.None)\n", variableName)
		} else {
			fmt.Fprintf(w, "if ~isa(%s, yardl.None)\n", variableName)
		}

		common.WriteBlockBody(w, func() {
			writeCore(declarationIdentifier)
		})
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		writeCore("")
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func writeTypeConversion(w *formatting.IndentedWriter, t dsl.Type, next func()) {
	getWrapper := func(t dsl.Type) (string, string) {
		switch t := t.(type) {
		case *dsl.SimpleType:
			switch t := t.ResolvedDefinition.(type) {
			case dsl.PrimitiveDefinition:
				switch t {
				case dsl.Bool:
					return "logical(", ")"
				case dsl.Int8:
					return "int8(", ")"
				case dsl.Uint8:
					return "uint8(", ")"
				case dsl.Int16:
					return "int32(", ")"
				case dsl.Uint16:
					return "uint16(", ")"
				case dsl.Int32:
					return "int32(", ")"
				case dsl.Uint32:
					return "uint32(", ")"
				case dsl.Int64:
					return "int64(", ")"
				case dsl.Uint64, dsl.Size:
					return "uint64(", ")"
				case dsl.Float32:
					return "single(", ")"
				case dsl.Float64:
					return "double(", ")"
				case dsl.ComplexFloat32:
					return "complex(single(", "))"
				case dsl.ComplexFloat64:
					return "complex(double(", "))"
				case dsl.String:
					return "string(", ")"
				case dsl.Date, dsl.Time, dsl.DateTime:
					return "datetime(", ")"
				}
			}
		}
		panic(fmt.Sprintf("Unsupported type '%s'", t))
	}

	open, close := getWrapper(t)
	w.WriteString(open)
	next()
	w.WriteString(close)
}
