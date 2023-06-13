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

func WriteTypes(ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import dataclasses
import datetime
import enum
import typing
import numpy as np
from . import yardl_types as yardl
`)

	writeTypes(w, ns)

	definitionsPath := path.Join(packageDir, "types.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeTypes(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	typeVars := make(map[string]any)
	for _, td := range ns.TypeDefinitions {
		for _, tp := range td.GetDefinitionMeta().TypeParameters {
			identifier := common.TypeIdentifierName(tp.Name)
			if _, ok := typeVars[identifier]; !ok {
				typeVars[identifier] = nil
				fmt.Fprintf(w, "%s = typing.TypeVar('%s')\n", identifier, identifier)
			}
		}
	}
	if len(typeVars) > 0 {
		w.WriteStringln("")
	}

	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.EnumDefinition:
			writeEnum(w, td)
		case *dsl.RecordDefinition:
			writeRecord(w, td)
		case *dsl.NamedType:
			writeNamedType(w, td)
		default:
			panic(fmt.Sprintf("unsupported type definition: %T", td))
		}
	}
}

func writeNamedType(w *formatting.IndentedWriter, td *dsl.NamedType) {
	fmt.Fprintf(w, "%s = %s\n", common.TypeDefinitionSyntax(td, td.Namespace), common.TypeSyntax(td.Type, td.Namespace))
	common.WriteDocstring(w, td.Comment)
	w.Indent().WriteStringln("")
}

func writeRecord(w *formatting.IndentedWriter, rec *dsl.RecordDefinition) {
	w.WriteStringln("@dataclasses.dataclass(slots=True, kw_only=True)")
	fmt.Fprintf(w, "class %s%s:\n", common.TypeDefinitionSyntax(rec, rec.Namespace), GetGenericBase(rec))
	w.Indented(func() {
		common.WriteDocstring(w, rec.Comment)
		for _, field := range rec.Fields {
			fmt.Fprintf(w, "%s: %s", common.FieldIdentifierName(field.Name), common.TypeSyntax(field.Type, rec.Namespace))
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

func writeEnum(w *formatting.IndentedWriter, enum *dsl.EnumDefinition) {
	var baseType string
	if enum.IsFlags {
		baseType = "enum.Flag"
	} else {
		baseType = "enum.Enum"
	}
	fmt.Fprintf(w, "class %s(%s):\n", common.TypeDefinitionSyntax(enum, enum.Namespace), baseType)

	w.Indented(func() {
		common.WriteDocstring(w, enum.Comment)
		for _, value := range enum.Values {
			fmt.Fprintf(w, "%s = %d\n", common.EnumValueIdentifierName(value.Symbol), &value.IntegerValue)
			common.WriteDocstring(w, value.Comment)
		}
	})
	w.WriteStringln("")
}
