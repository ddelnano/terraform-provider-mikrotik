package client

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func getRouterOSMajorVersion(systemResources SystemResources) (majorVersion int, err error) {
	if len(systemResources.Version) == 0 {
		return 0, errors.New("RouterOS system resources returned empty string")
	}
	majorVersion, err = strconv.Atoi(string(systemResources.Version[0]))
	return
}

func SkipIfRouterOSV6OrEarlier(t *testing.T, systemResources SystemResources) {
	majorVersion, err := getRouterOSMajorVersion(systemResources)
	if err != nil {
		t.Errorf("failed to get the system resource major version: %v", err)
	}
	if majorVersion <= 6 {
		t.Skip()
	}
}

func SkipIfRouterOSV7OrLater(t *testing.T, systemResources SystemResources) {
	majorVersion, err := getRouterOSMajorVersion(systemResources)
	if err != nil {
		t.Errorf("failed to get the system resource major version: %v", err)
	}
	if majorVersion >= 7 {
		t.Skip()
	}
}

// RandomString returns a random string
func RandomString() string {
	// a naive implementation with all-digits for now
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}
