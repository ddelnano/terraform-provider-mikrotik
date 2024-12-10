package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDnsRecord_onNonExistantDnsRecord(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "dns record does not exist"
	_, err := c.FindDnsRecord(name)

	require.Truef(t, IsNotFoundError(err),
		"Expecting to receive NotFound error for dns record %q", name)
}

func TestDnsRecord_basic(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	recordName := "new_record"
	record := &DnsRecord{
		Name:    recordName,
		Address: "10.10.10.200",
		Ttl:     300,
		Comment: "new record from test",
	}

	created, err := c.Add(record)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	found, err := c.Find(&DnsRecord{Id: created.ID()})
	require.NoError(t, err)

	if !reflect.DeepEqual(created, found) {
		t.Error("expected created and found resources to be equal, but they don't")
	}

	created.(*DnsRecord).Comment = "updated comment"
	_, err = c.Update(created)
	require.NoError(t, err)
	found, err = c.Find(&DnsRecord{Id: created.ID()})
	require.NoError(t, err)
	assert.Equal(t, created, found)

	err = c.Delete(found)
	assert.NoError(t, err)

	_, err = c.Find(&DnsRecord{Id: created.ID()})
	require.Error(t, err)

	require.True(t, IsNotFoundError(err),
		"expected to get NotFound error")
}

func TestDns_Regexp(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	recordName := "new_record"
	record := &DnsRecord{
		Name:    recordName,
		Regexp:  ".*\\.domain\\.com",
		Address: "10.10.10.200",
		Ttl:     300,
		Comment: "new record from test",
	}

	_, err := c.Add(record)
	require.Error(t, err, "usage of 'name' and 'regexp' at the same type should result in error")

	regexRecord := &DnsRecord{
		Address: "10.10.10.201",
		Ttl:     300,
		Regexp:  ".+\\.domain\\.com",
		Comment: "new record from test",
	}
	regexCreated, err := c.Add(regexRecord)
	require.NoError(t, err)
	defer func() {
		_ = c.Delete(regexCreated)
	}()
	assert.Equal(t, regexRecord, regexCreated)
}
