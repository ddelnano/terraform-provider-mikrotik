//go:build tools
// +build tools

package tools

import (
	_ "github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/codegen"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
