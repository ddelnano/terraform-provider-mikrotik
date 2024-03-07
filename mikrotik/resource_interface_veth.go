package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
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
				Description: "Identifier of this resource assigned by RouterOS.",
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
				Computed:    true,
				Description: "Whether the interface is running.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *interfaceVeth) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel interfaceVethModel
	var mikrotikModel client.InterfaceVeth
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *interfaceVeth) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceVethModel
	var mikrotikModel client.InterfaceVeth
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *interfaceVeth) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceVethModel
	var mikrotikModel client.InterfaceVeth
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *interfaceVeth) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceVethModel
	var mikrotikModel client.InterfaceVeth
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (i *interfaceVeth) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type interfaceVethModel struct {
	Id         tftypes.String `tfsdk:"id"`
	Name       tftypes.String `tfsdk:"name"`
	Comment    tftypes.String `tfsdk:"comment"`
	Disabled   tftypes.Bool   `tfsdk:"disabled"`
	Address	   tftypes.String `tfsdk:"address"`
	Gateway	   tftypes.String `tfsdk:"gateway"`
	Gateway6   tftypes.String `tfsdk:"gateway6"`
	Running    tftypes.Bool   `tfsdk:"running"`
}
