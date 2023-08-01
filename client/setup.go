package client

import (
	"fmt"
	"os"
	"testing"
)

func SetupAndTestMainExec(m *testing.M, sysResources *SystemResources) {
	c := NewClient(GetConfigFromEnv())
	s, err := c.GetSystemResources()

	if err != nil {
		fmt.Printf("Unable to perform test setup, failed with error: %v\n", err)
		os.Exit(1)
	}

	*sysResources = *s

	os.Exit(m.Run())
}
