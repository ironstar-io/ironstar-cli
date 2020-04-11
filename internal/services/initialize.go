package services

import (
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

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
		Version: "1.0",
	}

	pn, err := StdinPrompt("Project Name: ")
	if err != nil {
		return errs.UnexpectedError
	}

	projConf.Project.Name = pn

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	err = SafeTouchConfigYAML(confPath)
	if err != nil {
		return err
	}

	fs.Replace(confPath, newMarhsalled)

	return nil
}
