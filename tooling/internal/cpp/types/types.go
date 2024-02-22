// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package types

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteTypes(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	err := writeHeaderFile(env, options)
	if err != nil {
		return err
	}

	return writeSourceFile(env, options)
}

func writeHeaderFile(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#pragma once
#include <array>
#include <complex>
#include <optional>
#include <unordered_map>
#include <variant>
#include <vector>

#include "yardl/yardl.h"
`)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
		writeNamespaceMembers(w, ns)
		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	definitionsPath := path.Join(options.SourcesOutputDir, "types.h")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeSourceFile(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include "types.h"`)
	w.WriteStringln(`#include "yardl/yardl.h"`)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))

		for _, td := range ns.TypeDefinitions {
			switch td := td.(type) {
			case *dsl.EnumDefinition:
				if td.IsFlags {
					typeName := common.TypeIdentifierName(td.Name)
					for _, v := range td.Values {
						fmt.Fprintf(w, "const %s %s::%s = %s(%s);\n", typeName, typeName, common.EnumValueIdentifierName(v.Symbol), typeName, common.EnumIntegerLiteral(td, v))
					}
				}
			}
		}
		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	definitionsPath := path.Join(options.SourcesOutputDir, "types.cc")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeNamespaceMembers(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.EnumDefinition:
			if td.IsFlags {
				common.WriteComment(w, td.Comment)

				typeName := common.TypeIdentifierName(td.Name)
				var valueTypeSyntax string
				if td.BaseType != nil {
					valueTypeSyntax = common.TypeSyntax(td.BaseType)
				} else {
					valueTypeSyntax = common.TypeSyntax(dsl.Int32Type)
				}

				fmt.Fprintf(w, "struct %s : yardl::BaseFlags<%s, %s> {\n", typeName, valueTypeSyntax, typeName)
				w.Indented(func() {
					w.WriteStringln("using BaseFlags::BaseFlags;")

					for _, v := range td.Values {
						common.WriteComment(w, v.Comment)
						fmt.Fprintf(w, "static const %s %s;\n", typeName, common.EnumValueIdentifierName(v.Symbol))
					}
				})
				fmt.Fprint(w, "};\n\n")

			} else {
				common.WriteComment(w, td.Comment)
				fmt.Fprintf(w, "enum class %s ", common.TypeIdentifierName(td.Name))
				if td.BaseType != nil {
					fmt.Fprintf(w, ": %s ", common.TypeSyntax(td.BaseType))
				}
				fmt.Fprintln(w, "{")
				w.Indented(func() {
					for _, enumValue := range td.Values {
						common.WriteComment(w, enumValue.Comment)
						fmt.Fprintf(w, "%s = %s,\n", common.EnumValueIdentifierName(enumValue.Symbol), common.EnumIntegerLiteral(td, enumValue))
					}
				})
				fmt.Fprint(w, "};\n\n")
			}

		case *dsl.NamedType:
			writeNamedTypeDefinition(w, td)
		case *dsl.RecordDefinition:
			common.WriteComment(w, td.Comment)
			common.WriteDefinitionTemplateSpec(w, td)
			fmt.Fprintf(w, "struct %s {\n", common.TypeIdentifierName(td.Name))
			w.Indented(func() {
				for _, field := range td.Fields {
					common.WriteComment(w, field.Comment)
					fmt.Fprintf(w, "%s %s{};\n", common.TypeSyntax(field.Type), common.FieldIdentifierName(field.Name))
				}

				w.WriteString("\n")

				for _, computedField := range td.ComputedFields {
					isRef := computedField.Expression.IsReference()
					refString := ""
					if isRef {
						refString = " const&"
					}
					common.WriteComment(w, computedField.Comment)
					expressionTypeSyntax := common.TypeSyntax(computedField.Expression.GetResolvedType())
					fieldName := common.ComputedFieldIdentifierName(computedField.Name)
					fmt.Fprintf(w, "%s%s %s() const {\n", expressionTypeSyntax, refString, fieldName)
					w.Indented(func() {
						w.WriteString("return ")
						writeComputedFieldExpression(w, computedField.Expression)
						w.WriteStringln(";")
					})
					fmt.Fprint(w, "}\n\n")

					if isRef {
						common.WriteComment(w, computedField.Comment)
						fmt.Fprintf(w, "%s& %s() {\n", expressionTypeSyntax, fieldName)
						w.Indented(func() {
							fmt.Fprintf(w, "return const_cast<%s&>(std::as_const(*this).%s());\n", expressionTypeSyntax, fieldName)
						})
						fmt.Fprint(w, "}\n\n")
					}
				}

				unused := ""
				if len(td.Fields) == 0 {
					unused = "[[maybe_unused]]"
				}

				fmt.Fprintf(w, "bool operator==(%sconst %s& other) const {\n", unused, common.TypeIdentifierName(td.Name))
				w.Indented(func() {
					w.WriteString("return ")
					if len(td.Fields) == 0 {
						w.WriteString("true")
					} else {
						formatting.Delimited(
							w.Indent(),
							" &&\n",
							td.Fields,
							func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
								fmt.Fprintf(w, "%s == other.%s", common.FieldIdentifierName(f.Name), common.FieldIdentifierName(f.Name))
							})
					}
					w.WriteStringln(";")
				})

				w.WriteString("}\n\n")

				fmt.Fprintf(w, "bool operator!=(%sconst %s& other) const {\n", unused, common.TypeIdentifierName(td.Name))
				w.Indented(func() {
					w.WriteString("return !(*this == other);\n")
				})
				w.WriteStringln("}")

			})
			fmt.Fprint(w, "};\n\n")
		}
	}

	// Write "compatibility" aliases for TypeDefinitions with changed names
	for _, versionLabel := range ns.Versions {
		if len(ns.DefinitionChanges[versionLabel]) > 0 {
			fmt.Fprintf(w, "// Compatibility aliases for version %s.\n\n", versionLabel)
			for _, ch := range ns.DefinitionChanges[versionLabel] {
				writeCompatibilityAliasDefinition(w, ch)
			}
		}
	}
}

