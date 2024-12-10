package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type wirelessInterface struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wirelessInterface{}
	_ resource.ResourceWithConfigure   = &wirelessInterface{}
	_ resource.ResourceWithImportState = &wirelessInterface{}
)

// NewWirelessInterfaceResource is a helper function to simplify the provider implementation.
func NewWirelessInterfaceResource() resource.Resource {
	return &wirelessInterface{}
}

func (r *wirelessInterface) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wirelessInterface) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_interface"
}

// Schema defines the schema for the resource.
func (s *wirelessInterface) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik WirelessInterface.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier for this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the interface.",
			},
			"master_interface": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Name of wireless interface that has virtual-ap capability. Virtual AP interface will only work if master interface is in ap-bridge, bridge, station or wds-slave mode. This property is only for virtual AP interfaces.",
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("station"),
				Description: "Selection between different station and access point (AP) modes.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether interface is disabled.",
			},
			"security_profile": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "Name of profile from security-profiles.",
			},
			"ssid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SSID (service set identifier) is a name that identifies wireless network.",
			},
			"hide_ssid": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "This property has an effect only in AP mode.",
			},
			"vlan_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				Description: "VLAN identification number.",
			},
			"vlan_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("no-tag"),
				Description: "Three VLAN modes are available: no-tag|use-service-tag|use-tag.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wirelessInterface) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wirelessInterfaceModel
	var mikrotikModel client.WirelessInterface
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wirelessInterface) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wirelessInterfaceModel
	var mikrotikModel client.WirelessInterface
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wirelessInterface) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wirelessInterfaceModel
	var mikrotikModel client.WirelessInterface
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wirelessInterface) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wirelessInterfaceModel
	var mikrotikModel client.WirelessInterface
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wirelessInterface) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type wirelessInterfaceModel struct {
	Id              tftypes.String `tfsdk:"id"`
	Name            tftypes.String `tfsdk:"name"`
	MasterInterface tftypes.String `tfsdk:"master_interface"`
	Mode            tftypes.String `tfsdk:"mode"`
	Disabled        tftypes.Bool   `tfsdk:"disabled"`
	SecurityProfile tftypes.String `tfsdk:"security_profile"`
	SSID            tftypes.String `tfsdk:"ssid"`
	HideSSID        tftypes.Bool   `tfsdk:"hide_ssid"`
	VlanID          tftypes.Int64  `tfsdk:"vlan_id"`
	VlanMode        tftypes.String `tfsdk:"vlan_mode"`
}
