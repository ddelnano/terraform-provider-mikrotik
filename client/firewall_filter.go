package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/internal/types"
	"github.com/go-routeros/routeros"
)

// FirewallFilterRule defines /ip/firewall/filter rule
type FirewallFilterRule struct {
	Id               string             `mikrotik:".id"`
	Action           string             `mikrotik:"action"`
	Chain            string             `mikrotik:"chain"`
	Comment          string             `mikrotik:"comment"`
	ConnectionState  types.MikrotikList `mikrotik:"connection-state"`
	DestPort         string             `mikrotik:"dst-port"`
	InInterface      string             `mikrotik:"in-interface"`
	InInterfaceList  string             `mikrotik:"in-interface-list"`
	OutInterfaceList string             `mikrotik:"out-interface-list"`
	Protocol         string             `mikrotik:"protocol"`
}

var _ Resource = (*FirewallFilterRule)(nil)

func (b *FirewallFilterRule) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/firewall/filter/add",
		Find:   "/ip/firewall/filter/print",
		Update: "/ip/firewall/filter/set",
		Delete: "/ip/firewall/filter/remove",
	}[a]
}

func (b *FirewallFilterRule) IDField() string {
	return ".id"
}

func (b *FirewallFilterRule) ID() string {
	return b.Id
}

func (b *FirewallFilterRule) SetID(id string) {
	b.Id = id
}

func (b *FirewallFilterRule) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (c Mikrotik) AddFirewallFilterRule(r *FirewallFilterRule) (*FirewallFilterRule, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallFilterRule), nil
}

func (c Mikrotik) UpdateFirewallFilterRule(r *FirewallFilterRule) (*FirewallFilterRule, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallFilterRule), nil
}

func (c Mikrotik) FindFirewallFilterRule(id string) (*FirewallFilterRule, error) {
	res, err := c.Find(&FirewallFilterRule{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*FirewallFilterRule), nil
}

func (c Mikrotik) DeleteFirewallFilterRule(id string) error {
	return c.Delete(&FirewallFilterRule{Id: id})
}
