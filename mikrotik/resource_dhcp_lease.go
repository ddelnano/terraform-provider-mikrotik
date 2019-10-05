package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLease() *schema.Resource {
	return &schema.Resource{
		Create: resourceLeaseCreate,
		Read:   resourceLeaseRead,
		Update: resourceLeaseUpdate,
		Delete: resourceLeaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"macaddress": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLeaseCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	macaddress := d.Get("macaddress").(string)
	comment := d.Get("comment").(string)

	c := m.(client.Mikrotik)

	lease, err := c.AddDhcpLease(address, macaddress, comment)
	if err != nil {
		return err
	}

	leaseToData(lease, d)
	return nil
}

func resourceLeaseRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	lease, err := c.FindDhcpLease(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	if lease == nil {
		d.SetId("")
		return nil
	}

	leaseToData(lease, d)
	return nil
}

func resourceLeaseUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	macaddress := d.Get("macaddress").(string)
	address := d.Get("address").(string)
	comment := d.Get("comment").(string)

	lease, err := c.UpdateDhcpLease(d.Id(), address, macaddress, comment)

	if err != nil {
		return err
	}

	leaseToData(lease, d)
	return nil
}

func resourceLeaseDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	err := c.DeleteDhcpLease(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func leaseToData(lease *client.DhcpLease, d *schema.ResourceData) error {
	d.SetId(lease.Id)
	d.Set("comment", lease.Comment)
	d.Set("address", lease.Address)
	d.Set("macaddress", lease.MacAddress)
	return nil
}
