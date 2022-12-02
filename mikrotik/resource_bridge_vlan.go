package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBridgeVlan() *schema.Resource {
	return &schema.Resource{
		Description: "New resource",

		CreateContext: resourceBridgeVlanCreate,
		ReadContext:   resourceBridgeVlanRead,
		UpdateContext: resourceBridgeVlanUpdate,
		DeleteContext: resourceBridgeVlanDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bridge": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"tagged": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "",
			},
			"untagged": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "",
			},
			"vlan_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "",
			},
		},
	}
}

func resourceBridgeVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r, err := c.AddBridgeVlan(dataToBridgeVlan(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.Id)

	return resourceBridgeVlanRead(ctx, d, m)
}

func resourceBridgeVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r, err := c.FindBridgeVlan(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return recordBridgeVlanToData(r, d)
}

func resourceBridgeVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r, err := c.UpdateBridgeVlan(dataToBridgeVlan(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return recordBridgeVlanToData(r, d)
}

func resourceBridgeVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	if err := c.DeleteBridgeVlan(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataToBridgeVlan(d *schema.ResourceData) *client.BridgeVlan {
	taggedInterface := d.Get("tagged").([]interface{})
	tagged := make([]string, len(taggedInterface))
	for i, v := range taggedInterface {
		tagged[i] = v.(string)
	}

	untaggedInterface := d.Get("untagged").([]interface{})
	untagged := make([]string, len(untaggedInterface))
	for i, v := range untaggedInterface {
		untagged[i] = v.(string)
	}

	vlanIDsInterface := d.Get("vlan_ids").([]interface{})
	vlanIDs := make([]int, len(vlanIDsInterface))
	for i, v := range vlanIDsInterface {
		vlanIDs[i] = v.(int)
	}

	return &client.BridgeVlan{
		Id:       d.Id(),
		Bridge:   d.Get("bridge").(string),
		Tagged:   tagged,
		Untagged: untagged,
		VlanIds:  vlanIDs,
	}
}

func recordBridgeVlanToData(r *client.BridgeVlan, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("bridge", r.Bridge); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("tagged", r.Tagged); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("untagged", r.Untagged); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("vlan_ids", r.VlanIds); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.SetId(r.Id)

	return diags
}
