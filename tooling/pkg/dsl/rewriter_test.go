package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
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
	newEnv := Rewrite(env, func(self Rewriter, node Node) Node {
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

	nsName := newEnv.Namespaces[0].Name
	qualifiedRecName := nsName + ".Rec"
	l := len(newEnv.SymbolTable[qualifiedRecName].(*RecordDefinition).Fields)
	assert.Equal(t, 2, l)

	l2 := len(newEnv.Namespaces[0].TypeDefinitions[1].(*RecordDefinition).Fields[0].Type.(*SimpleType).ResolvedDefinition.(*RecordDefinition).Fields)
	assert.Equal(t, 2, l2)
}
