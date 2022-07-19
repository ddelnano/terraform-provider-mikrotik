package client

import "testing"

func TestAddDhcpServerNetworkUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	netmask := "255.255.255.0"
	network := "192.168.99.0"
	dhcpServerNetwork, err := c.AddDhcpServerNetwork(&DhcpServerNetwork{
		Address: network + "/" + netmask,
		Netmask: netmask,
		Comment: "Created by terraform",
	})
	if err != nil {
		t.Fatal(err)
	}

	found, err := c.FindDhcpServerNetwork(dhcpServerNetwork.Id)
	if err != nil {
		t.Fatal(err)
	}

	if found.Address != dhcpServerNetwork.Address {
		t.Errorf("expected network address to be %q, got %q", dhcpServerNetwork.Address, found.Address)
	}

	dhcpServerNetwork.Comment = "updated network"
	updated, err := c.UpdateDhcpServerNetwork(dhcpServerNetwork)
	if err != nil {
		t.Error(err)
	}

	if updated.Comment != "updated network" {
		t.Errorf("expected comment to be %q, got %q", dhcpServerNetwork.Comment, updated.Comment)
	}

	// cleanup
	if err := c.DeleteDhcpServerNetwork(dhcpServerNetwork.Id); err != nil {
		t.Error(err)
	}

	_, err = c.FindDhcpServerNetwork(dhcpServerNetwork.Id)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
