// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"github.com/microsoft/yardl/tooling/internal/validation"
)

type EvolutionPass func(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) *Environment

func ValidateEvolution(env *Environment, predecessor *Environment) (*Environment, error) {

	errorSink := validation.ErrorSink{}

	passes := []EvolutionPass{
		ensureNoChanges,
	}

	for _, pass := range passes {
		env = pass(env, predecessor, &errorSink)
	}

	return env, errorSink.AsError()
}

func ensureNoChanges(env *Environment, predecessor *Environment, errorSink *validation.ErrorSink) *Environment {
	return env
}
