package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWirelessSecurityProgile_basic(t *testing.T) {
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
	_, err = c.UpdateWirelessSecurityProfile(updated)
	require.NoError(t, err)

	found, err := c.FindWirelessSecurityProfile(updated.Id)
	require.NoError(t, err)

	assert.Equal(t, updated, found)

}
