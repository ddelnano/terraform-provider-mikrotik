package client

import (
	"reflect"
	"testing"
)

func TestAddIpv6AddressAndDeleteIpv6Address(t *testing.T) {
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

	if err != nil {
		t.Errorf("Error creating an ipv6 address with: %v", err)
	}

	expectedIpv6Address.Id = ipv6addr.Id

	if !reflect.DeepEqual(ipv6addr, expectedIpv6Address) {
		t.Errorf("The ipv6 address does not match what we expected. actual: %v expected: %v", ipv6addr, expectedIpv6Address)
	}

	expectedIpv6Address.Comment = updatedComment
	ipv6addr, err = c.UpdateIpv6Address(expectedIpv6Address)

	if err != nil {
		t.Errorf("Error updating an ipv6 address with: %v", err)
	}
	if !reflect.DeepEqual(ipv6addr, expectedIpv6Address) {
		t.Errorf("The ipv6 address does not match what we expected. actual: %v expected: %v", ipv6addr, expectedIpv6Address)
	}

	foundIpv6Address, err := c.FindIpv6Address(ipv6addr.Id)

	if err != nil {
		t.Errorf("Error getting ipv6 address with: %v", err)
	}

	if !reflect.DeepEqual(ipv6addr, foundIpv6Address) {
		t.Errorf("Created ipv6 address and found ipv6 address do not match. actual: %v expected: %v", foundIpv6Address, ipv6addr)
	}

	err = c.DeleteIpv6Address(ipv6addr.Id)

	if err != nil {
		t.Errorf("Error deleting ipv6 address with: %v", err)
	}
}
