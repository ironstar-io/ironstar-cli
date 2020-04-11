package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func GetSubscription(authToken, hashOrAlias string) (types.Subscription, error) {
	empty := types.Subscription{}
	req := &Request{
		AuthToken:        authToken,
		Method:           "GET",
		Path:             "/subscription/" + hashOrAlias,
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var sub types.Subscription
	err = yaml.Unmarshal(res.Body, &sub)
	if err != nil {
		return empty, err
	}

	return sub, nil
}
