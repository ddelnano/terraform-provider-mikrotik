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
//
// The body is a copy-paste code from `GenericUpdateResource()` functions.
// It's done to support special case of 'unsetting' the 'next_pool' field.
// Since RouterOS API does not support empty value `""` for this field,
// a 'magic' value of 'none' is used.
// In that case, Terraform argues that planned value was `none` but actual (after Read() method) is `""`
// The only difference from `GenericUpdateResource()` in this implementation is checking of
// transition from some value to `""` for `next_pool` field. In that case, we simply change value to `none`,
// so API client can unset this value and subsequent `Read()` method will see `""` which is the same as config value.
//
// Be aware, that this hack prevents using `none` value explicitly in the configuration.
func (r *pool) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel, state poolModel
	var mikrotikModel client.Pool
	resp.Diagnostics.Append(req.Plan.Get(ctx, &terraformModel)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// if practitioner sets this value to `""` to unset field in remote system,
	// implicitly send `"none"` via API
	if !terraformModel.NextPool.Equal(state.NextPool) && terraformModel.NextPool.ValueString() == "" {
		terraformModel.NextPool = tftypes.StringValue("none")
	}

	if err := utils.TerraformModelToMikrotikStruct(&terraformModel, &mikrotikModel); err != nil {
		resp.Diagnostics.AddError("Cannot copy model: Terraform -> MikroTik", err.Error())
		return
	}
	updated, err := r.client.Update(&mikrotikModel)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	if err := utils.MikrotikStructToTerraformModel(updated, &terraformModel); err != nil {
		resp.Diagnostics.AddError("Cannot copy model: MikroTik -> Terraform", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, terraformModel)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pool) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel poolModel
	var mikrotikModel client.Pool
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *pool) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type poolModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Name     tftypes.String `tfsdk:"name"`
	Ranges   tftypes.String `tfsdk:"ranges"`
	NextPool tftypes.String `tfsdk:"next_pool"`
	Comment  tftypes.String `tfsdk:"comment"`
}
