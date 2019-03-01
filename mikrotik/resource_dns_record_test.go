package mikrotik

import (
	"testing"
)

func TestTtlToSeconds(t *testing.T) {
	tests := []struct {
		expected int
		input    string
	}{
		{300, "5m"},
		{59, "59s"},
		{301, "5m1s"},
		{55259, "15h20m59s"},
		{141659, "1d15h20m59s"},
		{228059, "2d15h20m59s"},
		{86400, "1d"},
	}

	for _, test := range tests {
		actual := ttlToSeconds(test.input)
		if test.expected != actual {
			t.Errorf("Input %s returned %d instead of %d", test.input, actual, test.expected)
		}
	}
}
