package mikrotik

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var origComment string = "testing"
var origListenPort int = 13231
var origMTU int = 1420
var updatedComment string = "new_comment"

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
					resource.TestCheckResourceAttr(resourceName, "comment", origComment),
					resource.TestCheckResourceAttr(resourceName, "listen_port", strconv.Itoa(origListenPort)),
					resource.TestCheckResourceAttr(resourceName, "mtu", strconv.Itoa(origMTU)),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false")),
			},
		},
	})
}

func TestAccMikrotikInterfaceWireguard_updatedComment(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-update-comment")

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
					resource.TestCheckResourceAttr(resourceName, "comment", origComment)),
			},
			{
				Config: testAccInterfaceWireguardUpdatedComment(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment)),
			},
		},
	})
}

func TestAccMikrotikInterfaceWireguard_import(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-import")

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
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return name, nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInterfaceWireguard(name string) string {
	return fmt.Sprintf(`
resource "mikrotik_interface_wireguard" "bar" {
    name = "%s"
	comment = "%s"
	listen_port = "%d"
	mtu = "%d"
}
`, name, origComment, origListenPort, origMTU)
}

func testAccInterfaceWireguardUpdatedComment(name string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_wireguard" "bar" {
		name = "%s"
		comment = "%s"
		listen_port = "%d"
		mtu = "%d"
	}
	`, name, updatedComment, origListenPort, origMTU)
}

func testAccCheckMikrotikInterfaceWireguardDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_wireguard" {
			continue
		}

		interfaceWireguard, err := c.FindInterfaceWireguard(rs.Primary.Attributes["name"])

		log.Printf("err type:  %T", err)
		log.Printf("err:  %v", err)
		if !client.IsNotFoundError(err) && err != nil {
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

		interfaceWireguard, err := c.FindInterfaceWireguard(rs.Primary.Attributes["name"])

		if err != nil {
			return fmt.Errorf("Unable to get the interface wireguard with error: %v", err)
		}

		if interfaceWireguard == nil {
			return fmt.Errorf("Unable to get the interface wireguard with name: %s", rs.Primary.Attributes["name"])
		}

		if interfaceWireguard.Name == rs.Primary.Attributes["name"] {
			return nil
		}
		return nil
	}
}
