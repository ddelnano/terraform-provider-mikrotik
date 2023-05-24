package client

import (
	"errors"
	"reflect"
	"testing"
)

func TestFindInterfaceWireguard_onNonExistantInterfaceWireguard(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	name := "Interface wireguard does not exist"
	_, err := c.FindInterfaceWireguard(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for Interface wireguard `%s`, instead error was nil.", name)
	}
}

func TestAddFindDeleteInterfaceWireguard(t *testing.T) {
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

	created, err := c.Add(interfaceWireguard)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer func() {
		err = c.Delete(interfaceWireguard)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expected := &NotFound{}
		if _, err := c.Find(interfaceWireguard); err == nil || !errors.As(err, &expected) {
			t.Error(err)
		}
	}()
	findInterface := &InterfaceWireguard{}
	findInterface.Name = name
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
