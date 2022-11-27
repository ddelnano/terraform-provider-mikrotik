package client

import "github.com/go-routeros/routeros"

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
