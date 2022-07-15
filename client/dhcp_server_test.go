package client

import "testing"

func TestAddDhcpServerUpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "myserver"
	disabled := true
	dhcpServer, err := c.AddDhcpServer(&DhcpServer{
		Name:     name,
		Disabled: disabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	foundServer, err := c.FindDhcpServer(name)
	if err != nil {
		t.Fatal(err)
	}

	if foundServer.Name != name {
		t.Errorf("expected server name to be %q, got %q", name, foundServer.Name)
	}

	dhcpServer.Name = dhcpServer.Name + "updated"
	updatedServer, err := c.UpdateDhcpServer(dhcpServer)
	if err != nil {
		t.Error(err)
	}

	if updatedServer.Name != dhcpServer.Name {
		t.Errorf("expected name to be %q, got %q", dhcpServer.Name, updatedServer.Name)
	}

	// cleanup
	if err := c.DeleteDhcpServer(dhcpServer.Id); err != nil {
		t.Error(err)
	}

	_, err = c.FindDhcpServer(name)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
