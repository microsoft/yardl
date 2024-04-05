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
	return fw.WriteFile(mockWriterName(p), func(w *formatting.IndentedWriter) {
		abstractWriterName := fmt.Sprintf("%s.%s", common.NamespaceIdentifierName(p.Namespace), common.AbstractWriterName(p))
		fmt.Fprintf(w, "classdef %s < matlab.mixin.Copyable & %s\n", mockWriterName(p), abstractWriterName)
		common.WriteBlockBody(w, func() {
			w.WriteStringln("properties")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("testCase_")
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "%swritten\n", common.ProtocolWriteImplMethodName(step))
				}
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				fmt.Fprintf(w, "function obj = %s(testCase)\n", mockWriterName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("obj.testCase_ = testCase;")
					for _, step := range p.Sequence {
						fmt.Fprintf(w, "obj.%swritten = yardl.None;\n", common.ProtocolWriteImplMethodName(step))
					}
				})
				w.WriteStringln("")

				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function expect_%s(obj, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "if obj.%swritten.has_value()\n", common.ProtocolWriteImplMethodName(step))
						w.Indented(func() {
							fmt.Fprintf(w, "last_dim = ndims(value);\n")
							fmt.Fprintf(w, "obj.%swritten = yardl.Optional(cat(last_dim, obj.%swritten.value, value));\n", common.ProtocolWriteImplMethodName(step), common.ProtocolWriteImplMethodName(step))
						})
						w.WriteStringln("else")
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "obj.%swritten = yardl.Optional(value);\n", common.ProtocolWriteImplMethodName(step))
						})
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function verify(obj)")
				common.WriteBlockBody(w, func() {
					for _, step := range p.Sequence {
						diagnostic := fmt.Sprintf("Expected call to %s was not received", common.ProtocolWriteImplMethodName(step))
						fmt.Fprintf(w, "obj.testCase_.verifyEqual(obj.%swritten, yardl.None, \"%s\");\n", common.ProtocolWriteImplMethodName(step), diagnostic)
					}
				})
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function %s(obj, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "obj.testCase_.verifyTrue(obj.%swritten.has_value(), \"Unexpected call to %s\");\n", common.ProtocolWriteImplMethodName(step), common.ProtocolWriteImplMethodName(step))
						fmt.Fprintf(w, "expected = obj.%swritten.value;\n", common.ProtocolWriteImplMethodName(step))
						fmt.Fprintf(w, "obj.testCase_.verifyEqual(value, expected, \"Unexpected argument value for call to %s\");\n", common.ProtocolWriteImplMethodName(step))
						fmt.Fprintf(w, "obj.%swritten = yardl.None;\n", common.ProtocolWriteImplMethodName(step))
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function close_(obj)")
				w.WriteStringln("end")

				w.WriteStringln("function end_stream_(obj)")
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
				fmt.Fprintf(w, "function obj = %s(testCase, format, create_writer, create_reader)\n", testWriterName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("obj.filename_ = tempname();")
					w.WriteStringln("obj.format_ = format;")
					w.WriteStringln("obj.writer_ = create_writer(obj.filename_);")
					w.WriteStringln("obj.create_reader_ = create_reader;")
					mockWriterName := fmt.Sprintf("%s.testing.%s", common.NamespaceIdentifierName(p.Namespace), mockWriterName(p))
					fmt.Fprintf(w, "obj.mock_writer_ = %s(testCase);\n", mockWriterName)
					w.WriteStringln("obj.close_called_ = false;")
				})
				w.WriteStringln("")

				w.WriteStringln("function delete(obj)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("delete(obj.filename_);")
					w.WriteStringln("if ~obj.close_called_")
					common.WriteBlockBody(w, func() {
						common.WriteComment(w, "ADD_FAILURE() << ...;")
						fmt.Fprintf(w, "throw(yardl.RuntimeError(\"Close() must be called on '%s' to verify mocks\"));\n", testWriterName(p))
					})
				})
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "function %s(obj, value)\n", common.ProtocolWriteImplMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "obj.writer_.%s(value);\n", common.ProtocolWriteMethodName(step))
						fmt.Fprintf(w, "obj.mock_writer_.expect_%s(value);\n", common.ProtocolWriteImplMethodName(step))
					})
					w.WriteStringln("")
				}

				w.WriteStringln("function close_(obj)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("obj.close_called_ = true;")
					w.WriteStringln("obj.writer_.close();")
					w.WriteStringln("mock_copy = copy(obj.mock_writer_);")

					w.WriteStringln("")
					w.WriteStringln("reader = obj.create_reader_(obj.filename_);")
					w.WriteStringln("reader.copy_to(obj.mock_writer_);")
					w.WriteStringln("reader.close();")
					w.WriteStringln("obj.mock_writer_.verify();")
					w.WriteStringln("obj.mock_writer_.close();")

					w.WriteStringln("")
					w.WriteStringln("translated = invoke_translator(obj.filename_, obj.format_, obj.format_);")
					w.WriteStringln("reader = obj.create_reader_(translated);")
					w.WriteStringln("reader.copy_to(mock_copy);")
					w.WriteStringln("reader.close();")
					w.WriteStringln("mock_copy.verify();")
					w.WriteStringln("mock_copy.close();")
					w.WriteStringln("delete(translated);")
				})
				w.WriteStringln("")

				w.WriteStringln("function end_stream_(obj)")
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
