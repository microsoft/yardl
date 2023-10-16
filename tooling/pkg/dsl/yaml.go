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

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/microsoft/yardl/tooling/pkg/dsl/parser"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"gopkg.in/yaml.v3"
)

func ParsePackageContents(pkgInfo packaging.PackageInfo) (*Namespace, error) {
	return ParseYamlInDir(pkgInfo.PackageDir(), pkgInfo.Namespace)
}

// Parses all model YAML files, combining them into a single Namespace
// path can be a single YAML file or a directory containing YAML files
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
				case *DefinitionMeta:
					for _, p := range node.TypeParameters {
						self.Visit(p)
					}
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

	parsedTypeString, err := parser.ParseType(meta.Name)
	if err != nil {
		return parseError(value, err.Error())
	}

	parsedType := convertType(parsedTypeString, meta.NodeMeta)
	simpleParsedType, ok := parsedType.(*SimpleType)
	if !ok {
		return parseError(value, "not a valid type declaration name")
	}

	meta.Name = simpleParsedType.Name

	for _, t := range simpleParsedType.TypeArguments {
		sa, ok := t.(*SimpleType)
		if !ok {
			return parseError(value, "invalid type parameter name")
		}

		if len(sa.TypeArguments) > 0 {
			return parseError(value, "generic type parameters cannot themselves have generic type parameters")
		}

		meta.TypeParameters = append(meta.TypeParameters, &GenericTypeParameter{NodeMeta: meta.NodeMeta, Name: sa.Name})
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

	// There cannot be an empty line beteen a comment and the Yardl
	// element for it to be considered a documentation comment.
	i := len(lines) - 1
	for ; i >= 0; i-- {
		if !strings.HasPrefix(lines[i], "#") {
			break
		}
	}

	lines = lines[i+1:]

	for i := range lines {
		lines[i] = strings.TrimPrefix(strings.TrimPrefix(lines[i], "# "), "#")
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
		return ParseExpression(value.Value, value.Line, value.Column)
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
		patAst, err := parser.ParsePattern(patternNode.Value)
		if err != nil {
			return nil, parseError(patternNode, err.Error())
		}

		pat := convertPattern(patAst, createNodeMeta(patternNode))

		// check for (invalid) declaration pattern:
		// null <var>:
		// We only accept
		// null:
		// but we don't want an error saying that the type is invalid, rather the error should
		// say that a variable of type null is not allowed.
		if dp, ok := pat.(*DeclarationPattern); ok {
			if st, ok := dp.TypePattern.Type.(*SimpleType); ok && st.Name == "null" {
				dp.Type = nil
			}
		}

		return pat, nil
	default:
		return nil, parseError(patternNode, "expected pattern to be a string")
	}
}

func convertType(ast *parser.Type, node NodeMeta) Type {
	nodeWithPositionUpdated := node
	nodeWithPositionUpdated.Column += ast.Pos.Offset

	var t Type
	if ast.Named != nil {
		simpleType := SimpleType{NodeMeta: nodeWithPositionUpdated, Name: ast.Named.Name}
		for _, typeArg := range ast.Named.TypeArgs {
			simpleType.TypeArguments = append(simpleType.TypeArguments, convertType(typeArg, node))
		}
		t = &simpleType
	} else if ast.Sub != nil {
		t = convertType(ast.Sub, node)
	} else {
		panic("unreachable")
	}

	for _, tt := range ast.Tails {
		t = applyTypeTail(t, tt)
	}

	return t
}

func applyTypeTail(inner Type, tail parser.TypeTail) Type {
	nodeMeta := *inner.GetNodeMeta()
	gt := GeneralizedType{
		NodeMeta: nodeMeta,
		Cases:    TypeCases{&TypeCase{NodeMeta: nodeMeta, Type: inner}},
	}

	if tail.Optional {
		gt.Cases = append(TypeCases{&TypeCase{NodeMeta: nodeMeta}}, gt.Cases...)
	} else if tail.MapValue != nil {
		gt.Cases = TypeCases{&TypeCase{NodeMeta: nodeMeta, Type: convertType(tail.MapValue, nodeMeta)}}
		gt.Dimensionality = &Map{
			NodeMeta: nodeMeta,
			KeyType:  inner,
		}
	} else if tail.Vector != nil {
		gt.Dimensionality = &Vector{
			NodeMeta: nodeMeta,
			Length:   tail.Vector.Length,
		}
	} else if tail.Array != nil {
		a := Array{NodeMeta: nodeMeta}

		if len(tail.Array.Dimensions) > 0 {
			dims := ArrayDimensions{}
			for _, dim := range tail.Array.Dimensions {
				dims = append(dims, &ArrayDimension{NodeMeta: nodeMeta, Name: dim.Name, Length: dim.Length})
			}

			a.Dimensions = &dims
		}

		gt.Dimensionality = &a
	} else {
		panic("unreachable")
	}

	return &gt
}

func convertPattern(pat *parser.Pattern, node NodeMeta) Pattern {
	if pat.Discard {
		return &DiscardPattern{NodeMeta: node}
	}
	if pat.Type != nil {
		tp := TypePattern{NodeMeta: node, Type: convertType(pat.Type, node)}
		if pat.Variable != nil {
			return &DeclarationPattern{TypePattern: tp, Identifier: *pat.Variable}
		}

		return &tp
	}

	panic("unreachable")
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

func UnmarshalMapYAML(value *yaml.Node) (*GeneralizedType, error) {
	if value.Kind != yaml.MappingNode {
		return nil, parseError(value, "a !map must be specified with fields `keys` and `values`")
	}

	m := &Map{NodeMeta: createNodeMeta(value)}
	t := &GeneralizedType{Dimensionality: m, NodeMeta: m.NodeMeta}

	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		switch k.Value {
		case "keys":
			keyType, err := UnmarshalTypeYAML(v)
			if err != nil {
				return nil, err
			}
			m.KeyType = keyType
		case "values":
			cases, err := UnmarshalTypeCases(v)
			if err != nil {
				return nil, err
			}
			t.Cases = cases
		default:
			return nil, parseError(k, "field '%s' is not valid on a !map specification", k.Value)
		}
	}

	if m.KeyType == nil {
		return nil, parseError(value, "`keys` must be specified on a !map")
	}

	if t.Cases == nil {
		return nil, parseError(value, "`values` must be specified on a !map")
	}

	return t, nil
}

