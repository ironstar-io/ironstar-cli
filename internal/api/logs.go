package api

import (
	"encoding/json"

	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func QueryEnvironmentLogs(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string, payload map[string]interface{}) ([]types.CWLogResponse, error) {
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/logs",
		MapStringPayload: payload,
	}

	resp, err := req.ArimaSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIQueryLogsErrorMsg)
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(resp.Body))
	}

	var cwLog []types.CWLogResponse
	err = json.Unmarshal(resp.Body, &cwLog)
	if err != nil {
		return nil, err
	}

	return cwLog, nil
}

func GetEnvironmentLogStreams(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string) ([]types.CWLogStreamsResponse, error) {
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/log-streams",
		MapStringPayload: map[string]interface{}{},
	}

	resp, err := req.ArimaSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIQueryLogsErrorMsg)
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(resp.Body))
	}

	var cwLog []types.CWLogStreamsResponse
	err = json.Unmarshal(resp.Body, &cwLog)
	if err != nil {
		return nil, err
	}

	return cwLog, nil
}
