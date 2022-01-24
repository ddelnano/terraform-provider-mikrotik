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
