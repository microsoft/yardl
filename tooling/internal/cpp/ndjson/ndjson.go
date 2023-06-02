package ndjson

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

type jsonDataType int

const (
	jsonNull jsonDataType = 1 << iota
	jsonBoolean
	jsonNumber
	jsonString
	jsonArray
	jsonObject
)

func (t jsonDataType) typeCheck(varName string) string {
	options := make([]string, 0, 1)
	if t&jsonNull != 0 {
		options = append(options, fmt.Sprintf("%s.is_null()", varName))
	}
	if t&jsonBoolean != 0 {
		options = append(options, fmt.Sprintf("%s.is_boolean()", varName))
	}
	if t&jsonNumber != 0 {
		options = append(options, fmt.Sprintf("%s.is_number()", varName))
	}
	if t&jsonString != 0 {
		options = append(options, fmt.Sprintf("%s.is_string()", varName))
	}
	if t&jsonArray != 0 {
		options = append(options, fmt.Sprintf("%s.is_array()", varName))
	}
	if t&jsonObject != 0 {
		options = append(options, fmt.Sprintf("%s.is_object()", varName))
	}
	return fmt.Sprintf("(%s)", strings.Join(options, " || "))
}

func WriteNdJson(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	options = options.ChangeOutputDir("ndjson")
	if err := os.MkdirAll(options.SourcesOutputDir, 0775); err != nil {
		return err
	}

	err := writeHeaderFile(env, options)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#include "../yardl/detail/ndjson/serializers.h"
#include "protocols.h"
`)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
		w.WriteString("using ordered_json = nlohmann::ordered_json;\n\n")

		for _, t := range ns.TypeDefinitions {
			switch t := t.(type) {
			case *dsl.EnumDefinition, *dsl.RecordDefinition:
				typeName := common.TypeDefinitionSyntax(t)
				typeParameters := t.GetDefinitionMeta().TypeParameters
				var templateDeclarationBuilder strings.Builder
				if len(typeParameters) > 0 {
					templateDeclarationBuilder.WriteString("template <")
					for i, tp := range typeParameters {
						if i > 0 {
							templateDeclarationBuilder.WriteString(", ")
						}
						templateDeclarationBuilder.WriteString(fmt.Sprintf("typename %s", tp.Name))
					}
					templateDeclarationBuilder.WriteString(">\n")
				}
				w.WriteString(templateDeclarationBuilder.String())
				fmt.Fprintf(w, "void to_json(ordered_json& j, %s const& value);\n", typeName)
				w.WriteString(templateDeclarationBuilder.String())
				fmt.Fprintf(w, "void from_json(ordered_json const& j, %s& value);\n\n", typeName)
			}
		}

		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
	}

	unionsBySyntax := make(map[string]*dsl.GeneralizedType)
	for _, ns := range env.Namespaces {
		for _, p := range ns.Protocols {
			dsl.Visit(p, func(self dsl.Visitor, node dsl.Node) {
				switch t := node.(type) {
				case *dsl.SimpleType:
					self.Visit(t.ResolvedDefinition)
				case *dsl.GeneralizedType:
					if t.Cases.IsUnion() {
						scalarType := t.ToScalar().(*dsl.GeneralizedType)
						typeSyntax := common.TypeSyntax(scalarType)
						if _, ok := unionsBySyntax[typeSyntax]; !ok {
							if len(unionsBySyntax) == 0 {
								w.WriteStringln("NLOHMANN_JSON_NAMESPACE_BEGIN\n")
							}
							unionsBySyntax[typeSyntax] = t
							writeUnionConverters(w, scalarType)
						}
					}
				}

				self.VisitChildren(node)
			})
		}
	}

	if len(unionsBySyntax) > 0 {
		w.WriteStringln("NLOHMANN_JSON_NAMESPACE_END\n")
	}

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))

		w.WriteString("using ordered_json = nlohmann::ordered_json;\n\n")

		for _, t := range ns.TypeDefinitions {
			switch t := t.(type) {
			case *dsl.EnumDefinition:
				writeEnumValuesMap(w, t)
				if t.IsFlags {
					writeFlagsConverters(w, t)
				} else {
					writeEnumConverters(w, t)
				}
			case *dsl.RecordDefinition:
				writeRecordConverters(w, t)
			}
		}

		fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))

		fmt.Fprintf(w, "namespace %s::ndjson {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {
			writeProtocolMethods(w, protocol)
		}
		fmt.Fprintf(w, "} // namespace %s::ndjson", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.cc")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)

}

func writeEnumValuesMap(w *formatting.IndentedWriter, t *dsl.EnumDefinition) {
	w.WriteStringln("namespace {")
	fmt.Fprintf(w, "std::unordered_map<std::string, %s> const %s = {\n", common.TypeDefinitionSyntax(t), enumValuesMapName(t))
	for _, v := range t.Values {
		fmt.Fprintf(w, "  {\"%s\", %s::%s},\n", v.Symbol, common.TypeDefinitionSyntax(t), common.EnumValueIdentifierName(v.Symbol))
	}
	w.WriteStringln("};")

	w.WriteStringln("} //namespace\n")
}

func enumValuesMapName(enum *dsl.EnumDefinition) string {
	return fmt.Sprintf("__%s_values", common.TypeIdentifierName(enum.Name))
}

func writeHeaderFile(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#pragma once
#include <array>
#include <complex>
#include <memory>
#include <optional>
#include <variant>
#include <vector>

#include "../yardl/detail/ndjson/reader_writer.h"
#include "../protocols.h"
#include "../types.h"
`)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s::ndjson {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {
			common.WriteComment(w, fmt.Sprintf("NDJSON writer for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			writerClassName := NDJsonWriterClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::ndjson::NDJsonWriter {\n", writerClassName, common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				w.WriteStringln("public:")
				fmt.Fprintf(w, "%s(std::ostream& stream)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonWriter(stream, schema_) {")
					})
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "%s(std::string file_name)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonWriter(file_name, schema_) {")
					})
				})
				w.WriteStringln("}\n")

				w.WriteString("void Flush() override;\n\n")

				w.WriteStringln("protected:")
				for _, step := range protocol.Sequence {
					endMethodName := common.ProtocolWriteEndImplMethodName(step)
					common.WriteComment(w, step.Comment)

					fmt.Fprintf(w, "void %s(%s const& value) override;\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))

					if step.IsStream() {
						fmt.Fprintf(w, "void %s() override {}\n", endMethodName)
					}
				}

				w.WriteString("void CloseImpl() override;\n")
			})
			fmt.Fprint(w, "};\n\n")

			common.WriteComment(w, fmt.Sprintf("NDJSON reader for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			readerClassName := NDJsonReaderClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::ndjson::NDJsonReader {\n", readerClassName, common.QualifiedAbstractReaderName(protocol))
			w.Indented(func() {
				fmt.Fprintln(w, "public:")
				fmt.Fprintf(w, "%s(std::istream& stream)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonReader(stream, schema_) {")
					})
				})
				w.WriteStringln("}\n")

				fmt.Fprintf(w, "%s(std::string file_name)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonReader(file_name, schema_) {")
					})
				})
				w.WriteStringln("}\n")

				w.WriteStringln("protected:")
				for _, step := range protocol.Sequence {
					returnType := "void"
					if step.IsStream() {
						returnType = "bool"
					}
					fmt.Fprintf(w, "%s %s(%s& value) override;\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
				}

				w.WriteString("void CloseImpl() override;\n")
			})
			fmt.Fprint(w, "};\n\n")
		}
		w.WriteStringln("}")
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.h")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)
}

