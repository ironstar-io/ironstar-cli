package utils

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils/goos"
)

// OpenSite - Split for multi platform
func OpenSite(url string) (string, error) {
	return goos.OpenSite(url)
}
