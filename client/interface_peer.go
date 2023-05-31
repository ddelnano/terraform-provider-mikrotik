package client

import (
	"github.com/go-routeros/routeros"
)

type InterfacePeer struct {
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
	CurrentEndpointPort    string `mikrotik:"current-endpoint-port,readonly"`
	LastHandshake          string `mikrotik:"last-handshake,readonly"`
	Rx                     string `mikrotik:"rx,readonly"`
	Tx                     string `mikrotik:"tx,readonly"`
}

func (i *InterfacePeer) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wireguard/peers/add", //is this correct?
		Find:   "/interface/wireguard/peers/print",
		List:   "/interface/wireguard/peers/print",
		Update: "/interface/wireguard/peers/set",
		Delete: "/interface/wireguard/peers/remove",
	}[action]
}

func (i *InterfacePeer) IDField() string {
	return ".id"
}

func (i *InterfacePeer) ID() string {
	return i.Id
}

func (i *InterfacePeer) SetID(id string) {
	i.Id = id
}

func (i *InterfacePeer) AfterAddHook(r *routeros.Reply) {
	i.Id = r.Done.Map["ret"]
}

func (i *InterfacePeer) FindField() string {
	return ".id"
}

func (i *InterfacePeer) FindFieldValue() string {
	return i.Id
}

func (i *InterfacePeer) DeleteField() string {
	return "numbers"
}

func (i *InterfacePeer) DeleteFieldValue() string {
	return i.Id
}

func (client Mikrotik) AddInterfacePeer(i *InterfacePeer) (*InterfacePeer, error) {
	res, err := client.Add(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfacePeer), nil
}

func (client Mikrotik) FindInterfacePeer(id string) (*InterfacePeer, error) {
	res, err := client.Find(&InterfacePeer{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*InterfacePeer), nil
}

func (client Mikrotik) UpdateInterfacePeer(i *InterfacePeer) (*InterfacePeer, error) {
	res, err := client.Update(i)
	if err != nil {
		return nil, err
	}

	return res.(*InterfacePeer), nil
}

func (client Mikrotik) DeleteInterfacePeer(id string) error {
	return client.Delete(&InterfacePeer{Id: id})
}
