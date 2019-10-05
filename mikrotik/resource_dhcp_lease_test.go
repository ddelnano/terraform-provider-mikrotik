package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)


var originalIpAddress string = "1.1.1.1"
var originalMacAddress string = "11:11:11:11:11:11"
var updatedIpAddress string = "2.2.2.2"
var updatedMacAddress string = "22:22:22:22:22:22"

func TestAccMikrotikDhcpLease_create(t *testing.T) {
	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "address", originalIpAddress),
					resource.TestCheckResourceAttr(resourceName, "macaddress", originalMacAddress),
				),
			},
		},
	})
}

func TestAccMikrotikDhcpLease_updateAddress(t *testing.T) {
	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", originalIpAddress),
					resource.TestCheckResourceAttr(resourceName, "macaddress", originalMacAddress),
				),
			},
			{
				Config: testAccDhcpLeaseUpdatedIpAddress(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedIpAddress),
					resource.TestCheckResourceAttr(resourceName, "macaddress", originalMacAddress),
				),
			},
			{
				Config: testAccDhcpLeaseUpdatedMacAddress(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", originalIpAddress),
					resource.TestCheckResourceAttr(resourceName, "macaddress", updatedMacAddress),
				),
			},
		},
	})
}

func TestAccMikrotikDhcpLease_import(t *testing.T) {
	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
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

func testAccDhcpLease() string {
	return fmt.Sprintf(`
resource "mikrotik_dhcp_lease" "bar" {
    comment = "bar"
    address = "%s"
    macaddress = "%s"
}
`, originalIpAddress, originalMacAddress)
}

func testAccDhcpLeaseUpdatedIpAddress() string {
	return fmt.Sprintf(`
resource "mikrotik_dhcp_lease" "bar" {
    comment = "bar"
    address = "%s"
    macaddress = "%s"
}
`, updatedIpAddress, originalMacAddress)
}

func testAccDhcpLeaseUpdatedMacAddress() string {
	return fmt.Sprintf(`
resource "mikrotik_dhcp_lease" "bar" {
    comment = "bar"
    address = "%s"
    macaddress = "%s"
}
`, originalIpAddress, updatedMacAddress)
}

func testAccDhcpLeaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_dhcp_lease does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dhcpLease, err := c.FindDhcpLease(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the dhcp lease with error: %v", err)
		}

		if dhcpLease == nil {
			return fmt.Errorf("Unable to get the dhcp lease")
		}

		if dhcpLease.Id == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikDhcpLeaseDestroyNow(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dhcp lease Id is set")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dhcpLease, err := c.FindDhcpLease(rs.Primary.ID)

		if err != nil {
			return err
		}

		err = c.DeleteDhcpLease(dhcpLease.Id)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckMikrotikDhcpLeaseDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dhcp_lease" {
			continue
		}

		dhcpLease, err := c.FindDhcpLease(rs.Primary.ID)

		if err != nil {
			return err
		}

		if dhcpLease != nil {
			return fmt.Errorf("dhcp lease (%s) still exists", dhcpLease.Id)
		}
	}
	return nil
}
