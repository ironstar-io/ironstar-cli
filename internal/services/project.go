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

func GetProjectData() (types.ProjectConfig, error) {
	empty := types.ProjectConfig{}

	pr := fs.ProjectRoot()
	proj, err := ReadInProjectConfig(pr)
	if err != nil {
		return empty, errors.Wrap(err, errs.NoProjectFoundErrorMsg)
	}

	return proj, nil
}

func LinkSubscriptionToProject(config types.ProjectConfig, sub types.Subscription) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	projConf, err := ReadInProjectConfig(wd)
	if err != nil {
		return errors.Wrap(err, errs.APISubLinkErrorMsg)
	}

	if projConf.Project.Name == "" {
		pn, err := StdinPrompt("Project Name: ")
		if err != nil {
			return errors.Wrap(err, errs.APISubLinkErrorMsg)
		}

		projConf.Project.Name = pn
	}

	projConf.Version = "1.0"
	projConf.Subscription = sub

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	py := filepath.Join(wd, ".ironstar", "config.yml")
	fs.Replace(py, newMarhsalled)

	return nil
}
