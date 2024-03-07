// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package protocols

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/cpp/common"
	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
)

func WriteProtocols(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	err := writeHeader(env, options)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteString("#include \"protocols.h\"\n\n")
	w.WriteStringln("#ifdef _MSC_VER")
	w.WriteStringln("#define unlikely(x) x")
	w.WriteStringln("#else")
	w.WriteString("#define unlikely(x) __builtin_expect((x), 0)\n")
	w.WriteString("#endif\n\n")

	for _, ns := range env.Namespaces {
		if ns.IsTopLevel {
			fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
			writeDefinitions(w, ns, env.SymbolTable)
			fmt.Fprintf(w, "} // namespace %s\n", common.NamespaceIdentifierName(ns.Name))
		}
	}

	definitionsPath := path.Join(options.SourcesOutputDir, "protocols.cc")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeHeader(env *dsl.Environment, options packaging.CppCodegenOptions) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "  ")
	common.WriteGeneratedFileHeader(w)

	w.WriteStringln(`#pragma once
#include "types.h"
`)

	for _, ns := range env.Namespaces {
		if ns.IsTopLevel {
			fmt.Fprintf(w, "namespace %s {\n", common.NamespaceIdentifierName(ns.Name))
			writeDeclarations(w, ns)
			fmt.Fprintf(w, "} // namespace %s\n", common.NamespaceIdentifierName(ns.Name))
		}
	}

	definitionsPath := path.Join(options.SourcesOutputDir, "protocols.h")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeDeclarations(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	fmt.Fprintf(w, "enum class Version {\n")
	w.Indented(func() {
		for _, v := range ns.Versions {
			fmt.Fprintf(w, "%s,\n", v)
		}
		fmt.Fprintln(w, "Current")
	})
	fmt.Fprintln(w, "};")

	formatting.Delimited(w, "\n", ns.Protocols, func(w *formatting.IndentedWriter, i int, p *dsl.ProtocolDefinition) {

		// Writer
		common.WriteComment(w, fmt.Sprintf("Abstract writer for the %s protocol.", p.Name))
		common.WriteComment(w, p.Comment)
		fmt.Fprintf(w, "class %s {\n", common.AbstractWriterName(p))
		w.Indented(func() {
			fmt.Fprintln(w, "public:")
			for i, step := range p.Sequence {
				endMethodName := common.ProtocolWriteEndMethodName(step)
				common.WriteComment(w, fmt.Sprintf("Ordinal %d.", i))
				common.WriteComment(w, step.Comment)
				if step.IsStream() {
					common.WriteComment(w, fmt.Sprintf("Call this method for each element of the `%s` stream, then call `%s() when done.`", step.Name, endMethodName))
				}

				fmt.Fprintf(w, "void %s(%s const& value);\n\n", common.ProtocolWriteMethodName(step), common.TypeSyntax(step.Type))

				if step.IsStream() {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d.", i))
					common.WriteComment(w, step.Comment)
					if step.IsStream() {
						common.WriteComment(w, fmt.Sprintf("Call this method to write many values to the `%s` stream, then call `%s()` when done.", step.Name, endMethodName))
					}

					fmt.Fprintf(w, "void %s(std::vector<%s> const& values);\n\n", common.ProtocolWriteMethodName(step), common.TypeSyntax(step.Type))

					common.WriteComment(w, fmt.Sprintf("Marks the end of the `%s` stream.", step.Name))
					fmt.Fprintf(w, "void %s();\n\n", endMethodName)
				}
			}

			common.WriteComment(w, "Optionaly close this writer before destructing. Validates that all steps were completed.")
			w.WriteString("void Close();\n\n")
			fmt.Fprintf(w, "virtual ~%s() = default;\n\n", common.AbstractWriterName(p))

			common.WriteComment(w, "Flushes all buffered data.")
			w.WriteString("virtual void Flush() {}\n\n")

			w.WriteStringln("protected:")
			for _, step := range p.Sequence {
				fmt.Fprintf(w, "virtual void %s(%s const& value) = 0;\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))

				if step.IsStream() {
					fmt.Fprintf(w, "virtual void %s(std::vector<%s> const& value);\n", common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
					fmt.Fprintf(w, "virtual void %s() = 0;\n", common.ProtocolWriteEndImplMethodName(step))
				}
			}
			w.WriteString("virtual void CloseImpl() {}\n\n")

			w.WriteString("static std::string schema_;\n\n")
			w.WriteString("static std::vector<std::string> previous_schemas_;\n\n")
			w.WriteString("static std::string SchemaFromVersion(Version version);\n\n")

			w.WriteStringln("private:")
			w.WriteString("uint8_t state_ = 0;\n\n")

			fmt.Fprintf(w, "friend class %s;\n", common.AbstractReaderName(p))
		})
		fmt.Fprint(w, "};\n\n")

		// Reader
		common.WriteComment(w, fmt.Sprintf("Abstract reader for the %s protocol.", p.Name))
		common.WriteComment(w, p.Comment)
		fmt.Fprintf(w, "class %s {\n", common.AbstractReaderName(p))
		w.Indented(func() {
			w.WriteString("public:\n")
			for i, step := range p.Sequence {
				common.WriteComment(w, fmt.Sprintf("Ordinal %d.", i))
				common.WriteComment(w, step.Comment)

				returnType := "void"
				if step.IsStream() {
					returnType = "[[nodiscard]] bool"
				}
				fmt.Fprintf(w, "%s %s(%s& value);\n\n", returnType, common.ProtocolReadMethodName(step), common.TypeSyntax(step.Type))

				if step.IsStream() {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d.", i))
					common.WriteComment(w, step.Comment)
					fmt.Fprintf(w, "%s %s(std::vector<%s>& values);\n\n", returnType, common.ProtocolReadMethodName(step), common.TypeSyntax(step.Type))
				}
			}

			common.WriteComment(w, "Optionaly close this writer before destructing. Validates that all steps were completely read.")
			w.WriteString("void Close();\n\n")

			fmt.Fprintf(w, "void CopyTo(%s& writer", common.AbstractWriterName(p))
			for _, s := range p.Sequence {
				if s.IsStream() {
					fmt.Fprintf(w, ", size_t %s_buffer_size = 1", formatting.ToSnakeCase(s.Name))
				}
			}
			w.WriteString(");\n\n")

			fmt.Fprintf(w, "virtual ~%s() = default;\n\n", common.AbstractReaderName(p))

			w.WriteStringln("protected:")
			for _, step := range p.Sequence {
				returnType := "void"
				if step.IsStream() {
					returnType = "bool"
				}
				fmt.Fprintf(w, "virtual %s %s(%s& value) = 0;\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))

				if step.IsStream() {
					fmt.Fprintf(w, "virtual %s %s(std::vector<%s>& values);\n", returnType, common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
				}
			}

			w.WriteString("virtual void CloseImpl() {}\n")

			w.WriteString("static std::string schema_;\n\n")
			w.WriteString("static std::vector<std::string> previous_schemas_;\n\n")
			w.WriteString("static Version VersionFromSchema(const std::string& schema);\n\n")

			w.WriteStringln("private:")
			w.WriteStringln("uint8_t state_ = 0;")
		})
		fmt.Fprint(w, "};\n")
	})
}

func writeDefinitions(w *formatting.IndentedWriter, ns *dsl.Namespace, symbolTable dsl.SymbolTable) {
	formatting.Delimited(w, "\n", ns.Protocols, func(w *formatting.IndentedWriter, i int, p *dsl.ProtocolDefinition) {
		w.WriteString("namespace {\n")
		writeInvalidWriterStateMethod(w, p)
		writeInvalidReaderStateMethod(w, p)

		w.WriteStringln("} // namespace \n")

		// Writers

		fmt.Fprintf(w, "std::string %s::schema_ = R\"(%s)\";\n\n", common.AbstractWriterName(p), dsl.GetProtocolSchemaString(p, symbolTable))
		fmt.Fprintf(w, "std::vector<std::string> %s::previous_schemas_ = {\n", common.AbstractWriterName(p))
		w.Indented(func() {
			for _, versionLabel := range ns.Versions {
				change := p.Versions[versionLabel]
				if change != nil {
					fmt.Fprintf(w, "R\"(%s)\",\n", change.PreviousSchema)
				} else {
					fmt.Fprintf(w, "%s::schema_,\n", common.AbstractWriterName(p))
				}
			}
		})
		fmt.Fprintf(w, "};\n\n")

		fmt.Fprintf(w, "std::string %s::SchemaFromVersion(Version version) {\n", common.AbstractWriterName(p))
		w.Indented(func() {
			w.WriteStringln("switch (version) {")
			for i, versionLabel := range ns.Versions {
				fmt.Fprintf(w, "case Version::%s: return previous_schemas_[%d]; break;\n", versionLabel, i)
			}
			fmt.Fprintf(w, "case Version::Current: return %s::schema_; break;\n", common.AbstractWriterName(p))
			fmt.Fprintf(w, "default: throw std::runtime_error(\"The version does not correspond to any schema supported by protocol %s.\");\n", p.Name)
			w.WriteStringln("}")
			w.WriteStringln("")
		})
		fmt.Fprintln(w, "}")

		for i, step := range p.Sequence {
			writeWriteMethod := func(signature string, variableName string) {
				w.WriteString(signature)
				w.Indented(func() {
					fmt.Fprintf(w, "if (unlikely(state_ != %d)) {\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "%s(%d, false, state_);\n", invalidWriterStateMethodName(p), i)
					})
					w.WriteString("}\n\n")

					fmt.Fprintf(w, "%s(%s);\n", common.ProtocolWriteImplMethodName(step), variableName)
					if !step.IsStream() {
						fmt.Fprintf(w, "state_ = %d;\n", i+1)
					}
				})
				w.WriteString("}\n\n")
			}

			writeWriteMethod(fmt.Sprintf("void %s::%s(%s const& value) {\n", common.AbstractWriterName(p), common.ProtocolWriteMethodName(step), common.TypeSyntax(step.Type)), "value")

			if step.IsStream() {
				writeWriteMethod(fmt.Sprintf("void %s::%s(std::vector<%s> const& values) {\n",
					common.AbstractWriterName(p), common.ProtocolWriteMethodName(step),
					common.TypeSyntax(step.Type)), "values")

				fmt.Fprintf(w, "void %s::%s() {\n", common.AbstractWriterName(p), common.ProtocolWriteEndMethodName(step))
				w.Indented(func() {
					fmt.Fprintf(w, "if (unlikely(state_ != %d)) {\n", i)
					w.Indented(func() {
						fmt.Fprintf(w, "%s(%d, true, state_);\n", invalidWriterStateMethodName(p), i)
					})
					w.WriteString("}\n\n")
					fmt.Fprintf(w, "%s();\n", common.ProtocolWriteEndImplMethodName(step))
					fmt.Fprintf(w, "state_ = %d;\n", i+1)
				})
				w.WriteString("}\n\n")

				common.WriteComment(w, "fallback implementation")
				fmt.Fprintf(w, "void %s::%s(std::vector<%s> const& values) {\n", common.AbstractWriterName(p), common.ProtocolWriteImplMethodName(step), common.TypeSyntax(step.Type))
				w.Indented(func() {
					w.WriteString("for (auto const& v : values) {\n")
					w.Indented(func() {
						fmt.Fprintf(w, "%s(v);\n", common.ProtocolWriteImplMethodName(step))
					})
					w.WriteString("}\n")
				})
				w.WriteString("}\n\n")
			}
		}

		fmt.Fprintf(w, "void %s::Close() {\n", common.AbstractWriterName(p))
		w.Indented(func() {
			fmt.Fprintf(w, "if (unlikely(state_ != %d)) {\n", len(p.Sequence))
			w.Indented(func() {
				fmt.Fprintf(w, "%s(%d, false, state_);\n", invalidWriterStateMethodName(p), len(p.Sequence))
			})
			w.WriteString("}\n\n")
			fmt.Fprintf(w, "CloseImpl();\n")
		})
		w.WriteString("}\n\n")

		// Readers

		fmt.Fprintf(w, "std::string %s::schema_ = %s::schema_;\n\n", common.AbstractReaderName(p), common.AbstractWriterName(p))
		fmt.Fprintf(w, "std::vector<std::string> %s::previous_schemas_ = %s::previous_schemas_;\n\n", common.AbstractReaderName(p), common.AbstractWriterName(p))

		fmt.Fprintf(w, "Version %s::VersionFromSchema(std::string const& schema) {\n", common.AbstractReaderName(p))
		w.Indented(func() {
			fmt.Fprintf(w, "if (schema == %s::schema_) {\n", common.AbstractWriterName(p))
			w.Indented(func() {
				w.WriteStringln("return Version::Current;")
			})
			w.WriteStringln("}")
			for i, versionLabel := range ns.Versions {
				fmt.Fprintf(w, "else if (schema == previous_schemas_[%d]) {\n", i)
				w.Indented(func() {
					fmt.Fprintf(w, "return Version::%s;\n", versionLabel)
				})
				w.WriteStringln("}")
			}
			fmt.Fprintf(w, "throw std::runtime_error(\"The schema does not match any version supported by protocol %s.\");\n", p.Name)
		})
		fmt.Fprintln(w, "}")

		for i, step := range p.Sequence {
			returnType := "void"
			if step.IsStream() {
				returnType = "bool"
			}

			fmt.Fprintf(w, "%s %s::%s(%s& value) {\n", returnType, common.AbstractReaderName(p), common.ProtocolReadMethodName(step), common.TypeSyntax(step.Type))
			w.Indented(func() {
				writeReaderStateCheckIfStatement(w, p, i, false)
				if step.IsStream() {
					w.WriteString("bool result = ")
				}
				fmt.Fprintf(w, "%s(value);\n", common.ProtocolReadImplMethodName(step))
				if step.IsStream() {
					w.WriteString("if (!result) {\n")
					w.Indented(func() {
						fmt.Fprintf(w, "state_ = %d;\n", 2*(i+1))
					})
					w.WriteStringln("}")
					w.WriteStringln("return result;")
				} else {
					fmt.Fprintf(w, "state_ = %d;\n", 2*(i+1))
				}
			})
			w.WriteString("}\n\n")

			if step.IsStream() {
				fmt.Fprintf(w, "%s %s::%s(std::vector<%s>& values) {\n", returnType, common.AbstractReaderName(p), common.ProtocolReadMethodName(step), common.TypeSyntax(step.Type))
				w.Indented(func() {
					w.WriteStringln("if (values.capacity() == 0) {")
					w.Indented(func() {
						w.WriteStringln("throw std::runtime_error(\"vector must have a nonzero capacity.\");")
					})
					w.WriteStringln("}")
					writeReaderStateCheckIfStatement(w, p, i, true)

					fmt.Fprintf(w, "if (!%s(values)) {\n", common.ProtocolReadImplMethodName(step))
					w.Indented(func() {
						fmt.Fprintf(w, "state_ = %d;\n", 2*i+1)
						w.WriteStringln("return values.size() > 0;")
					})
					w.WriteStringln("}")
					w.WriteStringln("return true;")
				})
				w.WriteString("}\n\n")

				common.WriteComment(w, "fallback implementation")
				fmt.Fprintf(w, "%s %s::%s(std::vector<%s>& values) {\n", returnType, common.AbstractReaderName(p), common.ProtocolReadImplMethodName(step), common.TypeSyntax(step.Type))
				w.Indented(func() {
					w.WriteStringln("size_t i = 0;")
					w.WriteStringln("while (true) {")
					w.Indented(func() {
						w.WriteStringln("if (i == values.size()) {")
						w.Indented(func() {
							w.WriteStringln("values.resize(i + 1);")
						})
						w.WriteStringln("}")
						fmt.Fprintf(w, "if (!%s(values[i])) {\n", common.ProtocolReadImplMethodName(step))
						w.Indented(func() {
							w.WriteStringln("values.resize(i);")
							w.WriteStringln("return false;")
						})
						w.WriteStringln("}")
						w.WriteStringln("i++;")
						w.WriteStringln("if (i == values.capacity()) {")
						w.Indented(func() {
							w.WriteStringln("return true;")
						})
						w.WriteStringln("}")
					})
					w.WriteStringln("}")
				})
				w.WriteStringln("}\n")
			}
		}

		fmt.Fprintf(w, "void %s::Close() {\n", common.AbstractReaderName(p))
		w.Indented(func() {
			expectedState := len(p.Sequence) * 2
			fmt.Fprintf(w, "if (unlikely(state_ != %d)) {\n", expectedState)

			w.Indented(func() {
				writeReaderStateUnobservedCompletionCheck(w, p, len(p.Sequence)-1, expectedState)
			})
			w.WriteString("}\n\n")
			fmt.Fprintf(w, "CloseImpl();\n")
		})
		w.WriteString("}\n")

		writeProtocolCopyToMethod(w, p)
	})
}

func writeReaderStateUnobservedCompletionCheck(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition, prevStepIndex, expectedState int) {
	if prevStepIndex >= 0 && p.Sequence[prevStepIndex].IsStream() {
		previousUnobservedcompletionState := expectedState - 1
		fmt.Fprintf(w, "if (state_ == %d) {\n", previousUnobservedcompletionState)
		w.Indented(func() {
			fmt.Fprintf(w, "state_ = %d;\n", expectedState)
		})
		w.WriteStringln("} else {")
		w.Indented(func() {
			fmt.Fprintf(w, "%s(%d, state_);\n", invalidReaderStateMethodName(p), expectedState)
		})
		w.WriteStringln("}")
	} else {
		fmt.Fprintf(w, "%s(%d, state_);\n", invalidReaderStateMethodName(p), expectedState)
	}
}

func writeReaderStateCheckIfStatement(w *formatting.IndentedWriter, protocol *dsl.ProtocolDefinition, stepIndex int, isBatchOverload bool) {
	step := protocol.Sequence[stepIndex]
	expectedState := 2 * stepIndex
	unobservedCompletionState := expectedState + 1
	nextState := expectedState + 2

	fmt.Fprintf(w, "if (unlikely(state_ != %d)) {\n", expectedState)
	w.Indented(func() {
		if step.IsStream() {
			fmt.Fprintf(w, "if (state_ == %d) {\n", unobservedCompletionState)
			w.Indented(func() {
				fmt.Fprintf(w, "state_ = %d;\n", nextState)
				if isBatchOverload {
					w.WriteStringln("values.clear();")
				}
				w.WriteStringln("return false;")
			})
			w.WriteStringln("}")
		}

		writeReaderStateUnobservedCompletionCheck(w, protocol, stepIndex-1, expectedState)
	})
	w.WriteString("}\n\n")
}

func writeInvalidWriterStateMethod(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	fmt.Fprintf(w, "void %s(uint8_t attempted, [[maybe_unused]] bool end, uint8_t current) {\n", invalidWriterStateMethodName(p))
	w.Indented(func() {
		w.WriteStringln("std::string expected_method;")
		w.WriteStringln("switch (current) {")
		for i, step := range p.Sequence {
			methodName := fmt.Sprintf("%s()", common.ProtocolWriteMethodName(step))
			if step.IsStream() {
				methodName = fmt.Sprintf("%s or %s()", methodName, common.ProtocolWriteEndMethodName(step))
			}
			fmt.Fprintf(w, "case %d: expected_method = \"%s\"; break;\n", i, methodName)
		}
		w.WriteStringln("}")

		w.WriteStringln("std::string attempted_method;")
		w.WriteStringln("switch (attempted) {")
		for i, step := range p.Sequence {
			if step.IsStream() {
				fmt.Fprintf(w, "case %d: attempted_method = end ? \"%s()\" : \"%s()\"; break;\n", i, common.ProtocolWriteEndMethodName(step), common.ProtocolWriteMethodName(step))
			} else {
				fmt.Fprintf(w, "case %d: attempted_method = \"%s()\"; break;\n", i, common.ProtocolWriteMethodName(step))
			}
		}
		fmt.Fprintf(w, "case %d: attempted_method = \"%s\"; break;\n", len(p.Sequence), "Close()")
		w.WriteStringln("}")

		fmt.Fprintf(w, "throw std::runtime_error(\"Expected call to \" + expected_method + \" but received call to \" + attempted_method + \" instead.\");\n")
	})

	w.WriteString("}\n\n")
}

func writeInvalidReaderStateMethod(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	fmt.Fprintf(w, "void %s(uint8_t attempted, uint8_t current) {\n", invalidReaderStateMethodName(p))
	w.Indented(func() {
		w.WriteString("auto f = [](uint8_t i) -> std::string {\n")
		w.Indented(func() {
			w.WriteStringln("switch (i/2) {")
			for i, step := range p.Sequence {
				fmt.Fprintf(w, "case %d: return \"%s()\";\n", i, common.ProtocolReadMethodName(step))
			}
			fmt.Fprintf(w, "case %d: return \"Close()\";\n", len(p.Sequence))
			w.WriteStringln("default: return \"<unknown>\";")
			w.WriteStringln("}")
		})
		w.WriteString("};\n")

		fmt.Fprintf(w, "throw std::runtime_error(\"Expected call to \" + f(current) + \" but received call to \" + f(attempted) + \" instead.\");\n")
	})

	w.WriteString("}\n\n")
}

func invalidWriterStateMethodName(p *dsl.ProtocolDefinition) string {
	writerInvalidStateMethod := fmt.Sprintf("%sInvalidState", common.AbstractWriterName(p))
	return writerInvalidStateMethod
}

func invalidReaderStateMethodName(p *dsl.ProtocolDefinition) string {
	writerInvalidStateMethod := fmt.Sprintf("%sInvalidState", common.AbstractReaderName(p))
	return writerInvalidStateMethod
}

func writeProtocolCopyToMethod(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition) {
	fmt.Fprintf(w, "void %s::CopyTo(%s& writer", common.AbstractReaderName(p), common.AbstractWriterName(p))
	for _, s := range p.Sequence {
		if s.IsStream() {
			fmt.Fprintf(w, ", size_t %s_buffer_size", formatting.ToSnakeCase(s.Name))
		}
	}
	w.WriteStringln(") {")

	w.Indented(func() {
		if len(p.Sequence) == 0 {
			w.WriteStringln("(void)(writer);")
		} else {
			for _, s := range p.Sequence {
				if s.IsStream() {
					bufferSizeParameterName := fmt.Sprintf("%s_buffer_size", formatting.ToSnakeCase(s.Name))
					fmt.Fprintf(w, "if (%s > 1) {\n", bufferSizeParameterName)
					w.Indented(func() {
						fmt.Fprintf(w, "std::vector<%s> values;\n", common.TypeSyntax(s.Type))
						fmt.Fprintf(w, "values.reserve(%s);\n", bufferSizeParameterName)
						fmt.Fprintf(w, "while(%s(values)) {\n", common.ProtocolReadMethodName(s))
						fmt.Fprintf(w.Indent(), "writer.%s(values);\n", common.ProtocolWriteMethodName(s))
						w.WriteStringln("}")
						fmt.Fprintf(w, "writer.%s();\n", common.ProtocolWriteEndMethodName(s))
					})
					w.WriteStringln("} else {")
					w.Indented(func() {
						fmt.Fprintf(w, "%s value;\n", common.TypeSyntax(s.Type))
						fmt.Fprintf(w, "while(%s(value)) {\n", common.ProtocolReadMethodName(s))
						fmt.Fprintf(w.Indent(), "writer.%s(value);\n", common.ProtocolWriteMethodName(s))
						w.WriteStringln("}")
						fmt.Fprintf(w, "writer.%s();\n", common.ProtocolWriteEndMethodName(s))
					})
					w.WriteString("}\n")
				} else {
					w.WriteStringln("{")
					w.Indented(func() {
						fmt.Fprintf(w, "%s value;\n", common.TypeSyntax(s.Type))
						fmt.Fprintf(w, "%s(value);\n", common.ProtocolReadMethodName(s))
						fmt.Fprintf(w, "writer.%s(value);\n", common.ProtocolWriteMethodName(s))
					})
					w.WriteString("}\n")
				}
			}
		}
	})
	w.WriteString("}\n")
}
