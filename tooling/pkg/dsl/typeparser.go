// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"errors"
	"fmt"
	"strings"
)

type simpleTypeTree struct {
	Name           string           `json:"name"`
	TypeArguments  []simpleTypeTree `json:"args,omitempty"`
	Optional       bool             `json:"optional,omitempty"`
	PositionOffset int              `json:"positionOffset,omitempty"`
}

func (pt *simpleTypeTree) String() string {
	if len(pt.TypeArguments) == 0 {
		return pt.Name
	}
	args := make([]string, len(pt.TypeArguments))
	for i, typeArg := range pt.TypeArguments {
		args[i] = typeArg.String()
	}

	return fmt.Sprintf("%s<%s>", pt.Name, strings.Join(args, ", "))
}

type typeParser struct {
	position       int
	remainingInput string
}

func (tp *typeParser) skipWhitespace() {
	for i := 0; i < len(tp.remainingInput); i++ {
		if tp.remainingInput[i] != ' ' {
			tp.position += i
			tp.remainingInput = tp.remainingInput[i:]
			return
		}
	}

	tp.position += len(tp.remainingInput)
	tp.remainingInput = ""
}

func (tp *typeParser) advance(count int) string {
	tp.position += count
	rtn := tp.remainingInput[:count]
	tp.remainingInput = tp.remainingInput[count:]
	return rtn
}

func (tp *typeParser) consumeIdentifier() string {
	idx := strings.IndexAny(tp.remainingInput, ",<>? ")
	if idx == -1 {
		rtn := tp.remainingInput
		tp.position += len(tp.remainingInput)
		tp.remainingInput = ""
		return rtn
	}

	return tp.advance(idx)
}

func (tp *typeParser) parseTypeString() (simpleTypeTree, error) {
	parsed := simpleTypeTree{}
	tp.skipWhitespace()
	parsed.PositionOffset = tp.position
	parsed.Name = tp.consumeIdentifier()
	tp.skipWhitespace()

	if tp.remainingInput == "" {
		return parsed, nil
	}

	if tp.remainingInput[0] == '?' {
		tp.advance(1)
		tp.skipWhitespace()

		if tp.remainingInput != "" && tp.remainingInput[0] == '<' {
			return parsed, fmt.Errorf("'?' at position %d must appear after generic type arguments", tp.position)
		}

		parsed.Optional = true
		return parsed, nil
	}

	if tp.remainingInput[0] != '<' {
		return parsed, nil
	}

	tp.advance(1)

	parsed.TypeArguments = make([]simpleTypeTree, 0)

	for {
		arg, err := tp.parseTypeString()
		if err != nil {
			return parsed, err
		}
		parsed.TypeArguments = append(parsed.TypeArguments, arg)
		tp.skipWhitespace()
		if tp.remainingInput == "" {
			return parsed, errors.New("missing '>' in type string")
		}
		if arg.Name == "" {
			return parsed, fmt.Errorf("the type parameter name cannot be empty at position %d", tp.position+1)
		}

		switch tp.remainingInput[0] {
		case '>':
			tp.advance(1)
			tp.skipWhitespace()
			if tp.remainingInput != "" && tp.remainingInput[0] == '?' {
				tp.advance(1)
				parsed.Optional = true
			}
			return parsed, nil
		case ',':
			tp.advance(1)
			continue
		default:
			return parsed, fmt.Errorf("unexpected '%s' in type string at position %d", tp.remainingInput, tp.position+1)
		}
	}
}

func parseSimpleTypeString(typeString string) (simpleTypeTree, error) {
	typeTree, remaining, err := parseSimpleTypeStringAllowingRemaining(typeString)
	if err != nil {
		return typeTree, err
	}

	if remaining != "" {
		return typeTree, fmt.Errorf("unexpected trailing '%s' in type string", remaining)
	}

	return typeTree, nil
}

func parseSimpleTypeStringAllowingRemaining(typeString string) (typeTree simpleTypeTree, remaining string, err error) {
	parser := typeParser{remainingInput: typeString}
	parsed, err := parser.parseTypeString()
	if err != nil {
		return parsed, "", err
	}
	if parsed.Name == "" {
		return parsed, "", errors.New("the type name cannot be empty")
	}
	parser.skipWhitespace()

	return parsed, parser.remainingInput, nil
}

func (tree simpleTypeTree) ToType(node NodeMeta) Type {
	nodeWithPositionUpdated := node
	nodeWithPositionUpdated.Column += tree.PositionOffset

	simpleType := SimpleType{NodeMeta: nodeWithPositionUpdated, Name: tree.Name}
	for _, typeArg := range tree.TypeArguments {
		simpleType.TypeArguments = append(simpleType.TypeArguments, typeArg.ToType(node))
	}

	if tree.Optional {
		return &GeneralizedType{
			NodeMeta: node,
			Cases: TypeCases{
				&TypeCase{NodeMeta: node},
				&TypeCase{NodeMeta: node, Type: &simpleType},
			}}
	}

	return &simpleType
}
