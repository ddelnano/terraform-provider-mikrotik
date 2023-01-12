package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDurationUnmarshal(t *testing.T) {

	testCases := []struct {
		name        string
		in          string
		expected    MikrotikDuration
		expectError bool
	}{
		{
			name:     "single unit",
			in:       "23s",
			expected: MikrotikDuration(23),
		},
		{
			name:     "units below second are zeroed",
			in:       "20ms",
			expected: MikrotikDuration(0),
		},
		{
			name:     "parse week",
			in:       "2w",
			expected: MikrotikDuration(time.Hour.Seconds() * 24 * 7 * 2),
		},
		{
			name:     "multiple units",
			in:       "2h17m01s",
			expected: MikrotikDuration(time.Hour.Seconds()*2 + time.Minute.Seconds()*17 + 1),
		},
		{
			name:        "no-unit produces error",
			in:          "17",
			expectError: true,
		},
		{
			name:        "unit and no-unit produces error",
			in:          "2h17",
			expectError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := MikrotikDuration(0)
			err := (&m).UnmarshalMikrotik(tc.in)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, m)
		})
	}
}
