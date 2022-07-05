package client

import (
	"github.com/go-routeros/routeros"
)

type DhcpLease struct {
	Id          string `mikrotik:".id"`
	Address     string `mikrotik:"address"`
	MacAddress  string `mikrotik:"mac-address"`
	Comment     string `mikrotik:"comment"`
	BlockAccess bool   `mikrotik:"block-access"`
	Dynamic     bool   // TODO:  don't see this listed as a param https://wiki.mikrotik.com/wiki/Manual:IP/DHCP_Server, but our docs list it as one
	Hostname    string
}

var dhcpLeaseWrapper *resourceWrapper = &resourceWrapper{
	idField:       ".id",
	idFieldDelete: ".id",
	actionsMap: map[string]string{
		"add":    "/ip/dhcp-server/lease/add",
		"find":   "/ip/dhcp-server/lease/print",
		"list":   "/ip/dhcp-server/lease/print",
		"update": "/ip/dhcp-server/lease/set",
		"delete": "/ip/dhcp-server/lease/remove",
	},
	targetStruct:          &DhcpLease{},
	addIDExtractorFunc:    func(r *routeros.Reply, _ interface{}) string { return r.Done.Map["ret"] },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*DhcpLease).Id },
}

func (client Mikrotik) AddDhcpLease(l *DhcpLease) (*DhcpLease, error) {
	r, err := dhcpLeaseWrapper.Add(l, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DhcpLease), nil
}

func (client Mikrotik) ListDhcpLeases() ([]DhcpLease, error) {
	r, err := dhcpLeaseWrapper.List(client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.([]DhcpLease), nil

}

func (client Mikrotik) FindDhcpLease(id string) (*DhcpLease, error) {
	r, err := dhcpLeaseWrapper.Find(id, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DhcpLease), nil
}

func (client Mikrotik) UpdateDhcpLease(l *DhcpLease) (*DhcpLease, error) {
	r, err := dhcpLeaseWrapper.Update(l, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*DhcpLease), nil

}

func (client Mikrotik) DeleteDhcpLease(id string) error {
	return dhcpLeaseWrapper.Delete(id, client.getMikrotikClient)
}
