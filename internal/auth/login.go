package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	res, err := postLogin(email, password, flg.LockSessionToIP)
	if err != nil {
		return err
	}

	c, err := loginRedirectChecks(flg, res.Body, email)
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

	expDiff := strconv.Itoa(int(math.RoundToEven(c.Expiry.Sub(time.Now().UTC()).Hours())))
	fmt.Println(c.Expiry.String() + " (" + expDiff + " hours)")

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

func postLogin(email, password string, lockSessionToIP bool) (*api.RawResponse, error) {
	req := &api.Request{
		Retries:         3,
		RunTokenRefresh: false,
		Credentials:     types.Keylink{},
		Method:          "POST",
		Path:            "/auth/login",
		MapStringPayload: map[string]interface{}{
			"email":              email,
			"password":           password,
			"lock_session_to_ip": lockSessionToIP,
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
		Retries:         3,
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: MFAAuthToken,
		},
		Method: "POST",
		Path:   "/auth/mfa/validate",
		MapStringPayload: map[string]interface{}{
			"passcode": passcode,
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

func loginRedirectChecks(flg flags.Accumulator, body []byte, email string) (*types.AuthResponseBody, error) {
	b := &types.AuthResponseBody{}
	err := json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}

	switch b.RedirectEndpoint {
	case "/auth/mfa/validate":
		return validateMFAPasscode(flg, b.IDToken, email)
	case "/auth/mfa/enable":
		fmt.Println()
		fmt.Println()
		color.Yellow(`NOTICE!`)
		color.Yellow(`Usage of the Ironstar API now requires all users to have MFA enabled`)
		fmt.Println()

		cnt := services.ConfirmationPrompt("Would like to enable MFA now?", "y", false)
		if !cnt {
			return nil, errors.New("Unable to proceed without enabling MFA.")
		}

		return MFAEnable(flg, types.Keylink{
			Login:     email,
			AuthToken: b.IDToken,
			Expiry:    b.Expiry,
		}, email)
	}

	return b, nil
}

func mfaValidateRedirectChecks(flg flags.Accumulator, body []byte, email string) (*types.AuthResponseBody, error) {
	b := &types.AuthResponseBody{}
	err := json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}

	switch b.RedirectEndpoint {
	case "/auth/password-reset":
		return resetUserPassword(flg, email, b.IDToken, b.MFAStatus)
	case "/auth/accept-api-usage-terms":
		return acceptAPIUsageTerms(b.IDToken)
	}

	return b, nil
}

func resetUserPassword(flg flags.Accumulator, email, PWResetAuthToken, mfaStatus string) (*types.AuthResponseBody, error) {
	color.Yellow("Your password has expired! Please provide a new password.")
	fmt.Println()

	password, err := services.GetCLIPassword("")
	if err != nil {
		return nil, err
	}

	var passcode string
	if mfaStatus == "VIRTUAL_MFA_ENABLED" {
		pc, err := services.GetCLIMFAPasscode()
		if err != nil {
			return nil, err
		}
		passcode = pc
	}

	req := &api.Request{
		Retries:         3,
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

	lres, err := postLogin(email, password, flg.LockSessionToIP)
	if err != nil {
		return nil, err
	}

	b := &types.AuthResponseBody{}
	err = json.Unmarshal(lres.Body, b)
	if err != nil {
		return nil, err
	}

	if b.MFAStatus == "VIRTUAL_MFA_ENABLED" {
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

func validateMFAPasscodeWithRetries(flg flags.Accumulator, MFAAuthToken, email string) (*types.AuthResponseBody, error) {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		c, err := validateMFAPasscode(flg, MFAAuthToken, email)
		if err != nil {
			fmt.Println(err.Error() + ". Retry...")

			continue
		}

		return c, nil
	}

	return nil, errors.New("Unable to verify your MFA passcode. Code: FTevOS")
}

func validateMFAPasscode(flg flags.Accumulator, MFAAuthToken, email string) (*types.AuthResponseBody, error) {
	passcode, err := services.GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}
	fmt.Println()

	res, err := postMFAValidate(MFAAuthToken, passcode)
	if err != nil {
		return nil, err
	}

	return mfaValidateRedirectChecks(flg, res.Body, email)
}

func acceptAPIUsageTerms(APIUsageAcceptAuthToken string) (*types.AuthResponseBody, error) {
	aut, err := api.GetCurrentAPIUsageTerms()
	if err != nil {
		return nil, err
	}

	fmt.Println()
	color.Yellow(`NOTICE!`)
	color.Yellow(`Your use of the Ironstar API, CLI and Console is subject to the "Ironstar API Usage Terms".`)
	fmt.Println()

	terms, err := base64.StdEncoding.DecodeString(aut.Terms)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(terms))
	fmt.Println()

	cnt := services.ConfirmationPrompt("Do you accept these usage terms?", "y", false)
	if !cnt {
		return nil, errors.New("Unable to proceed without accepting API usage terms.")
	}

	res, err := api.PostAcceptAPIUsageTerms(APIUsageAcceptAuthToken)
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
