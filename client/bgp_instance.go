package client

import (
	"fmt"
	"log"
	"strings"
)

type LegacyBgpUnsupported struct{}

func (LegacyBgpUnsupported) Error() string {
	return "Your RouterOS version does not support /routing/bgp/{instance,peer} commands"
}

func legacyBgpUnsupported(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "no such command prefix") {
			return true
		}
	}
	return false
}

// BgpInstance Mikrotik resource
//go:generate gen
type BgpInstance struct {
	ID                       string `mikrotik:".id" gen:"-,mikrotikID"`
	Name                     string `mikrotik:"name" gen:"name,id,required"`
	As                       int    `mikrotik:"as" gen:"as,required"`
	ClientToClientReflection bool   `mikrotik:"client-to-client-reflection" gen:"client_to_client_reflection,optional,default=true"`
	Comment                  string `mikrotik:"comment" gen:"comment,optional"`
	ConfederationPeers       string `mikrotik:"confederation-peers" gen:"confederation_peers,optional"`
	Disabled                 bool   `mikrotik:"disabled" gen:"disabled,optional"`
	IgnoreAsPathLen          bool   `mikrotik:"ignore-as-path-len" gen:"ignore_as_path_len,optional"`
	OutFilter                string `mikrotik:"out-filter" gen:"out_filter,optional"`
	RedistributeConnected    bool   `mikrotik:"redistribute-connected" gen:"redistribute_connected,optional"`
	RedistributeOspf         bool   `mikrotik:"redistribute-ospf" gen:"redistribute_ospf,optional"`
	RedistributeOtherBgp     bool   `mikrotik:"redistribute-other-bgp" gen:"redistribute_other_bgp,optional"`
	RedistributeRip          bool   `mikrotik:"redistribute-rip" gen:"redistribute_rip,optional"`
	RedistributeStatic       bool   `mikrotik:"redistribute-static" gen:"redistribute_static,optional"`
	RouterID                 string `mikrotik:"router-id" gen:"router_id,required"`
	RoutingTable             string `mikrotik:"routing-table" gen:"routing_table,optional"`
	ClusterID                string `mikrotik:"cluster-id" gen:"cluster_id,optional"`
	Confederation            int    `mikrotik:"confederation" gen:"confederation,optional"`
}

// AddBgpInstance Mikrotik resource
func (client Mikrotik) AddBgpInstance(b *BgpInstance) (*BgpInstance, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := Marshal("/routing/bgp/instance/add", b)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /routing/bgp/instance/add returned %v", r)

	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return client.FindBgpInstance(b.Name)
}

// FindBgpInstance Mikrotik resource
func (client Mikrotik) FindBgpInstance(name string) (*BgpInstance, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/bgp/instance/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	log.Printf("[DEBUG] Find bgp instance: `%v`", cmd)

	bgpInstance := BgpInstance{}

	err = Unmarshal(*r, &bgpInstance)

	if err != nil {
		return nil, err
	}

	if bgpInstance.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("bgp instance `%s` not found", name))
	}

	return &bgpInstance, nil
}

// UpdateBgpInstance Mikrotik resource
func (client Mikrotik) UpdateBgpInstance(b *BgpInstance) (*BgpInstance, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	// compose mikrotik command
	cmd := Marshal("/routing/bgp/instance/set", b)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return client.FindBgpInstance(b.Name)
}

// DeleteBgpInstance Mikrotik resource
func (client Mikrotik) DeleteBgpInstance(name string) error {
	c, err := client.getMikrotikClient()

	bgpInstance, err := client.FindBgpInstance(name)

	if err != nil {
		return err
	}

	cmd := []string{"/routing/bgp/instance/remove", "=numbers=" + bgpInstance.Name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] Remove bgp instance via mikrotik api: %v", r)

	return err
}
