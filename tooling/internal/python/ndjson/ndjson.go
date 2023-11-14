// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

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

	fmt.Fprintf(w, "from %s import _ndjson\n", relativePath)
	fmt.Fprintf(w, "from %s import yardl_types as yardl\n\n", relativePath)

	writeConverters(w, ns)
	if ns.IsTopLevel {
		writeProtocols(w, ns)
	}

	ndjsonPath := path.Join(packageDir, "ndjson.py")
	return iocommon.WriteFileIfNeeded(ndjsonPath, b.Bytes(), 0644)
}

func writeConverters(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, t := range ns.TypeDefinitions {
		switch t := t.(type) {
		case *dsl.EnumDefinition:
			writeEnumMaps(t, w, ns)
		case *dsl.RecordDefinition:
			writeRecordConverter(t, w, ns)
		}

	}
}

func writeRecordConverter(td *dsl.RecordDefinition, w *formatting.IndentedWriter, ns *dsl.Namespace) {
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

	fmt.Fprintf(w, "class %s(%s_ndjson.JsonConverter[%s, np.void]):\n", recordConverterClassName(td, ns.Name), genericSpec, typeSyntax)
	w.Indented(func() {
		if len(td.TypeParameters) > 0 {
			typeParamSerializers := make([]string, 0, len(td.TypeParameters))
			for _, tp := range td.TypeParameters {
				typeParamSerializers = append(
					typeParamSerializers,
					fmt.Sprintf("%s: _ndjson.JsonConverter[%s, %s]", typeDefinitionConverter(tp, ns.Name), common.TypeParameterSyntax(tp, false), common.TypeParameterSyntax(tp, true)))
			}

			fmt.Fprintf(w, "def __init__(self, %s) -> None:\n", strings.Join(typeParamSerializers, ", "))
		} else {
			w.WriteStringln("def __init__(self) -> None:")
		}
		w.Indented(func() {
			for _, f := range td.Fields {
				fieldId := common.FieldIdentifierName(f.Name)
				fmt.Fprintf(w, "self._%s_converter = %s\n", fieldId, typeConverter(f.Type, ns.Name, nil))
				if isGenericParameterReference(f.Type) {
					fmt.Fprintf(w, "self._%s_supports_none = self._%s_converter.supports_none()\n", fieldId, fieldId)
				}
			}

			fmt.Fprintf(w, "super().__init__(np.dtype([\n")
			w.Indented(func() {
				for _, f := range td.Fields {
					fmt.Fprintf(w, "(\"%s\", self._%s_converter.overall_dtype()),\n", common.FieldIdentifierName(f.Name), common.FieldIdentifierName(f.Name))
				}
			})
			fmt.Fprintf(w, "]))\n")

		})
		w.WriteStringln("")

		fmt.Fprintf(w, "def to_json(self, value: %s) -> object:\n", typeSyntax)
		w.Indented(func() {
			fmt.Fprintf(w, "if not isinstance(value, %s): # pyright: ignore [reportUnnecessaryIsInstance]\n", common.TypeSyntaxWithoutTypeParameters(td, ns.Name))
			w.Indented(func() {
				fmt.Fprintf(w, "raise TypeError(\"Expected '%s' instance\")\n", typeSyntax)
			})
			w.WriteStringln("json_object = {}\n")

			for _, f := range td.Fields {
				fieldId := common.FieldIdentifierName(f.Name)
				if g, ok := f.Type.(*dsl.GeneralizedType); ok && g.Cases.HasNullOption() && g.Dimensionality == nil {
					fmt.Fprintf(w, "if value.%s is not None:\n", fieldId)
					w.Indented(func() {
						fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.to_json(value.%s)\n", f.Name, fieldId, fieldId)
					})
				} else if isGenericParameterReference(f.Type) {
					fmt.Fprintf(w, "if not self._%s_supports_none or value.%s is not None:\n", fieldId, fieldId)
					w.Indented(func() {
						fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.to_json(value.%s)\n", f.Name, fieldId, fieldId)
					})
				} else {
					fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.to_json(value.%s)\n", f.Name, fieldId, fieldId)
				}
			}
			fmt.Fprintf(w, "return json_object\n")
		})
		w.WriteStringln("")

		fmt.Fprintf(w, "def numpy_to_json(self, value: np.void) -> object:\n")
		w.Indented(func() {
			fmt.Fprintf(w, "if not isinstance(value, np.void): # pyright: ignore [reportUnnecessaryIsInstance]\n")
			w.Indented(func() {
				fmt.Fprintf(w, "raise TypeError(\"Expected 'np.void' instance\")\n")
			})
			w.WriteStringln("json_object = {}\n")

			for _, f := range td.Fields {
				fieldId := common.FieldIdentifierName(f.Name)
				if g, ok := f.Type.(*dsl.GeneralizedType); ok && g.Cases.HasNullOption() && g.Dimensionality == nil {
					fmt.Fprintf(w, "if (field_val := value[\"%s\"]) is not None:\n", fieldId)
					w.Indented(func() {
						fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.numpy_to_json(field_val)\n", f.Name, fieldId)
					})
				} else if isGenericParameterReference(f.Type) {
					fmt.Fprintf(w, "if not self._%s_supports_none or value[\"%s\"] is not None:\n", fieldId, fieldId)
					w.Indented(func() {
						fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.numpy_to_json(value[\"%s\"])\n", f.Name, fieldId, fieldId)
					})
				} else {
					fmt.Fprintf(w, "json_object[\"%s\"] = self._%s_converter.numpy_to_json(value[\"%s\"])\n", f.Name, fieldId, fieldId)
				}
			}
			fmt.Fprintf(w, "return json_object\n")
		})
		w.WriteStringln("")

		fmt.Fprintf(w, "def from_json(self, json_object: object) -> %s:\n", typeSyntax)
		w.Indented(func() {
			fmt.Fprintf(w, "if not isinstance(json_object, dict):\n")
			w.Indented(func() {
				fmt.Fprintf(w, "raise TypeError(\"Expected 'dict' instance\")\n")
			})
			fmt.Fprintf(w, "return %s(\n", typeSyntax)
			w.Indented(func() {
				for _, f := range td.Fields {
					fieldId := common.FieldIdentifierName(f.Name)
					if g, ok := f.Type.(*dsl.GeneralizedType); ok && g.Cases.HasNullOption() && g.Dimensionality == nil {
						fmt.Fprintf(w, "%s=self._%s_converter.from_json(json_object.get(\"%s\")),\n", fieldId, fieldId, f.Name)
					} else if isGenericParameterReference(f.Type) {
						fmt.Fprintf(w, "%s=self._%s_converter.from_json(json_object.get(\"%s\") if self._%s_supports_none else json_object[\"%s\"]),\n", fieldId, fieldId, f.Name, fieldId, f.Name)
					} else {
						fmt.Fprintf(w, "%s=self._%s_converter.from_json(json_object[\"%s\"],),\n", fieldId, fieldId, f.Name)
					}
				}
			})
			fmt.Fprintf(w, ")\n")
		})
		w.WriteStringln("")

		fmt.Fprintf(w, "def from_json_to_numpy(self, json_object: object) -> np.void:\n")
		w.Indented(func() {
			fmt.Fprintf(w, "if not isinstance(json_object, dict):\n")
			w.Indented(func() {
				fmt.Fprintf(w, "raise TypeError(\"Expected 'dict' instance\")\n")
			})
			fmt.Fprintf(w, "return (\n")
			w.Indented(func() {
				for _, f := range td.Fields {
					fieldId := common.FieldIdentifierName(f.Name)
					if g, ok := f.Type.(*dsl.GeneralizedType); ok && g.Cases.HasNullOption() && g.Dimensionality == nil {
						fmt.Fprintf(w, "self._%s_converter.from_json_to_numpy(json_object.get(\"%s\")),\n", fieldId, f.Name)
					} else if isGenericParameterReference(f.Type) {
						fmt.Fprintf(w, "self._%s_converter.from_json_to_numpy(json_object.get(\"%s\") if self._%s_supports_none else json_object[\"%s\"]),\n", fieldId, f.Name, fieldId, f.Name)
					} else {
						fmt.Fprintf(w, "self._%s_converter.from_json_to_numpy(json_object[\"%s\"]),\n", fieldId, f.Name)
					}
				}
			})
			fmt.Fprintf(w, ") # type:ignore \n")
		})
		w.WriteStringln("")
	})

	w.WriteStringln("")
}

func writeEnumMaps(t *dsl.EnumDefinition, w *formatting.IndentedWriter, ns *dsl.Namespace) {
	name_to_value_map_name := enumNameToValueMapName(t, ns.Name)
	fmt.Fprintf(w, "%s = {\n", name_to_value_map_name)
	w.Indented(func() {
		for _, v := range t.Values {
			fmt.Fprintf(w, "\"%s\": %s.%s,\n", v.Symbol, common.TypeSyntax(t, ns.Name), common.EnumValueIdentifierName(v.Symbol))
		}
	})
	fmt.Fprintf(w, "}\n")

	value_to_name_map_name := enumValueToNameMapName(t, ns.Name)
	fmt.Fprintf(w, "%s = {v: n for n, v in %s.items()}\n\n", value_to_name_map_name, name_to_value_map_name)
}

func enumNameToValueMapName(t *dsl.EnumDefinition, contextNamespace string) string {
	name := fmt.Sprintf("%s_name_to_value_map", formatting.ToSnakeCase(t.Name))
	if t.Namespace != contextNamespace {
		name = fmt.Sprintf("%s.ndjson.%s", common.NamespaceIdentifierName(t.Namespace), name)
	}
	return name
}

func enumValueToNameMapName(t *dsl.EnumDefinition, contextNamespace string) string {
	name := fmt.Sprintf("%s_value_to_name_map", formatting.ToSnakeCase(t.Name))
	if t.Namespace != contextNamespace {
		name = fmt.Sprintf("%s.ndjson.%s", common.NamespaceIdentifierName(t.Namespace), name)
	}
	return name
}

func recordConverterClassName(record *dsl.RecordDefinition, contextNamespace string) string {
	className := fmt.Sprintf("%sConverter", formatting.ToPascalCase(record.Name))
	if record.Namespace != contextNamespace {
		className = fmt.Sprintf("%s.ndjson.%s", common.NamespaceIdentifierName(record.Namespace), className)
	}
	return className
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
						fmt.Fprintf(w, "converter = %s\n", converter)
						fmt.Fprintf(w, "json_value = converter.to_json(value)\n")
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
						fmt.Fprintf(w, "while (json_object := self._read_json_line(\"%s\", False)) is not _ndjson.MISSING_SENTINEL:\n", step.Name)
						w.Indented(func() {
							fmt.Fprintf(w, "yield converter.from_json(json_object)\n")
						})
					})

				} else {
					fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
					w.Indented(func() {
						fmt.Fprintf(w, "json_object = self._read_json_line(\"%s\", True)\n", step.Name)
						fmt.Fprintf(w, "converter = %s\n", typeConverter(step.Type, ns.Name, nil))
						fmt.Fprintf(w, "return converter.from_json(json_object)\n")
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

		return fmt.Sprintf("%s(%s, %s, %s, %s)", className, common.TypeSyntax(t, contextNamespace), common.TypeDTypeSyntax(baseType), enumNameToValueMapName(t, contextNamespace), enumValueToNameMapName(t, contextNamespace))
	case *dsl.RecordDefinition:
		converterName := recordConverterClassName(t, contextNamespace)
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

func isGenericParameterReference(t dsl.Type) bool {
	if st, ok := t.(*dsl.SimpleType); ok {
		if _, ok := st.ResolvedDefinition.(*dsl.GenericTypeParameter); ok {
			return true
		}
	}

	return false
}
