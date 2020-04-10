package services

import (
	// "os"
	// "time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func ResolveUserCredentials(project types.ProjectConfig) (types.Keylink, error) {
	// empty := types.Credentials{}
	// ist := os.Getenv("IRONSTAR_SUBSCRIPTION_TOKEN")
	// if ist != "" {
	// 	return types.Credentials{
	// 		AuthToken: ist,
	// 		Login:     "IRONSTAR_SUBSCRIPTION_TOKEN",
	// 		Expiry:    time.Now().AddDate(25, 0, 0),
	// 		// These keys don't expire, add an arbitrary time
	// 	}, nil
	// }

	// if project == (types.ProjectConfig{}) {
	// 	globals, err := ReadInGlobals()
	// 	if err != nil {
	// 		return empty, errors.Wrap(err, errs.NoSuitableCredsMsg)
	// 	}

	// 	return PullLoginFromCredentials(globals.DefaultLogin)
	// }

	return PullLoginFromCredentials(project.Login)
}

func PullLoginFromCredentials(login string) (types.Keylink, error) {
	empty := types.Keylink{}
	if login == "" {
		return empty, errs.NoSuitableCreds
	}

	creds, err := ReadInCredentials()
	if err != nil {
		return empty, errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	var credMatch types.Keylink
	for _, cred := range creds.Keychain {
		if cred.Login == login {
			credMatch = cred
		}
	}

	if credMatch == empty {
		return empty, errs.NoSuitableCreds
	}

	return credMatch, nil
}
