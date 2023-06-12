package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	defer func() {
		c.DeleteBridgePort(bridgePort.Id)
		require.NoError(t, err)

		_, err = c.FindBridgePort(bridgePort.Id)
		require.True(t, IsNotFoundError(err), "expected to get NotFound error")
	}()

	expected := &BridgePort{
		Id:        bridgePort.Id,
		Bridge:    "test_bridge",
		Interface: "*0",
		PVId:      1,
		Comment:   bridgePort.Comment,
	}
	if !reflect.DeepEqual(expected, bridgePort) {
		t.Errorf(`expected and actual bridge port objects are not equal:
		want: %+v,
		got: %+v
	`, expected, bridgePort)
	}
}
