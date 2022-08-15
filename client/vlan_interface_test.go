package client

import "testing"

func TestAddVlanInterfaceUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "vlan-20"
	vlanID := 20

	iface, err := c.AddVlanInterface(&VlanInterface{
		Name:      name,
		Disabled:  false,
		Interface: "*0",
		VlanId:    vlanID,
	})
	if err != nil {
		t.Fatal(err)
	}

	foundInterface, err := c.FindVlanInterface(name)
	if err != nil {
		t.Fatal(err)
	}

	if foundInterface.Name != name {
		t.Errorf("expected name to be %q, got %q", name, foundInterface.Name)
	}
	if foundInterface.VlanId != vlanID {
		t.Errorf("expected VlanID to be %d, got %d", vlanID, foundInterface.VlanId)
	}

	iface.Name = foundInterface.Name + "updated"
	updatedIface, err := c.UpdateVlanInterface(iface)
	if err != nil {
		t.Error(err)
	}

	if updatedIface.Name != iface.Name {
		t.Errorf("expected name to be %q, got %q", iface.Name, updatedIface.Name)
	}

	// cleanup
	if err := c.DeleteVlanInterface(iface.Name); err != nil {
		t.Error(err)
	}

	_, err = c.FindVlanInterface(name)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
