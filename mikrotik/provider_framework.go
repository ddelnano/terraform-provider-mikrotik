package mikrotik

import (
	"context"
	"os"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/types/defaultaware"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderFramework struct {
	predefinedAPIClient *client.Mikrotik
}

var (
	_ provider.Provider = (*ProviderFramework)(nil)
)

func NewProviderFramework(c *client.Mikrotik) provider.Provider {
	return &ProviderFramework{
		predefinedAPIClient: c,
	}
}

func (p *ProviderFramework) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mikrotik"
}

func (p *ProviderFramework) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname of the MikroTik router",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "User account for MikroTik api",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for MikroTik api",
			},
			"tls": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to use TLS when connecting to MikroTik or not",
			},
			"ca_certificate": schema.StringAttribute{
				Optional:    true,
				Description: "Path to MikroTik's certificate authority",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Insecure connection does not verify MikroTik's TLS certificate",
			},
		},
	}
}

func (p *ProviderFramework) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.predefinedAPIClient != nil {
		resp.DataSourceData = p.predefinedAPIClient
		resp.ResourceData = p.predefinedAPIClient

		return
	}

	var data mikrotikProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if configuration sets provider configuration fields, the values must be known during provider configuration
	// otherwise, it is not possible to setup the client
	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown value for MikroTik host",
			"The provider cannot create MikroTik API client as the 'host' is unknown at this moment. "+
				"Either target apply the source of the value first, set it statically or use MIKROTIK_HOST environment variable.",
		)
	}
	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown value for MikroTik username",
			"The provider cannot create MikroTik API client as the 'username' is unknown at this moment. "+
				"Either target apply the source of the value first, set it statically or use MIKROTIK_USER environment variable.",
		)
	}
	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown value for MikroTik password",
			"The provider cannot create MikroTik API client as the 'password' is unknown at this moment. "+
				"Either target apply the source of the value first, set it statically or use MIKROTIK_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var mikrotikHost, mikrotikUser, mikrotikPassword, mikrotikCACertificates string
	var mikrotikTLS, mikrotikInsecure bool

	mikrotikHost = data.Host.ValueString()
	if v := os.Getenv("MIKROTIK_HOST"); v != "" {
		mikrotikHost = v
	}

	mikrotikUser = data.Username.ValueString()
	if v := os.Getenv("MIKROTIK_USER"); v != "" {
		mikrotikUser = v
	}

	mikrotikPassword = data.Password.ValueString()
	if v := os.Getenv("MIKROTIK_PASSWORD"); v != "" {
		mikrotikPassword = v
	}

	if !data.Tls.IsUnknown() {
		mikrotikTLS = data.Tls.ValueBool()
	}
	if v := os.Getenv("MIKROTIK_TLS"); v != "" {
		tls, err := utils.ParseBool(v)
		if err != nil {
			resp.Diagnostics.AddError("Could not parse MIKROTIK_TLS environment variable", err.Error())
		}
		mikrotikTLS = tls
	}

	mikrotikCACertificates = data.CACertificate.ValueString()
	if v := os.Getenv("MIKROTIK_CA_CERTIFICATE"); v != "" {
		mikrotikCACertificates = v
	}

	if !data.Insecure.IsUnknown() {
		mikrotikInsecure = data.Insecure.ValueBool()
	}
	if v := os.Getenv("MIKROTIK_INSECURE"); v != "" {
		insecure, err := utils.ParseBool(v)
		if err != nil {
			resp.Diagnostics.AddError("Could not parse MIKROTIK_INSECURE environment variable", err.Error())
		}
		mikrotikInsecure = insecure
	}

	if mikrotikHost == "" {
		resp.Diagnostics.AddError("Mikrotik 'host' is missing in configuration",
			"Provide it via 'host' provider configuration attribute or MIKROTIK_HOST environment variable")
	}

	if mikrotikUser == "" {
		resp.Diagnostics.AddError("Mikrotik 'username' is missing in configuration",
			"Provide it via 'host' provider configuration attribute or MIKROTIK_USER environment variable")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := client.NewClient(mikrotikHost, mikrotikUser, mikrotikPassword,
		mikrotikTLS, mikrotikCACertificates, mikrotikInsecure)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *ProviderFramework) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ProviderFramework) Resources(ctx context.Context) []func() resource.Resource {
	return defaultaware.WrapResources([]func() resource.Resource{
		NewBgpInstanceResource,
		NewBgpPeerResource,
		NewBridgePortResource,
		NewBridgeResource,
		NewBridgeVlanResource,
		NewDhcpLeaseResource,
		NewDhcpServerNetworkResource,
		NewDhcpServerResource,
		NewDnsRecordResource,
		NewFirewallFilterRuleResource,
		NewInterfaceListMemberResource,
		NewInterfaceListResource,
		NewInterfaceWireguardPeerResource,
		NewInterfaceWireguardResource,
		NewIpAddressResource,
		NewIpv6AddressResource,
		NewPoolResource,
		NewSchedulerResource,
		NewScriptResource,
		NewVlanInterfaceResource,
	},
	)
}

type mikrotikProviderModel struct {
	Host          types.String `tfsdk:"host"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Tls           types.Bool   `tfsdk:"tls"`
	CACertificate types.String `tfsdk:"ca_certificate"`
	Insecure      types.Bool   `tfsdk:"insecure"`
}
