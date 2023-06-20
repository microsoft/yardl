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
	w.WriteStringln(`# pyright: reportUnusedClass=false

import collections.abc
import datetime
import io
import typing
import numpy as np
import numpy.typing as npt

from . import *
from . import _binary
from . import yardl_types as yardl
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
			fmt.Fprintf(w, "class %s(_binary.RecordDescriptor[%s]):\n", recordDescriptorClassName(td), typeSyntax)

			w.Indented(func() {
				if len(td.TypeParameters) > 0 {
					typeParamDescriptors := make([]string, len(td.TypeParameters))
					for i, tp := range td.TypeParameters {
						typeParamDescriptors[i] = fmt.Sprintf("%s: _binary.TypeDescriptor[%s]", typeDefinitionDescriptor(tp, ns.Name), common.TypeDefinitionSyntax(tp, ns.Name, true))
					}

					fmt.Fprintf(w, "def __init__(self, %s) -> None:\n", strings.Join(typeParamDescriptors, ", "))
				} else {
					w.WriteStringln("def __init__(self) -> None:")
				}

				w.Indented(func() {
					fieldDescriptors := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldDescriptors[i] = fmt.Sprintf(`("%s", %s)`, common.FieldIdentifierName(field.Name), typeDescriptor(field.Type, ns.Name))
					}
					fmt.Fprintf(w, "super().__init__([%s])\n", strings.Join(fieldDescriptors, ", "))
				})
				w.WriteStringln("")

				fmt.Fprintf(w, "def write(self, stream: _binary.CodedOutputStream, value: %s) -> None:\n", typeSyntax)
				w.Indented(func() {
					fieldAccesses := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldAccesses[i] = fmt.Sprintf("value.%s", common.FieldIdentifierName(field.Name))
					}
					fmt.Fprintf(w, "self._write(stream, %s)\n", strings.Join(fieldAccesses, ", "))
				})
				w.WriteStringln("")

				fmt.Fprintf(w, "def read(self, stream: _binary.CodedInputStream, read_as_numpy: Types) -> %s:\n", typeSyntax)
				w.Indented(func() {
					w.WriteStringln("field_values = self._read(stream, read_as_numpy)")
					args := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						args[i] = fmt.Sprintf("%s=field_values[%d]", common.FieldIdentifierName(field.Name), i)
					}

					fmt.Fprintf(w, "return %s(%s)\n", typeSyntax, strings.Join(args, ", "))

				})
				w.WriteStringln("")
			})
			w.WriteStringln("")
		}
	}
}

func recordDescriptorClassName(record *dsl.RecordDefinition) string {
	return fmt.Sprintf("_%sDescriptor", formatting.ToPascalCase(record.Name))
}

func writeProtocols(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, p := range ns.Protocols {

		// writer
		fmt.Fprintf(w, "class %s(_binary.BinaryProtocolWriter, %s):\n", BinaryWriterName(p), common.AbstractWriterName(p))
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
					descriptor := typeDescriptor(step.Type, ns.Name)
					fmt.Fprintf(w, "%s.write(self._stream, value)\n", descriptor)
				})
				w.WriteStringln("")
			}
		})

		w.WriteStringln("")

		// reader
		fmt.Fprintf(w, "class %s(_binary.BinaryProtocolReader, %s):\n", BinaryReaderName(p), common.AbstractReaderName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Binary writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			w.WriteStringln("def __init__(self, stream: io.BufferedReader | str, read_as_numpy: Types) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self, read_as_numpy)\n", common.AbstractReaderName(p))
				fmt.Fprintf(w, "_binary.BinaryProtocolReader.__init__(self, stream, %s.schema)\n", common.AbstractReaderName(p))
			})
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name, false)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}

				fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
				w.Indented(func() {
					descriptor := typeDescriptor(step.Type, ns.Name)
					fmt.Fprintf(w, "return %s.read(self._stream, self._read_as_numpy)\n", descriptor)
				})
				w.WriteStringln("")
			}
		})
	}
}

func typeDefinitionDescriptor(t dsl.TypeDefinition, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return fmt.Sprintf("_binary.%s_descriptor", strings.ToLower(string(t)))
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		elementDescriptor := typeDescriptor(baseType, contextNamespace)
		return fmt.Sprintf("_binary.EnumDescriptor(%s, %s)", elementDescriptor, common.TypeDefinitionSyntax(t, contextNamespace, false))
	case *dsl.RecordDefinition:
		rwClassName := recordDescriptorClassName(t)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", rwClassName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, len(t.TypeArguments))
		for i, arg := range t.TypeArguments {
			typeArguments[i] = typeDescriptor(arg, "")
		}
		return fmt.Sprintf("%s(%s)", rwClassName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_descriptor", formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeDescriptor(t.Type, contextNamespace)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeDescriptor(t dsl.Type, contextNamespace string) string {
	switch t := t.(type) {
	case nil:
		return "_binary.none_descriptor"
	case *dsl.SimpleType:
		return typeDefinitionDescriptor(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		scalarDescriptor := func() string {
			if t.Cases.IsSingle() {
				return typeDescriptor(t.Cases[0].Type, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("_binary.OptionalDescriptor(%s)", typeDescriptor(t.Cases[1].Type, contextNamespace))
			}

			options := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				var typeSyntax string
				if c.Type == nil {
					typeSyntax = "None.__class__"
				} else {
					typeSyntax = common.TypeSyntax(c.Type, contextNamespace, true)
				}
				options[i] = fmt.Sprintf("(%s, %s)", typeSyntax, typeDescriptor(c.Type, contextNamespace))
			}

			return fmt.Sprintf("_binary.UnionDescriptor([%s])", strings.Join(options, ", "))

		}()
		switch td := t.Dimensionality.(type) {
		case nil:
			return scalarDescriptor
		case *dsl.Stream:
			return fmt.Sprintf("_binary.StreamDescriptor(%s)", scalarDescriptor)
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("_binary.FixedVectorDescriptor(%s, %d)", scalarDescriptor, *td.Length)
			}

			return fmt.Sprintf("_binary.VectorDescriptor(%s)", scalarDescriptor)
		case *dsl.Array:
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("_binary.FixedNDArrayDescriptor(%s, (%s,))", scalarDescriptor, strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("_binary.NDArrayDescriptor(%s, %d)", scalarDescriptor, len(*td.Dimensions))
			}

			return fmt.Sprintf("_binary.DynamicNDArrayDescriptor(%s)", scalarDescriptor)

		case *dsl.Map:
			keyDescriptor := typeDescriptor(td.KeyType, contextNamespace)
			valueDescriptor := typeDescriptor(t.ToScalar(), contextNamespace)

			return fmt.Sprintf("_binary.MapDescriptor(%s, %s)", keyDescriptor, valueDescriptor)
		default:
			panic(fmt.Sprintf("Not implemented %T", t.Dimensionality))
		}
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func BinaryWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sWriter", formatting.ToPascalCase(p.Name))
}

func BinaryReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Binary%sReader", formatting.ToPascalCase(p.Name))
}
