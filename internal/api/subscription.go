package api

import (
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func GetSubscription(creds types.Keylink, hashOrAlias string) (types.Subscription, error) {
	empty := types.Subscription{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + hashOrAlias,
		MapStringPayload: map[string]string{},
	}

	res, err := req.NankaiSend()
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

func GetUserSubscriptions(creds types.Keylink) ([]types.UserAccessResponse, error) {
	empty := []types.UserAccessResponse{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/user/subscriptions",
		MapStringPayload: map[string]string{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var uar []types.UserAccessResponse
	err = yaml.Unmarshal(res.Body, &uar)
	if err != nil {
		return empty, err
	}

	return uar, nil
}

func GetSubscriptionContext(creds types.Keylink, flg flags.Accumulator) (types.Subscription, error) {
	empty := types.Subscription{}
	if flg.Subscription != "" {
		sub, err := GetSubscription(creds, flg.Subscription)
		if err != nil {
			return empty, err
		}

		return sub, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return empty, err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if !exists {
		createNewProj := services.ConfirmationPrompt("Couldn't find a project configuration in this directory. Would you like to create one?", "y", flg.AutoAccept)
		if createNewProj == true {
			err = services.InitializeIronstarProject()
			if err != nil {
				return empty, err
			}
		} else {
			return empty, errors.New("This command requires a project to be configured.")
		}

	}

	pr := fs.ProjectRoot()
	proj, err := services.ReadInProjectConfig(pr)
	if err != nil {
		return empty, errors.Wrap(err, errs.NoProjectFoundErrorMsg)
	}

	return proj.Subscription, nil
}

func GetSubscriptionEnvironmentContext(creds types.Keylink, flg flags.Accumulator) (types.SubscriptionEnvironment, error) {
	empty := types.SubscriptionEnvironment{}

	sub, err := GetSubscriptionContext(creds, flg)
	if err != nil {
		return empty, err
	}
	if sub == (types.Subscription{}) {
		return empty, errors.New("No subscription was able to found")
	}

	env, err := GetEnvironmentContext(creds, flg, sub.HashedID)
	if err != nil {
		return empty, err
	}
	if env == (types.Environment{}) {
		return empty, errors.New("No environment was able to found")
	}

	return types.SubscriptionEnvironment{
		Subscription: sub,
		Environment:  env,
	}, nil
}
