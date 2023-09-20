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

type vlanInterface struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vlanInterface{}
	_ resource.ResourceWithConfigure   = &vlanInterface{}
	_ resource.ResourceWithImportState = &vlanInterface{}
)

// NewVlanInterfaceResource is a helper function to simplify the provider implementation.
func NewVlanInterfaceResource() resource.Resource {
	return &vlanInterface{}
}

func (r *vlanInterface) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *vlanInterface) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vlan_interface"
}

// Schema defines the schema for the resource.
func (s *vlanInterface) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Virtual Local Area Network (VLAN) interfaces.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of the resource.",
			},
			"interface": defaultaware.StringAttribute(
				schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("*0"),
					Description: "Name of physical interface on top of which VLAN will work.",
				},
			),
			"mtu": defaultaware.Int64Attribute(
				schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(1500),
					Description: "Layer3 Maximum transmission unit.",
				},
			),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Interface name.",
			},
			"disabled": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether to create the interface in disabled state.",
				},
			),
			"use_service_tag": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "802.1ad compatible Service Tag.",
				},
			),
			"vlan_id": defaultaware.Int64Attribute(
				schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(1),
					Description: "Virtual LAN identifier or tag that is used to distinguish VLANs. Must be equal for all computers that belong to the same VLAN.",
				},
			),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vlanInterface) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel vlanInterfaceModel
	var mikrotikModel client.VlanInterface
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *vlanInterface) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel vlanInterfaceModel
	var mikrotikModel client.VlanInterface
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vlanInterface) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel vlanInterfaceModel
	var mikrotikModel client.VlanInterface
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vlanInterface) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel vlanInterfaceModel
	var mikrotikModel client.VlanInterface
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *vlanInterface) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type vlanInterfaceModel struct {
	Id            tftypes.String `tfsdk:"id"`
	Interface     tftypes.String `tfsdk:"interface"`
	Mtu           tftypes.Int64  `tfsdk:"mtu"`
	Name          tftypes.String `tfsdk:"name"`
	Disabled      tftypes.Bool   `tfsdk:"disabled"`
	UseServiceTag tftypes.Bool   `tfsdk:"use_service_tag"`
	VlanId        tftypes.Int64  `tfsdk:"vlan_id"`
}
