package auth

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func IronstarAPILogin(args []string, passwordFlag string) error {
	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	password, err := services.GetCLIPassword(passwordFlag)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	req := &api.Request{
		AuthToken: "",
		Method:    "POST",
		Path:      "/auth/login",
		MapStringPayload: map[string]string{
			"email":    email,
			"password": password,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	c, err := mfaCredentialCheck(res.Body, email)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	err = services.UpdateCredentialsFile(types.Credentials{
		Login:     email,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	})
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	pr := fs.ProjectRoot()
	if pr != constants.ProjectRootNotFound {
		err = services.UpdateGlobalProjectLogin(pr, email)
		if err != nil {
			fmt.Println()
			color.Yellow("Authentication succeeded, but Tokaido was unable to update global credentials: ", err.Error())
		}
	}

	fmt.Println()
	color.Green("Ironstar API authentication succeeded!")
	fmt.Println()
	color.Green("User: ")
	fmt.Println(email)

	if c.Expiry.IsZero() {
		// Should always return expiry, but check anyway for safety
		return nil
	}

	fmt.Println()
	color.Green("Expiry: ")

	expDiff := strconv.Itoa(int(math.RoundToEven(c.Expiry.Sub(time.Now().UTC()).Hours() / 24)))
	fmt.Println(c.Expiry.String() + " (" + expDiff + " days)")

	return nil
}

func mfaCredentialCheck(body []byte, email string) (*types.AuthLoginBody, error) {
	b := &types.AuthLoginBody{}
	err := json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}

	// If this is set, user is an MFA user
	if b.RedirectEndpoint != "" {
		c, err := validateMFAPasscode(b)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return b, nil
}

func validateMFAPasscode(logResBody *types.AuthLoginBody) (*types.AuthLoginBody, error) {
	passcode, err := services.GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}

	req := &api.Request{
		AuthToken: logResBody.IDToken,
		Method:    "POST",
		Path:      "/auth/mfa/validate",
		MapStringPayload: map[string]string{
			"passcode": passcode,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	m := &types.AuthLoginBody{}
	err = json.Unmarshal(res.Body, m)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return m, nil
}
