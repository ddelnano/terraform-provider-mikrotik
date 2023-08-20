package client

import (
	"log"

	"github.com/go-routeros/routeros"
)

type DhcpLease struct {
	Id          string `mikrotik:".id"`
	Address     string `mikrotik:"address"`
	MacAddress  string `mikrotik:"mac-address"`
	Comment     string `mikrotik:"comment"`
	BlockAccess bool   `mikrotik:"block-access"`
	Dynamic     bool   // TODO:  don't see this listed as a param https://wiki.mikrotik.com/wiki/Manual:IP/DHCP_Server, but our docs list it as one
	Hostname    string
}

func (client Mikrotik) ListDhcpLeases() ([]DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/dhcp-server/lease/print"}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found dhcp leases: %v", r)

	leases := []DhcpLease{}

	err = Unmarshal(*r, &leases)

	if err != nil {
		return nil, err
	}

	return leases, nil
}

var _ Resource = (*DhcpLease)(nil)

func (b *DhcpLease) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/dhcp-server/lease/add",
		Find:   "/ip/dhcp-server/lease/print",
		Update: "/ip/dhcp-server/lease/set",
		Delete: "/ip/dhcp-server/lease/remove",
	}[a]
}

func (b *DhcpLease) IDField() string {
	return ".id"
}

func (b *DhcpLease) ID() string {
	return b.Id
}

func (b *DhcpLease) SetID(id string) {
	b.Id = id
}

func (b *DhcpLease) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddDhcpLease(r *DhcpLease) (*DhcpLease, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpLease), nil
}

func (c Mikrotik) UpdateDhcpLease(r *DhcpLease) (*DhcpLease, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*DhcpLease), nil
}

func (c Mikrotik) FindDhcpLease(id string) (*DhcpLease, error) {
	res, err := c.Find(&DhcpLease{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*DhcpLease), nil
}

func (client Mikrotik) ListDhcpLease() ([]DhcpLease, error) {
	res, err := client.List(&DhcpLease{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]DhcpLease, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*DhcpLease))
	}

	return returnSlice, nil
}

func (c Mikrotik) DeleteDhcpLease(id string) error {
	return c.Delete(&DhcpLease{Id: id})
}
