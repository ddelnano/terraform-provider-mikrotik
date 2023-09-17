package client

import (
	"github.com/go-routeros/routeros"
)

// Bridge defines /bridge resource
type Bridge struct {
	Id            string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name          string `mikrotik:"name" codegen:"name,required,terraformID"`
	FastForward   bool   `mikrotik:"fast-forward" codegen:"fast_forward"`
	VlanFiltering bool   `mikrotik:"vlan-filtering" codegen:"vlan_filtering"`
	Comment       string `mikrotik:"comment" codegen:"comment"`
}

var _ Resource = (*Bridge)(nil)

func (b *Bridge) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/bridge/add",
		Find:   "/interface/bridge/print",
		Update: "/interface/bridge/set",
		Delete: "/interface/bridge/remove",
	}[a]
}

func (b *Bridge) IDField() string {
	return ".id"
}

func (b *Bridge) ID() string {
	return b.Id
}

func (b *Bridge) SetID(id string) {
	b.Id = id
}

func (b *Bridge) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *Bridge) FindField() string {
	return "name"
}

func (b *Bridge) FindFieldValue() string {
	return b.Name
}

func (b *Bridge) DeleteField() string {
	return "numbers"
}

func (b *Bridge) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBridge(r *Bridge) (*Bridge, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*Bridge), nil
}

func (c Mikrotik) UpdateBridge(r *Bridge) (*Bridge, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*Bridge), nil
}

func (c Mikrotik) FindBridge(name string) (*Bridge, error) {
	res, err := c.Find(&Bridge{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*Bridge), nil
}

func (c Mikrotik) DeleteBridge(name string) error {
	return c.Delete(&Bridge{Name: name})
}