func UnmarshalTypeDefinition(value *yaml.Node, definitionMeta *DefinitionMeta) (TypeDefinition, error) {
	switch value.Tag {
	case "!record":
		rec := &RecordDefinition{DefinitionMeta: definitionMeta}
		err := value.DecodeWithOptions(rec, yaml.DecodeOptions{KnownFields: true})
		return rec, err
	case "!enum", "!flags":
		enum := &EnumDefinition{DefinitionMeta: definitionMeta, IsFlags: value.Tag == "!flags"}
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
		parsedTypeTree, err := parser.ParseType(value.Value)
		if err != nil {
			return nil, parseError(value, err.Error())
		}

		return convertType(parsedTypeTree, createNodeMeta(value)), nil
	case "!generic":
		return UnmarshalGenericNode(value)
	case "!!seq":
		cases, err := UnmarshalTypeCases(value)
		return &GeneralizedType{NodeMeta: createNodeMeta(value), Cases: cases, Dimensionality: nil}, err
	case "!vector":
		return UnmarshalVectorYAML(value)
	case "!array":
		return UnmarshalArrayYAML(value)
	case "!map":
		return UnmarshalMapYAML(value)
	case "!union":
		return UnmarshalUnionYAML(value)
	case "!stream":
		return UnmarshalStreamYAML(value)
	default:
		return nil, parseError(value, "unrecognized type kind '%s'", value.Tag)
	}
}

func UnmarshalUnionYAML(value *yaml.Node) (*GeneralizedType, error) {
	if value.Kind != yaml.MappingNode {
		return nil, parseError(value, "a !union must be specified as a map from tag to type")
	}

	cases := TypeCases{}
	for i := 0; i < len(value.Content); i += 2 {
		tagNode := value.Content[i]
		typeNode := value.Content[i+1]

		if tagNode.Tag != "!!str" && tagNode.Tag != "!!null" {
			return nil, parseError(tagNode, "tag must be a string")
		}

		tag := tagNode.Value
		parsedType, err := UnmarshalTypeYAML(typeNode)
		if err != nil {
			return nil, err
		}

		cases = append(cases, &TypeCase{Tag: tag, ExplicitTag: true, Type: parsedType, NodeMeta: createNodeMeta(tagNode)})
	}

	return &GeneralizedType{NodeMeta: createNodeMeta(value), Cases: cases, Dimensionality: nil}, nil
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
	if value.Tag == "!!str" {
		dimension.Name = &value.Value
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
			vals, err := UnmarshalEnumValues(enum.IsFlags, v)
			if err != nil {
				return err
			}

			enum.Values = *vals
		default:
			return parseError(k, "field '%s' is not valid on an !enum specification", k.Value)
		}
	}

	return nil
}

func UnmarshalEnumValues(flags bool, value *yaml.Node) (*EnumValues, error) {
	vals := EnumValues{}
	switch value.Tag {
	case "!!seq":
		for i, v := range value.Content {
			if v.Tag != "!!str" {
				goto err
			}
			var integerValue big.Int
			if flags {
				integerValue.SetBit(&integerValue, i, 1)
			} else {
				integerValue.SetInt64(int64(i))
			}

			ev := &EnumValue{
				NodeMeta:     createNodeMeta(v),
				Comment:      normalizeComment(v.HeadComment),
				Symbol:       v.Value,
				IntegerValue: integerValue,
			}
			vals = append(vals, ev)
		}

		return &vals, nil

	case "!!map":
		for i := 0; i < len(value.Content); i += 2 {
			k := value.Content[i]
			v := value.Content[i+1]
			if k.Tag != "!!str" && v.Tag != "!!int" {
				goto err
			}

			if v.Kind != yaml.ScalarNode {
				return nil, parseError(v, "enum or flag value must be an integer or empty")
			}

			val := &EnumValue{
				NodeMeta: createNodeMeta(k),
				Comment:  normalizeComment(k.HeadComment),
				Symbol:   k.Value,
			}

			if v.Value == "" {
				if flags {
					if i == 0 {
						val.IntegerValue.SetInt64(1)
					} else {
						prevVal := vals[i/2-1].IntegerValue
						if prevVal.Sign() < 0 {
							return nil, parseError(v, "flag value following a negative value must be explicitly specified")
						}
						newVal := big.NewInt(1)
						for ; newVal.Cmp(&prevVal) <= 0; newVal.Lsh(newVal, 1) {
						}
						val.IntegerValue = *newVal
					}
				} else {
					if i == 0 {
						val.IntegerValue.SetInt64(0)
					} else {
						prevVal := vals[i/2-1].IntegerValue
						var newVal big.Int
						if prevVal.Sign() < 0 {
							newVal.Sub(&prevVal, big.NewInt(1))
						} else {
							newVal.Add(&prevVal, big.NewInt(1))
						}
						val.IntegerValue = newVal
					}
				}
			} else {
				if err := val.IntegerValue.UnmarshalText([]byte(v.Value)); err != nil {
					return nil, err
				}
			}

			vals = append(vals, val)
		}
		return &vals, nil
	}
err:
	return nil, parseError(value, "invalid enum or flag specification")
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
