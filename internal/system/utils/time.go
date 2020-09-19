package utils

import (
	"strconv"
)

func CalculateFriendlyETA(eta int) string {
	if eta <= 60 {
		return strconv.Itoa(eta) + " seconds "
	}

	if eta <= 3600 {
		minutes := eta / 60
		seconds := eta % 60

		var minStr = " minute "
		if minutes > 1 {
			minStr = " minutes "
		}
		var secStr = " second "
		if seconds > 1 {
			secStr = " seconds "
		}
		return strconv.Itoa(minutes) + minStr + strconv.Itoa(seconds) + secStr
	}

	minutes := eta / 60
	hours := minutes / 60
	rMins := minutes % 60

	var hrStr = " hour "
	if hours > 1 {
		hrStr = " hours "
	}
	var minStr = " minute "
	if rMins > 1 {
		minStr = " minutes "
	}
	return strconv.Itoa(hours) + hrStr + strconv.Itoa(minutes) + minStr
}
