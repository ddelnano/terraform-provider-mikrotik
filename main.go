package main

import (
	"context"
	"flag"
	"log"

	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serverOpts := []tf5server.ServeOpt{}
	if debugMode {
		serverOpts = append(serverOpts, tf5server.WithManagedDebug())
	}

	providers := []func() tfprotov5.ProviderServer{
		func() tfprotov5.ProviderServer { return mikrotik.NewProvider().GRPCProvider() },
	}
	ctx := context.Background()

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	err = tf5server.Serve("registry.terraform.io/ddelnano/mikrotik", muxServer.ProviderServer, serverOpts...)
	if err != nil {
		log.Fatal(err)
	}
}
