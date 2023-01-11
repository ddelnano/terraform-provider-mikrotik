package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
)

// BridgeVlan defines vlan filtering in bridge resource
type BridgeVlan struct {
	Id       string                `mikrotik:".id"`
	Bridge   string                `mikrotik:"bridge"`
	Tagged   types.MikrotikList    `mikrotik:"tagged"`
	Untagged types.MikrotikList    `mikrotik:"untagged"`
	VlanIds  types.MikrotikIntList `mikrotik:"vlan-ids"`
}

var _ Resource = (*BridgeVlan)(nil)

func (b *BridgeVlan) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/bridge/vlan/add",
		Find:   "/interface/bridge/vlan/print",
		Update: "/interface/bridge/vlan/set",
		Delete: "/interface/bridge/vlan/remove",
	}[a]
}

func (b *BridgeVlan) IDField() string {
	return ".id"
}

func (b *BridgeVlan) ID() string {
	return b.Id
}

func (b *BridgeVlan) SetID(id string) {
	b.Id = id
}

func (b *BridgeVlan) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (c Mikrotik) AddBridgeVlan(r *BridgeVlan) (*BridgeVlan, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BridgeVlan), nil
}

func (c Mikrotik) UpdateBridgeVlan(r *BridgeVlan) (*BridgeVlan, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BridgeVlan), nil
}

func (c Mikrotik) FindBridgeVlan(id string) (*BridgeVlan, error) {
	res, err := c.Find(&BridgeVlan{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*BridgeVlan), nil
}

func (c Mikrotik) DeleteBridgeVlan(id string) error {
	return c.Delete(&BridgeVlan{Id: id})
}
