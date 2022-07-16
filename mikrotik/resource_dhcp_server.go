package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDhcpServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDhcpServer,
		ReadContext:   readDhcpServer,
		UpdateContext: updateDhcpServer,
		DeleteContext: deleteDhcpServer,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_arp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"address_pool": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "static-only",
			},
			"authoritative": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "yes",
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"interface": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*0",
			},
			"lease_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createDhcpServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	dhcpServer, err := c.AddDhcpServer(dataToDhcpServer(d))
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpServerToData(dhcpServer, d)
	d.SetId(dhcpServer.Name)

	return readDhcpServer(ctx, d, m)
}

func readDhcpServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	dhcpServer, err := c.FindDhcpServer(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpServerToData(dhcpServer, d)

	return diags
}

func updateDhcpServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	dhcpServer := dataToDhcpServer(d)
	_, err := c.UpdateDhcpServer(dhcpServer)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteDhcpServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	err := c.DeleteDhcpServer(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToDhcpServer(d *schema.ResourceData) *client.DhcpServer {
	return &client.DhcpServer{
		Id:            d.Id(),
		AddArp:        d.Get("add_arp").(bool),
		AddressPool:   d.Get("address_pool").(string),
		Authoritative: d.Get("authoritative").(string),
		Disabled:      d.Get("disabled").(bool),
		Interface:     d.Get("interface").(string),
		LeaseScript:   d.Get("lease_script").(string),
		Name:          d.Get("name").(string),
	}
}

func dhcpServerToData(dhcpServer *client.DhcpServer, d *schema.ResourceData) {
	d.Set("add_arp", dhcpServer.AddArp)
	d.Set("address_pool", dhcpServer.AddressPool)
	d.Set("authoritative", dhcpServer.Authoritative)
	d.Set("disabled", dhcpServer.Disabled)
	d.Set("interface", dhcpServer.Interface)
	d.Set("lease_script", dhcpServer.LeaseScript)
	d.Set("name", dhcpServer.Name)
}
