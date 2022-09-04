package auth

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Logout(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	if creds.Login == "IRONSTAR_SUBSCRIPTION_TOKEN" {
		email, err := services.GetCLIEmail(args)
		if err != nil {
			return errors.Wrap(err, errs.APILoginErrorMsg)
		}

		revCreds, err := services.ResolveUserCredentials(email)
		if err != nil {
			return err
		}

		creds = revCreds
	} else {
		fmt.Println()
		confirmLogout := services.ConfirmationPrompt("Would you like to log out of your Ironstar session for <"+creds.Login+">?", "y", flg.AutoAccept)
		if !confirmLogout {
			fmt.Println("Exiting...")
			return nil
		}
	}

	_, err = postLogout(creds)
	if err != nil {
		return errors.Wrap(err, errs.APILogoutErrorMsg)
	}

	err = services.RemoveCredentials(creds)
	if err != nil {
		return errors.Wrap(err, errs.APILogoutErrorMsg)
	}

	fmt.Println()
	color.Green("Ironstar API logout succeeded! The session token has been destroyed.")

	return nil
}

func postLogout(creds types.Keylink) (*api.RawResponse, error) {
	req := &api.Request{
		RunTokenRefresh:  false,
		Credentials:      creds,
		Method:           "POST",
		Path:             "/auth/logout",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APILoginErrorMsg)
	}

	if res.StatusCode != 204 {
		return nil, res.HandleFailure()
	}

	return res, nil
}
