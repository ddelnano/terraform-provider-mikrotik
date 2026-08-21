package utils

import (
	"context"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyStruct(t *testing.T) {
	type (
		simpleStruct struct {
			String     string
			ExtraField int
			Int        int
			Int8       int8
			Int16      int16
			Int32      int32
			Int64      int64
			Uint       uint
			Uint8      uint8
			Uint16     uint16
			Uint32     uint32
			Uint64     uint64
			Boolean    bool
			Float32    float32
			Float64    float64
			IntList    []int
			UintList   []uint
			StringList []string
		}

		terraformStruct struct {
			String        tftypes.String `tfsdk:"string"`
			Int           tftypes.Int64  `tfsdk:"int"`
			Int8          tftypes.Int64  `tfsdk:"int_8"`
			Int16         tftypes.Int64  `tfsdk:"int_16"`
			Int32         tftypes.Int64  `tfsdk:"int_32"`
			Int64         tftypes.Int64  `tfsdk:"int_64"`
			Uint          tftypes.Int64  `tfsdk:"uint"`
			Uint8         tftypes.Int64  `tfsdk:"uint_8"`
			Uint16        tftypes.Int64  `tfsdk:"uint_16"`
			Uint32        tftypes.Int64  `tfsdk:"uint_32"`
			Uint64        tftypes.Int64  `tfsdk:"uint_64"`
			UnmappedField tftypes.String `tfsdk:"unmapped_field"`
			Boolean       tftypes.Bool   `tfsdk:"boolean"`
			IntList       tftypes.List   `tfsdk:"int_list"`
			UintList      tftypes.List   `tfsdk:"uint_list"`
			StringList    tftypes.List   `tfsdk:"string_list"`
		}
	)

	testCases := []struct {
		name           string
		src            any
		dest           any
		expected       any
		resourceSchema schema.Schema
		expectError    bool
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
			name: "different case",
			src: struct {
				NAME        string
				SourceField int
				itEMS       []string
			}{
				NAME:        "src field name",
				SourceField: 10,
				itEMS:       []string{"one", "two"},
			},
			dest: &struct {
				Name             string
				DestinationField int
				Items            []string
			}{
				Name:             "dest field name",
				DestinationField: 20,
				Items:            []string{"one", "two", "three"},
			},
			expected: &struct {
				Name             string
				DestinationField int
				Items            []string
			}{
				Name:             "src field name",
				DestinationField: 20,
				Items:            []string{"one", "two", "three"},
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
		{
			name: "core type to core type",
			src: struct {
				String      string
				Int         int
				Int8        int8
				Int16       int16
				Int32       int32
				Int64       int64
				Uint        uint
				Uint8       uint8
				Uint16      uint16
				Uint32      uint32
				Uint64      uint64
				IntSlice    []int
				UintSlice   []uint
				StringSlice []string
			}{
				String:      "source field name",
				Int:         10,
				Int8:        20,
				Int16:       20_000,
				Int32:       20_000_000,
				Int64:       20_000_000_000_000,
				Uint:        10,
				Uint8:       200,
				Uint16:      20_000,
				Uint32:      20_000_000,
				Uint64:      20_000_000_000_000,
				StringSlice: []string{"one", "two"},
				IntSlice:    []int{1, 2, 3},
				UintSlice:   []uint{5, 30},
			},
			dest: &struct {
				String      string
				Int         int
				Int8        int8
				Int16       int16
				Int32       int32
				Int64       int64
				Uint        uint
				Uint8       uint8
				Uint16      uint16
				Uint32      uint32
				Uint64      uint64
				IntSlice    []int
				UintSlice   []uint
				StringSlice []string
			}{},
			expected: &struct {
				String      string
				Int         int
				Int8        int8
				Int16       int16
				Int32       int32
				Int64       int64
				Uint        uint
				Uint8       uint8
				Uint16      uint16
				Uint32      uint32
				Uint64      uint64
				IntSlice    []int
				UintSlice   []uint
				StringSlice []string
			}{
				String:      "source field name",
				Int:         10,
				Int8:        20,
				Int16:       20_000,
				Int32:       20_000_000,
				Int64:       20_000_000_000_000,
				Uint:        10,
				Uint8:       200,
				Uint16:      20_000,
				Uint32:      20_000_000,
				Uint64:      20_000_000_000_000,
				StringSlice: []string{"one", "two"},
				IntSlice:    []int{1, 2, 3},
				UintSlice:   []uint{5, 30},
			},
		},
		{
			name: "core type to terraform type",
			src: simpleStruct{
				String:     "name new",
				ExtraField: 30,
				Int:        10,
				Int8:       20,
				Int16:      20_000,
				Int32:      20_000_000,
				Int64:      20_000_000_000_000,
				Uint:       20,
				Uint8:      200,
				Uint16:     20_000,
				Uint32:     20_000_000,
				Uint64:     20_000_000_000_000,
				Boolean:    true,
				IntList:    []int{10, 20, 30},
				UintList:   []uint{10, 20, 30},
				StringList: []string{"new value"},
			},
			dest: &terraformStruct{
				String:        tftypes.StringValue("field name"),
				Int:           tftypes.Int64Value(20),
				UnmappedField: tftypes.StringValue("unmapped field"),
				Boolean:       tftypes.BoolValue(false),
				IntList: tftypes.ListValueMust(tftypes.Int64Type,
					[]attr.Value{
						tftypes.Int64Value(2),
						tftypes.Int64Value(4),
						tftypes.Int64Value(5),
					}),
				StringList: tftypes.ListValueMust(tftypes.StringType,
					[]attr.Value{
						tftypes.StringValue("old value 1"),
						tftypes.StringValue("old value 2"),
					}),
			},
			resourceSchema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"string":         schema.StringAttribute{},
					"int":            schema.Int64Attribute{},
					"int_8":          schema.Int64Attribute{},
					"int_16":         schema.Int64Attribute{},
					"int_32":         schema.Int64Attribute{},
					"int_64":         schema.Int64Attribute{},
					"uint":           schema.Int64Attribute{},
					"uint_8":         schema.Int64Attribute{},
					"uint_16":        schema.Int64Attribute{},
					"uint_32":        schema.Int64Attribute{},
					"uint_64":        schema.Int64Attribute{},
					"unmapped_field": schema.StringAttribute{},
					"boolean":        schema.BoolAttribute{},
					"int_list":       schema.ListAttribute{},
					"uint_list":      schema.ListAttribute{},
					"string_list":    schema.ListAttribute{},
				},
			},
			expected: &terraformStruct{
				String:        tftypes.StringValue("name new"),
				Int:           tftypes.Int64Value(10),
				Int8:          tftypes.Int64Value(20),
				Int16:         tftypes.Int64Value(20_000),
				Int32:         tftypes.Int64Value(20_000_000),
				Int64:         tftypes.Int64Value(20_000_000_000_000),
				Uint:          tftypes.Int64Value(20),
				Uint8:         tftypes.Int64Value(200),
				Uint16:        tftypes.Int64Value(20_000),
				Uint32:        tftypes.Int64Value(20_000_000),
				Uint64:        tftypes.Int64Value(20_000_000_000_000),
				UnmappedField: tftypes.StringValue("unmapped field"),
				Boolean:       tftypes.BoolValue(true),
				IntList: tftypes.ListValueMust(tftypes.Int64Type,
					[]attr.Value{
						tftypes.Int64Value(10),
						tftypes.Int64Value(20),
						tftypes.Int64Value(30),
					}),
				UintList: tftypes.ListValueMust(tftypes.Int64Type,
					[]attr.Value{
						tftypes.Int64Value(10),
						tftypes.Int64Value(20),
						tftypes.Int64Value(30),
					}),
				StringList: tftypes.ListValueMust(tftypes.StringType,
					[]attr.Value{
						tftypes.StringValue("new value"),
					}),
			},
		},
		{
			name: "terraform type to core type",
			src: struct {
				String        tftypes.String
				Int           tftypes.Int64
				Int8          tftypes.Int64
				Int16         tftypes.Int64
				Int32         tftypes.Int64
				Int64         tftypes.Int64
				Uint          tftypes.Int64
				Uint8         tftypes.Int64
				Uint16        tftypes.Int64
				Uint32        tftypes.Int64
				Uint64        tftypes.Int64
				UnmappedField tftypes.String
				Boolean       tftypes.Bool
				IntList       tftypes.List
				UintList      tftypes.List
				StringList    tftypes.List
			}{
				String:        tftypes.StringValue("new field name"),
				Int:           tftypes.Int64Value(20),
				Int8:          tftypes.Int64Value(20),
				Int16:         tftypes.Int64Value(20_000),
				Int32:         tftypes.Int64Value(20_000_000),
				Int64:         tftypes.Int64Value(20_000_000_000_000),
				UnmappedField: tftypes.StringValue("unmapped field"),
				Boolean:       tftypes.BoolValue(true),
				IntList: tftypes.ListValueMust(tftypes.Int64Type,
					[]attr.Value{
						tftypes.Int64Value(2),
						tftypes.Int64Value(4),
						tftypes.Int64Value(5),
					}),
				UintList: tftypes.ListValueMust(tftypes.Int64Type,
					[]attr.Value{
						tftypes.Int64Value(2),
						tftypes.Int64Value(4),
						tftypes.Int64Value(5),
					}),
				StringList: tftypes.ListValueMust(tftypes.StringType,
					[]attr.Value{
						tftypes.StringValue("new value 1"),
						tftypes.StringValue("new value 2"),
					}),
			},
			dest: &struct {
				String     string
				Int        int
				Int8       int
				Int16      int
				Int32      int
				Int64      int
				ExtraField int
				Boolean    bool
				Float32    float32
				Float64    float64
				IntList    []int
				UintList   []uint
				StringList []string
			}{
				String:     "name old",
				Int:        10,
				ExtraField: 30,
				Boolean:    false,
				IntList:    []int{10, 20, 30},
				StringList: []string{"old value"},
			},
			expected: &struct {
				String     string
				Int        int
				Int8       int
				Int16      int
				Int32      int
				Int64      int
				ExtraField int
				Boolean    bool
				Float32    float32
				Float64    float64
				IntList    []int
				UintList   []uint
				StringList []string
			}{
				String:     "new field name",
				Int:        20,
				Int8:       20,
				Int16:      20_000,
				Int32:      20_000_000,
				Int64:      20_000_000_000_000,
				ExtraField: 30,
				Boolean:    true,
				IntList:    []int{2, 4, 5},
				UintList:   []uint{2, 4, 5},
				StringList: []string{"new value 1", "new value 2"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := copyStruct(context.TODO(), tc.src, tc.dest, tc.resourceSchema)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, tc.dest)
		})
	}
}

