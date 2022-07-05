package client

import (
	"strings"

	"github.com/go-routeros/routeros"
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

var bgpInstanceWrapper *resourceWrapper = &resourceWrapper{
	idField:       "name",
	idFieldDelete: "numbers",
	actionsMap: map[string]string{
		"add":    "/routing/bgp/instance/add",
		"find":   "/routing/bgp/instance/print",
		"update": "/routing/bgp/instance/set",
		"delete": "/routing/bgp/instance/remove",
	},
	targetStruct:          &BgpInstance{},
	addIDExtractorFunc:    func(_ *routeros.Reply, resource interface{}) string { return resource.(*BgpInstance).Name },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*BgpInstance).Name },
}

// AddBgpInstance Mikrotik resource
func (client Mikrotik) AddBgpInstance(b *BgpInstance) (*BgpInstance, error) {
	r, err := bgpInstanceWrapper.Add(b, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpInstance), nil
}

// FindBgpInstance Mikrotik resource
func (client Mikrotik) FindBgpInstance(name string) (*BgpInstance, error) {
	r, err := bgpInstanceWrapper.Find(name, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpInstance), nil
}

// UpdateBgpInstance Mikrotik resource
func (client Mikrotik) UpdateBgpInstance(b *BgpInstance) (*BgpInstance, error) {
	r, err := bgpInstanceWrapper.Update(b, client.getMikrotikClient)
	if err != nil {
		if legacyBgpUnsupported(err) {
			return nil, LegacyBgpUnsupported{}
		}
		return nil, err
	}

	return r.(*BgpInstance), nil
}

// DeleteBgpInstance Mikrotik resource
func (client Mikrotik) DeleteBgpInstance(name string) error {
	return bgpInstanceWrapper.Delete(name, client.getMikrotikClient)
}
