package client

import (
	"testing"
)

var sysResources SystemResources

func TestMain(m *testing.M) {
	SetupAndTestMainExec(m, &sysResources)
}
