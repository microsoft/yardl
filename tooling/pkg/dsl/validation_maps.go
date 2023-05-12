// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"github.com/microsoft/yardl/tooling/internal/validation"
)

func validateMaps(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		m, ok := node.(*Map)
		if !ok {
			self.VisitChildren(node)
			return
		}

		t := GetUnderlyingType(m.KeyType)
		if st, ok := t.(*SimpleType); ok {
			switch st.ResolvedDefinition.(type) {
			case nil, PrimitiveDefinition:
				return
			}
		}

		errorSink.Add(validationError(m, "map key type must be a primitive scalar type"))
	})

	return env
}
