package client

import (
	"reflect"
	"testing"
)

func TestAddIpAddressAndDeleteIpAddress(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	address := "1.1.1.1/24"
	comment := "terraform-acc-test"
	disabled := false
	network := "1.1.1.0"
	ifname := "ether1"
	updatedComment := "terraform acc test updated"

	expectedIpAddress := &IpAddress{
		Address:   address,
		Comment:   comment,
		Disabled:  disabled,
		Interface: ifname,
		Network:   network,
	}

	ipaddr, err := c.AddIpAddress(expectedIpAddress)

	if err != nil {
		t.Errorf("Error creating an ip address with: %v", err)
	}

	expectedIpAddress.Id = ipaddr.Id

	if !reflect.DeepEqual(ipaddr, expectedIpAddress) {
		t.Errorf("The ip address does not match what we expected. actual: %v expected: %v", ipaddr, expectedIpAddress)
	}

	expectedIpAddress.Comment = updatedComment
	ipaddr, err = c.UpdateIpAddress(expectedIpAddress)

	if err != nil {
		t.Errorf("Error updating an ip address with: %v", err)
	}
	if !reflect.DeepEqual(ipaddr, expectedIpAddress) {
		t.Errorf("The ip address does not match what we expected. actual: %v expected: %v", ipaddr, expectedIpAddress)
	}

	foundIpAddress, err := c.FindIpAddress(ipaddr.Id)

	if err != nil {
		t.Errorf("Error getting ip address with: %v", err)
	}

	if !reflect.DeepEqual(ipaddr, foundIpAddress) {
		t.Errorf("Created ip address and found ip address do not match. actual: %v expected: %v", foundIpAddress, ipaddr)
	}

	err = c.DeleteIpAddress(ipaddr.Id)

	if err != nil {
		t.Errorf("Error deleting ip address with: %v", err)
	}
}
