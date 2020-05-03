package client

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros"
)

type DnsRecord struct {
	Id      string `mikrotik:".id"`
	Name    string
	Ttl     int `mikrotik:"ttl,ttlToSeconds"`
	Address string
}

func ttlToSeconds(ttl string) int {
	parts := strings.Split(ttl, "d")

	idx := 0
	days := 0
	var err error
	if len(parts) == 2 {
		idx = 1
		days, err = strconv.Atoi(parts[0])

		// We should be parsing an ascii number
		// if this fails we should fail loudly
		if err != nil {
			panic(err)
		}

		// In the event we just get days parts[1] will be an
		// empty string. Just coerce that into 0 seconds.
		if parts[1] == "" {
			parts[1] = "0s"
		}
	}
	d, err := time.ParseDuration(parts[idx])

	// We should never receive a duration greater than
	// 23h59m59s. So this should always parse.
	if err != nil {
		panic(err)
	}
	return 86400*days + int(d)/int(math.Pow10(9))

}

func (client Mikrotik) AddDnsRecord(name, address string, ttl int) (*routeros.Reply, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dns/static/add =name=%s =address=%s =ttl=%d", name, address, ttl), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	return r, err
}

func (client Mikrotik) FindDnsRecord(name string) (*DnsRecord, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dns/static/print ?name=%s", name), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Found dns record: %v", r)

	record := DnsRecord{}
	err = Unmarshal(*r, &record)

	if err != nil {
		return nil, err
	}

	if record.Name == "" {
		return nil, nil
	}

	return &record, nil
}

func (client Mikrotik) UpdateDnsRecord(id, name, address string, ttl int) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dns/static/set =numbers=%s =name=%s =address=%s =ttl=%d", id, name, address, ttl), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}

func (client Mikrotik) DeleteDnsRecord(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dns/static/remove =numbers=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
