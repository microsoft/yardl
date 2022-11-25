// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/microsoft/yardl/tooling/pkg/dsl/expressions"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"gopkg.in/yaml.v3"
)

func ParseYamlInDir(path string, namespaceName string) (*Namespace, error) {
	errorSink := validation.ErrorSink{}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var paths []string

	if fileInfo.IsDir() {
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() &&
					(strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) &&
					info.Name() != packaging.PackageFileName {
					paths = append(paths, path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}

		sort.Slice(paths, func(i, j int) bool { return paths[i] < paths[j] })

	} else {
		paths = []string{path}
	}

	combinedNamespace := &Namespace{Name: namespaceName}

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		d := yaml.NewDecoder(f)
		d.KnownFields(true)

		ns := Namespace{Name: namespaceName}
		for {
			if err := d.Decode(&ns); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				errorSink.Add(validation.NewValidationError(err, path))
				break
			}
		}

		if len(errorSink.Errors) == 0 {
			// augment nodes with file path information
			Visit(&ns, func(self Visitor, node Node) {
				switch node := node.(type) {
				case *Namespace:
				default:
					nodeMeta := node.GetNodeMeta()
					nodeMeta.File = path
					if nodeMeta.Line == 0 || nodeMeta.Column == 0 {
						log.Panicf("node %T is missing line/column information. Line: %d, Column: %d", node, nodeMeta.Line, nodeMeta.Column)
					}
				}

				self.VisitChildren(node)
			})
		}

		combinedNamespace.TypeDefinitions = append(combinedNamespace.TypeDefinitions, ns.TypeDefinitions...)
		combinedNamespace.Protocols = append(combinedNamespace.Protocols, ns.Protocols...)
	}

	return combinedNamespace, errorSink.AsError()
}

func (meta *DefinitionMeta) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!str" {
		return parseError(value, "the name of a type is required to be a string")
	}

	meta.NodeMeta = createNodeMeta(value)
	meta.Name = value.Value

	parsedTypeString, err := parseSimpleTypeString(meta.Name)
	if err != nil {
		return parseError(value, err.Error())
	}

	meta.Name = parsedTypeString.Name
	if len(parsedTypeString.TypeArguments) > 0 {
		meta.TypeParameters = make([]*GenericTypeParameter, len(parsedTypeString.TypeArguments))
		for i, arg := range parsedTypeString.TypeArguments {
			if len(arg.TypeArguments) > 0 {
				return parseError(value, "generic type parameters cannot themselves have generic type parameters")
			}

			meta.TypeParameters[i] =
				&GenericTypeParameter{
					NodeMeta: createNodeMeta(value),
					Name:     arg.Name,
				}
		}
	}

	meta.Comment = normalizeComment(value.HeadComment)

	return nil
}

func (rec *RecordDefinition) UnmarshalYAML(value *yaml.Node) error {
	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "fields":
			err := v.DecodeWithOptions(&rec.Fields, yaml.DecodeOptions{KnownFields: true})
			if err != nil {
				return err
			}
		case "computedFields":
			err := v.DecodeWithOptions(&rec.ComputedFields, yaml.DecodeOptions{KnownFields: true})
			if err != nil {
				return err
			}
		default:
			return parseError(k, "field '%s' is not valid on a !record specification", k.Value)
		}
	}

	return nil
}

func (ns *Namespace) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!map" {
		return parseError(value, "expected a mapping from <typename>: <type definition>")
	}

	for i := 0; i < len(value.Content); i += 2 {
		nameNode := value.Content[i]
		typeNode := value.Content[i+1]

		meta := &DefinitionMeta{}
		if err := nameNode.DecodeWithOptions(&meta, yaml.DecodeOptions{KnownFields: true}); err != nil {
			return err
		}

		typeDef, err := UnmarshalTypeDefinition(typeNode, meta)
		if err != nil {
			return err
		}

		if protocol, ok := typeDef.(*ProtocolDefinition); ok {
			ns.Protocols = append(ns.Protocols, protocol)
		} else {
			ns.TypeDefinitions = append(ns.TypeDefinitions, typeDef)
		}
	}

	return nil
}

