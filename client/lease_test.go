package client

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddLeaseAndDeleteLease(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	address := "1.1.1.1"
	macaddress := "11:11:11:11:11:11"
	comment := "terraform-acc-test"
	blocked := false
	updatedMacaddress := "11:11:11:11:11:12"
	updatedComment := "terraform acc test updated"

	expectedLease := &DhcpLease{
		Address:     address,
		MacAddress:  macaddress,
		Comment:     comment,
		BlockAccess: blocked,
	}
	lease, err := c.AddDhcpLease(expectedLease)

	if err != nil {
		t.Errorf("Error creating a lease with: %v", err)
	}

	expectedLease.Id = lease.Id

	if !reflect.DeepEqual(lease, expectedLease) {
		t.Errorf("The dhcp lease does not match what we expected. actual: %v expected: %v", lease, expectedLease)
	}

	expectedLease.Comment = updatedComment
	expectedLease.MacAddress = updatedMacaddress
	lease, err = c.UpdateDhcpLease(expectedLease)

	if err != nil {
		t.Errorf("Error updating a lease with: %v", err)
	}
	if !reflect.DeepEqual(lease, expectedLease) {
		t.Errorf("The dhcp lease does not match what we expected. actual: %v expected: %v", lease, expectedLease)
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

func TestFindDhcpLease_forNonExistantLease(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	leaseId := "Invalid id"
	_, err := c.FindDhcpLease(leaseId)

	expectedErrStr := fmt.Sprintf("dhcp lease `%s` not found", leaseId)
	if err == nil || err.Error() != expectedErrStr {
		t.Errorf("client should have received error indicating the following dns record `%s` was not found. Instead error was nil", leaseId)
	}
}
