package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	expectedLease.Id = lease.Id
	assert.Equal(t, expectedLease, lease)

	expectedLease.Comment = updatedComment
	expectedLease.MacAddress = updatedMacaddress

	lease, err = c.UpdateDhcpLease(expectedLease)
	assert.NoError(t, err)
	assert.Equal(t, expectedLease, lease)

	foundLease, err := c.FindDhcpLease(lease.Id)
	assert.NoError(t, err)
	assert.Equal(t, lease, foundLease)

	err = c.DeleteDhcpLease(lease.Id)
	assert.NoError(t, err)
}

func TestFindDhcpLease_forNonExistantLease(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	leaseId := "Invalid id"
	_, err := c.FindDhcpLease(leaseId)

	assert.Error(t, err)
	assert.True(t, IsNotFoundError(err), "expected error to be of NotFound type")
}
