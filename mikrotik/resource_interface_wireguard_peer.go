package mikrotik

import (
	"context"
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type interfaceWireguardPeer struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceWireguardPeer{}
	_ resource.ResourceWithConfigure   = &interfaceWireguardPeer{}
	_ resource.ResourceWithImportState = &interfaceWireguardPeer{}
)

// NewInterfaceWireguardPeerResource is a helper function to simplify the provider implementation.
func NewInterfaceWireguardPeerResource() resource.Resource {
	return &interfaceWireguardPeer{}

}

func (i *interfaceWireguardPeer) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	i.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (i *interfaceWireguardPeer) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_wireguard_peer"
}

// Schema defines the schema for the resource.
func (i *interfaceWireguardPeer) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Mikrotik Interface Wireguard Peer only supported by RouterOS v7+.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Identifier of this resource assigned by RouterOS",
			},
			"allowed_address": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Description: "List of IP (v4 or v6) addresses with CIDR masks from which incoming traffic for this peer is allowed and to which outgoing traffic for this peer is directed. The catch-all 0.0.0.0/0 may be specified for matching all IPv4 addresses, and ::/0 may be specified for matching all IPv6 addresses.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Short description of the peer.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Boolean for whether or not the interface peer is disabled.",
			},
			"endpoint_address": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Description: "An endpoint IP or hostname can be left blank to allow remote connection from any address.",
			},
			"endpoint_port": schema.Int64Attribute{
				Optional: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
				Description: "An endpoint port can be left blank to allow remote connection from any port.",
			},
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Name of the WireGuard interface the peer belongs to.",
			},
			"persistent_keepalive": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
				Description: "A seconds interval, between 1 and 65535 inclusive, of how often to send an authenticated empty packet to the peer for the purpose of keeping a stateful firewall or NAT mapping valid persistently. For example, if the interface very rarely sends traffic, but it might at anytime receive traffic from a peer, and it is behind NAT, the interface might benefit from having a persistent keepalive interval of 25 seconds.",
			},
			"preshared_key": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Description: "A base64 preshared key. Optional, and may be omitted. This option adds an additional layer of symmetric-key cryptography to be mixed into the already existing public-key cryptography, for post-quantum resistance.",
			},
			"public_key": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Description: "The remote peer's calculated public key.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (i *interfaceWireguardPeer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan interfaceWireguardPeerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := i.client.AddInterfaceWireguardPeer(modelToInterfaceWireguardPeer(&plan))
	if err != nil {
		resp.Diagnostics.AddError("creation failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceWireguardPeerToModel(created, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (i *interfaceWireguardPeer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state interfaceWireguardPeerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := i.client.FindInterfaceWireguardPeer(state.Interface.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remote resource",
			fmt.Sprintf("Could not read interfaceWireguardPeer with interface name %q", state.Interface.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(interfaceWireguardPeerToModel(resource, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (i *interfaceWireguardPeer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan interfaceWireguardPeerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := i.client.UpdateInterfaceWireguardPeer(modelToInterfaceWireguardPeer(&plan))
	if err != nil {
		resp.Diagnostics.AddError("update failed", err.Error())
		return
	}

	resp.Diagnostics.Append(interfaceWireguardPeerToModel(updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (i *interfaceWireguardPeer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state interfaceWireguardPeerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := i.client.DeleteInterfaceWireguardPeer(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Could not delete interfaceWireguardPeer", err.Error())
		return
	}
}

func (i *interfaceWireguardPeer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("interface"), req, resp)
}

type interfaceWireguardPeerModel struct {
	ID                     tftypes.String `tfsdk:"id"`
	AllowedAddress         tftypes.String `tfsdk:"allowed_address"`
	Comment                tftypes.String `tfsdk:"comment"`
	Disabled               tftypes.Bool   `tfsdk:"disabled"`
	EndpointAddress        tftypes.String `tfsdk:"endpoint_address"`
	EndpointPort           tftypes.Int64  `tfsdk:"endpoint_port"`
	Interface              tftypes.String `tfsdk:"interface"`
	PersistentKeepalive    tftypes.Int64  `tfsdk:"persistent_keepalive"`
	PresharedKey           tftypes.String `tfsdk:"preshared_key"`
	PublicKey              tftypes.String `tfsdk:"public_key"`
	CurrentEndpointAddress tftypes.String `tfsdk:"current_endpoint_address"`
	CurrentEndpointPort    tftypes.Int64  `tfsdk:"current_endpoint_port"`
}

func interfaceWireguardPeerToModel(i *client.InterfaceWireguardPeer, m *interfaceWireguardPeerModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if i == nil {
		diags.AddError("Interface Wireguard Peer cannot be nil", "Cannot build model from nil object")
		return diags
	}
	m.ID = tftypes.StringValue(i.Id)
	m.AllowedAddress = tftypes.StringValue(i.AllowedAddress)
	m.Comment = tftypes.StringValue(i.Comment)
	m.Disabled = tftypes.BoolValue(i.Disabled)
	m.EndpointAddress = tftypes.StringValue(i.EndpointAddress)
	m.EndpointPort = tftypes.Int64Value(i.EndpointPort)
	m.Interface = tftypes.StringValue(i.Interface)
	m.PersistentKeepalive = tftypes.Int64Value(i.PersistentKeepalive)
	m.PresharedKey = tftypes.StringValue(i.PresharedKey)
	m.PublicKey = tftypes.StringValue(i.PublicKey)

	return diags
}

func modelToInterfaceWireguardPeer(m *interfaceWireguardPeerModel) *client.InterfaceWireguardPeer {
	return &client.InterfaceWireguardPeer{
		Id:                  m.ID.ValueString(),
		AllowedAddress:      m.AllowedAddress.ValueString(),
		Comment:             m.Comment.ValueString(),
		Disabled:            m.Disabled.ValueBool(),
		EndpointAddress:     m.EndpointAddress.ValueString(),
		EndpointPort:        m.EndpointPort.ValueInt64(),
		Interface:           m.Interface.ValueString(),
		PersistentKeepalive: m.PersistentKeepalive.ValueInt64(),
		PresharedKey:        m.PresharedKey.ValueString(),
		PublicKey:           m.PublicKey.ValueString(),
	}
}
