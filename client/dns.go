package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
)

type DnsRecord struct {
	Id      string                 `mikrotik:".id" codegen:"id,mikrotikID"`
	Name    string                 `mikrotik:"name" codegen:"name,terraformID,required"`
	Ttl     types.MikrotikDuration `mikrotik:"ttl" codegen:"ttl"`
	Address string                 `mikrotik:"address" codegen:"address,required"`
	Comment string                 `mikrotik:"comment" codegen:"comment"`
}

func (d *DnsRecord) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/ip/dns/static/add",
		Find:   "/ip/dns/static/print",
		List:   "/ip/dns/static/print",
		Update: "/ip/dns/static/set",
		Delete: "/ip/dns/static/remove",
	}[action]
}

func (d *DnsRecord) IDField() string {
	return ".id"
}

func (d *DnsRecord) ID() string {
	return d.Id
}

func (d *DnsRecord) SetID(id string) {
	d.Id = id
}

func (d *DnsRecord) AfterAddHook(r *routeros.Reply) {
	d.Id = r.Done.Map["ret"]
}

func (d *DnsRecord) FindField() string {
	return "name"
}

func (d *DnsRecord) FindFieldValue() string {
	return d.Name
}

func (d *DnsRecord) DeleteField() string {
	return "numbers"
}

func (d *DnsRecord) DeleteFieldValue() string {
	return d.Id
}

func (client Mikrotik) AddDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	res, err := client.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*DnsRecord), nil
}

func (client Mikrotik) FindDnsRecord(name string) (*DnsRecord, error) {
	res, err := client.Find(&DnsRecord{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*DnsRecord), nil
}

func (client Mikrotik) UpdateDnsRecord(d *DnsRecord) (*DnsRecord, error) {
	res, err := client.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*DnsRecord), nil
}

func (client Mikrotik) DeleteDnsRecord(id string) error {
	return client.Delete(&DnsRecord{Id: id})
}
