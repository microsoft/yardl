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
				fmt.Fprintf(w, "function obj = %s()\n", common.AbstractWriterName(p))
				common.WriteBlockBody(w, func() { w.WriteStringln("obj.state_ = 0;") })
				w.WriteStringln("")

				// Destructor
				// w.WriteStringln("function delete(obj)")
				// common.WriteBlockBody(w, func() {
				// })
				// w.WriteStringln("")

				// Close method
				w.WriteStringln("function close(obj)")
				common.WriteBlockBody(w, func() {
					if len(p.Sequence) > 0 && p.Sequence[len(p.Sequence)-1].IsStream() {
						fmt.Fprintf(w, "if obj.state_ == %d\n", len(p.Sequence)*2-1)
						common.WriteBlockBody(w, func() {
							w.WriteStringln("obj.end_stream_();")
							w.WriteStringln("obj.close_();")
							w.WriteStringln("return")
						})
					}
					w.WriteStringln("obj.close_();")
					fmt.Fprintf(w, "if obj.state_ ~= %d\n", len(p.Sequence)*2)
					common.WriteBlockBody(w, func() {
						w.WriteStringln("expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));")
						w.WriteStringln(`throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));`)
					})
				})
				w.WriteStringln("")

				// Public write methods
				for i, step := range p.Sequence {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d", i))
					common.WriteComment(w, step.Comment)
					fmt.Fprintf(w, "function %s(obj, value)\n", common.ProtocolWriteMethodName(step))
					common.WriteBlockBody(w, func() {
						prevIsStream := i > 0 && p.Sequence[i-1].IsStream()
						if prevIsStream {
							fmt.Fprintf(w, "if obj.state_ == %d\n", i*2-1)
							w.Indented(func() {
								w.WriteStringln("obj.end_stream_();")
								fmt.Fprintf(w, "obj.state_ = %d;\n", i*2)
							})
							w.WriteString("else")
						}

						if step.IsStream() {
							fmt.Fprintf(w, "if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= %d\n", i*2)
						} else {
							fmt.Fprintf(w, "if obj.state_ ~= %d\n", i*2)
						}
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "obj.raise_unexpected_state_(%d);\n", i*2)
						})
						w.WriteStringln("")
						fmt.Fprintf(w, "obj.%s(value);\n", common.ProtocolWriteImplMethodName(step))
						if step.IsStream() {
							fmt.Fprintf(w, "obj.state_ = %d;\n", i*2+1)
						} else {
							fmt.Fprintf(w, "obj.state_ = %d;\n", (i+1)*2)
						}
					})

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
					fmt.Fprintf(w, "%s(obj, value)\n", common.ProtocolWriteImplMethodName(step))
				}
				w.WriteStringln("")

				// end_stream method
				w.WriteStringln("end_stream_(obj)")
				// underlying close method
				w.WriteStringln("close_(obj)")
			})
			w.WriteStringln("")

			// Private methods
			w.WriteStringln("methods (Access=private)")
			common.WriteBlockBody(w, func() {
				// _raise_unexpected_state method
				w.WriteStringln("function raise_unexpected_state_(obj, actual)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("expected_method = obj.state_to_method_name_(obj.state_);")
					w.WriteStringln("actual_method = obj.state_to_method_name_(actual);")
					w.WriteStringln(`throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));`)
				})
				w.WriteStringln("")

				w.WriteStringln("function name = state_to_method_name_(obj, state)")
				common.WriteBlockBody(w, func() {
					for i, step := range p.Sequence {
						fmt.Fprintf(w, "if state == %d\n", i*2)
						w.Indented(func() {
							fmt.Fprintf(w, "name = '%s';\n", common.ProtocolWriteMethodName(step))
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
			})
			w.WriteStringln("")

			w.WriteStringln("methods")
			common.WriteBlockBody(w, func() {
				// Constructor
				fmt.Fprintf(w, "function obj = %s()\n", common.AbstractReaderName(p))
				common.WriteBlockBody(w, func() {
					w.WriteStringln("obj.state_ = 0;")
				})
				w.WriteStringln("")

				// Destructor
				// w.WriteStringln("function delete(obj)")
				// common.WriteBlockBody(w, func() {
				// })
				// w.WriteStringln("")

				// Close method
				w.WriteStringln("function close(obj)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("obj.close_();")
					fmt.Fprintf(w, "if obj.state_ ~= %d\n", len(p.Sequence)*2)
					common.WriteBlockBody(w, func() {
						w.WriteStringln("if mod(obj.state_, 2) == 1")
						w.Indented(func() {
							w.WriteStringln("previous_method = obj.state_to_method_name_(obj.state_ - 1);")
							w.WriteStringln(`throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. The iterable returned by '%s' was not fully consumed.", previous_method));`)
						})
						w.WriteStringln("else")
						common.WriteBlockBody(w, func() {
							w.WriteStringln("expected_method = obj.state_to_method_name_(obj.state_);")
							w.WriteStringln(`throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));`)
						})
					})
				})
				w.WriteStringln("")

				// Public read methods
				for i, step := range p.Sequence {
					common.WriteComment(w, fmt.Sprintf("Ordinal %d", i))
					common.WriteComment(w, step.Comment)
					fmt.Fprintf(w, "function value = %s(obj)\n", common.ProtocolReadMethodName(step))
					common.WriteBlockBody(w, func() {
						fmt.Fprintf(w, "if obj.state_ ~= %d\n", i*2)
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "obj.raise_unexpected_state_(%d);\n", i*2)
						})
						w.WriteStringln("")

						fmt.Fprintf(w, "value = obj.%s();\n", common.ProtocolReadImplMethodName(step))
						if step.IsStream() {
							// fmt.Fprintf(w, "obj.state_ = %d;\n", i*2+1)
							// fmt.Fprintf(w, "value = obj.wrap_iterable_(value, %d);\n", (i+1)*2)
							fmt.Fprintf(w, "obj.state_ = %d;\n", (i+1)*2)
						} else {
							fmt.Fprintf(w, "obj.state_ = %d;\n", (i+1)*2)
						}
					})
					w.WriteStringln("")
				}

				// copy_to method
				fmt.Fprintf(w, "function copy_to(obj, writer)\n")
				common.WriteBlockBody(w, func() {
					for _, step := range p.Sequence {
						fmt.Fprintf(w, "writer.%s(obj.%s());\n", common.ProtocolWriteMethodName(step), common.ProtocolReadMethodName(step))
					}
				})
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Static)")
			common.WriteBlockBody(w, func() {
				w.WriteStringln("function res = schema()")
				common.WriteBlockBody(w, func() {
					fmt.Fprintf(w, "res = %s.schema;\n", common.AbstractWriterName(p))
				})
			})
			w.WriteStringln("")

			// Protected abstract methods
			w.WriteStringln("methods (Abstract, Access=protected)")
			common.WriteBlockBody(w, func() {
				for _, step := range p.Sequence {
					fmt.Fprintf(w, "%s(obj, value)\n", common.ProtocolReadImplMethodName(step))
				}

				w.WriteStringln("")
				w.WriteStringln("close_(obj)")
			})
			w.WriteStringln("")

			w.WriteStringln("methods (Access=private)")
			common.WriteBlockBody(w, func() {
				// wrap_iterable method
				w.WriteStringln("function value = wrap_iterable_(obj, iterable, final_state)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("% This is a no-op... In python, it's yield from iterable")
					w.WriteStringln("value = iterable;")
					w.WriteStringln("obj.state_ = final_state;")
				})
				w.WriteStringln("")

				// raise_unexpected_state method
				w.WriteStringln("function raise_unexpected_state_(obj, actual)")
				common.WriteBlockBody(w, func() {
					w.WriteStringln("actual_method = obj.state_to_method_name_(actual);")
					w.WriteStringln("if mod(obj.state_, 2) == 1")
					w.Indented(func() {
						w.WriteStringln("previous_method = obj.state_to_method_name_(obj.state_ - 1);")
						w.WriteStringln(`throw(yardl.ProtocolError("Received call to '%s' but the iterable returned by '%s' was not fully consumed.", actual_method, previous_method));`)
					})
					w.WriteStringln("else")
					common.WriteBlockBody(w, func() {
						w.WriteStringln("expected_method = obj.state_to_method_name_(obj.state_);")
						w.WriteStringln(`throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));`)
					})
				})
				w.WriteStringln("")

				// state_to_method_name method
				w.WriteStringln("function name = state_to_method_name_(obj, state)")
				common.WriteBlockBody(w, func() {
					for i, step := range p.Sequence {
						fmt.Fprintf(w, "if state == %d\n", i*2)
						common.WriteBlockBody(w, func() {
							fmt.Fprintf(w, "name = '%s';\n", common.ProtocolReadMethodName(step))
						})
					}
					w.WriteStringln("name = '<unknown>';")
				})
			})
		})
	})
}
