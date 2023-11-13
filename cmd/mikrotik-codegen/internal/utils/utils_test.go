package utils

import (
	"testing"
)

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

func TestFirstLower(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "title case",
			input:    "ClientResourceName",
			expected: "clientResourceName",
		},
		{
			name:     "kebab case",
			input:    "clientResourceName",
			expected: "clientResourceName",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FirstLower(tc.input)
			if result != tc.expected {
				t.Errorf(`
				expected %s,
				got %s`, tc.expected, result)
			}
		})
	}
}
func TestPascalCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already PascalCase",
			input:    "FieldNameInProperCase",
			expected: "FieldNameInProperCase",
		},
		{
			name:     "dashes",
			input:    "field-name-with-dashes",
			expected: "FieldNameWithDashes",
		},
		{
			name:     "dashes, underscores",
			input:    "field-name_with_dashes-and___underscores",
			expected: "FieldNameWithDashesAndUnderscores",
		},
		{
			name:     "other symbols",
			input:    "field/name  with+++++different||||symbols",
			expected: "FieldNameWithDifferentSymbols",
		},
		{
			name:     "consecutive upper-cased if one-letter word",
			input:    "field/name  with-a/b-testing",
			expected: "FieldNameWithABTesting",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PascalCase(tc.input)
			if result != tc.expected {
				t.Errorf(`
				expected %s,
				got %s`, tc.expected, result)
			}
		})
	}
}
