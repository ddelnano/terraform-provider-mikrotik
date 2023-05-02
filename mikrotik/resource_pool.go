package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a Mikrotik IP Pool.",

		CreateContext: resourcePoolCreate,
		ReadContext:   resourcePoolRead,
		UpdateContext: resourcePoolUpdate,
		DeleteContext: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: utils.ImportStateContextUppercaseWrapper(schema.ImportStatePassthroughContext),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of IP pool.",
			},
			"ranges": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP range(s) of the pool. Multiple ranges can be specified, separated by commas: `172.16.0.6-172.16.0.12,172.16.0.50-172.16.0.60`.",
			},
			"next_pool": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(i interface{}) string {
					v := i.(string)
					// handle special case for 'none' string:
					// it behaves the same as an empty string - unsets the value
					// and MikroTik API will return an empty string, but we don't wont diff on '' != 'none'
					if v == "none" {
						return ""
					}

					return v
				},
				Description: "The IP pool to pick next address from if current is exhausted.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment of the IP Pool to be created.",
			},
		},
	}
}

func resourcePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := preparePool(d)

	c := m.(*client.Mikrotik)

	pool, err := c.AddPool(p)
	if err != nil {
		return diag.FromErr(err)
	}

	return poolToData(pool, d)
}

func resourcePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	pool, err := c.FindPool(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return poolToData(pool, d)
}

func resourcePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	p := preparePool(d)
	p.Id = d.Id()

	pool, err := c.UpdatePool(p)

	if err != nil {
		return diag.FromErr(err)
	}

	return poolToData(pool, d)
}

func resourcePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeletePool(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func poolToData(pool *client.Pool, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"name":      pool.Name,
		"ranges":    pool.Ranges,
		"next_pool": pool.NextPool,
		"comment":   pool.Comment,
	}

	d.SetId(pool.Id)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func preparePool(d *schema.ResourceData) *client.Pool {
	pool := new(client.Pool)

	pool.Name = d.Get("name").(string)
	pool.NextPool = d.Get("next_pool").(string)
	pool.Ranges = d.Get("ranges").(string)
	pool.Comment = d.Get("comment").(string)

	return pool
}
