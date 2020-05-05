package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
)

// ActivateSavedVersion - Copies the specified version (which may be downloaded previously) into /usr/local/bin on macOS
func ActivateSavedVersion(version string) (bool, error) {
	// Check that the version is downloaded already
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathLinux, version, "iron")
	if !fs.CheckExists(p) {
		fmt.Println("Ironstar CLI version [" + version + "] was not found in ~/.ironstar/bin, downloading a new copy...")

		_, err := DownloadCLIBinary(version)
		if err != nil {
			return false, errors.Wrap(err, "Unexpected error downloading that version")
		}
	}

	// Remove any existing copy of Ironstar CLI at /usr/local/bin/iron
	fs.Remove(constants.ActiveBinaryPathDarwin)

	// Copy the specified version to /usr/local/bin/iron
	fs.Copy(p, constants.ActiveBinaryPathDarwin)

	// Make sure the version is executable
	err := os.Chmod(constants.ActiveBinaryPathDarwin, 0777)
	if err != nil {
		return false, errors.Wrap(err, "Unexpected error granting execute permissions to ["+constants.ActiveBinaryPathDarwin+"]: "+err.Error())
	}

	return true
}
