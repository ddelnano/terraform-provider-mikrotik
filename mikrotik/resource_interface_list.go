package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInterfaceList() *schema.Resource {
	return &schema.Resource{
		Description: "Allows to define set of interfaces for easier interface management.",

		CreateContext: resourceInterfaceListCreate,
		ReadContext:   resourceInterfaceListRead,
		UpdateContext: resourceInterfaceListUpdate,
		DeleteContext: resourceInterfaceListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the interface list.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment to this list.",
			},
		},
	}
}

func resourceInterfaceListCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	r := dataToInterfaceList(d)
	record, err := c.AddInterfaceList(r)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(record.Name)

	return resourceInterfaceListRead(ctx, d, m)
}

func resourceInterfaceListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	record, err := c.FindInterfaceList(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return recordInterfaceListToData(record, d)
}

func resourceInterfaceListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	currentRecord, err := c.FindInterfaceList(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	r := dataToInterfaceList(d)
	r.Id = currentRecord.Id

	_, err = c.UpdateInterfaceList(r)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.Name)

	return resourceInterfaceListRead(ctx, d, m)
}

func resourceInterfaceListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)
	err := c.DeleteInterfaceList(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataToInterfaceList(d *schema.ResourceData) *client.InterfaceList {
	return &client.InterfaceList{
		Id:      d.Id(),
		Name:    d.Get("name").(string),
		Comment: d.Get("comment").(string),
	}
}

func recordInterfaceListToData(r *client.InterfaceList, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", r.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("comment", r.Comment); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.SetId(r.Name)

	return diags
}
