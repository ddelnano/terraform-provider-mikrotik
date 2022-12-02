package types

import (
	"strconv"
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

// MikrotikIntList type translates slice of ints to comma separated list and back
type MikrotikIntList []int

func (m MikrotikIntList) MarshalMikrotik() string {
	if len(m) == 0 {
		return ""
	}
	if len(m) == 1 {
		return strconv.Itoa(m[0])
	}

	buf := strings.Builder{}
	buf.WriteString(strconv.Itoa(m[0]))
	for i := range m[1:] {
		buf.WriteRune(',')
		buf.WriteString(strconv.Itoa(m[i+1]))
	}

	return buf.String()
}

func (m *MikrotikIntList) UnmarshalMikrotik(value string) error {
	if len(value) == 0 {
		return nil
	}
	stringSlice := strings.Split(value, ",")
	res := []int{}
	for _, s := range stringSlice {
		elem, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		res = append(res, elem)
	}
	*m = res

	return nil
}
