package dsl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestSymbolTableUpdatedWhenFieldChanged(t *testing.T) {
	src := `
Rec: !record
  fields:
    f: int
`
	env, err := parseAndValidate(t, src)
	assert.NoError(t, err)
	newEnv := Rewrite(env, func(self *Rewriter, node Node) Node {
		switch node := node.(type) {
		case *SimpleType:
			return &SimpleType{
				NodeMeta:           node.NodeMeta,
				Name:               String,
				ResolvedDefinition: PrimitiveString,
			}
		}
		return self.DefaultRewrite(node)
	}).(*Environment)

	assert.Equal(t, PrimitiveString, newEnv.SymbolTable["test.Rec"].(*RecordDefinition).Fields[0].Type.(*SimpleType).ResolvedDefinition)
	assert.Equal(t, PrimitiveInt32, env.SymbolTable["test.Rec"].(*RecordDefinition).Fields[0].Type.(*SimpleType).ResolvedDefinition)
}

func TestResolvedReferencesUpdatedWhenTargetRewritten(t *testing.T) {
	src := `
Rec: !record
  fields:
    f: int

OtherRec: !record
  fields:
    r: Rec
`
	env, err := parseAndValidate(t, src)
	assert.NoError(t, err)
	newEnv := Rewrite(env, func(self *Rewriter, node Node) Node {
		switch node := node.(type) {
		case *RecordDefinition:

			if node.Name == "Rec" {
				newRec := *node
				newRec.Fields = append(newRec.Fields, &Field{Name: "g", Type: node.Fields[0].Type})
				return &newRec
			}
		}
		return self.DefaultRewrite(node)
	}).(*Environment)

	assert.Equal(t, 2, len(newEnv.Namespaces[0].TypeDefinitions[1].(*RecordDefinition).Fields[0].Type.(*SimpleType).ResolvedDefinition.(*RecordDefinition).Fields))
}

func TestFoo(t *testing.T) {
	src := `
Rec<T>: !record
  fields:
    arr: T[]

OtherRec: !record
  fields:
    r: Rec<float>
`
	env, err := parseAndValidate(t, src)
	assert.NoError(t, err)
	newEnv := RewriteWithContext(env, nil, func(node Node, context *TypeParameterBag, self *RewriterWithContext[*TypeParameterBag]) Node {
		switch node := node.(type) {
		case *ProtocolDefinition, *EnumDefinition:
			return node
		case *GenericTypeParameter:
			return context.GetParameter(node.Name, context.InArray)
		case TypeDefinition:
			meta := node.GetDefinitionMeta()
			if len(meta.TypeParameters) == 0 {
				return node
			}

			typeParameterBag := &TypeParameterBag{
				Parameters:     make(map[string]*GenericTypeParameter),
				UsedParameters: make(map[string]any),
			}
			for _, tp := range meta.TypeParameters {
				typeParameterBag.Parameters[tp.Name] = tp
			}

			rewritten := self.DefaultRewrite(node, typeParameterBag).(TypeDefinition)

			newMeta := *rewritten.GetDefinitionMeta()

			newTypeParameters := make([]*GenericTypeParameter, 0, len(meta.TypeParameters))
			for _, tp := range meta.TypeParameters {
				if _, ok := typeParameterBag.UsedParameters[tp.Name]; ok {
					newTypeParameters = append(newTypeParameters, tp)
				}

				if _, ok := typeParameterBag.UsedParameters[numPyTypeParameterName(tp.Name)]; ok {
					newTypeParameters = append(newTypeParameters, tp)
				}
			}
			newMeta.TypeParameters = newTypeParameters
			switch rewritten := rewritten.(type) {
			case *RecordDefinition:
				rewritten.DefinitionMeta = &newMeta
			case *NamedType:
				rewritten.DefinitionMeta = &newMeta
			default:
				panic(fmt.Sprintf("Unexpected type %T", rewritten))
			}

			return rewritten

		case *GeneralizedType:
			switch node.Dimensionality.(type) {
			case *Array:
				context.InArray = true
				newNode := self.DefaultRewrite(node, context)
				context.InArray = false
				return newNode
			}
		}
		return self.DefaultRewrite(node, context)
	}).(*Environment)

	for _, gtp := range newEnv.Namespaces[0].TypeDefinitions[0].(*RecordDefinition).DefinitionMeta.TypeParameters {
		assert.Equal(t, gtp.Name+"_NP", gtp.Name)
	}

	fmt.Println(newEnv)
}

type TypeParameterBag struct {
	UsedParameters map[string]any
	Parameters     map[string]*GenericTypeParameter
	InArray        bool
}

func (t *TypeParameterBag) GetParameter(name string, numpy bool) *GenericTypeParameter {
	normalParameter := t.Parameters[name]
	if normalParameter == nil {
		panic(fmt.Sprintf("Type parameter %s not found", name))
	}

	if !numpy {
		t.UsedParameters[name] = true
		return normalParameter
	}

	numpyName := numPyTypeParameterName(name)

	if p, ok := t.Parameters[numpyName]; ok {
		return p
	}

	p := *normalParameter
	p.Name = numpyName
	t.Parameters[numpyName] = &p
	t.UsedParameters[numpyName] = true
	return &p
}

func numPyTypeParameterName(typeParameterName string) string {
	return typeParameterName + "_NP"
}
