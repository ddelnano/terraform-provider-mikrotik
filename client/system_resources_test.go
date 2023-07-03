package client

import (
	"strings"
	"testing"
)

func TestGetSystemResources(t *testing.T) {
	c := NewClient(GetConfigFromEnv())
	sysResources, err := c.GetSystemResources()

	if err != nil {
		t.Fatalf("failed to get system resources with error: %v", err)
	}

	if sysResources.Uptime <= 0 {
		t.Fatalf("expected uptime > 0, instead received '%d'", sysResources.Uptime)
	}

	version := sysResources.Version
	if strings.Index(version, "6") != 0 && strings.Index(version, "7") != 0 {
		t.Errorf("expected RouterOS version to start with a '7' or '6' major release, instead received '%s'", version)
	}
}
