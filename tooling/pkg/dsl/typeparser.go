// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// type simpleTypeTree struct {
// 	Name           string                 `json:"name"`
// 	Dimensionality typeTreeDimensionality `json:"dimensionality,omitempty"`
// 	TypeArguments  []simpleTypeTree       `json:"args,omitempty"`
// 	Optional       bool                   `json:"optional,omitempty"`
// 	PositionOffset int                    `json:"positionOffset,omitempty"`
// }

// type typeTreeDimensionality interface {
// 	_typeTreeDimension()
// }

// type mapTypeTreeDimensionality struct {
// 	Key simpleTypeTree `json:"key"`
// }

// func (*mapTypeTreeDimensionality) _typeTreeDimension() {}

// func (m *mapTypeTreeDimensionality) MarshalJSON() ([]byte, error) {
// 	var i = struct {
// 		Map *simpleTypeTree `json:"map"`
// 	}{Map: &m.Key}
// 	return json.Marshal(i)
// }

// type vectorTreeDimensionality struct {
// 	Length   *int `json:"length,omitempty"`
// 	Optional bool `json:"optional,omitempty"`
// }

// func (*vectorTreeDimensionality) _typeTreeDimension() {}

// func (v *vectorTreeDimensionality) MarshalJSON() ([]byte, error) {
// 	type innerVector vectorTreeDimensionality
// 	var i = struct {
// 		Vector *innerVector `json:"vector"`
// 	}{Vector: (*innerVector)(v)}
// 	return json.Marshal(i)
// }

// type arrayTreeDimensionality struct {
// 	Length *int `json:"length,omitempty"`
// }

// func (pt *simpleTypeTree) String() string {
// 	if len(pt.TypeArguments) == 0 {
// 		return pt.Name
// 	}
// 	args := make([]string, len(pt.TypeArguments))
// 	for i, typeArg := range pt.TypeArguments {
// 		args[i] = typeArg.String()
// 	}

// 	return fmt.Sprintf("%s<%s>", pt.Name, strings.Join(args, ", "))
// }

// type typeParser struct {
// 	position       int
// 	remainingInput string
// }

// func (tp *typeParser) skipWhitespace() {
// 	for i := 0; i < len(tp.remainingInput); i++ {
// 		if tp.remainingInput[i] != ' ' {
// 			tp.position += i
// 			tp.remainingInput = tp.remainingInput[i:]
// 			return
// 		}
// 	}

// 	tp.position += len(tp.remainingInput)
// 	tp.remainingInput = ""
// }

// func (tp *typeParser) advance(count int) string {
// 	tp.position += count
// 	rtn := tp.remainingInput[:count]
// 	tp.remainingInput = tp.remainingInput[count:]
// 	return rtn
// }

// func (tp *typeParser) tryParseInteger() *int {
// 	var val *int

// 	for tp.remainingInput != "" && tp.remainingInput[0] >= '0' && tp.remainingInput[0] <= '9' {
// 		d := int(tp.remainingInput[0] - '0')
// 		if val == nil {
// 			val = &d
// 		} else {
// 			*val = *val*10 + d
// 		}
// 		tp.advance(1)
// 	}

// 	return val

// }

// func (tp *typeParser) consumeIdentifier() string {
// 	idx := strings.IndexAny(tp.remainingInput, ",<>?-  *")
// 	if idx == -1 {
// 		rtn := tp.remainingInput
// 		tp.position += len(tp.remainingInput)
// 		tp.remainingInput = ""
// 		return rtn
// 	}

// 	return tp.advance(idx)
// }

// func (tp *typeParser) parseTypeString() (simpleTypeTree, error) {
// 	parsed := simpleTypeTree{}
// 	tp.skipWhitespace()
// 	parsed.PositionOffset = tp.position
// 	parsed.Name = tp.consumeIdentifier()
// 	tp.skipWhitespace()

// 	if tp.remainingInput == "" {
// 		return parsed, nil
// 	}

// 	if tp.remainingInput[0] == '?' {
// 		tp.advance(1)
// 		tp.skipWhitespace()

// 		if tp.remainingInput != "" && tp.remainingInput[0] == '<' {
// 			return parsed, fmt.Errorf("'?' at position %d must appear after generic type arguments", tp.position)
// 		}

// 		parsed.Optional = true
// 		goto doneGenericArgs
// 	}

// 	if tp.remainingInput[0] == '<' {
// 		tp.advance(1)

// 		parsed.TypeArguments = make([]simpleTypeTree, 0)

// 		for {
// 			arg, err := tp.parseTypeString()
// 			if err != nil {
// 				return parsed, err
// 			}
// 			parsed.TypeArguments = append(parsed.TypeArguments, arg)
// 			tp.skipWhitespace()
// 			if tp.remainingInput == "" {
// 				return parsed, errors.New("missing '>' in type string")
// 			}
// 			if arg.Name == "" {
// 				return parsed, fmt.Errorf("the type parameter name cannot be empty at position %d", tp.position+1)
// 			}

