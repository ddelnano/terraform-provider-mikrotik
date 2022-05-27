package client

import (
	"fmt"
	"log"
)

type Ipv6Address struct {
	Id        string `mikrotik:".id"`
	Address   string `mikrotik:"address"`
	Advertise bool   `mikrotik:"advertise"`
	Comment   string `mikrotik:"comment"`
	Disabled  bool   `mikrotik:"disabled"`
	Eui64     bool   `mikrotik:"eui-64"`
	FromPool  string `mikrotik:"from-pool"`
	Interface string `mikrotik:"interface"`
	NoDad     bool   `mikrotik:"no-dad"`
}

func (client Mikrotik) AddIpv6Address(addr *Ipv6Address) (*Ipv6Address, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	cmd := Marshal("/ipv6/address/add", addr)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] ipv6 address creation response: `%v`", r)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.FindIpv6Address(id)
}

func (client Mikrotik) ListIpv6Address() ([]Ipv6Address, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/ipv6/address/print"}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] found ipv6 address: %v", r)

	ipv6addr := []Ipv6Address{}

	err = Unmarshal(*r, &ipv6addr)

	if err != nil {
		return nil, err
	}

	return ipv6addr, nil
}
func (client Mikrotik) FindIpv6Address(id string) (*Ipv6Address, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/ipv6/address/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] ipv6 address response: %v", r)
	if err != nil {
		return nil, err
	}

	ipv6addr := Ipv6Address{}
	err = Unmarshal(*r, &ipv6addr)
	if err != nil {
		return nil, err
	}

	if ipv6addr.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("ipv6 address `%s` not found", id))
	}

	return &ipv6addr, nil
}

func (client Mikrotik) UpdateIpv6Address(addr *Ipv6Address) (*Ipv6Address, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ipv6/address/set", addr)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindIpv6Address(addr.Id)
}

func (client Mikrotik) DeleteIpv6Address(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := []string{"/ipv6/address/remove", "=.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
