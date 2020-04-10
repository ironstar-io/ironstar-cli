package auth

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/pkg/errors"
)

func SetDefaultCredentials(args []string) error {
	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	return services.UpdateGlobalDefaultLogin(email)
}
