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

func TestAccMikrotikPool_create(t *testing.T) {
	name := acctest.RandomWithPrefix("pool-create")
	ranges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())

	resourceName := "mikrotik_pool.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(name, ranges),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
				),
			},
		},
	})
}

func TestAccMikrotikPool_createNextPool(t *testing.T) {
	name := acctest.RandomWithPrefix("pool-create")
	nextPoolName := acctest.RandomWithPrefix("next_ip_pool")
	ranges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())

	resourceName := "mikrotik_pool.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolWithNextPool(name, ranges, "", nextPoolName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
					resource.TestCheckResourceAttr(resourceName, "next_pool", ""),
				),
			}, {
				Config: testAccPoolWithNextPool(name, ranges, nextPoolName, nextPoolName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
					resource.TestCheckResourceAttr(resourceName, "next_pool", nextPoolName),
				),
			},
			{
				Config: testAccPoolWithNextPool(name, ranges, "", "next_ip_pool"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
					resource.TestCheckResourceAttr(resourceName, "next_pool", ""),
				),
			},
		},
	})
}

func TestAccMikrotikPool_createAndPlanWithNonExistantPool(t *testing.T) {
	name := acctest.RandomWithPrefix("pool-plan")
	ranges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())

	resourceName := "mikrotik_pool.bar"
	removePool := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		pool, err := c.FindPoolByName(name)
		if err != nil {
			t.Fatalf("Error finding the pool by name: %s", err)
		}
		err = c.DeletePool(pool.Id)
		if err != nil {
			t.Fatalf("Error removing the pool: %s", err)
		}
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(name, ranges),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removePool,
				Config:             testAccPool(name, ranges),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikPool_updatePool(t *testing.T) {
	name := acctest.RandomWithPrefix("pool-update-1")
	updatedName := acctest.RandomWithPrefix("pool-update-2")
	ranges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())
	updatedRanges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_pool.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(name, ranges),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
				),
			},
			{
				Config: testAccPool(updatedName, ranges),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
				),
			},
			{
				Config: testAccPool(name, updatedRanges),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", updatedRanges),
				),
			},
			{
				Config: testAccPoolWithComment(name, ranges, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ranges", ranges),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccMikrotikPool_import(t *testing.T) {
	name := acctest.RandomWithPrefix("pool-import")
	ranges := fmt.Sprintf("%s,%s", internal.GetNewIpAddrRange(10), internal.GetNewIpAddr())

	resourceName := "mikrotik_pool.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPool(name, ranges),
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

func testAccPool(name, ranges string) string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
}
`, name, ranges)
}

func testAccPoolWithNextPool(name, ranges, nextPoolToUse, nextPoolName string) string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = %q
    ranges = %q
    next_pool = %q
    depends_on = [mikrotik_pool.next_pool]
}

resource "mikrotik_pool" "next_pool" {
    name = %q
    ranges = "10.10.10.10-10.10.10.20"
}
`, name, ranges, nextPoolToUse, nextPoolName)
}

func testAccPoolWithComment(name, ranges, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_pool" "bar" {
    name = "%s"
    ranges = "%s"
    comment = "%s"
}
`, name, ranges, comment)
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

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if pool != nil {
			return fmt.Errorf("pool (%s) still exists", pool.Id)
		}
	}
	return nil
}
