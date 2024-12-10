package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
)

const (
	WirelessAuthenticationTypeWpaPsk  = "wpa-psk"
	WirelessAuthenticationTypeWpa2Psk = "wpa2-psk"
	WirelessAuthenticationTypeWpaEap  = "wpa-eap"
	WirelessAuthenticationTypeWpa2Eap = "wpa2-eap"

	WirelessModeNone               = "none"
	WirelessModeStaticKeysOptional = "static-keys-optional"
	WirelessModeStaticKeysRequired = "static-keys-required"
	WirelessModeDynamicKeys        = "dynamic-keys"
)

// WirelessSecurityProfile defines resource
type WirelessSecurityProfile struct {
	Id                  string             `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                string             `mikrotik:"name" codegen:"name,required"`
	Mode                string             `mikrotik:"mode" codegen:"mode,optional"`
	AuthenticationTypes types.MikrotikList `mikrotik:"authentication-types" codegen:"authentication_types,optional"`
	WPA2PreSharedKey    string             `mikrotik:"wpa2-pre-shared-key" codegen:"wpa2_pre_shared_key"`
}

var _ Resource = (*WirelessSecurityProfile)(nil)

func (b *WirelessSecurityProfile) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/wireless/security-profiles/add",
		Find:   "/interface/wireless/security-profiles/print",
		Update: "/interface/wireless/security-profiles/set",
		Delete: "/interface/wireless/security-profiles/remove",
	}[a]
}

func (b *WirelessSecurityProfile) IDField() string {
	return ".id"
}

func (b *WirelessSecurityProfile) ID() string {
	return b.Id
}

func (b *WirelessSecurityProfile) SetID(id string) {
	b.Id = id
}

func (b *WirelessSecurityProfile) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddWirelessSecurityProfile(r *WirelessSecurityProfile) (*WirelessSecurityProfile, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*WirelessSecurityProfile), nil
}

func (c Mikrotik) UpdateWirelessSecurityProfile(r *WirelessSecurityProfile) (*WirelessSecurityProfile, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*WirelessSecurityProfile), nil
}

func (c Mikrotik) FindWirelessSecurityProfile(id string) (*WirelessSecurityProfile, error) {
	res, err := c.Find(&WirelessSecurityProfile{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*WirelessSecurityProfile), nil
}

func (c Mikrotik) ListWirelessSecurityProfile() ([]WirelessSecurityProfile, error) {
	res, err := c.List(&WirelessSecurityProfile{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]WirelessSecurityProfile, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*WirelessSecurityProfile))
	}

	return returnSlice, nil
}

func (c Mikrotik) DeleteWirelessSecurityProfile(id string) error {
	return c.Delete(&WirelessSecurityProfile{Id: id})
}