func normalizeComment(comments string) string {
	if comments == "" {
		return comments
	}

	lines := strings.Split(comments, "\n")
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "# ")
	}
	return strings.Join(lines, "\n")
}

func (fields *Fields) UnmarshalYAML(value *yaml.Node) error {
	unpacked := []*Field(*fields)
	err := UnmarshalFieldsOrProtocolStepsYAML(&unpacked, value)
	*fields = Fields(unpacked)
	return err
}

func (computedFields *ComputedFields) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!map" {
		return parseError(value, "expected computed fields to be a map")
	}

	for i := 0; i < len(value.Content); i += 2 {
		fieldKey := value.Content[i]
		fieldValue := value.Content[i+1]
		if fieldKey.Tag != "!!str" {
			return parseError(value, "expected computed field name to be a string")
		}

		fieldName := fieldKey.Value

		expression, err := UnmarshalExpression(fieldValue)
		if err != nil {
			return err
		}

		field := &ComputedField{
			Name:       fieldName,
			Comment:    normalizeComment(fieldKey.HeadComment),
			NodeMeta:   createNodeMeta(fieldKey),
			Expression: expression,
		}

		*computedFields = append(*computedFields, field)
	}

	return nil
}

func UnmarshalExpression(value *yaml.Node) (Expression, error) {
	switch value.Tag {
	case "!!null":
		return nil, parseError(value, "An expression cannot be null")
	case "!!str", "!!int", "!!float", "!!bool", "!switch":
		exp, err := expressions.ParseExpression(value.Value)
		if err != nil {
			if err, ok := err.(expressions.ParseError); ok {
				line := value.Line + err.Position().Line - 1
				column := value.Column + err.Position().Column - 1
				return nil, validation.ValidationError{
					Message: errors.New(err.Message()),
					Line:    &line,
					Column:  &column,
				}
			}

			panic("not parse error")
		}

		return ConvertExpression(*exp, value), nil
	case "!!map":
		if len(value.Content) == 2 {
			key := value.Content[0]
			value := value.Content[1]
			if key.Tag == "!switch" {
				if value.Tag != "!!map" {
					return nil, parseError(value, "expected a mapping from <case>: <expression>")
				}
				return UnmarshalSwitchExpression(key, value.Content)
			}
		}
		return nil, parseError(value, "expected a !switch expression")
	default:
		return nil, parseError(value, "unsupported expression type: %s", value.Tag)
	}
}

func UnmarshalSwitchExpression(targetNode *yaml.Node, caseNodes []*yaml.Node) (Expression, error) {
	switchExpression := SwitchExpression{
		NodeMeta: createNodeMeta(targetNode),
		Cases:    make([]*SwitchCase, len(caseNodes)/2),
	}

	target, err := UnmarshalExpression(targetNode)
	if err != nil {
		return nil, err
	}

	switchExpression.Target = target

	for i := 0; i < len(caseNodes); i += 2 {
		patternNode := caseNodes[i]
		exprNode := caseNodes[i+1]
		switchCase := &SwitchCase{
			NodeMeta: createNodeMeta(patternNode),
		}

		pattern, err := UnmarshalPattern(patternNode)
		if err != nil {
			return nil, err
		}
		switchCase.Pattern = pattern

		expr, err := UnmarshalExpression(exprNode)
		if err != nil {
			return nil, err
		}
		switchCase.Expression = expr

		switchExpression.Cases[i/2] = switchCase
	}

	return &switchExpression, nil
}

