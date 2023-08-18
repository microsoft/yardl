package ndjson

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/ndjsoncommon"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteNDJson(ns *dsl.Namespace, packageDir string) error {
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
from . import _ndjson
from . import yardl_types as yardl
`)

	common.WriteTypeVars(w, ns)
	writeConverters(w, ns)
	writeProtocols(w, ns)

	ndjsonPath := path.Join(packageDir, "ndjson.py")
	return iocommon.WriteFileIfNeeded(ndjsonPath, b.Bytes(), 0644)
}

func writeConverters(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, t := range ns.TypeDefinitions {
		switch t := t.(type) {
		case *dsl.EnumDefinition:
			writeEnumMaps(t, w, ns)
		}
	}
}

func writeEnumMaps(t *dsl.EnumDefinition, w *formatting.IndentedWriter, ns *dsl.Namespace) {
	name_to_value_map_name := enumNameToValueMapName(t)
	fmt.Fprintf(w, "%s = {\n", name_to_value_map_name)
	w.Indented(func() {
		for _, v := range t.Values {
			fmt.Fprintf(w, "\"%s\": %s.%s,\n", v.Symbol, common.TypeSyntax(t, ns.Name), common.EnumValueIdentifierName(v.Symbol))
		}
	})
	fmt.Fprintf(w, "}\n")

	value_to_name_map_name := enumValueToNameMapName(t)
	fmt.Fprintf(w, "%s = {v: n for n, v in %s.items()}\n\n", value_to_name_map_name, name_to_value_map_name)
}

func enumNameToValueMapName(t *dsl.EnumDefinition) string {
	return fmt.Sprintf("_%s_name_to_value_map", formatting.ToSnakeCase(t.Name))
}

func enumValueToNameMapName(t *dsl.EnumDefinition) string {
	return fmt.Sprintf("_%s_value_to_name_map", formatting.ToSnakeCase(t.Name))
}

func recordConverterClassName(record *dsl.RecordDefinition) string {
	return fmt.Sprintf("_%sConverter", formatting.ToPascalCase(record.Name))
}

func writeProtocols(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, p := range ns.Protocols {

		// writer
		fmt.Fprintf(w, "class %s(_ndjson.NDJsonProtocolWriter, %s):\n", NDJsonWriterName(p), common.AbstractWriterName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("NDJson writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			w.WriteStringln("def __init__(self, stream: typing.Union[typing.TextIO, str]) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self)\n", common.AbstractWriterName(p))
				fmt.Fprintf(w, "_ndjson.NDJsonProtocolWriter.__init__(self, stream, %s.schema)\n", common.AbstractWriterName(p))
			})
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					fmt.Fprintf(w, "def %s(self, value: collections.abc.Iterable[%s]) -> None:\n", common.ProtocolWriteImplMethodName(step), valueType)
					w.Indented(func() {
						fmt.Fprintf(w, "converter = %s\n", typeConverter(step.Type, ns.Name, nil))
						fmt.Fprintf(w, "for item in value:\n")
						w.Indented(func() {
							fmt.Fprintf(w, "json_item = converter.to_json(item)\n")
							fmt.Fprintf(w, "self._write_json_line({\"%s\": json_item})\n", step.Name)
						})
					})
				} else {
					fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteImplMethodName(step), valueType)
					w.Indented(func() {
						converter := typeConverter(step.Type, ns.Name, nil)
						fmt.Fprintf(w, "json_value = %s.to_json(value)\n", converter)
						fmt.Fprintf(w, "self._write_json_line({\"%s\": json_value})\n", step.Name)
					})
				}
				w.WriteStringln("")
			}
		})

		w.WriteStringln("")

		// reader
		fmt.Fprintf(w, "class %s(_ndjson.NDJsonProtocolReader, %s):\n", NDJsonReaderName(p), common.AbstractReaderName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("NDJson writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			w.WriteStringln("def __init__(self, stream: typing.Union[io.BufferedReader, typing.TextIO, str]) -> None:")
			w.Indented(func() {
				fmt.Fprintf(w, "%s.__init__(self)\n", common.AbstractReaderName(p))
				fmt.Fprintf(w, "_ndjson.NDJsonProtocolReader.__init__(self, stream, %s.schema)\n", common.AbstractReaderName(p))
			})
			w.WriteStringln("")

			for _, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					fmt.Fprintf(w, "def %s(self) -> collections.abc.Iterable[%s]:\n", common.ProtocolReadImplMethodName(step), valueType)
					w.Indented(func() {
						fmt.Fprintf(w, "converter = %s\n", typeConverter(step.Type, ns.Name, nil))
						fmt.Fprintf(w, "while (json_object := self._read_json_line(\"%s\", False)) is not None:\n", step.Name)
						w.Indented(func() {
							fmt.Fprintf(w, "yield converter.from_json(json_object)\n")
						})
					})

				} else {
					fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
					w.Indented(func() {
						converter := typeConverter(step.Type, ns.Name, nil)
						fmt.Fprintf(w, "json_object = self._read_json_line(\"%s\", True)\n", step.Name)
						fmt.Fprintf(w, "return %s.from_json(json_object)\n", converter)
					})
				}

				w.WriteStringln("")
			}
		})
	}
}

func typeDefinitionConverter(t dsl.TypeDefinition, contextNamespace string) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		return fmt.Sprintf("_ndjson.%s_converter", strings.ToLower(string(t)))
	case *dsl.EnumDefinition:
		var baseType dsl.Type
		if t.BaseType != nil {
			baseType = t.BaseType
		} else {
			baseType = dsl.Int32Type
		}

		var className string
		if t.IsFlags {
			className = "_ndjson.FlagsConverter"
		} else {
			className = "_ndjson.EnumConverter"
		}

		return fmt.Sprintf("%s(%s, %s, %s, %s)", className, common.TypeSyntax(t, contextNamespace), common.TypeSyntax(baseType, contextNamespace), enumNameToValueMapName(t), enumValueToNameMapName(t))
	case *dsl.RecordDefinition:
		converterName := recordConverterClassName(t)
		if len(t.TypeParameters) == 0 {
			return fmt.Sprintf("%s()", converterName)
		}
		if len(t.TypeArguments) == 0 {
			panic("Expected type arguments")
		}

		typeArguments := make([]string, 0, len(t.TypeArguments))
		for _, arg := range t.TypeArguments {
			typeArguments = append(typeArguments, typeConverter(arg, contextNamespace, nil))
		}

		if len(typeArguments) == 0 {
			return fmt.Sprintf("%s()", converterName)
		}

		return fmt.Sprintf("%s(%s)", converterName, strings.Join(typeArguments, ", "))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s_converter", formatting.ToSnakeCase(t.Name))
	case *dsl.NamedType:
		return typeConverter(t.Type, contextNamespace, t)
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func typeConverter(t dsl.Type, contextNamespace string, namedType *dsl.NamedType) string {
	switch t := t.(type) {
	case nil:
		return "_ndjson.none_converter"
	case *dsl.SimpleType:
		return typeDefinitionConverter(t.ResolvedDefinition, contextNamespace)
	case *dsl.GeneralizedType:
		getScalarConverter := func() string {
			if t.Cases.IsSingle() {
				return typeConverter(t.Cases[0].Type, contextNamespace, namedType)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("_ndjson.OptionalConverter(%s)", typeConverter(t.Cases[1].Type, contextNamespace, namedType))
			}

			unionClassName, typeParameters := common.UnionClassName(t)
			if namedType != nil {
				unionClassName = namedType.Name
			}

			var classSyntax string
			if len(typeParameters) == 0 {
				classSyntax = unionClassName
			} else {
				classSyntax = fmt.Sprintf("%s[%s]", unionClassName, typeParameters)
			}

			simplfied := "True"
			var possibleTypes ndjsoncommon.JsonDataType
			options := make([]string, len(t.Cases))
			for i, c := range t.Cases {
				if c.Type == nil {
					options[i] = "None"
				} else {
					jsonTypes := ndjsoncommon.GetJsonDataType(c.Type)
					jsonTypeStrings := make([]string, 0, 1)
					if jsonTypes&ndjsoncommon.JsonNull != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "None")
					}
					if jsonTypes&ndjsoncommon.JsonBoolean != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "bool")
					}
					if jsonTypes&ndjsoncommon.JsonNumber != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "int", "float")
					}
					if jsonTypes&ndjsoncommon.JsonString != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "str")
					}
					if jsonTypes&ndjsoncommon.JsonArray != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "list")
					}
					if jsonTypes&ndjsoncommon.JsonObject != 0 {
						jsonTypeStrings = append(jsonTypeStrings, "dict")
					}

					if jsonTypes&possibleTypes != 0 {
						simplfied = "False"
					}
					possibleTypes |= jsonTypes
					options[i] = fmt.Sprintf("(%s.%s, %s, [%s])", classSyntax, formatting.ToPascalCase(c.Tag), typeConverter(c.Type, contextNamespace, namedType), strings.Join(jsonTypeStrings, ", "))
				}
			}

			return fmt.Sprintf("_ndjson.UnionConverter(%s, [%s], %s)", unionClassName, strings.Join(options, ", "), simplfied)
		}
		switch td := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return getScalarConverter()
		case *dsl.Vector:
			if td.Length != nil {
				return fmt.Sprintf("_ndjson.FixedVectorConverter(%s, %d)", getScalarConverter(), *td.Length)
			}

			return fmt.Sprintf("_ndjson.VectorConverter(%s)", getScalarConverter())
		case *dsl.Array:
			if td.IsFixed() {
				dims := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					dims[i] = strconv.FormatUint(*d.Length, 10)
				}

				return fmt.Sprintf("_ndjson.FixedNDArrayConverter(%s, (%s,))", getScalarConverter(), strings.Join(dims, ", "))
			}

			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("_ndjson.NDArrayConverter(%s, %d)", getScalarConverter(), len(*td.Dimensions))
			}

			return fmt.Sprintf("_ndjson.DynamicNDArrayConverter(%s)", getScalarConverter())

		case *dsl.Map:
			keyConverter := typeConverter(td.KeyType, contextNamespace, namedType)
			valueConverter := typeConverter(t.ToScalar(), contextNamespace, namedType)

			return fmt.Sprintf("_ndjson.MapConverter(%s, %s)", keyConverter, valueConverter)
		default:
			panic(fmt.Sprintf("Not implemented %T", t.Dimensionality))
		}
	default:
		panic(fmt.Sprintf("Not implemented %T", t))
	}
}

func NDJsonWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("NDJson%sWriter", formatting.ToPascalCase(p.Name))
}

func NDJsonReaderName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("NDJson%sReader", formatting.ToPascalCase(p.Name))
}
