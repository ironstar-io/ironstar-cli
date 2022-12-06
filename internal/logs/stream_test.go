package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLogLabelFlags(t *testing.T) {
	testCases := []struct {
		filenames  []string
		sources    []string
		errMessage string
	}{
		{[]string{"file1", "file2"}, []string{"source1", "source2"}, "The flags 'filenames' and 'sources' must not be set simultaneously"},
		{[]string{}, []string{"source1", "source2"}, ""},
		{[]string{"file1", "file2"}, []string{}, ""},
		{[]string{}, []string{}, ""},
	}

	for _, tc := range testCases {
		err := checkLogLabelFlags(tc.filenames, tc.sources)
		if tc.errMessage == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tc.errMessage)
		}
	}
}
