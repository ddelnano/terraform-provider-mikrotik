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

var instanceName string = "default"
var peerTTL string = "default"
var addressFamilies string = "ip"
var defaultOriginate string = "never"
var holdTime string = "3m"
var nextHopChoice string = "default"

var maxPrefixRestartTime string = "1w3d"
var tcpMd5Key string = "test-tcp-md5-key"
var updatedTTL string = "255"
var updatedUseBfd string = "true"

func TestAccMikrotikBgpPeer_create(t *testing.T) {
	client.SkipLegacyBgpIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-create")
	remoteAs := acctest.RandIntRange(1, 65535)
	remoteAddress, _ := acctest.RandIpAddress("192.168.0.0/24")

	resourceName := "mikrotik_bgp_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(name, remoteAs, remoteAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "remote_as", strconv.Itoa(remoteAs)),
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
	client.SkipLegacyBgpIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-create_with_plan")
	remoteAs := acctest.RandIntRange(1, 65535)
	remoteAddress, _ := acctest.RandIpAddress("192.168.1.0/24")

	resourceName := "mikrotik_bgp_peer.bar"
	removeBgpPeer := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		bgpPeer, err := c.FindBgpPeer(name)
		if err != nil {
			t.Fatalf("Error finding the bgp peer by name: %s", err)
		}
		err = c.DeleteBgpPeer(bgpPeer.Name)
		if err != nil {
			t.Fatalf("Error removing the bgp peer: %s", err)
		}
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(name, remoteAs, remoteAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removeBgpPeer,
				Config:             testAccBgpPeer(name, remoteAs, remoteAddress),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikBgpPeer_updateBgpPeer(t *testing.T) {
	client.SkipLegacyBgpIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-update")
	remoteAs := acctest.RandIntRange(1, 65535)
	remoteAddress, _ := acctest.RandIpAddress("192.168.3.0/24")

	resourceName := "mikrotik_bgp_peer.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(name, remoteAs, remoteAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "remote_as", strconv.Itoa(remoteAs)),
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
				Config: testAccBgpPeerUpdatedUseBfdTCPMd5KeyTTLAndMaxPrefixRestartTime(name, remoteAs, remoteAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpPeerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "remote_as", strconv.Itoa(remoteAs)),
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
	client.SkipLegacyBgpIfUnsupported(t)
	name := acctest.RandomWithPrefix("tf-acc-import")
	remoteAs := acctest.RandIntRange(1, 65535)
	remoteAddress, _ := acctest.RandIpAddress("192.168.4.0/24")

	resourceName := "mikrotik_bgp_peer.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeer(name, remoteAs, remoteAddress),
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

func testAccBgpPeer(name string, remoteAs int, remoteAddress string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_peer" "bar" {
    name = "%s"
    remote_as = %d
    remote_address = "%s"
    instance = "%s"
    ttl = "%s"
    address_families = "%s"
    default_originate = "%s"
    hold_time = "%s"
    nexthop_choice = "%s"
}
`, name, remoteAs, remoteAddress, instanceName, peerTTL, addressFamilies, defaultOriginate, holdTime, nextHopChoice)
}

func testAccBgpPeerUpdatedUseBfdTCPMd5KeyTTLAndMaxPrefixRestartTime(name string, remoteAs int, remoteAddress string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_peer" "bar" {
    name = "%s"
    remote_as = %d
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
`, name, remoteAs, remoteAddress, instanceName, updatedTTL, addressFamilies, defaultOriginate, holdTime, nextHopChoice, maxPrefixRestartTime, tcpMd5Key)
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

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if bgpPeer != nil {
			return fmt.Errorf("bgp peer (%s) still exists", bgpPeer.Name)
		}
	}
	return nil
}
