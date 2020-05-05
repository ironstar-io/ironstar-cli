package auth

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func MFADisable(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	color.Yellow("Removing MFA for this account...")
	fmt.Println()

	passcode, err := services.GetCLIMFAPasscode()
	if err != nil {
		return err
	}

	req := &api.Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/auth/mfa/remove",
		MapStringPayload: map[string]string{
			"passcode": passcode,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure()
	}

	fmt.Println()
	fmt.Println()
	color.Green("Successfully disabled MFA for this account")

	return nil
}
