package mikrotik

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: RecordImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)

	c := m.(*routeros.Client)
	// TODO: Provide some basic validation here
	r, err := c.RunArgs(strings.Split(fmt.Sprintf("/ip/dns/static/add =name=%s =address=%s =ttl=%d", name, address, ttl), " "))
	if err != nil {
		return err
	}

	// If API is successful we should only get a single sentence and list back like so
	// 2019/02/28 20:13:15 !done @ [{`ret` `*14`}]
	var id string
	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Key == ".id" {
				id = item.Value
			}
		}
	}

	record := &dnsRecord{
		id:      id,
		address: address,
		name:    name,
		ttl:     ttl,
	}
	recordToData(record, d)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*routeros.Client)
	record, err := findDnsRecord(c, d.Id())

	// TODO: Ignoring this error can cause all resources to think they
	// need to be created. We should more appropriately handle this. The
	// error where the DNS record is not found is not actually an error and
	// needs to be disambiguated from real failures
	if err != nil {
		d.SetId("")
		return nil
	}

	recordToData(record, d)
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	// return resourceServerRead(d, m)
	c := m.(*routeros.Client)
	address := d.Get("address").(string)
	ttl := d.Get("ttl").(int)
	name := d.Id()

	record, err := findDnsRecord(c, name)

	if err != nil {
		return err
	}

	_, err = c.RunArgs(strings.Split(fmt.Sprintf("/ip/dns/static/set =numbers=%s =name=%s =address=%s =ttl=%d", record.id, name, address, ttl), " "))

	if err != nil {
		return err
	}

	recordToData(record, d)
	return nil
}

// !re @ [{`.id` `*2`} {`name` `radarr`} {`address` `192.168.88.254`} {`ttl` `59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*3`} {`name` `sonarr`} {`address` `192.168.88.254`} {`ttl` `59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*4`} {`name` `sabnzbd`} {`address` `192.168.88.254`} {`ttl` `59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*5`} {`name` `kodi`} {`address` `192.168.88.244`} {`ttl` `59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*6`} {`name` `osmc`} {`address` `192.168.88.244`} {`ttl` `59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*7`} {`name` `radarr.internal.ddelnano.com`} {`address` `192.168.88.254`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*8`} {`name` `sonarr.internal.ddelnano.com`} {`address` `192.168.88.254`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*9`} {`name` `sabnzbd.internal.ddelnano.com`} {`address` `192.168.88.254`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*A`} {`name` `router.internal.ddelnano.com`} {`address` `192.168.88.1`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*B`} {`name` `kodi.internal.ddelnano.com`} {`address` `192.168.88.244`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*D`} {`name` `osmc.internal.ddelnano.com`} {`address` `192.168.88.244`} {`ttl` `1m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*E`} {`name` `switch.internal.ddelnano.com`} {`address` `192.168.88.90`} {`ttl` `5m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*F`} {`name` `xen.internal.ddelnano.com`} {`address` `192.168.88.117`} {`ttl` `5m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`} {`comment` `DNS for the hypervisor`}]
// !re @ [{`.id` `*10`} {`name` `xoa.internal.ddelnano.com`} {`address` `192.168.88.86`} {`ttl` `5m`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]
// !re @ [{`.id` `*15`} {`name` `test`} {`address` `10.0.0.1`} {`ttl` `1d15h20m59s`} {`dynamic` `false`} {`regexp` `false`} {`disabled` `false`}]

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	c := m.(*routeros.Client)
	record, err := findDnsRecord(c, name)

	if err != nil {
		return err
	}
	_, err = c.RunArgs(strings.Split(fmt.Sprintf("/ip/dns/static/remove =numbers=%s", record.id), " "))

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func recordToData(record *dnsRecord, d *schema.ResourceData) error {
	d.SetId(record.name)
	d.Set("numerical_id", record.id)
	d.Set("name", record.name)
	d.Set("address", record.address)
	d.Set("ttl", record.ttl)
	return nil
}

type dnsRecord struct {
	// .id field that mikrotik uses as the 'real' ID
	id      string
	name    string
	ttl     int
	address string
}

// TODO: Why does /print seem to return ID's as hex but /add seems to always return the ID as a decimal number
func findDnsRecord(c *routeros.Client, name string) (*dnsRecord, error) {
	r, err := c.Run("/ip/dns/static/print")
	found := false
	var sentence *proto.Sentence

	if err != nil {
		return nil, err
	}

	fmt.Println(r)
	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Value == name {
				found = true
				sentence = reply
			}
		}
	}

	if !found {
		return nil, errors.New("Resource was not found")
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

	return &dnsRecord{
		id:      id,
		address: address,
		name:    name,
		ttl:     ttlToSeconds(ttl),
	}, nil
}

// RecordImport - import record from existing mikrotik api. ID is specified by the address (google.com)
func RecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	c := m.(*routeros.Client)
	record, err := findDnsRecord(c, name)

	if err != nil {
		return nil, err
	}
	recordToData(record, d)
	return []*schema.ResourceData{d}, nil
}

func ttlToSeconds(ttl string) int {
	parts := strings.Split(ttl, "d")

	idx := 0
	days := 0
	var err error
	fmt.Println(parts)
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
