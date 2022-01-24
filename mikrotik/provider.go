package mikrotik

import (
	"context"

	mt "github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider(client *mt.Mikrotik) *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_HOST", nil),
				Description: "Hostname of the mikrotik router",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_USER", nil),
				Description: "User account for mikrotik api",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_PASSWORD", ""),
				Description: "Password for mikrotik api",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_TLS", false),
				Description: "Whether to use TLS when connecting to MikroTik or not",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_CA_CERTIFICATE", ""),
				Description: "Path to MikroTik's certificate authority",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_INSECURE", false),
				Description: "Insecure connection does not verify MikroTik's TLS certificate",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mikrotik_dns_record":   resourceRecord(),
			"mikrotik_dhcp_lease":   resourceLease(),
			"mikrotik_ip_address":   resourceIpAddress(),
			"mikrotik_scheduler":    resourceScheduler(),
			"mikrotik_script":       resourceScript(),
			"mikrotik_pool":         resourcePool(),
			"mikrotik_bgp_instance": resourceBgpInstance(),
			"mikrotik_bgp_peer":     resourceBgpPeer(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		if client != nil {
			return client, nil
		}

		address := d.Get("host").(string)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		tls := d.Get("tls").(bool)
		caCertificate := d.Get("ca_certificate").(string)
		insecure := d.Get("insecure").(bool)

		return mt.NewClient(address, username, password, tls, caCertificate, insecure), nil
	}

	return provider
}

func NewProvider() *schema.Provider {
	return Provider(nil)
}
