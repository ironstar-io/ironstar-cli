package goos

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/cmd"
)

// OpenSite - Linux Root executable
func OpenSite(url string) (string, error) {
	return cmd.ChildProcess("xdg-open", url)
}
