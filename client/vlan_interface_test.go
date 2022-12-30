package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Fatal(err)
	}

	expectedIface.Id = iface.Id

	foundInterface, err := c.FindVlanInterface(expectedIface.Name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedIface, foundInterface)

	expectedIface.Name = expectedIface.Name + "updated"
	expectedIface.Mtu = expectedIface.Mtu - 100
	updatedIface, err := c.UpdateVlanInterface(expectedIface)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedIface, updatedIface)
	// cleanup
	if err := c.DeleteVlanInterface(iface.Name); err != nil {
		t.Error(err)
	}

	_, err = c.FindVlanInterface(expectedIface.Name)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
