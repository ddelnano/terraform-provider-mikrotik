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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bridge struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bridge{}
	_ resource.ResourceWithConfigure   = &bridge{}
	_ resource.ResourceWithImportState = &bridge{}
)

// NewBridgeResource is a helper function to simplify the provider implementation.
func NewBridgeResource() resource.Resource {
	return &bridge{}
}

func (r *bridge) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *bridge) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bridge"
}

// Schema defines the schema for the resource.
func (s *bridge) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a bridge resource on remote MikroTik device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID for the instance.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the bridge interface",
			},
			"fast_forward": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
					Description: "Special and faster case of FastPath which works only on bridges with 2 interfaces (enabled by default only for new bridges).",
				},
			),
			"vlan_filtering": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Globally enables or disables VLAN functionality for bridge.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Short description of the interface.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bridge) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bridgeModel
	var mikrotikModel client.Bridge
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bridge) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bridgeModel
	var mikrotikModel client.Bridge
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bridge) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bridgeModel
	var mikrotikModel client.Bridge
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bridge) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bridgeModel
	var mikrotikModel client.Bridge
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bridge) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bridgeModel struct {
	Id            tftypes.String `tfsdk:"id"`
	Name          tftypes.String `tfsdk:"name"`
	FastForward   tftypes.Bool   `tfsdk:"fast_forward"`
	VlanFiltering tftypes.Bool   `tfsdk:"vlan_filtering"`
	Comment       tftypes.String `tfsdk:"comment"`
}
