package binary

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
import numpy.typing as npt

from . import *
from . import _binary
from . import yardl_types as yardl

# pyright: reportUnusedClass=false
`)

	common.WriteTypeVars(w, ns)
	writeRecordSerializers(w, ns)
	writeProtocols(w, ns)

	definitionsPath := path.Join(packageDir, "binary.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeRecordSerializers(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.RecordDefinition:
			typeSyntax := common.TypeDefinitionSyntax(td, ns.Name, true)
			fmt.Fprintf(w, "class %s(_binary.RecordWriter[%s]):\n", recordRwClassName(td, true), typeSyntax)

			w.Indented(func() {
				if len(td.TypeParameters) > 0 {
					typeParamWriters := make([]string, len(td.TypeParameters))
					for i, tp := range td.TypeParameters {
						typeParamWriters[i] = fmt.Sprintf("%s: _binary.Writer[%s]", typeDefinitionRwCallable(tp, true, ns.Name), common.TypeDefinitionSyntax(tp, ns.Name, true))
					}

					fmt.Fprintf(w, "def __init__(self, %s) -> None:\n", strings.Join(typeParamWriters, ", "))
				} else {
					w.WriteStringln("def __init__(self) -> None:")
				}

				w.Indented(func() {
					fieldWriters := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldWriters[i] = typeRwCallable(field.Type, true, ns.Name)
					}
					fmt.Fprintf(w, "super().__init__([%s])\n", strings.Join(fieldWriters, ", "))
				})
				w.WriteStringln("")

				fmt.Fprintf(w, "def __call__(self, stream: _binary.CodedOutputStream, value: %s) -> None:\n", typeSyntax)
				w.Indented(func() {
					fieldAccesses := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldAccesses[i] = fmt.Sprintf("value.%s", common.FieldIdentifierName(field.Name))
					}
					fmt.Fprintf(w, "self._write(stream, %s)\n", strings.Join(fieldAccesses, ", "))
				})
				w.WriteStringln("")
			})
			w.WriteStringln("")
		}
	}
}

func recordRwClassName(record *dsl.RecordDefinition, write bool) string {
	return fmt.Sprintf("_%s%s", formatting.ToPascalCase(record.Name), noun(write))
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
				valueType := common.TypeSyntax(step.Type, ns.Name, false)
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
				valueType := common.TypeSyntax(step.Type, ns.Name, false)
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

func typeDefinitionRwCallable(t dsl.TypeDefinition, write bool, contextNamespace string) string {
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
	case *dsl.RecordDefinition:
		rwClassName := recordRwClassName(t, write)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", rwClassName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, len(t.TypeArguments))
		for i, arg := range t.TypeArguments {
			typeArguments[i] = typeRwCallable(arg, write, "")
		}
		return fmt.Sprintf("%s(%s)", rwClassName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_%s", verb(write), formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeRwCallable(t.Type, write, contextNamespace)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeRwCallable(t dsl.Type, write bool, contextNamespace string) string {
	switch t := t.(type) {
	case nil:
		return fmt.Sprintf("_binary.%s_none", verb(write))
	case *dsl.SimpleType:
		return typeDefinitionRwCallable(t.ResolvedDefinition, write, contextNamespace)
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
					typeSyntax = common.TypeSyntax(c.Type, contextNamespace, true)
				}
				options[i] = fmt.Sprintf("(%s, %s)", typeSyntax, typeRwCallable(c.Type, write, contextNamespace))
			}

			return fmt.Sprintf("_binary.Union%s([%s])", noun(write), strings.Join(options, ", "))

		}()
		switch td := t.Dimensionality.(type) {
		case nil:
			return scalarCallable
		case *dsl.Stream:
			return fmt.Sprintf("_binary.Stream%s(%s)", noun(write), scalarCallable)
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("_binary.FixedVector%s(%s, %d)", noun(write), scalarCallable, *td.Length)
			}

			return fmt.Sprintf("_binary.Vector%s(%s)", noun(write), scalarCallable)
		case *dsl.Array:
			dtype := common.TypeDTypeSyntax(t.ToScalar())
			triviallySerializable := boolSyntax(isTypePotentiallyTriviallySerializable(t.ToScalar()))
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("_binary.FixedNDArray%s(%s, %s, %s, (%s,))", noun(write), scalarCallable, dtype, triviallySerializable, strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("_binary.NDArray%s(%s, %s, %s, %d)", noun(write), scalarCallable, dtype, triviallySerializable, len(*td.Dimensions))
			}

			return fmt.Sprintf("_binary.DynamicNDArray%s(%s, %s, %s)", noun(write), scalarCallable, dtype, triviallySerializable)

		case *dsl.Map:
			keyCallable := typeRwCallable(td.KeyType, write, contextNamespace)
			valueCallable := typeRwCallable(t.ToScalar(), write, contextNamespace)

			return fmt.Sprintf("_binary.Map%s(%s, %s)", noun(write), keyCallable, valueCallable)
		default:
			panic(fmt.Sprintf("Not implemented %T", t.Dimensionality))
		}
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func isTypeDefinitionPotentiallyTriviallySerializable(t dsl.TypeDefinition) bool {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		switch t {
		case dsl.Uint8, dsl.Int8, dsl.Float32, dsl.Float64, dsl.ComplexFloat32, dsl.ComplexFloat64:
			return true
		}
	case *dsl.EnumDefinition:
		if t.BaseType != nil {
			return isTypePotentiallyTriviallySerializable(t.BaseType)
		}
	case *dsl.RecordDefinition:
		for _, f := range t.Fields {
			if !isTypePotentiallyTriviallySerializable(f.Type) {
				return false
			}
		}

		return true
	}

	return false
}

func isTypePotentiallyTriviallySerializable(t dsl.Type) bool {
	switch t := t.(type) {
	case *dsl.SimpleType:
		return isTypeDefinitionPotentiallyTriviallySerializable(t.ResolvedDefinition)
	}

	return false
}

func boolSyntax(b bool) string {
	if b {
		return "True"
	}
	return "False"
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
