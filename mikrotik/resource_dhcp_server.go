package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type dhcpServer struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpServer{}
	_ resource.ResourceWithConfigure   = &dhcpServer{}
	_ resource.ResourceWithImportState = &dhcpServer{}
)

// NewDhcpServerResource is a helper function to simplify the provider implementation.
func NewDhcpServerResource() resource.Resource {
	return &dhcpServer{}
}

func (r *dhcpServer) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *dhcpServer) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_server"
}

// Schema defines the schema for the resource.
func (s *dhcpServer) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DHCP server resource within MikroTik device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Reference name.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Disable this DHCP server instance.",
			},
			"add_arp": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to add dynamic ARP entry. If set to no either ARP mode should be enabled on that interface or static ARP entries should be administratively defined.",
			},
			"address_pool": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("static-only"),
				Description: "IP pool, from which to take IP addresses for the clients. If set to static-only, then only the clients that have a static lease (added in lease submenu) will be allowed.",
			},
			"authoritative": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("yes"),
				Description: "Option changes the way how server responds to DHCP requests.",
			},
			"interface": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("*0"),
				Description: "Interface on which server will be running.",
			},
			"lease_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Script that will be executed after lease is assigned or de-assigned. Internal \"global\" variables that can be used in the script.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dhcpServer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel dhcpServerModel
	var mikrotikModel client.DhcpServer
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpServer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel dhcpServerModel
	var mikrotikModel client.DhcpServer
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpServer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel dhcpServerModel
	var mikrotikModel client.DhcpServer
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpServer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel dhcpServerModel
	var mikrotikModel client.DhcpServer
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *dhcpServer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type dhcpServerModel struct {
	Id            tftypes.String `tfsdk:"id"`
	Name          tftypes.String `tfsdk:"name"`
	Disabled      tftypes.Bool   `tfsdk:"disabled"`
	AddArp        tftypes.Bool   `tfsdk:"add_arp"`
	AddressPool   tftypes.String `tfsdk:"address_pool"`
	Authoritative tftypes.String `tfsdk:"authoritative"`
	Interface     tftypes.String `tfsdk:"interface"`
	LeaseScript   tftypes.String `tfsdk:"lease_script"`
}
