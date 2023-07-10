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
	writeProtocols(w, ns)
	writeRecordSerializers(w, ns)

	definitionsPath := path.Join(packageDir, "binary.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeRecordSerializers(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.RecordDefinition:
			writeInit := func(numpy bool) {
				if len(td.TypeParameters) > 0 {
					typeParamSerializers := make([]string, 0, len(td.TypeParameters))
					for _, tp := range td.TypeParameters {
						typeParamSerializers = append(typeParamSerializers, fmt.Sprintf("%s: _binary.TypeSerializer[%s, %s]", typeDefinitionSerializer(tp, false, ns.Name), common.TypeParameterSyntax(tp, false), common.TypeParameterSyntax(tp, true)))
					}

					fmt.Fprintf(w, "def __init__(self, %s) -> None:\n", strings.Join(typeParamSerializers, ", "))
				} else {
					w.WriteStringln("def __init__(self) -> None:")
				}

				w.Indented(func() {
					fieldSerializers := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldSerializers[i] = fmt.Sprintf(`("%s", %s)`, common.FieldIdentifierName(field.Name), typeSerializer(field.Type, false, ns.Name))
					}
					fmt.Fprintf(w, "super().__init__([%s])\n", strings.Join(fieldSerializers, ", "))
				})
				w.WriteStringln("")
			}

			typeSyntax := common.TypeSyntax(td, ns.Name)
			var genericSpec string
			if len(td.TypeParameters) > 0 {
				params := make([]string, 2*len(td.TypeParameters))
				for i, tp := range td.TypeParameters {
					params[2*i] = common.TypeParameterSyntax(tp, false)
					params[2*i+1] = common.TypeParameterSyntax(tp, true)
				}
				genericSpec = fmt.Sprintf("typing.Generic[%s], ", strings.Join(params, ", "))
			} else {
				genericSpec = ""
			}

			fmt.Fprintf(w, "class %s(%s_binary.RecordSerializer[%s]):\n", recordSerializerClassName(td, false), genericSpec, typeSyntax)
			w.Indented(func() {
				writeInit(false)

				fmt.Fprintf(w, "def write(self, stream: _binary.CodedOutputStream, value: %s) -> None:\n", typeSyntax)
				w.Indented(func() {
					w.WriteStringln("if isinstance(value, np.void):")
					w.Indented(func() {
						w.WriteStringln("self.write_numpy(stream, value)")
						w.WriteStringln("return")
					})

					fieldAccesses := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldAccesses[i] = fmt.Sprintf("value.%s", common.FieldIdentifierName(field.Name))
					}
					fmt.Fprintf(w, "self._write(stream, %s)\n", strings.Join(fieldAccesses, ", "))
				})
				w.WriteStringln("")

				fmt.Fprintf(w, "def write_numpy(self, stream: _binary.CodedOutputStream, value: np.void) -> None:\n")
				w.Indented(func() {
					fieldAccesses := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldAccesses[i] = fmt.Sprintf(`value['%s']`, common.FieldIdentifierName(field.Name))
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

func recordSerializerClassName(record *dsl.RecordDefinition, numpy bool) string {
	if numpy {
		return fmt.Sprintf("_%s_NumpySerializer", formatting.ToPascalCase(record.Name))
	}
	return fmt.Sprintf("_%sSerializer", formatting.ToPascalCase(record.Name))
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
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}
				fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteImplMethodName(step), valueType)
				w.Indented(func() {
					serializer := typeSerializer(step.Type, false, ns.Name)
					fmt.Fprintf(w, "%s.write(self._stream, value)\n", serializer)
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

			w.WriteStringln("def __init__(self, stream: io.BufferedReader | io.BytesIO | typing.BinaryIO | str, read_as_numpy: Types = Types.NONE) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self, read_as_numpy)\n", common.AbstractReaderName(p))
				fmt.Fprintf(w, "_binary.BinaryProtocolReader.__init__(self, stream, %s.schema)\n", common.AbstractReaderName(p))
			})
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}

				fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
				w.Indented(func() {
					serializer := typeSerializer(step.Type, false, ns.Name)
					fmt.Fprintf(w, "return %s.read(self._stream, self._read_as_numpy)\n", serializer)
				})
				w.WriteStringln("")
			}
		})
	}
}

func typeDefinitionSerializer(t dsl.TypeDefinition, numpy bool, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return fmt.Sprintf("_binary.%s_serializer", strings.ToLower(string(t)))
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		elementSerializer := typeSerializer(baseType, false, contextNamespace)
		return fmt.Sprintf("_binary.EnumSerializer(%s, %s)", elementSerializer, common.TypeSyntax(t, contextNamespace))
	case *dsl.RecordDefinition:
		rwClassName := recordSerializerClassName(t, false)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", rwClassName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, 0, len(t.TypeArguments))
		for _, arg := range t.TypeArguments {
			typeArguments = append(typeArguments, typeSerializer(arg, false, contextNamespace))
		}

		if len(typeArguments) == 0 {
			return fmt.Sprintf("%s()", rwClassName)
		}

		return fmt.Sprintf("%s(%s)", rwClassName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_serializer", formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeSerializer(t.Type, false, contextNamespace)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeSerializer(t dsl.Type, numpy bool, contextNamespace string) string {
	switch t := t.(type) {
	case nil:
		return "_binary.none_serializer"
	case *dsl.SimpleType:
		return typeDefinitionSerializer(t.ResolvedDefinition, numpy, contextNamespace)
	case *dsl.GeneralizedType:
		getScalarSerializer := func(numpy bool) string {
			if t.Cases.IsSingle() {
				return typeSerializer(t.Cases[0].Type, numpy, contextNamespace)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("_binary.OptionalSerializer(%s)", typeSerializer(t.Cases[1].Type, numpy, contextNamespace))
			}

			options := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				if c.Type == nil {
					options[i] = "None"
				} else {
					options[i] = fmt.Sprintf("(\"%s\", %s)", c.Tag, typeSerializer(c.Type, numpy, contextNamespace))
				}
			}

			return fmt.Sprintf("_binary.UnionSerializer([%s])", strings.Join(options, ", "))

		}
		switch td := t.Dimensionality.(type) {
		case nil:
			return getScalarSerializer(numpy)
		case *dsl.Stream:
			return fmt.Sprintf("_binary.StreamSerializer(%s)", getScalarSerializer(numpy))
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("_binary.FixedVectorSerializer(%s, %d)", getScalarSerializer(numpy), *td.Length)
			}

			return fmt.Sprintf("_binary.VectorSerializer(%s)", getScalarSerializer(numpy))
		case *dsl.Array:
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("_binary.FixedNDArraySerializer(%s, (%s,))", getScalarSerializer(true), strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("_binary.NDArraySerializer(%s, %d)", getScalarSerializer(true), len(*td.Dimensions))
			}

			return fmt.Sprintf("_binary.DynamicNDArraySerializer(%s)", getScalarSerializer(true))

		case *dsl.Map:
			keySerializer := typeSerializer(td.KeyType, numpy, contextNamespace)
			valueSerializer := typeSerializer(t.ToScalar(), numpy, contextNamespace)

			return fmt.Sprintf("_binary.MapSerializer(%s, %s)", keySerializer, valueSerializer)
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
