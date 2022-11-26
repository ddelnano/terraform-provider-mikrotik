package client

import (
	"fmt"
	"log"
)

// Bridge defines /bridge resource
type Bridge struct {
	Id            string `mikrotik:".id"`
	Name          string `mikrotik:"name"`
	FastForward   bool   `mikrotik:"fast-forward"`
	VlanFiltering bool   `mikrotik:"vlan-filtering"`
	Comment       string `mikrotik:"comment"`
}

func (client Mikrotik) AddBridge(r *Bridge) (*Bridge, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/add", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	response, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /interface/bridge/add returned %v", response)
	if err != nil {
		return nil, err
	}

	return client.FindBridge(r.Name)
}

func (client Mikrotik) FindBridge(id string) (*Bridge, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/bridge/print", "?name=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found bridge: %v", r)

	record := Bridge{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}
	if record.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("bridge `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateBridge(r *Bridge) (*Bridge, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/set", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindBridge(r.Name)
}

func (client Mikrotik) DeleteBridge(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}
	cmd := []string{"/interface/bridge/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	return err
}
