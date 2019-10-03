package client

import (
	"fmt"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type DhcpLease struct {
	Id		string
	Address		string
	MacAddress	string
}

func (client Mikrotik) AddDhcpLease(address, macaddress string) (*routeros.Reply, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/add =address=%s =mac-address=%s", address, macaddress), " ")
	fmt.Println(fmt.Sprintf("[INFO] Running the mikrotik command: `%s`", cmd))
	r, err := c.RunArgs(cmd)
	return r, err
}

func (client Mikrotik) FindDhcpLease(address string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := "/ip/dhcp-server/lease/print"
	fmt.Println(fmt.Sprintf("[INFO] Running the mikrotik command: `%s`", cmd))
	r, err := c.Run(cmd)
	found := false
	var sentence *proto.Sentence

	if err != nil {
		return nil, err
	}

	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Value == address {
				found = true
				sentence = reply
				fmt.Println(fmt.Sprintf("[DEBUG] Found dhcp lease we were looking for: %v", sentence))
			}
		}
	}

	if !found {
		return nil, nil
	}

	// TODO: Add error checking

	macaddress := ""
	id := ""
	for _, pair := range sentence.List {
		if pair.Key == ".id" {
			id = pair.Value
		}
		if pair.Key == "mac-address" {
			macaddress = pair.Value
		}
	}

	return &DhcpLease{
		Id:      id,
		MacAddress: macaddress,
		Address:    address,
	}, nil
}

func (client Mikrotik) UpdateDhcpLease(id, address, macaddress string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/set =numbers=%s =address=%s =mac-address=%s", id, address, macaddress), " ")
	fmt.Println(fmt.Sprintf("[INFO] Running the mikrotik command: `%s`", cmd))
	_, err = c.RunArgs(cmd)
	return err
}

func (client Mikrotik) DeleteDhcpLease(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/remove =numbers=%s", id), " ")
	fmt.Println(fmt.Sprintf("[INFO] Running the mikrotik command: `%s`", cmd))
	_, err = c.RunArgs(cmd)
	return err
}
