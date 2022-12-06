package logs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatLogTimestamp(t *testing.T) {
	testCases := []struct {
		timestamp string
		expected  string
	}{
		{timestamp: "1623471213000", expected: "2021-06-12T04:13:33Z"},
		{timestamp: "invalid", expected: "invalid"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s -> %s", tc.timestamp, tc.expected), func(t *testing.T) {
			actual := formatLogTimestamp(tc.timestamp)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
