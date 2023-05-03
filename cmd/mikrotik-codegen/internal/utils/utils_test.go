package utils

import "testing"

func TestToSnakeCase(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "title case",
			input:    "ClientResourceName",
			expected: "client_resource_name",
		},
		{
			name:     "kebab case",
			input:    "clientResourceName",
			expected: "client_resource_name",
		},
		{
			name:     "all lowercase",
			input:    "clientresourcename",
			expected: "clientresourcename",
		},
		{
			name:     "several uppercase at beginning",
			input:    "IPAddress",
			expected: "ipaddress",
		},
		{
			name:     "several uppercase inside",
			input:    "DefaultHTTPConfig",
			expected: "default_httpconfig",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToSnakeCase(tc.input)
			if result != tc.expected {
				t.Errorf(`
				expected %s,
				got %s`, tc.expected, result)
			}
		})
	}
}
