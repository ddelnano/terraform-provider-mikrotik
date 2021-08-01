package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var originalName string = "test-pool"
var originalRanges string = "172.16.0.1-172.16.0.8,172.16.0.10"
var updatedName string = "test pool updated"
var updatedRanges string = "172.16.0.11-172.16.0.12"
var updatedComment string = "updated"

func TestAccMikrotikPool_create(t *testing.T) {
	resourceName := "mikrotik_pool.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", originalName),
					resource.TestCheckResourceAttr(resourceName, "ranges", originalRanges),
				),
			},
		},
	})
}

func TestAccMikrotikPool_createAndPlanWithNonExistantPool(t *testing.T) {
	resourceName := "mikrotik_pool.bar"
	removePool := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		pool, err := c.FindPoolByName(originalName)
		if err != nil {
			t.Fatalf("Error finding the pool by name: %s", err)
		}
		err = c.DeletePool(pool.Id)
		if err != nil {
			t.Fatalf("Error removing the pool: %s", err)
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removePool,
				Config:             testAccPool(),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikPool_updatePool(t *testing.T) {
	resourceName := "mikrotik_pool.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalName),
					resource.TestCheckResourceAttr(resourceName, "ranges", originalRanges),
				),
			},
			{
				Config: testAccPoolUpdatedName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "ranges", originalRanges),
				),
			},
			{
				Config: testAccPoolUpdatedRanges(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalName),
					resource.TestCheckResourceAttr(resourceName, "ranges", updatedRanges),
				),
			},
			{
				Config: testAccPoolUpdatedComment(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", originalName),
					resource.TestCheckResourceAttr(resourceName, "ranges", originalRanges),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccMikrotikPool_import(t *testing.T) {
	resourceName := "mikrotik_pool.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
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

func testAccPool() string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
}
`, originalName, originalRanges)
}

func testAccPoolUpdatedName() string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
}
`, updatedName, originalRanges)
}

func testAccPoolUpdatedRanges() string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
}
`, originalName, updatedRanges)
}

func testAccPoolUpdatedComment() string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
    comment = "%s"
}
`, originalName, originalRanges, updatedComment)
}

func testAccPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_pool does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		pool, err := c.FindPool(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the pool with error: %v", err)
		}

		if pool == nil {
			return fmt.Errorf("Unable to get the pool")
		}

		if pool.Id == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikPoolDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_pool" {
			continue
		}

		pool, err := c.FindPool(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if pool != nil {
			return fmt.Errorf("pool (%s) still exists", pool.Id)
		}
	}
	return nil
}
