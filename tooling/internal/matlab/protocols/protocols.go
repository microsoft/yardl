// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package protocols

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/matlab/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteProtocols(fw *common.MatlabFileWriter, ns *dsl.Namespace, st dsl.SymbolTable) error {
	for _, p := range ns.Protocols {
		if err := writeAbstractWriter(fw, p, ns, st); err != nil {
			return err
		}
		if err := writeAbstractReader(fw, p, ns); err != nil {
			return err
		}
	}

	return nil
}

func writeAbstractWriter(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition, ns *dsl.Namespace, st dsl.SymbolTable) error {
	return fw.WriteFile(common.AbstractWriterName(p), func(w *formatting.IndentedWriter) {
		common.WriteComment(w, fmt.Sprintf("Abstract writer for protocol %s", p.Name))
		common.WriteComment(w, p.Comment)
		fmt.Fprintf(w, "classdef (Abstract) %s < handle\n", common.AbstractWriterName(p))

		common.WriteBlockBody(w, func() {

			w.WriteStringln("properties (Access=protected)")
			common.WriteBlockBody(w, func() { w.WriteStringln("state_") })
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				// Constructor
				fmt.Fprintf(w, "function self = %s()\n", common.AbstractWriterName(p))
				common.WriteBlockBody(w, func() { w.WriteStringln("self.state_ = 0;") })
				w.WriteStringln("")

				// Close method
				w.WriteStringln("function close(self)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("self.close_();")
					fmt.Fprintf(w, "if self.state_ ~= %d\n", len(p.Sequence))
					common.WriteBlockBody(w, func() {
						w.WriteStringln("expected_method = self.state_to_method_name_(self.state_);")
						w.WriteStringln(`throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));`)
					})
				})
				w.WriteStringln("")

				// Public write methods
				for i, step := range p.Sequence {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d", i))
					common.WriteComment(w, step.Comment)
					fmt.Fprintf(w, "function %s(self, value)\n", common.ProtocolWriteMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "if self.state_ ~= %d\n", i)
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "self.raise_unexpected_state_(%d);\n", i)
						})
						w.WriteStringln("")

						fmt.Fprintf(w, "self.%s(value);\n", common.ProtocolWriteImplMethodName(step))
						if !step.IsStream() {
							fmt.Fprintf(w, "self.state_ = %d;\n", i+1)
						}
					})

					if step.IsStream() {
						// End stream method
						w.WriteStringln("")
						fmt.Fprintf(w, "function %s(self)\n", common.ProtocolEndMethodName(step))
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "if self.state_ ~= %d\n", i)
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "self.raise_unexpected_state_(%d);\n", i)
							})
							w.WriteStringln("")
							fmt.Fprintf(w, "self.end_stream_();\n")
							fmt.Fprintf(w, "self.state_ = %d;\n", i+1)
						})
					}

					if i < len(p.Sequence)-1 {
						w.WriteStringln("")
					}
				}
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Static)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("function res = schema()")
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "res = string('%s');\n", dsl.GetProtocolSchemaString(p, st))
				})
			})
			w.WriteStringln("")

			// Protected abstract write methods
			w.WriteStringln("methods (Abstract, Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "%s(self, value)\n", common.ProtocolWriteImplMethodName(step))
				}
				w.WriteStringln("")

				// end_stream method
				w.WriteStringln("end_stream_(self)")
				// underlying close method
				w.WriteStringln("close_(self)")
			})
			w.WriteStringln("")

			// Private methods
			w.WriteStringln("methods (Access=private)")
			common.WriteBlockBody(w, func() {
				// raise_unexpected_state method
				w.WriteStringln("function raise_unexpected_state_(self, actual)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("expected_method = self.state_to_method_name_(self.state_);")
					w.WriteStringln("actual_method = self.state_to_method_name_(actual);")
					w.WriteStringln(`throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));`)
				})
				w.WriteStringln("")

				w.WriteStringln("function name = state_to_method_name_(self, state)")
				common.WriteBlockBody(w, func() {
					for i, step := range p.Sequence {
						fmt.Fprintf(w, "if state == %d\n", i)
						w.Indented(func() {
							if step.IsStream() {
								fmt.Fprintf(w, "name = \"%s or %s\";\n", common.ProtocolWriteMethodName(step), common.ProtocolEndMethodName(step))
							} else {
								fmt.Fprintf(w, "name = \"%s\";\n", common.ProtocolWriteMethodName(step))
							}
						})
						w.WriteString("else")
					}
					w.WriteStringln("")
					common.WriteBlockBody(w, func() {
						w.WriteStringln("name = '<unknown>';")
					})
				})
			})
		})
	})
}

