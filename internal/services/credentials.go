package services

import (
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func ResolveUserCredentials(loginFlag string) (types.Keylink, error) {
	// Priority #1 Return the credentials for the --login=x flag
	if loginFlag != "" {
		return PullCredentialsByLogin(loginFlag)
	}

	// Priority #2 Return the credentials contained in the IRONSTAR_SUBSCRIPTION_TOKEN env var
	ist := os.Getenv("IRONSTAR_SUBSCRIPTION_TOKEN")
	if ist != "" {
		return types.Keylink{
			AuthToken: ist,
			Login:     "IRONSTAR_SUBSCRIPTION_TOKEN",
			Expiry:    time.Now().AddDate(25, 0, 0),
			// These keys don't expire, add an arbitrary time
		}, nil
	}

	// Priority #3 Return the credential set that has been set active
	return PullActiveCredentials()
}

func PullActiveCredentials() (types.Keylink, error) {
	empty := types.Keylink{}
	creds, err := ReadInCredentials()
	if err != nil {
		return empty, errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	var credMatch types.Keylink
	for _, cred := range creds.Keychain {
		if cred.Login == creds.Active {
			credMatch = cred
		}
	}

	if credMatch == empty {
		return empty, errs.NoSuitableCreds
	}

	return credMatch, nil
}

func PullCredentialsByLogin(login string) (types.Keylink, error) {
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
