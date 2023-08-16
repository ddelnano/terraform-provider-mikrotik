package client

import (
	"github.com/go-routeros/routeros"
)

// Ipv6Address defines resource
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

var _ Resource = (*Ipv6Address)(nil)

func (b *Ipv6Address) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ipv6/address/add",
		Find:   "/ipv6/address/print",
		Update: "/ipv6/address/set",
		Delete: "/ipv6/address/remove",
	}[a]
}

func (b *Ipv6Address) IDField() string {
	return ".id"
}

func (b *Ipv6Address) ID() string {
	return b.Id
}

func (b *Ipv6Address) SetID(id string) {
	b.Id = id
}

func (b *Ipv6Address) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddIpv6Address(r *Ipv6Address) (*Ipv6Address, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*Ipv6Address), nil
}

func (c Mikrotik) UpdateIpv6Address(r *Ipv6Address) (*Ipv6Address, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*Ipv6Address), nil
}

func (c Mikrotik) ListIpv6Address() ([]Ipv6Address, error) {
	res, err := c.List(&Ipv6Address{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]Ipv6Address, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*Ipv6Address))
	}

	return returnSlice, nil
}

func (c Mikrotik) FindIpv6Address(id string) (*Ipv6Address, error) {
	res, err := c.Find(&Ipv6Address{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*Ipv6Address), nil
}

func (c Mikrotik) DeleteIpv6Address(id string) error {
	return c.Delete(&Ipv6Address{Id: id})
}
