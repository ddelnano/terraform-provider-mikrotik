package mikrotik

import (
	"context"
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/types/defaultaware"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type interfaceWireguard struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceWireguard{}
	_ resource.ResourceWithConfigure   = &interfaceWireguard{}
	_ resource.ResourceWithImportState = &interfaceWireguard{}
)

// NewInterfaceWireguardResource is a helper function to simplify the provider implementation.
func NewInterfaceWireguardResource() resource.Resource {
	return &interfaceWireguard{}

}

func (i *interfaceWireguard) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	i.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (i *interfaceWireguard) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_wireguard"
}

// Schema defines the schema for the resource.
func (i *interfaceWireguard) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik interface wireguard only supported by RouterOS v7+.",
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
				Description: "Name of the interface wireguard.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Comment associated with interface wireguard.",
			},
			"disabled": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Boolean for whether or not the interface wireguard is disabled.",
				},
			),
			"listen_port": defaultaware.Int64Attribute(
				schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(13231),
					Description: "Port for WireGuard service to listen on for incoming sessions.",
				},
			),
			"mtu": defaultaware.Int64Attribute(
				schema.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(1420),
					Validators: []validator.Int64{
						int64validator.Between(0, 65536),
					},
					Description: "Layer3 Maximum transmission unit.",
				},
			),
			"private_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "A base64 private key. If not specified, it will be automatically generated upon interface creation.",
			},
			"public_key": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "A base64 public key is calculated from the private key.",
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
func (i *interfaceWireguard) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan interfaceWireguardModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := i.client.AddInterfaceWireguard(modelToInterfaceWireguard(&plan))
	if err != nil {
		resp.Diagnostics.AddError("creation failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceWireguardToModel(created, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (i *interfaceWireguard) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state interfaceWireguardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := i.client.FindInterfaceWireguard(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remote resource",
			fmt.Sprintf("Could not read interfaceWireguard with name %q", state.Name.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(interfaceWireguardToModel(resource, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (i *interfaceWireguard) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan interfaceWireguardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := i.client.UpdateInterfaceWireguard(modelToInterfaceWireguard(&plan))
	if err != nil {
		resp.Diagnostics.AddError("update failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceWireguardToModel(updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (i *interfaceWireguard) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state interfaceWireguardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := i.client.DeleteInterfaceWireguard(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Could not delete interfaceWireguard", err.Error())
		return
	}
}

func (i *interfaceWireguard) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type interfaceWireguardModel struct {
	ID         tftypes.String `tfsdk:"id"`
	Name       tftypes.String `tfsdk:"name"`
	Comment    tftypes.String `tfsdk:"comment"`
	Disabled   tftypes.Bool   `tfsdk:"disabled"`
	ListenPort tftypes.Int64  `tfsdk:"listen_port"`
	Mtu        tftypes.Int64  `tfsdk:"mtu"`
	PrivateKey tftypes.String `tfsdk:"private_key"`
	PublicKey  tftypes.String `tfsdk:"public_key"`
	Running    tftypes.Bool   `tfsdk:"running"`
}

func interfaceWireguardToModel(i *client.InterfaceWireguard, m *interfaceWireguardModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if i == nil {
		diags.AddError("Interface Wireguard cannot be nil", "Cannot build model from nil object")
		return diags
	}
	m.ID = tftypes.StringValue(i.Id)
	m.Name = tftypes.StringValue(i.Name)
	m.Comment = tftypes.StringValue(i.Comment)
	m.Disabled = tftypes.BoolValue(i.Disabled)
	m.ListenPort = tftypes.Int64Value(int64(i.ListenPort))
	m.Mtu = tftypes.Int64Value(int64(i.Mtu))
	m.PrivateKey = tftypes.StringValue(i.PrivateKey)
	m.PublicKey = tftypes.StringValue(i.PublicKey)
	m.Running = tftypes.BoolValue(i.Running)

	return diags
}

func modelToInterfaceWireguard(m *interfaceWireguardModel) *client.InterfaceWireguard {
	return &client.InterfaceWireguard{
		Id:         m.ID.ValueString(),
		Name:       m.Name.ValueString(),
		Comment:    m.Comment.ValueString(),
		Disabled:   m.Disabled.ValueBool(),
		ListenPort: int(m.ListenPort.ValueInt64()),
		Mtu:        int(m.Mtu.ValueInt64()),
		PrivateKey: m.PrivateKey.ValueString(),
		PublicKey:  m.PublicKey.ValueString(),
		Running:    m.Running.ValueBool(),
	}
}
