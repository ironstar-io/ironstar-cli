package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
)

// ActivateSavedVersion - Copies the specified version (which may be downloaded previously) into /usr/local/bin on macOS
func ActivateSavedVersion(version string) bool {
	// Check that the version is downloaded already
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathDarwin, version, "tok")
	if fs.CheckExists(p) != true {
		fmt.Println("Ironstar CLI version [" + version + "] was not found in ~/.tok/bin, downloading a new copy...")

		_, err := DownloadTokBinary(version)
		if err != nil {
			fmt.Println("Unexpected error downloading that version: " + err.Error())
			os.Exit(1)
		}
	}

	// Remove any existing copy of Ironstar CLI at /usr/local/bin/tok
	fs.Remove(constants.ActiveBinaryPathDarwin)

	// Copy the specified version to /usr/local/bin/tok
	fs.Copy(p, constants.ActiveBinaryPathDarwin)

	// Make sure the version is executable
	err := os.Chmod(constants.ActiveBinaryPathDarwin, 0777)
	if err != nil {
		fmt.Println("Unexpected error granting execute permissions to [" + constants.ActiveBinaryPathDarwin + "]: " + err.Error())
		os.Exit(1)
	}

	return true
}
