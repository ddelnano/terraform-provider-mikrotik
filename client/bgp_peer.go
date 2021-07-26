package client

import (
	"fmt"
	"log"
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

// AddBgpPeer Mikrotik resource
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
		return nil, err
	}

	return client.FindBgpPeer(b.Name)
}

// FindBgpPeer Mikrotik resource
func (client Mikrotik) FindBgpPeer(name string) (*BgpPeer, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/bgp/peer/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
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

// UpdateBgpPeer Mikrotik resource
func (client Mikrotik) UpdateBgpPeer(b *BgpPeer) (*BgpPeer, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	// compose mikrotik command
	cmd := Marshal("/routing/bgp/peer/set", b)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindBgpPeer(b.Name)
}

// DeleteBgpPeer Mikrotik resource
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
