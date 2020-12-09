package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var originalBgpName string = "test-bgp-instance"
var originalConfederation string = "8"
var originalAs string = "65532"
var updatedAs string = "65533"
var originalRouterId string = "172.21.16.1"
var originalClusterId string = "172.21.17.1"
var updatedRouterId string = "172.21.16.2"
var commentBgpInstance string = "test-comment"

func TestAccMikrotikBgpInstance_create(t *testing.T) {
	resourceName := "mikrotik_bgp_instance.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", originalBgpName),
					resource.TestCheckResourceAttr(resourceName, "as", originalAs),
					resource.TestCheckResourceAttr(resourceName, "router_id", originalRouterId),
				),
			},
		},
	})
}

func TestAccMikrotikBgpInstance_createAndPlanWithNonExistantBgpInstance(t *testing.T) {
	resourceName := "mikrotik_bgp_instance.bar"
	removeBgpInstance := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		bgpInstance, err := c.FindBgpInstance(originalBgpName)
		if err != nil {
			t.Fatalf("Error finding the bgp instance by name: %s", err)
		}
		err = c.DeleteBgpInstance(bgpInstance.Name)
		if err != nil {
			t.Fatalf("Error removing the bgp instance: %s", err)
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removeBgpInstance,
				Config:             testAccBgpInstance(),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikBgpInstance_updateBgpInstance(t *testing.T) {
	resourceName := "mikrotik_bgp_instance.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalBgpName),
					resource.TestCheckResourceAttr(resourceName, "as", originalAs),
					resource.TestCheckResourceAttr(resourceName, "router_id", originalRouterId),
				),
			},
			{
				Config: testAccBgpInstanceUpdatedAsAndRouterId(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalBgpName),
					resource.TestCheckResourceAttr(resourceName, "as", updatedAs),
					resource.TestCheckResourceAttr(resourceName, "router_id", updatedRouterId),
				),
			},
			{
				Config: testAccBgpInstanceUpdatedOptionalFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalBgpName),
					resource.TestCheckResourceAttr(resourceName, "as", updatedAs),
					resource.TestCheckResourceAttr(resourceName, "router_id", updatedRouterId),
					resource.TestCheckResourceAttr(resourceName, "comment", commentBgpInstance),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", originalClusterId),
					resource.TestCheckResourceAttr(resourceName, "client_to_client_reflection", "false"),
					resource.TestCheckResourceAttr(resourceName, "confederation", originalConfederation),
				),
			},
		},
	})
}

func TestAccMikrotikBgpInstance_import(t *testing.T) {
	resourceName := "mikrotik_bgp_instance.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			// TODO:  figure out why this fails
			{
				ResourceName: resourceName,
				// tried adding this field, but didn't help
				ImportStateId:     originalBgpName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBgpInstance() string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = 65532
    router_id = "%s"
}
`, originalBgpName, originalRouterId)
}

func testAccBgpInstanceUpdatedAsAndRouterId() string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = 65533
    router_id = "%s"
}
`, originalBgpName, updatedRouterId)
}

func testAccBgpInstanceUpdatedOptionalFields() string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = 65533
    router_id = "%s"
    comment = "%s"
    cluster_id = "%s"
    client_to_client_reflection = false
    confederation = 8
}
`, originalBgpName, updatedRouterId, commentBgpInstance, originalClusterId)
}

func testAccBgpInstanceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_bgp_instance does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		bgpInstance, err := c.FindBgpInstance(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the bgp instance with error: %v", err)
		}

		if bgpInstance == nil {
			return fmt.Errorf("Unable to get the bgp instance")
		}

		if bgpInstance.Name == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikBgpInstanceDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_instance" {
			continue
		}

		bgpInstance, err := c.FindBgpInstance(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if bgpInstance != nil {
			return fmt.Errorf("bgp instance (%s) still exists", bgpInstance.Name)
		}
	}
	return nil
}