func writeRecordConverters(w *formatting.IndentedWriter, t *dsl.RecordDefinition) {
	typeName := common.TypeDefinitionSyntax(t)
	typeParameters := t.GetDefinitionMeta().TypeParameters
	var templateDeclarationBuilder strings.Builder
	if len(typeParameters) > 0 {
		templateDeclarationBuilder.WriteString("template <")
		for i, tp := range typeParameters {
			if i > 0 {
				templateDeclarationBuilder.WriteString(", ")
			}
			templateDeclarationBuilder.WriteString(fmt.Sprintf("typename %s", tp.Name))
		}
		templateDeclarationBuilder.WriteString(">\n")
	}

	w.WriteString(templateDeclarationBuilder.String())
	fmt.Fprintf(w, "void to_json(ordered_json& j, %s const& value) {\n", typeName)
	w.Indented(func() {
		w.WriteStringln("j = ordered_json::object();")
		for _, field := range t.Fields {
			fmt.Fprintf(w, "if (yardl::ndjson::ShouldSerializeFieldValue(value.%s)) {\n", common.FieldIdentifierName(field.Name))
			w.Indented(func() {
				fmt.Fprintf(w, "j.push_back({\"%s\", value.%s});\n", field.Name, common.FieldIdentifierName(field.Name))
			})
			w.WriteStringln("}")
		}
	})
	w.WriteStringln("}\n")

	w.WriteString(templateDeclarationBuilder.String())
	fmt.Fprintf(w, "void from_json(ordered_json const& j, %s& value) {\n", typeName)
	w.Indented(func() {
		for _, field := range t.Fields {
			fmt.Fprintf(w, "if (auto it = j.find(\"%s\"); it != j.end()) {\n", field.Name)
			w.Indented(func() {
				fmt.Fprintf(w, "it->get_to(value.%s);\n", common.FieldIdentifierName(field.Name))
			})
			w.WriteStringln("}")
		}
	})
	w.WriteStringln("}\n")
}

