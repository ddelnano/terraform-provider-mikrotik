package client

import (
	"github.com/go-routeros/routeros"
)

type InterfaceWireguard struct {
	Id         string `mikrotik:".id"`
	Name       string `mikrotik:"name"`
	Comment    string `mikrotik:"comment"`
	Disabled   bool   `mikrotik:"disabled"`
	ListenPort int    `mikrotik:"listen-port"`
	Mtu        int    `mikrotik:"mtu"`
	PrivateKey string `mikrotik:"private-key"`
	PublicKey  string `mikrotik:"public-key"`       //read only property
	Running    bool   `mikrotik:"running,readonly"` //read only property
}

func (i *InterfaceWireguard) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wireguard/add",
		Find:   "/interface/wireguard/print",
		List:   "/interface/wireguard/print",
		Update: "/interface/wireguard/set",
		Delete: "/interface/wireguard/remove",
	}[action]
}

func (i *InterfaceWireguard) IDField() string {
	return ".id"
}

func (i *InterfaceWireguard) ID() string {
	return i.Id
}

func (i *InterfaceWireguard) SetID(id string) {
	i.Id = id
}

func (i *InterfaceWireguard) AfterAddHook(r *routeros.Reply) {
	i.Id = r.Done.Map["ret"]
}

func (i *InterfaceWireguard) FindField() string {
	return ".id"
}

func (i *InterfaceWireguard) FindFieldValue() string {
	return i.Id
}

func (i *InterfaceWireguard) DeleteField() string {
	return ".id"
}

func (i *InterfaceWireguard) DeleteFieldValue() string {
	return i.Id
}

func (client Mikrotik) AddInterfaceWireguard(i *InterfaceWireguard) (*InterfaceWireguard, error) {
	res, err := client.Add(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguard), nil
}

func (client Mikrotik) FindInterfaceWireguard(id string) (*InterfaceWireguard, error) {
	res, err := client.Find(&InterfaceWireguard{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguard), nil
}

func (client Mikrotik) UpdateInterfaceWireguard(i *InterfaceWireguard) (*InterfaceWireguard, error) {
	res, err := client.Update(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguard), nil
}

func (client Mikrotik) DeleteInterfaceWireguard(id string) error {
	return client.Delete(&InterfaceWireguard{Id: id})
}
