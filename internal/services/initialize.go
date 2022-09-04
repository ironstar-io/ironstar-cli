package services

import (
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func InitializeIronstarProject() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if exists {
		return errors.New("A config file has already been initialized in this directory")
	}

	projConf := types.ProjectConfig{
		Version: "1",
		Package: types.PackageConfig{
			Exclude: []string{".git", "private", "web/sites/default/files", "docroot/sites/default/files", "docroot/sites/default/settings.tok.php", "web/sites/default/settings.tok.php"},
		},
	}

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	err = SafeTouchConfigYAML(confPath)
	if err != nil {
		return err
	}

	fs.Replace(confPath, newMarhsalled, 0400)

	return nil
}
