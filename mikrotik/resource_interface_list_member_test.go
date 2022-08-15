package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestInterfaceListMember_basic(t *testing.T) {
	resourceName := "mikrotik_interface_list_member.list_member"

	listName1 := "interface_list1"
	listName2 := "interface_list2"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckInterfaceListMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceListMember(listName1, listName2, "list1", "*0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceListMemberExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "list", listName1),
					resource.TestCheckResourceAttr(resourceName, "interface", "*0"),
				),
			},
			{
				Config: testAccInterfaceListMember(listName1, listName2, "list2", "*0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceListMemberExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "list", listName2),
					resource.TestCheckResourceAttr(resourceName, "interface", "*0"),
				),
			},
		},
	})
}

func testAccInterfaceListMemberExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("%s does not exist in the statefile", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())
		record, err := c.FindInterfaceListMember(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Unable to get remote record for %s: %v", resourceName, err)
		}

		if record == nil {
			return fmt.Errorf("Unable to get the remote record %s", resourceName)
		}

		return nil
	}
}

func testAccCheckInterfaceListMemberDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_list_member" {
			continue
		}

		remoteRecord, err := c.FindInterfaceListMember(rs.Primary.ID)
		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if remoteRecord != nil {
			return fmt.Errorf("remote record (%s) still exists", remoteRecord.Id)
		}
	}

	return nil
}

func testAccInterfaceListMember(listName1, listName2, listToUse string, iface string) string {
	return fmt.Sprintf(`
		resource mikrotik_interface_list "list1" {
			name = %q
		}

		resource mikrotik_interface_list "list2" {
			name = %q
		}

		resource mikrotik_interface_list_member "list_member" {
			interface = %q
			list      = mikrotik_interface_list.%s.name
		}
	`, listName1, listName2, iface, listToUse)
}