func writeEnumConverters(w *formatting.IndentedWriter, t *dsl.EnumDefinition) {
	typeName := common.TypeDefinitionSyntax(t)
	fmt.Fprintf(w, "void to_json(ordered_json& j, %s const& value) {\n", typeName)
	w.Indented(func() {
		w.WriteStringln("switch (value) {")
		w.Indented(func() {
			for _, v := range t.Values {
				fmt.Fprintf(w, "case %s::%s:\n", typeName, common.EnumValueIdentifierName(v.Symbol))
				w.Indented(func() {
					fmt.Fprintf(w, "j = \"%s\";\n", v.Symbol)
					w.WriteStringln("break;")
				})
			}
			w.WriteStringln("default:")
			w.Indented(func() {
				fmt.Fprintf(w, "using underlying_type = typename std::underlying_type<%s>::type;\n", typeName)
				w.WriteStringln("j = static_cast<underlying_type>(value);")
				w.WriteStringln("break;")
			})
		})
		w.WriteStringln("}")
	})
	w.WriteStringln("}\n")

	fmt.Fprintf(w, "void from_json(ordered_json const& j, %s& value) {\n", common.TypeDefinitionSyntax(t))
	w.Indented(func() {
		w.WriteStringln("if (j.is_string()) {")
		w.Indented(func() {
			w.WriteStringln("auto symbol = j.get<std::string>();")
			fmt.Fprintf(w, "if (auto res = %s.find(symbol); res != %s.end()) {\n", enumValuesMapName(t), enumValuesMapName(t))
			w.Indented(func() {
				fmt.Fprintf(w, "value = res->second;\n")
				w.WriteStringln("return;")
			})
			w.WriteStringln("}")
			fmt.Fprintf(w, "throw std::runtime_error(\"Invalid enum value '\" + symbol + \"' for enum %s\");\n", typeName)
		})
		w.WriteStringln("}")
		fmt.Fprintf(w, "using underlying_type = typename std::underlying_type<%s>::type;\n", typeName)
		fmt.Fprintf(w, "value = static_cast<%s>(j.get<underlying_type>());\n", typeName)
	})
	w.WriteStringln("}\n")
}

