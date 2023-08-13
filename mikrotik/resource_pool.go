package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type pool struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pool{}
	_ resource.ResourceWithConfigure   = &pool{}
	_ resource.ResourceWithImportState = &pool{}
)

// NewPoolResource is a helper function to simplify the provider implementation.
func NewPoolResource() resource.Resource {
	return &pool{}
}

func (r *pool) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *pool) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool"
}

// Schema defines the schema for the resource.
func (s *pool) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik IP Pool.",

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
				Description: "The name of IP pool.",
			},
			"ranges": schema.StringAttribute{
				Required:    true,
				Description: "The IP range(s) of the pool. Multiple ranges can be specified, separated by commas: `172.16.0.6-172.16.0.12,172.16.0.50-172.16.0.60`.",
			},
			"next_pool": schema.StringAttribute{
				// CustomType: noneStringType{},
				// todo(maksym): handle special case of  "none"
				// which equals to an empty string
				// to supress diff on "" != "none"
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The IP pool to pick next address from if current is exhausted.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The comment of the IP Pool to be created.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *pool) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel poolModel
	var mikrotikModel client.Pool

	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *pool) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel poolModel
	var mikrotikModel client.Pool

	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pool) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel poolModel
	var mikrotikModel client.Pool

	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pool) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel poolModel
	var mikrotikModel client.Pool
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *pool) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type poolModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Name     tftypes.String `tfsdk:"name"`
	Ranges   tftypes.String `tfsdk:"ranges"`
	NextPool tftypes.String `tfsdk:"next_pool"`
	Comment  tftypes.String `tfsdk:"comment"`
}
