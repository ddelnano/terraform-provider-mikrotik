package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_PASSWORD", nil),
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
			"mikrotik_scheduler":    resourceScheduler(),
			"mikrotik_script":       resourceScript(),
			"mikrotik_pool":         resourcePool(),
			"mikrotik_bgp_instance": resourceBgpInstance(),
			"mikrotik_bgp_peer":     resourceBgpPeer(),
		},
		ConfigureContextFunc: mikrotikConfigure,
	}
}

func mikrotikConfigure(ctx context.Context, d *schema.ResourceData) (c interface{}, diags diag.Diagnostics) {
	address := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	tls := d.Get("tls").(bool)
	caCertificate := d.Get("ca_certificate").(string)
	insecure := d.Get("insecure").(bool)
	c = client.NewClient(address, username, password, tls, caCertificate, insecure)
	return
}
