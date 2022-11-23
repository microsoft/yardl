// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package validation

import (
	"errors"
	"sort"
	"strings"
)

type ErrorSink struct {
	Errors []ValidationError
}

func (e *ErrorSink) Add(err ValidationError) {
	e.Errors = append(e.Errors, err)
}

func (e *ErrorSink) AsError() error {
	if len(e.Errors) == 0 {
		return nil
	}

	// sort errors by filename, line number, column number, then error message
	sort.Slice(e.Errors, func(i, j int) bool {
		iErr := e.Errors[i]
		jErr := e.Errors[j]

		if iErr.File != jErr.File {
			return iErr.File < jErr.File
		}

		if iLine, jLine := pointerValueOrDefault(iErr.Line, 0), pointerValueOrDefault(jErr.Line, 0); iLine != jLine {
			return iLine < jLine
		}

		if iColumn, jColumn := pointerValueOrDefault(iErr.Column, 0), pointerValueOrDefault(jErr.Column, 0); iColumn != jColumn {
			return iColumn < jColumn
		}

		return iErr.Message.Error() < jErr.Message.Error()
	})

	messages := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		messages[i] = err.Error()
	}

	return errors.New(strings.Join(messages, "\n"))
}

func pointerValueOrDefault[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}

	return defaultValue
}
