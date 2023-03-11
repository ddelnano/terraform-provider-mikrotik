package client

import (
	"errors"
	"reflect"
	"testing"
)

func TestBridgePort_basic(t *testing.T) {
	c := NewClient(GetConfigFromEnv())
	bridge, err := c.AddBridge(&Bridge{
		Name: "test_bridge",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		if err := c.DeleteBridge(bridge.Name); err != nil {
			t.Error(err)
		}
	}()

	bridgePort, err := c.AddBridgePort(&BridgePort{
		Bridge:    bridge.Name,
		Interface: "*0",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		if err := c.DeleteBridgePort(bridgePort.Id); err != nil {
			t.Error(err)

		}
		expected := &NotFound{}
		if _, err := c.FindBridgePort(bridgePort.Id); err == nil || !errors.As(err, &expected) {
			t.Error(err)
		}
	}()

	expected := &BridgePort{
		Id:         bridgePort.Id,
		Bridge:     "test_bridge",
		Interface:  "*0",
		PVId:       1,
		Comment:    bridgePort.Comment,
		FrameTypes: "admit-all",
	}
	if !reflect.DeepEqual(expected, bridgePort) {
		t.Errorf(`expected and actual bridge port objects are not equal:
		want: %+v,
		got: %+v
	`, expected, bridgePort)
	}
}
