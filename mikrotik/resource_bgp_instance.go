package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBgpInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a Mikrotik BGP Instance.",

		CreateContext: resourceBgpInstanceCreate,
		ReadContext:   resourceBgpInstanceRead,
		UpdateContext: resourceBgpInstanceUpdate,
		DeleteContext: resourceBgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the BGP instance.",
			},
			"as": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The 32-bit BGP autonomous system number. Must be a value within 0 to 4294967295.",
			},
			"client_to_client_reflection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The comment of the IP Pool to be created.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment of the BGP instance to be created.",
			},
			"confederation_peers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of AS numbers internal to the [local] confederation. For example: `10,20,30-50`.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether instance is disabled.",
			},
			"ignore_as_path_len": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to ignore AS_PATH attribute in BGP route selection algorithm.",
			},
			"out_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Output routing filter chain used by all BGP peers belonging to this instance.",
			},
			"redistribute_connected": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, this BGP instance will redistribute the information about connected routes.",
			},
			"redistribute_ospf": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, this BGP instance will redistribute the information about routes learned by OSPF.",
			},
			"redistribute_other_bgp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, this BGP instance will redistribute the information about routes learned by other BGP instances.",
			},
			"redistribute_rip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, this BGP instance will redistribute the information about routes learned by RIP.",
			},
			"redistribute_static": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the router will redistribute the information about static routes added to its routing database.",
			},
			"router_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "BGP Router ID (for this instance). If set to 0.0.0.0, BGP will use one of router's IP addresses.",
			},
			"routing_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of routing table this BGP instance operates on. ",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "In case this instance is a route reflector: cluster ID of the router reflector cluster this instance belongs to.",
			},
			"confederation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "In case of BGP confederations: autonomous system number that identifies the [local] confederation as a whole.",
			},
		},
	}
}

func resourceBgpInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instance := prepareBgpInstance(d)

	c := m.(*client.Mikrotik)

	bgpInstance, err := c.AddBgpInstance(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	bgpInstance, err := c.FindBgpInstance(d.Id())

	if _, ok := err.(client.LegacyBgpUnsupported); ok {
		return diag.FromErr(err)
	}

	if client.IsNotFoundError(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	currentBgpInstance, err := c.FindBgpInstance(d.Get("name").(string))
	if _, ok := err.(client.LegacyBgpUnsupported); ok {
		return diag.FromErr(err)
	}

	instance := prepareBgpInstance(d)
	instance.ID = currentBgpInstance.ID

	bgpInstance, err := c.UpdateBgpInstance(instance)

	if err != nil {
		return diag.FromErr(err)
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteBgpInstance(d.Get("name").(string))
	if _, ok := err.(client.LegacyBgpUnsupported); ok {
		return diag.FromErr(err)
	}

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
	return &client.BgpInstance{
		Name:                     d.Get("name").(string),
		As:                       d.Get("as").(int),
		ClientToClientReflection: d.Get("client_to_client_reflection").(bool),
		Comment:                  d.Get("comment").(string),
		ConfederationPeers:       d.Get("confederation_peers").(string),
		Disabled:                 d.Get("disabled").(bool),
		IgnoreAsPathLen:          d.Get("ignore_as_path_len").(bool),
		OutFilter:                d.Get("out_filter").(string),
		RedistributeConnected:    d.Get("redistribute_connected").(bool),
		RedistributeOspf:         d.Get("redistribute_ospf").(bool),
		RedistributeOtherBgp:     d.Get("redistribute_other_bgp").(bool),
		RedistributeRip:          d.Get("redistribute_rip").(bool),
		RedistributeStatic:       d.Get("redistribute_static").(bool),
		RouterID:                 d.Get("router_id").(string),
		RoutingTable:             d.Get("routing_table").(string),
		ClusterID:                d.Get("cluster_id").(string),
		Confederation:            d.Get("confederation").(int),
	}
}
