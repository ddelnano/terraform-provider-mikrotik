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
		},
		ResourcesMap: map[string]*schema.Resource{
			"mikrotik_dns_record": resourceRecord(),
			"mikrotik_dhcp_lease": resourceLease(),
			"mikrotik_script":     resourceScript(),
		},
		ConfigureFunc: mikrotikConfigure,
	}
}

func mikrotikConfigure(d *schema.ResourceData) (c interface{}, err error) {
	address := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	c = client.NewClient(address, username, password)
	return
}

type mikrotikConn struct {
	host     string
	username string
	password string
}
