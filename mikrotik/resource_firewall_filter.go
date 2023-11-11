package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type firewallFilterRule struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallFilterRule{}
	_ resource.ResourceWithConfigure   = &firewallFilterRule{}
	_ resource.ResourceWithImportState = &firewallFilterRule{}
)

// NewFirewallFilterRuleResource is a helper function to simplify the provider implementation.
func NewFirewallFilterRuleResource() resource.Resource {
	return &firewallFilterRule{}
}

func (r *firewallFilterRule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *firewallFilterRule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_filter_rule"
}

// Schema defines the schema for the resource.
func (s *firewallFilterRule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik FirewallFilterRule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("accept"),
				Description: "Action to take if packet is matched by the rule.",
			},
			"chain": schema.StringAttribute{
				Required:    true,
				Description: "Specifies to which chain rule will be added. If the input does not match the name of an already defined chain, a new chain will be created.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Comment to the rule.",
			},
			"connection_state": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: tftypes.StringType,
				Description: "Interprets the connection tracking analysis data for a particular packet.",
			},
			"dst_port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of destination port numbers or port number ranges.",
			},
			"in_interface": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Interface the packet has entered the router.",
			},
			"in_interface_list": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set of interfaces defined in interface list. Works the same as in-interface.",
			},
			"out_interface_list": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set of interfaces defined in interface list. Works the same as out-interface.",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("tcp"),
				Description: "Matches particular IP protocol specified by protocol name or number.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *firewallFilterRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel firewallFilterRuleModel
	var mikrotikModel client.FirewallFilterRule
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *firewallFilterRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel firewallFilterRuleModel
	var mikrotikModel client.FirewallFilterRule
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *firewallFilterRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel firewallFilterRuleModel
	var mikrotikModel client.FirewallFilterRule
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallFilterRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel firewallFilterRuleModel
	var mikrotikModel client.FirewallFilterRule
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *firewallFilterRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type firewallFilterRuleModel struct {
	Id               tftypes.String `tfsdk:"id"`
	Action           tftypes.String `tfsdk:"action"`
	Chain            tftypes.String `tfsdk:"chain"`
	Comment          tftypes.String `tfsdk:"comment"`
	ConnectionState  tftypes.Set    `tfsdk:"connection_state"`
	DestPort         tftypes.String `tfsdk:"dst_port"`
	InInterface      tftypes.String `tfsdk:"in_interface"`
	InInterfaceList  tftypes.String `tfsdk:"in_interface_list"`
	OutInterfaceList tftypes.String `tfsdk:"out_interface_list"`
	Protocol         tftypes.String `tfsdk:"protocol"`
}
