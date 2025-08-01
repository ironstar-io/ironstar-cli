package api

import (
	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetSubscriptionEnvironments(creds types.Keylink, output, hashOrAlias string) ([]types.Environment, error) {
	empty := []types.Environment{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + hashOrAlias + "/environments",
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var envs []types.Environment
	err = json.Unmarshal(res.Body, &envs)
	if err != nil {
		return empty, err
	}

	return envs, nil
}

func GetSubscriptionEnvironment(creds types.Keylink, output, subHashOrAlias, envHashOrAlias string) (types.Environment, error) {
	empty := types.Environment{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias,
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var env types.Environment
	err = json.Unmarshal(res.Body, &env)
	if err != nil {
		return empty, err
	}

	return env, nil
}

func PatchEnvironment(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, restorePermission string) error {
	req := &Request{
		Retries:         3,
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
		return res.HandleFailure(output)
	}

	return nil
}

func PostEnvironmentHook(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, hookName string) error {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/hook",
		MapStringPayload: map[string]interface{}{
			"hook": hookName,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIUpdateEnvironmentErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure(output)
	}

	return nil
}

func DeleteEnvironmentHook(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, hookName string) error {
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "DELETE",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/hook/" + hookName,
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIUpdateEnvironmentErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure(output)
	}

	return nil
}

func GetEnvironmentContext(creds types.Keylink, flg flags.Accumulator, subHashOrAlias string) (types.Environment, error) {
	empty := types.Environment{}

	envID := flg.Environment
	if envID == "" {
		env, err := services.StdinPrompt("Environment ID or Name: ")
		if err != nil {
			return empty, nil
		}

		envID = env
	}

	env, err := GetSubscriptionEnvironment(creds, flg.Output, subHashOrAlias, envID)
	if err != nil {
		return empty, err
	}

	return env, nil
}
