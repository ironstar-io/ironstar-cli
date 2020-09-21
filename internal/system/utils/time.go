package utils

import (
	"time"
)

func CalculateFriendlyETA(eta int) string {
	return (time.Duration(int64(eta)) * time.Second).Round(time.Second).String()
}
