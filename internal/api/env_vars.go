package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"
	"github.com/pkg/errors"
)

func GetEnvironmentEnvVars(creds types.Keylink, subHashOrAlias, envHashOrAlias string) ([]types.EnvVars, error) {
	empty := []types.EnvVars{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/env-vars",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetEnvironmentErrorMsg)
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
