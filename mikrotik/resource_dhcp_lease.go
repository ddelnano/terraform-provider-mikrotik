package mikrotik

import (
	"context"
	"strconv"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLeaseCreate,
		ReadContext:   resourceLeaseRead,
		UpdateContext: resourceLeaseUpdate,
		DeleteContext: resourceLeaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"macaddress": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"blocked": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
			"dynamic": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceLeaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dhcpLease := prepareDhcpLease(d)

	c := m.(client.Mikrotik)

	lease, err := c.AddDhcpLease(dhcpLease)
	if err != nil {
		return diag.FromErr(err)
	}

	err = leaseToData(lease, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLeaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	err = leaseToData(lease, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLeaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	dhcpLease := prepareDhcpLease(d)
	dhcpLease.Id = d.Id()

	lease, err := c.UpdateDhcpLease(dhcpLease)
	lease.Dynamic = dhcpLease.Dynamic

	if err != nil {
		return diag.FromErr(err)
	}

	err = leaseToData(lease, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLeaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	err := c.DeleteDhcpLease(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func leaseToData(lease *client.DhcpLease, d *schema.ResourceData) error {
	d.SetId(lease.Id)
	d.Set("blocked", strconv.FormatBool(lease.BlockAccess))
	d.Set("comment", lease.Comment)
	d.Set("address", lease.Address)
	d.Set("macaddress", lease.MacAddress)
	d.Set("hostname", lease.Hostname)
	d.Set("dynamic", lease.Dynamic)
	return nil
}

func prepareDhcpLease(d *schema.ResourceData) *client.DhcpLease {
	lease := new(client.DhcpLease)

	lease.BlockAccess, _ = strconv.ParseBool(d.Get("blocked").(string))
	lease.Comment = d.Get("comment").(string)
	lease.Address = d.Get("address").(string)
	lease.MacAddress = d.Get("macaddress").(string)
	lease.Hostname = d.Get("hostname").(string)
	lease.Dynamic = d.Get("dynamic").(bool)

	return lease
}
