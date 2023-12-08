package client

import (
	"github.com/go-routeros/routeros"
)

type InterfaceVeth struct {
	Id         string `mikrotik:".id"`
	Name       string `mikrotik:"name"`
	Comment    string `mikrotik:"comment"`
	Disabled   bool   `mikrotik:"disabled"`
	Address	   string `mikrotik:"address"`
	Gateway	   string `mikrotik:"gateway"`
	Gateway6   string `mikrotik:"gateway6"`
	Running    bool   `mikrotik:"running,readonly"`    //read only property
}

func (i *InterfaceVeth) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/veth/add",
		Find:   "/interface/veth/print",
		List:   "/interface/veth/print",
		Update: "/interface/veth/set",
		Delete: "/interface/veth/remove",
	}[action]
}

func (i *InterfaceVeth) IDField() string {
	return ".id"
}

func (i *InterfaceVeth) ID() string {
	return i.Id
}

func (i *InterfaceVeth) SetID(id string) {
	i.Id = id
}

func (i *InterfaceVeth) AfterAddHook(r *routeros.Reply) {
	i.Id = r.Done.Map["ret"]
}

func (i *InterfaceVeth) FindField() string {
	return "name"
}

func (i *InterfaceVeth) FindFieldValue() string {
	return i.Name
}

func (i *InterfaceVeth) DeleteField() string {
	return "numbers"
}

func (i *InterfaceVeth) DeleteFieldValue() string {
	return i.Name
}

func (client Mikrotik) AddInterfaceVeth(i *InterfaceVeth) (*InterfaceVeth, error) {
	res, err := client.Add(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVeth), nil
}

func (client Mikrotik) FindInterfaceVeth(name string) (*InterfaceVeth, error) {
	res, err := client.Find(&InterfaceVeth{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVeth), nil
}

func (client Mikrotik) UpdateInterfaceVeth(i *InterfaceVeth) (*InterfaceVeth, error) {
	res, err := client.Update(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVeth), nil
}

func (client Mikrotik) DeleteInterfaceVeth(name string) error {
	return client.Delete(&InterfaceVeth{Name: name})
}
