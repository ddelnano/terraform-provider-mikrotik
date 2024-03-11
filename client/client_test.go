package client

import (
	"reflect"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	type testStruct struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
		RunCount      int    `mikrotik:"run-count"`
		RunCount8     int8   `mikrotik:"run-count8"`
		RunCount16    int16  `mikrotik:"run-count16"`
		RunCount32    int32  `mikrotik:"run-count32"`
		RunCount64    int64  `mikrotik:"run-count64"`
		CountUint     uint   `mikrotik:"run-count-uint"`
		CountUint8    uint8  `mikrotik:"run-count-uint8"`
		CountUint16   uint16 `mikrotik:"run-count-uint16"`
		CountUint32   uint32 `mikrotik:"run-count-uint32"`
		CountUint64   uint64 `mikrotik:"run-count-uint64"`

		Allowed  bool
		Schedule types.MikrotikList
	}

	testCases := []struct {
		name           string
		expectedResult testStruct
		reply          routeros.Reply
	}{
		{
			name: "basic types only",
			reply: routeros.Reply{
				Re: []*proto.Sentence{
					{
						Word: "!re",
						List: []proto.Pair{
							{
								Key:   "name",
								Value: "testing script",
							},
							{
								Key:   "owner",
								Value: "admin",
							},
							{
								Key:   "run-count",
								Value: "3",
							},
							{
								Key:   "run-count8",
								Value: "-3",
							},
							{
								Key:   "run-count16",
								Value: "12000",
							},
							{
								Key:   "run-count32",
								Value: "12000000",
							},
							{
								Key:   "run-count64",
								Value: "12000000000000",
							},
							{
								Key:   "run-count-uint",
								Value: "500",
							},
							{
								Key:   "run-count-uint8",
								Value: "5",
							},
							{
								Key:   "run-count-uint16",
								Value: "15000",
							},
							{
								Key:   "run-count-uint32",
								Value: "15000000",
							},
							{
								Key:   "run-count-uint64",
								Value: "15000000000000000",
							},
							{
								Key:   "allowed",
								Value: "true",
							},
						},
					},
				},
			},
			expectedResult: testStruct{
				Name:          "testing script",
				NotNamedOwner: "admin",
				RunCount:      3,
				RunCount8:     -3,
				RunCount16:    12000,
				RunCount32:    12000000,
				RunCount64:    12000000000000,
				CountUint:     500,
				CountUint8:    5,
				CountUint16:   15000,
				CountUint32:   15000000,
				CountUint64:   15000000000000000,
				Allowed:       true,
			},
		},
		{
			name: "MikrotikList type",
			reply: routeros.Reply{
				Re: []*proto.Sentence{
					{
						Word: "!re",
						List: []proto.Pair{
							{
								Key:   "owner",
								Value: "admin",
							},
							{
								Key:   "schedule",
								Value: "mon,wed,fri",
							},
						},
					},
				},
			},
			expectedResult: testStruct{
				NotNamedOwner: "admin",
				Schedule:      []string{"mon", "wed", "fri"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			targetStruct := testStruct{}
			err := Unmarshal(tc.reply, &targetStruct)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, targetStruct)
		})
	}
}

