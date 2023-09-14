package api

import (
	"net/http"

	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetRemoteCommands(creds types.Keylink, output, subHashOrAlias, envHashOrAlias string) ([]types.RemoteCommandResponse, error) {
	empty := []types.RemoteCommandResponse{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/remote-cmds",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetRemoteCommandsErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var rcs []types.RemoteCommandResponse
	err = json.Unmarshal(res.Body, &rcs)
	if err != nil {
		return empty, err
	}

	return rcs, nil
}

func PostRemoteCommandDrushCacheRebuild(creds types.Keylink, output, subHashOrAlias, envHashOrAlias string, envVars []types.RemoteCommandEnvironmentVariable, timeout int) (*types.RemoteCommandResponse, error) {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/remote-cmds/drush-cache-rebuild",
		MapStringPayload: map[string]interface{}{
			"environment_variables": envVars,
			"timeout":               timeout,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIPostRemoteCommandsErrorMsg)
	}

	if res.StatusCode != 201 {
		return nil, res.HandleFailure(output)
	}

	var rcs *types.RemoteCommandResponse
	err = json.Unmarshal(res.Body, &rcs)
	if err != nil {
		return nil, err
	}

	return rcs, nil
}

func PostRemoteCommandDrush(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, args string, envVars []types.RemoteCommandEnvironmentVariable, timeout int) (*types.RemoteCommandResponse, error) {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/remote-cmds/drush",
		MapStringPayload: map[string]interface{}{
			"args":                  args,
			"environment_variables": envVars,
			"timeout":               timeout,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIPostRemoteCommandsErrorMsg)
	}

	if res.StatusCode != 201 {
		return nil, res.HandleFailure(output)
	}

	var rcs *types.RemoteCommandResponse
	err = json.Unmarshal(res.Body, &rcs)
	if err != nil {
		return nil, err
	}

	return rcs, nil
}

func PostRemoteCommandShell(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, workDir, command string, envVars []types.RemoteCommandEnvironmentVariable, timeout int) (*types.RemoteCommandResponse, error) {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/remote-cmds/shell",
		MapStringPayload: map[string]interface{}{
			"command":               command,
			"work_dir":              workDir,
			"environment_variables": envVars,
			"timeout":               timeout,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIPostRemoteCommandsErrorMsg)
	}

	if res.StatusCode != 201 {
		return nil, res.HandleFailure(output)
	}

	var rcs *types.RemoteCommandResponse
	err = json.Unmarshal(res.Body, &rcs)
	if err != nil {
		return nil, err
	}

	return rcs, nil
}

func GetRemoteCommand(creds types.Keylink, output, subHashOrAlias, envHashOrAlias, rcni string) (*types.RemoteCommandResponse, error) {
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          http.MethodGet,
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/remote-cmds/" + rcni,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIGetRemoteCommandsErrorMsg)
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure(output)
	}

	var rcs *types.RemoteCommandResponse
	err = json.Unmarshal(res.Body, &rcs)
	if err != nil {
		return nil, err
	}

	return rcs, nil
}