func TestCopyTerraformToMikrotik(t *testing.T) {
	testCases := []struct {
		name        string
		src         any
		dest        client.Resource
		expected    any
		expectError bool
	}{
		{
			name: "pass",
			src: struct {
				Id       tftypes.String
				Bridge   tftypes.String
				Tagged   tftypes.Set
				Untagged tftypes.List
				VlanIds  tftypes.List
			}{
				Id:     tftypes.StringValue("new id field"),
				Bridge: tftypes.StringValue("new bridge"),
				Tagged: tftypes.SetValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new tagged 3"),
				}),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new untagged 5"),
				}),

				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(2),
					tftypes.Int64Value(5),
					tftypes.Int64Value(10),
				}),
			},
			dest: &client.BridgeVlan{
				Id:       "old id field",
				Bridge:   "old bridge",
				Tagged:   types.MikrotikList{"old tagged 1", "old tagged 2"},
				Untagged: types.MikrotikList{"old untagged 1"},
				VlanIds:  types.MikrotikIntList{1, 3},
			},
			expected: &client.BridgeVlan{
				Id:       "new id field",
				Bridge:   "new bridge",
				Tagged:   types.MikrotikList{"new tagged 3"},
				Untagged: types.MikrotikList{"new untagged 5"},
				VlanIds:  types.MikrotikIntList{2, 5, 10},
			},
		},
		{
			name: "skip readonly fields",
			src: struct {
				Id       tftypes.String
				Hostname tftypes.String
			}{
				Id:       tftypes.StringValue("new id"),
				Hostname: tftypes.StringValue("should not be copied"),
			},
			dest: &client.DhcpLease{
				Id:       "old id",
				Hostname: "old hostname",
			},
			expected: &client.DhcpLease{
				Id:       "new id",
				Hostname: "old hostname", // should not be changed because it's readonly
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := TerraformModelToMikrotikStruct(context.TODO(), tc.src, tc.dest)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, tc.dest)
		})
	}
}