func writeAbstractReader(fw *common.MatlabFileWriter, p *dsl.ProtocolDefinition, ns *dsl.Namespace) error {
	return fw.WriteFile(common.AbstractReaderName(p), func(w *formatting.IndentedWriter) {
		common.WriteComment(w, p.Comment)
		fmt.Fprintf(w, "classdef %s < handle\n", common.AbstractReaderName(p))

		common.WriteBlockBody(w, func() {

			w.WriteStringln("properties (Access=protected)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("state_")
				w.WriteStringln("skip_completed_check_")
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				// Constructor
				fmt.Fprintf(w, "function self = %s(options)\n", common.AbstractReaderName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("arguments")
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "options.skip_completed_check (1,1) logical = false\n")
					})
					w.WriteStringln("self.state_ = 0;")
					w.WriteStringln("self.skip_completed_check_ = options.skip_completed_check;")
				})
				w.WriteStringln("")

				// Close method
				w.WriteStringln("function close(self)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("self.close_();")
					fmt.Fprintf(w, "if ~self.skip_completed_check_ && self.state_ ~= %d\n", len(p.Sequence))
					common.WriteBlockBody(w, func() {
						w.WriteStringln("expected_method = self.state_to_method_name_(self.state_);")
						w.WriteStringln(`throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));`)
					})
				})
				w.WriteStringln("")

				// Public has/read methods
				for i, step := range p.Sequence {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d", i))
					if step.IsStream() {
						fmt.Fprintf(w, "function more = %s(self)\n", common.ProtocolHasMoreMethodName(step))
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "if self.state_ ~= %d\n", i)
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "self.raise_unexpected_state_(%d);\n", i)
							})
							w.WriteStringln("")

							fmt.Fprintf(w, "more = self.%s();\n", common.ProtocolHasMoreImplMethodName(step))
							w.WriteStringln("if ~more")
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "self.state_ = %d;\n", i+1)
							})
						})
						w.WriteStringln("")
					}
					common.WriteComment(w, step.Comment)
					fmt.Fprintf(w, "function value = %s(self)\n", common.ProtocolReadMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "if self.state_ ~= %d\n", i)
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "self.raise_unexpected_state_(%d);\n", i)
						})
						w.WriteStringln("")

						fmt.Fprintf(w, "value = self.%s();\n", common.ProtocolReadImplMethodName(step))
						if !step.IsStream() {
							fmt.Fprintf(w, "self.state_ = %d;\n", i+1)
						}
					})
					w.WriteStringln("")
				}

				// copy_to method
				fmt.Fprintf(w, "function copy_to(self, writer)\n")
				common.WriteBlockBody(w, func() {
					for _, step := range p.Sequence {
						if step.IsStream() {
							fmt.Fprintf(w, "while self.%s()\n", common.ProtocolHasMoreMethodName(step))
							common.WriteBlockBody(w, func() {
								fmt.Fprintf(w, "item = self.%s();\n", common.ProtocolReadMethodName(step))
								fmt.Fprintf(w, "writer.%s({item});\n", common.ProtocolWriteMethodName(step))
							})
							fmt.Fprintf(w, "writer.%s();\n", common.ProtocolEndMethodName(step))
						} else {
							fmt.Fprintf(w, "writer.%s(self.%s());\n", common.ProtocolWriteMethodName(step), common.ProtocolReadMethodName(step))
						}
					}
				})
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Static)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("function res = schema()")
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "res = %s.%s.schema;\n", common.NamespaceIdentifierName(ns.Name), common.AbstractWriterName(p))
				})
			})
			w.WriteStringln("")

			// Protected abstract methods
			w.WriteStringln("methods (Abstract, Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					if step.IsStream() {
						fmt.Fprintf(w, "%s(self)\n", common.ProtocolHasMoreImplMethodName(step))
					}
					fmt.Fprintf(w, "%s(self)\n", common.ProtocolReadImplMethodName(step))
				}

				w.WriteStringln("")
				w.WriteStringln("close_(self)")
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=private)")
			common.WriteBlockBody(w, func() {
				// raise_unexpected_state method
				w.WriteStringln("function raise_unexpected_state_(self, actual)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("actual_method = self.state_to_method_name_(actual);")
					w.WriteStringln("expected_method = self.state_to_method_name_(self.state_);")
					w.WriteStringln(`throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));`)
				})
				w.WriteStringln("")

				// state_to_method_name method
				w.WriteStringln("function name = state_to_method_name_(self, state)")
				common.WriteBlockBody(w, func() {
					for i, step := range p.Sequence {
						fmt.Fprintf(w, "if state == %d\n", i)
						w.Indented(func() {
							fmt.Fprintf(w, "name = \"%s\";\n", common.ProtocolReadMethodName(step))
						})
						w.WriteString("else")
					}
					w.WriteStringln("")
					common.WriteBlockBody(w, func() {
						w.WriteStringln("name = \"<unknown>\";")
					})
				})
			})
		})
	})
}
