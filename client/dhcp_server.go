package client

import (
	"github.com/go-routeros/routeros"
)

// DhcpServer represents DHCP server resource
type DhcpServer struct {
	Id            string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name          string `mikrotik:"name" codegen:"name,terraformID,required"`
	Disabled      bool   `mikrotik:"disabled" codegen:"disabled"`
	AddArp        bool   `mikrotik:"add-arp" codegen:"add_arp"`
	AddressPool   string `mikrotik:"address-pool" codegen:"address_pool"`
	Authoritative string `mikrotik:"authoritative" codegen:"authoritative"`
	Interface     string `mikrotik:"interface" codegen:"interface"`
	LeaseScript   string `mikrotik:"lease-script" codegen:"lease_script"`
}

var _ Resource = (*DhcpServer)(nil)

func (b *DhcpServer) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/dhcp-server/add",
		Find:   "/ip/dhcp-server/print",
		Update: "/ip/dhcp-server/set",
		Delete: "/ip/dhcp-server/remove",
	}[a]
}

func (b *DhcpServer) IDField() string {
	return ".id"
}

func (b *DhcpServer) ID() string {
	return b.Id
}

func (b *DhcpServer) SetID(id string) {
	b.Id = id
}

func (b *DhcpServer) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *DhcpServer) FindField() string {
	return "name"
}

func (b *DhcpServer) FindFieldValue() string {
	return b.Name
}

func (b *DhcpServer) DeleteField() string {
	return "numbers"
}

func (b *DhcpServer) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddDhcpServer(r *DhcpServer) (*DhcpServer, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServer), nil
}

func (c Mikrotik) UpdateDhcpServer(r *DhcpServer) (*DhcpServer, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServer), nil
}

func (c Mikrotik) FindDhcpServer(name string) (*DhcpServer, error) {
	res, err := c.Find(&DhcpServer{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServer), nil
}

func (c Mikrotik) DeleteDhcpServer(name string) error {
	return c.Delete(&DhcpServer{Name: name})
}
