package mikrotik

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Provider name for single configuration testing
	ProviderNameMikrotik = "mikrotik"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameMikrotik: func() (*schema.Provider, error) { return Provider(), nil },
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
