package defaultaware

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/stretchr/testify/require"
)

func TestStringWrapper(t *testing.T) {
	testCases := []struct {
		name           string
		description    string
		defaultValue   defaults.String
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
			defaultValue:   stringdefault.StaticString("some value"),
			expectedResult: "Attribute description. Default: `some value`.",
		}, {
			name:           "with empty string default value",
			description:    "Attribute description.",
			defaultValue:   stringdefault.StaticString(""),
			expectedResult: "Attribute description. Default: `\"\"`.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr := StringAttribute(
				schema.StringAttribute{
					Description: tc.description,
					Default:     tc.defaultValue,
				},
			)
			require.Equal(t, tc.expectedResult, attr.GetDescription())
		})
	}
}
