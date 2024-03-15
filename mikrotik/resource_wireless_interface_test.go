package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestWirelessInterface_basic(t *testing.T) {
	resourceName := "mikrotik_wireless_interface.testacc"
	name := acctest.RandomWithPrefix("ssid")
	resource.Test(t,
		resource.TestCase{
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
					resource "mikrotik_wireless_interface" "testacc" {
						name = %q
						mode = %q
						ssid = %q
						vlan_id = 2
						hide_ssid = false
					}`, name, client.WirelessInterfaceModeAPBridge, name+"-ssid"),

					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", name),
						resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "mode", client.WirelessInterfaceModeAPBridge),
						resource.TestCheckResourceAttr(resourceName, "ssid", name+"-ssid"),
						resource.TestCheckResourceAttr(resourceName, "hide_ssid", "false"),
						resource.TestCheckResourceAttr(resourceName, "vlan_id", "2"),
					),
				},
				{
					Config: fmt.Sprintf(`
					resource mikrotik_wireless_interface testacc {
						name = %q
						mode = %q
						disabled = false
						ssid = %q
						hide_ssid = true
					}`, name, client.WirelessInterfaceModeAPBridge, name+"-ssid"),

					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "name", name),
						resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "mode", client.WirelessInterfaceModeAPBridge),
						resource.TestCheckResourceAttr(resourceName, "ssid", name+"-ssid"),
						resource.TestCheckResourceAttr(resourceName, "hide_ssid", "true"),
					),
				},
				{
					ImportState:       true,
					ImportStateVerify: true,
					ResourceName:      resourceName,
				},
			},
		},
	)
}
