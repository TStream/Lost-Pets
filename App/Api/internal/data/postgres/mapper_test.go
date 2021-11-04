package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnderscore(t *testing.T) {
	type Test struct {
		Input  string
		Output string
	}

	tests := []Test{
		Test{
			"CamelCase",
			"camel_case",
		},
		Test{
			"CamelCaSe",
			"camel_ca_se",
		},
		Test{
			"camelcase",
			"camelcase",
		},
		Test{
			"CAMELCASE",
			"camelcase",
		},
		Test{
			"CAMELCase",
			"camel_case",
		},
		Test{
			"camelCase",
			"camel_case",
		},
	}

	for i, test := range tests {
		str := underscore(test.Input)
		assert.Equal(t, test.Output, str, "underscore failed at %d", i)
	}
}
