package services

import (
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
)

func SafeTouchConfigYAML(path string) error {
	// Initialise the config file if it doesn't exist
	var _, errf = os.Stat(path)
	if os.IsNotExist(errf) {
		// The global .ironstar path requires appropriate permissions
		np := filepath.Dir(path)
		if !fs.CheckExists(np) {
			err := os.MkdirAll(np, 0700)
			if err != nil {
				return err
			}
		}
		return fs.TouchEmpty(path, 0400)
	}

	return nil
}

func SafeTouchCredentialsYAML() error {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")
	exists := fs.CheckExists(cp)
	if !exists {
		np := filepath.Dir(cp)
		if !fs.CheckExists(np) {
			err := os.MkdirAll(np, 0700)
			if err != nil {
				return err
			}
		}
		return fs.TouchEmpty(cp, 0400)
	}

	return nil
}
