package mikrotik

import (
	"context"
	"fmt"
	"os"
	"strings"

	mt "github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		// add default value in description, if it was declared in the resource's schema.
		if s.Default != nil {
			if s.Default == "" {
				desc += " Default: `\"\"`."
			} else {
				desc += fmt.Sprintf(" Default: `%v`.", s.Default)
			}
		}

		return strings.TrimSpace(desc)
	}
}

func Provider(client *mt.Mikrotik) *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname of the MikroTik router",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User account for MikroTik api",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for MikroTik api",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use TLS when connecting to MikroTik or not",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to MikroTik's certificate authority",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Insecure connection does not verify MikroTik's TLS certificate",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mikrotik_bgp_instance":         resourceBgpInstance(),
			"mikrotik_bridge_vlan":          resourceBridgeVlan(),
			"mikrotik_dhcp_server_network":  resourceDhcpServerNetwork(),
			"mikrotik_firewall_filter_rule": resourceFirewallFilterRule(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		if client != nil {
			return client, nil
		}
		var diags diag.Diagnostics

		address := d.Get("host").(string)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		tls := d.Get("tls").(bool)
		caCertificate := d.Get("ca_certificate").(string)
		insecure := d.Get("insecure").(bool)

		if v := os.Getenv("MIKROTIK_HOST"); v != "" {
			address = v
		}
		if v := os.Getenv("MIKROTIK_USER"); v != "" {
			username = v
		}
		if v := os.Getenv("MIKROTIK_PASSWORD"); v != "" {
			password = v
		}
		if v := os.Getenv("MIKROTIK_TLS"); v != "" {
			tlsValue, err := utils.ParseBool(v)
			if err != nil {
				diags = append(diags,
					diag.FromErr(fmt.Errorf("could not parse MIKROTIK_TLS environment variable: %w", err))...)
			}
			tls = tlsValue
		}
		if v := os.Getenv("MIKROTIK_CA_CERTIFICATE"); v != "" {
			caCertificate = v
		}
		if v := os.Getenv("MIKROTIK_INSECURE"); v != "" {
			insecureValue, err := utils.ParseBool(v)
			if err != nil {
				diags = append(diags,
					diag.FromErr(fmt.Errorf("could not parse MIKROTIK_INSECURE environment variable: %w", err))...)
			}
			insecure = insecureValue
		}

		return mt.NewClient(address, username, password, tls, caCertificate, insecure), diags
	}

	return provider
}

func NewProvider() *schema.Provider {
	return Provider(nil)
}
