package mikrotik

import (
	"context"
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type scheduler struct {
	client *client.Mikrotik
}

func schedulerToModel(s *client.Scheduler, m *schedulerModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if s == nil {
		diags.AddError("Scheduler cannot be nil", "Cannot build model from nil object")
		return diags
	}

	m.ID = tftypes.StringValue(s.Id)
	m.Interval = tftypes.Int64Value(int64(s.Interval))
	m.Name = tftypes.StringValue(s.Name)
	m.OnEvent = tftypes.StringValue(s.OnEvent)
	m.StartDate = tftypes.StringValue(s.StartDate)
	m.StartTime = tftypes.StringValue(s.StartTime)

	return diags
}

func modelToScheduler(m *schedulerModel) *client.Scheduler {
	return &client.Scheduler{
		Id:        m.ID.ValueString(),
		Name:      m.Name.ValueString(),
		OnEvent:   m.OnEvent.ValueString(),
		StartDate: m.StartDate.ValueString(),
		StartTime: m.StartTime.ValueString(),
		Interval:  types.MikrotikDuration(m.Interval.ValueInt64()),
	}
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
		resp.Diagnostics.AddError(
			"Missing ProviderData",
			"Provider data (API client) is not properly configured for this resource.",
		)
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the task.",
			},
			"on_event": schema.StringAttribute{
				Required:    true,
				Description: "Name of the script to execute. It must exist `/system script`.",
			},
			"start_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date of the first script execution.",
			},
			"start_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time of the first script execution.",
			},
			"interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between two script executions, if time interval is set to zero, the script is only executed at its start time, otherwise it is executed repeatedly at the time interval is specified.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (s *scheduler) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schedulerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := s.client.AddScheduler(modelToScheduler(&plan))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("creation failed"), err.Error())
		return
	}

	plan.ID = tftypes.StringValue(created.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *scheduler) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schedulerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := s.client.FindScheduler(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remote resource",
			fmt.Sprintf("Could not read scheduler with id %q", state.ID.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(schedulerToModel(resource, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (s *scheduler) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schedulerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := s.client.UpdateScheduler(modelToScheduler(&plan))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("update failed"), err.Error())
		return
	}

	resp.Diagnostics.Append(schedulerToModel(updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *scheduler) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schedulerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := s.client.DeleteScheduler(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Could not delete scheduler", err.Error())
		return
	}
}

func (s *scheduler) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type schedulerModel struct {
	ID        tftypes.String `tfsdk:"id"`
	Name      tftypes.String `tfsdk:"name"`
	OnEvent   tftypes.String `tfsdk:"on_event"`
	StartDate tftypes.String `tfsdk:"start_date"`
	StartTime tftypes.String `tfsdk:"start_time"`
	Interval  tftypes.Int64  `tfsdk:"interval"`
}
