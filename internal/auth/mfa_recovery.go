package auth

import (
	"encoding/json"
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func MFARecovery(args []string, flg flags.Accumulator) error {
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
		return errors.Wrap(err, errs.APILoginErrorMsg)
	}

	c := &types.AuthResponseBody{}
	err = json.Unmarshal(res.Body, c)
	if err != nil {
		return err
	}

	creds := types.Keylink{
		Login:     email,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	}

	rcode, err := services.StdinSecret("Recovery Code: ")
	if err != nil {
		return err
	}

	req := &api.Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/auth/mfa/recovery",
		MapStringPayload: map[string]interface{}{
			"recovery_code": rcode,
		},
	}

	rres, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if rres.StatusCode != 200 {
		return rres.HandleFailure()
	}

	rrc := &types.AuthResponseBody{}
	err = json.Unmarshal(rres.Body, rrc)
	if err != nil {
		return err
	}

	enableCreds := types.Keylink{
		Login:     email,
		AuthToken: rrc.IDToken,
		Expiry:    rrc.Expiry,
	}

	fmt.Println()
	fmt.Println()
	color.Yellow(`NOTICE!`)
	color.Yellow(`MFA has now been disabled for this account, however usage of the Ironstar API now requires all users to have MFA enabled`)
	fmt.Println()

	cnt := services.ConfirmationPrompt("Would like to re-enable MFA now?", "y", false)
	if !cnt {
		return errors.New("Recovery successful, MFA not re-enabled.")
	}

	_, err = MFAEnable(flg, enableCreds, email)
	if err != nil {
		return err
	}

	return nil
}
