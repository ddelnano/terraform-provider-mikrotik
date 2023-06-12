package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBridge_basic(t *testing.T) {
	rName := "testacc_bridge"
	bridge := client.Bridge{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBridgeConfig(rName, true, false, "testacc bridge"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBridgeExists("mikrotik_bridge.testacc", &bridge),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "name", rName),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "fast_forward", "true"),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "vlan_filtering", "false"),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "comment", "testacc bridge"),
				),
			},
			{
				Config: testAccBridgeConfig(rName+"_updated", false, true, "updated bridge"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBridgeExists("mikrotik_bridge.testacc", &bridge),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "name", rName+"_updated"),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "fast_forward", "false"),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "vlan_filtering", "true"),
					resource.TestCheckResourceAttr("mikrotik_bridge.testacc", "comment", "updated bridge"),
				),
			},
		},
	})
}

func testAccBridgeExists(resource string, record *client.Bridge) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resource)
		}
		if r.Primary.ID == "" {
			return fmt.Errorf("resource %q has empty primary ID in state", resource)
		}
		c := client.NewClient(client.GetConfigFromEnv())
		remoteRecord, err := c.FindBridge(r.Primary.ID)
		if err != nil {
			return err
		}
		*record = *remoteRecord

		return nil
	}
}

func testAccBridgeDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bridge" {
			continue
		}

		remoteRecord, err := c.FindBridge(rs.Primary.ID)
		if err != nil && !client.IsNotFoundError(err) {
			return fmt.Errorf("expected not found error, got %+#v", err)
		}

		if remoteRecord != nil {
			return fmt.Errorf("bridge %q (%s) still exists in remote system", remoteRecord.Name, remoteRecord.Id)
		}
	}

	return nil
}

func testAccBridgeConfig(name string, fastForward, vlanFiltering bool, comment string) string {
	return fmt.Sprintf(`
		resource "mikrotik_bridge" "testacc" {
			name = %q
			fast_forward = %t
			vlan_filtering = %t
			comment = %q
		}
	`, name, fastForward, vlanFiltering, comment)
}
