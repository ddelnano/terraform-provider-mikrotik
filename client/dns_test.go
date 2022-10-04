package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDnsRecord_onNonExistantDnsRecord(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "dns record does not exist"
	_, err := c.FindDnsRecord(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for dns record `%s`, instead error was nil.", name)
	}
}

func TestAddFindDeleteDnsRecord(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	recordName := "new_record"
	record := &DnsRecord{
		Name:    recordName,
		Address: "10.10.10.200",
		Ttl:     300,
		Comment: "new record from test",
	}

	created, err := c.Add(record)
	require.NoError(t, err)

	findRecord := &DnsRecord{}
	findRecord.Name = recordName
	found, err := c.Find(findRecord)
	require.NoError(t, err)
	assert.Implements(t, (*Resource)(nil), found)

	assert.Equal(t, created, found)
	err = c.Delete(found.(Resource))
	assert.NoError(t, err)

	_, err = c.Find(findRecord)
	assert.Error(t, err)
	assert.IsType(t, &NotFound{}, err)
}
