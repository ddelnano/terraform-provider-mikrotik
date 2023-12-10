package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindInterfaceVeth_onNonExistantInterfaceVeth(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	name := "Interface Veth does not exist"
	_, err := c.FindInterfaceVeth(name)

	require.Truef(t, IsNotFoundError(err),
		"Expecting to receive NotFound error for Interface Veth %q.", name)
}

func TestAddFindDeleteInterfaceVeth(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	expectedIface := &InterfaceVeth{
		Name:      "veth-test-interface",
		Disabled:   false,
		Address: 	"192.168.88.2/24",
		Gateway:    "192.168.88.1",
		Comment:    "new interface from test",
		Running:	true,
	}

	iface, err := c.AddInterfaceVeth(&InterfaceVeth{
		Name:       expectedIface.Name,
		Disabled:   expectedIface.Disabled,
		Address: 	expectedIface.Address,
		Gateway:    expectedIface.Gateway,
		Comment:    expectedIface.Comment,
	})
	require.NoError(t, err)

	expectedIface.Id = iface.Id

	foundInterface, err := c.FindInterfaceVeth(expectedIface.Name)
	require.NoError(t, err)
	assert.Equal(t, expectedIface, foundInterface)

	expectedIface.Name = expectedIface.Name + "updated"
	expectedIface.Address = "192.168.188.2/24"
	expectedIface.Gateway = "192.168.188.1"
	expectedIface.Comment = expectedIface.Comment + " with updated comment"

	updatedIface, err := c.UpdateInterfaceVeth(expectedIface)
	require.NoError(t, err)
	assert.Equal(t, expectedIface, updatedIface)
	// cleanup
	err = c.DeleteInterfaceVeth(iface.Name)
	assert.NoError(t, err)

	_, err = c.FindInterfaceVeth(expectedIface.Name)
	assert.Error(t, err)
}