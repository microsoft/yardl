// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		intput   string
		expected string
	}{
		{"", ""},
		{"A", "A"},
		{"a", "A"},
		{"aB", "AB"},
		{"99", "99"},
		{"ioReader", "IoReader"},
		{"IoReader", "IoReader"},
		{"IOReader", "IOReader"},
		{"snake_case", "SnakeCase"},
		{"snake__case", "SnakeCase"},
		{"_snake_case_", "SnakeCase"},
		{"kebab-case", "KebabCase"},
		{"kebab--case", "KebabCase"},
		{"apple banana", "AppleBanana"},
	}
	for _, tt := range tests {
		t.Run(tt.intput, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToPascalCase(tt.intput))
		})
	}
}

func TestPascalOrCamelToSnakeCase(t *testing.T) {
	tests := []struct {
		intput   string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"Aa", "aa"},
		{"aAa", "a_aa"},
		{"ioReader", "io_reader"},
		{"IoReader", "io_reader"},
		{"IOReader", "io_reader"},
		{"_IOReader", "_io_reader"},
		{"IO_Reader", "io_reader"},
		{"IO_READER", "io_reader"},
		{"readA", "read_a"},
		{"DynamicNDArray", "dynamic_nd_array"},
		{"parseHTMLString", "parse_html_string"},
		{"getElementById", "get_element_by_id"},
		{"CSSSelectorsList", "css_selectors_list"},
		{"iD", "i_d"},
		{"tEST", "t_est"},
		{"convertM4AToMP3", "convert_m4a_to_mp3"},
		{"snake_case", "snake_case"},
		{"Capital_Snake_Case", "capital_snake_case"},
		{"YAML", "yaml"},
		{"YAML2", "yaml2"},
		{"yaml2Spec", "yaml2_spec"},
		{"YAML2Spec", "yaml2_spec"},
		{"YAML42Spec", "yaml42_spec"},
	}
	for _, tt := range tests {
		t.Run(tt.intput, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToSnakeCase(tt.intput))
		})
	}
}
