package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMikrotikBgpInstance_create(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-create")
	routerId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(1, 65535)

	resourceName := "mikrotik_bgp_instance.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(name, as, routerId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(as)),
					resource.TestCheckResourceAttr(resourceName, "router_id", routerId),
				),
			},
		},
	})
}

func TestAccMikrotikBgpInstance_createAndPlanWithNonExistantBgpInstance(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-create_with_plan")
	routerId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(1, 65535)

	resourceName := "mikrotik_bgp_instance.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(name, as, routerId),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckResourceDisappears(testAccProvider, resourceBgpInstance(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMikrotikBgpInstance_updateBgpInstance(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-update")
	routerId := internal.GetNewIpAddr()
	updatedRouterId := internal.GetNewIpAddr()
	clusterId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(1, 65535)
	updatedAs := acctest.RandIntRange(1, 65535)
	comment := acctest.RandomWithPrefix("test comment ")
	confederation := 8

	resourceName := "mikrotik_bgp_instance.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(name, as, routerId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(as)),
					resource.TestCheckResourceAttr(resourceName, "router_id", routerId),
				),
			},
			{
				Config: testAccBgpInstanceUpdatedAsAndRouterId(name, updatedAs, updatedRouterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(updatedAs)),
					resource.TestCheckResourceAttr(resourceName, "router_id", updatedRouterId),
				),
			},
			{
				Config: testAccBgpInstanceUpdatedOptionalFields(name, updatedAs, updatedRouterId, comment, clusterId, confederation),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(updatedAs)),
					resource.TestCheckResourceAttr(resourceName, "router_id", updatedRouterId),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterId),
					resource.TestCheckResourceAttr(resourceName, "client_to_client_reflection", "false"),
					resource.TestCheckResourceAttr(resourceName, "confederation", strconv.Itoa(confederation)),
				),
			},
		},
	})
}

func TestAccMikrotikBgpInstance_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-import")
	routerId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(1, 65535)

	resourceName := "mikrotik_bgp_instance.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikBgpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstance(name, as, routerId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceExists(resourceName),
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

func testAccBgpInstance(name string, as int, routerId string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = %d
    router_id = "%s"
}
`, name, as, routerId)
}

func testAccBgpInstanceUpdatedAsAndRouterId(name string, as int, routerId string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = %d
    router_id = "%s"
}
`, name, as, routerId)
}

func testAccBgpInstanceUpdatedOptionalFields(name string, as int, routerId, comment, clusterId string, confederation int) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance" "bar" {
    name = "%s"
    as = %d
    router_id = "%s"
    comment = "%s"
    cluster_id = "%s"
    client_to_client_reflection = false
    confederation = %d
}
`, name, as, routerId, comment, clusterId, confederation)
}

func testAccBgpInstanceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_bgp_instance does not exist in the statefile")
		}

		bgpInstance, err := apiClient.FindBgpInstance(rs.Primary.ID)

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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_instance" {
			continue
		}

		bgpInstance, err := apiClient.FindBgpInstance(rs.Primary.ID)

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
