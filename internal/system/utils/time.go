package utils

import (
	"time"
)

func CalculateFriendlyETA(eta int) string {
	return (time.Duration(int64(eta)) * time.Second).Round(time.Second).String()
}

func UnixMilliseconds(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}
