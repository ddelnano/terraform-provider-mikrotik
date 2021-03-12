package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_HOST", nil),
				Description: "Hostname of the mikrotik router",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_USER", nil),
				Description: "User account for mikrotik api",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_PASSWORD", nil),
				Description: "Password for mikrotik api",
			},
			"tls": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_TLS", false),
				Description: "Whether to use TLS when connecting to MikroTik or not",
			},
			"ca": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_CA", nil),
				Description: "Path to MikroTik's certificate authority",
			},
			"verify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MIKROTIK_VERIFY", false),
				Description: "Whether to verify TLS certification or not",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mikrotik_dns_record":   resourceRecord(),
			"mikrotik_dhcp_lease":   resourceLease(),
			"mikrotik_scheduler":    resourceScheduler(),
			"mikrotik_script":       resourceScript(),
			"mikrotik_pool":         resourcePool(),
			"mikrotik_bgp_instance": resourceBgpInstance(),
		},
		ConfigureFunc: mikrotikConfigure,
	}
}

func mikrotikConfigure(d *schema.ResourceData) (c interface{}, err error) {
	address := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	tls := d.Get("tls").(bool)
	ca := d.Get("ca").(string)
	verify := d.Get("verify").(bool)
	c = client.NewClient(address, username, password, tls, ca, verify)
	return
}

type mikrotikConn struct {
	host     string
	username string
	password string
}
