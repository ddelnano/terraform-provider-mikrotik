package client

import (
	"github.com/go-routeros/routeros"
)

type Ipv6Address struct {
	Id        string `mikrotik:".id"`
	Address   string `mikrotik:"address"`
	Advertise bool   `mikrotik:"advertise"`
	Comment   string `mikrotik:"comment"`
	Disabled  bool   `mikrotik:"disabled"`
	Eui64     bool   `mikrotik:"eui-64"`
	FromPool  string `mikrotik:"from-pool"`
	Interface string `mikrotik:"interface"`
	NoDad     bool   `mikrotik:"no-dad"`
}

var ipv6Wrapper *resourceWrapper = &resourceWrapper{
	idField:       ".id",
	idFieldDelete: ".id",
	actionsMap: map[string]string{
		"add":    "/ipv6/address/add",
		"find":   "/ipv6/address/print",
		"list":   "/ipv6/address/print",
		"update": "/ipv6/address/set",
		"delete": "/ipv6/address/remove",
	},
	targetStruct:          &Ipv6Address{},
	addIDExtractorFunc:    func(r *routeros.Reply, _ interface{}) string { return r.Done.Map["ret"] },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*Ipv6Address).Id },
}

func (client Mikrotik) AddIpv6Address(addr *Ipv6Address) (*Ipv6Address, error) {
	r, err := ipv6Wrapper.Add(addr, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}
	return r.(*Ipv6Address), nil
}

func (client Mikrotik) FindIpv6Address(id string) (*Ipv6Address, error) {
	r, err := ipv6Wrapper.Find(id, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}
	return r.(*Ipv6Address), nil
}

func (client Mikrotik) ListIpv6Address() ([]Ipv6Address, error) {
	r, err := ipv6Wrapper.List(client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.([]Ipv6Address), nil
}

func (client Mikrotik) UpdateIpv6Address(addr *Ipv6Address) (*Ipv6Address, error) {
	r, err := ipv6Wrapper.Update(addr, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Ipv6Address), nil
}

func (client Mikrotik) DeleteIpv6Address(id string) error {
	return ipv6Wrapper.Delete(id, client.getMikrotikClient)

}
