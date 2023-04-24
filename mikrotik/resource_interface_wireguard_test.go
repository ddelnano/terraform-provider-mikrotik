package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var origComment string = "testing"
var origListenPort int = 13231
var origMTU int = 1420

func TestAccMikrotikInterfaceWireguard_create(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-create")

	resourceName := "mikrotik_interface_wireguard.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguard(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "comment")),
				// resource.TestCheckResourceAttrSet(resourceName, "listen_port"),
				// resource.TestCheckResourceAttrSet(resourceName, "mtu"),
				// resource.TestCheckResourceAttrSet(resourceName, "disabled")),
			},
		},
	})
}

func testAccInterfaceWireguard(name string) string {
	return fmt.Sprintf(`
resource "mikrotik_interface_wireguard" "bar" {
    name = "%s"
	comment = "%s"
}
`, name, origComment)
}

func testAccCheckMikrotikInterfaceWireguardDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_wireguard" {
			continue
		}

		interfaceWireguard, err := c.FindInterfaceWireguard(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if interfaceWireguard != nil {
			return fmt.Errorf("interface wireguard (%s) still exists", interfaceWireguard.Name)
		}
	}
	return nil
}

func testAccInterfaceWireguardExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_interface_wireguard does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		interfaceWireguard, err := c.FindInterfaceWireguard(rs.Primary.ID)

		_, ok = err.(*client.NotFound)
		if !ok && err != nil {
			return fmt.Errorf("Unable to get the interface wireguard with error: %v", err)
		}

		if interfaceWireguard == nil {
			return fmt.Errorf("Unable to get the interface wireguard with id: %s", rs.Primary.ID)
		}

		if interfaceWireguard.Id == rs.Primary.ID {
			return nil
		}
		return nil
	}
}
