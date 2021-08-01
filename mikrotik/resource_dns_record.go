package mikrotik

import (
	"context"
	"log"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	record := prepareDnsRecord(d)

	c := m.(client.Mikrotik)

	dnsRecord, err := c.AddDnsRecord(record)
	if err != nil {
		return diag.FromErr(err)
	}

	err = recordToData(dnsRecord, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	err = recordToData(record, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	currentRecord, err := c.FindDnsRecord(d.Id())
	record := prepareDnsRecord(d)
	record.Id = currentRecord.Id

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] About to update dns record with %v", record)
	dnsRecord, err := c.UpdateDnsRecord(record)
	if err != nil {
		return diag.FromErr(err)
	}

	err = recordToData(dnsRecord, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()

	c := m.(client.Mikrotik)

	record, err := c.FindDnsRecord(name)

	if err != nil {
		return diag.FromErr(err)
	}
	err = c.DeleteDnsRecord(record.Id)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func recordToData(record *client.DnsRecord, d *schema.ResourceData) error {
	d.SetId(record.Name)
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
