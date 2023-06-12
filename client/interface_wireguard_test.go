package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindInterfaceWireguard_onNonExistantInterfaceWireguard(t *testing.T) {
	SkipInterfaceWireguardIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	name := "Interface wireguard does not exist"
	_, err := c.FindInterfaceWireguard(name)

	require.Truef(t, IsNotFoundError(err),
		"Expecting to receive NotFound error for Interface wireguard %q.", name)
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
		require.NoError(t, err)

		_, err := c.Find(interfaceWireguard)
		require.True(t, IsNotFoundError(err), "expected to get NotFound error")
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
