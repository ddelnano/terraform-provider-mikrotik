package client

import (
	"fmt"
	"log"
)

// DhcpServer represents DHCP server resource
type DhcpServer struct {
	Id            string `mikrotik:".id"`
	Name          string `mikrotik:"name"`
	Disabled      bool   `mikrotik:"disabled"`
	AddArp        bool   `mikrotik:"add-arp"`
	AddressPool   string `mikrotik:"address-pool"`
	Authoritative string `mikrotik:"authoritative"`
	Interface     string `mikrotik:"interface"`
	LeaseScript   string `mikrotik:"lease-script"`
}

func (client Mikrotik) AddDhcpServer(d *DhcpServer) (*DhcpServer, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindDhcpServer(d.Name)
}

func (client Mikrotik) UpdateDhcpServer(d *DhcpServer) (*DhcpServer, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/dhcp-server/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindDhcpServer(d.Name)
}

func (client Mikrotik) FindDhcpServer(name string) (*DhcpServer, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/dhcp-server/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Found record: %v", r)
	record := DhcpServer{}
	err = Unmarshal(*r, &record)

	if err != nil {
		return nil, err
	}

	if record.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("record `%s` not found", name))
	}

	return &record, nil
}

func (client Mikrotik) DeleteDhcpServer(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/ip/dhcp-server/remove", "=numbers=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Command returned: %v", r)

	return nil
}
