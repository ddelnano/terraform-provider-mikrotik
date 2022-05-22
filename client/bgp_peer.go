package client

import (
	"fmt"
	"log"
)

// BgpPeer Mikrotik resource
//go:generate gen
type BgpPeer struct {
	ID                   string `mikrotik:".id" gen:"-,mikrotikID"`
	Name                 string `mikrotik:"name" gen:"name,required,id"`
	AddressFamilies      string `mikrotik:"address-families" gen:"address_families,optional,default=ip"`
	AllowAsIn            int    `mikrotik:"allow-as-in" gen:"allow_as_in,optional"`
	AsOverride           bool   `mikrotik:"as-override" gen:"as_override,optional"`
	CiscoVplsNlriLenFmt  string `mikrotik:"cisco-vpls-nlri-len-fmt" gen:"cisco_vpls_nlri_len_fmt,optional"`
	Comment              string `mikrotik:"comment" gen:"comment,optional"`
	DefaultOriginate     string `mikrotik:"default-originate" gen:"default_originate,optional,default=never"`
	Disabled             bool   `mikrotik:"disabled" gen:"disabled,optional"`
	HoldTime             string `mikrotik:"hold-time" gen:"hold_time,optional,default=3m"`
	InFilter             string `mikrotik:"in-filter" gen:"in_filter,optional"`
	Instance             string `mikrotik:"instance" gen:"instance,required"`
	KeepAliveTime        string `mikrotik:"keepalive-time" gen:"keepalive_time,optional"`
	MaxPrefixLimit       int    `mikrotik:"max-prefix-limit" gen:"max_prefix_limit,optional"`
	MaxPrefixRestartTime string `mikrotik:"max-prefix-restart-time" gen:"max_prefix_restart_time,optional"`
	Multihop             bool   `mikrotik:"multihop" gen:"multihop,optional"`
	NexthopChoice        string `mikrotik:"nexthop-choice" gen:"nexthop_choice,optional,default=default"`
	OutFilter            string `mikrotik:"out-filter" gen:"out_filter,optional"`
	Passive              bool   `mikrotik:"passive" gen:"passive,optional"`
	RemoteAddress        string `mikrotik:"remote-address" gen:"remote_address,required"`
	RemoteAs             int    `mikrotik:"remote-as" gen:"remote_as,required"`
	RemotePort           int    `mikrotik:"remote-port" gen:"remote_port,optional"`
	RemovePrivateAs      bool   `mikrotik:"remove-private-as" gen:"remove_private_as,optional"`
	RouteReflect         bool   `mikrotik:"route-reflect" gen:"route_reflect,optional"`
	TCPMd5Key            string `mikrotik:"tcp-md5-key" gen:"tcp_md5_key,optional"`
	TTL                  string `mikrotik:"ttl" gen:"ttl,optional,default=default"`
	UpdateSource         string `mikrotik:"update-source" gen:"update_source,optional"`
	UseBfd               bool   `mikrotik:"use-bfd" gen:"use_bfd,optional"`
}

func (client Mikrotik) AddBgpPeer(b *BgpPeer) (*BgpPeer, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/bgp/peer/add", b)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /routing/bgp/peer/add returned %v", r)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return client.FindBgpPeer(b.Name)
}

func (client Mikrotik) FindBgpPeer(name string) (*BgpPeer, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/bgp/peer/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	log.Printf("[DEBUG] Find bgp peer: `%v`", cmd)

	bgpPeer := BgpPeer{}

	err = Unmarshal(*r, &bgpPeer)

	if err != nil {
		return nil, err
	}

	if bgpPeer.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("bgp peer `%s` not found", name))
	}

	return &bgpPeer, nil
}

func (client Mikrotik) UpdateBgpPeer(b *BgpPeer) (*BgpPeer, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/bgp/peer/set", b)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return client.FindBgpPeer(b.Name)
}

func (client Mikrotik) DeleteBgpPeer(name string) error {
	c, err := client.getMikrotikClient()

	bgpPeer, err := client.FindBgpPeer(name)

	if err != nil {
		return err
	}

	cmd := []string{"/routing/bgp/peer/remove", "=numbers=" + bgpPeer.Name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] Remove bgp peer via mikrotik api: %v", r)

	return err
}
