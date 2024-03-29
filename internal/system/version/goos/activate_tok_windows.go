package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"

	"github.com/pkg/errors"
)

// ActivateSavedVersion - Copies the specified version (which may be downloaded previously) into /usr/local/bin on macOS
func ActivateSavedVersion(version string) error {
	// activePath is where the currently active version of Ironstar CLI is found, such as /c/Users/Frank/bin/iron
	activePath := filepath.Join(fs.HomeDir(), "bin", "iron")

	// savePath is where to save the Ironstar CLI binary, such as /c/Users/Frank/AppData/Local/Ironstar/CLI/{version}/iron
	savePath := filepath.Join(fs.HomeDir(), constants.BaseInstallPathWindows, version)

	// Check if the requested version is not downloaded already
	p := filepath.Join(savePath, "iron")
	if fs.CheckExists(p) != true {
		fmt.Println("Ironstar CLI version [" + version + "] was not found at [" + p + "], downloading a new copy...")

		_, err := DownloadCLIBinary(version)
		if err != nil {
			return errors.Wrap(err, "Unexpected error downloading that version")
		}
	}

	// Remove any existing global binary
	fmt.Println("removing any existing Ironstar CLI version at [" + activePath + "]")
	err := os.Remove(activePath)
	if err != nil {
		fmt.Println(err)
	}

	// Remove any existing copy of Ironstar CLI at ~/bin/iron
	fs.Remove(activePath)

	// Copy the specified version to ~/bin/iron
	fs.Copy(p, activePath)

	// Make sure the version is executable
	err = os.Chmod(activePath, 0777)
	if err != nil {
		return errors.Wrap(err, "Unexpected error granting execute permissions to ["+activePath+"]")
	}

	return nil
}
