// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

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

	common.WriteComment(w, "pyright: reportUnusedClass=false")
	common.WriteComment(w, "pyright: reportUnusedImport=false")
	common.WriteComment(w, "pyright: reportUnknownArgumentType=false")
	common.WriteComment(w, "pyright: reportUnknownMemberType=false")
	common.WriteComment(w, "pyright: reportUnknownVariableType=false")

	w.WriteStringln(`
import collections.abc
import io
import typing

import numpy as np
import numpy.typing as npt

from .types import *
`)

	relativePath := ".."
	if ns.IsTopLevel {
		relativePath = "."
		w.WriteStringln("from .protocols import *")
	}

	fmt.Fprintf(w, "from %s import _binary\n", relativePath)
	fmt.Fprintf(w, "from %s import yardl_types as yardl\n\n", relativePath)

	if ns.IsTopLevel {
		writeProtocols(w, ns)
	}
	writeRecordSerializers(w, ns)

	binaryPath := path.Join(packageDir, "binary.py")
	return iocommon.WriteFileIfNeeded(binaryPath, b.Bytes(), 0644)
}

func writeRecordSerializers(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, td := range ns.TypeDefinitions {
		switch td := td.(type) {
		case *dsl.RecordDefinition:
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

			fmt.Fprintf(w, "class %s(%s_binary.RecordSerializer[%s]):\n", recordSerializerClassName(td, ns.Name), genericSpec, typeSyntax)
			w.Indented(func() {
				if len(td.TypeParameters) > 0 {
					typeParamSerializers := make([]string, 0, len(td.TypeParameters))
					for _, tp := range td.TypeParameters {
						typeParamSerializers = append(
							typeParamSerializers,
							fmt.Sprintf("%s: _binary.TypeSerializer[%s, %s]", typeDefinitionSerializer(tp, ns.Name), common.TypeParameterSyntax(tp, false), common.TypeParameterSyntax(tp, true)))
					}

					fmt.Fprintf(w, "def __init__(self, %s) -> None:\n", strings.Join(typeParamSerializers, ", "))
				} else {
					w.WriteStringln("def __init__(self) -> None:")
				}

				w.Indented(func() {
					fieldSerializers := make([]string, len(td.Fields))
					for i, field := range td.Fields {
						fieldSerializers[i] = fmt.Sprintf(`("%s", %s)`, common.FieldIdentifierName(field.Name), typeSerializer(field.Type, ns.Name, nil))
					}
					fmt.Fprintf(w, "super().__init__([%s])\n", strings.Join(fieldSerializers, ", "))
				})
				w.WriteStringln("")

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

				fmt.Fprintf(w, "def read(self, stream: _binary.CodedInputStream) -> %s:\n", typeSyntax)
				w.Indented(func() {
					w.WriteStringln("field_values = self._read(stream)")
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

func recordSerializerClassName(record *dsl.RecordDefinition, contextNamespace string) string {
	className := fmt.Sprintf("%sSerializer", formatting.ToPascalCase(record.Name))
	if record.Namespace != contextNamespace {
		className = fmt.Sprintf("%s.binary.%s", common.NamespaceIdentifierName(record.Namespace), className)
	}
	return className
}

func writeProtocols(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, p := range ns.Protocols {

		// writer
		fmt.Fprintf(w, "class %s(_binary.BinaryProtocolWriter, %s):\n", BinaryWriterName(p), common.AbstractWriterName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Binary writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			w.WriteStringln("def __init__(self, stream: typing.Union[typing.BinaryIO, str]) -> None:")
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
					serializer := typeSerializer(step.Type, ns.Name, nil)
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

			w.WriteStringln("def __init__(self, stream: typing.Union[io.BufferedReader, io.BytesIO, typing.BinaryIO, str], skip_completed_check: bool = False) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self, skip_completed_check)\n", common.AbstractReaderName(p))
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
					serializer := typeSerializer(step.Type, ns.Name, nil)
					fmt.Fprintf(w, "return %s.read(self._stream)\n", serializer)
				})
				w.WriteStringln("")
			}
		})
	}
}

func typeDefinitionSerializer(t dsl.TypeDefinition, contextNamespace string) string {
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

		elementSerializer := typeSerializer(baseType, contextNamespace, nil)
		return fmt.Sprintf("_binary.EnumSerializer(%s, %s)", elementSerializer, common.TypeSyntax(t, contextNamespace))
	case *dsl.RecordDefinition:
		serializerName := recordSerializerClassName(t, contextNamespace)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, 0, len(t.TypeArguments))
		for _, arg := range t.TypeArguments {
			typeArguments = append(typeArguments, typeSerializer(arg, contextNamespace, nil))
		}

		if len(typeArguments) == 0 {
			return fmt.Sprintf("%s()", serializerName)
		}

		return fmt.Sprintf("%s(%s)", serializerName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_serializer", formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeSerializer(t.Type, contextNamespace, t)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeSerializer(t dsl.Type, contextNamespace string, namedType *dsl.NamedType) string {
	switch t := t.(type) {
	case nil:
		return "_binary.none_serializer"
	case *dsl.SimpleType:
		return typeDefinitionSerializer(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		getScalarSerializer := func() string {
			if t.Cases.IsSingle() {
				return typeSerializer(t.Cases[0].Type, contextNamespace, namedType)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("_binary.OptionalSerializer(%s)", typeSerializer(t.Cases[1].Type, contextNamespace, namedType))
			}

			unionClassName, typeParameters := common.UnionClassName(t)
			if namedType != nil {
				unionClassName = namedType.Name
				if namedType.Namespace != contextNamespace {
					unionClassName = fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(namedType.Namespace), unionClassName)
				}
			}

			var classSyntax string
			if len(typeParameters) == 0 {
				classSyntax = unionClassName
			} else {
				classSyntax = fmt.Sprintf("%s[%s]", unionClassName, typeParameters)
			}
			options := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				if c.Type == nil {
					options[i] = "None"
				} else {
					options[i] = fmt.Sprintf("(%s.%s, %s)", classSyntax, formatting.ToPascalCase(c.Tag), typeSerializer(c.Type, contextNamespace, namedType))
				}
			}

			return fmt.Sprintf("_binary.UnionSerializer(%s, [%s])", unionClassName, strings.Join(options, ", "))

		}
		switch td := t.Dimensionality.(type) {
		case nil:
			return getScalarSerializer()
		case *dsl.Stream:
			return fmt.Sprintf("_binary.StreamSerializer(%s)", getScalarSerializer())
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("_binary.FixedVectorSerializer(%s, %d)", getScalarSerializer(), *td.Length)
			}

			return fmt.Sprintf("_binary.VectorSerializer(%s)", getScalarSerializer())
		case *dsl.Array:
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("_binary.FixedNDArraySerializer(%s, (%s,))", getScalarSerializer(), strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("_binary.NDArraySerializer(%s, %d)", getScalarSerializer(), len(*td.Dimensions))
			}

			return fmt.Sprintf("_binary.DynamicNDArraySerializer(%s)", getScalarSerializer())

		case *dsl.Map:
			keySerializer := typeSerializer(td.KeyType, contextNamespace, namedType)
			valueSerializer := typeSerializer(t.ToScalar(), contextNamespace, namedType)

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
