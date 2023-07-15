package client

import (
	"github.com/go-routeros/routeros"
)

// BgpPeer Mikrotik resource
type BgpPeer struct {
	Id                   string `mikrotik:".id"`
	Name                 string `mikrotik:"name"`
	AddressFamilies      string `mikrotik:"address-families"`
	AllowAsIn            int    `mikrotik:"allow-as-in"`
	AsOverride           bool   `mikrotik:"as-override"`
	CiscoVplsNlriLenFmt  string `mikrotik:"cisco-vpls-nlri-len-fmt"`
	Comment              string `mikrotik:"comment"`
	DefaultOriginate     string `mikrotik:"default-originate"`
	Disabled             bool   `mikrotik:"disabled"`
	HoldTime             string `mikrotik:"hold-time"`
	InFilter             string `mikrotik:"in-filter"`
	Instance             string `mikrotik:"instance"`
	KeepAliveTime        string `mikrotik:"keepalive-time"`
	MaxPrefixLimit       int    `mikrotik:"max-prefix-limit"`
	MaxPrefixRestartTime string `mikrotik:"max-prefix-restart-time"`
	Multihop             bool   `mikrotik:"multihop"`
	NexthopChoice        string `mikrotik:"nexthop-choice"`
	OutFilter            string `mikrotik:"out-filter"`
	Passive              bool   `mikrotik:"passive"`
	RemoteAddress        string `mikrotik:"remote-address"`
	RemoteAs             int    `mikrotik:"remote-as"`
	RemotePort           int    `mikrotik:"remote-port"`
	RemovePrivateAs      bool   `mikrotik:"remove-private-as"`
	RouteReflect         bool   `mikrotik:"route-reflect"`
	TCPMd5Key            string `mikrotik:"tcp-md5-key"`
	TTL                  string `mikrotik:"ttl"`
	UpdateSource         string `mikrotik:"update-source"`
	UseBfd               bool   `mikrotik:"use-bfd"`
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
