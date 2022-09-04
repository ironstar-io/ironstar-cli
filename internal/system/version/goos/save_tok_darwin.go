package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"

	"github.com/pkg/errors"
)

// SaveCLIBinary - Saves the running instance of the Ironstar CLI binary to a persistent path in the user's /usr/bin folder
func SaveCLIBinary(version string) (string, error) {
	var empty string

	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathDarwin, version)
	b := filepath.Join(p, "iron")

	fmt.Println("Saving the running Ironstar CLI version [" + version + "] to [" + b + "]")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return empty, errors.Wrap(err, "There was an error creating the install directory")
	}

	ex, err := os.Executable()
	if err != nil {
		return empty, errors.Wrap(err, "Unexpected error obtaining this executable's path")
	}

	fs.Copy(ex, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		return empty, errors.Wrap(err, "Could not ensure correct permissions on ["+b+"]")
	}

	// Change file ownership.
	err = os.Chown(b, os.Getuid(), os.Getgid())
	if err != nil {
		return empty, errors.Wrap(err, "Could not ensure correct ownership on ["+b+"]")
	}

	return b, nil
}
