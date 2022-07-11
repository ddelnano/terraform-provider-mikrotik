package client

import (
	"reflect"
	"testing"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

func TestNewTargetStruct(t *testing.T) {
	type targetStruct struct {
		ID       string `mikrotik:".id"`
		Name     string `mikrotik:"name"`
		Value    int    `mikrotik:"value"`
		Comments string `mikrotik:"comments"`
		Enabled  bool   `mikrotik:"enabled"`
	}

	wrapper := resourceWrapper{
		targetStruct: &targetStruct{},
	}

	reply := routeros.Reply{
		Re: []*proto.Sentence{
			{
				Word: "!re",
				List: []proto.Pair{
					{
						Key:   ".id",
						Value: "recordID",
					},
					{
						Key:   "name",
						Value: "recordName",
					},
					{
						Key:   "value",
						Value: "42",
					},
					{
						Key:   "enabled",
						Value: "true",
					},
				},
			},
		},
	}

	expectedResult := targetStruct{
		ID:      "recordID",
		Name:    "recordName",
		Value:   42,
		Enabled: true,
	}

	newStruct := wrapper.newTargetStruct().Interface()
	err := Unmarshal(reply, newStruct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(expectedResult, *newStruct.(*targetStruct)) {
		t.Errorf(`expected and actual results are different:
				want: %+#v

				got: %+#v`, expectedResult, *newStruct.(*targetStruct))
	}
}

func TestNewTargetStructsList(t *testing.T) {
	type targetStruct struct {
		ID       string `mikrotik:".id"`
		Name     string `mikrotik:"name"`
		Value    int    `mikrotik:"value"`
		Comments string `mikrotik:"comments"`
		Enabled  bool   `mikrotik:"enabled"`
	}

	wrapper := resourceWrapper{
		targetStruct: &targetStruct{},
	}

	reply := routeros.Reply{
		Re: []*proto.Sentence{
			{
				Word: "!re",
				List: []proto.Pair{
					{
						Key:   ".id",
						Value: "id1",
					},
					{
						Key:   "name",
						Value: "name 1",
					},
					{
						Key:   "value",
						Value: "42",
					},
					{
						Key:   "enabled",
						Value: "true",
					},
				},
			},
			{
				Word: "!re",
				List: []proto.Pair{
					{
						Key:   ".id",
						Value: "id2",
					},
					{
						Key:   "name",
						Value: "name 2",
					},
					{
						Key:   "value",
						Value: "43",
					},
				},
			},
		},
	}

	expectedResult := []targetStruct{
		{
			ID:      "id1",
			Name:    "name 1",
			Value:   42,
			Enabled: true,
		},
		{
			ID:      "id2",
			Name:    "name 2",
			Value:   43,
			Enabled: false,
		},
	}

	newStructs := *wrapper.newListOfTargetStructs().Interface().(*[]targetStruct)
	err := Unmarshal(reply, &newStructs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(expectedResult) != len(newStructs) {
		t.Fatalf("expected and resulting list have different length: %d != %d", len(expectedResult), len(newStructs))
	}

	for i := range expectedResult {
		if !reflect.DeepEqual(expectedResult[i], newStructs[i]) {
			t.Errorf(`expected and actual element #%d:
				want: %+#v
				got: %+#v

				expected list: %+#v

				actual list: %+#v
				`, i, expectedResult[i], newStructs[i],
				expectedResult, newStructs,
			)
		}
	}
}
