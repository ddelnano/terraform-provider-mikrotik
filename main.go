package main

import (
	"context"
	"flag"
	"log"

	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	downgradedProviderFramework, err := tf6to5server.DowngradeServer(ctx, providerserver.NewProtocol6(mikrotik.NewProviderFramework()))
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov5.ProviderServer{
		mikrotik.NewProvider().GRPCProvider,
		func() tfprotov5.ProviderServer { return downgradedProviderFramework },
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	serverOpts := []tf5server.ServeOpt{}
	if debugMode {
		serverOpts = append(serverOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve("registry.terraform.io/ddelnano/mikrotik", muxServer.ProviderServer, serverOpts...)
	if err != nil {
		log.Fatal(err)
	}
}
