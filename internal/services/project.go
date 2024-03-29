package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func GetProjectData(autoAccept bool) (types.ProjectConfig, error) {
	empty := types.ProjectConfig{}
	wd, err := os.Getwd()
	if err != nil {
		return empty, err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if !exists {
		createNewProj := ConfirmationPrompt("Couldn't find a project configuration in this directory. Would you like to create one?", "y", autoAccept)
		if createNewProj == true {
			err = InitializeIronstarProject()
			if err != nil {
				return empty, err
			}
		} else {
			return empty, errors.New("This command requires a project to be configured.")
		}

	}

	pr := fs.ProjectRoot()
	proj, err := ReadInProjectConfig(pr)
	if err != nil {
		return empty, errors.Wrap(err, errs.NoProjectFoundErrorMsg)
	}

	return proj, nil
}

func GetProjectDataSkipNew() (types.ProjectConfig, error) {
	empty := types.ProjectConfig{}
	wd, err := os.Getwd()
	if err != nil {
		return empty, err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if !exists {
		return types.ProjectConfig{}, nil
	}

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

	if projConf.Subscription != (types.Subscription{}) && projConf.Subscription.Alias != "" && projConf.Subscription.Alias != sub.Alias {
		color.Yellow("This project was previously linked to subscription [" + projConf.Subscription.Alias + "]. The link will be replaced with the subscription [" + sub.Alias + "].")
		fmt.Println()
	}

	projConf.Version = "1.0"
	projConf.Subscription = sub

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	py := filepath.Join(wd, ".ironstar", "config.yml")
	fs.Replace(py, newMarhsalled, 0400)

	return nil
}
