package goos

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
)

// SaveCLIBinary - Saves the running instance of the Ironstar CLI binary to a persistent path in the user's /usr/bin folder
func SaveCLIBinary(version string) (string, error) {
	p := filepath.Join(fs.HomeDir(), constants.BaseInstallPathDarwin, version)
	b := filepath.Join(p, "iron")

	fmt.Println("Saving the running Ironstar CLI version [" + version + "] to [" + b + "]")

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fmt.Println("There was an error creating the install directory: ", err.Error())
		os.Exit(1)
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Unexpected error obtaining this executable's path: ", err.Error())
		os.Exit(1)
	}

	fs.Copy(ex, b)
	// Change file permission bit
	err = os.Chmod(b, 0755)
	if err != nil {
		fmt.Println("Could not ensure correct permissions on ["+b+"]: ", err.Error())
		os.Exit(1)
	}

	// Change file ownership.
	err = os.Chown(b, os.Getuid(), os.Getgid())
	if err != nil {
		fmt.Println("Could not ensure correct ownership on ["+b+"]: ", err.Error())
		os.Exit(1)
	}

	return b, nil
}
