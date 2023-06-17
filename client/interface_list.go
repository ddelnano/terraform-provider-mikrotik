package client

import (
	"github.com/go-routeros/routeros"
)

// InterfaceList manages a list of interfaces
type InterfaceList struct {
	Id      string `mikrotik:".id"`
	Comment string `mikrotik:"comment"`
	Name    string `mikrotik:"name"`
}

var _ Resource = (*InterfaceList)(nil)

func (b *InterfaceList) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/list/add",
		Find:   "/interface/list/print",
		Update: "/interface/list/set",
		Delete: "/interface/list/remove",
	}[a]
}

func (b *InterfaceList) IDField() string {
	return ".id"
}

func (b *InterfaceList) ID() string {
	return b.Id
}

func (b *InterfaceList) SetID(id string) {
	b.Id = id
}

// Uncomment extra methods to satisfy more interfaces

func (b *InterfaceList) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *InterfaceList) FindField() string {
	return "name"
}

func (b *InterfaceList) FindFieldValue() string {
	return b.Name
}

func (b *InterfaceList) DeleteField() string {
	return "numbers"
}

func (b *InterfaceList) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddInterfaceList(r *InterfaceList) (*InterfaceList, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceList), nil
}

func (c Mikrotik) UpdateInterfaceList(r *InterfaceList) (*InterfaceList, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceList), nil
}

func (c Mikrotik) FindInterfaceList(name string) (*InterfaceList, error) {
	res, err := c.Find(&InterfaceList{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceList), nil
}

func (c Mikrotik) DeleteInterfaceList(name string) error {
	return c.Delete(&InterfaceList{Name: name})
}
