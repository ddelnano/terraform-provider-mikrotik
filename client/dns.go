package client

import (
	"github.com/go-routeros/routeros"
)

type DnsRecord struct {
	Id      string `mikrotik:".id"`
	Name    string `mikrotik:"name"`
	Ttl     int    `mikrotik:"ttl,ttlToSeconds"`
	Address string `mikrotik:"address"`
	Comment string `mikrotik:"comment"`
}

var _ Resource = (*DnsRecord)(nil)

func (d *DnsRecord) ActionCommand(action string) string {
	return map[string]string{
		"add":    "/ip/dns/static/add",
		"find":   "/ip/dns/static/print",
		"list":   "/ip/dns/static/print",
		"update": "/ip/dns/static/set",
		"delete": "/ip/dns/static/remove",
	}[action]
}

func (d *DnsRecord) AddIDExtractionFunc(r *routeros.Reply) string {
	return d.Name
}

func (d *DnsRecord) IDField() string {
	return "name"
}

func (d *DnsRecord) DeleteIDField() string {
	return "numbers"
}

func (d *DnsRecord) SetID(id string) {
	d.Name = id
}

func (d *DnsRecord) ID() string {
	return d.Name
}

func (d *DnsRecord) SetDeleteID(fieldName string) {
	d.Name = fieldName
}

func (client Mikrotik) AddDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	r, err := client.Add(d)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil
}

func (client Mikrotik) FindDnsRecord(name string) (*DnsRecord, error) {
	d := &DnsRecord{Name: name}
	r, err := client.Find(d)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil
}

func (client Mikrotik) UpdateDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	r, err := client.Update(d)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil
}

func (client Mikrotik) DeleteDnsRecord(id string) error {
	d := &DnsRecord{}
	d.SetDeleteID(id)
	return client.Delete(d)
}
