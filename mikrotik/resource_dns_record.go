package mikrotik

import (
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
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

	c := m.(client.Mikrotik)

	r, err := c.AddDnsRecord(name, address, ttl)
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

	record := &client.DnsRecord{
		Id:      id,
		Address: address,
		Name:    name,
		Ttl:     ttl,
	}
	recordToData(record, d)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(d.Id())

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
	c := m.(client.Mikrotik)

	address := d.Get("address").(string)
	ttl := d.Get("ttl").(int)
	name := d.Id()

	record, err := c.FindDnsRecord(name)

	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] About to update dns record with %v", record)
	err = c.UpdateDnsRecord(record.Id, name, address, ttl)

	if err != nil {
		return err
	}

	// TODO: the c.UpdateDnsRecord call should return a
	// new DnsRecord instead of mutating the current one.
	record.Address = address
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

	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(name)

	if err != nil {
		return err
	}
	err = c.DeleteDnsRecord(record.Id)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func recordToData(record *client.DnsRecord, d *schema.ResourceData) error {
	d.SetId(record.Name)
	d.Set("numerical_id", record.Id)
	d.Set("name", record.Name)
	d.Set("address", record.Address)
	d.Set("ttl", record.Ttl)
	return nil
}

func RecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(name)

	if err != nil {
		return nil, err
	}
	recordToData(record, d)
	return []*schema.ResourceData{d}, nil
}
