package main

import (
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mikrotik.Provider,
	})
}
