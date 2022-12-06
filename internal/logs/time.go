package logs

import (
	"strconv"
	"time"
)

func formatLogTimestamp(ts string) string {
	tsint, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ts
	}

	// Format the time value as an RFC 3339 string
	return time.Unix(0, tsint*int64(time.Millisecond)).UTC().Format(time.RFC3339)
}
