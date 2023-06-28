package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDhcpServerNetwork_basic(t *testing.T) {

	resourceName := "mikrotik_dhcp_server_network.testacc"

	netmask := "24"
	address := "10.10.10.0/" + netmask
	gateway := "10.10.10.2"
	dnsServer := "10.10.10.3"
	comment := "Terraform managed"
	dnsServerUpdated := "192.168.5.3"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDhcpServerNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpServerNetwork(address, netmask, gateway, dnsServer, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpServerNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "netmask", netmask),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "dns_server", dnsServer),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccDhcpServerNetwork(address, netmask, gateway, dnsServerUpdated, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpServerNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "netmask", netmask),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "dns_server", dnsServerUpdated),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func testAccDhcpServerNetwork(address, netmask, gateway, dns_server, comment string) string {
	return fmt.Sprintf(`
resource mikrotik_dhcp_server_network "testacc" {
	address    = %q
	netmask    = %q
	gateway    = %q
	dns_server = %q
	comment    = %q
}
`, address, netmask, gateway, dns_server, comment)
}

func testAccCheckDhcpServerNetworkDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dhcp_server_network" {
			continue
		}

		remoteRecord, err := c.FindDhcpServerNetwork(rs.Primary.ID)

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if remoteRecord != nil {
			return fmt.Errorf("remote recrod (%s) still exists", remoteRecord.Id)
		}

	}
	return nil
}

func testAccDhcpServerNetworkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("%s does not exist in the statefile", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())
		record, err := c.FindDhcpServerNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Unable to get remote record for %s: %v", resourceName, err)
		}

		if record == nil {
			return fmt.Errorf("Unable to get the remote record %s", resourceName)
		}

		return nil
	}
}
