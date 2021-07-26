package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBgpPeerCreate,
		Read:   resourceBgpPeerRead,
		Update: resourceBgpPeerUpdate,
		Delete: resourceBgpPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_as": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"remote_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address_families": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ip",
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"default_originate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "never",
			},
			"hold_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "3m",
			},
			"nexthop_choice": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"out_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"in_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_as_in": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"as_override": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cisco_vpls_nlri_len_fmt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keepalive_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_prefix_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_prefix_restart_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"multihop": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"passive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"remote_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"remove_private_as": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"route_reflect": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tcp_md5_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_bfd": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBgpPeerCreate(d *schema.ResourceData, m interface{}) error {
	instance := prepareBgpPeer(d)

	c := m.(client.Mikrotik)

	bgpPeer, err := c.AddBgpPeer(instance)
	if err != nil {
		return err
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	bgpPeer, err := c.FindBgpPeer(d.Id())

	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	currentBgpPeer, err := c.FindBgpPeer(d.Get("name").(string))

	instance := prepareBgpPeer(d)
	instance.ID = currentBgpPeer.ID

	bgpPeer, err := c.UpdateBgpPeer(instance)

	if err != nil {
		return err
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	err := c.DeleteBgpPeer(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func bgpPeerToData(b *client.BgpPeer, d *schema.ResourceData) error {
	d.SetId(b.Name)

	if err := d.Set("name", b.Name); err != nil {
		return err
	}
	if err := d.Set("address_families", b.AddressFamilies); err != nil {
		return err
	}
	if err := d.Set("allow_as_in", b.AllowAsIn); err != nil {
		return err
	}
	if err := d.Set("as_override", b.AsOverride); err != nil {
		return err
	}
	if err := d.Set("cisco_vpls_nlri_len_fmt", b.CiscoVplsNlriLenFmt); err != nil {
		return err
	}
	if err := d.Set("comment", b.Comment); err != nil {
		return err
	}
	if err := d.Set("default_originate", b.DefaultOriginate); err != nil {
		return err
	}
	if err := d.Set("disabled", b.Disabled); err != nil {
		return err
	}
	if err := d.Set("hold_time", b.HoldTime); err != nil {
		return err
	}
	if err := d.Set("in_filter", b.InFilter); err != nil {
		return err
	}
	if err := d.Set("instance", b.Instance); err != nil {
		return err
	}
	if err := d.Set("keepalive_time", b.KeepAliveTime); err != nil {
		return err
	}
	if err := d.Set("max_prefix_limit", b.MaxPrefixLimit); err != nil {
		return err
	}
	if err := d.Set("max_prefix_restart_time", b.MaxPrefixRestartTime); err != nil {
		return err
	}
	if err := d.Set("multihop", b.Multihop); err != nil {
		return err
	}
	if err := d.Set("nexthop_choice", b.NexthopChoice); err != nil {
		return err
	}
	if err := d.Set("out_filter", b.OutFilter); err != nil {
		return err
	}
	if err := d.Set("passive", b.Passive); err != nil {
		return err
	}
	if err := d.Set("remote_address", b.RemoteAddress); err != nil {
		return err
	}
	if err := d.Set("remote_as", b.RemoteAs); err != nil {
		return err
	}
	if err := d.Set("remote_port", b.RemotePort); err != nil {
		return err
	}
	if err := d.Set("remove_private_as", b.RemovePrivateAs); err != nil {
		return err
	}
	if err := d.Set("route_reflect", b.RouteReflect); err != nil {
		return err
	}
	if err := d.Set("tcp_md5_key", b.TCPMd5Key); err != nil {
		return err
	}
	if err := d.Set("ttl", b.TTL); err != nil {
		return err
	}
	if err := d.Set("update_source", b.UpdateSource); err != nil {
		return err
	}
	if err := d.Set("use_bfd", b.UseBfd); err != nil {
		return err
	}
	return nil
}

func prepareBgpPeer(d *schema.ResourceData) *client.BgpPeer {
	bgpPeer := new(client.BgpPeer)

	bgpPeer.Name = d.Get("name").(string)
	bgpPeer.AddressFamilies = d.Get("address_families").(string)
	bgpPeer.AllowAsIn = d.Get("allow_as_in").(int)
	bgpPeer.AsOverride = d.Get("as_override").(bool)
	bgpPeer.CiscoVplsNlriLenFmt = d.Get("cisco_vpls_nlri_len_fmt").(string)
	bgpPeer.Comment = d.Get("comment").(string)
	bgpPeer.DefaultOriginate = d.Get("default_originate").(string)
	bgpPeer.Disabled = d.Get("disabled").(bool)
	bgpPeer.HoldTime = d.Get("hold_time").(string)
	bgpPeer.InFilter = d.Get("in_filter").(string)
	bgpPeer.Instance = d.Get("instance").(string)
	bgpPeer.KeepAliveTime = d.Get("keepalive_time").(string)
	bgpPeer.MaxPrefixLimit = d.Get("max_prefix_limit").(int)
	bgpPeer.MaxPrefixRestartTime = d.Get("max_prefix_restart_time").(string)
	bgpPeer.Multihop = d.Get("multihop").(bool)
	bgpPeer.NexthopChoice = d.Get("nexthop_choice").(string)
	bgpPeer.OutFilter = d.Get("out_filter").(string)
	bgpPeer.Passive = d.Get("passive").(bool)
	bgpPeer.RemoteAddress = d.Get("remote_address").(string)
	bgpPeer.RemoteAs = d.Get("remote_as").(int)
	bgpPeer.RemotePort = d.Get("remote_port").(int)
	bgpPeer.RemovePrivateAs = d.Get("remove_private_as").(bool)
	bgpPeer.RouteReflect = d.Get("route_reflect").(bool)
	bgpPeer.TCPMd5Key = d.Get("tcp_md5_key").(string)
	bgpPeer.TTL = d.Get("ttl").(string)
	bgpPeer.UpdateSource = d.Get("update_source").(string)
	bgpPeer.UseBfd = d.Get("use_bfd").(bool)

	return bgpPeer
}
