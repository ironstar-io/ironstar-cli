package services

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func GetProjectData() (types.ProjectConfig, error) {
	// empty := types.ProjectConfig{}

	// pr := fs.ProjectRoot()
	// globals, err := ReadInGlobals()
	// if err != nil {
	// 	return empty, errors.Wrap(err, errs.SetCredentialsErrorMsg)
	// }

	// var projMatch types.ProjectConfig
	// if pr == constants.ProjectRootNotFound {
	// 	for _, proj := range globals.Projects {
	// 		if proj.Path == pr {
	// 			projMatch = proj
	// 		}
	// 	}
	// } else {
	// 	pn, err := GetCLIProjectName()
	// 	if err != nil {
	// 		return empty, errors.Wrap(err, errs.GetCredentialsErrorMsg)
	// 	}

	// 	for _, proj := range globals.Projects {
	// 		if proj.Name == pn {
	// 			projMatch = proj
	// 		}
	// 	}
	// }

	// return projMatch, nil
	return types.ProjectConfig{}, nil
}

func SetProjectSubscription(args []string) error {
	pr := fs.ProjectRoot()
	if pr == constants.ProjectRootNotFound {
		return errors.New(errs.NoProjectFoundErrorMsg)
	}

	return nil
	// var credMatch types.Credentials
	// for _, cred := range creds {
	// 	if cred.Login == email {
	// 		credMatch = cred
	// 	}
	// }

	// if (types.Credentials{}) == credMatch {
	// 	return errors.Wrap(errs.NoCredentialMatch, errs.SetCredentialsErrorMsg)
	// }

	// return utils.UpdateGlobalProjectLogin(pr, email)
}
