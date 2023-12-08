package client

import (
	"reflect"
	"testing"

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

	name := "new_interface_veth"
	interfaceVeth := &InterfaceVeth{
		Name:       name,
		Disabled:   false,
		Address: 	"192.168.88.2/24",
		Gateway:    "192.168.88.1",
		Comment:    "new interface from test",
	}

	created, err := c.Add(interfaceVeth)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer func() {
		err = c.Delete(interfaceVeth)
		require.NoError(t, err)

		_, err := c.Find(interfaceVeth)
		require.True(t, IsNotFoundError(err), "expected to get NotFound error")
	}()

	findInterface := &InterfaceVeth{}
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