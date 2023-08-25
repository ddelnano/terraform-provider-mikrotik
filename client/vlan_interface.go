package client

import (
	"github.com/go-routeros/routeros"
)

// VlanInterface represents vlan interface resource
type VlanInterface struct {
	Id            string `mikrotik:".id"`
	Interface     string `mikrotik:"interface"`
	Mtu           int    `mikrotik:"mtu"`
	Name          string `mikrotik:"name"`
	Disabled      bool   `mikrotik:"disabled"`
	UseServiceTag bool   `mikrotik:"use-service-tag"`
	VlanId        int    `mikrotik:"vlan-id"`
}

var _ Resource = (*VlanInterface)(nil)

func (b *VlanInterface) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/vlan/add",
		Find:   "/interface/vlan/print",
		Update: "/interface/vlan/set",
		Delete: "/interface/vlan/remove",
	}[a]
}

func (b *VlanInterface) IDField() string {
	return ".id"
}

func (b *VlanInterface) ID() string {
	return b.Id
}

func (b *VlanInterface) SetID(id string) {
	b.Id = id
}

func (b *VlanInterface) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *VlanInterface) FindField() string {
	return "name"
}

func (b *VlanInterface) FindFieldValue() string {
	return b.Name
}

func (b *VlanInterface) DeleteField() string {
	return "numbers"
}

func (b *VlanInterface) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddVlanInterface(r *VlanInterface) (*VlanInterface, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*VlanInterface), nil
}

func (c Mikrotik) UpdateVlanInterface(r *VlanInterface) (*VlanInterface, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*VlanInterface), nil
}

func (c Mikrotik) FindVlanInterface(name string) (*VlanInterface, error) {
	res, err := c.Find(&VlanInterface{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*VlanInterface), nil
}

func (c Mikrotik) ListVlanInterface() ([]VlanInterface, error) {
	res, err := c.List(&VlanInterface{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]VlanInterface, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*VlanInterface))
	}

	return returnSlice, nil
}

func (c Mikrotik) DeleteVlanInterface(name string) error {
	return c.Delete(&VlanInterface{Name: name})
}
