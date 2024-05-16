// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package mocks

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteMocks(fw *common.MatlabFileWriter, ns *dsl.Namespace) error {
	for _, p := range ns.Protocols {
		if err := writeProtocolMock(fw, p); err != nil {
			return err
		}
		if err := writeProtocolTestWriter(fw, p); err != nil {
			return err
		}
	}
	return nil
}

func writeProtocolMock(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition) error {
	expectedName := func(step *dsl.ProtocolStep) string {
		return fmt.Sprintf("expected_%s", formatting.ToSnakeCase(step.Name))
	}
	return fw.WriteFile(mockWriterName(p), func(w *formatting.IndentedWriter) {
		abstractWriterName := fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(p.Namespace), common.AbstractWriterName(p))
		fmt.Fprintf(w, "classdef %s < matlab.mixin.Copyable & %s\n", mockWriterName(p), abstractWriterName)
		common.WriteBlockBody(w, func() {
			w.WriteStringln("properties")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("testCase_")
				for _, step := range p.Sequence {
					w.WriteStringln(expectedName(step))
				}
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "function self = %s(testCase)\n", mockWriterName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("self.testCase_ = testCase;")
					for _, step := range p.Sequence {
						if step.IsStream() {
							fmt.Fprintf(w, "self.%s = {};\n", expectedName(step))
						} else {
							fmt.Fprintf(w, "self.%s = yardl.None;\n", expectedName(step))
						}
					}
				})
				w.WriteStringln("")

				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function expect_%s(self, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						if step.IsStream() {
							w.WriteStringln("if iscell(value)")
							common.WriteBlockBody(w, func() {
								w.WriteStringln("for n = 1:numel(value)")
								common.WriteBlockBody(w, func() {
									fmt.Fprintf(w, "self.%s{end+1} = value{n};\n", expectedName(step))
								})
								w.WriteStringln("return;")
							})

							w.WriteStringln("shape = size(value);")
							w.WriteStringln("lastDim = ndims(value);")
							w.WriteStringln("count = shape(lastDim);")
							w.WriteStringln("index = repelem({':'}, lastDim-1);")
							w.WriteStringln("for n = 1:count")
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "self.%s{end+1} = value(index{:}, n);\n", expectedName(step))
							})
						} else {
							fmt.Fprintf(w, "self.%s = yardl.Optional(value);\n", expectedName(step))
						}
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function verify(self)")
				common.WriteBlockBody(w, func() {
					for _, step := range p.Sequence {
						diagnostic := fmt.Sprintf("Expected call to %s was not received", common.ProtocolWriteImplMethodName(step))
						if step.IsStream() {
							fmt.Fprintf(w, "self.testCase_.verifyTrue(isempty(self.%s), \"%s\");\n", expectedName(step), diagnostic)
						} else {
							fmt.Fprintf(w, "self.testCase_.verifyEqual(self.%s, yardl.None, \"%s\");\n", expectedName(step), diagnostic)
						}
					}
				})
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function %s(self, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						if step.IsStream() {
							w.WriteStringln("assert(iscell(value));")
							w.WriteStringln("assert(isscalar(value));")
							fmt.Fprintf(w, "self.testCase_.verifyFalse(isempty(self.%s), \"Unexpected call to %s\");\n", expectedName(step), common.ProtocolWriteImplMethodName(step))
							fmt.Fprintf(w, "self.testCase_.verifyEqual(value{1}, self.%s{1}, \"Unexpected argument value for call to %s\");\n", expectedName(step), common.ProtocolWriteImplMethodName(step))
							fmt.Fprintf(w, "self.%s = self.%s(2:end);\n", expectedName(step), expectedName(step))
						} else {
							fmt.Fprintf(w, "self.testCase_.verifyTrue(self.%s.has_value(), \"Unexpected call to %s\");\n", expectedName(step), common.ProtocolWriteImplMethodName(step))
							fmt.Fprintf(w, "self.testCase_.verifyEqual(value, self.%s.value, \"Unexpected argument value for call to %s\");\n", expectedName(step), common.ProtocolWriteImplMethodName(step))
							fmt.Fprintf(w, "self.%s = yardl.None;\n", expectedName(step))
						}
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function close_(self)")
				w.WriteStringln("end")

				w.WriteStringln("function end_stream_(self)")
				w.WriteStringln("end")
			})
		})
	})
}

