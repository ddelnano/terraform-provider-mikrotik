package client

import (
	"fmt"
	"log"
)

type DnsRecord struct {
	Id      string `mikrotik:".id"`
	Name    string `mikrotik:"name"`
	Ttl     int    `mikrotik:"ttl,ttlToSeconds"`
	Address string `mikrotik:"address"`
}

func (client Mikrotik) AddDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := Marshal("/ip/dns/static/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /ip/dns/static/add returned %v", r)

	if err != nil {
		return nil, err
	}

	return client.FindDnsRecord(d.Name)
}

func (client Mikrotik) FindDnsRecord(name string) (*DnsRecord, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/dns/static/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Found dns record: %v", r)

	record := DnsRecord{}
	err = Unmarshal(*r, &record)

	if err != nil {
		return nil, err
	}

	if record.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("dns record `%s` not found", name))
	}

	return &record, nil
}

func (client Mikrotik) UpdateDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := Marshal("/ip/dns/static/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	return client.FindDnsRecord(d.Name)
}

func (client Mikrotik) DeleteDnsRecord(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}
	cmd := []string{"/ip/dns/static/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
