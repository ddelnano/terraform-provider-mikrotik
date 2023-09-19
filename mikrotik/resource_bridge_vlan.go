package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bridgeVlan struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bridgeVlan{}
	_ resource.ResourceWithConfigure   = &bridgeVlan{}
	_ resource.ResourceWithImportState = &bridgeVlan{}
)

// NewBridgeVlanResource is a helper function to simplify the provider implementation.
func NewBridgeVlanResource() resource.Resource {
	return &bridgeVlan{}
}

func (r *bridgeVlan) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *bridgeVlan) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bridge_vlan"
}

// Schema defines the schema for the resource.
func (s *bridgeVlan) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik BridgeVlan.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "A unique ID for this resource.",
			},
			"bridge": schema.StringAttribute{
				Required:    true,
				Description: "The bridge interface which the respective VLAN entry is intended for.",
			},
			"tagged": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: tftypes.StringType,
				Description: "Interface list with a VLAN tag adding action in egress.",
			},
			"untagged": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: tftypes.StringType,
				Description: "Interface list with a VLAN tag removing action in egress. ",
			},
			"vlan_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: tftypes.Int64Type,
				Description: "The list of VLAN IDs for certain port configuration. Ranges are not supported yet.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bridgeVlan) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bridgeVlanModel
	var mikrotikModel client.BridgeVlan
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bridgeVlan) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bridgeVlanModel
	var mikrotikModel client.BridgeVlan
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bridgeVlan) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bridgeVlanModel
	var mikrotikModel client.BridgeVlan
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bridgeVlan) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bridgeVlanModel
	var mikrotikModel client.BridgeVlan
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bridgeVlan) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx,
		path.Root("id"),
		req,
		resp,
	)
}

type bridgeVlanModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Bridge   tftypes.String `tfsdk:"bridge"`
	Tagged   tftypes.Set    `tfsdk:"tagged"`
	Untagged tftypes.Set    `tfsdk:"untagged"`
	VlanIds  tftypes.Set    `tfsdk:"vlan_ids"`
}
