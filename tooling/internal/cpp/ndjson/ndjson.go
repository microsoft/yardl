package ndjson

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

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

	w.WriteStringln(`#include "protocols.h"
`)
	for _, ns := range env.Namespaces {
		w.WriteString("using json = nlohmann::ordered_json;\n\n")

		if len(ns.TypeDefinitions) > 0 {
			fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
			for _, typeDef := range ns.TypeDefinitions {
				writeConverters(w, typeDef)
			}
			fmt.Fprintf(w, "} // namespace %s\n\n", common.NamespaceIdentifierName(ns.Name))
		}

		fmt.Fprintf(w, "namespace %s::ndjson {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {
			writeProtocolMethods(w, protocol)
		}
		fmt.Fprintf(w, "} // namespace %s::ndjson", common.NamespaceIdentifierName(ns.Name))
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.cc")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)

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
			common.WriteComment(w, fmt.Sprintf("Json writer for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			writerClassName := JsonWriterClassName(protocol)
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

			common.WriteComment(w, fmt.Sprintf("Json reader for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			readerClassName := JsonReaderClassName(protocol)
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
				hasStream := false
				for _, step := range protocol.Sequence {
					if step.IsStream() {
						hasStream = true
					}

					returnType := "void"
					if step.IsStream() {
						returnType = "bool"
					}
					fmt.Fprintf(w, "%s %s(%s& value) override;\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
				}

				w.WriteString("void CloseImpl() override;\n")
				if hasStream {
					w.WriteStringln("\nprivate:")
					w.WriteStringln("size_t current_block_remaining_ = 0;")
				}
			})
			fmt.Fprint(w, "};\n\n")
		}
		w.WriteStringln("}")
	}

	filePath := path.Join(options.SourcesOutputDir, "protocols.h")
	return iocommon.WriteFileIfNeeded(filePath, b.Bytes(), 0644)
}

func writeConverters(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	typeName := common.TypeDefinitionSyntax(t)
	fmt.Fprintf(w, "void to_json(json& j, %s const& value) {\n", typeName)
	w.Indented(func() {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			w.WriteStringln("j = json{")
			w.Indented(func() {
				for _, field := range t.Fields {
					fmt.Fprintf(w, "{\"%s\", value.%s},\n", field.Name, common.FieldIdentifierName(field.Name))
				}
			})
			w.WriteStringln("};")
		case *dsl.EnumDefinition:
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
		default:
			panic(fmt.Sprintf("Unsupported type: %T", t))
		}
	})
	w.WriteStringln("}\n")

	fmt.Fprintf(w, "void from_json(json const& j, %s& value) {\n", common.TypeDefinitionSyntax(t))
	w.Indented(func() {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			for _, field := range t.Fields {
				fmt.Fprintf(w, "j.at(\"%s\").get_to(value.%s);\n", field.Name, common.FieldIdentifierName(field.Name))
			}
		case *dsl.EnumDefinition:
			w.WriteStringln("if (j.is_string()) {")
			w.Indented(func() {
				w.WriteStringln("std::string_view symbol = j.get<std::string_view>();")
				for _, v := range t.Values {
					fmt.Fprintf(w, "if (symbol == \"%s\") {\n", v.Symbol)
					w.Indented(func() {
						fmt.Fprintf(w, "value = %s::%s;\n", typeName, common.EnumValueIdentifierName(v.Symbol))
						w.WriteStringln("return;")
					})
					w.WriteStringln("}")
				}
				fmt.Fprintf(w, "throw std::runtime_error(\"Invalid enum value '\" + std::string(symbol) + \"' for enum %s\");\n", typeName)
			})
			w.WriteStringln("}")
			fmt.Fprintf(w, "using underlying_type = typename std::underlying_type<%s>::type;\n", typeName)
			fmt.Fprintf(w, "value = static_cast<%s>(j.get<underlying_type>());\n", typeName)

		default:
			panic(fmt.Sprintf("Unsupported type: %T", t))
		}
	})
	w.WriteStringln("}\n")
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writerClassName := JsonWriterClassName(p)

	for _, step := range p.Sequence {
		fmt.Fprintf(w, "void %s::%s([[maybe_unused]]%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			w.WriteStringln("json json_value = value;")
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

	readerClassName := JsonReaderClassName(p)
	for _, step := range p.Sequence {
		returnType := "void"
		if step.IsStream() {
			returnType = "bool"
		}

		fmt.Fprintf(w, "%s %s::%s([[maybe_unused]]%s& value) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
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

func JsonWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriter", p.Name)
}

func QualifiedJsonWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::ndjson::%s", common.TypeNamespaceIdentifierName(p), JsonWriterClassName(p))
}

func JsonReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReader", p.Name)
}

func QualifiedJsonReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::ndjson::%s", common.TypeNamespaceIdentifierName(p), JsonReaderClassName(p))
}
