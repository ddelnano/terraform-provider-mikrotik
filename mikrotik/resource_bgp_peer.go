package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBgpPeer() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a MikroTik BGP Peer.",

		CreateContext: resourceBgpPeerCreate,
		ReadContext:   resourceBgpPeerRead,
		UpdateContext: resourceBgpPeerUpdate,
		DeleteContext: resourceBgpPeerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the BGP peer.",
			},
			"remote_as": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The 32-bit AS number of the remote peer.",
			},
			"remote_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The address of the remote peer",
			},
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the instance this peer belongs to. See Mikrotik bgp instance resource.",
			},
			"address_families": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ip",
				Description: "The list of address families about which this peer will exchange routing information.",
			},
			"ttl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Time To Live, the hop limit for TCP connection. This is a `string` field that can be 'default' or '0'-'255'.",
			},
			"default_originate": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "never",
				Description: "The comment of the BGP peer to be created.",
			},
			"hold_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "3m",
				Description: "Specifies the BGP Hold Time value to use when negotiating with peer",
			},
			"nexthop_choice": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Affects the outgoing NEXT_HOP attribute selection, either: 'default', 'force-self', or 'propagate'",
			},
			"out_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the routing filter chain that is applied to the outgoing routing information. ",
			},
			"in_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the routing filter chain that is applied to the incoming routing information.",
			},
			"allow_as_in": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How many times to allow own AS number in AS-PATH, before discarding a prefix.",
			},
			"as_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set, then all instances of remote peer's AS number in BGP AS PATH attribute are replaced with local AS number before sending route update to that peer.",
			},
			"cisco_vpls_nlri_len_fmt": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPLS NLRI length format type.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment of the BGP peer to be created.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether peer is disabled.",
			},
			"keepalive_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_prefix_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of prefixes to accept from a specific peer.",
			},
			"max_prefix_restart_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Minimum time interval after which peers can reestablish BGP session.",
			},
			"multihop": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether the remote peer is more than one hop away.",
			},
			"passive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Name of the routing filter chain that is applied to the outgoing routing information.",
			},
			"remote_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Remote peers port to establish tcp session.",
			},
			"remove_private_as": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set, then BGP AS-PATH attribute is removed before sending out route update if attribute contains only private AS numbers.",
			},
			"route_reflect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether this peer is route reflection client.",
			},
			"tcp_md5_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key used to authenticate the connection with TCP MD5 signature as described in RFC 2385.",
			},
			"update_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If address is specified, this address is used as the source address of the outgoing TCP connection.",
			},
			"use_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to use BFD protocol for fast state detection.",
			},
		},
	}
}

func resourceBgpPeerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	peer := prepareBgpPeer(d)

	c := m.(*client.Mikrotik)

	bgpPeer, err := c.AddBgpPeer(peer)
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	bgpPeer, err := c.FindBgpPeer(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	currentBgpPeer, err := c.FindBgpPeer(d.Get("name").(string))

	peer := prepareBgpPeer(d)
	peer.ID = currentBgpPeer.ID

	bgpPeer, err := c.UpdateBgpPeer(peer)
	if err != nil {
		return diag.FromErr(err)
	}

	return bgpPeerToData(bgpPeer, d)
}

func resourceBgpPeerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteBgpPeer(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func bgpPeerToData(b *client.BgpPeer, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"name":                    b.Name,
		"address_families":        b.AddressFamilies,
		"allow_as_in":             b.AllowAsIn,
		"as_override":             b.AsOverride,
		"cisco_vpls_nlri_len_fmt": b.CiscoVplsNlriLenFmt,
		"comment":                 b.Comment,
		"default_originate":       b.DefaultOriginate,
		"disabled":                b.Disabled,
		"hold_time":               b.HoldTime,
		"in_filter":               b.InFilter,
		"instance":                b.Instance,
		"keepalive_time":          b.KeepAliveTime,
		"max_prefix_limit":        b.MaxPrefixLimit,
		"max_prefix_restart_time": b.MaxPrefixRestartTime,
		"multihop":                b.Multihop,
		"nexthop_choice":          b.NexthopChoice,
		"out_filter":              b.OutFilter,
		"passive":                 b.Passive,
		"remote_address":          b.RemoteAddress,
		"remote_as":               b.RemoteAs,
		"remote_port":             b.RemotePort,
		"remove_private_as":       b.RemovePrivateAs,
		"route_reflect":           b.RouteReflect,
		"tcp_md5_key":             b.TCPMd5Key,
		"ttl":                     b.TTL,
		"update_source":           b.UpdateSource,
		"use_bfd":                 b.UseBfd,
	}

	d.SetId(b.Name)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
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
