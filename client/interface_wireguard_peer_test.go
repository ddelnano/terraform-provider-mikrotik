package client

import (
	"reflect"
	"testing"
)

func TestFindInterfaceWireguardPeer_onNonExistantInterfacePeer(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	id := "Interface peer does not exist"
	_, err := c.FindInterfaceWireguardPeer(id)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for Interface peer `%s`, instead error was nil.", id)
	}
}

func TestWireguardInterfacePeer_Crud(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	name := "new_interface_wireguard"
	interfaceWireguard := &InterfaceWireguard{
		Name:       name,
		Disabled:   false,
		ListenPort: 10000,
		Mtu:        10001,
		PrivateKey: "YOi0P0lTTiN8hAQvuRET23Srb+U7C52iOZokj0CCSkM=",
		Comment:    "new interface from test",
	}

	created_interface, err := c.Add(interfaceWireguard)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer func() {
		err = c.Delete(interfaceWireguard)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}()

	interfaceWireguardPeer := &InterfaceWireguardPeer{
		Interface: created_interface.(*InterfaceWireguard).Name,
		Disabled:  false,
		Comment:   "new interface from test",
	}

	created, err := c.Add(interfaceWireguardPeer)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer func() {
		err = c.Delete(interfaceWireguardPeer)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}()
	findInterface := &InterfaceWireguardPeer{}
	findInterface.Interface = created_interface.(*InterfaceWireguard).Name
	found, err := c.Find(findInterface)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if _, ok := found.(Resource); !ok {
		t.Error("expected found resource to implement Resource interface, but it doesn't")
		return
	}
	if !reflect.DeepEqual(created, found) {
		t.Error("expected created and found resources to be equal, but they don't")
	}
}
