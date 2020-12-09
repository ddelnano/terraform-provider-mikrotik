package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBgpInstanceCreate,
		Read:   resourceBgpInstanceRead,
		Update: resourceBgpInstanceUpdate,
		Delete: resourceBgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"as": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"client_to_client_reflection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"confederation_peers": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ignore_as_path_len": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"out_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"redistribute_connected": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_ospf": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_other_bgp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_rip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redistribute_static": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_table": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"confederation": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceBgpInstanceCreate(d *schema.ResourceData, m interface{}) error {
	instance := prepareBgpInstance(d)

	c := m.(client.Mikrotik)

	bgpInstance, err := c.AddBgpInstance(instance)
	if err != nil {
		return err
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	bgpInstance, err := c.FindBgpInstance(d.Id())

	// TODO: Ignoring this error can cause all resources to think they
	// need to be created. We should more appropriately handle this. The
	// error where the record is not found is not actually an error and
	// needs to be disambiguated from real failures
	if err != nil {
		d.SetId("")
		return nil
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	currentBgpInstance, err := c.FindBgpInstance(d.Get("name").(string))

	instance := prepareBgpInstance(d)
	instance.ID = currentBgpInstance.ID

	bgpInstance, err := c.UpdateBgpInstance(instance)

	if err != nil {
		return err
	}

	return bgpInstanceToData(bgpInstance, d)
}

func resourceBgpInstanceDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	err := c.DeleteBgpInstance(d.Get("name").(string))

	if err != nil {
		return err
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
