package client

import (
	"os"
	"testing"
)

func SkipBgpIfUnsupported(t *testing.T) {
	if os.Getenv("LEGACY_BGP_SUPPORT") != "true" {
		t.Skip()
	}
}
