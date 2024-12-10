package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type dnsRecord struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &dnsRecord{}
	_ resource.ResourceWithConfigure        = &dnsRecord{}
	_ resource.ResourceWithConfigValidators = &dnsRecord{}
	_ resource.ResourceWithImportState      = &dnsRecord{}
)

// NewDnsRecordResource is a helper function to simplify the provider implementation.
func NewDnsRecordResource() resource.Resource {
	return &dnsRecord{}
}

func (r *dnsRecord) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *dnsRecord) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

// Schema defines the schema for the resource.
func (s *dnsRecord) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a DNS record on the MikroTik device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the DNS hostname to be created.",
			},
			"regexp": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Regular expression against which domain names should be verified.",
			},
			"ttl": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The ttl of the DNS record.",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "The A record to be returend from the DNS hostname.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The comment text associated with the DNS record.",
			},
		},
	}
}

func (r *dnsRecord) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("name"),
			path.MatchRoot("regexp"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dnsRecord) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel dnsRecordModel
	var mikrotikModel client.DnsRecord
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsRecord) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel dnsRecordModel
	var mikrotikModel client.DnsRecord
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsRecord) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel dnsRecordModel
	var mikrotikModel client.DnsRecord
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsRecord) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel dnsRecordModel
	var mikrotikModel client.DnsRecord
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *dnsRecord) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type dnsRecordModel struct {
	Id      tftypes.String `tfsdk:"id"`
	Name    tftypes.String `tfsdk:"name"`
	Regexp  tftypes.String `tfsdk:"regexp"`
	Ttl     tftypes.Int64  `tfsdk:"ttl"`
	Address tftypes.String `tfsdk:"address"`
	Comment tftypes.String `tfsdk:"comment"`
}
