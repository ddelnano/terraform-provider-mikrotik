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

var dnsRecordWrapper *resourceWrapper = &resourceWrapper{
	idField: "name",
	actionsMap: map[string]string{
		"add":    "/ip/dns/static/add",
		"find":   "/ip/dns/static/print",
		"list":   "/ip/dns/static/print",
		"update": "/ip/dns/static/set",
		"delete": "/ip/dns/static/remove",
	},
	targetStruct:          &DnsRecord{},
	addIDExtractorFunc:    func(r *routeros.Reply) string { return r.Done.Map["ret"] },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*DnsRecord).Id },
}

func (client Mikrotik) AddDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	r, err := dnsRecordWrapper.Add(d, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil
}

func (client Mikrotik) FindDnsRecord(name string) (*DnsRecord, error) {
	r, err := dnsRecordWrapper.Find(name, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil

}

func (client Mikrotik) UpdateDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	r, err := dnsRecordWrapper.Update(d, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DnsRecord), nil

}

func (client Mikrotik) DeleteDnsRecord(id string) error {
	return dnsRecordWrapper.Delete(id, client.getMikrotikClient)
}
