package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceScriptCreate,
		Read:   resourceScriptRead,
		Update: resourceScriptUpdate,
		Delete: resourceScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dont_require_permissions": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceScriptCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	source := d.Get("source").(string)
	policy := d.Get("policy").([]interface{})
	policies := []string{}
	for _, p := range policy {
		policies = append(policies, p.(string))
	}
	dontReqPerms := d.Get("dont_require_permissions").(bool)

	c := m.(client.Mikrotik)

	script, err := c.CreateScript(
		name,
		owner,
		source,
		policies,
		dontReqPerms,
	)
	if err != nil {
		return err
	}

	scriptToData(script, d)
	return nil
}

func scriptToData(s client.Script, d *schema.ResourceData) error {
	d.SetId(s.Name)
	d.Set("name", s.Name)
	d.Set("owner", s.Owner)
	d.Set("source", s.Source)
	err := d.Set("policy", s.Policy())
	if err != nil {
		return err
	}
	d.Set("dont_require_permissions", s.DontRequirePermissions)
	return nil
}

func resourceScriptRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	script, err := c.FindScript(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}
	scriptToData(script, d)
	return nil
}
func resourceScriptUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	source := d.Get("source").(string)
	dontReqPerms := d.Get("dont_require_permissions").(bool)
	policy := d.Get("policy").([]interface{})
	policies := []string{}
	for _, p := range policy {
		str, ok := p.(string)
		if ok {
			policies = append(policies, str)
		}
	}

	c := m.(client.Mikrotik)

	script, err := c.UpdateScript(name, owner, source, policies, dontReqPerms)
	if err != nil {
		return err
	}

	scriptToData(script, d)

	return nil
}
func resourceScriptDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	c := m.(client.Mikrotik)

	err := c.DeleteScript(name)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
