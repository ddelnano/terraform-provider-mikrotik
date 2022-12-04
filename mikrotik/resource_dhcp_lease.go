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
		Description: "Creates a DHCP lease on the mikrotik device.",

		CreateContext: resourceLeaseCreate,
		ReadContext:   resourceLeaseRead,
		UpdateContext: resourceLeaseUpdate,
		DeleteContext: resourceLeaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address of the DHCP lease to be created.",
			},
			"macaddress": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MAC addreess of the DHCP lease to be created.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment of the DHCP lease to be created.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the device",
			},
			"blocked": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "Whether to block access for this DHCP client (true|false).",
			},
			"dynamic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the dhcp lease is static or dynamic. Dynamic leases are not guaranteed to continue to be assigned to that specific device. Defaults to false.",
			},
		},
	}
}

func resourceLeaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dhcpLease := prepareDhcpLease(d)

	c := m.(*client.Mikrotik)

	lease, err := c.AddDhcpLease(dhcpLease)
	if err != nil {
		return diag.FromErr(err)
	}

	return leaseToData(lease, d)
}

func resourceLeaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	lease, err := c.FindDhcpLease(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if lease == nil {
		d.SetId("")
		return nil
	}

	return leaseToData(lease, d)
}

func resourceLeaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	dhcpLease := prepareDhcpLease(d)
	dhcpLease.Id = d.Id()

	lease, err := c.UpdateDhcpLease(dhcpLease)
	if err != nil {
		return diag.FromErr(err)
	}
	lease.Dynamic = dhcpLease.Dynamic

	return leaseToData(lease, d)
}

func resourceLeaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteDhcpLease(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func leaseToData(lease *client.DhcpLease, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"blocked":    strconv.FormatBool(lease.BlockAccess),
		"comment":    lease.Comment,
		"address":    lease.Address,
		"macaddress": lease.MacAddress,
		"hostname":   lease.Hostname,
		"dynamic":    lease.Dynamic,
	}

	d.SetId(lease.Id)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
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
