package client

import (
	"errors"
	"reflect"
	"testing"
)

func TestFindInterfacePeer_onNonExistantInterfacePeer(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	id := "Interface peer does not exist"
	_, err := c.FindInterfacePeer(id)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for Interface peer `%s`, instead error was nil.", id)
	}
}

func TestAddFindDeleteInterfacePeer(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	id := "new_interface_peer"
	interfacePeer := &InterfacePeer{
		Id:       id,
		Disabled: false,
		Comment:  "new interface from test",
	}

	created, err := c.Add(interfacePeer)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer func() {
		err = c.Delete(interfacePeer)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expected := &NotFound{}
		if _, err := c.Find(interfacePeer); err == nil || !errors.As(err, &expected) {
			t.Error(err)
		}
	}()
	findInterface := &InterfacePeer{}
	findInterface.Id = id
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
