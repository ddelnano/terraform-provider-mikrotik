package defaultaware

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/stretchr/testify/require"
)

func TestBoolWrapper(t *testing.T) {
	testCases := []struct {
		name           string
		description    string
		defaultValue   defaults.Bool
		expectedResult string
	}{
		{
			name:           "no default value",
			description:    "Attribute description.",
			expectedResult: "Attribute description.",
		},
		{
			name:           "true default value",
			description:    "Attribute description.",
			defaultValue:   booldefault.StaticBool(true),
			expectedResult: "Attribute description. Default: `true`.",
		}, {
			name:           "false default value",
			description:    "Attribute description.",
			defaultValue:   booldefault.StaticBool(false),
			expectedResult: "Attribute description. Default: `false`.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr := BoolAttribute(
				schema.BoolAttribute{
					Description: tc.description,
					Default:     tc.defaultValue,
				},
			)
			require.Equal(t, tc.expectedResult, attr.GetDescription())
		})
	}
}
