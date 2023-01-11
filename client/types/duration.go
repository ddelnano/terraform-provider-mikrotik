package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MikrotikDuration type represents a RouterOS durations [w,d] in seconds
type MikrotikDuration int

func (m MikrotikDuration) MarshalMikrotik() string {
	return strconv.Itoa(int(m))
}

func (m *MikrotikDuration) UnmarshalMikrotik(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return errors.New("cannot unmarshal empty value")
	}
	d, err := parseDuration(value)
	if err != nil {
		return err
	}
	*m = MikrotikDuration(d.Seconds())

	return nil
}

func parseDuration(s string) (time.Duration, error) {
	var digitsStartIndex, unitStartIndex int
	var nanoseconds int64

	parsePart := func(s string, unitStart int) (int64, error) {
		var ret int64
		digits, err := strconv.Atoi(s[:unitStart])
		if err != nil {
			return 0, err
		}

		unit := s[unitStart:]
		switch unit {
		case "ns":
			ret = int64(digits)
		case "us":
			ret = int64(digits) * time.Microsecond.Nanoseconds()
		case "ms":
			ret = int64(digits) * time.Millisecond.Nanoseconds()
		case "s":
			ret = int64(digits) * time.Second.Nanoseconds()
		case "m":
			ret = int64(digits) * time.Minute.Nanoseconds()
		case "h":
			ret = int64(digits) * time.Hour.Nanoseconds()
		case "d":
			ret = int64(digits) * time.Hour.Nanoseconds() * 24
		case "w":
			ret = int64(digits) * time.Hour.Nanoseconds() * 24 * 7
		default:
			return 0, fmt.Errorf("unknown unit: %q", unit)
		}
		return ret, nil
	}

	for i := 0; i < len(s); i++ {
		char := string(s[i])
		if char >= "0" && char <= "9" {
			if unitStartIndex > digitsStartIndex {
				parsed, err := parsePart(s[digitsStartIndex:i], unitStartIndex-digitsStartIndex)
				if err != nil {
					return 0, err
				}
				nanoseconds += parsed
				digitsStartIndex = i
				unitStartIndex = i
			}
			continue
		}
		if digitsStartIndex == unitStartIndex {
			unitStartIndex = i
		}
		continue
	}
	if digitsStartIndex == unitStartIndex {
		return 0, errors.New("duration without unit is not supported")
	}
	parsed, err := parsePart(s[digitsStartIndex:], unitStartIndex-digitsStartIndex)
	if err != nil {
		return 0, err
	}
	nanoseconds += parsed

	return time.Duration(nanoseconds), nil
}
