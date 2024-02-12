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
				fmt.Fprintf(w, "function res = %s(obj)\n", fieldName)
				common.WriteBlockBody(w, func() {
					writeComputedFieldExpression(w, computedField.Expression, rec.Namespace)
					w.WriteStringln("")
				})
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
	panic(fmt.Sprintf("How about type equality expression for %s", dsl.TypeToShortSyntax(t, false)))
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

func writeComputedFieldExpression(w *formatting.IndentedWriter, expression dsl.Expression, contextNamespace string) {
	// TODO
}
