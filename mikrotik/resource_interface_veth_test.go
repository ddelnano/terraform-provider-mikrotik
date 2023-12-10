package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)



func TestAccMikrotikInterfaceVeth_create(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-create")
	address := "192.168.88.15/24"
	gateway := "192.168.88.1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")

	resourceName := "mikrotik_interface_veth.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceVethDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceVeth(name, address, gateway, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false")),
			},
		},
	})
}


func TestAccMikrotikInterfaceVeth_updateAddr(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-create")
	address := "192.168.88.15/24"
	gateway := "192.168.88.1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")
	disabled := "false"
    updatedAddress := "192.168.188.15/24"
	updatedGateway := "192.168.188.1"
	updatedComment := acctest.RandomWithPrefix("tf-acc-comment")
	updatedDisabled := "true"

	resourceName := "mikrotik_interface_veth.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceVethDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceVeth(name, address, gateway, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccInterfaceVeth(name, updatedAddress, updatedGateway, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "address", updatedAddress),
					resource.TestCheckResourceAttr(resourceName, "gateway", updatedGateway),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccInterfaceVeth(name, address, gateway, updatedComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
					resource.TestCheckResourceAttr(resourceName, "disabled", disabled),
				),
			},
			{
				Config: testAccInterfaceVethUpdatedDisabled(name, address, gateway, comment, updatedDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "gateway", gateway),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "disabled", updatedDisabled),
				),
			},
		},
	})
}

func TestAccMikrotikInterfaceVeth_import(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-import")
	address := "192.168.88.15/24"
	gateway := "192.168.88.1"
	comment := acctest.RandomWithPrefix("tf-acc-comment")
	
	resourceName := "mikrotik_interface_veth.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikInterfaceVethDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceVeth(name, address, gateway, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccInterfaceVethExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateId: name,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInterfaceVeth(name string, address string, gateway string, comment string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_veth" "bar" {
		name = "%s"
		address = "%s"
		gateway = "%s"
		comment = "%s"
	}
	`, name, address, gateway, comment)
}

func testAccInterfaceVethUpdatedDisabled(name string, address string, gateway string, comment string, disabled string) string {
	return fmt.Sprintf(`
	resource "mikrotik_interface_veth" "bar" {
		name = "%s"
		address = "%s"
		gateway = "%s"
		comment = "%s"
		disabled = "%s"
	}
	`, name, address, gateway, comment, disabled)
}

func testAccCheckMikrotikInterfaceVethDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_interface_veth" {
			continue
		}

		interfaceVeth, err := c.FindInterfaceVeth(rs.Primary.Attributes["name"])

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if interfaceVeth != nil {
			return fmt.Errorf("interface veth (%s) still exists", interfaceVeth.Name)
		}
	}
	return nil
}

func testAccInterfaceVethExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_interface_veth does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		interfaceVeth, err := c.FindInterfaceVeth(rs.Primary.Attributes["name"])

		if err != nil {
			return fmt.Errorf("Unable to get the interface veth with error: %v", err)
		}

		if interfaceVeth == nil {
			return fmt.Errorf("Unable to get the interface veth with name: %s", rs.Primary.Attributes["name"])
		}

		if interfaceVeth.Name == rs.Primary.Attributes["name"] {
			return nil
		}
		return nil
	}
}
