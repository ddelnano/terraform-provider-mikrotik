package client

import (
	"github.com/go-routeros/routeros"
)

type IpAddress struct {
	Id        string `mikrotik:".id" codegen:"id,mikrotikID"`
	Address   string `mikrotik:"address" codegen:"address,required"`
	Comment   string `mikrotik:"comment" codegen:"comment"`
	Disabled  bool   `mikrotik:"disabled" codegen:"disabled"`
	Interface string `mikrotik:"interface" codegen:"interface,required"`
	Network   string `mikrotik:"network" codegen:"network,computed"`
}

var _ Resource = (*IpAddress)(nil)

func (b *IpAddress) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/address/add",
		Find:   "/ip/address/print",
		Update: "/ip/address/set",
		Delete: "/ip/address/remove",
	}[a]
}

func (b *IpAddress) IDField() string {
	return ".id"
}

func (b *IpAddress) ID() string {
	return b.Id
}

func (b *IpAddress) SetID(id string) {
	b.Id = id
}

func (b *IpAddress) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddIpAddress(r *IpAddress) (*IpAddress, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*IpAddress), nil
}

func (c Mikrotik) UpdateIpAddress(r *IpAddress) (*IpAddress, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*IpAddress), nil
}

func (c Mikrotik) FindIpAddress(id string) (*IpAddress, error) {
	res, err := c.Find(&IpAddress{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*IpAddress), nil
}

func (client Mikrotik) ListIpAddress() ([]IpAddress, error) {
	res, err := client.List(&IpAddress{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]IpAddress, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*IpAddress))
	}

	return returnSlice, nil
}

func (c Mikrotik) DeleteIpAddress(id string) error {
	return c.Delete(&IpAddress{Id: id})
}
