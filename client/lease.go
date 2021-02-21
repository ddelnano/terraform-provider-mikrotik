package client

import (
	"fmt"
	"log"
	"strings"
)

type DhcpLease struct {
	Id          string `mikrotik:".id"`
	Address     string
	MacAddress  string `mikrotik:"mac-address"`
	Comment     string
	BlockAccess string `mikrotik:"blocked"`
	Hostname    string `mikrotik:"host-name"`
	Dynamic     bool
}

func (client Mikrotik) AddDhcpLease(address, macaddress, name string, blocked string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{
		"/ip/dhcp-server/lease/add",
		fmt.Sprintf("=address=%s", address),
		fmt.Sprintf("=mac-address=%s", macaddress),
		fmt.Sprintf("=comment=%s", name),
		fmt.Sprintf("=block-access=%s", blocked),
	}
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
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/print ?.id=%s", id), " ")
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

func (client Mikrotik) UpdateDhcpLease(id, address, macaddress, comment string, blocked string, dynamic bool) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := []string{
		"/ip/dhcp-server/lease/set",
		fmt.Sprintf("=.id=%s", id),
		fmt.Sprintf("=address=%s", address),
		fmt.Sprintf("=mac-address=%s", macaddress),
		fmt.Sprintf("=comment=%s", comment),
		fmt.Sprintf("=block-access=%s", blocked),
	}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindDhcpLease(id)
}

func (client Mikrotik) DeleteDhcpLease(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/remove =.id=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
