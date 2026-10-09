// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression test for https://github.com/microsoft/yardl/issues/298 — when two
// namespaces define different records with the same unqualified name, the
// embedded schema must still uniquely identify each definition so consumers
// that resolve by name (e.g. a dynamic/schema-only reader) can disambiguate.
func TestGetProtocolSchema_QualifiesImportedDefinitionNames(t *testing.T) {
	importedSrc := `
Foo: !record
  fields:
    a: int32`
	baseSrc := `
Foo: !record
  fields:
    b: string

P: !protocol
  sequence:
    localFoo: Foo
    importedFoo: Imported.Foo`

	root := t.TempDir()
	importedDir := path.Join(root, "imported")
	baseDir := path.Join(root, "base")
	require.NoError(t, os.MkdirAll(importedDir, 0755))
	require.NoError(t, os.MkdirAll(baseDir, 0755))
	require.NoError(t, os.WriteFile(path.Join(importedDir, "model.yaml"), []byte(importedSrc), 0644))
	require.NoError(t, os.WriteFile(path.Join(baseDir, "model.yaml"), []byte(baseSrc), 0644))

	importedNs, err := ParseYamlInDir(importedDir, "Imported")
	require.NoError(t, err)
	baseNs, err := ParseYamlInDir(baseDir, "Base")
	require.NoError(t, err)

	env, err := Validate([]*Namespace{importedNs, baseNs})
	require.NoError(t, err)

	var protocol *ProtocolDefinition
	for _, ns := range env.Namespaces {
		if ns.Name == "Base" {
			require.Len(t, ns.Protocols, 1)
			protocol = ns.Protocols[0]
		}
	}
	require.NotNil(t, protocol)

	schemaStr := GetProtocolSchemaString(protocol, env.SymbolTable)

	var schema struct {
		Types []map[string]any `json:"types"`
	}
	require.NoError(t, json.Unmarshal([]byte(schemaStr), &schema))

	names := make(map[string]bool)
	for _, entry := range schema.Types {
		if name, ok := entry["name"].(string); ok {
			names[name] = true
			continue
		}
		for _, kind := range []string{"record", "enum", "flags", "alias"} {
			if v, ok := entry[kind]; ok {
				if def, ok := v.(map[string]any); ok {
					if name, ok := def["name"].(string); ok {
						names[name] = true
					}
				}
			}
		}
	}

	assert.True(t, names["Foo"], "base-namespace definition should keep unqualified name; got %v", names)
	assert.True(t, names["Imported.Foo"], "imported definition should be qualified; got %v", names)
}
