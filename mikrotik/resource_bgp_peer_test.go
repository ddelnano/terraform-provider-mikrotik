package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var bgpPeerName string = "test-peer"
var remoteAs string = "65533"
var remoteAddress string = "172.21.16.0"
var instanceName string = "test"
var peerTTL string = "default"
var addressFamilies string = "ip"
var defaultOriginate string = "never"
var holdTime string = "3m"
var nextHopChoice string = "default"
var commentBgpPeer string = "test-comment"

var maxPrefixRestartTime string = "1w3d"
var tcpMd5Key string = "test-tcp-md5-key"
var updatedTTL string = "255"
var updatedUseBfd string = "true"

func TestAccMikrotikBgpPeer_create(t *testing.T) {
	resourceName := "mikrotik_bgp_peer.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", bgpPeerName),
					resource.TestCheckResourceAttr(resourceName, "remote_as", remoteAs),
					resource.TestCheckResourceAttr(resourceName, "instance", instanceName),
					resource.TestCheckResourceAttr(resourceName, "ttl", peerTTL),
					resource.TestCheckResourceAttr(resourceName, "address_families", addressFamilies),
					resource.TestCheckResourceAttr(resourceName, "default_originate", defaultOriginate),
					resource.TestCheckResourceAttr(resourceName, "hold_time", holdTime),
					resource.TestCheckResourceAttr(resourceName, "nexthop_choice", nextHopChoice),
					resource.TestCheckResourceAttr(resourceName, "as_override", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "passive", "false"),
					resource.TestCheckResourceAttr(resourceName, "remove_private_as", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_reflect", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", "false"),
				),
			},
		},
	})
}

func TestAccMikrotikBgpPeer_createAndPlanWithNonExistantBgpPeer(t *testing.T) {
	resourceName := "mikrotik_bgp_peer.bar"
	removeBgpPeer := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		bgpPeer, err := c.FindBgpPeer(bgpPeerName)
		if err != nil {
			t.Fatalf("Error finding the bgp peer by name: %s", err)
		}
		err = c.DeleteBgpPeer(bgpPeer.Name)
		if err != nil {
			t.Fatalf("Error removing the bgp peer: %s", err)
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removeBgpPeer,
				Config:             testAccBgpPeer(),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikBgpPeer_updateBgpPeer(t *testing.T) {
	resourceName := "mikrotik_bgp_peer.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", bgpPeerName),
					resource.TestCheckResourceAttr(resourceName, "remote_as", remoteAs),
					resource.TestCheckResourceAttr(resourceName, "instance", instanceName),
					resource.TestCheckResourceAttr(resourceName, "ttl", peerTTL),
					resource.TestCheckResourceAttr(resourceName, "address_families", addressFamilies),
					resource.TestCheckResourceAttr(resourceName, "default_originate", defaultOriginate),
					resource.TestCheckResourceAttr(resourceName, "hold_time", holdTime),
					resource.TestCheckResourceAttr(resourceName, "nexthop_choice", nextHopChoice),
					resource.TestCheckResourceAttr(resourceName, "as_override", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "passive", "false"),
					resource.TestCheckResourceAttr(resourceName, "remove_private_as", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_reflect", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", "false"),
				),
			},
			{
				Config: testAccBgpPeerUpdatedUseBfdTCPMd5KeyTTLAndMaxPrefixRestartTime(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", bgpPeerName),
					resource.TestCheckResourceAttr(resourceName, "remote_as", remoteAs),
					resource.TestCheckResourceAttr(resourceName, "instance", instanceName),
					resource.TestCheckResourceAttr(resourceName, "address_families", addressFamilies),
					resource.TestCheckResourceAttr(resourceName, "default_originate", defaultOriginate),
					resource.TestCheckResourceAttr(resourceName, "hold_time", holdTime),
					resource.TestCheckResourceAttr(resourceName, "nexthop_choice", nextHopChoice),
					resource.TestCheckResourceAttr(resourceName, "as_override", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "passive", "false"),
					resource.TestCheckResourceAttr(resourceName, "remove_private_as", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_reflect", "false"),
					resource.TestCheckResourceAttr(resourceName, "ttl", updatedTTL),
					resource.TestCheckResourceAttr(resourceName, "max_prefix_restart_time", maxPrefixRestartTime),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", updatedUseBfd),
					resource.TestCheckResourceAttr(resourceName, "tcp_md5_key", tcpMd5Key),
				),
			},
		},
	})
}

func TestAccMikrotikBgpPeer_import(t *testing.T) {
	resourceName := "mikrotik_bgp_peer.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBgpPeer() string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_peer" "bar" {
    name = "%s"
    remote_as = 65533
    remote_address = "%s"
    instance = "%s"
    ttl = "%s"
    address_families = "%s"
    default_originate = "%s"
    hold_time = "%s"
    nexthop_choice = "%s"
}
`, bgpPeerName, remoteAddress, instanceName, peerTTL, addressFamilies, defaultOriginate, holdTime, nextHopChoice)
}

func testAccBgpPeerUpdatedUseBfdTCPMd5KeyTTLAndMaxPrefixRestartTime() string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_peer" "bar" {
    name = "%s"
    remote_as = 65533
    remote_address = "%s"
    instance = "%s"
    ttl = "%s"
    address_families = "%s"
    default_originate = "%s"
    hold_time = "%s"
    nexthop_choice = "%s"

    max_prefix_restart_time = "%s"
    use_bfd = true
    tcp_md5_key = "%s"
}
`, bgpPeerName, remoteAddress, instanceName, updatedTTL, addressFamilies, defaultOriginate, holdTime, nextHopChoice, maxPrefixRestartTime, tcpMd5Key)
}

func testAccBgpPeerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_bgp_peer does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		bgpPeer, err := c.FindBgpPeer(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the bgp peer with error: %v", err)
		}

		if bgpPeer == nil {
			return fmt.Errorf("Unable to get the bgp peer")
		}

		if bgpPeer.Name == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikBgpPeerDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_peer" {
			continue
		}

		bgpPeer, err := c.FindBgpPeer(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if bgpPeer != nil {
			return fmt.Errorf("bgp peer (%s) still exists", bgpPeer.Name)
		}
	}
	return nil
}
