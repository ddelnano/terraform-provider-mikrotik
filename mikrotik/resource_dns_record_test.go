package mikrotik

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMikrotikDnsRecord_create(t *testing.T) {
	dnsName := internal.GetNewDnsName()
	ipAddr := internal.GetNewIpAddr()

	resourceName := "mikrotik_dns_record.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(dnsName, ipAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				Config: `
					resource "mikrotik_dns_record" "bar" {
						address = "10.10.200.100"
						regexp  = ".+\\.domain\\.com"
						ttl     = "300"
					}
				`,
				ExpectError: regexp.MustCompile("only name or regexp allowed"),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_createRegexp(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "mikrotik_dns_record" "bar" {
						address = "10.10.200.100"
						regexp  = ".+\\.domain\\.com"
						ttl     = "300"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr("mikrotik_dns_record.bar", "regexp", ".+\\.domain\\.com"),
				),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_createAndPlanWithNonExistantRecord(t *testing.T) {
	dnsName := internal.GetNewDnsName()
	ipAddr := internal.GetNewIpAddr()

	resourceName := "mikrotik_dns_record.bar"
	removeDnsRecord := func() {
		c := client.NewClient(client.GetConfigFromEnv())
		dns, err := c.FindDnsRecord(dnsName)

		if err != nil {
			t.Fatalf("Error finding the DNS record: %s", err)
		}
		err = c.DeleteDnsRecord(dns.Id)
		if err != nil {
			t.Fatalf("Error removing the DNS record: %s", err)
		}

	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(dnsName, ipAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removeDnsRecord,
				Config:             testAccDnsRecord(dnsName, ipAddr),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikDnsRecord_updateAddress(t *testing.T) {
	dnsName := internal.GetNewDnsName()
	ipAddr := internal.GetNewIpAddr()
	updatedIpAddr := internal.GetNewIpAddr()

	resourceName := "mikrotik_dns_record.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(dnsName, ipAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", ipAddr),
				),
			},
			{
				Config: testAccDnsRecord(dnsName, updatedIpAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedIpAddr)),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_updateComment(t *testing.T) {
	dnsName := internal.GetNewDnsName()
	ipAddr := internal.GetNewIpAddr()
	comment := "Initial comment"
	updatedComment := "new comment"

	resourceName := "mikrotik_dns_record.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordWithComment(dnsName, ipAddr, comment),

				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccDnsRecordWithComment(dnsName, ipAddr, updatedComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment)),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_import(t *testing.T) {
	dnsName := internal.GetNewDnsName()
	ipAddr := internal.GetNewIpAddr()

	resourceName := "mikrotik_dns_record.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(dnsName, ipAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				ImportState:       true,
				ResourceName:      resourceName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsRecord(dnsName, ipAddr string) string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
}
`, dnsName, ipAddr)
}

func testAccDnsRecordWithComment(dnsName, ipAddr, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
    comment = "%s"
}
`, dnsName, ipAddr, comment)
}

func testAccDnsRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_dns_record does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dnsRecord, err := c.FindDnsRecord(rs.Primary.Attributes["name"])

		if err != nil {
			return fmt.Errorf("Unable to get the dns record with error: %v", err)
		}

		if dnsRecord == nil {
			return fmt.Errorf("Unable to get the dns record with name: %s", rs.Primary.Attributes["name"])
		}

		if dnsRecord.Name == rs.Primary.Attributes["name"] {
			return nil
		}

		return nil
	}
}

func testAccCheckMikrotikDnsRecordDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dns_record" {
			continue
		}

		dnsRecord, err := c.FindDnsRecord(rs.Primary.Attributes["name"])

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if dnsRecord != nil {
			return fmt.Errorf("dns recrod (%s) still exists", dnsRecord.Id)
		}
	}
	return nil
}