func UnmarshalPattern(patternNode *yaml.Node) (Pattern, error) {
	switch patternNode.Tag {
	case "!!null":
		return &TypePattern{
			NodeMeta: createNodeMeta(patternNode),
		}, nil

	case "!!str":
		if patternNode.Value == "_" {
			return &DiscardPattern{
				NodeMeta: createNodeMeta(patternNode),
			}, nil
		}

		parsedTypeTree, remaining, err := parseSimpleTypeStringAllowingRemaining(patternNode.Value)
		if err != nil {
			return nil, parseError(patternNode, err.Error())
		}

		var typePatternType Type

		// check for (invalid) declaration pattern:
		// null <var>:
		// We only accept
		// null:
		// but we don't want an error saying that the type is invalid, rather the error should
		// say that a variable of type null is not allowed.
		if parsedTypeTree.Name != "null" || parsedTypeTree.Optional || len(parsedTypeTree.TypeArguments) != 0 {
			typePatternType = parsedTypeTree.ToType(createNodeMeta(patternNode))
		}

		typePattern := TypePattern{
			NodeMeta: createNodeMeta(patternNode),
			Type:     typePatternType,
		}

		if remaining == "" {
			return &typePattern, nil
		}

		if len(strings.Fields(remaining)) == 1 {
			return &DeclarationPattern{
				TypePattern: typePattern,
				Identifier:  remaining,
			}, nil
		}

		return nil, parseError(patternNode, "unable to parse pattern. Expected a type name, or a type name and an identifier, or a discard `_`")
	default:
		return nil, parseError(patternNode, "expected pattern to be a string")
	}
}

func ConvertExpression(expression expressions.Expression, hostNode *yaml.Node) Expression {
	createNodeMeta := func(pos lexer.Position) NodeMeta {
		return NodeMeta{
			Line:   hostNode.Line + pos.Line - 1,
			Column: hostNode.Column + pos.Column - 1,
		}
	}

	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		exp := &IntegerLiteralExpression{
			NodeMeta: createNodeMeta(expression.Pos),
		}

		exp.Value.UnmarshalText([]byte(expression.Value))
		return exp

	case expressions.StringLiteral:
		return &StringLiteralExpression{
			NodeMeta: createNodeMeta(expression.Pos),
			Value:    expression.Value,
		}
	case expressions.PathExpr:
		var target Expression
		for i, part := range expression.Parts {
			curr := &MemberAccessExpression{
				NodeMeta: createNodeMeta(part.Pos),
				Member:   part.Name,
			}

			if i > 0 {
				curr.Target = target
			}
			target = curr

			for _, indexAst := range part.Indexes {
				index := &IndexExpression{
					NodeMeta: createNodeMeta(part.Pos),
					Target:   target,
				}

				for _, arg := range indexAst.IndexArgs {
					convertedArg := &IndexArgument{
						NodeMeta: createNodeMeta(arg.Pos),
						Value:    ConvertExpression(arg.Value, hostNode),
					}

					if arg.Label != nil {
						convertedArg.Label = arg.Label.Name
						convertedArg.NodeMeta = createNodeMeta(arg.Label.Pos)
					} else {
						convertedArg.NodeMeta = *convertedArg.Value.GetNodeMeta()
					}

					index.Arguments = append(index.Arguments, convertedArg)
				}

				target = index
			}
		}

		return target

	case expressions.FunctionCall:
		args := make([]Expression, len(expression.Arguments))
		for i, arg := range expression.Arguments {
			args[i] = ConvertExpression(arg, hostNode)
		}

		return &FunctionCallExpression{
			NodeMeta:     createNodeMeta(expression.Pos),
			FunctionName: expression.FunctionName,
			Arguments:    args,
		}
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", expression))
	}
}

func (protocol *ProtocolDefinition) UnmarshalYAML(value *yaml.Node) error {
	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "sequence":
			err := v.DecodeWithOptions(&protocol.Sequence, yaml.DecodeOptions{KnownFields: true})
			if err != nil {
				return err
			}
		default:
			return parseError(k, "field '%s' is not valid on a !protocol specification", k.Value)
		}
	}

	return nil
}

func (steps *ProtocolSteps) UnmarshalYAML(value *yaml.Node) error {
	unpacked := []*ProtocolStep(*steps)
	err := UnmarshalFieldsOrProtocolStepsYAML(&unpacked, value)
	*steps = ProtocolSteps(unpacked)
	return err
}