func writeNamedTypeDefinition(w *formatting.IndentedWriter, nt *dsl.NamedType) {
	common.WriteComment(w, nt.Comment)
	common.WriteDefinitionTemplateSpec(w, nt)
	fmt.Fprintf(w, "using %s = %s;\n\n", common.TypeIdentifierName(nt.Name), common.TypeSyntax(nt.Type))
}

func writeCompatibilityAliasDefinition(w *formatting.IndentedWriter, ch dsl.DefinitionChange) {
	oldDef := ch.PreviousDefinition()
	newDef := ch.LatestDefinition()
	common.WriteDefinitionTemplateSpec(w, oldDef)
	switch ch := ch.(type) {
	case *dsl.AliasRemoved:
		if nt, ok := ch.LatestDefinition().(*dsl.NamedType); ok {
			fmt.Fprintf(w, "using %s = %s;\n\n", common.TypeIdentifierName(oldDef.GetDefinitionMeta().Name), common.TypeSyntax(nt.Type))
		}
	default:
		fmt.Fprintf(w, "using %s = %s;\n\n", common.TypeIdentifierName(oldDef.GetDefinitionMeta().Name), common.TypeDefinitionSyntax(newDef))
	}
}

func writeComputedFieldExpression(w *formatting.IndentedWriter, expression dsl.Expression) {
	dsl.Visit(expression, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.UnaryExpression:
			if t.Operator != dsl.UnaryOpNegate {
				panic(fmt.Sprintf("unexpected unary operator %d", t.Operator))
			}
			w.WriteString("-(")
			self.Visit(t.Expression)
			w.WriteString(")")
		case *dsl.BinaryExpression:
			if t.Operator == dsl.BinaryOpPow {
				w.WriteString("std::pow(")
				self.Visit(t.Left)
				w.WriteString(", ")
				self.Visit(t.Right)
				w.WriteString(")")
				return
			}

			requiresParentheses := false
			if l, ok := t.Left.(*dsl.BinaryExpression); ok && l.Operator.Precedence() < t.Operator.Precedence() {
				requiresParentheses = true
			}

			if requiresParentheses {
				w.WriteString("(")
			}
			self.Visit(t.Left)
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
				w.WriteString("/")
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
			self.Visit(t.Right)
			if requiresParentheses {
				w.WriteString(")")
			}
		case *dsl.IntegerLiteralExpression:
			w.Write([]byte(common.IntegerLiteral(t.Value, t.ResolvedType)))
		case *dsl.FloatingPointLiteralExpression:
			w.WriteString(t.Value)
			if t.ResolvedType == dsl.Float32Type {
				w.WriteString("f")
			}
		case *dsl.StringLiteralExpression:
			fmt.Fprintf(w, "%q", t.Value)
		case *dsl.MemberAccessExpression:
			if t.Target != nil {
				self.Visit(t.Target)
				w.WriteString(".")
			}
			if t.Kind == dsl.MemberAccessComputedField {
				fmt.Fprintf(w, "%s()", common.ComputedFieldIdentifierName(t.Member))
			} else {
				w.WriteString(common.FieldIdentifierName(t.Member))
			}
		case *dsl.SubscriptExpression:
			self.Visit(t.Target)
			w.WriteString(".at(")
			formatting.Delimited(w, ", ", t.Arguments, func(w *formatting.IndentedWriter, i int, a *dsl.SubscriptArgument) {
				self.Visit(a.Value)
			})
			w.WriteString(")")
		case *dsl.Vector:
		case *dsl.FunctionCallExpression:
			switch t.FunctionName {
			case dsl.FunctionSize:
				self.Visit(t.Arguments[0])
				switch dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(type) {
				case *dsl.Vector, *dsl.Map:
					fmt.Fprintf(w, ".size(")
				case *dsl.Array:
					if len(t.Arguments) == 1 {
						fmt.Fprintf(w, ".size(")
					} else {
						fmt.Fprintf(w, ".shape(")
					}
				}
				if len(t.Arguments) > 1 {
					remainingArgs := t.Arguments[1:]

					formatting.Delimited(w, ", ", remainingArgs, func(w *formatting.IndentedWriter, i int, arg dsl.Expression) {
						self.Visit(arg)
					})
				}
				fmt.Fprint(w, ")")

			case dsl.FunctionDimensionIndex:
				dims := dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Arguments[0].GetResolvedType())).Dimensionality.(*dsl.Array).Dimensions
				fmt.Fprintf(w, "([](std::string dim_name) {\n")
				w.Indented(func() {
					for i, d := range *dims {
						fmt.Fprintf(w, "if (dim_name == \"%s\") return %d;\n", *d.Name, i)
					}
					fmt.Fprintf(w, "throw std::invalid_argument(\"Unknown dimension name: \" + dim_name);\n")
				})
				w.WriteString("})(")
				self.Visit(t.Arguments[1])
				w.WriteString(")")

			case dsl.FunctionDimensionCount:
				self.Visit(t.Arguments[0])
				fmt.Fprintf(w, ".dimension()")
			default:
				panic(fmt.Sprintf("Unknown function '%s'", t.FunctionName))
			}
		case *dsl.SwitchExpression:
			targetType := dsl.ToGeneralizedType(dsl.GetUnderlyingType(t.Target.GetResolvedType()))
			if targetType.Cases.IsUnion() {
				w.WriteStringln("std::visit(")
				w.Indented(func() {
					fmt.Fprintf(w, "[&](auto&& __case_arg__) -> %s {\n", common.TypeSyntax(t.ResolvedType))
					w.Indented(func() {
						for _, switchCase := range t.Cases {
							writeSwitchCaseOverUnion(w, switchCase, "__case_arg__", self)
						}
					})
					w.WriteStringln("},")
					self.Visit(t.Target)
					w.WriteString(")")
				})
				return
			}
			if targetType.Cases.IsOptional() {
				fmt.Fprintf(w, "[](auto&& __case_arg__) -> %s {\n", common.TypeSyntax(t.ResolvedType))
				w.Indented(func() {
					for i, switchCase := range t.Cases {
						writeSwitchCaseOverOptional(w, switchCase, "__case_arg__", i == len(targetType.Cases)-1, self)
					}
				})
				w.WriteString("}(")
				self.Visit(t.Target)
				w.WriteString(")")
				return
			}

			// this is over a single type
			if len(t.Cases) != 1 {
				panic("switch expression over a single type expected to have exactly one case")
			}
			switch pattern := t.Cases[0].Pattern.(type) {
			case *dsl.DeclarationPattern:
				fmt.Fprintf(w, "[]([[maybe_unused]] %s const& %s) -> %s {\n", common.TypeSyntax(pattern.Type), common.FieldIdentifierName(pattern.Identifier), common.TypeSyntax(t.ResolvedType))
				w.Indented(func() {
					w.WriteString("return ")
					self.Visit(t.Cases[0].Expression)
					w.WriteStringln(";")
				})
				w.WriteString("}(")
				self.Visit(t.Target)
				w.WriteString(")")
			case *dsl.TypePattern, *dsl.DiscardPattern:
				self.Visit(t.Target)
			default:
				panic(fmt.Sprintf("Unexpected pattern type %T", t.Cases[0].Pattern))
			}

		case *dsl.TypeConversionExpression:
			fmt.Fprintf(w, "static_cast<%s>(", common.TypeSyntax(t.Type))
			self.Visit(t.Expression)
			w.WriteString(")")
		default:
			panic(fmt.Sprintf("Unknown expression type '%T'", t))
		}
	})
}

