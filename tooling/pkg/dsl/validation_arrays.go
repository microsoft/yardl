// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import "github.com/microsoft/yardl/tooling/internal/validation"

func validateArrayAndVectorDimensions(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		switch t := node.(type) {
		case *Array:
			if t.Dimensions != nil && len(*t.Dimensions) > 0 {
				nullLengthCount := 0
				notNullLengthCount := 0
				dimensionNames := make(map[string]bool)
				for _, dim := range *t.Dimensions {
					if dim.Length == nil {
						nullLengthCount++
					} else {
						notNullLengthCount++
					}

					if dim.Name != nil {
						if !memberNameRegex.MatchString(*dim.Name) {
							errorSink.Add(validationError(t, "dimension name '%s' must match the format %s", *dim.Name, typeNameRegex.String()))
						}

						if _, found := dimensionNames[*dim.Name]; found {
							errorSink.Add(validationError(t, "a dimension with the name '%s' is already defined on the array", *dim.Name))
						} else {
							dimensionNames[*dim.Name] = true
						}
					}
				}

				if (notNullLengthCount > 0) == (nullLengthCount > 0) {
					errorSink.Add(validationError(node, "lengths must either be specified on all dimensions or none of them"))
				}
			}
		}

		self.VisitChildren(node)
	})

	return env
}
