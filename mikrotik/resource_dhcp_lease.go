package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/types/defaultaware"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type dhcpLease struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpLease{}
	_ resource.ResourceWithConfigure   = &dhcpLease{}
	_ resource.ResourceWithImportState = &dhcpLease{}
)

// NewDhcpLeaseResource is a helper function to simplify the provider implementation.
func NewDhcpLeaseResource() resource.Resource {
	return &dhcpLease{}
}

func (r *dhcpLease) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *dhcpLease) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_lease"
}

// Schema defines the schema for the resource.
func (s *dhcpLease) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a DHCP lease on the MikroTik device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique resource identifier.",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "The IP address of the DHCP lease to be created.",
			},
			"macaddress": schema.StringAttribute{
				Required:    true,
				Description: "The MAC addreess of the DHCP lease to be created.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The comment of the DHCP lease to be created.",
			},
			"blocked": defaultaware.BoolAttribute(
				schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether to block access for this DHCP client (true|false).",
				},
			),
			"dynamic": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Whether the dhcp lease is static or dynamic. Dynamic leases are not guaranteed to continue to be assigned to that specific device. Defaults to false.",
			},
			"hostname": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The hostname of the device",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dhcpLease) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel dhcpLeaseModel
	var mikrotikModel client.DhcpLease
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpLease) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel dhcpLeaseModel
	var mikrotikModel client.DhcpLease
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpLease) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel dhcpLeaseModel
	var mikrotikModel client.DhcpLease
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpLease) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel dhcpLeaseModel
	var mikrotikModel client.DhcpLease
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *dhcpLease) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	utils.ImportUppercaseWrapper(resource.ImportStatePassthroughID)(ctx, path.Root("id"), req, resp)
}

type dhcpLeaseModel struct {
	Id          tftypes.String `tfsdk:"id"`
	Address     tftypes.String `tfsdk:"address"`
	MacAddress  tftypes.String `tfsdk:"macaddress"`
	Comment     tftypes.String `tfsdk:"comment"`
	BlockAccess tftypes.Bool   `tfsdk:"blocked"`
	Dynamic     tftypes.Bool   `tfsdk:"dynamic"`
	Hostname    tftypes.String `tfsdk:"hostname"`
}