func writeProtocolTestWriter(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition) error {
	return fw.WriteFile(testWriterName(p), func(w *formatting.IndentedWriter) {
		abstractWriterName := fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(p.Namespace), common.AbstractWriterName(p))
		fmt.Fprintf(w, "classdef %s < %s\n", testWriterName(p), abstractWriterName)
		common.WriteBlockBody(w, func() {
			w.WriteStringln("properties (Access = private)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("writer_")
				w.WriteStringln("create_reader_")
				w.WriteStringln("mock_writer_")
				w.WriteStringln("close_called_")
				w.WriteStringln("filename_")
				w.WriteStringln("format_")
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "function self = %s(testCase, format, create_writer, create_reader)\n", testWriterName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("self.filename_ = tempname();")
					w.WriteStringln("self.format_ = format;")
					w.WriteStringln("self.writer_ = create_writer(self.filename_);")
					w.WriteStringln("self.create_reader_ = create_reader;")
					mockWriterName := fmt.Sprintf("%s.testing.%s", common.NamespaceIdentifierName(p.Namespace), mockWriterName(p))
					fmt.Fprintf(w, "self.mock_writer_ = %s(testCase);\n", mockWriterName)
					w.WriteStringln("self.close_called_ = false;")
				})
				w.WriteStringln("")

				w.WriteStringln("function delete(self)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("delete(self.filename_);")
					w.WriteStringln("if ~self.close_called_")
					common.WriteBlockBody(w, func() {
						common.WriteComment(w, "ADD_FAILURE() << ...;")
						fmt.Fprintf(w, "throw(yardl.RuntimeError(\"Close() must be called on '%s' to verify mocks\"));\n", testWriterName(p))
					})
				})

				for _, step := range p.Sequence {
					if step.IsStream() {
						fmt.Fprintf(w, "function %s(self)\n", common.ProtocolEndMethodName(step))
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "%s@%s(self);\n", common.ProtocolEndMethodName(step), abstractWriterName)
							fmt.Fprintf(w, "self.writer_.%s();\n", common.ProtocolEndMethodName(step))
						})
						w.WriteStringln("")
					}
				}
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function %s(self, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "self.writer_.%s(value);\n", common.ProtocolWriteMethodName(step))
						fmt.Fprintf(w, "self.mock_writer_.expect_%s(value);\n", common.ProtocolWriteImplMethodName(step))
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function close_(self)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("self.close_called_ = true;")
					w.WriteStringln("self.writer_.close();")
					w.WriteStringln("mock_copy = copy(self.mock_writer_);")

					w.WriteStringln("")
					w.WriteStringln("reader = self.create_reader_(self.filename_);")
					w.WriteStringln("reader.copy_to(self.mock_writer_);")
					w.WriteStringln("reader.close();")
					w.WriteStringln("self.mock_writer_.verify();")
					w.WriteStringln("self.mock_writer_.close();")

					w.WriteStringln("")
					w.WriteStringln("translated = invoke_translator(self.filename_, self.format_, self.format_);")
					w.WriteStringln("reader = self.create_reader_(translated);")
					w.WriteStringln("reader.copy_to(mock_copy);")
					w.WriteStringln("reader.close();")
					w.WriteStringln("mock_copy.verify();")
					w.WriteStringln("mock_copy.close();")
					w.WriteStringln("delete(translated);")
				})
				w.WriteStringln("")

				w.WriteStringln("function end_stream_(self)")
				common.WriteBlockBody(w, func() {})
			})
		})
	})
}

func mockWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Mock%sWriter", formatting.ToPascalCase(p.Name))
}

func testWriterName(p *dsl.ProtocolDefinition) string {
	return fmt.Sprintf("Test%sWriter", formatting.ToPascalCase(p.Name))
}
