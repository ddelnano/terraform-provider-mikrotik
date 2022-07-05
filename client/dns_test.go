package client

import "testing"

func TestFindDnsRecord_onNonExistantDnsRecord(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "dns record does not exist"
	_, err := c.FindDnsRecord(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for dns record `%s`, instead error was nil.", name)
	}
}

func TestAddDnsRecordAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	r, err := c.AddDnsRecord(&DnsRecord{
		Name:    "test-record",
		Address: "192.168.10.22",
	})

	t.Log(r)
	if err != nil {
		t.Error(err)
	}

	err = c.DeleteDnsRecord(r.Id)
	if err != nil {
		t.Error(err)
	}

}
