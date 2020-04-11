package services

import (
	"os"
	"path/filepath"
)

func InitializeIronstarProject() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = SafeTouchConfigYAML(filepath.Join(wd, ".ironstar", "config.yml"))
	if err != nil {
		return err
	}

	return nil
}