// 			switch tp.remainingInput[0] {
// 			case '>':
// 				tp.advance(1)
// 				tp.skipWhitespace()
// 				if tp.remainingInput != "" && tp.remainingInput[0] == '?' {
// 					tp.advance(1)
// 					parsed.Optional = true
// 				}
// 				goto doneGenericArgs
// 			case ',':
// 				tp.advance(1)
// 				continue
// 			default:
// 				return parsed, fmt.Errorf("unexpected '%s' in type string at position %d", tp.remainingInput, tp.position+1)
// 			}
// 		}
// 	}
// doneGenericArgs:

// 	if len(tp.remainingInput) > 1 && tp.remainingInput[0] == '-' && tp.remainingInput[1] == '>' {
// 		// map type
// 		tp.advance(2)
// 		tp.skipWhitespace()
// 		if tp.remainingInput == "" {
// 			return parsed, errors.New("missing type name after '->'")
// 		}
// 		value, err := tp.parseTypeString()
// 		if err != nil {
// 			return parsed, err
// 		}

// 		value.Dimensionality = &mapTypeTreeDimensionality{Key: parsed}
// 		parsed = value
// 		// return value, nil
// 	}

// 	tp.skipWhitespace()
// 	if tp.remainingInput == "" {
// 		return parsed, nil
// 	}

// 	if tp.remainingInput[0] == '*' {
// 		// vector type
// 		tp.advance(1)
// 		tp.skipWhitespace()

// 		d := &vectorTreeDimensionality{
// 			Length: tp.tryParseInteger(),
// 		}

// 		if tp.remainingInput != "" && tp.remainingInput[0] == '?' {
// 			tp.advance(1)
// 			d.Optional = true
// 		}

// 		parsed.Dimensionality = d
// 		return parsed, nil
// 	}

// 	return parsed, nil
// }

// func parseSimpleTypeString(typeString string) (simpleTypeTree, error) {
// 	typeTree, remaining, err := parseSimpleTypeStringAllowingRemaining(typeString)
// 	if err != nil {
// 		return typeTree, err
// 	}

// 	if remaining != "" {
// 		return typeTree, fmt.Errorf("unexpected trailing '%s' in type string", remaining)
// 	}

// 	return typeTree, nil
// }

// func parseSimpleTypeStringAllowingRemaining(typeString string) (typeTree simpleTypeTree, remaining string, err error) {
// 	parser := typeParser{remainingInput: typeString}
// 	parsed, err := parser.parseTypeString()
// 	if err != nil {
// 		return parsed, "", err
// 	}
// 	if parsed.Name == "" {
// 		return parsed, "", errors.New("the type name cannot be empty")
// 	}
// 	parser.skipWhitespace()

// 	return parsed, parser.remainingInput, nil
// }

// func (tree simpleTypeTree) ToType(node NodeMeta) Type {
// 	nodeWithPositionUpdated := node
// 	nodeWithPositionUpdated.Column += tree.PositionOffset

// 	simpleType := SimpleType{NodeMeta: nodeWithPositionUpdated, Name: tree.Name}
// 	for _, typeArg := range tree.TypeArguments {
// 		simpleType.TypeArguments = append(simpleType.TypeArguments, typeArg.ToType(node))
// 	}

// 	if tree.Optional || tree.Dimensionality != nil {
// 		gt := &GeneralizedType{NodeMeta: nodeWithPositionUpdated}
// 		if tree.Optional {
// 			gt.Cases = TypeCases{
// 				&TypeCase{NodeMeta: node},
// 				&TypeCase{NodeMeta: node, Type: &simpleType},
// 			}
// 		} else {
// 			gt.Cases = TypeCases{
// 				&TypeCase{NodeMeta: node, Type: &simpleType},
// 			}
// 		}

// 		switch d := tree.Dimensionality.(type) {
// 		case nil:
// 		case *mapTypeTreeDimensionality:
// 			gt.Dimensionality = &Map{NodeMeta: node, KeyType: d.Key.ToType(node)}
// 		default:
// 			panic(fmt.Sprintf("unknown dimensionality type %T", d))
// 		}

// 		return gt
// 	}

// 	return &simpleType
// }

// type typeAst interface {
// 	_typeAst()
// }

// type nameAst struct {
// 	Name     string `@Ident`
// 	TypeArgs []typeAst
// }

// func (nameAst) _typeAst() {}

// type optionalAst struct {
// 	Element typeAst `@@ '?'`
// }

// func (optionalAst) _typeAst() {}

// type mapAst struct {
// 	Key   typeAst `@@ '-' '>'`
// 	Value typeAst `@@`
// }

