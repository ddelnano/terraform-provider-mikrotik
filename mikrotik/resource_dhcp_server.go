package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDhcpServer() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a DHCP server resource within MikroTik device.",

		CreateContext: createDhcpServer,
		ReadContext:   readDhcpServer,
		UpdateContext: updateDhcpServer,
		DeleteContext: deleteDhcpServer,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_arp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to add dynamic ARP entry. If set to no either ARP mode should be enabled on that interface or static ARP entries should be administratively defined.",
			},
			"address_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "static-only",
				Description: "IP pool, from which to take IP addresses for the clients. If set to static-only, then only the clients that have a static lease (added in lease submenu) will be allowed.",
			},
			"authoritative": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yes",
				Description: "Option changes the way how server responds to DHCP requests.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Disable this DHCP server instance.",
			},
			"interface": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*0",
				Description: "Interface on which server will be running.",
			},
			"lease_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script that will be executed after lease is assigned or de-assigned. Internal \"global\" variables that can be used in the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference name.",
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
