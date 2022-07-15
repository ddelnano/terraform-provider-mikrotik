package client

import (
	"fmt"
	"log"
)

// DhcpServerNetwork describes network configuration for DHCP server
type DhcpServerNetwork struct {
	Id        string `mikrotik:".id"`
	Comment   string `mikrotik:"comment"`
	Address   string `mikrotik:"address"`
	Netmask   string `mikrotik:"netmask"`
	Gateway   string `mikrotik:"gateway"`
	DnsServer string `mikrotik:"dns-server"`
}

func (client Mikrotik) AddDhcpServerNetwork(d *DhcpServerNetwork) (*DhcpServerNetwork, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/network/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	id := r.Done.Map["ret"]
	return client.FindDhcpServerNetwork(id)
}

func (client Mikrotik) FindDhcpServerNetwork(id string) (*DhcpServerNetwork, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/dhcp-server/network/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found record: %v", r)

	record := DhcpServerNetwork{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}

	if record.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("record `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateDhcpServerNetwork(d *DhcpServerNetwork) (*DhcpServerNetwork, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/network/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindDhcpServerNetwork(d.Id)
}

func (client Mikrotik) DeleteDhcpServerNetwork(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/ip/dhcp-server/network/remove", "=.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Command returned: %v", r)

	return nil
}
