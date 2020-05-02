package client

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type Mikrotik struct {
	Host     string
	Username string
	Password string
}

type DnsRecord struct {
	// .id field that mikrotik uses as the 'real' ID
	Id      string
	Name    string
	Ttl     int
	Address string
}

func Unmarshal(reply routeros.Reply, v interface{}) error {
	rv := reflect.ValueOf(v)
	elem := rv.Elem()

	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal cannot work without a pointer")
	}

	switch elem.Kind() {
	case reflect.Slice:
	case reflect.Struct:
		if len(reply.Re) < 1 {
			// This is an empty message
			return nil
		}
		if len(reply.Re) > 1 {
			msg := fmt.Sprintf("Failed to decode reply: %v", reply)
			return errors.New(msg)
		}

		for i := 0; i < elem.NumField(); i++ {
			field := elem.Field(i)
			fieldType := elem.Type().Field(i)
			tag := fieldType.Tag.Get("mikrotik")

			path := strings.ToLower(fieldType.Name)

			for _, pair := range reply.Re[0].List {
				if strings.Compare(pair.Key, path) == 0 || strings.Compare(pair.Key, tag) == 0 {
					switch fieldType.Type.Kind() {
					case reflect.String:
						field.SetString(pair.Value)
					case reflect.Bool:
						b, _ := strconv.ParseBool(pair.Value)
						field.SetBool(b)
					}

				}
			}
		}
	}

	return nil
}

func NewClient(host, username, password string) Mikrotik {
	return Mikrotik{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func GetConfigFromEnv() (host, username, password string) {
	host = os.Getenv("MIKROTIK_HOST")
	username = os.Getenv("MIKROTIK_USER")
	password = os.Getenv("MIKROTIK_PASSWORD")
	if host == "" || username == "" || password == "" {
		// panic("Unable to find the MIKROTIK_HOST, MIKROTIK_USER or MIKROTIK_PASSWORD environment variable")
	}
	return host, username, password
}

func (client Mikrotik) getMikrotikClient() (c *routeros.Client, err error) {
	address := client.Host
	username := client.Username
	password := client.Password
	c, err = routeros.Dial(address, username, password)

	if err != nil {
		log.Printf("[ERROR] Failed to login to routerOS with error: %v", err)
	}

	return
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
	cmd := "/ip/dns/static/print"
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.Run(cmd)
	found := false
	var sentence *proto.Sentence

	if err != nil {
		return nil, err
	}

	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Value == name {
				found = true
				sentence = reply
				log.Printf("[DEBUG] Found dns record we were looking for: %v", sentence)
			}
		}
	}

	if !found {
		return nil, nil
	}

	// TODO: Add error checking

	address := ""
	ttl := ""
	id := ""
	for _, pair := range sentence.List {
		if pair.Key == ".id" {
			id = pair.Value
		}
		if pair.Key == "address" {
			address = pair.Value
		}

		if pair.Key == "ttl" {
			ttl = pair.Value
		}
	}

	return &DnsRecord{
		Id:      id,
		Address: address,
		Name:    name,
		Ttl:     ttlToSeconds(ttl),
	}, nil
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
