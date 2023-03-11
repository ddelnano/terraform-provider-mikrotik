package mikrotik

import (
	"context"
	"fmt"
	"strings"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/exp/slices"
)

func resourceBridgePort() *schema.Resource {
	return &schema.Resource{
		Description: "Manages ports in bridge associations.",

		CreateContext: resourceBridgePortCreate,
		ReadContext:   resourceBridgePortRead,
		UpdateContext: resourceBridgePortUpdate,
		DeleteContext: resourceBridgePortDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bridge": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bridge interface the respective interface is grouped in.",
			},
			"interface": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*0",
				Description: "Name of the interface.",
			},
			"pvid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 4094),
				Description:  "Port VLAN ID (pvid) specifies which VLAN the untagged ingress traffic is assigned to. This property only has effect when vlan-filtering is set to yes.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Short description for this association.",
			},
			"frame_types": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "admit-all",
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					expected := []string{"admit-all", "admit-only-untagged-and-priority-tagged", "admit-only-vlan-tagged"}
					value := v.(string)

					var diags diag.Diagnostics
					if !slices.Contains(expected, value) {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Wrong value",
							Detail:   fmt.Sprintf("%q is not part of the possible values: %q", value, strings.Join(expected, ",")),
						}
						diags = append(diags, diag)
					}
					return diags
				},
				Description: "Can be used to filter out packets whether they have a VLAN tag or not.",
			},
		},
	}
}

func resourceBridgePortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	bridgePort, err := c.AddBridgePort(dataToBridgePort(d))
	if err != nil {
		return diag.FromErr(err)
	}
	recordBridgePortToData(bridgePort, d)

	return resourceBridgePortRead(ctx, d, m)
}

func resourceBridgePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	bridgePort, err := c.FindBridgePort(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return recordBridgePortToData(bridgePort, d)
}

func resourceBridgePortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	bridgePort := dataToBridgePort(d)
	_, err := c.UpdateBridgePort(bridgePort)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceBridgePortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Mikrotik)
	err := c.DeleteBridgePort(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToBridgePort(d *schema.ResourceData) *client.BridgePort {
	return &client.BridgePort{
		Id:         d.Id(),
		Bridge:     d.Get("bridge").(string),
		Interface:  d.Get("interface").(string),
		PVId:       d.Get("pvid").(int),
		Comment:    d.Get("comment").(string),
		FrameTypes: d.Get("frame_types").(string),
	}
}

func recordBridgePortToData(r *client.BridgePort, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("bridge", r.Bridge); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("interface", r.Interface); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("pvid", r.PVId); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("comment", r.Comment); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("frame_types", r.FrameTypes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	d.SetId(r.Id)

	return diags
}
