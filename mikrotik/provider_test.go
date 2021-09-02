package mikrotik

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	// Provider name for single configuration testing
	ProviderNameMikrotik = "mikrotik"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

var apiClient *client.Mikrotik

func init() {
	apiClient = client.NewClient(os.Getenv("MIKROTIK_HOST"), os.Getenv("MIKROTIK_USER"), os.Getenv("MIKROTIK_PASSWORD"), false, "", true)

	testAccProvider = Provider(apiClient)
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameMikrotik: func() (*schema.Provider, error) { return testAccProvider, nil },
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MIKROTIK_HOST"); v == "" {
		t.Fatal("The MIKROTIK_HOST environment variable must be set")
	}
	if v := os.Getenv("MIKROTIK_USER"); v == "" {
		t.Fatal("The MIKROTIK_USER environment variable must be set")
	}
	if _, exists := os.LookupEnv("MIKROTIK_PASSWORD"); !exists {
		t.Fatal("The MIKROTIK_PASSWORD environment variable must be set")
	}
}

func testAccDeleteResource(resource *schema.Resource, d *schema.ResourceData, meta interface{}) error {
	if resource.DeleteContext != nil {
		var diags diag.Diagnostics

		diags = resource.DeleteContext(context.Background(), d, meta)

		for i := range diags {
			if diags[i].Severity == diag.Error {
				return fmt.Errorf("error deleting resource: %s", diags[i].Summary)
			}
		}

		return nil
	}

	return resource.Delete(d, meta)
}

func testAccCheckResourceDisappears(provider *schema.Provider, resource *schema.Resource, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		return testAccDeleteResource(resource, resource.Data(resourceState.Primary), provider.Meta())
	}
}
