package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDhcpServer_basic(t *testing.T) {
	rName := "dhcp-server"
	rLeaseScript := ":put 123"
	dhcpServer := client.DhcpServer{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDhcpServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpServerConfig(rName, true, rLeaseScript),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpServerResourceExists("mikrotik_dhcp_server.testacc", &dhcpServer),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "name", rName),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "disabled", "true"),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "lease_script", rLeaseScript),
				),
			},
			{
				Config: testAccDhcpServerConfig(rName, false, ":put updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpServerResourceExists("mikrotik_dhcp_server.testacc", &dhcpServer),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "name", rName),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "disabled", "false"),
					resource.TestCheckResourceAttr("mikrotik_dhcp_server.testacc", "lease_script", ":put updated"),
				),
			},
		},
	})
}

func testAccDhcpServerResourceExists(resource string, record *client.DhcpServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		r, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resource)
		}
		if r.Primary.ID == "" {
			return fmt.Errorf("resource %q has empty primary ID in state", resource)
		}
		c := client.NewClient(client.GetConfigFromEnv())
		dhcpServer, err := c.FindDhcpServer(r.Primary.ID)
		if err != nil {
			return err
		}
		*record = *dhcpServer

		return nil
	}
}

func testAccCheckDhcpServerDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dhcp_server" {
			continue
		}

		dhcpServer, err := c.FindDhcpServer(rs.Primary.ID)
		_, ok := err.(*client.NotFound)
		if err != nil && !ok {
			return fmt.Errorf("expected not found error, got %+#v", err)
		}

		if dhcpServer != nil {
			return fmt.Errorf("dhcp-server %q (%s) still exists in remote system", dhcpServer.Name, dhcpServer.Id)
		}
	}

	return nil
}

func testAccDhcpServerConfig(name string, disabled bool, leaseScript string) string {
	return fmt.Sprintf(`
		resource "mikrotik_dhcp_server" "testacc" {
			name = %q
			disabled = %t
			lease_script = %q
		}
	`, name, disabled, leaseScript)
}
