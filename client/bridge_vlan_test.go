package client

import (
	"reflect"
	"testing"
)

func TestBridgeVlanBasic(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	bridge1Name := "test_bridge1"
	bridge1 := &Bridge{
		Name:          bridge1Name,
		FastForward:   false,
		VlanFiltering: false,
		Comment:       "a test bridge",
	}
	_, err := c.AddBridge(bridge1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		if err = c.DeleteBridge(bridge1Name); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}()

	bridge2Name := "test_bridge2"
	bridge2 := &Bridge{
		Name:          bridge2Name,
		FastForward:   false,
		VlanFiltering: false,
		Comment:       "a test bridge",
	}
	_, err = c.AddBridge(bridge2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		if err = c.DeleteBridge(bridge2Name); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}()

	bridgeVlan := &BridgeVlan{
		Bridge:  bridge1.Name,
		VlanIds: []string{"10", "20"},
	}

	createdBridgeVlan, err := c.AddBridgeVlan(bridgeVlan)
	if err != nil {
		t.Fatal(err)
	}

	expectedBridgeVlan := &BridgeVlan{
		Id:     createdBridgeVlan.Id,
		Bridge: bridge1Name,
	}
	if !reflect.DeepEqual(createdBridgeVlan, expectedBridgeVlan) {
		t.Fatal("created bridge vlan differs")
	}

	createdBridgeVlan.Bridge = bridge2Name
	updatedBridgeVlan, err := c.UpdateBridgeVlan(createdBridgeVlan)
	if err != nil {
		t.Fatal(err)
	}
	expectedBridgeVlan = &BridgeVlan{
		Id:     createdBridgeVlan.Id,
		Bridge: bridge2Name,
	}
	if !reflect.DeepEqual(updatedBridgeVlan, expectedBridgeVlan) {
		t.Fatal("updated bridge vlan differs")
	}
}
