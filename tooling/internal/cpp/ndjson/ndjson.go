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
		fmt.Fprintf(w, "namespace %s::ndjson {\n", common.NamespaceIdentifierName(ns.Name))
		// writeNamespaceDefinitions(w, ns)
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
				common.WriteComment(w, "The stream_arg parameter can either be a std::string filename")
				common.WriteComment(w, "or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::ostream.")
				w.WriteStringln("template <typename TStreamArg>")
				fmt.Fprintf(w, "%s(TStreamArg&& stream_arg)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonWriter(std::forward<TStreamArg>(stream_arg), schema_) {")
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
				common.WriteComment(w, "The stream_arg parameter can either be a std::string filename")
				common.WriteComment(w, "or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::istream.")
				w.WriteStringln("template <typename TStreamArg>")
				fmt.Fprintf(w, "%s(TStreamArg&& stream_arg)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::ndjson::NDJsonReader(std::forward<TStreamArg>(stream_arg), schema_) {")
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
