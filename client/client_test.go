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
	allowed := "true"
	testStruct := struct {
		Name          string
		NotNamedOwner string `mikrotik:"owner"`
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
		t.Errorf("Failed to unmarshal name '%s' and testStruct.name '%s' should match", name, testStruct.Name)
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
