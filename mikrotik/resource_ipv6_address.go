package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type ipv6Address struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ipv6Address{}
	_ resource.ResourceWithConfigure   = &ipv6Address{}
	_ resource.ResourceWithImportState = &ipv6Address{}
)

// NewIpv6AddressResource is a helper function to simplify the provider implementation.
func NewIpv6AddressResource() resource.Resource {
	return &ipv6Address{}
}

func (r *ipv6Address) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *ipv6Address) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipv6_address"
}

// Schema defines the schema for the resource.
func (s *ipv6Address) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik Ipv6Address.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier for this resource.",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "The IPv6 address and prefix length of the interface using slash notation.",
			},
			"advertise": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),

				Description: "Whether to enable stateless address configuration. The prefix of that address is automatically advertised to hosts using ICMPv6 protocol. The option is set by default for addresses with prefix length 64.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The comment for the IPv6 address assignment.",
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),

				Description: "Whether to disable IPv6 address.",
			},
			"eui_64": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),

				Description: "Whether to calculate EUI-64 address and use it as last 64 bits of the IPv6 address.",
			},
			"from_pool": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Name of the pool from which prefix will be taken to construct IPv6 address taking last part of the address from address property.",
			},
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "The interface on which the IPv6 address is assigned.",
			},
			"no_dad": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),

				Description: "If set indicates that address is anycast address and Duplicate Address Detection should not be performed.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ipv6Address) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel ipv6AddressModel
	var mikrotikModel client.Ipv6Address
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *ipv6Address) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel ipv6AddressModel
	var mikrotikModel client.Ipv6Address
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ipv6Address) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel ipv6AddressModel
	var mikrotikModel client.Ipv6Address
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ipv6Address) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel ipv6AddressModel
	var mikrotikModel client.Ipv6Address
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *ipv6Address) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type ipv6AddressModel struct {
	Id        tftypes.String `tfsdk:"id"`
	Address   tftypes.String `tfsdk:"address"`
	Advertise tftypes.Bool   `tfsdk:"advertise"`
	Comment   tftypes.String `tfsdk:"comment"`
	Disabled  tftypes.Bool   `tfsdk:"disabled"`
	Eui64     tftypes.Bool   `tfsdk:"eui_64"`
	FromPool  tftypes.String `tfsdk:"from_pool"`
	Interface tftypes.String `tfsdk:"interface"`
	NoDad     tftypes.Bool   `tfsdk:"no_dad"`
}
