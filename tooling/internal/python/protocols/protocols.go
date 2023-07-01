package protocols

import (
	"bytes"
	"fmt"
	"path"

	"github.com/microsoft/yardl/tooling/internal/formatting"
	"github.com/microsoft/yardl/tooling/internal/iocommon"
	"github.com/microsoft/yardl/tooling/internal/python/common"
	"github.com/microsoft/yardl/tooling/pkg/dsl"
)

func WriteProtocols(ns *dsl.Namespace, st dsl.SymbolTable, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import abc
import collections.abc
import datetime
import numpy as np
import numpy.typing as npt
import typing
from . import *
from . import yardl_types as yardl
`)

	for _, p := range ns.Protocols {
		writeAbstractWriter(w, p, st, ns)
		writeAbstractReader(w, p, ns)
	}

	writeExceptions(w)

	definitionsPath := path.Join(packageDir, "protocols.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeAbstractWriter(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition, st dsl.SymbolTable, ns *dsl.Namespace) {
	fmt.Fprintf(w, "class %s(abc.ABC):\n", common.AbstractWriterName(p))
	w.Indented(func() {
		common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Abstract writer for the %s protocol.", p.Name), p.Comment)
		w.WriteStringln("")

		// __init__
		w.WriteStringln("def __init__(self) -> None:")
		w.Indented(func() {
			w.WriteStringln("self._state = 0")
		})
		w.WriteStringln("")

		// schema field
		fmt.Fprintf(w, `schema = """%s"""`, dsl.GetProtocolSchemaString(p, st))
		w.WriteStringln("\n")

		// dunder methods
		w.WriteStringln("def __enter__(self):")
		w.Indented(func() {
			w.WriteStringln("return self")
		})
		w.WriteStringln("")

		w.WriteStringln("def __exit__(self, exc_type: type[BaseException] | None, exc: BaseException | None, traceback: typing.Any | None) -> None:")
		w.Indented(func() {
			w.WriteStringln("self.close()")
			fmt.Fprintf(w, "if exc is None and self._state != %d:\n", len(p.Sequence))
			w.Indented(func() {
				w.WriteStringln("expected_method = self._state_to_method_name(self._state)")
				w.WriteStringln(`raise ProtocolException(f"Protocol writer closed before all steps were called. Expected to call to '{expected_method}'.")`)
			})
		})
		w.WriteStringln("")

		// public write methods
		for i, step := range p.Sequence {
			valueType := common.TypeSyntax(step.Type, ns.Name)
			if step.IsStream() {
				valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
			}

			fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteMethodName(step), valueType)
			w.Indented(func() {
				common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Ordinal %d", i), step.Comment)
				fmt.Fprintf(w, "if self._state != %d:\n", i)
				w.Indented(func() {
					fmt.Fprintf(w, "self._raise_unexpected_state(%d)\n", i)
				})
				w.WriteStringln("")
				fmt.Fprintf(w, "self.%s(value)\n", common.ProtocolWriteImplMethodName(step))
				fmt.Fprintf(w, "self._state = %d\n", i+1)
			})
			w.WriteStringln("")
		}

		// protected abstract write methods
		for _, step := range p.Sequence {
			valueType := common.TypeSyntax(step.Type, ns.Name)
			if step.IsStream() {
				valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
			}

			w.WriteStringln("@abc.abstractmethod")
			fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteImplMethodName(step), valueType)
			w.Indented(func() {
				w.WriteStringln("raise NotImplementedError()")
			})
			w.WriteStringln("")
		}

		// close method
		w.WriteStringln("@abc.abstractmethod")
		w.WriteStringln("def close(self) -> None:")
		w.Indented(func() {
			w.WriteStringln("raise NotImplementedError()")
		})
		w.WriteStringln("")

		// _raise_unexpected_state method
		w.WriteStringln("def _raise_unexpected_state(self, actual: int) -> None:")
		w.Indented(func() {
			w.WriteStringln("expected_method = self._state_to_method_name(self._state)")
			w.WriteStringln("actual_method = self._state_to_method_name(actual)")
			w.WriteStringln(`raise ProtocolException(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")`)
		})
		w.WriteStringln("")

		// _state_to_method_name method
		w.WriteStringln("def _state_to_method_name(self, state: int) -> str:")
		w.Indented(func() {
			for i, step := range p.Sequence {
				fmt.Fprintf(w, "if state == %d:\n", i)
				w.Indented(func() {
					fmt.Fprintf(w, "return '%s'\n", common.ProtocolWriteMethodName(step))
				})
			}
			w.WriteStringln(`return "<unknown>"`)
		})
		w.WriteStringln("")
	})
}

