package auth

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func SetProjectCredentials(args []string) error {
	pr := fs.ProjectRoot()
	if pr == constants.ProjectRootNotFound {
		return errors.New(errs.NoProjectFoundErrorMsg)
	}

	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	creds, err := services.ReadInCredentials()
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	var credMatch types.Credentials
	for _, cred := range creds {
		if cred.Login == email {
			credMatch = cred
		}
	}

	if (types.Credentials{}) == credMatch {
		return errors.Wrap(errs.NoCredentialMatch, errs.SetCredentialsErrorMsg)
	}

	return services.UpdateGlobalProjectLogin(pr, email)
}
