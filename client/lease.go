package client

import (
	"fmt"
	"log"
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

func (client Mikrotik) AddDhcpLease(l *DhcpLease) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/lease/add", l)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Dhcp lease creation response: `%v`", r)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.FindDhcpLease(id)
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

func (client Mikrotik) FindDhcpLease(id string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/dhcp-server/lease/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Dhcp lease response: %v", r)

	if err != nil {
		return nil, err
	}

	lease := DhcpLease{}
	err = Unmarshal(*r, &lease)

	if err != nil {
		return nil, err
	}

	if lease.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("dhcp lease `%s` not found", id))
	}

	return &lease, nil
}

func (client Mikrotik) UpdateDhcpLease(l *DhcpLease) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/lease/set", l)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindDhcpLease(l.Id)
}

func (client Mikrotik) DeleteDhcpLease(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := []string{"/ip/dhcp-server/lease/remove", "=.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
