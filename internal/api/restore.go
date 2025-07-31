package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func PostRestoreRequest(creds types.Keylink, output string, payload types.PostRestoreRequestParams) (types.RestoreRequest, error) {
	empty := types.RestoreRequest{}
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + payload.SubscriptionID + "/environment/" + payload.EnvironmentID + "/restore-request",
		MapStringPayload: map[string]interface{}{
			"name":       payload.Name,
			"strategy":   payload.Strategy,
			"iteration":  payload.Backup,
			"components": payload.Components,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIPostRestoreErrorMsg)
	}

	if res.StatusCode != 201 {
		return empty, res.HandleFailure(output)
	}

	var rr types.RestoreRequest
	err = json.Unmarshal(res.Body, &rr)
	if err != nil {
		return empty, err
	}

	return rr, nil
}

func GetSubscriptionRestoreIterations(creds types.Keylink, output, subAliasOrHashedID string) ([]types.RestoreRequest, error) {
	empty := []types.RestoreRequest{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/restores",
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetRestoreErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var ris []types.RestoreRequest
	err = json.Unmarshal(res.Body, &ris)
	if err != nil {
		return empty, err
	}

	return ris, nil
}

func GetEnvironmentRestoreIterations(creds types.Keylink, output, subAliasOrHashedID, envNameOrHashedID string) ([]types.RestoreRequest, error) {
	empty := []types.RestoreRequest{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/restore-requests",
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetRestoreErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var ris []types.RestoreRequest
	err = json.Unmarshal(res.Body, &ris)
	if err != nil {
		return empty, err
	}

	return ris, nil
}

func GetEnvironmentRestore(creds types.Keylink, output, subAliasOrHashedID, envNameOrHashedID, restoreName string) (types.RestoreRequest, error) {
	empty := types.RestoreRequest{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/restore-requests/" + restoreName,
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetRestoreErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var r types.RestoreRequest
	err = json.Unmarshal(res.Body, &r)
	if err != nil {
		return empty, err
	}

	return r, nil
}