// func (mapAst) _typeAst() {}

// type vectorAst struct {
// 	Element typeAst `@@ '*'`
// 	Length  *int    `@Int?`
// }

// func (vectorAst) _typeAst() {}

// var (
// 	expressionParser = participle.MustBuild[typeAst](
// 		participle.Union[typeAst](
// 			optionalAst{},
// 			mapAst{},
// 			nameAst{},
// 			vectorAst{}))
// )

type typeAst struct {
	Pos   lexer.Position
	Named *nameAst   `(@@`
	Sub   *typeAst   ` | '(' @@ ')')`
	Tails []typeTail `@@*`
}

func (t typeAst) String() string {
	var val string
	if t.Named != nil {
		val = t.Named.String()
	} else {
		val = t.Sub.String()
	}
	for _, tail := range t.Tails {
		val = tail.String(val)
	}

	return val
}

type nameAst struct {
	Pos      lexer.Position
	Name     string    `@Ident`
	TypeArgs []typeAst `('<' @@ (',' @@)* '>')?`
}

func (n nameAst) String() string {
	if len(n.TypeArgs) == 0 {
		return fmt.Sprintf("'%s'", n.Name)
	}

	var args []string
	for _, arg := range n.TypeArgs {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("(Generic '%s' %s)", n.Name, strings.Join(args, " "))
}

type typeTail struct {
	Optional *string    `  @'?'`
	MapValue *typeAst   `| '-' '>' @@`
	Vector   *vectorAst `| '*' @@`
}

func (t typeTail) String(target string) string {
	if t.Optional != nil {
		return fmt.Sprintf("(Optional %s)", target)
	}
	if t.MapValue != nil {
		return fmt.Sprintf("(Map %s %s)", target, t.MapValue)
	}
	if t.Vector != nil {
		if t.Vector.Length != nil {
			return fmt.Sprintf("(Vector[%d] %s)", *t.Vector.Length, target)
		}
		return fmt.Sprintf("(Vector %s)", target)
	}

	panic("unreachable")
}

type vectorAst struct {
	Length *uint64 `@Int?`
}

type patternAst struct {
	Discard     *string         `  @'_'`
	Declaration *declarationAst `| @@`
}

type declarationAst struct {
	Type     *typeAst `@@`
	Variable string   `@Ident?`
}

var (
	typeParser    = participle.MustBuild[typeAst]()
	patternParser = participle.MustBuild[patternAst]()
)

func parseSimpleType2(input string) (*typeAst, error) {
	return typeParser.ParseString("", input)
}

func parsePattern(input string, node NodeMeta) (Pattern, error) {
	pat, err := patternParser.ParseString("", input)
	if err != nil {
		return nil, err
	}
	if pat.Discard != nil {
		return &DiscardPattern{NodeMeta: node}, nil
	}
	if pat.Declaration != nil {
		t := pat.Declaration.Type.ToType(node)
		tp := TypePattern{NodeMeta: node, Type: t}
		if pat.Declaration.Variable != "" {
			return &DeclarationPattern{TypePattern: tp, Identifier: pat.Declaration.Variable}, nil
		}

		return &tp, nil
	}

	panic("unreachable")
}

func (ast typeAst) ToType(node NodeMeta) Type {
	nodeWithPositionUpdated := node
	nodeWithPositionUpdated.Column += ast.Pos.Offset

	var t Type
	if ast.Named != nil {
		simpleType := SimpleType{NodeMeta: nodeWithPositionUpdated, Name: ast.Named.Name}
		for _, typeArg := range ast.Named.TypeArgs {
			simpleType.TypeArguments = append(simpleType.TypeArguments, typeArg.ToType(node))
		}
		t = &simpleType
	} else {
		t = ast.Sub.ToType(node)
	}

	for _, tt := range ast.Tails {
		t = applyTail(t, tt)
	}

	return t
}

func applyTail(inner Type, tail typeTail) Type {
	nodeMeta := *inner.GetNodeMeta()
	gt := GeneralizedType{
		NodeMeta: nodeMeta,
		Cases:    TypeCases{&TypeCase{NodeMeta: nodeMeta, Type: inner}},
	}

	if tail.Optional != nil {
		gt.Cases = append(TypeCases{&TypeCase{NodeMeta: nodeMeta}}, gt.Cases...)
	} else if tail.MapValue != nil {
		gt.Cases = TypeCases{&TypeCase{NodeMeta: nodeMeta, Type: tail.MapValue.ToType(nodeMeta)}}
		gt.Dimensionality = &Map{
			NodeMeta: nodeMeta,
			KeyType:  inner,
		}
	} else if tail.Vector != nil {
		gt.Dimensionality = &Vector{
			NodeMeta: nodeMeta,
			Length:   tail.Vector.Length,
		}
	} else {
		panic("unreachable")
	}

	return &gt
}
