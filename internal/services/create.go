package services

import (
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
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
		return fs.TouchEmpty(path)
	}

	return nil
}
