package auth

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func IronstarAPILogin(args []string, flg flags.Accumulator) error {
	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	password, err := services.GetCLIPassword(flg.Password)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	req := &api.Request{
		RunTokenRefresh: false,
		Credentials:     types.Keylink{},
		Method:          "POST",
		Path:            "/auth/login",
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

	err = services.UpdateCredentialsFile(types.Keylink{
		Login:     email,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	})
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	_, err = services.UpdateActiveCredentials(email)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
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

	proj, err := services.GetProjectDataSkipNew()
	if err != nil {
		return nil
	}
	if proj.Subscription != (types.Subscription{}) && proj.Subscription.Alias == "" {
		fmt.Println()
		color.Yellow("You have logged in successfully and can now link this project to an Ironstar subscription. Run `iron subscription list` to see a list of available subscriptions")
	}

	return nil
}

func mfaCredentialCheck(body []byte, email string) (*types.AuthResponseBody, error) {
	b := &types.AuthResponseBody{}
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

func validateMFAPasscode(logResBody *types.AuthResponseBody) (*types.AuthResponseBody, error) {
	passcode, err := services.GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}

	req := &api.Request{
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: logResBody.IDToken,
		},
		Method: "POST",
		Path:   "/auth/mfa/validate",
		MapStringPayload: map[string]string{
			"passcode": passcode,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	m := &types.AuthResponseBody{}
	err = json.Unmarshal(res.Body, m)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return m, nil
}
