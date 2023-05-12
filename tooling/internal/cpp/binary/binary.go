// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package binary

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteBinary(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	options = options.ChangeOutputDir("binary")
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

#include <cstddef>

#include "../yardl/detail/binary/coded_stream.h"
#include "../yardl/detail/binary/serializers.h"
`)
	writeIsTriviallySerializableSpecializations(w, env)
	writeUnionSerializers(w, env)
	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s::binary {\n", common.NamespaceIdentifierName(ns.Name))
		writeNamespaceDefinitions(w, ns)
		fmt.Fprintf(w, "} // namespace %s::binary", common.NamespaceIdentifierName(ns.Name))
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

#include "../yardl/detail/binary/reader_writer.h"
#include "../protocols.h"
#include "../types.h"
`)

	for _, ns := range env.Namespaces {
		fmt.Fprintf(w, "namespace %s::binary {\n", common.NamespaceIdentifierName(ns.Name))
		for _, protocol := range ns.Protocols {
			common.WriteComment(w, fmt.Sprintf("Binary writer for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			writerClassName := BinaryWriterClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::binary::BinaryWriter {\n", writerClassName, common.QualifiedAbstractWriterName(protocol))
			w.Indented(func() {
				w.WriteStringln("public:")
				common.WriteComment(w, "The stream_arg parameter can either be a std::string filename")
				common.WriteComment(w, "or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::ostream.")
				w.WriteStringln("template <typename TStreamArg>")
				fmt.Fprintf(w, "%s(TStreamArg&& stream_arg)\n", writerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::binary::BinaryWriter(std::forward<TStreamArg>(stream_arg), schema_) {")
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

			common.WriteComment(w, fmt.Sprintf("Binary reader for the %s protocol.", protocol.Name))
			common.WriteComment(w, protocol.Comment)
			readerClassName := BinaryReaderClassName(protocol)
			fmt.Fprintf(w, "class %s : public %s, yardl::binary::BinaryReader {\n", readerClassName, common.QualifiedAbstractReaderName(protocol))
			w.Indented(func() {
				fmt.Fprintln(w, "public:")
				common.WriteComment(w, "The stream_arg parameter can either be a std::string filename")
				common.WriteComment(w, "or a reference, std::unique_ptr, or std::shared_ptr to a stream-like object, such as std::istream.")
				w.WriteStringln("template <typename TStreamArg>")
				fmt.Fprintf(w, "%s(TStreamArg&& stream_arg)\n", readerClassName)
				w.Indented(func() {
					w.Indented(func() {
						w.WriteStringln(": yardl::binary::BinaryReader(std::forward<TStreamArg>(stream_arg), schema_) {")
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

func writeIsTriviallySerializableSpecializations(w *formatting.IndentedWriter, env *dsl.Environment) {
	fmt.Fprintf(w, "namespace yardl::binary {\n")
	w.WriteStringln("#ifndef _MSC_VER")
	common.WriteComment(w, "Values of offsetof() are only used if types are standard-layout.")
	w.WriteStringln("#pragma GCC diagnostic push")
	w.WriteStringln("#pragma GCC diagnostic ignored \"-Winvalid-offsetof\"")
	w.WriteStringln("#endif\n")

	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {
			writeIsTriviallySerializableSpecialization(w, td)
		}
	}
	w.WriteStringln("#ifndef _MSC_VER")
	w.WriteStringln("#pragma GCC diagnostic pop // #pragma GCC diagnostic ignored \"-Winvalid-offsetof\" ")
	w.WriteStringln("#endif")
	fmt.Fprintf(w, "} //namespace yardl::binary \n\n")
}

func writeIsTriviallySerializableSpecialization(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	switch t.(type) {
	case *dsl.RecordDefinition:
		break
	default:
		return
	}

	meta := t.GetDefinitionMeta()

	w.WriteString("template <")
	formatting.Delimited(w, ", ", meta.TypeParameters, func(w *formatting.IndentedWriter, i int, item *dsl.GenericTypeParameter) {
		fmt.Fprintf(w, "typename %s", common.TypeDefinitionSyntax(item))
	})
	w.WriteStringln(">")
	fmt.Fprintf(w, "struct IsTriviallySerializable<%s> {\n", common.TypeDefinitionSyntax(t))

	w.Indented(func() {
		fmt.Fprintf(w, "using __T__ = %s;\n", common.TypeDefinitionSyntax(t))
		w.WriteString("static constexpr bool value = \n")
		w.Indented(func() {
			w.WriteStringln("std::is_standard_layout_v<__T__> &&")

			switch t := t.(type) {
			case *dsl.RecordDefinition:
				formatting.Delimited(w, " &&\n", t.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
					fmt.Fprintf(w, "IsTriviallySerializable<decltype(__T__::%s)>::value", common.FieldIdentifierName(f.Name))
				})

				if len(t.Fields) > 0 {

					w.WriteStringln(" &&")

					fmt.Fprintf(w, "(sizeof(__T__) == (")
					formatting.Delimited(w, " + ", t.Fields, func(w *formatting.IndentedWriter, i int, f *dsl.Field) {
						fmt.Fprintf(w, "sizeof(__T__::%s)", common.FieldIdentifierName(f.Name))
					})

					w.WriteString("))")

					if len(t.Fields) > 1 {
						w.WriteStringln(" &&")
					}

					for i, f := range t.Fields {
						if i > 0 {
							if i > 1 {
								w.WriteString(" && ")
							}
							fmt.Fprintf(w, "offsetof(__T__, %s) < offsetof(__T__, %s)", common.FieldIdentifierName(common.FieldIdentifierName(t.Fields[i-1].Name)), common.FieldIdentifierName(f.Name))
						}
					}
				}
			}
		})

		w.WriteStringln(";")
	})
	fmt.Fprintf(w, "};\n\n")
}

func collectUnionArities(env *dsl.Environment) []int {
	arities := make(map[int]any)

	dsl.Visit(env, func(self dsl.Visitor, node dsl.Node) {
		switch t := node.(type) {
		case *dsl.GeneralizedType:
			if t.Cases.IsUnion() {
				arities[len(t.Cases)] = nil
			}

			self.VisitChildren(node)
		default:
			self.VisitChildren(node)
		}

	})

	keys := make([]int, 0, len(arities))
	for k := range arities {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func writeUnionSerializers(w *formatting.IndentedWriter, env *dsl.Environment) {
	arites := collectUnionArities(env)
	if len(arites) == 0 {
		return
	}

	w.WriteString("namespace {\n")
	formatting.Delimited(w, "\n", arites, func(w *formatting.IndentedWriter, i int, arity int) {
		elements := make([]int, arity)
		for i := 0; i < arity; i++ {
			elements[i] = i
		}

		writeTemplateParameters := make([]string, 2*arity)
		for i := range elements {
			writeTemplateParameters[2*i] = fmt.Sprintf("typename T%d", i)
			writeTemplateParameters[2*i+1] = fmt.Sprintf("yardl::binary::Writer<T%d> WriteT%d", i, i)
		}

		readTemplateParameters := make([]string, 2*arity)
		for i := range elements {
			readTemplateParameters[2*i] = fmt.Sprintf("typename T%d", i)
			readTemplateParameters[2*i+1] = fmt.Sprintf("yardl::binary::Reader<T%d> ReadT%d", i, i)
		}

		variantTemplateArguments := make([]string, arity)
		for i := range elements {
			variantTemplateArguments[i] = fmt.Sprintf("T%d", i)
		}

		fmt.Fprintf(w, "template<%s>\n", strings.Join(writeTemplateParameters, ", "))
		fmt.Fprintf(w, "void WriteUnion(yardl::binary::CodedOutputStream& stream, std::variant<%s> const& value) {\n", strings.Join(variantTemplateArguments, ", "))
		w.Indented(func() {
			w.WriteStringln("yardl::binary::WriteInteger(stream, value.index());")
			fmt.Fprintf(w, "switch (value.index()) {\n")
			for i := range elements {
				fmt.Fprintf(w, "case %d: {\n", i)
				w.Indented(func() {
					fmt.Fprintf(w, "T%d const& v = std::get<%d>(value);\n", i, i)
					fmt.Fprintf(w, "WriteT%d(stream, v);\n", i)
					w.WriteStringln("break;")
				})
				w.WriteStringln("}")
			}
			fmt.Fprintf(w, "default: throw std::runtime_error(\"Invalid union index.\");\n")
			fmt.Fprintf(w, "}\n")
		})
		fmt.Fprintf(w, "}\n\n")

		fmt.Fprintf(w, "template<%s>\n", strings.Join(readTemplateParameters, ", "))
		fmt.Fprintf(w, "void ReadUnion(yardl::binary::CodedInputStream& stream, std::variant<%s>& value) {\n", strings.Join(variantTemplateArguments, ", "))
		w.Indented(func() {
			w.WriteStringln("size_t index;")
			w.WriteStringln("yardl::binary::ReadInteger(stream, index);")
			w.WriteStringln("switch (index) {")
			w.Indented(func() {
				for i := range elements {
					fmt.Fprintf(w, "case %d: {\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "T%d v;\n", i)
						fmt.Fprintf(w, "ReadT%d(stream, v);\n", i)
						fmt.Fprintf(w, "value = std::move(v);\n")
						w.WriteStringln("break;")
					})
					w.WriteStringln("}")
				}
				fmt.Fprintf(w, "default: throw std::runtime_error(\"Invalid union index.\");\n")
			})
			w.WriteStringln("}")
		})
		fmt.Fprintf(w, "}\n")

	})

	w.WriteString("} // namespace\n\n")
}

func writeNamespaceDefinitions(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	if len(ns.TypeDefinitions) > 0 {
		w.WriteStringln("namespace {")
		for _, typeDef := range ns.TypeDefinitions {
			writeSerializers(w, typeDef)
		}
		w.WriteString("} // namespace\n\n")
	}

	for _, protocol := range ns.Protocols {
		writeProtocolMethods(w, protocol)
	}
}

func writeSerializers(w *formatting.IndentedWriter, t dsl.TypeDefinition) {
	switch t.(type) {
	case *dsl.EnumDefinition:
		return
	}

	writeFallbackBody := func(write bool) {
		switch t := t.(type) {
		case *dsl.RecordDefinition:
			for _, field := range t.Fields {
				fmt.Fprintf(w, "%s(stream, value.%s);\n", typeRwFunction(field.Type, write), common.FieldIdentifierName(field.Name))
			}
		case *dsl.NamedType:
			fmt.Fprintf(w, "%s(stream, value);\n", typeRwFunction(t.Type, write))
		default:
			panic(fmt.Sprintf("Unexpected type %T", t))
		}
	}

	writeRwFunctionSignature(t, w, true)
	w.Indented(func() {
		fmt.Fprintf(w, "if constexpr (yardl::binary::IsTriviallySerializable<%s>::value) {\n", common.TypeDefinitionSyntax(t))
		w.Indented(func() {
			fmt.Fprintf(w, "yardl::binary::WriteTriviallySerializable(stream, value);\n")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}\n")
		writeFallbackBody(true)
	})
	w.WriteString("}\n\n")

	writeRwFunctionSignature(t, w, false)
	w.Indented(func() {
		fmt.Fprintf(w, "if constexpr (yardl::binary::IsTriviallySerializable<%s>::value) {\n", common.TypeDefinitionSyntax(t))
		w.Indented(func() {
			fmt.Fprintf(w, "yardl::binary::ReadTriviallySerializable(stream, value);\n")
			w.WriteStringln("return;")
		})
		w.WriteStringln("}\n")
		writeFallbackBody(false)
	})
	w.WriteString("}\n\n")
}

func verb(write bool) string {
	if write {
		return "Write"
	}
	return "Read"
}

func writeRwFunctionSignature(t dsl.TypeDefinition, w *formatting.IndentedWriter, write bool) {
	writeRwFunctionTemplateDeclaration(t, w, write)
	if write {
		fmt.Fprintf(w, "[[maybe_unused]] void Write%s(yardl::binary::CodedOutputStream& stream, %s const& value) {\n", t.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(t))
	} else {
		fmt.Fprintf(w, "[[maybe_unused]] void Read%s(yardl::binary::CodedInputStream& stream, %s& value) {\n", t.GetDefinitionMeta().Name, common.TypeDefinitionSyntax(t))
	}
}

func writeRwFunctionTemplateDeclaration(t dsl.TypeDefinition, w *formatting.IndentedWriter, write bool) {
	meta := t.GetDefinitionMeta()
	if len(meta.TypeParameters) > 0 {
		functionPointerName := "Reader"
		if write {
			functionPointerName = "Writer"
		}
		templateParameters := make([]string, 2*len(meta.TypeParameters))
		for i, p := range meta.TypeParameters {
			templateParameters[2*i] = "typename " + common.TypeDefinitionSyntax(p)
			templateParameters[2*i+1] = fmt.Sprintf("yardl::binary::%s<%s> %s%s", functionPointerName, common.TypeDefinitionSyntax(p), verb(write), common.TypeDefinitionSyntax(p))
		}
		fmt.Fprintf(w, "template<%s>\n", strings.Join(templateParameters, ", "))
	}
}

func writeProtocolMethods(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	writerClassName := BinaryWriterClassName(p)
	for _, step := range p.Sequence {
		fmt.Fprintf(w, "void %s::%s(%s const& value) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			if step.IsStream() {
				w.WriteString("yardl::binary::WriteInteger(stream_, 1U);\n")
			}
			fmt.Fprintf(w, "%s(stream_, value);\n", typeRwFunction(step.Type, true))
		})
		w.WriteString("}\n\n")

		if step.IsStream() {
			fmt.Fprintf(w, "void %s::%s(std::vector<%s> const& values) {\n", writerClassName, common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				w.WriteStringln("if (!values.empty()) {")
				w.Indented(func() {
					vectorType := *step.Type.(*dsl.GeneralizedType)
					vectorType.Dimensionality = &dsl.Vector{}
					fmt.Fprintf(w, "%s(stream_, values);\n", typeRwFunction(&vectorType, true))
				})
				w.WriteStringln("}")
			})
			w.WriteString("}\n\n")

			fmt.Fprintf(w, "void %s::%s() {\n", writerClassName, common.ProtocolWriteEndImplMethodName(step))
			w.Indented(func() {
				w.WriteString("yardl::binary::WriteInteger(stream_, 0U);\n")
			})
			w.WriteString("}\n\n")
		}
	}

	fmt.Fprintf(w, "void %s::Flush() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.Flush();\n")
	})
	w.WriteString("}\n\n")

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", writerClassName)
	w.Indented(func() {
		w.WriteString("stream_.Flush();\n")
	})
	w.WriteString("}\n\n")

	readerClassName := BinaryReaderClassName(p)
	for _, step := range p.Sequence {
		returnType := "void"
		if step.IsStream() {
			returnType = "bool"
		}

		fmt.Fprintf(w, "%s %s::%s(%s& value) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
		w.Indented(func() {
			if step.IsStream() {
				w.WriteStringln("if (current_block_remaining_ == 0) {")
				w.Indented(func() {
					w.WriteStringln("yardl::binary::ReadInteger(stream_, current_block_remaining_);")
					w.WriteStringln("if (current_block_remaining_ == 0) {")
					w.Indented(func() {
						w.WriteStringln("return false;")
					})
					w.WriteStringln("}")
				})
				w.WriteStringln("}")
			}

			fmt.Fprintf(w, "%s(stream_, value);\n", typeRwFunction(step.Type, false))
			if step.IsStream() {
				w.WriteStringln("current_block_remaining_--;")
				w.WriteStringln("return true;")
			}
		})
		w.WriteString("}\n\n")

		if step.IsStream() {
			fmt.Fprintf(w, "%s %s::%s(std::vector<%s>& values) {\n", returnType, readerClassName, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				scalarType := step.Type.(*dsl.GeneralizedType).ToScalar()

				fmt.Fprintf(w, "yardl::binary::ReadBlocksIntoVector<%s, %s>(stream_, current_block_remaining_, values);\n", common.TypeSyntax(scalarType), typeRwFunction(scalarType, false))
				w.WriteStringln("return current_block_remaining_ != 0;")
			})
			w.WriteString("}\n\n")

		}
	}

	fmt.Fprintf(w, "void %s::CloseImpl() {\n", readerClassName)
	w.Indented(func() {
		w.WriteString("stream_.VerifyFinished();\n")
	})
	w.WriteString("}\n\n")
}

func typeDefinitionRwFunction(t dsl.TypeDefinition, write bool) string {
	switch t := t.(type) {
	case dsl.PrimitiveDefinition:
		suffix := func() string {
			switch t {
			case dsl.Bool, dsl.Int8, dsl.Int16, dsl.Int32, dsl.Int64, dsl.Uint8, dsl.Uint16, dsl.Uint32, dsl.Uint64, dsl.Size:
				return "Integer"
			case dsl.Float32, dsl.Float64, dsl.ComplexFloat32, dsl.PrimitiveComplexFloat64:
				return "FloatingPoint"
			case dsl.String:
				return "String"
			case dsl.Date:
				return "Date"
			case dsl.Time:
				return "Time"
			case dsl.DateTime:
				return "DateTime"
			default:
				panic(fmt.Sprintf("Unknown primitive type %s", t))
			}
		}()
		if write {
			return "yardl::binary::Write" + suffix
		}
		return "yardl::binary::Read" + suffix
	case *dsl.EnumDefinition:
		return fmt.Sprintf("yardl::binary::%sEnum<%s>", verb(write), common.TypeDefinitionSyntax(t))
	case *dsl.GenericTypeParameter:
		return fmt.Sprintf("%s%s", verb(write), common.TypeDefinitionSyntax(t))
	default:
		meta := t.GetDefinitionMeta()
		verb := verb(write)
		typeArgumentsString := ""
		if len(meta.TypeParameters) > 0 {
			typeArguments := make([]string, 2*len(meta.TypeParameters))
			if len(meta.TypeArguments) > 0 {
				for i, p := range meta.TypeArguments {
					typeArguments[2*i] = common.TypeSyntax(p)
					typeArguments[2*i+1] = typeRwFunction(p, write)
				}
			} else {
				for i, p := range meta.TypeParameters {
					typeArguments[2*i] = common.TypeDefinitionSyntax(p)
					typeArguments[2*i+1] = fmt.Sprintf("%s%s", verb, common.TypeDefinitionSyntax(p))
				}
			}
			typeArgumentsString = fmt.Sprintf("<%s>", strings.Join(typeArguments, ", "))
		}
		return fmt.Sprintf("%s::binary::%s%s%s", common.NamespaceIdentifierName(meta.Namespace), verb, meta.Name, typeArgumentsString)
	}
}

func typeRwFunction(t dsl.Type, write bool) string {
	switch t := t.(type) {
	case nil:
		return fmt.Sprintf("yardl::binary::%sMonostate", verb(write))
	case *dsl.SimpleType:
		return typeDefinitionRwFunction(t.ResolvedDefinition, write)
	case *dsl.GeneralizedType:
		scalarType := t.ToScalar()
		scalarFunction := func() string {
			if t.Cases.IsSingle() {
				return typeRwFunction(t.Cases[0].Type, write)
			}
			if t.Cases.IsOptional() {
				return fmt.Sprintf("yardl::binary::%sOptional<%s, %s>", verb(write), common.TypeSyntax(t.Cases[1].Type), typeRwFunction(t.Cases[1].Type, write))
			}

			templateArguments := make([]string, 2*len(t.Cases))
			for i, c := range t.Cases {
				templateArguments[2*i] = common.TypeSyntax(c.Type)
				templateArguments[2*i+1] = typeRwFunction(c.Type, write)
			}

			return fmt.Sprintf("%sUnion<%s>", verb(write), strings.Join(templateArguments, ", "))
		}()

		switch td := t.Dimensionality.(type) {
		case nil, *dsl.Stream:
			return scalarFunction
		case *dsl.Vector:
			if td.Length == nil {
				return fmt.Sprintf("yardl::binary::%sVector<%s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction)
			}
			return fmt.Sprintf("yardl::binary::%sArray<%s, %s, %d>", verb(write), common.TypeSyntax(scalarType), scalarFunction, *td.Length)
		case *dsl.Array:
			if td.IsFixed() {
				lengths := make([]string, len(*td.Dimensions))
				for i, d := range *td.Dimensions {
					lengths[i] = strconv.FormatUint(*d.Length, 10)
				}
				return fmt.Sprintf("yardl::binary::%sFixedNDArray<%s, %s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction, strings.Join(lengths, ", "))
			}
			if td.HasKnownNumberOfDimensions() {
				return fmt.Sprintf("yardl::binary::%sNDArray<%s, %s, %d>", verb(write), common.TypeSyntax(scalarType), scalarFunction, len(*td.Dimensions))
			}

			return fmt.Sprintf("yardl::binary::%sDynamicNDArray<%s, %s>", verb(write), common.TypeSyntax(scalarType), scalarFunction)
		case *dsl.Map:
			return fmt.Sprintf("yardl::binary::%sMap<%s, %s, %s, %s>", verb(write), common.TypeSyntax(td.KeyType), common.TypeSyntax(scalarType), typeRwFunction(td.KeyType, write), scalarFunction)
		default:
			panic(fmt.Sprintf("Unknown dimensionality type %T", td))
		}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
}

func BinaryWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sWriter", p.Name)
}

func QualifiedBinaryWriterClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::binary::%s", common.TypeNamespaceIdentifierName(p), BinaryWriterClassName(p))
}

func BinaryReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%sReader", p.Name)
}

func QualifiedBinaryReaderClassName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("%s::binary::%s", common.TypeNamespaceIdentifierName(p), BinaryReaderClassName(p))
}
