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

	findInterface := &InterfaceWireguard{}
	findInterface.Id = created.(InterfaceWireguard).Id
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
	err = c.Delete(found.(Resource))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = c.Find(findInterface)
	if err == nil {
		t.Errorf("expected error, got nothing")
		return
	}

	target := &NotFound{}
	if !errors.As(err, &target) {
		t.Errorf("expected error to be of type %T, got %T", &NotFound{}, err)
	}
}