func writeFlagsConverters(w *formatting.IndentedWriter, t *dsl.EnumDefinition) {
	// If the value is not a combination of the defined flags, we write it as the integer value.
	// Otherwise, we write it as an array of strings.
	// If the value is zero and there is a zero value defined, we write it as the zero value,
	// otherwise, the array will be empty.

	typeName := common.TypeDefinitionSyntax(t)
	zero := t.GetZeroValue()
	fmt.Fprintf(w, "void to_json(ordered_json& j, %s const& value) {\n", typeName)
	w.Indented(func() {
		w.WriteStringln("auto arr = ordered_json::array();")
		w.WriteStringln("if (value == 0) {")
		w.Indented(func() {
			if zero != nil {
				fmt.Fprintf(w, "arr.push_back(\"%s\");\n", zero.Symbol)
			}
			w.WriteStringln("j = arr;")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}")

		w.WriteStringln("auto remaining = value;")
		for _, v := range t.Values {
			if v == zero {
				continue
			}
			fmt.Fprintf(w, "if (remaining.HasFlags(%s::%s)) {\n", typeName, common.EnumValueIdentifierName(v.Symbol))
			w.Indented(func() {
				fmt.Fprintf(w, "remaining.UnsetFlags(%s::%s);\n", typeName, common.EnumValueIdentifierName(v.Symbol))
				fmt.Fprintf(w, "arr.push_back(\"%s\");\n", v.Symbol)
				w.WriteStringln("if (remaining == 0) {")
				w.Indented(func() {
					w.WriteStringln("j = arr;")
					w.WriteStringln("return;")
				})
				w.WriteStringln("}")
			})
			w.WriteStringln("}")
		}

		w.WriteStringln("j = value.Value();")
	})
	w.WriteStringln("}\n")

	fmt.Fprintf(w, "void from_json(ordered_json const& j, %s& value) {\n", common.TypeDefinitionSyntax(t))
	w.Indented(func() {
		w.WriteStringln("if (j.is_number()) {")
		w.Indented(func() {
			fmt.Fprintf(w, "using underlying_type = typename %s::value_type;\n", typeName)
			w.WriteStringln("value = j.get<underlying_type>();")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}")
		w.WriteStringln("std::vector<std::string> arr = j;")
		w.WriteStringln("value = {};")
		w.WriteStringln("for (auto const& item : arr) {")
		w.Indented(func() {
			fmt.Fprintf(w, "if (auto res = %s.find(item); res != %s.end()) {\n", enumValuesMapName(t), enumValuesMapName(t))
			w.Indented(func() {
				fmt.Fprintf(w, "value |= res->second;\n")
				w.WriteStringln("continue;")
			})
			w.WriteStringln("}")
			fmt.Fprintf(w, "throw std::runtime_error(\"Invalid enum value '\" + item + \"' for enum %s\");\n", typeName)
		})
		w.WriteStringln("}")
	})
	w.WriteStringln("}\n")
}

func writeUnionConverters(w *formatting.IndentedWriter, unionType *dsl.GeneralizedType) {
	simplfied := true
	var possibleTypes jsonDataType
	for _, c := range unionType.Cases {
		thisType := getJsonDataType(c.Type)
		if thisType&possibleTypes != 0 {
			simplfied = false
		}
		possibleTypes |= thisType
	}

	unionTypeSyntax := common.TypeSyntax(unionType)

	w.WriteStringln("template<>")
	fmt.Fprintf(w, "struct adl_serializer<%s> {\n", unionTypeSyntax)
	w.Indented(func() {

		fmt.Fprintf(w, "static void to_json(ordered_json& j, %s const& value) {\n", unionTypeSyntax)
		w.Indented(func() {
			if simplfied {
				w.WriteStringln("std::visit([&j](auto const& v) {j = v;}, value);")
			} else {
				w.WriteStringln("switch (value.index()) {")
				w.Indented(func() {
					for i, c := range unionType.Cases {
						fmt.Fprintf(w, "case %d:\n", i)
						w.Indented(func() {
							fmt.Fprintf(w, "j = ordered_json{ {\"%s\", std::get<%s>(value)} };\n", c.Label, common.TypeSyntax(c.Type))
							w.WriteStringln("break;")
						})
					}
					w.WriteStringln("default:")
					w.Indented(func() {
						w.WriteStringln("throw std::runtime_error(\"Invalid union value\");")
					})
				})
				w.WriteStringln("}")
			}

		})
		w.WriteStringln("}\n")

		fmt.Fprintf(w, "static void from_json(ordered_json const& j, %s& value) {\n", unionTypeSyntax)
		w.Indented(func() {
			if simplfied {
				for _, c := range unionType.Cases {
					dt := getJsonDataType(c.Type)
					fmt.Fprintf(w, "if (%s) {\n", dt.typeCheck("j"))
					w.Indented(func() {
						fmt.Fprintf(w, "value = j.get<%s>();\n", common.TypeSyntax(c.Type))
						w.WriteStringln("return;")
					})
					w.WriteStringln("}")
				}

				w.WriteStringln("throw std::runtime_error(\"Invalid union value\");")
			} else {
				w.WriteStringln("auto it = j.begin();")
				w.WriteStringln("std::string label = it.key();")
				for _, v := range unionType.Cases {
					fmt.Fprintf(w, "if (label == \"%s\") {\n", v.Label)
					w.Indented(func() {
						fmt.Fprintf(w, "value = it.value().get<%s>();\n", common.TypeSyntax(v.Type))
						w.WriteStringln("return;")
					})
					w.WriteStringln("}")
				}
			}
		})
		w.WriteStringln("}")
	})
	w.WriteStringln("};\n")

}

func getJsonDataType(t dsl.Type) jsonDataType {
	if t == nil {
		return jsonNull
	}
	gt := dsl.ToGeneralizedType(t)
	switch d := gt.Dimensionality.(type) {
	case *dsl.Vector:
		return jsonArray
	case *dsl.Array:
		if d.IsFixed() {
			return jsonArray
		}
		return jsonObject
	case *dsl.Map:
		if p, ok := dsl.GetPrimitiveType(d.KeyType); ok && p == dsl.String {
			return jsonObject
		}
		return jsonArray
	}

	if len(gt.Cases) > 1 {
		panic("unexpected union type")
	}

	scalarType := gt.Cases[0].Type.(*dsl.SimpleType)
	switch td := scalarType.ResolvedDefinition.(type) {
	case dsl.PrimitiveDefinition:
		switch td {
		case dsl.String:
			return jsonString
		case dsl.Int8, dsl.Int16, dsl.Int32, dsl.Int64, dsl.Uint8, dsl.Uint16, dsl.Uint32, dsl.Uint64, dsl.Size, dsl.Float32, dsl.Float64:
			return jsonNumber
		case dsl.Bool:
			return jsonBoolean
		case dsl.ComplexFloat32, dsl.ComplexFloat64:
			return jsonArray
		case dsl.Date, dsl.Time, dsl.DateTime:
			return jsonNumber
		default:
			panic(fmt.Sprintf("unexpected primitive type %s", td))
		}
	case *dsl.EnumDefinition:
		if td.IsFlags {
			return jsonArray
		}
		return jsonString | jsonNumber
	case *dsl.RecordDefinition:
		return jsonObject
	case *dsl.GenericTypeParameter:
		return jsonObject
	case *dsl.NamedType:
		return getJsonDataType(td.Type)
	default:
		panic(fmt.Sprintf("unexpected type %T", td))
	}
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writerClassName := NDJsonWriterClassName(p)

	for _, step := range p.Sequence {
		fmt.Fprintf(w, "void %s::%s(%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			w.WriteStringln("ordered_json json_value = value;")
			fmt.Fprintf(w, "yardl::ndjson::WriteProtocolValue(stream_, \"%s\", json_value);", step.Name)
		})
		w.WriteString("}\n\n")
	}

	fmt.Fprintf(w, "void %s::Flush() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.flush();\n")
	})
	w.WriteString("}\n\n")

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.flush();\n")
	})
	w.WriteString("}\n\n")

	readerClassName := NDJsonReaderClassName(p)
	for _, step := range p.Sequence {
		returnType := "void"
		if step.IsStream() {
			returnType = "bool"
		}

		fmt.Fprintf(w, "%s %s::%s(%s& value) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			if step.IsStream() {
				w.WriteString("return ")
			}
			fmt.Fprintf(w, "yardl::ndjson::ReadProtocolValue(stream_, line_, \"%s\", %t, unused_step_, value);\n", step.Name, !step.IsStream())

		})
		w.WriteString("}\n\n")

	}

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", readerClassName)
	w.Indented(func() {
		w.WriteString("VerifyFinished();\n")
	})
	w.WriteString("}\n\n")
}

func NDJsonWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriter", p.Name)
}

func QualifiedNDJsonWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::ndjson::%s", common.TypeNamespaceIdentifierName(p), NDJsonWriterClassName(p))
}

func NDJsonReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReader", p.Name)
}

func QualifiedNDJsonReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::ndjson::%s", common.TypeNamespaceIdentifierName(p), NDJsonReaderClassName(p))
}
