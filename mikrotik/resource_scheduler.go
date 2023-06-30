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

type scheduler struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &scheduler{}
	_ resource.ResourceWithConfigure   = &scheduler{}
	_ resource.ResourceWithImportState = &scheduler{}
)

// NewSchedulerResource is a helper function to simplify the provider implementation.
func NewSchedulerResource() resource.Resource {
	return &scheduler{}
}

func (s *scheduler) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	s.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (s *scheduler) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduler"
}

// Schema defines the schema for the resource.
func (s *scheduler) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik scheduler.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Identifier of this resource assigned by RouterOS",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the task.",
			},
			"on_event": schema.StringAttribute{
				Required:    true,
				Description: "Name of the script to execute. It must exist `/system script`.",
			},
			"start_date": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Date of the first script execution.",
			},
			"start_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Time of the first script execution.",
			},
			"interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Interval between two script executions, if time interval is set to zero, the script is only executed at its start time, otherwise it is executed repeatedly at the time interval is specified.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (s *scheduler) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel schedulerModel
	var mikrotikModel client.Scheduler
	GenericCreateResource(&terraformModel, &mikrotikModel, s.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (s *scheduler) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel schedulerModel
	var mikrotikModel client.Scheduler
	GenericReadResource(&terraformModel, &mikrotikModel, s.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (s *scheduler) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel schedulerModel
	var mikrotikModel client.Scheduler
	GenericUpdateResource(&terraformModel, &mikrotikModel, s.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *scheduler) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel schedulerModel
	var mikrotikModel client.Scheduler
	GenericDeleteResource(&terraformModel, &mikrotikModel, s.client)(ctx, req, resp)
}

func (s *scheduler) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type schedulerModel struct {
	Id        tftypes.String `tfsdk:"id"`
	Name      tftypes.String `tfsdk:"name"`
	OnEvent   tftypes.String `tfsdk:"on_event"`
	StartDate tftypes.String `tfsdk:"start_date"`
	StartTime tftypes.String `tfsdk:"start_time"`
	Interval  tftypes.Int64  `tfsdk:"interval"`
}
