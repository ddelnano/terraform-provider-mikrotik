package client

import (
	"github.com/go-routeros/routeros"
)

// BgpPeer Mikrotik resource
type BgpPeer struct {
	Id                   string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                 string `mikrotik:"name" codegen:"name,required,terraformID"`
	AddressFamilies      string `mikrotik:"address-families" codegen:"address_families,optional,computed"`
	AllowAsIn            int    `mikrotik:"allow-as-in" codegen:"allow_as_in"`
	AsOverride           bool   `mikrotik:"as-override" codegen:"as_override"`
	CiscoVplsNlriLenFmt  string `mikrotik:"cisco-vpls-nlri-len-fmt" codegen:"cisco_vpls_nlri_len_fmt"`
	Comment              string `mikrotik:"comment" codegen:"comment"`
	DefaultOriginate     string `mikrotik:"default-originate" codegen:"default_originate,optional,computed"`
	Disabled             bool   `mikrotik:"disabled" codegen:"disabled"`
	HoldTime             string `mikrotik:"hold-time" codegen:"hold_time,optional,computed"`
	InFilter             string `mikrotik:"in-filter" codegen:"in_filter"`
	Instance             string `mikrotik:"instance" codegen:"instance"`
	KeepAliveTime        string `mikrotik:"keepalive-time" codegen:"keepalive_time"`
	MaxPrefixLimit       int    `mikrotik:"max-prefix-limit" codegen:"max_prefix_limit"`
	MaxPrefixRestartTime string `mikrotik:"max-prefix-restart-time" codegen:"max_prefix_restart_time"`
	Multihop             bool   `mikrotik:"multihop" codegen:"multihop"`
	NexthopChoice        string `mikrotik:"nexthop-choice" codegen:"nexthop_choice,optional,computed"`
	OutFilter            string `mikrotik:"out-filter" codegen:"out_filter"`
	Passive              bool   `mikrotik:"passive" codegen:"passive"`
	RemoteAddress        string `mikrotik:"remote-address" codegen:"remote_address,required"`
	RemoteAs             int    `mikrotik:"remote-as" codegen:"remote_as,required"`
	RemotePort           int    `mikrotik:"remote-port" codegen:"remote_port"`
	RemovePrivateAs      bool   `mikrotik:"remove-private-as" codegen:"remove_private_as"`
	RouteReflect         bool   `mikrotik:"route-reflect" codegen:"route_reflect"`
	TCPMd5Key            string `mikrotik:"tcp-md5-key" codegen:"tcp_md5_key"`
	TTL                  string `mikrotik:"ttl" codegen:"ttl,optional,computed"`
	UpdateSource         string `mikrotik:"update-source" codegen:"update_source"`
	UseBfd               bool   `mikrotik:"use-bfd" codegen:"use_bfd"`
}

var _ Resource = (*BgpPeer)(nil)

func (b *BgpPeer) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/routing/bgp/peer/add",
		Find:   "/routing/bgp/peer/print",
		Update: "/routing/bgp/peer/set",
		Delete: "/routing/bgp/peer/remove",
	}[a]
}

func (b *BgpPeer) IDField() string {
	return ".id"
}

func (b *BgpPeer) ID() string {
	return b.Id
}

func (b *BgpPeer) SetID(id string) {
	b.Id = id
}

func (b *BgpPeer) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BgpPeer) FindField() string {
	return "name"
}

func (b *BgpPeer) FindFieldValue() string {
	return b.Name
}

func (b *BgpPeer) DeleteField() string {
	return "numbers"
}

func (b *BgpPeer) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBgpPeer(r *BgpPeer) (*BgpPeer, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpPeer), nil
}

func (c Mikrotik) UpdateBgpPeer(r *BgpPeer) (*BgpPeer, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpPeer), nil
}

func (c Mikrotik) FindBgpPeer(name string) (*BgpPeer, error) {
	res, err := c.Find(&BgpPeer{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpPeer), nil
}

func (c Mikrotik) DeleteBgpPeer(name string) error {
	return c.Delete(&BgpPeer{Name: name})
}
