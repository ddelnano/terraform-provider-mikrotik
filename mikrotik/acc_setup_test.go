package mikrotik

import (
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
)

var sysResources client.SystemResources

func TestMain(m *testing.M) {
	client.SetupAndTestMainExec(m, &sysResources)
}
