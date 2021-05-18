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