func UnmarshalFieldsOrProtocolStepsYAML[T fieldOrProtocolStep](elements *[]*T, value *yaml.Node) error {
	if value.Tag != "!!map" {
		return parseError(value, "expected field map")
	}

	for i := 0; i < len(value.Content); i += 2 {
		fieldKey := value.Content[i]
		fieldValue := value.Content[i+1]
		if fieldKey.Tag != "!!str" {
			return parseError(value, "expected field name to be a string")
		}
		if fieldValue.Tag == "!!null" {
			return parseError(value, "a field or protocol step cannot be null")
		}

		fieldName := fieldKey.Value

		t, err := UnmarshalTypeYAML(fieldValue)
		if err != nil {
			return err
		}

		e := &T{
			Name:     fieldName,
			Comment:  normalizeComment(fieldKey.HeadComment),
			Type:     t,
			NodeMeta: createNodeMeta(fieldKey),
		}
		*elements = append(*elements, e)
	}

	return nil
}

type fieldOrProtocolStep interface {
	Field | ProtocolStep
}

func UnmarshalVectorYAML(value *yaml.Node) (*GeneralizedType, error) {
	if value.Kind != yaml.MappingNode {
		return nil, parseError(value, "a !vector must be specified with field `items` and optionally `length`")
	}

	vector := &Vector{NodeMeta: createNodeMeta(value)}
	t := &GeneralizedType{Dimensionality: vector, NodeMeta: vector.NodeMeta}

	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "items":
			cases, err := UnmarshalTypeCases(v)
			if err != nil {
				return nil, err
			}
			t.Cases = cases
		case "length":
			var length big.Int
			if err := length.UnmarshalText([]byte(v.Value)); err != nil {
				return nil, err
			}
			if length.Sign() < 0 {
				return nil, parseError(v, "vector length cannot be negative")
			}
			asUint64 := length.Uint64()
			vector.Length = &asUint64
		default:
			return nil, parseError(k, "field '%s' is not valid on a !vector specification", k.Value)
		}
	}

	if t.Cases == nil {
		return nil, parseError(value, "`items` must be specified on a !vector")
	}

	return t, nil
}

func UnmarshalArrayYAML(value *yaml.Node) (*GeneralizedType, error) {
	if value.Kind != yaml.MappingNode {
		return nil, parseError(value, "an !array must be specified with field `items` and optionally `dimensions`")
	}

	array := &Array{NodeMeta: createNodeMeta(value)}
	nt := &GeneralizedType{Dimensionality: array, NodeMeta: createNodeMeta(value)}

	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "items":
			cases, err := UnmarshalTypeCases(v)
			if err != nil {
				return nil, err
			}
			nt.Cases = cases
		case "dimensions":
			switch v.Tag {
			case "!!null":
			case "!!int":
				var ndims int

				if err := v.DecodeWithOptions(&ndims, yaml.DecodeOptions{KnownFields: true}); err != nil {
					return nil, err
				}

				dims := make(ArrayDimensions, ndims)
				for i := range dims {
					dims[i] = &ArrayDimension{NodeMeta: createNodeMeta(v)}
				}
				array.Dimensions = &dims
			case "!!seq":
				array.Dimensions = &ArrayDimensions{}
				for i := 0; i < len(v.Content); i++ {
					dim := &ArrayDimension{Comment: normalizeComment(v.Content[i].HeadComment), NodeMeta: createNodeMeta(v.Content[i])}
					if v.Content[i].Tag != "!!null" {
						if err := v.Content[i].DecodeWithOptions(&dim, yaml.DecodeOptions{KnownFields: true}); err != nil {
							return nil, err
						}
					}
					*array.Dimensions = append(*array.Dimensions, dim)
				}
			case "!!map":
				array.Dimensions = &ArrayDimensions{}
				for i := 0; i < len(v.Content); i += 2 {
					k := v.Content[i]
					v := v.Content[i+1]
					dim := ArrayDimension{Name: &k.Value, Comment: normalizeComment(k.HeadComment), NodeMeta: createNodeMeta(k)}
					if err := v.DecodeWithOptions(&dim, yaml.DecodeOptions{KnownFields: true}); err != nil {
						return nil, err
					}
					*array.Dimensions = append(*array.Dimensions, &dim)
				}

			default:
				return nil, parseError(v, "dimensions must be specified as a list of dimension specifications or the number of dimensions")
			}
		default:
			return nil, parseError(k, "field '%s' is not valid on an !array specification", k.Value)
		}
	}

	if nt.Cases == nil {
		return nil, parseError(value, "items must be specified on an !array")
	}

	return nt, nil
}

