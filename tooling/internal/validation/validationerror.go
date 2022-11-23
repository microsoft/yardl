// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ValidationError struct {
	Message error
	File    string
	Line    *int
	Column  *int
}

func NewValidationError(underlyingError error, file string) ValidationError {
	validationError, ok := underlyingError.(ValidationError)

	if ok {
		validationError.File = file
	} else {
		validationError = ValidationError{
			Message: underlyingError,
			File:    file,
		}

		// extract the line number from the message
		r := regexp.MustCompile(`^yaml: (unmarshal errors:\s*)?line (\d+): `)

		if groups := r.FindStringSubmatch(underlyingError.Error()); groups != nil {
			line, _ := strconv.Atoi(groups[2])
			validationError.Line = &line
			validationError.Message = errors.New(underlyingError.Error()[len(groups[0]):])
		}
	}

	return validationError
}

func (e ValidationError) Error() string {
	prefix := fmt.Sprintf("❌ %s:", e.File)
	if e.Line != nil {
		prefix = fmt.Sprintf("%s%d:", prefix, *e.Line)
	}
	if e.Column != nil {
		prefix = fmt.Sprintf("%s%d:", prefix, *e.Column)
	}

	return fmt.Sprintf("%s %v", prefix, e.Message)
}
