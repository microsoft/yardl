// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package validation

import (
	"fmt"
)

type ValidationWarning struct {
	Message string
	File    string
	Line    *int
	Column  *int
}

func (e ValidationWarning) String() string {
	prefix := fmt.Sprintf("⚠️  %s:", e.File)
	if e.Line != nil {
		prefix = fmt.Sprintf("%s%d:", prefix, *e.Line)
	}
	if e.Column != nil {
		prefix = fmt.Sprintf("%s%d:", prefix, *e.Column)
	}

	return fmt.Sprintf("%s %v", prefix, e.Message)
}
