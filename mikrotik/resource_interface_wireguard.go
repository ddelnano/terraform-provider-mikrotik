package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
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
				Default:     stringdefault.StaticString(""),
				Description: "Comment associated with interface wireguard.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Boolean for whether or not the interface wireguard is disabled.",
			},
			"listen_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(13231),
				Description: "Port for WireGuard service to listen on for incoming sessions.",
			},
			"mtu": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1420),
				Validators: []validator.Int64{
					int64validator.Between(0, 65536),
				},
				Description: "Layer3 Maximum transmission unit.",
			},
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
	var terraformModel interfaceWireguardModel
	var mikrotikModel client.InterfaceWireguard
	GenericCreateResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (i *interfaceWireguard) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceWireguardModel
	var mikrotikModel client.InterfaceWireguard
	GenericReadResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (i *interfaceWireguard) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceWireguardModel
	var mikrotikModel client.InterfaceWireguard
	GenericUpdateResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (i *interfaceWireguard) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceWireguardModel
	var mikrotikModel client.InterfaceWireguard
	GenericDeleteResource(&terraformModel, &mikrotikModel, i.client)(ctx, req, resp)
}

func (i *interfaceWireguard) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type interfaceWireguardModel struct {
	Id         tftypes.String `tfsdk:"id"`
	Name       tftypes.String `tfsdk:"name"`
	Comment    tftypes.String `tfsdk:"comment"`
	Disabled   tftypes.Bool   `tfsdk:"disabled"`
	ListenPort tftypes.Int64  `tfsdk:"listen_port"`
	Mtu        tftypes.Int64  `tfsdk:"mtu"`
	PrivateKey tftypes.String `tfsdk:"private_key"`
	PublicKey  tftypes.String `tfsdk:"public_key"`
	Running    tftypes.Bool   `tfsdk:"running"`
}
