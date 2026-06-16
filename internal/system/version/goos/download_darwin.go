package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/system/console"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/pkg/errors"
)

// GetInstallPath - Check if Ironstar CLI version is installed or not
func GetInstallPath(version string) string {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathDarwin, version, "iron")
	if fs.CheckExists(p) {
		return p
	}

	return ""
}

// DownloadCLIBinary - Install a selected Ironstar CLI version and returns install path
func DownloadCLIBinary(version string) (string, error) {
	var empty string

	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathDarwin, version)
	b := filepath.Join(p, "iron")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return empty, errors.Wrap(err, "There was an error creating the install directory")
	}

	// macOS ships a single signed, notarised universal binary that runs on
	// both Intel and Apple Silicon.
	bin := constants.BinaryNameMacOS

	fmt.Println()

	w := console.SpinStart("Downloading the specified release from GitHub.")
	err = utils.DownloadFile(b, constants.BaseBinaryURL+version+"/"+bin)
	if err != nil {
		return empty, err
	}
	console.SpinPersist(w, "🚉", "Download complete!")

	return b, nil
}
