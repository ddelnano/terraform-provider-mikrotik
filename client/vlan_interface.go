package client

import (
	"fmt"
	"log"
)

// VlanInterface represents vlan interface resource
type VlanInterface struct {
	Id            string `mikrotik:".id"`
	Interface     string `mikrotik:"interface"`
	Mtu           int    `mikrotik:"mtu"`
	Name          string `mikrotik:"name"`
	Disabled      bool   `mikrotik:"disabled"`
	UseServiceTag bool   `mikrotik:"use-service-tag"`
	VlanId        int    `mikrotik:"vlan-id"`
}

func (client Mikrotik) AddVlanInterface(d *VlanInterface) (*VlanInterface, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/vlan/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindVlanInterface(d.Name)
}

func (client Mikrotik) UpdateVlanInterface(d *VlanInterface) (*VlanInterface, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/vlan/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindVlanInterface(d.Name)
}

func (client Mikrotik) FindVlanInterface(name string) (*VlanInterface, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/vlan/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Found record: %v", r)
	record := VlanInterface{}
	err = Unmarshal(*r, &record)

	if err != nil {
		return nil, err
	}

	if record.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("vlan interface `%s` not found", name))
	}

	return &record, nil
}

func (client Mikrotik) DeleteVlanInterface(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/interface/vlan/remove", "=numbers=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Command returned: %v", r)

	return nil
}
