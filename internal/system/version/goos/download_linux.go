package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/console"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
)

// GetInstallPath - Check if tok version is installed or not
func GetInstallPath(version string) string {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathLinux, version, "tok")
	if fs.CheckExists(p) == true {
		return p
	}

	return ""
}

// DownloadTokBinary - Install a selected tok version and returns install path
func DownloadTokBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathLinux, version)
	b := filepath.Join(p, "tok")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory: ", err.Error())
		os.Exit(1)
	}

	fmt.Println()
	w := console.SpinStart("Downloading the specified release from GitHub.")
	err = utils.DownloadFile(b, constants.BaseBinaryURL+version+"/"+constants.BinaryNameLinux)
	if err != nil {
		return "", err
	}
	console.SpinPersist(w, "🚉", "Download complete!")

	return b, nil
}
