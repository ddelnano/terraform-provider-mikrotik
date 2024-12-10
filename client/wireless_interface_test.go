package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWirelessInterface_basic(t *testing.T) {
	// This test is skipped, until we find a way to include required packages.
	//
	// Since RouterOS 7.13, 'wireless' package is separate from the main system package
	// and there is no easy way to install it in Docker during tests.
	// see https://help.mikrotik.com/docs/spaces/ROS/pages/40992872/Packages#Packages-RouterOSpackages
	SkipIfRouterOSV7OrLater(t, sysResources)

	randSuffix := RandomString()
	c := NewClient(GetConfigFromEnv())
	expected := &WirelessInterface{
		Name:            "wireless-" + randSuffix,
		SSID:            "ssid-" + randSuffix,
		MasterInterface: "*0",
	}
	created, err := c.AddWirelessInterface(expected)
	require.NoError(t, err)
	defer c.DeleteWirelessInterface(created.Id)

	assert.Equal(t, expected.Name, created.Name)
	assert.Equal(t, expected.SSID, created.SSID)
	assert.Equal(t, false, created.Disabled)

	created.Disabled = true
	created.Name = "wireless-updated-" + randSuffix
	updated, err := c.UpdateWirelessInterface(created)
	require.NoError(t, err)
	assert.Equal(t, created, updated)

	found, err := c.FindWirelessInterface(updated.Id)
	require.NoError(t, err)
	assert.Equal(t, updated, found)
}
