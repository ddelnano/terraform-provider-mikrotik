package client

import (
	"fmt"
	"log"
)

// InterfaceListMember manages an interface list's members
type InterfaceListMember struct {
	Id        string `mikrotik:".id"`
	Interface string `mikrotik:"interface"`
	List      string `mikrotik:"list"`
}

func (client Mikrotik) AddInterfaceListMember(d *InterfaceListMember) (*InterfaceListMember, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/list/member/add", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)
	id := r.Done.Map["ret"]

	return client.FindInterfaceListMember(id)
}

func (client Mikrotik) FindInterfaceListMember(id string) (*InterfaceListMember, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/list/member/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found record: %v", r)

	record := InterfaceListMember{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}

	if record.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("interface list member `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateInterfaceListMember(d *InterfaceListMember) (*InterfaceListMember, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/list/member/set", d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] command returned: %v", r)

	return client.FindInterfaceListMember(d.Id)
}

func (client Mikrotik) DeleteInterfaceListMember(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/interface/list/member/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Command returned: %v", r)

	return nil
}