func writeAbstractReader(w *formatting.IndentedWriter, p *dsl.ProtocolDefinition, ns *dsl.Namespace) {
	fmt.Fprintf(w, "class %s(abc.ABC):\n", common.AbstractReaderName(p))
	w.Indented(func() {
		common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Abstract reader for the %s protocol.", p.Name), p.Comment)
		w.WriteStringln("")

		// init method
		w.WriteStringln("def __init__(self, read_as_numpy: Types = Types.NONE) -> None:")
		w.Indented(func() {
			w.WriteStringln("self._read_as_numpy = read_as_numpy")
			w.WriteStringln("self._state = 0")
		})
		w.WriteStringln("")

		// schema field
		fmt.Fprintf(w, `schema = %s.schema`, common.AbstractWriterName(p))
		w.WriteStringln("\n")

		// dunder methods
		w.WriteStringln("def __enter__(self):")
		w.Indented(func() {
			w.WriteStringln("return self")
		})
		w.WriteStringln("")

		w.WriteStringln("def __exit__(self, exc_type: type[BaseException] | None, exc: BaseException | None, traceback: typing.Any | None) -> None:")
		w.Indented(func() {
			w.WriteStringln("self.close()")
			fmt.Fprintf(w, "if exc is None and self._state != %d:\n", len(p.Sequence)*2)
			w.Indented(func() {
				w.WriteStringln(`if self._state % 2 == 1:
    previous_method = self._state_to_method_name(self._state - 1)
    raise ProtocolException(f"Protocol reader closed before all data was consumed. The iterable returned by '{previous_method}' was not fully consumed.")
else:
    expected_method = self._state_to_method_name(self._state)
    raise ProtocolException(f"Protocol reader closed before all data was consumed. Expected call to '{expected_method}'.")
	`)
			})
		})
		w.WriteStringln("")

		w.WriteStringln("@abc.abstractmethod")
		w.WriteStringln("def close(self) -> None:")
		w.Indented(func() {
			w.WriteStringln("raise NotImplementedError()")
		})
		w.WriteStringln("")

		// public read methods
		for i, step := range p.Sequence {
			valueType := common.TypeSyntax(step.Type, ns.Name)
			if step.IsStream() {
				valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
			}

			fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadMethodName(step), valueType)
			w.Indented(func() {
				common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Ordinal %d", i), step.Comment)
				fmt.Fprintf(w, "if self._state != %d:\n", i*2)
				w.Indented(func() {
					fmt.Fprintf(w, "self._raise_unexpected_state(%d)\n", i*2)
				})
				w.WriteStringln("")

				fmt.Fprintf(w, "value = self.%s()\n", common.ProtocolReadImplMethodName(step))
				if step.IsStream() {
					fmt.Fprintf(w, "self._state = %d\n", i*2+1)
					fmt.Fprintf(w, "return self._wrap_iterable(value, %d)\n", (i+1)*2)
				} else {
					fmt.Fprintf(w, "self._state = %d\n", (i+1)*2)
					w.WriteStringln("return value")
				}
			})
			w.WriteStringln("")
		}

		// copy_to method
		fmt.Fprintf(w, "def copy_to(self, writer: %s) -> None:\n", common.AbstractWriterName(p))
		w.Indented(func() {
			for _, step := range p.Sequence {
				fmt.Fprintf(w, "writer.%s(self.%s())\n", common.ProtocolWriteMethodName(step), common.ProtocolReadMethodName(step))
			}
		})
		w.WriteStringln("")

		// protected abstract read methods
		for _, step := range p.Sequence {
			valueType := common.TypeSyntax(step.Type, ns.Name)
			if step.IsStream() {
				valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
			}

			w.WriteStringln("@abc.abstractmethod")
			fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadImplMethodName(step), valueType)
			w.Indented(func() {
				w.WriteStringln("raise NotImplementedError()")
			})
			w.WriteStringln("")
		}

		// _wrap_iterable method
		w.WriteStringln("T = typing.TypeVar('T')")
		w.WriteStringln("def _wrap_iterable(self, iterable: collections.abc.Iterable[T], final_state: int) -> collections.abc.Iterable[T]:")
		w.Indented(func() {
			w.WriteStringln("yield from iterable")
			w.WriteStringln("self._state = final_state")
		})
		w.WriteStringln("")

		// _raise_unexpected_state method
		w.WriteStringln("def _raise_unexpected_state(self, actual: int) -> None:")
		w.Indented(func() {
			w.WriteStringln("actual_method = self._state_to_method_name(actual)")
			w.WriteStringln(`if self._state % 2 == 1:
    previous_method = self._state_to_method_name(self._state - 1)
    raise ProtocolException(f"Received call to '{actual_method}' but the iterable returned by '{previous_method}' was not fully consumed.")
else:
    expected_method = self._state_to_method_name(self._state)
    raise ProtocolException(f"Expected to call to '{expected_method}' but received call to '{actual_method}'.")
	`)
		})

		// _state_to_method_name method
		w.WriteStringln("def _state_to_method_name(self, state: int) -> str:")
		w.Indented(func() {
			for i, step := range p.Sequence {
				fmt.Fprintf(w, "if state == %d:\n", i*2)
				w.Indented(func() {
					fmt.Fprintf(w, "return '%s'\n", common.ProtocolReadMethodName(step))
				})
			}
			w.WriteStringln(`return "<unknown>"`)
		})
		w.WriteStringln("")
	})
}

func writeExceptions(w *formatting.IndentedWriter) {
	w.WriteStringln("class ProtocolException(Exception):")
	w.Indented(func() {
		w.WriteStringln(`"""Raised when the contract of a protocol is not respected."""`)
		w.WriteStringln("pass")
	})
}
