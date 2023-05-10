package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestImportStateContextUppercaseWrapper(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "input contains no letter, should be unchanged",
			in:       "*123",
			expected: "*123",
		},
		{
			name:     "input contains digits and only upper case letters, should be unchanged",
			in:       "*2E",
			expected: "*2E",
		},
		{
			name:     "input contains lower case letters, should be mapped to upper case",
			in:       "*f2",
			expected: "*F2",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual string
			f := ImportStateContextUppercaseWrapper(
				func(ctx context.Context, rd *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
					actual = rd.Id()
					return nil, nil
				},
			)
			rd := schema.ResourceData{}
			rd.SetId(tc.in)
			_, _ = f(context.TODO(), &rd, nil)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
