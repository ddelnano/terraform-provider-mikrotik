package client

import (
	"github.com/go-routeros/routeros"
)

type InterfaceWireguardPeer struct {
	Id                     string `mikrotik:".id"`
	AllowedAddress         string `mikrotik:"allowed-address"`
	Comment                string `mikrotik:"comment"`
	Disabled               bool   `mikrotik:"disabled"`
	EndpointAddress        string `mikrotik:"endpoint-address"`
	EndpointPort           int    `mikrotik:"endpoint-port"`
	Interface              string `mikrotik:"interface"`
	PersistentKeepalive    int    `mikrotik:"persistent-keepalive"`
	PresharedKey           string `mikrotik:"preshared-key"`
	PublicKey              string `mikrotik:"public-key"`
	CurrentEndpointAddress string `mikrotik:"current-endpoint-address,readonly"`
	CurrentEndpointPort    int    `mikrotik:"current-endpoint-port,readonly"`
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

func (i *InterfaceWireguardPeer) FindField() string {
	return "interface"
}

func (i *InterfaceWireguardPeer) FindFieldValue() string {
	return i.Interface
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

func (client Mikrotik) FindInterfaceWireguardPeer(interfaceName string) (*InterfaceWireguardPeer, error) {
	res, err := client.Find(&InterfaceWireguardPeer{Interface: interfaceName})
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
