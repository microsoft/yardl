package dsl

import (
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
