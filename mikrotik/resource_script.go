package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type script struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &script{}
	_ resource.ResourceWithConfigure   = &script{}
	_ resource.ResourceWithImportState = &script{}
)

// NewScriptResource is a helper function to simplify the provider implementation.
func NewScriptResource() resource.Resource {
	return &script{}
}

func (r *script) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *script) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the resource.
func (s *script) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik Script.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The name of script.",
			},
			"policy": schema.ListAttribute{
				Required:    true,
				ElementType: tftypes.StringType,
				Description: "What permissions the script has. This must be one of the following: ftp, reboot, read, write, policy, test, password, sniff, sensitive, romon.",
			},
			"dont_require_permissions": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If the script requires permissions or not.",
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "The source code of the script. See the [MikroTik docs](https://wiki.mikrotik.com/wiki/Manual:Scripting) for the scripting language.",
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner of the script.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *script) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scriptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	GenericCreateResource(&plan, &client.Script{}, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *script) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scriptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	GenericReadResource(&state, &client.Script{}, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *script) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scriptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	GenericUpdateResource(&plan, &client.Script{}, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *script) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scriptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	GenericDeleteResource(&state, &client.Script{}, r.client)(ctx, req, resp)
}

func (r *script) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type scriptModel struct {
	Id                     tftypes.String `tfsdk:"id"`
	Name                   tftypes.String `tfsdk:"name"`
	Owner                  tftypes.String `tfsdk:"owner"`
	Policy                 tftypes.List   `tfsdk:"policy"`
	DontRequirePermissions tftypes.Bool   `tfsdk:"dont_require_permissions"`
	Source                 tftypes.String `tfsdk:"source"`
}
