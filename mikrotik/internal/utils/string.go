package utils

import (
	"fmt"
	"strconv"
)

// ParseBool is wrapper around strconv.ParseBool to save few lines of code
func ParseBool(v string) (bool, error) {
	res, err := strconv.ParseBool(v)
	if err != nil {
		return res, fmt.Errorf("could not parse %q as bool: %w", v, err)
	}

	return res, nil
}
