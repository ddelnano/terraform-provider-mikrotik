package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallFilterRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manages firewall filter rules.",

		CreateContext: resourceFirewallFilterRuleCreate,
		ReadContext:   resourceFirewallFilterRuleRead,
		UpdateContext: resourceFirewallFilterRuleUpdate,
		DeleteContext: resourceFirewallFilterRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: firewallFilterRuleSchema(),
	}
}

func firewallFilterRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"action": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "accept",
			Description: "Action to take if packet is matched by the rule.",
		},
		"chain": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies to which chain rule will be added. If the input does not match the name of an already defined chain, a new chain will be created.",
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comment to the rule.",
		},
		"connection_state": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Interprets the connection tracking analysis data for a particular packet.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dst_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "List of destination port numbers or port number ranges.",
		},
		"in_interface": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Interface the packet has entered the router.",
		},
		"in_interface_list": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set of interfaces defined in interface list. Works the same as in-interface.",
		},
		"out_interface_list": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set of interfaces defined in interface list. Works the same as out-interface.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "tcp",
			Description: "Matches particular IP protocol specified by protocol name or number.",
		},
	}
}

func resourceFirewallFilterRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	created, err := c.AddFirewallFilterRule(dataToFirewallFilterRule(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.Id)

	return resourceFirewallFilterRuleRead(ctx, d, m)
}

func resourceFirewallFilterRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	found, err := c.FindFirewallFilterRule(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return recordFirewallFilterRuleToData(found, d)
}

func resourceFirewallFilterRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	_, err := c.UpdateFirewallFilterRule(dataToFirewallFilterRule(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallFilterRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	if err := c.DeleteFirewallFilterRule(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataToFirewallFilterRule(d *schema.ResourceData) *client.FirewallFilterRule {
	getStringList := func(data interface{}) []string {
		ret := []string{}
		if data == nil {
			return ret
		}

		for _, v := range data.([]interface{}) {
			if s, ok := v.(string); ok {
				ret = append(ret, s)
			}
		}

		return ret
	}

	return &client.FirewallFilterRule{
		Id:               d.Id(),
		Action:           d.Get("action").(string),
		Chain:            d.Get("chain").(string),
		Comment:          d.Get("comment").(string),
		ConnectionState:  getStringList(d.Get("connection_state")),
		DestPort:         d.Get("dst_port").(string),
		InInterface:      d.Get("in_interface").(string),
		InInterfaceList:  d.Get("in_interface_list").(string),
		OutInterfaceList: d.Get("out_interface_list").(string),
		Protocol:         d.Get("protocol").(string),
	}
}

func recordFirewallFilterRuleToData(r *client.FirewallFilterRule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("action", r.Action); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("chain", r.Chain); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("comment", r.Comment); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("connection_state", r.ConnectionState); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("dst_port", r.DestPort); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("in_interface", r.InInterface); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("in_interface_list", r.InInterfaceList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("out_interface_list", r.OutInterfaceList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("protocol", r.Protocol); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.SetId(r.Id)

	return diags
}
