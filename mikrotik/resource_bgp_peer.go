package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/types/defaultaware"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpPeer struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bgpPeer{}
	_ resource.ResourceWithConfigure   = &bgpPeer{}
	_ resource.ResourceWithImportState = &bgpPeer{}
)

// NewBgpPeerResource is a helper function to simplify the provider implementation.
func NewBgpPeerResource() resource.Resource {
	return &bgpPeer{}
}

func (r *bgpPeer) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *bgpPeer) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_peer"
}

// Schema defines the schema for the resource.
func (s *bgpPeer) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik BGP Peer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique MikroTik identifier.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the BGP peer.",
			},
			"address_families": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Default:     stringdefault.StaticString("ip"),
					Description: "The list of address families about which this peer will exchange routing information.",
				},
			),
			"allow_as_in": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "How many times to allow own AS number in AS-PATH, before discarding a prefix.",
			},
			"as_override": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If set, then all instances of remote peer's AS number in BGP AS PATH attribute are replaced with local AS number before sending route update to that peer.",
				},
			),
			"cisco_vpls_nlri_len_fmt": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "VPLS NLRI length format type.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The comment of the BGP peer to be created.",
			},
			"default_originate": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("never"),
					Description: "The comment of the BGP peer to be created.",
				},
			),
			"disabled": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether peer is disabled.",
				},
			),
			"hold_time": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("3m"),
					Description: "Specifies the BGP Hold Time value to use when negotiating with peer",
				},
			),
			"in_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the routing filter chain that is applied to the incoming routing information.",
			},
			"instance": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("default"),
					Description: "The name of the instance this peer belongs to. See Mikrotik bgp instance resource.",
				},
			),
			"keepalive_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"max_prefix_limit": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of prefixes to accept from a specific peer.",
			},
			"max_prefix_restart_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum time interval after which peers can reestablish BGP session.",
			},
			"multihop": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the remote peer is more than one hop away.",
			},
			"nexthop_choice": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("default"),
					Description: "Affects the outgoing NEXT_HOP attribute selection, either: 'default', 'force-self', or 'propagate'",
				},
			),
			"out_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the routing filter chain that is applied to the outgoing routing information. ",
			},
			"passive": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Name of the routing filter chain that is applied to the outgoing routing information.",
				},
			),
			"remote_address": schema.StringAttribute{
				Required:    true,
				Description: "The address of the remote peer",
			},
			"remote_as": schema.Int64Attribute{
				Required:    true,
				Description: "The 32-bit AS number of the remote peer.",
			},
			"remote_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Remote peers port to establish tcp session.",
			},
			"remove_private_as": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If set, then BGP AS-PATH attribute is removed before sending out route update if attribute contains only private AS numbers.",
			},
			"route_reflect": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether this peer is route reflection client.",
			},
			"tcp_md5_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Key used to authenticate the connection with TCP MD5 signature as described in RFC 2385.",
			},
			"ttl": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("default"),
					Description: "Time To Live, the hop limit for TCP connection. This is a `string` field that can be 'default' or '0'-'255'.",
				},
			),
			"update_source": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If address is specified, this address is used as the source address of the outgoing TCP connection.",
			},
			"use_bfd": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to use BFD protocol for fast state detection.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bgpPeer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bgpPeerModel
	var mikrotikModel client.BgpPeer
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bgpPeer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bgpPeerModel
	var mikrotikModel client.BgpPeer
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bgpPeer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bgpPeerModel
	var mikrotikModel client.BgpPeer
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bgpPeer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bgpPeerModel
	var mikrotikModel client.BgpPeer
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bgpPeer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bgpPeerModel struct {
	Id                   tftypes.String `tfsdk:"id"`
	Name                 tftypes.String `tfsdk:"name"`
	AddressFamilies      tftypes.String `tfsdk:"address_families"`
	AllowAsIn            tftypes.Int64  `tfsdk:"allow_as_in"`
	AsOverride           tftypes.Bool   `tfsdk:"as_override"`
	CiscoVplsNlriLenFmt  tftypes.String `tfsdk:"cisco_vpls_nlri_len_fmt"`
	Comment              tftypes.String `tfsdk:"comment"`
	DefaultOriginate     tftypes.String `tfsdk:"default_originate"`
	Disabled             tftypes.Bool   `tfsdk:"disabled"`
	HoldTime             tftypes.String `tfsdk:"hold_time"`
	InFilter             tftypes.String `tfsdk:"in_filter"`
	Instance             tftypes.String `tfsdk:"instance"`
	KeepAliveTime        tftypes.String `tfsdk:"keepalive_time"`
	MaxPrefixLimit       tftypes.Int64  `tfsdk:"max_prefix_limit"`
	MaxPrefixRestartTime tftypes.String `tfsdk:"max_prefix_restart_time"`
	Multihop             tftypes.Bool   `tfsdk:"multihop"`
	NexthopChoice        tftypes.String `tfsdk:"nexthop_choice"`
	OutFilter            tftypes.String `tfsdk:"out_filter"`
	Passive              tftypes.Bool   `tfsdk:"passive"`
	RemoteAddress        tftypes.String `tfsdk:"remote_address"`
	RemoteAs             tftypes.Int64  `tfsdk:"remote_as"`
	RemotePort           tftypes.Int64  `tfsdk:"remote_port"`
	RemovePrivateAs      tftypes.Bool   `tfsdk:"remove_private_as"`
	RouteReflect         tftypes.Bool   `tfsdk:"route_reflect"`
	TCPMd5Key            tftypes.String `tfsdk:"tcp_md5_key"`
	TTL                  tftypes.String `tfsdk:"ttl"`
	UpdateSource         tftypes.String `tfsdk:"update_source"`
	UseBfd               tftypes.Bool   `tfsdk:"use_bfd"`
}
