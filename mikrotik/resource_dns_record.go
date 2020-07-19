package mikrotik

import (
	"log"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

	record, err := c.AddDnsRecord(name, address, ttl)
	if err != nil {
		return err
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

	log.Printf("[DEBUG] About to update dns record with %v", record)
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
