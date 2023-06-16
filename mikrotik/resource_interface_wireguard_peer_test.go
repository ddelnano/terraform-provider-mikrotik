package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var origCommentPeer string = "testing"
var origAllowedAddress string = "192.168.8.1/32"
var origEndpointPort int = 13231
var updatedCommentPeer string = "new_comment"

func TestAccMikrotikInterfaceWireguardPeer_create(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)

	interfaceName := "tf-acc-interface-wireguard"
	publicKey := "/yZWgiYAgNNSy7AIcxuEewYwOVPqJJRKG90s9ypwfiM="
	resourceName := "mikrotik_interface_wireguard_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguardPeer(interfaceName, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "allowed_address", origAllowedAddress),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "endpoint_port", strconv.Itoa(origEndpointPort)),
					resource.TestCheckResourceAttr(resourceName, "interface", interfaceName)),
			},
		},
	})
}

func TestAccMikrotikInterfaceWireguardPeer_updatedComment(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)

	interfaceName := "tf-acc-interface-wireguard"
	publicKey := "/bTmUihbgNsSy2AIcxuEcwYwOVdqJJRKG51s4ypwfiM="
	interfaceNameUpdated := "tf-acc-interface-wireguard-updated"
	resourceName := "mikrotik_interface_wireguard_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguardPeer(interfaceName, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "interface", interfaceName),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "comment", origCommentPeer)),
			},
			{
				Config: testAccInterfaceWireguardPeerUpdatedComment(interfaceName, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "interface", interfaceNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedCommentPeer)),
			},
		},
	})
}

func TestAccMikrotikInterfaceWireguardPeer_import(t *testing.T) {
	client.SkipInterfaceWireguardIfUnsupported(t)

	interfaceName := "tf-acc-interface-wireguard-import"
	publicKey := "/zYaGiYbgNsSy8AIcxuEcwYwOVdqJJRKG91s9ypwfiM="
	resourceName := "mikrotik_interface_wireguard_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceWireguardPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceWireguardPeer(interfaceName, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceWireguardPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "interface"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return interfaceName, nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInterfaceWireguardPeer(interfaceName string, publicKey string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_wireguard" "bar" {
		name = "%s"
		comment = "test interface"
		listen_port = "12321"
		mtu = "1420"
	}
	resource "mikrotik_interface_wireguard_peer" "bar" {
		interface = mikrotik_interface_wireguard.bar.name
		public_key = "%s"
		comment = "%s"
		allowed_address = "%s"
		endpoint_port = "%d"
	}
	`, interfaceName, publicKey, origCommentPeer, origAllowedAddress, origEndpointPort)
}

func testAccInterfaceWireguardPeerUpdatedComment(interfaceName string, publicKey string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_wireguard" "bar" {
		name = "%s"
		comment = "test interface"
		listen_port = "12321"
		mtu = "1420"
	}
	resource "mikrotik_interface_wireguard_peer" "bar" {
		interface = mikrotik_interface_wireguard.bar.name
		public_key = "%s"
		comment = "%s"
		allowed_address = "%s"
		endpoint_port = "%d"
	}
	`, interfaceName, publicKey, updatedCommentPeer, origAllowedAddress, origEndpointPort)
}

func testAccCheckMikrotikInterfaceWireguardPeerDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_wireguard_peer" {
			continue
		}

		interfaceWireguardPeer, err := c.FindInterfaceWireguardPeer(rs.Primary.Attributes["interface"])

		if !client.IsNotFoundError(err) && err != nil {
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
