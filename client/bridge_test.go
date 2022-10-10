package client

import (
	"reflect"
	"testing"
)

func TestBridgeBasic(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "test_bridge"
	bridge := &Bridge{
		Name:          name,
		FastForward:   false,
		VlanFiltering: false,
		Comment:       "a test bridge",
	}
	_, err := c.AddBridge(bridge)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, err := c.FindBridge(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	bridge.Id = found.Id
	if !reflect.DeepEqual(bridge, found) {
		t.Fatalf("expected found resource to have pre-defined fields but it didn't")
	}

	updatedResource := &Bridge{
		Id:            found.Id,
		Name:          found.Name + "_updated",
		FastForward:   true,
		VlanFiltering: true,
		Comment:       "updated comment",
	}
	_, err = c.UpdateBridge(updatedResource)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	foundAfterUpdate, err := c.FindBridge(updatedResource.Name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(updatedResource, foundAfterUpdate) {
		t.Fatalf("expected found resource to have pre-defined fields but it didn't")
	}

	if err = c.DeleteBridge(name); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err = c.DeleteBridge(name); err == nil {
		t.Fatal("expected notfound error, got nothing")
	}
}
