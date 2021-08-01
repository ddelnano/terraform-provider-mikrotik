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

	err = bgpInstanceToData(bgpInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBgpInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	bgpInstance, err := c.FindBgpInstance(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	err = bgpInstanceToData(bgpInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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

	err = bgpInstanceToData(bgpInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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

func bgpInstanceToData(b *client.BgpInstance, d *schema.ResourceData) error {
	d.SetId(b.Name)

	if err := d.Set("name", b.Name); err != nil {
		return err
	}
	if err := d.Set("as", b.As); err != nil {
		return err
	}
	if err := d.Set("client_to_client_reflection", b.ClientToClientReflection); err != nil {
		return err
	}
	if err := d.Set("comment", b.Comment); err != nil {
		return err
	}
	if err := d.Set("confederation_peers", b.ConfederationPeers); err != nil {
		return err
	}
	if err := d.Set("disabled", b.Disabled); err != nil {
		return err
	}
	if err := d.Set("ignore_as_path_len", b.IgnoreAsPathLen); err != nil {
		return err
	}
	if err := d.Set("out_filter", b.OutFilter); err != nil {
		return err
	}
	if err := d.Set("redistribute_connected", b.RedistributeConnected); err != nil {
		return err
	}
	if err := d.Set("redistribute_ospf", b.RedistributeOspf); err != nil {
		return err
	}
	if err := d.Set("redistribute_other_bgp", b.RedistributeOtherBgp); err != nil {
		return err
	}
	if err := d.Set("redistribute_rip", b.RedistributeRip); err != nil {
		return err
	}
	if err := d.Set("redistribute_static", b.RedistributeStatic); err != nil {
		return err
	}
	if err := d.Set("router_id", b.RouterID); err != nil {
		return err
	}
	if err := d.Set("routing_table", b.RoutingTable); err != nil {
		return err
	}
	if err := d.Set("cluster_id", b.ClusterID); err != nil {
		return err
	}
	if err := d.Set("confederation", b.Confederation); err != nil {
		return err
	}
	return nil
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
