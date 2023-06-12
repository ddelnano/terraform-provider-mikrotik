package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNotFoundError(t *testing.T) {

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			expected: false,
		},
		{
			name:     "created via NewNotFoundError()",
			err:      NewNotFound("not found"),
			expected: true,
		},
		{
			name:     "created directly via struct initialization",
			err:      NotFound{},
			expected: true,
		},
		{
			name:     "chained with other errors",
			err:      fmt.Errorf("cannot load object info: %w", NewNotFound("no such object")),
			expected: true,
		},
		{
			name:     "chain of non-matching errors",
			err:      fmt.Errorf("cannot load object info: %w", errors.New("no such object")),
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("no such object"),
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsNotFoundError(tc.err)
			require.Equal(t, tc.expected, result)
		})
	}
}
