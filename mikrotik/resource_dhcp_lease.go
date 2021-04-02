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
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"blocked": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dynamic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceLeaseCreate(d *schema.ResourceData, m interface{}) error {
	dhcpLease := prepareDhcpLease(d)

	c := m.(client.Mikrotik)

	lease, err := c.AddDhcpLease(dhcpLease)
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

	currentLease, err := c.FindDhcpLease(d.Id())
	dhcpLease := prepareDhcpLease(d)
	dhcpLease.Id = currentLease.Id

	lease, err := c.UpdateDhcpLease(dhcpLease)
	lease.Dynamic = dhcpLease.Dynamic

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
	d.Set("blocked", lease.BlockAccess)
	d.Set("comment", lease.Comment)
	d.Set("address", lease.Address)
	d.Set("macaddress", lease.MacAddress)
	d.Set("hostname", lease.Hostname)
	d.Set("dynamic", lease.Dynamic)
	return nil
}

func prepareDhcpLease(d *schema.ResourceData) *client.DhcpLease {
	lease := new(client.DhcpLease)

	lease.BlockAccess = d.Get("blocked").(bool)
	lease.Comment = d.Get("comment").(string)
	lease.Address = d.Get("address").(string)
	lease.MacAddress = d.Get("macaddress").(string)
	lease.Hostname = d.Get("hostname").(string)
	lease.Dynamic = d.Get("dynamic").(bool)

	return lease
}
