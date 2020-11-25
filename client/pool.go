package client

import (
	"fmt"
	"log"
	"strings"
)

type Pool struct {
	Id      string `mikrotik:".id"`
	Name    string
	Ranges  string `mikrotik:"ranges"`
	Comment string
}

func (client Mikrotik) AddPool(name string, ranges string, comment string) (*Pool, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/pool/add =name=%s =ranges=%s =comment=%s", name, ranges, comment), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Pool creation response: `%v`", r)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.FindPool(id)
}

func (client Mikrotik) ListPools() ([]Pool, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ip/pool/print"}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found pools: %v", r)

	pools := []Pool{}

	err = Unmarshal(*r, &pools)

	if err != nil {
		return nil, err
	}

	return pools, nil
}

func (client Mikrotik) FindPool(id string) (*Pool, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/pool/print ?.id=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Pool response: %v", r)

	if err != nil {
		return nil, err
	}

	pool := Pool{}
	err = Unmarshal(*r, &pool)

	if err != nil {
		return nil, err
	}
	if pool.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("pool `%s` not found", id))
	}

	return &pool, nil
}

func (client Mikrotik) FindPoolByName(name string) (*Pool, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/pool/print ?name=%s", name), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Pool response: %v", r)

	if err != nil {
		return nil, err
	}

	pool := Pool{}
	err = Unmarshal(*r, &pool)

	if err != nil {
		return nil, err
	}

	if pool.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("pool `%s` not found", name))
	}

	return &pool, nil
}

func (client Mikrotik) UpdatePool(id, name string, ranges string, comment string) (*Pool, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := strings.Split(fmt.Sprintf("/ip/pool/set =.id=%s =name=%s =ranges=%s =comment=%s", id, name, ranges, comment), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindPool(id)
}

func (client Mikrotik) DeletePool(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := strings.Split(fmt.Sprintf("/ip/pool/remove =.id=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
