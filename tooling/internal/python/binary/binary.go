package binary

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

func WriteBinary(ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import abc
import collections.abc
import datetime
import typing
import numpy as np

from . import *
from . import _binary
from . import yardl_types as yardl
`)

	writeProtocols(w, ns)

	definitionsPath := path.Join(packageDir, "binary.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeProtocols(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, p := range ns.Protocols {

		// writer
		fmt.Fprintf(w, "class %s(%s, _binary.BinaryProtocolWriter):\n", BinaryWriterName(p), common.AbstractWriterName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Binary writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			w.WriteStringln("def __init__(self, stream: typing.BinaryIO | str) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self)\n", common.AbstractWriterName(p))
				fmt.Fprintf(w, "_binary.BinaryProtocolWriter.__init__(self, stream, %s.schema)\n", common.AbstractWriterName(p))
			})
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}
				fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteImplMethodName(step), valueType)
				w.Indented(func() {
					callable := typeRwCallable(step.Type, true, ns.Name)
					fmt.Fprintf(w, "%s(self._stream, value)\n", callable)
				})
				w.WriteStringln("")
			}
		})

		w.WriteStringln("")

		// reader
		fmt.Fprintf(w, "class %s(%s):\n", BinaryReaderName(p), common.AbstractReaderName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Binary writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}

				fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
				w.Indented(func() {
					w.WriteStringln("raise NotImplementedError()")
				})
				w.WriteStringln("")
			}
		})
	}
}

func typeDefinitionRwCallable(t dsl.TypeDefinition, write bool) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		suffix := strings.ToLower(string(t))
		return fmt.Sprintf("_binary.%s_%s", verb(write), suffix)
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		baseRwCallable := typeRwCallable(baseType, write, "")
		return fmt.Sprintf("_binary.Enum%s(%s)", noun(write), baseRwCallable)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeRwCallable(t dsl.Type, write bool, contextNamespace string) string {
	switch t := t.(type) {
	case nil:
		return fmt.Sprintf("_binary.%s_none", verb(write))
	case *dsl.SimpleType:
		return typeDefinitionRwCallable(t.ResolvedDefinition, write)
	case *dsl.GeneralizedType:
		scalarCallable := func() string {
			if t.Cases.IsSingle() {
				return typeRwCallable(t.Cases[0].Type, write, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("_binary.Optional%s(%s)", noun(write), typeRwCallable(t.Cases[1].Type, write, contextNamespace))
			}

			options := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				var typeSyntax string
				if c.Type == nil {
					typeSyntax = "None.__class__"
				} else {
					typeSyntax = common.TypeSyntax(c.Type, contextNamespace)
				}
				options[i] = fmt.Sprintf("(%s, %s)", typeSyntax, typeRwCallable(c.Type, write, contextNamespace))
			}

			return fmt.Sprintf("_binary.Union%s([%s])", noun(write), strings.Join(options, ", "))

		}()
		switch t.Dimensionality.(type) {
		case nil:
			return scalarCallable
		case *dsl.Stream:
			return fmt.Sprintf("_binary.Stream%s(%s)", noun(write), scalarCallable)
		default:
			panic(fmt.Sprintf("Not implemented %T", t.Dimensionality))
		}
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func verb(write bool) string {
	if write {
		return "write"
	}
	return "read"
}

func noun(write bool) string {
	if write {
		return "Writer"
	}
	return "Reader"
}

func BinaryWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sWriter", formatting.ToPascalCase(p.Name))
}

func BinaryReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sReader", formatting.ToPascalCase(p.Name))
}
