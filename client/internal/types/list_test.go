package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMikrotikList_marshal(t *testing.T) {
	testCases := []struct {
		name     string
		list     MikrotikList
		expected string
	}{
		{
			name:     "empty list",
			list:     MikrotikList{""},
			expected: "",
		},
		{
			name:     "nil list",
			list:     nil,
			expected: "",
		},
		{
			name:     "one element",
			list:     MikrotikList{"2"},
			expected: "2",
		},
		{
			name:     "several elements",
			list:     MikrotikList{"2", "3", "one more"},
			expected: "2,3,one more",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.list.MarshalMikrotik()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMikrotikList_unmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		in          string
		expected    MikrotikList
		expectError bool
	}{
		{
			name:     "empty list",
			in:       "",
			expected: MikrotikList{""},
		},
		{
			name:     "one element",
			in:       "one",
			expected: MikrotikList{"one"},
		},
		{
			name:     "several elements",
			in:       "one,two,three",
			expected: MikrotikList{"one", "two", "three"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := MikrotikList{}
			err := l.UnmarshalMikrotik(tc.in)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, l)
		})
	}
}

func TestMikrotikIntList_marshal(t *testing.T) {
	testCases := []struct {
		name     string
		list     MikrotikIntList
		expected string
	}{
		{
			name:     "empty list",
			list:     MikrotikIntList{},
			expected: "",
		},
		{
			name:     "nil list",
			list:     nil,
			expected: "",
		},
		{
			name:     "one element",
			list:     MikrotikIntList{2},
			expected: "2",
		},
		{
			name:     "several elements",
			list:     MikrotikIntList{2, 3, 5},
			expected: "2,3,5",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.list.MarshalMikrotik()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMikrotikIntList_unmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		in          string
		expected    MikrotikIntList
		expectError bool
	}{
		{
			name:     "empty list",
			in:       "",
			expected: MikrotikIntList{},
		},
		{
			name:     "one element",
			in:       "2",
			expected: MikrotikIntList{2},
		},
		{
			name:     "several elements",
			in:       "2,10,15",
			expected: MikrotikIntList{2, 10, 15},
		},
		{
			name:        "bogus element",
			in:          "2,10,not_an_integer,15",
			expectError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := MikrotikIntList{}
			err := l.UnmarshalMikrotik(tc.in)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, l)
		})
	}
}
