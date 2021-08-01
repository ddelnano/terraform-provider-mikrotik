package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var origName string = "terraform-acc-testing"

// TODO: Add dependent resources for owner
var origOwner string = "admin"
var origSource string = ":put testing"
var originalPolicy []string = []string{
	"ftp",
}
var updatedOwner string = "prometheus"
var updatedSource string = ":put updated"
var updatedPolicy []string = []string{
	"ftp", "dude",
}

func TestAccMikrotikScript_create(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
		},
	})
}

func TestAccMikrotikScript_updateSource(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source", origSource)),
			},
			{
				Config: testAccScriptRecordUpdatedSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source", updatedSource)),
			},
		},
	})
}

func TestAccMikrotikScript_updateOwner(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", origOwner)),
			},
			{
				Config: testAccScriptRecordUpdatedOwner(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", updatedOwner)),
			},
		},
	})
}

func TestAccMikrotikScript_updateDontReqPerms(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "false")),
			},
			{
				Config: testAccScriptRecordUpdatedDontReqPerms(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "true")),
			},
		},
	})
}

func TestAccMikrotikScript_updatePolicies(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy.0", "ftp"),
					resource.TestCheckResourceAttr(resourceName, "policy.1", "reboot")),
			},
			{
				Config: testAccScriptRecordUpdatedPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy.0", "ftp")),
			},
		},
	})
}

func TestAccMikrotikScript_import(t *testing.T) {
	resourceName := "mikrotik_script.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dont_require_permissions", "false")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccScriptRecord() string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    // These require a very specific order otherwise
    // the mikrotik command fails to create it.
    // TODO: Add an error to the client to fail this.
    policy = [
	"ftp", "reboot"
    ]
}
`, origName, origOwner, origSource)
}

func testAccScriptRecordUpdatedSource() string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = [
	"ftp"
    ]
}
`, origName, origOwner, updatedSource)
}

func testAccScriptRecordUpdatedOwner() string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = [
	"ftp"
    ]
}
`, origName, updatedOwner, origSource)
}

func testAccScriptRecordUpdatedDontReqPerms() string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = [
	"ftp"
    ]
    dont_require_permissions = true
}
`, origName, updatedOwner, origSource)
}

func testAccScriptRecordUpdatedPolicy() string {
	return fmt.Sprintf(`
resource "mikrotik_script" "bar" {
    name = "%s"
    owner = "%s"
    source = "%s"
    policy = [
	"ftp",
    ]
}
`, origName, origOwner, origSource)
}

func testAccCheckMikrotikScriptDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_script" {
			continue
		}

		script, err := c.FindScript(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
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

		script, err := c.FindScript(rs.Primary.ID)

		_, ok = err.(*client.NotFound)
		if !ok && err != nil {
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
