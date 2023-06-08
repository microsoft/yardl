// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package formatting

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dlclark/regexp2"
)

type IndentedWriter struct {
	indentLevel   int
	indentString  string
	indentPending *bool
	writer        io.Writer
	buf           bytes.Buffer
}

func NewIndentedWriter(writer io.Writer, indentString string) *IndentedWriter {
	indentPending := true
	return &IndentedWriter{indentString: indentString, writer: writer, indentPending: &indentPending}
}

func (w *IndentedWriter) Indent() *IndentedWriter {
	indentedWriter := *w
	indentedWriter.indentLevel++
	return &indentedWriter
}

func (w *IndentedWriter) GetIndentString() string {
	return w.indentString
}

func (w *IndentedWriter) Indented(f func()) {
	defer func() { w.indentLevel-- }()
	w.indentLevel++
	f()
}

func (w *IndentedWriter) Write(payload []byte) (int, error) {
	w.buf.Reset()

	for _, b := range payload {
		if *w.indentPending && b != '\n' {
			w.buf.WriteString(strings.Repeat(w.indentString, w.indentLevel))
			*w.indentPending = false
		}

		w.buf.WriteByte(b)

		if b == '\n' {
			*w.indentPending = true
		}
	}

	n, err := w.writer.Write(w.buf.Bytes())
	if err != nil {
		// never return more than original length to satisfy io.Writer interface
		if n > len(payload) {
			n = len(payload)
		}
		return n, err
	}

	// return original length to satisfy io.Writer interface
	return len(payload), nil
}

func (w *IndentedWriter) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w *IndentedWriter) WriteStringln(s string) (n int, err error) {
	return fmt.Fprintln(w, s)
}

func Delimited[T any](w *IndentedWriter, delimiter string, items []T, action func(w *IndentedWriter, i int, item T)) {
	first := true
	for i, v := range items {
		if first {
			first = false
		} else {
			w.WriteString(delimiter)
		}
		action(w, i, v)
	}
}

func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}

	if !strings.ContainsAny(s, "_ -") {
		if start := s[0]; start >= 'A' && start <= 'Z' {
			return s
		}

		return strings.ToUpper(s[:1]) + s[1:]
	}

	var b strings.Builder
	tokens := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
	for _, part := range tokens {
		b.WriteString(ToPascalCase(part))
	}

	return b.String()
}

var snakeCaseRegex = regexp2.MustCompile(`((?<=\p{Ll})(\p{Lu}))|(?<!(\b|_)\p{Ll})((?<=\p{Ll})(\d))|(?<!\b|_)(\p{Lu})(?=\p{Ll})`, regexp2.ExplicitCapture)

func delimitWithUnderscores(str string) string {
	s, err := snakeCaseRegex.Replace(str, `_$&`, -1, -1)
	if err != nil {
		panic(err)
	}
	return s
}

func ToSnakeCase(str string) string {
	str = delimitWithUnderscores(str)
	return strings.ToLower(str)
}

func ToUpperSnakeCase(str string) string {
	str = delimitWithUnderscores(str)
	return strings.ToUpper(str)
}
