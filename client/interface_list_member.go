package client

import (
	"github.com/go-routeros/routeros"
)

// InterfaceListMember manages an interface list's members
type InterfaceListMember struct {
	Id        string `mikrotik:".id"`
	Interface string `mikrotik:"interface"`
	List      string `mikrotik:"list"`
}

var _ Resource = (*InterfaceListMember)(nil)

func (b *InterfaceListMember) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/list/member/add",
		Find:   "/interface/list/member/print",
		Update: "/interface/list/member/set",
		Delete: "/interface/list/member/remove",
	}[a]
}

func (b *InterfaceListMember) IDField() string {
	return ".id"
}

func (b *InterfaceListMember) ID() string {
	return b.Id
}

func (b *InterfaceListMember) SetID(id string) {
	b.Id = id
}

func (b *InterfaceListMember) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *InterfaceListMember) DeleteField() string {
	return "numbers"
}

func (b *InterfaceListMember) DeleteFieldValue() string {
	return b.Id
}

// Typed wrappers
func (c Mikrotik) AddInterfaceListMember(r *InterfaceListMember) (*InterfaceListMember, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceListMember), nil
}

func (c Mikrotik) UpdateInterfaceListMember(r *InterfaceListMember) (*InterfaceListMember, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceListMember), nil
}

func (c Mikrotik) FindInterfaceListMember(id string) (*InterfaceListMember, error) {
	res, err := c.Find(&InterfaceListMember{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceListMember), nil
}

func (c Mikrotik) DeleteInterfaceListMember(id string) error {
	return c.Delete(&InterfaceListMember{Id: id})
}
