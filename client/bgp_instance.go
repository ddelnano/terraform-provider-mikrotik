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
	Id                       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                     string `mikrotik:"name" codegen:"name,required,terraformID"`
	As                       int    `mikrotik:"as" codegen:"as,required"`
	ClientToClientReflection bool   `mikrotik:"client-to-client-reflection" codegen:"client_to_client_reflection"`
	Comment                  string `mikrotik:"comment" codegen:"comment"`
	ConfederationPeers       string `mikrotik:"confederation-peers" codegen:"confederation_peers"`
	Disabled                 bool   `mikrotik:"disabled" codegen:"disabled"`
	IgnoreAsPathLen          bool   `mikrotik:"ignore-as-path-len" codegen:"ignore_as_path_len"`
	OutFilter                string `mikrotik:"out-filter" codegen:"out_filter"`
	RedistributeConnected    bool   `mikrotik:"redistribute-connected" codegen:"redistribute_connected"`
	RedistributeOspf         bool   `mikrotik:"redistribute-ospf" codegen:"redistribute_ospf"`
	RedistributeOtherBgp     bool   `mikrotik:"redistribute-other-bgp" codegen:"redistribute_other_bgp"`
	RedistributeRip          bool   `mikrotik:"redistribute-rip" codegen:"redistribute_rip"`
	RedistributeStatic       bool   `mikrotik:"redistribute-static" codegen:"redistribute_static"`
	RouterID                 string `mikrotik:"router-id" codegen:"router_id,required"`
	RoutingTable             string `mikrotik:"routing-table" codegen:"routing_table"`
	ClusterID                string `mikrotik:"cluster-id" codegen:"cluster_id"`
	Confederation            int    `mikrotik:"confederation" codegen:"confederation"`
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

// HandleError intercepts errors during CRUD operations.
// It is used to catch "no such command prefix" on RouterOS >= v7.0
func (b *BgpInstance) HandleError(err error) error {
	if legacyBgpUnsupported(err) {
		return LegacyBgpUnsupported{}
	}

	return err
}

// Typed wrappers
func (c Mikrotik) AddBgpInstance(r *BgpInstance) (*BgpInstance, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) UpdateBgpInstance(r *BgpInstance) (*BgpInstance, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) FindBgpInstance(name string) (*BgpInstance, error) {
	res, err := c.Find(&BgpInstance{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstance), nil
}

func (c Mikrotik) DeleteBgpInstance(name string) error {
	err := c.Delete(&BgpInstance{Name: name})
	return err
}
