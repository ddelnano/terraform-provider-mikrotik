package mikrotik

import (
	"context"
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type interfaceVeth struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceVeth{}
	_ resource.ResourceWithConfigure   = &interfaceVeth{}
	_ resource.ResourceWithImportState = &interfaceVeth{}
)

// NewInterfaceVethResource is a helper function to simplify the provider implementation.
func NewInterfaceVethResource() resource.Resource {
	return &interfaceVeth{}

}

func (i *interfaceVeth) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	i.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (i *interfaceVeth) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_veth"
}

// Schema defines the schema for the resource.
func (i *interfaceVeth) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik interface veth only supported by RouterOS v7+.",
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
				Description: "Name of the interface veth.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Comment associated with interface veth.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Boolean for whether or not the interface veth is disabled.",
			},
			"address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:	 stringdefault.StaticString("0.0.0.0/0"),
				Description: "Address assigned to interface veth.",
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Ipv4 Gateway address for interface veth.",
			},
			"gateway6": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Ipv6 Gateway address for interface veth.",
			},
			"running": schema.BoolAttribute{
				Optional:    false,
				Computed:    true,
				Description: "Whether the interface is running.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (i *interfaceVeth) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan interfaceVethModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := i.client.AddInterfaceVeth(modelToInterfaceVeth(&plan))
	if err != nil {
		resp.Diagnostics.AddError("creation failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceVethToModel(created, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (i *interfaceVeth) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state interfaceVethModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := i.client.FindInterfaceVeth(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remote resource",
			fmt.Sprintf("Could not read interfaceVeth with name %q", state.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(interfaceVethToModel(resource, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (i *interfaceVeth) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan interfaceVethModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := i.client.UpdateInterfaceVeth(modelToInterfaceVeth(&plan))
	if err != nil {
		resp.Diagnostics.AddError("update failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceVethToModel(updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (i *interfaceVeth) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state interfaceVethModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := i.client.DeleteInterfaceVeth(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Could not delete interfaceVeth", err.Error())
		return
	}
}

func (i *interfaceVeth) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type interfaceVethModel struct {
	ID         tftypes.String `tfsdk:"id"`
	Name       tftypes.String `tfsdk:"name"`
	Comment    tftypes.String `tfsdk:"comment"`
	Disabled   tftypes.Bool   `tfsdk:"disabled"`
	Address	   tftypes.String `tfsdk:"address"`
	Gateway	   tftypes.String `tfsdk:"gateway"`
	Gateway6   tftypes.String `tfsdk:"gateway6"`
	Running    tftypes.Bool   `tfsdk:"running"`
}

func interfaceVethToModel(i *client.InterfaceVeth, m *interfaceVethModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if i == nil {
		diags.AddError("Interface Veth cannot be nil", "Cannot build model from nil object")
		return diags
	}
	m.ID = tftypes.StringValue(i.Id)
	m.Name = tftypes.StringValue(i.Name)
	m.Comment = tftypes.StringValue(i.Comment)
	m.Disabled = tftypes.BoolValue(i.Disabled)
	m.Address = tftypes.StringValue(i.Address)
	m.Gateway = tftypes.StringValue(i.Gateway)
	m.Gateway6 = tftypes.StringValue(i.Gateway6)
	m.Running = tftypes.BoolValue(i.Running)

	return diags
}

func modelToInterfaceVeth(m *interfaceVethModel) *client.InterfaceVeth {
	return &client.InterfaceVeth{
		Id:         m.ID.ValueString(),
		Name:       m.Name.ValueString(),
		Comment:    m.Comment.ValueString(),
		Disabled:   m.Disabled.ValueBool(),
		Address:	m.Address.ValueString(),
		Gateway:	m.Gateway.ValueString(),
		Gateway6:	m.Gateway6.ValueString(),
		Running:    m.Running.ValueBool(),
	}
}