func TestCopyMikrotikToTerraform(t *testing.T) {
	type (
		terraformStruct struct {
			Id       tftypes.String `tfsdk:"id"`
			Bridge   tftypes.String `tfsdk:"bridge"`
			Tagged   tftypes.List   `tfsdk:"tagged"`
			Untagged tftypes.List   `tfsdk:"untagged"`
			VlanIds  tftypes.List   `tfsdk:"vlan_ids"`
		}
	)
	testCases := []struct {
		name           string
		src            client.Resource
		dest           any
		resourceSchema schema.Schema
		expected       any
		expectError    bool
	}{
		{
			name: "pass",
			src: &client.BridgeVlan{
				Id:       "new id field",
				Bridge:   "new bridge",
				Tagged:   types.MikrotikList{"new tagged 1", "new tagged 2"},
				Untagged: types.MikrotikList{"new untagged 1"},
				VlanIds:  types.MikrotikIntList{1, 3},
			},
			dest: &terraformStruct{
				Id:     tftypes.StringValue("old id field"),
				Bridge: tftypes.StringValue("old bridge"),
				Tagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("old tagged 3"),
				}),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("old untagged 5"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(2),
					tftypes.Int64Value(5),
					tftypes.Int64Value(10),
				}),
			},
			resourceSchema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":       schema.StringAttribute{},
					"bridge":   schema.StringAttribute{},
					"tagged":   schema.ListAttribute{},
					"untagged": schema.ListAttribute{},
					"vlan_ids": schema.ListAttribute{},
				},
			},
			expected: &terraformStruct{
				Id:     tftypes.StringValue("new id field"),
				Bridge: tftypes.StringValue("new bridge"),
				Tagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new tagged 1"),
					tftypes.StringValue("new tagged 2"),
				}),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new untagged 1"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(1),
					tftypes.Int64Value(3),
				}),
			},
		},
		{
			name: "no zero-value write to null value for non-computed",
			src: &client.BridgeVlan{
				Id:       "new id field",
				Bridge:   "new bridge",
				Tagged:   nil,
				Untagged: types.MikrotikList{"new untagged 1"},
				VlanIds:  types.MikrotikIntList{1, 3},
			},
			dest: &terraformStruct{
				Id:     tftypes.StringValue("old id field"),
				Bridge: tftypes.StringValue("old bridge"),
				Tagged: tftypes.ListNull(tftypes.StringType),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("old untagged 5"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(2),
					tftypes.Int64Value(5),
					tftypes.Int64Value(10),
				}),
			},
			resourceSchema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":     schema.StringAttribute{},
					"bridge": schema.StringAttribute{},
					"tagged": schema.ListAttribute{
						Optional: true,
					},
					"untagged": schema.ListAttribute{},
					"vlan_ids": schema.ListAttribute{},
				},
			},
			expected: &terraformStruct{
				Id:     tftypes.StringValue("new id field"),
				Bridge: tftypes.StringValue("new bridge"),
				Tagged: tftypes.ListNull(tftypes.StringType),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new untagged 1"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(1),
					tftypes.Int64Value(3),
				}),
			},
		},
		{
			name: "write if incoming value is non-zero",
			src: &client.BridgeVlan{
				Id:       "new id field",
				Bridge:   "new bridge",
				Tagged:   types.MikrotikList{"new tagged 1", "new tagged 2"},
				Untagged: types.MikrotikList{"new untagged 1"},
				VlanIds:  types.MikrotikIntList{1, 3},
			},
			dest: &terraformStruct{
				Id:     tftypes.StringValue("old id field"),
				Bridge: tftypes.StringNull(),
				Tagged: tftypes.ListNull(tftypes.StringType),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("old untagged 5"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(2),
					tftypes.Int64Value(5),
					tftypes.Int64Value(10),
				}),
			},
			resourceSchema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{},
					"bridge": schema.StringAttribute{
						Optional: true,
					},
					"tagged": schema.ListAttribute{
						Optional: true,
					},
					"untagged": schema.ListAttribute{},
					"vlan_ids": schema.ListAttribute{},
				},
			},
			expected: &terraformStruct{
				Id:     tftypes.StringValue("new id field"),
				Bridge: tftypes.StringValue("new bridge"),
				Tagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new tagged 1"),
					tftypes.StringValue("new tagged 2"),
				}),
				Untagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new untagged 1"),
				}),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(1),
					tftypes.Int64Value(3),
				}),
			},
		},
		{
			name: "write if attribute computed",
			src: &client.BridgeVlan{
				Id:       "new id field",
				Bridge:   "",
				Tagged:   types.MikrotikList{"new tagged 1", "new tagged 2"},
				Untagged: nil,
				VlanIds:  types.MikrotikIntList{1, 3},
			},
			dest: &terraformStruct{
				Id:     tftypes.StringValue("old id field"),
				Bridge: tftypes.StringValue("old bridge"),
				Tagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("old tagged 3"),
				}),
				Untagged: tftypes.ListNull(tftypes.StringType),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(2),
					tftypes.Int64Value(5),
					tftypes.Int64Value(10),
				}),
			},
			resourceSchema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{},
					"bridge": schema.StringAttribute{
						Computed: true,
					},
					"tagged": schema.ListAttribute{},
					"untagged": schema.ListAttribute{
						Computed: true,
					},
					"vlan_ids": schema.ListAttribute{},
				},
			},
			expected: &terraformStruct{
				Id:     tftypes.StringValue("new id field"),
				Bridge: tftypes.StringValue(""),
				Tagged: tftypes.ListValueMust(tftypes.StringType, []attr.Value{
					tftypes.StringValue("new tagged 1"),
					tftypes.StringValue("new tagged 2"),
				}),
				Untagged: tftypes.ListNull(tftypes.StringType),
				VlanIds: tftypes.ListValueMust(tftypes.Int64Type, []attr.Value{
					tftypes.Int64Value(1),
					tftypes.Int64Value(3),
				}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MikrotikStructToTerraformModel(context.TODO(), tc.src, tc.dest, tc.resourceSchema)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, tc.dest)
		})
	}
}
