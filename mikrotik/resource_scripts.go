package mikrotik

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/go-routeros/routeros"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: RecordImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// func RecordImportScript(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
// 	name := d.Id()
// 	c := m.(client.Mikrotik)

// 	record, err := c.FindScript(name)

// 	if err != nil {
// 		return nil, err
// 	}
// 	recordToData(record, d)
// 	return []*schema.ResourceData{d}, nil
// }

func findScript(c *routeros.Client, name string) (*client.DnsRecord, error) {
	// r, err := c.Run("/system/script/print")
	// found := false
	// var sentence *proto.Sentence

	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
