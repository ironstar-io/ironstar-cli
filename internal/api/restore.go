package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func PostRestoreRequest(creds types.Keylink, payload types.PostRestoreRequestParams) (types.RestoreRequest, error) {
	empty := types.RestoreRequest{}
	req := &Request{
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
		return empty, res.HandleFailure()
	}

	var rr types.RestoreRequest
	err = json.Unmarshal(res.Body, &rr)
	if err != nil {
		return empty, err
	}

	return rr, nil
}

func GetSubscriptionRestoreIterations(creds types.Keylink, subAliasOrHashedID string) ([]types.RestoreRequest, error) {
	empty := []types.RestoreRequest{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/restores",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var ris []types.RestoreRequest
	err = json.Unmarshal(res.Body, &ris)
	if err != nil {
		return empty, err
	}

	return ris, nil
}

func GetEnvironmentRestoreIterations(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string) ([]types.RestoreRequest, error) {
	empty := []types.RestoreRequest{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/restore-requests",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var ris []types.RestoreRequest
	err = json.Unmarshal(res.Body, &ris)
	if err != nil {
		return empty, err
	}

	return ris, nil
}

func GetEnvironmentRestore(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, restoreName string) (types.RestoreRequest, error) {
	empty := types.RestoreRequest{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/restore-requests/" + restoreName,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var r types.RestoreRequest
	err = json.Unmarshal(res.Body, &r)
	if err != nil {
		return empty, err
	}

	return r, nil
}
