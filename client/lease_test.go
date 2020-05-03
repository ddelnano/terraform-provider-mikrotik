package client

import (
	"reflect"
	"strings"
	"testing"
)

func TestAddLeaseAndDeleteLease(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	address := "1.1.1.1"
	macaddress := "11:11:11:11:11:11"
	comment := "terraform-acc-test"
	expectedLease := &DhcpLease{
		Address:    address,
		MacAddress: macaddress,
		Comment:    comment,
	}
	lease, err := c.AddDhcpLease(
		address,
		macaddress,
		comment,
	)

	if err != nil {
		t.Errorf("Error creating a lease with: %v", err)
	}

	if len(lease.Id) < 1 {
		t.Errorf("The created lease does not have an Id: %v", lease)
	}

	if strings.Compare(lease.Address, expectedLease.Address) != 0 {
		t.Errorf("The lease address fields do not match. actual: %v expected: %v", lease.Address, expectedLease.Address)
	}

	if strings.Compare(lease.MacAddress, expectedLease.MacAddress) != 0 {
		t.Errorf("The lease MacAddress fields do not match. actual: %v expected: %v", lease.MacAddress, expectedLease.MacAddress)
	}

	if strings.Compare(lease.Comment, expectedLease.Comment) != 0 {
		t.Errorf("The lease Comment fields do not match. actual: %v expected: %v", lease.Comment, expectedLease.Comment)
	}

	foundLease, err := c.FindDhcpLease(lease.Id)

	if err != nil {
		t.Errorf("Error getting lease with: %v", err)
	}

	if !reflect.DeepEqual(lease, foundLease) {
		t.Errorf("Created lease and found lease do not match. actual: %v expected: %v", foundLease, lease)
	}

	err = c.DeleteDhcpLease(lease.Id)

	if err != nil {
		t.Errorf("Error deleting lease with: %v", err)
	}
}
