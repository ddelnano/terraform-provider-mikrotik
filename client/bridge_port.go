package client

import (
	"github.com/go-routeros/routeros"
)

// BridgePort defines port-in-bridge association
type BridgePort struct {
	Id        string `mikrotik:".id" codegen:"id,mikrotikID"`
	Bridge    string `mikrotik:"bridge" codegen:"bridge"`
	Interface string `mikrotik:"interface" codegen:"interface"`
	PVId      int    `mikrotik:"pvid" codegen:"pvid"`
	Comment   string `mikrotik:"comment" codegen:"comment"`
}

var _ Resource = (*BridgePort)(nil)

func (b *BridgePort) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/bridge/port/add",
		Find:   "/interface/bridge/port/print",
		Update: "/interface/bridge/port/set",
		Delete: "/interface/bridge/port/remove",
	}[a]
}

func (b *BridgePort) IDField() string {
	return ".id"
}

func (b *BridgePort) ID() string {
	return b.Id
}

func (b *BridgePort) SetID(id string) {
	b.Id = id
}

func (b *BridgePort) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BridgePort) DeleteField() string {
	return "numbers"
}

func (b *BridgePort) DeleteFieldValue() string {
	return b.Id
}

// Typed wrappers
func (c Mikrotik) AddBridgePort(r *BridgePort) (*BridgePort, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BridgePort), nil
}

func (c Mikrotik) UpdateBridgePort(r *BridgePort) (*BridgePort, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BridgePort), nil
}

func (c Mikrotik) FindBridgePort(id string) (*BridgePort, error) {
	res, err := c.Find(&BridgePort{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*BridgePort), nil
}

func (c Mikrotik) DeleteBridgePort(id string) error {
	return c.Delete(&BridgePort{Id: id})
}
