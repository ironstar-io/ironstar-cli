package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetEnvironmentEnvVars(creds types.Keylink, subHashOrAlias, envHashOrAlias string) ([]types.EnvVars, error) {
	empty := []types.EnvVars{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/env-vars",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentVariablesErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var env_vars []types.EnvVars
	err = json.Unmarshal(res.Body, &env_vars)
	if err != nil {
		return empty, err
	}

	return env_vars, nil
}

func PostEnvironmentEnvVars(creds types.Keylink, subHashOrAlias, envHashOrAlias, key, value, varType string) error {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/env-var",
		MapStringPayload: map[string]interface{}{
			"key":      key,
			"value":    value,
			"var_type": varType,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIPostEnvironmentVariableErrorMsg)
	}

	if res.StatusCode != 201 {
		return res.HandleFailure()
	}

	return nil
}

func PatchEnvironmentEnvVar(creds types.Keylink, subHashOrAlias, envHashOrAlias, key, value, varType string) error {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "PATCH",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/env-var/" + key,
		MapStringPayload: map[string]interface{}{
			"value":    value,
			"var_type": varType,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIPostEnvironmentVariableErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure()
	}

	return nil
}

func DeleteEnvironmentEnvVar(creds types.Keylink, subHashOrAlias, envHashOrAlias, key string) error {
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "DELETE",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/env-var/" + key,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIPostEnvironmentVariableErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure()
	}

	return nil
}
