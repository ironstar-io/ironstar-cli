package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"
	"github.com/pkg/errors"
)

func GetSubscriptionEnvironments(creds types.Keylink, hashOrAlias string) ([]types.Environment, error) {
	empty := []types.Environment{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + hashOrAlias + "/environments",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var envs []types.Environment
	err = json.Unmarshal(res.Body, &envs)
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
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var env types.Environment
	err = json.Unmarshal(res.Body, &env)
	if err != nil {
		return empty, err
	}

	return env, nil
}

func PatchEnvironment(creds types.Keylink, subHashOrAlias, envHashOrAlias, restorePermission string) error {
	req := &Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "PATCH",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias,
		MapStringPayload: map[string]interface{}{
			"restore_permission": restorePermission,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIUpdateEnvironmentErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure()
	}

	return nil
}

func GetEnvironmentContext(creds types.Keylink, flg flags.Accumulator, subHashOrAlias string) (types.Environment, error) {
	empty := types.Environment{}
	if flg.Environment != "" {
		env, err := GetSubscriptionEnvironment(creds, subHashOrAlias, flg.Environment)
		if err != nil {
			return empty, err
		}

		return env, nil
	}

	return empty, errors.New(errs.NoEnvironmentFlagSupplied)
}
