package types

import (
	"strings"
)

// MikrotikList type translates slice of strings to comma separated list and back
//
// It is useful to seamless serialize/deserialize data during communication with RouterOS
type MikrotikList []string

func (m MikrotikList) MarshalMikrotik() string {
	return strings.Join(m, ",")
}

func (m *MikrotikList) UnmarshalMikrotik(value string) error {
	*m = strings.Split(value, ",")

	return nil
}
