package defaultaware

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/stretchr/testify/require"
)

func TestInt64Wrapper(t *testing.T) {
	testCases := []struct {
		name           string
		description    string
		defaultValue   defaults.Int64
		expectedResult string
	}{
		{
			name:           "no default value",
			description:    "Attribute description.",
			expectedResult: "Attribute description.",
		},
		{
			name:           "with default value",
			description:    "Attribute description.",
			defaultValue:   int64default.StaticInt64(2),
			expectedResult: "Attribute description. Default: `2`.",
		}, {
			name:           "with zero default value",
			description:    "Attribute description.",
			defaultValue:   int64default.StaticInt64(0),
			expectedResult: "Attribute description. Default: `0`.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr := Int64Attribute(
				schema.Int64Attribute{
					Description: tc.description,
					Default:     tc.defaultValue,
				},
			)
			require.Equal(t, tc.expectedResult, attr.GetDescription())
		})
	}
}