func UnmarshalStreamYAML(value *yaml.Node) (*GeneralizedType, error) {
	if value.Kind != yaml.MappingNode {
		return nil, parseError(value, "a !stream must be specified with field `items`")
	}

	nodeMeta := createNodeMeta(value)

	t := &GeneralizedType{
		Dimensionality: &Stream{NodeMeta: nodeMeta},
		NodeMeta:       nodeMeta,
	}

	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "items":
			cases, err := UnmarshalTypeCases(v)
			if err != nil {
				return nil, err
			}
			t.Cases = cases
		default:
			return nil, parseError(k, "field '%s' is not valid on a !stream specification", k.Value)
		}
	}

	if t.Cases == nil {
		return nil, parseError(value, "'items' must be specified on an !stream")
	}

	return t, nil
}

func UnmarshalTypeDefinition(value *yaml.Node, definitionMeta *DefinitionMeta) (TypeDefinition, error) {
	switch value.Tag {
	case "!record":
		rec := &RecordDefinition{DefinitionMeta: definitionMeta}
		err := value.DecodeWithOptions(rec, yaml.DecodeOptions{KnownFields: true})
		return rec, err
	case "!enum":
		enum := &EnumDefinition{DefinitionMeta: definitionMeta}
		err := value.DecodeWithOptions(enum, yaml.DecodeOptions{KnownFields: true})
		return enum, err
	case "!protocol":
		protocol := &ProtocolDefinition{DefinitionMeta: definitionMeta}
		err := value.DecodeWithOptions(protocol, yaml.DecodeOptions{KnownFields: true})
		return protocol, err
	default:
		namedType := &NamedType{DefinitionMeta: definitionMeta}
		underlyingType, err := UnmarshalTypeYAML(value)
		if err == nil && underlyingType == nil {
			return nil, parseError(value, "type cannot be empty")
		}
		namedType.Type = underlyingType
		return namedType, err
	}
}

func UnmarshalTypeYAML(value *yaml.Node) (Type, error) {
	switch value.Tag {
	case "!!null":
		return nil, nil
	case "!!str":
		parsedTypeTree, err := parseSimpleTypeString(value.Value)
		if err != nil {
			return nil, parseError(value, err.Error())
		}

		return parsedTypeTree.ToType(createNodeMeta(value)), nil
	case "!generic":
		return UnmarshalGenericNode(value)
	case "!!seq":
		cases, err := UnmarshalTypeCases(value)
		return &GeneralizedType{NodeMeta: createNodeMeta(value), Cases: cases, Dimensionality: nil}, err
	case "!vector":
		return UnmarshalVectorYAML(value)
	case "!array":
		return UnmarshalArrayYAML(value)
	case "!stream":
		return UnmarshalStreamYAML(value)
	default:
		return nil, parseError(value, "unrecognized type kind '%s'", value.Tag)
	}
}

func UnmarshalTypeCases(value *yaml.Node) (TypeCases, error) {
	cases := TypeCases{}
	switch value.Kind {
	case yaml.SequenceNode:
		for _, c := range value.Content {
			t, err := UnmarshalTypeYAML(c)
			if err != nil {
				return nil, err
			}

			cases = append(cases, &TypeCase{Type: t, NodeMeta: createNodeMeta(c)})
		}
	default:
		t, err := UnmarshalTypeYAML(value)
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, parseError(value, "type null is only supported in unions")
		}

		cases = append(cases, &TypeCase{Type: t, NodeMeta: createNodeMeta(value)})
	}

	return cases, nil
}

