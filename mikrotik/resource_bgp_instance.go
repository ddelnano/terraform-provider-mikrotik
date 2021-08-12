package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBgpInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBgpInstanceCreate,
		ReadContext:   resourceBgpInstanceRead,
		UpdateContext: resourceBgpInstanceUpdate,
		DeleteContext: resourceBgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"as": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"client_to_client_reflection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confederation_peers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ignore_as_path_len": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"out_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"redistribute_connected": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_ospf": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_other_bgp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_rip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_static": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_table": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confederation": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceBgpInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instance := prepareBgpInstance(d)

	c := m.(client.Mikrotik)

	bgpInstance, err := c.AddBgpInstance(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	bgpInstance, err := c.FindBgpInstance(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	currentBgpInstance, err := c.FindBgpInstance(d.Get("name").(string))

	instance := prepareBgpInstance(d)
	instance.ID = currentBgpInstance.ID

	bgpInstance, err := c.UpdateBgpInstance(instance)

	if err != nil {
		return diag.FromErr(err)
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	err := c.DeleteBgpInstance(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func bgpInstanceToData(b *client.BgpInstance, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"name":                        b.Name,
		"as":                          b.As,
		"client_to_client_reflection": b.ClientToClientReflection,
		"comment":                     b.Comment,
		"confederation_peers":         b.ConfederationPeers,
		"disabled":                    b.Disabled,
		"ignore_as_path_len":          b.IgnoreAsPathLen,
		"out_filter":                  b.OutFilter,
		"redistribute_connected":      b.RedistributeConnected,
		"redistribute_ospf":           b.RedistributeOspf,
		"redistribute_other_bgp":      b.RedistributeOtherBgp,
		"redistribute_rip":            b.RedistributeRip,
		"redistribute_static":         b.RedistributeStatic,
		"router_id":                   b.RouterID,
		"routing_table":               b.RoutingTable,
		"cluster_id":                  b.ClusterID,
		"confederation":               b.Confederation,
	}

	d.SetId(b.Name)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func prepareBgpInstance(d *schema.ResourceData) *client.BgpInstance {
	bgpInstance := new(client.BgpInstance)

	bgpInstance.Name = d.Get("name").(string)
	bgpInstance.As = d.Get("as").(int)
	bgpInstance.ClientToClientReflection = d.Get("client_to_client_reflection").(bool)
	bgpInstance.Comment = d.Get("comment").(string)
	bgpInstance.ConfederationPeers = d.Get("confederation_peers").(string)
	bgpInstance.Disabled = d.Get("disabled").(bool)
	bgpInstance.IgnoreAsPathLen = d.Get("ignore_as_path_len").(bool)
	bgpInstance.OutFilter = d.Get("out_filter").(string)
	bgpInstance.RedistributeConnected = d.Get("redistribute_connected").(bool)
	bgpInstance.RedistributeOspf = d.Get("redistribute_ospf").(bool)
	bgpInstance.RedistributeOtherBgp = d.Get("redistribute_other_bgp").(bool)
	bgpInstance.RedistributeRip = d.Get("redistribute_rip").(bool)
	bgpInstance.RedistributeStatic = d.Get("redistribute_static").(bool)
	bgpInstance.RouterID = d.Get("router_id").(string)
	bgpInstance.RoutingTable = d.Get("routing_table").(string)
	bgpInstance.ClusterID = d.Get("cluster_id").(string)
	bgpInstance.Confederation = d.Get("confederation").(int)

	return bgpInstance
}
