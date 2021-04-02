package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolCreate,
		Read:   resourcePoolRead,
		Update: resourcePoolUpdate,
		Delete: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {
	p := preparePool(d)

	c := m.(client.Mikrotik)

	pool, err := c.AddPool(p)
	if err != nil {
		return err
	}

	return poolToData(pool, d)
}

func resourcePoolRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	pool, err := c.FindPool(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	return poolToData(pool, d)
}

func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	p, err := c.FindPool(d.Id())
	if err != nil {
		return err
	}
	pool := preparePool(d)
	pool.Id = p.Id

	pool, err = c.UpdatePool(pool)

	if err != nil {
		return err
	}

	return poolToData(pool, d)
}

func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	err := c.DeletePool(d.Id())

	if err != nil {
		return err
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
