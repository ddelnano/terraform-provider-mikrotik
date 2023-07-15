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
	Id                       string `mikrotik:".id"`
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

var _ Resource = (*BgpInstance)(nil)

func (b *BgpInstance) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/routing/bgp/instance/add",
		Find:   "/routing/bgp/instance/print",
		Update: "/routing/bgp/instance/set",
		Delete: "/routing/bgp/instance/remove",
	}[a]
}

func (b *BgpInstance) IDField() string {
	return ".id"
}

func (b *BgpInstance) ID() string {
	return b.Id
}

func (b *BgpInstance) SetID(id string) {
	b.Id = id
}

func (b *BgpInstance) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BgpInstance) FindField() string {
	return "name"
}

func (b *BgpInstance) FindFieldValue() string {
	return b.Name
}

func (b *BgpInstance) DeleteField() string {
	return "numbers"
}

func (b *BgpInstance) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBgpInstance(r *BgpInstance) (*BgpInstance, error) {
	res, err := c.Add(r)
	if legacyBgpUnsupported(err) {
		return nil, LegacyBgpUnsupported{}
	}
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) UpdateBgpInstance(r *BgpInstance) (*BgpInstance, error) {
	res, err := c.Update(r)
	if legacyBgpUnsupported(err) {
		return nil, LegacyBgpUnsupported{}
	}
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) FindBgpInstance(name string) (*BgpInstance, error) {
	res, err := c.Find(&BgpInstance{Name: name})
	if legacyBgpUnsupported(err) {
		return nil, LegacyBgpUnsupported{}
	}

	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) DeleteBgpInstance(name string) error {
	err := c.Delete(&BgpInstance{Name: name})
	if legacyBgpUnsupported(err) {
		return LegacyBgpUnsupported{}
	}

	return err
}
