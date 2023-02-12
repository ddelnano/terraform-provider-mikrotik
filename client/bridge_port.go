package client

import (
	"fmt"
	"log"
)

// BridgePort defines port-in-bridge association
type BridgePort struct {
	Id         string `mikrotik:".id"`
	Bridge     string `mikrotik:"bridge"`
	Interface  string `mikrotik:"interface"`
	PVId       int    `mikrotik:"pvid"`
	Comment    string `mikrotik:"comment"`
	FrameTypes string `mikrotik:"frame-types"`
}

func (client Mikrotik) AddBridgePort(r *BridgePort) (*BridgePort, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/port/add", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	response, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /interface/bridge/port/add returned %v", response)
	if err != nil {
		return nil, err
	}
	r.Id = response.Done.Map["ret"]

	return client.FindBridgePort(r.Id)
}

func (client Mikrotik) FindBridgePort(id string) (*BridgePort, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := []string{"/interface/bridge/port/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found bridge port: %v", r)

	record := BridgePort{}
	err = Unmarshal(*r, &record)
	if err != nil {
		return nil, err
	}
	if record.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("bridge port `%s` not found", id))
	}

	return &record, nil
}

func (client Mikrotik) UpdateBridgePort(r *BridgePort) (*BridgePort, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := Marshal("/interface/bridge/port/set", r)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindBridgePort(r.Id)
}

func (client Mikrotik) DeleteBridgePort(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}
	cmd := []string{"/interface/bridge/port/remove", "=numbers=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	return err
}
