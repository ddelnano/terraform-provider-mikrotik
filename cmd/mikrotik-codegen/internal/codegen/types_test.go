package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeIs(t *testing.T) {
	testCases := []struct {
		name     string
		type1    Type
		type2    Type
		expected bool
	}{
		{
			name:     "int==int",
			type1:    Int64Type,
			type2:    Int64Type,
			expected: true,
		},
		{
			name:     "list==list",
			type1:    ListType,
			type2:    ListType,
			expected: true,
		},
		{
			name:  "int==string",
			type1: Int64Type,
			type2: StringType,
		},
		{
			name:  "list==string",
			type1: ListType,
			type2: StringType,
		},
		{
			name:  "list==set",
			type1: ListType,
			type2: SetType,
		},
		{
			name:  "bool==unknown",
			type1: BoolType,
			type2: UnknownType,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.type1.Is(tc.type2)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