func TestUnmarshalOnSlices(t *testing.T) {
	name := "testing script"
	owner := "admin"
	allowed := "true"
	type testStruct struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
		Allowed       bool
	}

	cases := []struct {
		name     string
		reply    routeros.Reply
		expected []testStruct
	}{
		{
			name: "reply with len > 1",
			reply: routeros.Reply{
				Re: []*proto.Sentence{
					{
						Word: "!re",
						List: []proto.Pair{
							{
								Key:   "name",
								Value: name,
							},
							{
								Key:   "owner",
								Value: owner,
							},
							{
								Key:   "allowed",
								Value: allowed,
							},
						},
					},
					{
						Word: "!re",
						List: []proto.Pair{
							{
								Key:   "name",
								Value: name + " 2",
							},
							{
								Key:   "owner",
								Value: owner + " 2",
							},
							{
								Key:   "allowed",
								Value: allowed,
							},
						},
					},
				},
			},
			expected: []testStruct{
				{
					Name:          name,
					NotNamedOwner: owner,
					Allowed:       true,
				},
				{
					Name:          name + " 2",
					NotNamedOwner: owner + " 2",
					Allowed:       true,
				},
			},
		},
		{
			name: "reply with single element",
			reply: routeros.Reply{
				Re: []*proto.Sentence{
					{
						Word: "!re",
						List: []proto.Pair{
							{
								Key:   "name",
								Value: name,
							},
							{
								Key:   "owner",
								Value: owner,
							},
							{
								Key:   "allowed",
								Value: allowed,
							},
						},
					},
				},
			},
			expected: []testStruct{
				{
					Name:          name,
					NotNamedOwner: owner,
					Allowed:       true,
				},
			},
		},
		{
			name: "empty reply",
			reply: routeros.Reply{
				Re: []*proto.Sentence{},
			},
			expected: []testStruct{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var result []testStruct
			err := Unmarshal(tc.reply, &result)

			if err != nil {
				t.Errorf("Failed to unmarshal with error: %v", err)
			}

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf(`unexpected result:
				want: %#v
				got: %#v
				`, tc.expected, result)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	action := "/test/owner/add"
	testCases := []struct {
		name        string
		testStruct  interface{}
		expectedCmd []string
	}{
		{
			name: "basic types",
			testStruct: struct {
				Name          string `mikrotik:"name"`
				NotNamedOwner string `mikrotik:"owner,extraTagNotUsed"`
				RunCount      int    `mikrotik:"run-count"`
				RunCount8     int8   `mikrotik:"run-count8"`
				RunCount16    int16  `mikrotik:"run-count16"`
				RunCount32    int32  `mikrotik:"run-count32"`
				RunCount64    int64  `mikrotik:"run-count64"`
				CountUint     uint   `mikrotik:"run-count-uint"`
				CountUint8    uint8  `mikrotik:"run-count-uint8"`
				CountUint16   uint16 `mikrotik:"run-count-uint16"`
				CountUint32   uint32 `mikrotik:"run-count-uint32"`
				CountUint64   uint64 `mikrotik:"run-count-uint64"`
				ReadOnlyProp  bool   `mikrotik:"read-only-prop,readonly"`
				Allowed       bool   `mikrotik:"allowed-or-not"`
			}{
				Name:          "test owner",
				NotNamedOwner: "admin",
				RunCount:      3,
				RunCount8:     10,
				RunCount16:    12000,
				RunCount32:    -12_000_000,
				RunCount64:    12_000_000_000_000_000,
				CountUint:     15000,
				CountUint8:    250,
				CountUint16:   15000,
				CountUint32:   15_000_000,
				CountUint64:   15_000_000_000_000_000,
				Allowed:       true,
			},
			expectedCmd: []string{
				"/test/owner/add",
				"=name=test owner",
				"=owner=admin",
				"=run-count=3",
				"=run-count8=10",
				"=run-count16=12000",
				"=run-count32=-12000000",
				"=run-count64=12000000000000000",
				"=run-count-uint=15000",
				"=run-count-uint8=250",
				"=run-count-uint16=15000",
				"=run-count-uint32=15000000",
				"=run-count-uint64=15000000000000000",
				"=allowed-or-not=yes",
			},
		},
		{
			name: "MikrotikList type",
			testStruct: struct {
				Name          string             `mikrotik:"name"`
				NotNamedOwner string             `mikrotik:"owner,extraTagNotUsed"`
				RunCount      int                `mikrotik:"run-count"`
				Allowed       bool               `mikrotik:"allowed-or-not"`
				Schedule      types.MikrotikList `mikrotik:"schedule"`
			}{
				Name:          "test owner",
				NotNamedOwner: "admin",
				RunCount:      3,
				Allowed:       true,
				Schedule:      []string{"mon", "tue", "fri"},
			},
			expectedCmd: []string{
				"/test/owner/add",
				"=name=test owner",
				"=owner=admin",
				"=run-count=3",
				"=allowed-or-not=yes",
				"=schedule=mon,tue,fri",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Marshal(action, tc.testStruct)
			if !reflect.DeepEqual(cmd, tc.expectedCmd) {
				t.Errorf("Failed to marshal: %v does not equal expected %v", cmd, tc.expectedCmd)
			}
		})
	}
}

func TestMarshalStructWithoutTags(t *testing.T) {
	action := "/test/owner/add"
	name := "test owner"
	owner := "admin"
	runCount := 3
	allowed := true
	retain := false
	testStruct := struct {
		Name           string `example:"name"`
		NotNamedOwner  string `json:"not-named"`
		RunCount       int
		Allowed        bool
		Retain         bool
		SecondaryOwner string
	}{name, owner, runCount, allowed, retain, ""}

	expectedCmd := []string{action}
	cmd := Marshal(action, &testStruct)

	if !reflect.DeepEqual(cmd, expectedCmd) {
		t.Errorf("Marshaling with a struct without tags should return the command action supplied: %v does not equal expected %v", cmd, expectedCmd)
	}
}
