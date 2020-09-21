package utils

// SliceEqual tells whether a and b contain the same elements in any order.
// A nil argument is equivalent to an empty slice.
func SliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, v := range a {
		var found bool = false
		for _, y := range b {
			if v == y {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// SliceEqual tells whether b contains at least the same elements of a in any order.
// This is to cater for the edge case there all items may be present in b that are required,
// but the check would otherwise fail if b contained an additional element not included in a
// A nil argument is equivalent to an empty slice.
func SliceIncludesAll(a, b []string) bool {
	if a == nil || b == nil {
		return false
	}
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	var c []string
	for _, v := range b {
		if SliceIncludes(a, v) {
			c = append(c, v)
		}
	}

	return SliceEqual(a, c)
}

// SliceIncludes tells whether val occurs one time in slice
func SliceIncludes(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func RemoveStringFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
