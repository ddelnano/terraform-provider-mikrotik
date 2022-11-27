package client

import (
	"errors"
	"reflect"
	"testing"
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
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	findRecord := &DnsRecord{}
	findRecord.Name = recordName
	found, err := c.Find(findRecord)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if _, ok := found.(Resource); !ok {
		t.Error("expected found resource to implement Resource interface, but it doesn't")
		return
	}
	if !reflect.DeepEqual(created, found) {
		t.Error("expected created and found resources to be equal, but they don't")
	}
	err = c.Delete(found.(Resource))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = c.Find(findRecord)
	if err == nil {
		t.Errorf("expected error, got nothing")
		return
	}

	target := &NotFound{}
	if !errors.As(err, &target) {
		t.Errorf("expected error to be of type %T, got %T", &NotFound{}, err)
	}
}
