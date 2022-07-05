package client

import (
	"github.com/go-routeros/routeros"
)

// BgpPeer Mikrotik resource
type BgpPeer struct {
	ID                   string `mikrotik:".id"`
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

var bgpPeerWrapper *resourceWrapper = &resourceWrapper{
	idField:       "name",
	idFieldDelete: "numbers",
	actionsMap: map[string]string{
		"add":    "/routing/bgp/peer/add",
		"find":   "/routing/bgp/peer/print",
		"update": "/routing/bgp/peer/set",
		"delete": "/routing/bgp/peer/remove",
	},
	targetStruct:          &BgpPeer{},
	addIDExtractorFunc:    func(_ *routeros.Reply, resource interface{}) string { return resource.(*BgpPeer).Name },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*BgpPeer).Name },
}

func (client Mikrotik) AddBgpPeer(b *BgpPeer) (*BgpPeer, error) {
	r, err := bgpPeerWrapper.Add(b, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpPeer), nil
}

func (client Mikrotik) FindBgpPeer(name string) (*BgpPeer, error) {
	r, err := bgpPeerWrapper.Find(name, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpPeer), nil
}

func (client Mikrotik) UpdateBgpPeer(b *BgpPeer) (*BgpPeer, error) {
	r, err := bgpPeerWrapper.Update(b, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpPeer), nil
}

func (client Mikrotik) DeleteBgpPeer(name string) error {
	return bgpPeerWrapper.Delete(name, client.getMikrotikClient)
}
