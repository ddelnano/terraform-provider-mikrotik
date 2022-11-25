package client

import (
	"reflect"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client/internal/types"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"github.com/stretchr/testify/assert"
)

func TestTtlToSeconds(t *testing.T) {
	tests := []struct {
		expected int
		input    string
	}{
		{300, "5m"},
		{59, "59s"},
		{301, "5m1s"},
		{55259, "15h20m59s"},
		{141659, "1d15h20m59s"},
		{228059, "2d15h20m59s"},
		{86400, "1d"},
	}

	for _, test := range tests {
		actual := ttlToSeconds(test.input)
		if test.expected != actual {
			t.Errorf("Input %s returned %d instead of %d", test.input, actual, test.expected)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	type testStruct struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
		RunCount      int    `mikrotik:"run-count"`
		Allowed       bool
		Schedule      types.MikrotikList
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

func TestUnmarshal_ttlToSeconds(t *testing.T) {
	ttlStr := "5m"
	expectedTtl := ttlToSeconds(ttlStr)
	testStruct := struct {
		Ttl int `mikrotik:"ttl,ttlToSeconds"`
	}{}
	reply := routeros.Reply{
		Re: []*proto.Sentence{
			{
				Word: "!re",
				List: []proto.Pair{
					{
						Key:   "ttl",
						Value: ttlStr,
					},
				},
			},
		},
	}
	err := Unmarshal(reply, &testStruct)

	if err != nil {
		t.Errorf("Failed to unmarshal with error: %v", err)
	}

	if testStruct.Ttl != expectedTtl {
		t.Errorf("Failed to unmarshal ttl field with `ttlToSeconds`. Expected: '%d', received: %v", expectedTtl, testStruct.Ttl)
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
				Allowed       bool   `mikrotik:"allowed-or-not"`
			}{
				Name:          "test owner",
				NotNamedOwner: "admin",
				RunCount:      3,
				Allowed:       true,
			},
			expectedCmd: []string{
				"/test/owner/add",
				"=name=test owner",
				"=owner=admin",
				"=run-count=3",
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
