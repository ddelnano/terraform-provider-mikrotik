package main

import "testing"

func TestStructNameToResourceFilename(t *testing.T) {
	cases := []struct {
		name       string
		structName string
		expected   string
	}{
		{
			name:       "happy easy",
			structName: "IpAddress",
			expected:   "resource_ip_address.go",
		},
		{
			name:       "all caps",
			structName: "DNS",
			expected:   "resource_dns.go",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := structNameToResourceFilename(tc.structName)
			if result != tc.expected {
				t.Errorf(`
				expected: %s
				got: %s
				`, tc.expected, result)
			}
		})
	}
}
