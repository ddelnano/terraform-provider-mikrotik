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

type interfaceListMember struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceListMember{}
	_ resource.ResourceWithConfigure   = &interfaceListMember{}
	_ resource.ResourceWithImportState = &interfaceListMember{}
)

// NewInterfaceListMemberResource is a helper function to simplify the provider implementation.
func NewInterfaceListMemberResource() resource.Resource {
	return &interfaceListMember{}
}

func (r *interfaceListMember) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *interfaceListMember) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_list_member"
}

// Schema defines the schema for the resource.
func (s *interfaceListMember) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to define set of interfaces for easier interface management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Name of the interface.",
			},
			"list": schema.StringAttribute{
				Required:    true,
				Description: " 	Name of the interface list",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *interfaceListMember) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel interfaceListMemberModel
	var mikrotikModel client.InterfaceListMember
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *interfaceListMember) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceListMemberModel
	var mikrotikModel client.InterfaceListMember
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *interfaceListMember) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceListMemberModel
	var mikrotikModel client.InterfaceListMember
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *interfaceListMember) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceListMemberModel
	var mikrotikModel client.InterfaceListMember
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceListMember) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type interfaceListMemberModel struct {
	Id        tftypes.String `tfsdk:"id"`
	Interface tftypes.String `tfsdk:"interface"`
	List      tftypes.String `tfsdk:"list"`
}
