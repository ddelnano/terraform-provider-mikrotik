package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolCreate,
		ReadContext:   resourcePoolRead,
		UpdateContext: resourcePoolUpdate,
		DeleteContext: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ranges": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := preparePool(d)

	c := m.(client.Mikrotik)

	pool, err := c.AddPool(p)
	if err != nil {
		return diag.FromErr(err)
	}

	err = poolToData(pool, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	pool, err := c.FindPool(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	err = poolToData(pool, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	p := preparePool(d)
	p.Id = d.Id()

	pool, err := c.UpdatePool(p)

	if err != nil {
		return diag.FromErr(err)
	}

	err = poolToData(pool, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	err := c.DeletePool(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func poolToData(pool *client.Pool, d *schema.ResourceData) error {
	d.SetId(pool.Id)
	if err := d.Set("name", pool.Name); err != nil {
		return err
	}
	if err := d.Set("ranges", pool.Ranges); err != nil {
		return err
	}
	if err := d.Set("comment", pool.Comment); err != nil {
		return err
	}
	return nil
}

func preparePool(d *schema.ResourceData) *client.Pool {
	pool := new(client.Pool)

	pool.Name = d.Get("name").(string)
	pool.Ranges = d.Get("ranges").(string)
	pool.Comment = d.Get("comment").(string)

	return pool
}
