// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"regexp"

	"github.com/microsoft/yardl/tooling/internal/validation"
)

type ValidationPass func(env *Environment, errorSink *validation.ErrorSink) *Environment

var (
	memberNameRegex = regexp.MustCompile("^[a-z][a-zA-Z0-9]{0,63}$")
	typeNameRegex   = regexp.MustCompile("^[A-Z][a-zA-Z0-9]{0,63}$")
)

func Validate(namespaces []*Namespace) (*Environment, error) {
	env := &Environment{
		Namespaces:  namespaces,
		SymbolTable: map[string]TypeDefinition{},
	}

	errorSink := validation.ErrorSink{}

	passes := []ValidationPass{
		validateTypeDefinitionNames,
		validateGenericTypeDefinitions,
		validateRecordFieldNames,
		validateProtocolSequenceNames,
		validateArrayAndVectorDimensions,
		validateMaps,
		validateStreams,
		buildSymbolTable,
		resolveTypes,
		assignUnionCaseLabels,
		topologicalSortTypes,
		convertGenericReferences,
		validateUnionCases,
		validateEnums,
		resolveComputedFields,
		validateGenericParametersUsed,
	}

	for _, pass := range passes {
		env = pass(env, &errorSink)
	}

	return env, errorSink.AsError()
}

func validationError(node Node, message string, args ...any) validation.ValidationError {
	return validation.ValidationError{
		Message: fmt.Errorf(message, args...),
		File:    node.GetNodeMeta().File,
		Line:    &node.GetNodeMeta().Line,
		Column:  &node.GetNodeMeta().Column,
	}
}

func validateTypeDefinitionNames(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		switch t := node.(type) {
		case TypeDefinition:
			name := t.GetDefinitionMeta().Name
			if !typeNameRegex.MatchString(name) {
				errorSink.Add(validationError(t, "type name '%s' must be PascalCased matching the format %s", name, typeNameRegex.String()))
			}

		default:
			self.VisitChildren(node)
		}
	})

	return env
}

func validateRecordFieldNames(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		record, ok := node.(*RecordDefinition)
		if !ok {
			self.VisitChildren(node)
			return
		}

		fields := make(map[string]bool)

		for _, field := range record.Fields {
			if !memberNameRegex.MatchString(field.Name) {
				errorSink.Add(validationError(field, "field name '%s' must be camelCased matching the format %s", field.Name, memberNameRegex.String()))
			}

			if _, found := fields[field.Name]; found {
				errorSink.Add(validationError(field, "a field with the name '%s' is already defined on the record '%s'", field.Name, record.Name))
			}

			fields[field.Name] = true
		}

		for _, field := range record.ComputedFields {
			if !memberNameRegex.MatchString(field.Name) {
				errorSink.Add(validationError(field, "computed field name '%s' must be camelCased matching the format %s", field.Name, memberNameRegex.String()))
			}

			if _, found := fields[field.Name]; found {
				errorSink.Add(validationError(field, "a field or computed field with the name '%s' is already defined on the record '%s'", field.Name, record.Name))
			}

			fields[field.Name] = true
		}
	})

	return env
}

func validateProtocolSequenceNames(env *Environment, errorSink *validation.ErrorSink) *Environment {
	Visit(env, func(self Visitor, node Node) {
		protocol, ok := node.(*ProtocolDefinition)
		if !ok {
			self.VisitChildren(node)
			return
		}

		steps := make(map[string]bool)

		for _, step := range protocol.Sequence {
			if !memberNameRegex.MatchString(step.Name) {
				errorSink.Add(validationError(step, "protocol step name '%s' must be camelCased matching the format %s", step.Name, memberNameRegex.String()))
			}

			if _, found := steps[step.Name]; found {
				errorSink.Add(validationError(step, "a sequence step with the name '%s' is already defined on the protocol '%s'", step.Name, protocol.Name))
			}

			steps[step.Name] = true
		}
	})

	return env
}

func validateStreams(env *Environment, errorSink *validation.ErrorSink) *Environment {
	VisitWithContext(env, nil, func(self VisitorWithContext[Node], node Node, context Node) {
		switch node.(type) {
		case TypeDefinition:
			self.VisitChildren(node, node)
		case *Stream:
			if _, isProtocol := (context).(*ProtocolDefinition); !isProtocol {
				errorSink.Add(validationError(node, "!streams can only be declared as top-level protocol sequence elements"))
			}

			self.VisitChildren(node, node)
		default:
			self.VisitChildren(node, context)
		}
	})
	return env
}
