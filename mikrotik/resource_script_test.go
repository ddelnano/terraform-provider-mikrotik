package mikrotik

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO: Add dependent resources for owner
var defaultOwner = "admin"
var defaultSource = ":put testing"
var defaultPolicies = []string{"ftp", "reboot"}

func TestAccMikrotikScript_create(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-create")

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
		},
	})
}

func TestAccMikrotikScript_updateSource(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-update-src")
	updatedSource := ":put updated"

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source", defaultSource)),
			},
			{
				Config: testAccScriptRecord(name, defaultOwner, updatedSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source", updatedSource)),
			},
		},
	})
}

func TestAccMikrotikScript_updateOwner(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-update-owner")
	updatedOwner := "prometheus"

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", defaultOwner)),
			},
			{
				Config: testAccScriptRecord(name, updatedOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", updatedOwner)),
			},
		},
	})
}

func TestAccMikrotikScript_updateDontReqPerms(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-update-perm")

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "false")),
			},
			{
				Config: testAccScriptRecordWithPerms(name, defaultOwner, defaultSource, defaultPolicies, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "true")),
			},
		},
	})
}

func TestAccMikrotikScript_updatePolicies(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-update-pol")
	updatedPolicies := []string{"ftp"}

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy.0", defaultPolicies[0]),
					resource.TestCheckResourceAttr(resourceName, "policy.1", defaultPolicies[1])),
			},
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, updatedPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy.0", updatedPolicies[0])),
			},
		},
	})
}

func TestAccMikrotikScript_import(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-import")

	resourceName := "mikrotik_script.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(name, defaultOwner, defaultSource, defaultPolicies),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "false")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     name,
			},
		},
	})
}

func testAccScriptRecord(name, owner, source string, policies []string) string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = ["%s"]
}
`, name, owner, source, strings.Join(policies, "\",\""))
}

func testAccScriptRecordWithPerms(name, owner, source string, policies []string, dontRequirePermissions bool) string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = ["%s"]
    dont_require_permissions = %t
}
`, name, owner, source, strings.Join(policies, "\",\""), dontRequirePermissions)
}

func testAccCheckMikrotikScriptDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_script" {
			continue
		}

		script, err := c.FindScript(rs.Primary.ID)

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if script != nil && script.Name != "" {
			return fmt.Errorf("script (%s) still exists", script.Name)
		}
	}
	return nil
}

func testAccScriptExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_script does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		script, err := c.FindScript(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Unable to get the script with error: %v", err)
		}

		if script.Name == "" {
			return fmt.Errorf("Unable to get the script with name: %s", rs.Primary.ID)
		}

		if script.Name == rs.Primary.ID {
			return nil
		}
		return nil
	}
}
