package services

import (
	"path/filepath"
	// "strings"

	// "gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	// "github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func UpdateCredentialsFile(newCreds types.Credentials) error {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	// Pull out matching login if it exists in the map
	var splicedCredSet []types.Credentials
	for _, cred := range credSet {
		if cred.Login != newCreds.Login {
			splicedCredSet = append(splicedCredSet, cred)
		}
	}

	// Replace/Add the new credentials to the struct slice
	newCredSet := append(splicedCredSet, newCreds)

	newMarhsalled, err := yaml.Marshal(newCredSet)
	if err != nil {
		return err
	}

	fs.Replace(cp, newMarhsalled)

	return nil
}

func UpdateGlobalProjectLogin(projectRoot, login string) error {
	// globals, err := ReadInGlobals()
	// if err != nil {
	// 	return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	// }

	// var newProjects []types.ProjectConfig
	// var projMatch bool = false
	// for _, proj := range globals.Projects {
	// 	if proj.Path == projectRoot {
	// 		newProjects = append(newProjects, types.ProjectConfig{
	// 			Name:  proj.Name,
	// 			Path:  proj.Path,
	// 			Login: login,
	// 		})

	// 		projMatch = true
	// 		continue
	// 	}

	// 	newProjects = append(newProjects, proj)
	// }

	// if !projMatch {
	// 	prs := strings.Split(projectRoot, "/")
	// 	if len(prs) == 0 {
	// 		return errors.New("An unexpected error occurred, exiting...")
	// 	}

	// 	name := prs[len(prs)-1]
	// 	newProjects = append(newProjects, types.ProjectConfig{
	// 		Name:  name,
	// 		Path:  projectRoot,
	// 		Login: login,
	// 	})
	// }

	// globals.Projects = newProjects

	// gp := filepath.Join(fs.HomeDir(), ".ironstar", "global.yml")
	// newMarhsalled, err := yaml.Marshal(globals)
	// if err != nil {
	// 	return err
	// }

	// fs.Replace(gp, newMarhsalled)

	return nil
}

func UpdateGlobalDefaultLogin(login string) error {
	// globals, err := ReadInGlobals()
	// if err != nil {
	// 	return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	// }

	// globals.DefaultLogin = login

	// gp := filepath.Join(fs.HomeDir(), ".ironstar", "global.yml")
	// newMarhsalled, err := yaml.Marshal(globals)
	// if err != nil {
	// 	return err
	// }

	// fs.Replace(gp, newMarhsalled)

	return nil
}
