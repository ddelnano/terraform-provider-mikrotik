package client

import (
	"github.com/go-routeros/routeros"
)

type InterfaceWireguardPeer struct {
	Id                  string `mikrotik:".id"`
	AllowedAddress      string `mikrotik:"allowed-address"`
	Comment             string `mikrotik:"comment"`
	Disabled            bool   `mikrotik:"disabled"`
	EndpointAddress     string `mikrotik:"endpoint-address"`
	EndpointPort        int64  `mikrotik:"endpoint-port"`
	Interface           string `mikrotik:"interface"`
	PersistentKeepalive int64  `mikrotik:"persistent-keepalive"`
	PresharedKey        string `mikrotik:"preshared-key"`
	PublicKey           string `mikrotik:"public-key"`
}

func (i *InterfaceWireguardPeer) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wireguard/peers/add",
		Find:   "/interface/wireguard/peers/print",
		List:   "/interface/wireguard/peers/print",
		Update: "/interface/wireguard/peers/set",
		Delete: "/interface/wireguard/peers/remove",
	}[action]
}

func (i *InterfaceWireguardPeer) IDField() string {
	return ".id"
}

func (i *InterfaceWireguardPeer) ID() string {
	return i.Id
}

func (i *InterfaceWireguardPeer) SetID(id string) {
	i.Id = id
}

func (i *InterfaceWireguardPeer) AfterAddHook(r *routeros.Reply) {
	i.Id = r.Done.Map["ret"]
}

func (i *InterfaceWireguardPeer) DeleteField() string {
	return "numbers"
}

func (i *InterfaceWireguardPeer) DeleteFieldValue() string {
	return i.Id
}

func (client Mikrotik) AddInterfaceWireguardPeer(i *InterfaceWireguardPeer) (*InterfaceWireguardPeer, error) {
	res, err := client.Add(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguardPeer), nil
}

func (client Mikrotik) FindInterfaceWireguardPeer(id string) (*InterfaceWireguardPeer, error) {
	res, err := client.Find(&InterfaceWireguardPeer{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguardPeer), nil
}

func (client Mikrotik) UpdateInterfaceWireguardPeer(i *InterfaceWireguardPeer) (*InterfaceWireguardPeer, error) {
	res, err := client.Update(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceWireguardPeer), nil
}

func (client Mikrotik) DeleteInterfaceWireguardPeer(id string) error {
	return client.Delete(&InterfaceWireguardPeer{Id: id})
}
