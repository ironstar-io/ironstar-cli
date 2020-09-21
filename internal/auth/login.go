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

	res, err := postLogin(email, password)
	if err != nil {
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	c, err := redirectChecks(res.Body, email)
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

func postLogin(email, password string) (*api.RawResponse, error) {
	req := &api.Request{
		RunTokenRefresh: false,
		Credentials:     types.Keylink{},
		Method:          "POST",
		Path:            "/auth/login",
		MapStringPayload: map[string]interface{}{
			"email":    email,
			"password": password,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APILoginErrorMsg)
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return res, nil
}

func postMFAValidate(MFAAuthToken, passcode string) (*api.RawResponse, error) {
	req := &api.Request{
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: MFAAuthToken,
		},
		Method: "POST",
		Path:   "/auth/mfa/validate",
		MapStringPayload: map[string]interface{}{
			"passcode": passcode,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return res, nil
}

func redirectChecks(body []byte, email string) (*types.AuthResponseBody, error) {
	b := &types.AuthResponseBody{}
	err := json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}

	switch b.RedirectEndpoint {
	case "/auth/password-reset":
		return resetUserPassword(email, b.IDToken, b.MFAStatus)
	case "/auth/mfa/validate":
		return validateMFAPasscode(b.IDToken)
	}

	return b, nil
}

func resetUserPassword(email, PWResetAuthToken, mfaStatus string) (*types.AuthResponseBody, error) {
	color.Yellow("Your password has expired! Please provide a new password.")
	fmt.Println()

	password, err := services.GetCLIPassword("")
	if err != nil {
		return nil, err
	}

	var passcode string
	if mfaStatus == "ENABLED" {
		pc, err := services.GetCLIMFAPasscode()
		if err != nil {
			return nil, err
		}
		passcode = pc
	}

	req := &api.Request{
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: PWResetAuthToken,
		},
		Method: "POST",
		Path:   "/auth/password-reset",
		MapStringPayload: map[string]interface{}{
			"password": password,
			"passcode": passcode,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 204 {
		return nil, res.HandleFailure()
	}

	fmt.Println()
	color.Green("Password reset completed successfully")

	lres, err := postLogin(email, password)
	if err != nil {
		return nil, errors.Wrap(err, errs.APILoginErrorMsg)
	}

	b := &types.AuthResponseBody{}
	err = json.Unmarshal(lres.Body, b)
	if err != nil {
		return nil, err
	}

	if b.MFAStatus == "ENABLED" {
		vres, err := postMFAValidate(b.IDToken, passcode)
		if err != nil {
			return nil, err
		}

		v := &types.AuthResponseBody{}
		err = json.Unmarshal(vres.Body, v)
		if err != nil {
			return nil, err
		}

		return v, nil
	}

	return b, nil
}

func validateMFAPasscodeWithRetries(MFAAuthToken string) (*types.AuthResponseBody, error) {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		c, err := validateMFAPasscode(MFAAuthToken)
		if err != nil {
			continue
		}

		return c, nil
	}

	return nil, errors.New("Unable to verify your MFA passcode. Code: FTevOS")
}

func validateMFAPasscode(MFAAuthToken string) (*types.AuthResponseBody, error) {
	passcode, err := services.GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}

	res, err := postMFAValidate(MFAAuthToken, passcode)
	if err != nil {
		return nil, err
	}

	m := &types.AuthResponseBody{}
	err = json.Unmarshal(res.Body, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
