package client

import "github.com/go-routeros/routeros"

// DhcpServerNetwork describes network configuration for DHCP server
type DhcpServerNetwork struct {
	Id        string `mikrotik:".id"`
	Comment   string `mikrotik:"comment"`
	Address   string `mikrotik:"address"`
	Netmask   string `mikrotik:"netmask"`
	Gateway   string `mikrotik:"gateway"`
	DnsServer string `mikrotik:"dns-server"`
}

var _ Resource = (*DhcpServerNetwork)(nil)

func (b *DhcpServerNetwork) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/dhcp-server/network/add",
		Find:   "/ip/dhcp-server/network/print",
		Update: "/ip/dhcp-server/network/set",
		Delete: "/ip/dhcp-server/network/remove",
	}[a]
}

func (b *DhcpServerNetwork) IDField() string {
	return ".id"
}

func (b *DhcpServerNetwork) ID() string {
	return b.Id
}

func (b *DhcpServerNetwork) SetID(id string) {
	b.Id = id
}

func (b *DhcpServerNetwork) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddDhcpServerNetwork(r *DhcpServerNetwork) (*DhcpServerNetwork, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServerNetwork), nil
}

func (c Mikrotik) UpdateDhcpServerNetwork(r *DhcpServerNetwork) (*DhcpServerNetwork, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServerNetwork), nil
}

func (c Mikrotik) FindDhcpServerNetwork(id string) (*DhcpServerNetwork, error) {
	res, err := c.Find(&DhcpServerNetwork{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*DhcpServerNetwork), nil
}

func (c Mikrotik) DeleteDhcpServerNetwork(id string) error {
	return c.Delete(&DhcpServerNetwork{Id: id})
}
