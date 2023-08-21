package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddIpv6AddressAndDeleteIpv6Address(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	address := "1:1:1:1:1:1:1:1/64"
	comment := "terraform-acc-test"
	disabled := false
	ifname := "ether1"
	updatedComment := "terraform acc test updated"

	expectedIpv6Address := &Ipv6Address{
		Address:   address,
		Comment:   comment,
		Disabled:  disabled,
		Interface: ifname,
	}

	ipv6addr, err := c.AddIpv6Address(expectedIpv6Address)
	require.NoError(t, err)

	expectedIpv6Address.Id = ipv6addr.Id
	assert.Equal(t, expectedIpv6Address, ipv6addr)

	expectedIpv6Address.Comment = updatedComment
	ipv6addr, err = c.UpdateIpv6Address(expectedIpv6Address)
	require.NoError(t, err)
	assert.Equal(t, expectedIpv6Address, ipv6addr)

	foundIpv6Address, err := c.FindIpv6Address(ipv6addr.Id)
	require.NoError(t, err)
	assert.Equal(t, ipv6addr, foundIpv6Address)

	err = c.DeleteIpv6Address(ipv6addr.Id)
	assert.NoError(t, err)
}
