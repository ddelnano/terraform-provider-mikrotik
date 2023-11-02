package client

import "github.com/go-routeros/routeros"

// WirelessInterface defines resource
type WirelessInterface struct {
	Id              string `mikrotik:".id"`
	Name            string `mikrotik:"name"`
	MasterInterface string `mikrotik:"master-insterface"`
	Mode            string `mikrotik:"mode"`
	Disabled        bool   `mikrotik:"disabled"`
	SecurityProfile string `mikrotik:"security-profile"`
	SSID            string `mikrotik:"ssid"`
	HideSSID        bool   `mikrotik:"hide-ssid"`
	VlanID          string `mikrotik:"vlan-id"`
	VlanMode        string `mikrotik:"vlan-mode"`
}

var _ Resource = (*WirelessInterface)(nil)

func (b *WirelessInterface) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/wireless/add",
		Find:   "/interface/wireless/print",
		Update: "/interface/wireless/set",
		Delete: "/interface/wireless/remove",
	}[a]
}

func (b *WirelessInterface) IDField() string {
	return ".id"
}

func (b *WirelessInterface) ID() string {
	return b.Id
}

func (b *WirelessInterface) SetID(id string) {
	b.Id = id
}

func (b *WirelessInterface) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddWirelessInterface(r *WirelessInterface) (*WirelessInterface, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*WirelessInterface), nil
}

func (c Mikrotik) UpdateWirelessInterface(r *WirelessInterface) (*WirelessInterface, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*WirelessInterface), nil
}

func (c Mikrotik) FindWirelessInterface(id string) (*WirelessInterface, error) {
	res, err := c.Find(&WirelessInterface{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*WirelessInterface), nil
}

func (c Mikrotik) ListWirelessInterface() ([]WirelessInterface, error) {
	res, err := c.List(&WirelessInterface{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]WirelessInterface, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*WirelessInterface))
	}

	return returnSlice, nil
}

func (c Mikrotik) DeleteWirelessInterface(id string) error {
	return c.Delete(&WirelessInterface{Id: id})
}
