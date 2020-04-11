package subscription

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func List(args []string, loginFlag string) error {
	creds, err := services.ResolveUserCredentials(loginFlag)
	if err != nil {
		return err
	}

	color.Green("Using login [" + creds.Login + "]")

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/user/subscriptions",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	err = services.OutputJSON(res.Body)
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	return nil
}
