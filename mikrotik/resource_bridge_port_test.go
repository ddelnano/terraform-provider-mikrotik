package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestBridgePort_basic(t *testing.T) {
	rStatePath := "mikrotik_bridge_port.testacc"
	bridgeName := "testacc_bridge"
	bridgeInterface := "*0"
	remoteBridgePort := client.BridgePort{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBridgePortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBridgePortConfig(bridgeName, bridgeInterface, 1, "acceptance test bridge port"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBridgePortExists(rStatePath, &remoteBridgePort),
					resource.TestCheckResourceAttrSet(rStatePath, "id"),
					resource.TestCheckResourceAttr(rStatePath, "bridge", bridgeName),
					resource.TestCheckResourceAttr(rStatePath, "interface", bridgeInterface),
					resource.TestCheckResourceAttr(rStatePath, "pvid", "1"),
					resource.TestCheckResourceAttr(rStatePath, "comment", "acceptance test bridge port"),
				),
			},
			{
				Config: testAccBridgePortConfig(bridgeName+"_updated", bridgeInterface, 2, "updated resource"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBridgePortExists(rStatePath, &remoteBridgePort),
					resource.TestCheckResourceAttrSet(rStatePath, "id"),
					resource.TestCheckResourceAttr(rStatePath, "bridge", bridgeName+"_updated"),
					resource.TestCheckResourceAttr(rStatePath, "interface", bridgeInterface),
					resource.TestCheckResourceAttr(rStatePath, "pvid", "2"),
					resource.TestCheckResourceAttr(rStatePath, "comment", "updated resource"),
				),
			},
		},
	})
}

func testAccBridgePortExists(resourceID string, record *client.BridgePort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceID)
		}
		if r.Primary.ID == "" {
			return fmt.Errorf("resource %q has empty primary ID in state", resourceID)
		}
		c := client.NewClient(client.GetConfigFromEnv())
		remoteRecord, err := c.FindBridgePort(r.Primary.ID)
		if err != nil {
			return err
		}
		*record = *remoteRecord

		return nil
	}
}

func testAccCheckBridgePortDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bridge_port" {
			continue
		}

		remoteRecord, err := c.FindBridgePort(rs.Primary.ID)
		if err != nil && !client.IsNotFoundError(err) {
			return fmt.Errorf("expected not found error, got %+#v", err)
		}

		if remoteRecord != nil {
			return fmt.Errorf("bridge port %q still exists in remote system", remoteRecord.Id)
		}
	}

	return nil
}

func testAccBridgePortConfig(bridgeName, bridgeInterface string, pvid int, comment string) string {
	return fmt.Sprintf(`
		resource mikrotik_bridge "testacc" {
			name = %q
		}

		resource mikrotik_bridge_port "testacc" {
			bridge    = mikrotik_bridge.testacc.name
			interface = %q
			pvid      = %d
			comment   = %q
		}
	`, bridgeName, bridgeInterface, pvid, comment)
}
