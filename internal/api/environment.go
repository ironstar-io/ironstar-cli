package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func GetSubscriptionEnvironments(creds types.Keylink, hashOrAlias string) ([]types.Environment, error) {
	empty := []types.Environment{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + hashOrAlias + "/environments",
		MapStringPayload: map[string]string{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var envs []types.Environment
	err = yaml.Unmarshal(res.Body, &envs)
	if err != nil {
		return empty, err
	}

	return envs, nil
}

func GetSubscriptionEnvironment(creds types.Keylink, subHashOrAlias, envHashOrAlias string) (types.Environment, error) {
	empty := types.Environment{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias,
		MapStringPayload: map[string]string{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var env types.Environment
	err = yaml.Unmarshal(res.Body, &env)
	if err != nil {
		return empty, err
	}

	return env, nil
}

// func GetSubscriptionContext(creds types.Keylink, flg flags.Accumulator) (types.Subscription, error) {
// 	empty := types.Subscription{}
// 	if flg.Subscription != "" {
// 		sub, err := GetSubscription(creds, flg.Subscription)
// 		if err != nil {
// 			return empty, err
// 		}

// 		return sub, nil
// 	}

// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return empty, err
// 	}

// 	confPath := filepath.Join(wd, ".ironstar", "config.yml")

// 	exists := fs.CheckExists(confPath)
// 	if !exists {
// 		createNewProj := services.ConfirmationPrompt("Couldn't find a project configuration in this directory. Would you like to create one?", "y", flg.AutoAccept)
// 		if createNewProj == true {
// 			err = services.InitializeIronstarProject()
// 			if err != nil {
// 				return empty, err
// 			}
// 		} else {
// 			return empty, errors.New("This command requires a project to be configured.")
// 		}

// 	}

// 	pr := fs.ProjectRoot()
// 	proj, err := services.ReadInProjectConfig(pr)
// 	if err != nil {
// 		return empty, errors.Wrap(err, errs.NoProjectFoundErrorMsg)
// 	}

// 	return proj.Subscription, nil
// }
