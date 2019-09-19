package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var originalDnsName string = "terraform"

// var updatedDnsName string = "terraform.updated"
var originalAddress string = "10.255.255.1"
var updatedAddress string = "10.0.0.1"

func TestAccXenorchestraCloudConfig_create(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckXenorchestraCloudConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConfigConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCloudConfigExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
		},
	})
}

func TestAccXenorchestraCloudConfig_updateAddress(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckXenorchestraCloudConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConfigConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCloudConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", originalAddress),
				),
			},
			{
				Config: testAccCloudConfigConfigUpdatedAddress(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCloudConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedAddress)),
			},
		},
	})
}

// TODO: Add a test for importing the resource

func testAccCloudConfigConfig() string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
}
`, originalDnsName, originalAddress)
}

func testAccCloudConfigConfigUpdatedAddress() string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
}
`, originalDnsName, updatedAddress)
}

func testAccCloudConfigExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudConfig Id is set")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the dns record with error: %v", err)
		}

		if dnsRecord == nil {
			return fmt.Errorf("Unable to get the dns record with name: %s", dnsRecord.Name)
		}

		if dnsRecord.Name == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckXenorchestraCloudConfigDestroyNow(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dns record Id is set")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		if err != nil {
			return err
		}
		err = c.DeleteDnsRecord(dnsRecord.Id)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckXenorchestraCloudConfigDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dns_record" {
			continue
		}

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		if err != nil {
			return err
		}

		if dnsRecord != nil {
			return fmt.Errorf("dns recrod (%s) still exists", dnsRecord.Id)
		}
	}
	return nil
}
