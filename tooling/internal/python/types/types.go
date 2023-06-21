package types

import (
	"bytes"
	"fmt"
	"path"
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
	w.WriteStringln(`
import dataclasses
import datetime
import enum
import typing
import numpy as np
import numpy.typing as npt
from . import yardl_types as yardl
`)

	writeTypes(w, st, ns)

	definitionsPath := path.Join(packageDir, "types.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeTypes(w *formatting.IndentedWriter, st dsl.SymbolTable, ns *dsl.Namespace) {
	common.WriteTypeVars(w, ns)

	typeDefinitionsUsedInArrays := getTypeDefinitionsReachableFromArrays(ns, st)

	for _, td := range ns.TypeDefinitions {
		generateDTypeMethod := typeDefinitionsUsedInArrays[td]
		switch td := td.(type) {
		case *dsl.EnumDefinition:
			writeEnum(w, td, generateDTypeMethod)
		case *dsl.RecordDefinition:
			writeRecord(w, td, generateDTypeMethod)
		case *dsl.NamedType:
			writeNamedType(w, td)
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}
	}
}

func getTypeDefinitionsReachableFromArrays(ns *dsl.Namespace, st dsl.SymbolTable) map[dsl.TypeDefinition]bool {
	typeDefinitionsUsedInArrays := make(map[dsl.TypeDefinition]bool)
	dsl.VisitWithContext(ns, false, func(self dsl.VisitorWithContext[bool], node dsl.Node, inArray bool) {
		switch t := node.(type) {
		case *dsl.ProtocolDefinition:
			break
		case dsl.PrimitiveDefinition:
			break
		case *dsl.GenericTypeParameter:
			break
		case dsl.TypeDefinition:
			if inArray {
				typeDefinitionsUsedInArrays[t] = true
			}
		case *dsl.GeneralizedType:
			switch t.Dimensionality.(type) {
			case *dsl.Array:
				inArray = true
			}

		case *dsl.SimpleType:
			self.Visit(st.GetGenericTypeDefinition(t.ResolvedDefinition), inArray)
			for _, typeArg := range t.ResolvedDefinition.GetDefinitionMeta().TypeParameters {
				self.Visit(st.GetGenericTypeDefinition(typeArg), inArray)
			}
		}

		self.VisitChildren(node, inArray)
	})
	return typeDefinitionsUsedInArrays
}

func writeNamedType(w *formatting.IndentedWriter, td *dsl.NamedType) {
	fmt.Fprintf(w, "%s = %s\n", common.TypeDefinitionSyntax(td, td.Namespace, false), common.TypeSyntax(td.Type, td.Namespace, false))
	common.WriteDocstring(w, td.Comment)
	w.Indent().WriteStringln("")
}

func writeRecord(w *formatting.IndentedWriter, rec *dsl.RecordDefinition, generateDTypeMethod bool) {
	w.WriteStringln("@dataclasses.dataclass(slots=True, kw_only=True)")
	fmt.Fprintf(w, "class %s%s:\n", common.TypeDefinitionSyntax(rec, rec.Namespace, false), GetGenericBase(rec))
	w.Indented(func() {
		common.WriteDocstring(w, rec.Comment)
		for _, field := range rec.Fields {
			fmt.Fprintf(w, "%s: %s", common.FieldIdentifierName(field.Name), common.TypeSyntax(field.Type, rec.Namespace, true))
			if gt, ok := field.Type.(*dsl.GeneralizedType); ok && gt.Cases.HasNullOption() {
				w.WriteStringln(" = None")
			} else {
				w.WriteStringln("")
			}

			common.WriteDocstring(w, field.Comment)
		}

		if len(rec.Fields) == 0 {
			w.WriteStringln("pass")
		}

		if generateDTypeMethod {
			writeDTypeMethod(w, rec)
		}
	})
	w.WriteStringln("")
}

func GetGenericBase(t dsl.TypeDefinition) string {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) == 0 {
		return ""
	}

	var typeParams []string
	for _, tp := range meta.TypeParameters {
		typeParams = append(typeParams, common.TypeIdentifierName(tp.Name))
	}

	return fmt.Sprintf("(typing.Generic[%s])", strings.Join(typeParams, ", "))
}

func writeEnum(w *formatting.IndentedWriter, enum *dsl.EnumDefinition, generateDTypeMethod bool) {
	var baseType string
	if enum.IsFlags {
		baseType = "enum.Flag"
	} else {
		baseType = "enum.Enum"
	}
	fmt.Fprintf(w, "class %s(%s):\n", common.TypeDefinitionSyntax(enum, enum.Namespace, false), baseType)

	w.Indented(func() {
		common.WriteDocstring(w, enum.Comment)
		for _, value := range enum.Values {
			fmt.Fprintf(w, "%s = %d\n", common.EnumValueIdentifierName(value.Symbol), &value.IntegerValue)
			common.WriteDocstring(w, value.Comment)
		}
	})
	w.WriteStringln("")

	if generateDTypeMethod {
		writeDTypeMethod(w, enum)
	}
}

func writeDTypeMethod(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	w.WriteStringln("")
	meta := t.GetDefinitionMeta()
	typeParams := make([]string, len(meta.TypeParameters))
	for i, tp := range t.GetDefinitionMeta().TypeParameters {
		typeParams[i] = fmt.Sprintf("%s_dtype: npt.DTypeLike", formatting.ToSnakeCase(tp.Name))
	}

	w.WriteStringln("@staticmethod")
	fmt.Fprintf(w, "def dtype(%s) -> npt.DTypeLike:\n", strings.Join(typeParams, ", "))
	w.Indented(func() {
		fmt.Fprintf(w, "return %s\n", common.TypeDefinitionDTypeSyntax(t))
	})
}
