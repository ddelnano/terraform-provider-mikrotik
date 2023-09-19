package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestVlanInterface_basic(t *testing.T) {
	resourceName := "mikrotik_vlan_interface.testacc"
	iface := "ether1"
	mtu := 1500
	name := "test-vlan"
	useServiceTag := false
	vlanID := 20
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVlanInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanInterface(iface, mtu, name, useServiceTag, vlanID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVlanInterfaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "mtu", strconv.Itoa(mtu)),
					resource.TestCheckResourceAttr(resourceName, "vlan_id", strconv.Itoa(vlanID)),
				),
			},
			{
				Config: testAccVlanInterface(iface, mtu, name+"updated", useServiceTag, vlanID+1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVlanInterfaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name+"updated"),
					resource.TestCheckResourceAttr(resourceName, "mtu", strconv.Itoa(mtu)),
					resource.TestCheckResourceAttr(resourceName, "vlan_id", strconv.Itoa(vlanID+1)),
				),
			},
		},
	})
}

func TestVlanInterface_noVlanID(t *testing.T) {
	resourceName := "mikrotik_vlan_interface.testacc"
	iface := "ether1"
	mtu := 1500
	name := "test-vlan"
	useServiceTag := false
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVlanInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanInterfaceNoVLANID(iface, mtu, name, useServiceTag),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVlanInterfaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "mtu", strconv.Itoa(mtu)),
					resource.TestCheckResourceAttr(resourceName, "vlan_id", "1"),
				),
			},
		},
	})
}

func testAccVlanInterfaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("%s does not exist in the statefile", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())
		record, err := c.FindVlanInterface(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Unable to get remote record for %s: %v", resourceName, err)
		}

		if record == nil {
			return fmt.Errorf("Unable to get the remote record %s", resourceName)
		}

		return nil
	}
}

func testAccCheckVlanInterfaceDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_vlan_interface" {
			continue
		}

		remoteRecord, err := c.FindVlanInterface(rs.Primary.Attributes["name"])

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if remoteRecord != nil {
			return fmt.Errorf("remote record (%s) still exists", remoteRecord.Id)
		}

	}
	return nil
}

func testAccVlanInterface(iface string, mtu int, name string, useServiceTag bool, vlanID int) string {
	return fmt.Sprintf(`
		resource "mikrotik_vlan_interface" "testacc" {
			interface = %q
			mtu = %d
			name = %q
			use_service_tag = %t
			vlan_id = %d
		}
	`, iface, mtu, name, useServiceTag, vlanID)
}

func testAccVlanInterfaceNoVLANID(iface string, mtu int, name string, useServiceTag bool) string {
	return fmt.Sprintf(`
		resource "mikrotik_vlan_interface" "testacc" {
			interface = %q
			mtu = %d
			name = %q
			use_service_tag = %t
		}
	`, iface, mtu, name, useServiceTag)
}
