package mikrotik

import (
	"context"
	"os"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Provider name for single configuration testing
	ProviderNameMikrotik = "mikrotik"
)

var testAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)
var testAccProvider *schema.Provider

var apiClient *client.Mikrotik

func init() {
	apiClient = client.NewClient(os.Getenv("MIKROTIK_HOST"), os.Getenv("MIKROTIK_USER"), os.Getenv("MIKROTIK_PASSWORD"), false, "", true)

	testAccProvider = Provider(apiClient)
	downgradedProviderFramework, _ := tf6to5server.DowngradeServer(
		context.Background(),
		providerserver.NewProtocol6(NewProviderFramework(apiClient)),
	)
	servers := []func() tfprotov5.ProviderServer{
		testAccProvider.GRPCProvider,
		func() tfprotov5.ProviderServer {
			return downgradedProviderFramework
		},
	}
	muxServer, _ := tf5muxserver.NewMuxServer(context.Background(), servers...)

	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		ProviderNameMikrotik: func() (tfprotov5.ProviderServer, error) {
			return muxServer, nil
		},
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
