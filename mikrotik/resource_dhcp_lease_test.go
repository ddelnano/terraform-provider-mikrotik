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

func TestAccMikrotikDhcpLease_create(t *testing.T) {
	ipAddr := internal.GetNewIpAddr()
	macAddr := internal.GetNewMacAddr()
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(ipAddr, macAddr, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", macAddr),
					resource.TestCheckResourceAttr(resourceName, "dynamic", "false"),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccMikrotikDhcpLease_updateLease(t *testing.T) {
	ipAddr := internal.GetNewIpAddr()
	updatedIpAddr := internal.GetNewIpAddr()
	macAddr := internal.GetNewMacAddr()
	updatedMacAddr := internal.GetNewMacAddr()
	comment := acctest.RandomWithPrefix("tf-acc-comment")
	updatedComment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(ipAddr, macAddr, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", macAddr),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccDhcpLease(updatedIpAddr, macAddr, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedIpAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", macAddr),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccDhcpLease(ipAddr, macAddr, updatedComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", macAddr),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
			{
				Config: testAccDhcpLease(ipAddr, updatedMacAddr, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", updatedMacAddr),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccDhcpLeaseUpdatedBlockAccess(ipAddr, macAddr, comment, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDhcpLeaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
					resource.TestCheckResourceAttr(resourceName, "macaddress", macAddr),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "blocked", "true"),
				),
			},
		},
	})
}

func TestAccMikrotikDhcpLease_import(t *testing.T) {
	ipAddr := internal.GetNewIpAddr()
	macAddr := internal.GetNewMacAddr()
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_dhcp_lease.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDhcpLeaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpLease(ipAddr, macAddr, comment),
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

func testAccDhcpLease(ipAddr, macAddr, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_dhcp_lease" "bar" {
    address = "%s"
    macaddress = "%s"
    comment = "%s"
}
`, ipAddr, macAddr, comment)
}

func testAccDhcpLeaseUpdatedBlockAccess(ipAddr, macAddr, comment string, blocked bool) string {
	return fmt.Sprintf(`
resource "mikrotik_dhcp_lease" "bar" {
    address = "%s"
    macaddress = "%s"
    blocked = "%t"
    comment = "%s"
}
`, ipAddr, macAddr, blocked, comment)
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

func testAccCheckMikrotikDhcpLeaseDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dhcp_lease" {
			continue
		}

		dhcpLease, err := c.FindDhcpLease(rs.Primary.ID)

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if dhcpLease != nil {
			return fmt.Errorf("dhcp lease (%s) still exists", dhcpLease.Id)
		}
	}
	return nil
}
