package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestInterfaceList_basic(t *testing.T) {
	resourceName := "mikrotik_interface_list.testacc"
	listName := "custom_list"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInterfaceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceList(listName, "Initial record"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", listName),
				),
			},
			{
				Config: testAccInterfaceList(listName+"_updated", "updated record"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", listName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "comment", "updated record"),
				),
			},
		},
	})
}

func testAccCheckInterfaceListDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_list" {
			continue
		}

		remoteRecord, err := c.FindInterfaceList(rs.Primary.ID)
		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if remoteRecord != nil {
			return fmt.Errorf("remote record (%s) still exists", remoteRecord.Id)
		}
	}

	return nil
}

func testAccInterfaceListExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("%s does not exist in the statefile", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())
		record, err := c.FindInterfaceList(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Unable to get remote record for %s: %v", resourceName, err)
		}

		if record == nil {
			return fmt.Errorf("Unable to get the remote record %s", resourceName)
		}

		return nil
	}
}

func testAccInterfaceList(name, comment string) string {
	return fmt.Sprintf(`
		resource "mikrotik_interface_list" "testacc" {
			name    = %q
			comment = %q
		}
`, name, comment)
}
