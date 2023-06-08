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

func WriteProtocols(ns *dsl.Namespace, packageDir string) error {
	b := bytes.Buffer{}
	w := formatting.NewIndentedWriter(&b, "    ")
	common.WriteGeneratedFileHeader(w)
	w.WriteStringln(`
import abc
import collections.abc
import datetime
import numpy as np
`)

	writeProtocols(w, ns)

	definitionsPath := path.Join(packageDir, "protocols.py")
	return iocommon.WriteFileIfNeeded(definitionsPath, b.Bytes(), 0644)
}

func writeProtocols(w *formatting.IndentedWriter, ns *dsl.Namespace) {
	for _, p := range ns.Protocols {
		// abstract writer
		fmt.Fprintf(w, "class %s(abc.ABC):\n", common.AbstractWriterName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Abstract writer for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			for i, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}

				fmt.Fprintf(w, "def %s(self, value: %s) -> None:\n", common.ProtocolWriteMethodName(step), valueType)
				w.Indented(func() {
					common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Ordinal %d", i), step.Comment)
					fmt.Fprintf(w, "self.%s(value)\n", common.ProtocolWriteImplMethodName(step))
				})
				w.WriteStringln("")
			}

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
		})

		// abstract reader
		fmt.Fprintf(w, "class %s(abc.ABC):\n", common.AbstractReaderName(p))
		w.Indented(func() {
			common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Abstract reader for the %s protocol.", p.Name), p.Comment)
			w.WriteStringln("")

			for i, step := range p.Sequence {
				valueType := common.TypeSyntax(step.Type, ns.Name)
				if step.IsStream() {
					valueType = fmt.Sprintf("collections.abc.Iterable[%s]", valueType)
				}

				fmt.Fprintf(w, "def %s(self) -> %s:\n", common.ProtocolReadMethodName(step), valueType)
				w.Indented(func() {
					common.WriteDocstringWithLeadingLine(w, fmt.Sprintf("Ordinal %d", i), step.Comment)
					fmt.Fprintf(w, "return self.%s()\n", common.ProtocolReadImplMethodName(step))
				})
				w.WriteStringln("")
			}

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
		})
	}
}
