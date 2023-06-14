package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var origComment_peer string = "testing"
var origAllowedAddress int = 10
var origEndpointAddress int = 192
var origEndpointPort int = 13231
var origInterface string = "test_interface"
var updatedComment_peer string = "new_comment"

func TestAccMikrotikInterfaceWireguardPeer_create(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)
	//name := acctest.RandomWithPrefix("tf-acc-create")

	resourceName := "mikrotik_interface_wireguard_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguardPeer(origInterface),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, ".id"),
					resource.TestCheckResourceAttr(resourceName, "allowed_address", strconv.Itoa(origAllowedAddress)),
					resource.TestCheckResourceAttr(resourceName, "endpoint_address", strconv.Itoa(origEndpointAddress)),
					resource.TestCheckResourceAttr(resourceName, "endpoint_port", strconv.Itoa(origEndpointPort)),
					resource.TestCheckResourceAttr(resourceName, "interface", origInterface)),
			},
		},
	})
}

func TestAccMikrotikInterfaceWireguardPeer_updatedComment(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-update-comment")

	resourceName := "mikrotik_interface_wireguard_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguardPeer(origInterface), //what parameter should I use here?
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					//resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", origComment_peer)),
			},
			{
				Config: testAccInterfaceWireguardPeerUpdatedComment(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					//resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment_peer)),
			},
		},
	})
}

// func TestAccMikrotikInterfaceWireguardPeer_import(t *testing.T) {
// 	client.SkipInterfaceWireguardIfUnsupported(t)
// 	name := acctest.RandomWithPrefix("tf-acc-import")

// 	resourceName := "mikrotik_interface_wireguard_peer.bar"
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
// 		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccInterfaceWireguardPeer(name),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					testAccInterfaceWireguardPeerExists(resourceName),
// 					resource.TestCheckResourceAttrSet(resourceName, ".id"),
// 					resource.TestCheckResourceAttrSet(resourceName, "allowed_address"),
// 				),
// 			},
// 			{
// 				ResourceName: resourceName,
// 				ImportState:  true,
// 				ImportStateIdFunc: func(s *terraform.State) (string, error) {
// 					return name, nil
// 				},
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

func testAccInterfaceWireguardPeer(tinterface string) string {
	return fmt.Sprintf(`
resource "mikrotik_interface_wireguard_peer" "bar" {
	comment = "%s"
	allowed_address = "%d"
	endpoint_address = "%d"
	endpoint_port = "%d"
	interface = "%s"
}
`, origComment_peer, origAllowedAddress, origEndpointAddress, origEndpointPort, tinterface)
}

func testAccInterfaceWireguardPeerUpdatedComment(tinterface string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_wireguard_peer" "bar" {
		comment = "%s"
		allowed_address = "%d"
		endpoint_address = "%d"
		endpoint_port = "%d"
		interface = "%s"
	}
	`, updatedComment_peer, origAllowedAddress, origEndpointAddress, origEndpointPort, tinterface)
}

func testAccCheckMikrotikInterfaceWireguardPeerDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_wireguard_peer" {
			continue
		}

		interfaceWireguardPeer, err := c.FindInterfaceWireguardPeer(rs.Primary.Attributes["interface"])

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if interfaceWireguardPeer != nil {
			return fmt.Errorf("interface wireguard peer (%s) still exists", interfaceWireguardPeer.Id)
		}
	}
	return nil
}

func testAccInterfaceWireguardPeerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_interface_wireguard_peer does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		interfaceWireguardPeer, err := c.FindInterfaceWireguardPeer(rs.Primary.Attributes["interface"])

		_, ok = err.(*client.NotFound)
		if !ok && err != nil {
			return fmt.Errorf("Unable to get the interface wireguard peer with error: %v", err)
		}

		if interfaceWireguardPeer == nil {
			return fmt.Errorf("Unable to get the interface wireguard peer with id: %s", rs.Primary.Attributes[".id"])
		}

		if interfaceWireguardPeer.Id == rs.Primary.Attributes[".id"] {
			return nil
		}
		return nil
	}
}
