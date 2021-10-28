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

func TestAccMikrotikResourceIpAddress_create(t *testing.T) {
	ipAddr := internal.GetNewIpAddr() + "/24"
	ifName := "ether1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_ip_address.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikIpAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpAddress(ipAddr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpAddressExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccMikrotikResourceIpAddress_updateAddr(t *testing.T) {
	ipAddr := internal.GetNewIpAddr() + "/24"
	updatedIpAddr := internal.GetNewIpAddr() + "/24"
	ifName := "ether1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")
	updatedComment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_ip_address.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikIpAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpAddress(ipAddr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccIpAddress(updatedIpAddr, ifName, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedIpAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccIpAddress(ipAddr, ifName, updatedComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
			{
				Config: testAccIpAddressUpdatedDisabled(ipAddr, ifName, comment, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccIpAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "interface", ifName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
		},
	})
}

func testAccIpAddress(ipAddr, ifName, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_ip_address" "test" {
	address = "%s"
	interface = "%s"
	comment = "%s"
}
`, ipAddr, ifName, comment)
}

func testAccIpAddressUpdatedDisabled(ipAddr, ifName, comment string, disabled bool) string {
	return fmt.Sprintf(`
resource "mikrotik_ip_address" "test" {
	address = "%s"
	interface = "%s"
	comment = "%s"
	disabled = "%t"
}
`, ipAddr, ifName, comment, disabled)
}

func testAccIpAddressExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_ip_address does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		ipaddr, err := c.FindIpAddress(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the ip address with error: %v", err)
		}

		if ipaddr == nil {
			return fmt.Errorf("Unable to get the ip address")
		}

		if ipaddr.Id == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikIpAddressDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_ip_address" {
			continue
		}

		ipaddr, err := c.FindIpAddress(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if ipaddr != nil {
			return fmt.Errorf("ip address (%s) still exists", ipaddr.Id)
		}
	}
	return nil
}
