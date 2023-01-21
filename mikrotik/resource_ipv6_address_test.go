package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMikrotikResourceIpv6Address_create(t *testing.T) {
	if client.IsLegacyBgpSupported() {
		t.Skip()
	}

	ipv6Addr := internal.GetNewIpv6Addr() + "/64"
	ifName := "ether1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_ipv6_address.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikIpv6AddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpv6Address(ipv6Addr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpv6AddressExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "address", ipv6Addr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccMikrotikResourceIpv6Address_updateAddr(t *testing.T) {
	if client.IsLegacyBgpSupported() {
		t.Skip()
	}

	ipAddr := internal.GetNewIpv6Addr() + "/64"
	updatedIpv6Addr := internal.GetNewIpv6Addr() + "/64"
	ifName := "ether1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")
	disabled := "false"
	updatedComment := acctest.RandomWithPrefix("tf-acc-comment")
	updatedDisabled := "true"

	resourceName := "mikrotik_ipv6_address.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikIpv6AddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpv6Address(ipAddr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpv6AddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccIpv6Address(updatedIpv6Addr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpv6AddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedIpv6Addr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccIpv6Address(ipAddr, ifName, updatedComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpv6AddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccIpv6AddressUpdatedDisabled(ipAddr, ifName, comment, updatedDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpv6AddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", updatedDisabled),
				),
			},
		},
	})
}

func testAccIpv6Address(ipAddr, ifName, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_ipv6_address" "test" {
	address = "%s"
	interface = "%s"
	comment = "%s"
}
`, ipAddr, ifName, comment)
}

func testAccIpv6AddressUpdatedDisabled(ipAddr, ifName, comment string, disabled string) string {
	return fmt.Sprintf(`
resource "mikrotik_ipv6_address" "test" {
	address = "%s"
	interface = "%s"
	comment = "%s"
	disabled = "%s"
}
`, ipAddr, ifName, comment, disabled)
}

func testAccIpv6AddressExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_ipv6_address does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		ipaddr, err := c.FindIpv6Address(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the ipv6 address with error: %v", err)
		}

		if ipaddr == nil {
			return fmt.Errorf("Unable to get the ipv6 address")
		}

		if ipaddr.Id == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikIpv6AddressDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_ipv6_address" {
			continue
		}

		ipaddr, err := c.FindIpv6Address(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if ipaddr != nil {
			return fmt.Errorf("ipv6 address (%s) still exists", ipaddr.Id)
		}
	}
	return nil
}
