package client

import "github.com/go-routeros/routeros"

const (
	WirelessInterfaceModeStation                   = "station"
	WirelessInterfaceModeStationWDS                = "station-wds"
	WirelessInterfaceModeAPBridge                  = "ap-bridge"
	WirelessInterfaceModeBridge                    = "bridge"
	WirelessInterfaceModeAlignmentOnly             = "alignment-only"
	WirelessInterfaceModeNstremeDualSlave          = "nstreme-dual-slave"
	WirelessInterfaceModeWDSSlave                  = "wds-slave"
	WirelessInterfaceModeStationPseudobridge       = "station-pseudobridge"
	WirelessInterfaceModeStationsPseudobridgeClone = "station-pseudobridge-clone"
	WirelessInterfaceModeStationBridge             = "station-bridge"
)

// WirelessInterface defines resource
type WirelessInterface struct {
	Id              string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name            string `mikrotik:"name" codegen:"name,required"`
	MasterInterface string `mikrotik:"master-insterface" codegen:"master_interface"`
	Mode            string `mikrotik:"mode" codegen:"mode"`
	Disabled        bool   `mikrotik:"disabled" codegen:"disabled"`
	SecurityProfile string `mikrotik:"security-profile" codegen:"security_profile"`
	SSID            string `mikrotik:"ssid" codegen:"ssid"`
	HideSSID        bool   `mikrotik:"hide-ssid" codegen:"hide_ssid"`
	VlanID          int    `mikrotik:"vlan-id" codegen:"vlan_id"`
	VlanMode        string `mikrotik:"vlan-mode" codegen:"vlan_mode"`
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
