package services

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
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
