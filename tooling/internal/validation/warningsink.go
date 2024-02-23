// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package validation

import (
	"sort"
)

type WarningSink struct {
	Warnings []ValidationWarning
}

func (e *WarningSink) Add(err ValidationWarning) {
	e.Warnings = append(e.Warnings, err)
}

func (e *WarningSink) AsStrings() []string {
	if len(e.Warnings) == 0 {
		return nil
	}

	// sort errors by filename, line number, column number, then error message
	sort.Slice(e.Warnings, func(i, j int) bool {
		iWrn := e.Warnings[i]
		jWrn := e.Warnings[j]

		if iWrn.File != jWrn.File {
			return iWrn.File < jWrn.File
		}

		if iLine, jLine := pointerValueOrDefault(iWrn.Line, 0), pointerValueOrDefault(jWrn.Line, 0); iLine != jLine {
			return iLine < jLine
		}

		if iColumn, jColumn := pointerValueOrDefault(iWrn.Column, 0), pointerValueOrDefault(jWrn.Column, 0); iColumn != jColumn {
			return iColumn < jColumn
		}

		return iWrn.Message < jWrn.Message
	})

	messages := make([]string, len(e.Warnings))
	for i, wrn := range e.Warnings {
		messages[i] = wrn.String()
	}

	return messages
}
