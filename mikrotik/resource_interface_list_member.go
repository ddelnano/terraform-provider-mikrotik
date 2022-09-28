package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInterfaceListMember() *schema.Resource {
	return &schema.Resource{
		Description: "Allows to define set of interfaces for easier interface management.",

		CreateContext: resourceInterfaceListMemberCreate,
		ReadContext:   resourceInterfaceListMemberRead,
		UpdateContext: resourceInterfaceListMemberUpdate,
		DeleteContext: resourceInterfaceListMemberDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
			},

			"list": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInterfaceListMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r := dataToInterfaceListMember(d)
	record, err := c.AddInterfaceListMember(r)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(record.Id)

	return resourceInterfaceListMemberRead(ctx, d, m)
}

func resourceInterfaceListMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	record, err := c.FindInterfaceListMember(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return recordInterfaceListMemberToData(record, d)
}

func resourceInterfaceListMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r := dataToInterfaceListMember(d)
	_, err := c.UpdateInterfaceListMember(r)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceInterfaceListMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	err := c.DeleteInterfaceListMember(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func dataToInterfaceListMember(d *schema.ResourceData) *client.InterfaceListMember {
	return &client.InterfaceListMember{
		Id:        d.Id(),
		Interface: d.Get("interface").(string),
		List:      d.Get("list").(string),
	}
}

func recordInterfaceListMemberToData(r *client.InterfaceListMember, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("interface", r.Interface); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("list", r.List); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	d.SetId(r.Id)

	return diags
}
