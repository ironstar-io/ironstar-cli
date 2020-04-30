package goos

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/cmd"
)

// OpenSite - Open a URL using PowerShell
func OpenSite(url string) (string, error) {
	return cmd.ChildProcess("powershell.exe", "start", url)
}
