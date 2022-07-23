package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIpv6Address() *schema.Resource {
	return &schema.Resource{
		Description: "Assigns an IPv6 address to an interface.",

		CreateContext: resourceIpv6AddressCreate,
		ReadContext:   resourceIpv6AddressRead,
		UpdateContext: resourceIpv6AddressUpdate,
		DeleteContext: resourceIpv6AddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IPv6 address and prefix length of the interface using slash notation.",
			},
			"advertise": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable stateless address configuration. The prefix of that address is automatically advertised to hosts using ICMPv6 protocol. The option is set by default for addresses with prefix length 64.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment for the IPv6 address assignment.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to disable IPv6 address.",
			},
			"eui_64": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to calculate EUI-64 address and use it as last 64 bits of the IPv6 address.",
			},
			"from_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the pool from which prefix will be taken to construct IPv6 address taking last part of the address from address property.",
			},
			"interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The interface on which the IPv6 address is assigned.",
			},
			"no_dad": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set indicates that address is anycast address and Duplicate Address Detection should not be performed.",
			},
		},
	}
}

func resourceIpv6AddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ipv6Address := prepareIpv6Address(d)

	c := m.(*client.Mikrotik)

	ipv6addr, err := c.AddIpv6Address(ipv6Address)

	if err != nil {
		return diag.FromErr(err)
	}

	return v6addrToData(ipv6addr, d)
}

func resourceIpv6AddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	ipv6addr, err := c.FindIpv6Address(d.Id())

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

	return v6addrToData(ipv6addr, d)
}

func resourceIpv6AddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	ipv6Address := prepareIpv6Address(d)
	ipv6Address.Id = d.Id()

	ipv6addr, err := c.UpdateIpv6Address(ipv6Address)

	if err != nil {
		return diag.FromErr(err)
	}

	return v6addrToData(ipv6addr, d)
}

func resourceIpv6AddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteIpv6Address(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func v6addrToData(ipv6addr *client.Ipv6Address, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"address":   ipv6addr.Address,
		"advertise": ipv6addr.Advertise,
		"comment":   ipv6addr.Comment,
		"disabled":  ipv6addr.Disabled,
		"eui_64":    ipv6addr.Eui64,
		"from_pool": ipv6addr.FromPool,
		"interface": ipv6addr.Interface,
		"no_dad":    ipv6addr.NoDad,
	}

	d.SetId(ipv6addr.Id)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func prepareIpv6Address(d *schema.ResourceData) *client.Ipv6Address {
	ipv6addr := new(client.Ipv6Address)

	ipv6addr.Address = d.Get("address").(string)
	ipv6addr.Advertise = d.Get("advertise").(bool)
	ipv6addr.Comment = d.Get("comment").(string)
	ipv6addr.Disabled = d.Get("disabled").(bool)
	ipv6addr.Eui64 = d.Get("eui_64").(bool)
	ipv6addr.FromPool = d.Get("from_pool").(string)
	ipv6addr.Interface = d.Get("interface").(string)
	ipv6addr.NoDad = d.Get("no_dad").(bool)

	return ipv6addr
}
