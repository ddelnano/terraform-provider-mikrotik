package utils

import (
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyStruct(t *testing.T) {

	testCases := []struct {
		name        string
		src         interface{}
		dest        interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "same fields",
			src: struct {
				Name         string
				AnotherField int
				Items        []string
			}{
				Name:         "source field name",
				AnotherField: 10,
				Items:        []string{"one", "two"},
			},
			dest: &struct {
				Name         string
				AnotherField int
				Items        []string
			}{
				Name:         "destination field name",
				AnotherField: 20,
				Items:        []string{"one", "two", "three"},
			},
			expected: &struct {
				Name         string
				AnotherField int
				Items        []string
			}{
				Name:         "source field name",
				AnotherField: 10,
				Items:        []string{"one", "two"},
			},
		},
		{
			name: "overlapping fields",
			src: struct {
				Name        string
				SourceField int
				Items       []string
			}{
				Name:        "field name",
				SourceField: 10,
				Items:       []string{"one", "two"},
			},
			dest: &struct {
				Name             string
				DestinationField int
				Items            []string
			}{
				Name:             "field name",
				DestinationField: 20,
				Items:            []string{"one", "two", "three"},
			},
			expected: &struct {
				Name             string
				DestinationField int
				Items            []string
			}{
				Name:             "field name",
				DestinationField: 20,
				Items:            []string{"one", "two"},
			},
		},
		{
			name: "custom field to regular",
			src: client.BridgeVlan{
				Id:       "identifier",
				Bridge:   "bridge1",
				Tagged:   types.MikrotikList{"tagged1", "tagged2"},
				Untagged: types.MikrotikList{"untagged1", "untagged2"},
				VlanIds:  types.MikrotikIntList{3, 4, 5},
			},
			dest: &struct {
				Id         string
				Bridge     string
				Tagged     []string
				Untagged   []string
				VlanIds    []int
				ExtraField string
			}{
				Id:         "identifier old",
				Bridge:     "bridge old",
				Tagged:     []string{"tagged old"},
				Untagged:   []string{"untagged old"},
				VlanIds:    []int{2},
				ExtraField: "unchanged",
			},
			expected: &struct {
				Id         string
				Bridge     string
				Tagged     []string
				Untagged   []string
				VlanIds    []int
				ExtraField string
			}{
				Id:         "identifier",
				Bridge:     "bridge1",
				Tagged:     []string{"tagged1", "tagged2"},
				Untagged:   []string{"untagged1", "untagged2"},
				VlanIds:    []int{3, 4, 5},
				ExtraField: "unchanged",
			},
		},
		{
			name: "regular type to custom field",
			src: struct {
				Id         string
				Bridge     string
				Tagged     []string
				Untagged   []string
				VlanIds    []int
				ExtraField string
			}{
				Id:         "identifier new",
				Bridge:     "bridge new",
				Tagged:     []string{"tagged new"},
				Untagged:   []string{"untagged new"},
				VlanIds:    []int{2},
				ExtraField: "extra field",
			},
			dest: &client.BridgeVlan{
				Id:       "identifier",
				Bridge:   "bridge1",
				Tagged:   types.MikrotikList{"tagged1", "tagged2"},
				Untagged: types.MikrotikList{"untagged1", "untagged2"},
				VlanIds:  types.MikrotikIntList{3, 4, 5},
			},
			expected: &client.BridgeVlan{
				Id:       "identifier new",
				Bridge:   "bridge new",
				Tagged:   types.MikrotikList{"tagged new"},
				Untagged: types.MikrotikList{"untagged new"},
				VlanIds:  types.MikrotikIntList{2},
			},
		},
		{
			name: "field type mismatch",
			src: struct {
				Name         string
				AnotherField float64
				Items        []string
			}{
				Name:         "field name",
				AnotherField: 10,
				Items:        []string{"one", "two"},
			},
			dest: &struct {
				Name         string
				AnotherField int
				Items        []string
			}{
				Name:         "field name",
				AnotherField: 10,
				Items:        []string{"one", "two"},
			},
			expectError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CopyStruct(tc.src, tc.dest)
			require.True(t, (err != nil) == tc.expectError, "expected err to be %v but got %v", tc.expectError, err)
			if tc.expectError {
				return
			}
			assert.Equal(t, tc.expected, tc.dest)
		})
	}
}
