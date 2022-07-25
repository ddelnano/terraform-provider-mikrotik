package client

import (
	"fmt"
	"log"
)

// InterfaceList manages a list of interfaces
type InterfaceList struct {
	Id      string `mikrotik:".id"`
	Comment string `mikrotik:"comment"`
	Name    string `mikrotik:"name"`
}

func (client Mikrotik) AddInterfaceList(d *InterfaceList) (*InterfaceList, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/list/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindInterfaceList(d.Name)
}

func (client Mikrotik) FindInterfaceList(id string) (*InterfaceList, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/list/print", "?name=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found record: %v", r)

	record := InterfaceList{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}

	if record.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("interface list `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateInterfaceList(d *InterfaceList) (*InterfaceList, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/list/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindInterfaceList(d.Name)
}

func (client Mikrotik) DeleteInterfaceList(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/interface/list/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Command returned: %v", r)

	return nil
}