func writeSwitchCaseOverUnion(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, visitor dsl.Visitor) {
	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		typeCheck := fmt.Sprintf("std::is_same_v<std::decay_t<decltype(%s)>, %s>", variableName, common.TypeSyntax(typePattern.Type))
		fmt.Fprintf(w, "if constexpr (%s) {\n", typeCheck)
		w.Indented(func() {
			if declarationIdentifier != "" {
				fmt.Fprintf(w, "%s const& %s = %s;\n", common.TypeSyntax(typePattern.Type), declarationIdentifier, variableName)
			}
			w.WriteString("return ")
			visitor.Visit(switchCase.Expression)
			w.WriteStringln(";")
		})
		w.WriteStringln("}")
	}

	switch t := switchCase.Pattern.(type) {
	case *dsl.TypePattern:
		writeTypeCase(t, "")
	case *dsl.DeclarationPattern:
		writeTypeCase(&t.TypePattern, common.FieldIdentifierName(t.Identifier))
	case *dsl.DiscardPattern:
		w.WriteString("return ")
		visitor.Visit(switchCase.Expression)
		w.WriteStringln(";")
	default:
		panic(fmt.Sprintf("Unknown pattern type '%T'", switchCase.Pattern))
	}
}

func writeSwitchCaseOverOptional(w *formatting.IndentedWriter, switchCase *dsl.SwitchCase, variableName string, isLastCase bool, visitor dsl.Visitor) {
	writeCore := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if declarationIdentifier != "" {
			fmt.Fprintf(w, "%s const& %s = %s.value();\n", common.TypeSyntax(typePattern.Type), declarationIdentifier, variableName)
		}
		w.WriteString("return ")
		visitor.Visit(switchCase.Expression)
		w.WriteStringln(";")
	}

	writeTypeCase := func(typePattern *dsl.TypePattern, declarationIdentifier string) {
		if isLastCase {
			writeCore(typePattern, declarationIdentifier)
			return
		}

		negationString := ""
		if typePattern.Type == nil {
			negationString = "!"
		}

		fmt.Fprintf(w, "if (%s%s.has_value()) {\n", negationString, variableName)
		w.Indented(func() {
			writeCore(typePattern, declarationIdentifier)
		})
		w.WriteStringln("}")
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
