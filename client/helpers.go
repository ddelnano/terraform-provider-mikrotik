package client

import (
	"os"
	"testing"
)

func SkipLegacyBgpIfUnsupported(t *testing.T) {
	if !IsLegacyBgpSupported() {
		t.Skip()
	}
}

func IsLegacyBgpSupported() bool {
	if os.Getenv("LEGACY_BGP_SUPPORT") == "true" {
		return true
	}
	return false
}

func SkipInterfaceWireguardIfUnsupported(t *testing.T) {
	if !IsInterfaceWireguardSupported() {
		t.Skip()
	}
}

func IsInterfaceWireguardSupported() bool {
	if os.Getenv("INTERFACE_WIREGUARD_SUPPORT") == "true" {
		return true
	}
	return false
}
