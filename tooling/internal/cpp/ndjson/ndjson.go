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
	// writeIsTriviallySerializableSpecializations(w, env)
	// writeUnionSerializers(w, env)
	for _, ns := range env.Namespaces {
		w.WriteStringln("using json = nlohmann::ordered_json;")

		fmt.Fprintf(w, "namespace %s::ndjson {\n", common.NamespaceIdentifierName(ns.Name))
		writeNamespaceDefinitions(w, ns)
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
						fmt.Fprintf(w, "void %s(std::vector<%s> const& values) override;\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
						fmt.Fprintf(w, "void %s() override;\n", endMethodName)
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
					if step.IsStream() {
						fmt.Fprintf(w, "bool %s(std::vector<%s>& values) override;\n", common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
					}
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

func writeNamespaceDefinitions(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	if len(ns.TypeDefinitions) > 0 {
		w.WriteStringln("namespace {")
		// for _, typeDef := range ns.TypeDefinitions {
		// 	writeSerializers(w, typeDef)
		// }
		w.WriteString("} // namespace\n\n")
	}

	for _, protocol := range ns.Protocols {
		writeProtocolMethods(w, protocol)
	}
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writerClassName := JsonWriterClassName(p)

	for _, step := range p.Sequence {
		fmt.Fprintf(w, "void %s::%s([[maybe_unused]]%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			w.WriteStringln("json json_value = value;")
			fmt.Fprintf(w, "yardl::ndjson::WriteProtocolValue(stream_, \"%s\", json_value);", step.Name)

			// if step.IsStream() {
			// 	w.WriteString("yardl::binary::WriteInteger(stream_, 1U);\n")
			// }
			// fmt.Fprintf(w, "%s(stream_, value);\n", typeRwFunction(step.Type, true))
		})
		w.WriteString("}\n\n")

		// if step.IsStream() {
		// 	fmt.Fprintf(w, "void %s::%s([[maybe_unused]]std::vector<%s> const& values) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		// 	// w.Indented(func() {
		// 	// 	w.WriteStringln("if (!values.empty()) {")
		// 	// 	w.Indented(func() {
		// 	// 		vectorType := *step.Type.(*dsl.GeneralizedType)
		// 	// 		vectorType.Dimensionality = &dsl.Vector{}
		// 	// 		fmt.Fprintf(w, "%s(stream_, values);\n", typeRwFunction(&vectorType, true))
		// 	// 	})
		// 	// 	w.WriteStringln("}")
		// 	// })
		// 	w.WriteString("}\n\n")

		// 	fmt.Fprintf(w, "void %s::%s() {\n", writerClassName, common.ProtocolWriteEndImplMethodName(step))
		// 	// w.Indented(func() {
		// 	// 	w.WriteString("yardl::binary::WriteInteger(stream_, 0U);\n")
		// 	// })
		// 	w.WriteString("}\n\n")
		// }
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
			fmt.Fprintf(w, "yardl::ndjson::ReadProtocolValue(stream_, line_, \"%s\", %t, unused_step_, value);\n", step.Name, step.IsStream())

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
