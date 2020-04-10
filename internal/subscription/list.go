package subscription

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/pkg/errors"
)

var APISubListErrorMsg = "Ironstar API failed to retrieve subscriptions"

func List(args []string, loginFlag string) error {
	user, err := services.ResolveUserCredentials(loginFlag)
	if err != nil {
		return err
	}

	req := &api.Request{
		AuthToken:        user.AuthToken,
		Method:           "GET",
		Path:             "/user/subscriptions",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	err = services.OutputJSON(res.Body)
	if err != nil {
		return errors.Wrap(err, APISubListErrorMsg)
	}

	return nil
}
