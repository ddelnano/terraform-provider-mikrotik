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
	record := prepareDnsRecord(d)

	c := m.(client.Mikrotik)

	dnsRecord, err := c.AddDnsRecord(record)
	if err != nil {
		return err
	}

	return recordToData(dnsRecord, d)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	return recordToData(record, d)
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	currentRecord, err := c.FindDnsRecord(d.Id())
	record := prepareDnsRecord(d)
	record.Id = currentRecord.Id

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] About to update dns record with %v", record)
	dnsRecord, err := c.UpdateDnsRecord(record)

	if err != nil {
		return err
	}

	return recordToData(dnsRecord, d)
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

func prepareDnsRecord(d *schema.ResourceData) *client.DnsRecord {
	dnsRecord := new(client.DnsRecord)

	dnsRecord.Name = d.Get("name").(string)
	dnsRecord.Ttl = d.Get("ttl").(int)
	dnsRecord.Address = d.Get("address").(string)

	return dnsRecord
}
