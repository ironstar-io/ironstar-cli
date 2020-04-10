package services

import (
	"path/filepath"
	// "strings"

	// "gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	// "github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func UpdateCredentialsFile(newCreds types.Keylink) error {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	// Pull out matching login if it exists in the map
	var splicedKeychain []types.Keylink
	for _, cred := range credSet.Keychain {
		if cred.Login != newCreds.Login {
			splicedKeychain = append(splicedKeychain, cred)
		}
	}

	// Replace/Add the new credentials to the struct slice
	newKeychain := append(splicedKeychain, newCreds)

	newMarhsalled, err := yaml.Marshal(&types.Credentials{
		Active:   credSet.Active,
		Keychain: newKeychain,
	})
	if err != nil {
		return err
	}

	fs.Replace(cp, newMarhsalled)

	return nil
}

func UpdateActiveCredentials(login string) (types.Keylink, error) {
	empty := types.Keylink{}
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return empty, err
	}

	var credMatch types.Keylink
	for _, cred := range credSet.Keychain {
		if cred.Login == login {
			credMatch = cred
		}
	}
	if credMatch == (types.Keylink{}) {
		return empty, errs.NoCredentialMatch
	}

	newMarhsalled, err := yaml.Marshal(&types.Credentials{
		Active:   login,
		Keychain: credSet.Keychain,
	})
	if err != nil {
		return empty, err
	}

	fs.Replace(cp, newMarhsalled)

	return credMatch, nil
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
