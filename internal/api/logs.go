package api

import (
	"encoding/json"
	"fmt"

	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func QueryEnvironmentLogs(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string, payload map[string]interface{}) (*types.CustomerLogsResponse, error) {
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		Path:             fmt.Sprintf("/subscription/%s/environment/%s/log-query", subAliasOrHashedID, envNameOrHashedID),
		MapStringPayload: payload,
	}

	resp, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIQueryLogsErrorMsg)
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(resp.Body))
	}

	custLog := &types.CustomerLogsResponse{}
	err = json.Unmarshal(resp.Body, custLog)
	if err != nil {
		return nil, err
	}

	return custLog, nil
}

func GetEnvironmentLogLabelValues(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, label string) ([]string, error) {
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             fmt.Sprintf("/subscription/%s/environment/%s/log-label-values?label=%s", subAliasOrHashedID, envNameOrHashedID, label),
		MapStringPayload: nil,
	}

	resp, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIQueryLogsErrorMsg)
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(resp.Body))
	}

	var lvs types.LogLabelValuesResponse
	err = json.Unmarshal(resp.Body, &lvs)
	if err != nil {
		return nil, err
	}

	return lvs.Result, nil
}
