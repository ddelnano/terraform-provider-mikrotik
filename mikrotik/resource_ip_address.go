package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIpAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAddressCreate,
		ReadContext:   resourceIpAddressRead,
		UpdateContext: resourceIpAddressUpdate,
		DeleteContext: resourceIpAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIpAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ipAddress := prepareIpAddress(d)

	c := m.(*client.Mikrotik)

	ipaddr, err := c.AddIpAddress(ipAddress)

	if err != nil {
		return diag.FromErr(err)
	}

	return addrToData(ipaddr, d)
}

func resourceIpAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	ipaddr, err := c.FindIpAddress(d.Id())

	// Clear the state if the error represents that the resource no longer exists
	_, resourceMissing := err.(*client.NotFound)
	if resourceMissing && err != nil {
		d.SetId("")
		return nil
	}

	// Make sure all other errors are propagated
	if err != nil {
		return diag.FromErr(err)
	}

	return addrToData(ipaddr, d)
}

func resourceIpAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	ipAddress := prepareIpAddress(d)
	ipAddress.Id = d.Id()

	ipaddr, err := c.UpdateIpAddress(ipAddress)

	if err != nil {
		return diag.FromErr(err)
	}

	return addrToData(ipaddr, d)
}

func resourceIpAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteIpAddress(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func addrToData(ipaddr *client.IpAddress, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"address":   ipaddr.Address,
		"comment":   ipaddr.Comment,
		"disabled":  ipaddr.Disabled,
		"interface": ipaddr.Interface,
		"network":   ipaddr.Network,
	}

	d.SetId(ipaddr.Id)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func prepareIpAddress(d *schema.ResourceData) *client.IpAddress {
	ipaddr := new(client.IpAddress)

	ipaddr.Comment = d.Get("comment").(string)
	ipaddr.Address = d.Get("address").(string)
	ipaddr.Disabled = d.Get("disabled").(bool)
	ipaddr.Interface = d.Get("interface").(string)
	ipaddr.Network = d.Get("network").(string)

	return ipaddr
}
