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
type BgpInstance struct {
	ID                       string `mikrotik:".id"`
	Name                     string `mikrotik:"name"`
	As                       int    `mikrotik:"as"`
	ClientToClientReflection bool   `mikrotik:"client-to-client-reflection"`
	Comment                  string `mikrotik:"comment"`
	ConfederationPeers       string `mikrotik:"confederation-peers"`
	Disabled                 bool   `mikrotik:"disabled"`
	IgnoreAsPathLen          bool   `mikrotik:"ignore-as-path-len"`
	OutFilter                string `mikrotik:"out-filter"`
	RedistributeConnected    bool   `mikrotik:"redistribute-connected"`
	RedistributeOspf         bool   `mikrotik:"redistribute-ospf"`
	RedistributeOtherBgp     bool   `mikrotik:"redistribute-other-bgp"`
	RedistributeRip          bool   `mikrotik:"redistribute-rip"`
	RedistributeStatic       bool   `mikrotik:"redistribute-static"`
	RouterID                 string `mikrotik:"router-id"`
	RoutingTable             string `mikrotik:"routing-table"`
	ClusterID                string `mikrotik:"cluster-id"`
	Confederation            int    `mikrotik:"confederation"`
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
	if err != nil {
		return err
	}

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
