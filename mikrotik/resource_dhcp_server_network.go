package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type dhcpServerNetwork struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpServerNetwork{}
	_ resource.ResourceWithConfigure   = &dhcpServerNetwork{}
	_ resource.ResourceWithImportState = &dhcpServerNetwork{}
)

// NewDhcpServerNetworkResource is a helper function to simplify the provider implementation.
func NewDhcpServerNetworkResource() resource.Resource {
	return &dhcpServerNetwork{}
}

func (r *dhcpServerNetwork) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *dhcpServerNetwork) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_server_network"
}

// Schema defines the schema for the resource.
func (s *dhcpServerNetwork) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DHCP network resource within Mikrotik device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The network DHCP server(s) will lease addresses from.",
			},
			"netmask": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0"),
				Description: "The actual network mask to be used by DHCP client. If set to '0' - netmask from network address will be used.",
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0.0.0.0"),
				Description: "The default gateway to be used by DHCP Client.",
			},
			"dns_server": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DHCP client will use these as the default DNS servers.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dhcpServerNetwork) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel dhcpServerNetworkModel
	var mikrotikModel client.DhcpServerNetwork
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpServerNetwork) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel dhcpServerNetworkModel
	var mikrotikModel client.DhcpServerNetwork
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpServerNetwork) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel dhcpServerNetworkModel
	var mikrotikModel client.DhcpServerNetwork
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpServerNetwork) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel dhcpServerNetworkModel
	var mikrotikModel client.DhcpServerNetwork
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *dhcpServerNetwork) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type dhcpServerNetworkModel struct {
	Id        tftypes.String `tfsdk:"id"`
	Comment   tftypes.String `tfsdk:"comment"`
	Address   tftypes.String `tfsdk:"address"`
	Netmask   tftypes.String `tfsdk:"netmask"`
	Gateway   tftypes.String `tfsdk:"gateway"`
	DnsServer tftypes.String `tfsdk:"dns_server"`
}
