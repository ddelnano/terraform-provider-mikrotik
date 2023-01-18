package mikrotik

import (
	"context"
	"os"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderFramework struct {
}

var _ provider.Provider = (*ProviderFramework)(nil)

func NewProviderFramework() provider.Provider {
	return &ProviderFramework{}
}

func (p *ProviderFramework) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mikrotik"
}

func (p *ProviderFramework) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:    true,
				Description: "Hostname of the MikroTik router",
			},
			"username": schema.StringAttribute{
				Required:    true,
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
	var data mikrotikProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

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

	if !resp.Diagnostics.HasError() {
		resp.ResourceData = client.NewClient(mikrotikHost, mikrotikUser, mikrotikPassword,
			mikrotikTLS, mikrotikCACertificates, mikrotikInsecure)
	}
}

func (p *ProviderFramework) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ProviderFramework) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

type mikrotikProviderModel struct {
	Host          types.String `tfsdk:"host"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Tls           types.Bool   `tfsdk:"tls"`
	CACertificate types.String `tfsdk:"ca_certificate"`
	Insecure      types.Bool   `tfsdk:"insecure"`
}
