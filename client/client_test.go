package client

import (
	"strconv"
	"strings"
	"testing"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
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
	name := "testing script"
	owner := "admin"
	runCount := "3"
	allowed := "true"
	testStruct := struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
		RunCount      int    `mikrotik:"run-count"`
		Allowed       bool
	}{}
	reply := routeros.Reply{
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
						Key:   "run-count",
						Value: runCount,
					},
					{
						Key:   "allowed",
						Value: allowed,
					},
				},
			},
		},
	}
	err := Unmarshal(reply, &testStruct)

	if err != nil {
		t.Errorf("Failed to unmarshal with error: %v", err)
	}

	if strings.Compare(name, testStruct.Name) != 0 {
		t.Errorf("Failed to unmarshal name '%s' and testStruct.name '%s' should match", name, testStruct.Name)
	}

	if strings.Compare(owner, testStruct.NotNamedOwner) != 0 {
		t.Errorf("Failed to unmarshal name '%s' and testStruct.Name '%s' should match", name, testStruct.Name)
	}

	intRunCount, err := strconv.Atoi(runCount)
	if intRunCount != testStruct.RunCount || err != nil {
		t.Errorf("Failed to unmarshal run-count '%s' and testStruct.RunCount '%d' should match", runCount, testStruct.RunCount)
	}

	b, _ := strconv.ParseBool(allowed)
	if testStruct.Allowed != b {
		t.Errorf("Failed to unmarshal Allowed '%v' and testStruct.Allowed'%v' should match", b, testStruct.Allowed)

	}
}

func TestUnmarshalOnSlices(t *testing.T) {
	name := "testing script"
	owner := "admin"
	allowed := "true"
	testStruct := []struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
		Allowed       bool
	}{{}}
	reply := routeros.Reply{
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
	}
	err := Unmarshal(reply, &testStruct)

	if err != nil {
		t.Errorf("Failed to unmarshal with error: %v", err)
	}

	if strings.Compare(name, testStruct[0].Name) != 0 {
		t.Errorf("Failed to unmarshal name '%s' and testStruct.name '%s' should match", name, testStruct[0].Name)
	}

	if strings.Compare(owner, testStruct[0].NotNamedOwner) != 0 {
		t.Errorf("Failed to unmarshal name '%s' and testStruct.name '%s' should match", name, testStruct[0].Name)
	}

	b, _ := strconv.ParseBool(allowed)
	if testStruct[0].Allowed != b {
		t.Errorf("Failed to unmarshal Allowed '%v' and testStruct.Allowed'%v' should match", b, testStruct[0].Allowed)

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
	name := "test owner"
	owner := "admin"
	runCount := 3
	allowed := true
	retain := false
	testStruct := struct {
		Name           string `mikrotik:"name"`
		NotNamedOwner  string `mikrotik:"owner,extraTagNotUsed"`
		RunCount       int    `mikrotik:"run-count"`
		Allowed        bool   `mikrotik:"allowed-or-not"`
		Retain         bool   `mikrotik:"retain"`
		SecondaryOwner string `mikrotik:"secondary-owner"`
	}{name, owner, runCount, allowed, retain, ""}

	expectedAttributes := "=name=test owner =owner=admin =run-count=3 =allowed-or-not=yes =retain=no"
	// Marshal by passing pointer to struct
	attributes := Marshal(&testStruct)

	if attributes != expectedAttributes {
		t.Errorf("Failed to marshal: %v does not equal expected %v", attributes, expectedAttributes)
	}

	// Marshal by passing by struct value
	attributes = Marshal(testStruct)

	if attributes != expectedAttributes {
		t.Errorf("Failed to marshal: %v does not equal expected %v", attributes, expectedAttributes)
	}
}

func TestMarshalStructWithoutTags(t *testing.T) {
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

	expectedAttributes := ""
	attributes := Marshal(&testStruct)

	if attributes != expectedAttributes {
		t.Errorf("Marshaling with a struct without tags shoudl return empty attributes for command: %v does not equal expected %v", attributes, expectedAttributes)
	}
}
