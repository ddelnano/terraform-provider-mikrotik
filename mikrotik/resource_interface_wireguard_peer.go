package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
// TODO: Reevaluate the Computed schema attributes and determine if that is correct
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
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "List of IP (v4 or v6) addresses with CIDR masks from which incoming traffic for this peer is allowed and to which outgoing traffic for this peer is directed. The catch-all 0.0.0.0/0 may be specified for matching all IPv4 addresses, and ::/0 may be specified for matching all IPv6 addresses.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "An endpoint IP or hostname can be left blank to allow remote connection from any address.",
			},
			"endpoint_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
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
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "A base64 preshared key. Optional, and may be omitted. This option adds an additional layer of symmetric-key cryptography to be mixed into the already existing public-key cryptography, for post-quantum resistance.",
			},
			"public_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The remote peer's calculated public key.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (i *interfaceWireguardPeer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel interfaceWireguardPeerModel
	var mikrotikModel client.InterfaceWireguardPeer
	GenericCreateResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (i *interfaceWireguardPeer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceWireguardPeerModel
	var mikrotikModel client.InterfaceWireguardPeer
	GenericReadResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (i *interfaceWireguardPeer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceWireguardPeerModel
	var mikrotikModel client.InterfaceWireguardPeer
	GenericUpdateResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (i *interfaceWireguardPeer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceWireguardPeerModel
	var mikrotikModel client.InterfaceWireguardPeer
	GenericDeleteResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

func (i *interfaceWireguardPeer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type interfaceWireguardPeerModel struct {
	Id                  tftypes.String `tfsdk:"id"`
	AllowedAddress      tftypes.String `tfsdk:"allowed_address"`
	Comment             tftypes.String `tfsdk:"comment"`
	Disabled            tftypes.Bool   `tfsdk:"disabled"`
	EndpointAddress     tftypes.String `tfsdk:"endpoint_address"`
	EndpointPort        tftypes.Int64  `tfsdk:"endpoint_port"`
	Interface           tftypes.String `tfsdk:"interface"`
	PersistentKeepalive tftypes.Int64  `tfsdk:"persistent_keepalive"`
	PresharedKey        tftypes.String `tfsdk:"preshared_key"`
	PublicKey           tftypes.String `tfsdk:"public_key"`
}
