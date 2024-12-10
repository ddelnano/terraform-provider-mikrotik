package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWirelessSecurityProfile_basic(t *testing.T) {
	// This test is skipped, until we find a way to include required packages.
	//
	// Since RouterOS 7.13, 'wireless' package is separate from the main system package
	// and there is no easy way to install it in Docker during tests.
	// see https://help.mikrotik.com/docs/spaces/ROS/pages/40992872/Packages#Packages-RouterOSpackages
	SkipIfRouterOSV7OrLater(t, sysResources)

	c := NewClient(GetConfigFromEnv())

	randSuffix := RandomString()
	expected := &WirelessSecurityProfile{
		Name: "test-profile-" + randSuffix,
		Mode: WirelessModeNone,
	}

	created, err := c.AddWirelessSecurityProfile(expected)
	require.NoError(t, err)
	defer c.DeleteWirelessSecurityProfile(created.Id)

	expected.Id = created.Id
	assert.Equal(t, expected, created)

	updated := &WirelessSecurityProfile{}
	*updated = *created
	updated.Name += "-updated"
	updated.Mode = WirelessModeDynamicKeys
	updated.AuthenticationTypes = []string{WirelessAuthenticationTypeWpa2Psk}
	updated.WPA2PreSharedKey = "1234567890"
	_, err = c.UpdateWirelessSecurityProfile(updated)
	require.NoError(t, err)

	found, err := c.FindWirelessSecurityProfile(updated.Id)
	require.NoError(t, err)

	assert.Equal(t, updated, found)

}