func UnmarshalGenericNode(value *yaml.Node) (Type, error) {
	simpleType := &SimpleType{NodeMeta: createNodeMeta(value)}

	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "name":
			simpleType.Name = v.Value
		case "args":
			simpleType.TypeArguments = make([]Type, 0)
			if v.Kind != yaml.SequenceNode {
				typeArg, err := UnmarshalTypeYAML(v)
				if err != nil {
					return nil, err
				}

				simpleType.TypeArguments = append(simpleType.TypeArguments, typeArg)
			} else {
				for _, c := range v.Content {
					typeArg, err := UnmarshalTypeYAML(c)
					if err != nil {
						return nil, err
					}

					simpleType.TypeArguments = append(simpleType.TypeArguments, typeArg)
				}
			}
		default:
			return nil, parseError(k, "field '%s' is not valid on an !generic specification ('name' and 'args' are expected)", k.Value)
		}
	}

	if simpleType.Name == "" {
		return nil, parseError(value, "the 'name' property of a !generic type must be specified and non-empty")
	}

	if simpleType.TypeArguments == nil {
		return nil, parseError(value, "the 'args' property of a !generic type must be specified")
	}

	return simpleType, nil
}

func (dimension *ArrayDimension) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!int" {
		var length big.Int
		if err := length.UnmarshalText([]byte(value.Value)); err != nil {
			return err
		}
		if length.Sign() < 0 {
			return parseError(value, "array dimension length cannot be negative")
		}
		asUnit64 := length.Uint64()
		dimension.Length = &asUnit64
		return nil
	}

	return parseError(value, "invalid dimension specification")
}

func (enum *EnumDefinition) UnmarshalYAML(value *yaml.Node) error {
	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "base":
			base, err := UnmarshalTypeYAML(v)
			if err != nil {
				return err
			}
			enum.BaseType = base
		case "values":
			err := v.DecodeWithOptions(&enum.Values, yaml.DecodeOptions{KnownFields: true})
			if err != nil {
				return err
			}
		default:
			return parseError(k, "field '%s' is not valid on an !enum specification", k.Value)
		}
	}

	return nil
}

func (evs *EnumValues) UnmarshalYAML(value *yaml.Node) error {
	switch value.Tag {
	case "!!seq":
		for i, v := range value.Content {
			if v.Tag != "!!str" {
				goto err
			}
			ev := &EnumValue{
				NodeMeta:     createNodeMeta(v),
				Comment:      normalizeComment(v.HeadComment),
				Symbol:       v.Value,
				IntegerValue: *big.NewInt(int64(i)),
			}
			*evs = append(*evs, ev)
		}

		return nil

	case "!!map":
		for i := 0; i < len(value.Content); i += 2 {
			k := value.Content[i]
			v := value.Content[i+1]
			if k.Tag != "!!str" && v.Tag != "!!int" {
				goto err
			}

			if v.Kind != yaml.ScalarNode {
				return parseError(v, "enum value must be an integer")
			}

			val := &EnumValue{
				NodeMeta: createNodeMeta(k),
				Comment:  normalizeComment(k.HeadComment),
				Symbol:   k.Value,
			}

			if err := val.IntegerValue.UnmarshalText([]byte(v.Value)); err != nil {
				return err
			}

			*evs = append(*evs, val)
		}
		return nil
	}
err:
	return parseError(value, "invalid enum specification")
}

func parseError(node *yaml.Node, message string, args ...any) validation.ValidationError {
	return validation.ValidationError{
		Message: fmt.Errorf(message, args...),
		Line:    &node.Line,
		Column:  &node.Column,
	}
}

func createNodeMeta(yamlNode *yaml.Node) NodeMeta {
	return NodeMeta{Line: yamlNode.Line, Column: yamlNode.Column}
}
