package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/types/defaultaware"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpInstance struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bgpInstance{}
	_ resource.ResourceWithConfigure   = &bgpInstance{}
	_ resource.ResourceWithImportState = &bgpInstance{}
)

// NewBgpInstanceResource is a helper function to simplify the provider implementation.
func NewBgpInstanceResource() resource.Resource {
	return &bgpInstance{}
}

func (r *bgpInstance) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *bgpInstance) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_instance"
}

// Schema defines the schema for the resource.
func (s *bgpInstance) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik BGP Instance.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the BGP instance.",
			},
			"as": schema.Int64Attribute{
				Required:    true,
				Description: "The 32-bit BGP autonomous system number. Must be a value within 0 to 4294967295.",
			},
			"client_to_client_reflection": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
					Description: "In case this instance is a route reflector: whether to redistribute routes learned from one routing reflection client to other clients.",
				},
			),
			"comment": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "The comment of the BGP instance to be created.",
				},
			),
			"confederation_peers": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "List of AS numbers internal to the [local] confederation. For example: `10,20,30-50`.",
				},
			),
			"disabled": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether instance is disabled.",
				},
			),
			"ignore_as_path_len": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether to ignore AS_PATH attribute in BGP route selection algorithm.",
				},
			),
			"out_filter": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "Output routing filter chain used by all BGP peers belonging to this instance.",
				},
			),
			"redistribute_connected": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If enabled, this BGP instance will redistribute the information about connected routes.",
				},
			),
			"redistribute_ospf": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If enabled, this BGP instance will redistribute the information about routes learned by OSPF.",
				},
			),
			"redistribute_other_bgp": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If enabled, this BGP instance will redistribute the information about routes learned by other BGP instances.",
				},
			),
			"redistribute_rip": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If enabled, this BGP instance will redistribute the information about routes learned by RIP.",
				},
			),
			"redistribute_static": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "If enabled, the router will redistribute the information about static routes added to its routing database.",
				},
			),
			"router_id": schema.StringAttribute{
				Required:    true,
				Description: "BGP Router ID (for this instance). If set to 0.0.0.0, BGP will use one of router's IP addresses.",
			},
			"routing_table": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "Name of routing table this BGP instance operates on. ",
				},
			),
			"cluster_id": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "In case this instance is a route reflector: cluster ID of the router reflector cluster this instance belongs to.",
				},
			),
			"confederation": defaultaware.Int64Attribute(
				schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
					Description: "In case of BGP confederations: autonomous system number that identifies the [local] confederation as a whole.",
				},
			),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bgpInstance) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bgpInstanceModel
	var mikrotikModel client.BgpInstance
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bgpInstance) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bgpInstanceModel
	var mikrotikModel client.BgpInstance
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bgpInstance) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bgpInstanceModel
	var mikrotikModel client.BgpInstance
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bgpInstance) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bgpInstanceModel
	var mikrotikModel client.BgpInstance
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bgpInstance) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bgpInstanceModel struct {
	Id                       tftypes.String `tfsdk:"id"`
	Name                     tftypes.String `tfsdk:"name"`
	As                       tftypes.Int64  `tfsdk:"as"`
	ClientToClientReflection tftypes.Bool   `tfsdk:"client_to_client_reflection"`
	Comment                  tftypes.String `tfsdk:"comment"`
	ConfederationPeers       tftypes.String `tfsdk:"confederation_peers"`
	Disabled                 tftypes.Bool   `tfsdk:"disabled"`
	IgnoreAsPathLen          tftypes.Bool   `tfsdk:"ignore_as_path_len"`
	OutFilter                tftypes.String `tfsdk:"out_filter"`
	RedistributeConnected    tftypes.Bool   `tfsdk:"redistribute_connected"`
	RedistributeOspf         tftypes.Bool   `tfsdk:"redistribute_ospf"`
	RedistributeOtherBgp     tftypes.Bool   `tfsdk:"redistribute_other_bgp"`
	RedistributeRip          tftypes.Bool   `tfsdk:"redistribute_rip"`
	RedistributeStatic       tftypes.Bool   `tfsdk:"redistribute_static"`
	RouterID                 tftypes.String `tfsdk:"router_id"`
	RoutingTable             tftypes.String `tfsdk:"routing_table"`
	ClusterID                tftypes.String `tfsdk:"cluster_id"`
	Confederation            tftypes.Int64  `tfsdk:"confederation"`
}
