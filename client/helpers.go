package client

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func SkipLegacyBgpIfUnsupported(t *testing.T) {
	if !IsLegacyBgpSupported() {
		t.Skip()
	}
}

func IsLegacyBgpSupported() bool {
	return os.Getenv("LEGACY_BGP_SUPPORT") == "true"
}

func SkipInterfaceWireguardIfUnsupported(t *testing.T) {
	if !IsInterfaceWireguardSupported() {
		t.Skip()
	}
}

func IsInterfaceWireguardSupported() bool {
	return os.Getenv("INTERFACE_WIREGUARD_SUPPORT") == "true"
}

// RandomString returns a random string
func RandomString() string {
	// a naive implementation with all-digits for now
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}
