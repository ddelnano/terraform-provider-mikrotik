package client

import (
	"fmt"
	"log"
)

// BridgeVlan defines vlan filtering in bridge resource
type BridgeVlan struct {
	Id       string   `mikrotik:".id"`
	Bridge   string   `mikrotik:"bridge"`
	Tagged   []string `mikrotik:"tagged"`
	Untagged []string `mikrotik:"untagged"`
	VlanIds  []string `mikrotik:"vlan-ids"`
}

func (client Mikrotik) AddBridgeVlan(r *BridgeVlan) (*BridgeVlan, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/vlan/add", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	response, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /interface/bridge/vlan/add returned %v", response)
	if err != nil {
		return nil, err
	}
	id := response.Done.Map["ret"]

	return client.FindBridgeVlan(id)
}

func (client Mikrotik) FindBridgeVlan(id string) (*BridgeVlan, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/bridge/vlan/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found bridge vlan: %v", r)

	record := BridgeVlan{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}
	if record.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("bridge vlan `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateBridgeVlan(r *BridgeVlan) (*BridgeVlan, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/vlan/set", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindBridgeVlan(r.Id)
}

func (client Mikrotik) DeleteBridgeVlan(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}
	cmd := []string{"/interface/bridge/vlan/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	return err
}
