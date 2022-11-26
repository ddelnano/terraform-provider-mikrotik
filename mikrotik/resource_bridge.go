package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBridge() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a bridge resource on remote MikroTik device.",

		CreateContext: resourceBridgeCreate,
		ReadContext:   resourceBridgeRead,
		UpdateContext: resourceBridgeUpdate,
		DeleteContext: resourceBridgeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the bridge interface",
			},
			"fast_forward": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Special and faster case of FastPath which works only on bridges with 2 interfaces (enabled by default only for new bridges).",
			},
			"vlan_filtering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Globally enables or disables VLAN functionality for bridge.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Short description of the interface.",
			},
		},
	}
}

func resourceBridgeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	bridge, err := c.AddBridge(dataToBridge(d))
	if err != nil {
		return diag.FromErr(err)
	}

	recordBridgeToData(bridge, d)
	d.SetId(bridge.Name)

	return resourceBridgeRead(ctx, d, m)
}

func resourceBridgeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	bridge, err := c.FindBridge(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	recordBridgeToData(bridge, d)

	return diags
}

func resourceBridgeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	bridge := dataToBridge(d)
	updatedBridge, err := c.UpdateBridge(bridge)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(updatedBridge.Name)

	return diags
}

func resourceBridgeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	err := c.DeleteBridge(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToBridge(d *schema.ResourceData) *client.Bridge {
	return &client.Bridge{
		Id:            d.Id(),
		Name:          d.Get("name").(string),
		FastForward:   d.Get("fast_forward").(bool),
		VlanFiltering: d.Get("vlan_filtering").(bool),
		Comment:       d.Get("comment").(string),
	}
}

func recordBridgeToData(r *client.Bridge, d *schema.ResourceData) {
	d.Set("name", r.Name)
	d.Set("fast_forward", r.FastForward)
	d.Set("vlan_filtering", r.VlanFiltering)
	d.Set("comment", r.Comment)
}
