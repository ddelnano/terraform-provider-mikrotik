package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var terraformResourceTypeFirewallFilterRule string = "mikrotik_firewall_filter_rule"

func TestFirewallFilterRule_basic(t *testing.T) {

	resourceName := terraformResourceTypeFirewallFilterRule + ".testacc"

	action := "accept"
	chain := "testChain"
	connectionState := []string{"new"}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFirewallFilterRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallFilterRuleConfigBasic(action, chain, connectionState, "80", "tcp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "chain", chain),
					resource.TestCheckResourceAttr(resourceName, "dst_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
				),
			},
			{
				Config: testAccFirewallFilterRuleConfigBasic(action, chain, connectionState, "68", "udp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "chain", chain),
					resource.TestCheckResourceAttr(resourceName, "dst_port", "68"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
				),
			},
		},
	})
}

func testAccCheckFirewallFilterRuleDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != terraformResourceTypeFirewallFilterRule {
			continue
		}

		remoteRecord, err := c.FindFirewallFilterRule(rs.Primary.ID)
		_, ok := err.(*client.NotFound)
		if err != nil && !ok {
			return fmt.Errorf("expected not found error, got %+#v", err)
		}

		if remoteRecord != nil {
			return fmt.Errorf("resource %T with id %q still exists in remote system", remoteRecord, remoteRecord.Id)
		}
	}
	return nil
}

func testAccFirewallFilterRuleConfigBasic(action, chain string, connectionState []string, destPort, protocol string) string {
	return fmt.Sprintf(`
		resource "mikrotik_firewall_filter_rule" "testacc" {
			action             = %q
			chain              = %q
			connection_state   = [%s]
			dst_port           = %q
			protocol           = %q
		}
	`, action, chain, internal.JoinStringsToString(connectionState, ","), destPort, protocol)
}
