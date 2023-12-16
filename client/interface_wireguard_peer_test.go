package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindInterfaceWireguardPeer_onNonExistantInterfacePeer(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	peerID := "Interface peer does not exist"
	_, err := c.FindInterfaceWireguardPeer(peerID)

	require.Truef(t, IsNotFoundError(err),
		"Expecting to receive NotFound error for Interface peer `%q`, instead error was nil.", peerID)
}

func TestInterfaceWireguardPeer_Crud(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
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

	createdInterface, err := c.Add(interfaceWireguard)
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
		Interface:      createdInterface.(*InterfaceWireguard).Name,
		Disabled:       false,
		AllowedAddress: "0.0.0.0/0",
		Comment:        "new interface from test",
		PublicKey:      "/yZWgiYAgNNSy7AIcxuEewYwOVPqJJRKG90s9ypwfiM=",
	}

	created, err := c.Add(interfaceWireguardPeer)
	require.NoError(t, err)
	defer func() {
		err = c.Delete(interfaceWireguardPeer)
		assert.NoError(t, err)
	}()

	findPeer := &InterfaceWireguardPeer{}
	findPeer.Id = created.(*InterfaceWireguardPeer).Id
	foundPeer, err := c.Find(findPeer)
	require.NoError(t, err)

	assert.Equal(t, created, foundPeer)
}
