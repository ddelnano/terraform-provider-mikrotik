package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVlanInterfaceUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedIface := &VlanInterface{
		Name:      "vlan-20",
		VlanId:    20,
		Mtu:       1000,
		Interface: "*0",
		Disabled:  false,
	}

	iface, err := c.AddVlanInterface(&VlanInterface{
		Name:      expectedIface.Name,
		Disabled:  expectedIface.Disabled,
		Interface: expectedIface.Interface,
		VlanId:    expectedIface.VlanId,
		Mtu:       expectedIface.Mtu,
	})
	require.NoError(t, err)

	expectedIface.Id = iface.Id

	foundInterface, err := c.FindVlanInterface(expectedIface.Name)
	require.NoError(t, err)
	assert.Equal(t, expectedIface, foundInterface)

	expectedIface.Name = expectedIface.Name + "updated"
	expectedIface.Mtu = expectedIface.Mtu - 100
	updatedIface, err := c.UpdateVlanInterface(expectedIface)
	require.NoError(t, err)
	assert.Equal(t, expectedIface, updatedIface)
	// cleanup
	err = c.DeleteVlanInterface(iface.Name)
	assert.NoError(t, err)

	_, err = c.FindVlanInterface(expectedIface.Name)
	assert.Error(t, err)
}
