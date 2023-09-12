package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func PostSyncRequest(creds types.Keylink, output string, payload types.PostSyncRequestParams) (types.SyncRequest, error) {
	empty := types.SyncRequest{}
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + payload.SubscriptionID + "/sync-request",
		MapStringPayload: map[string]interface{}{
			"restore_strategy": payload.RestoreStrategy,
			"src_environment":  payload.SrcEnvironment,
			"dest_environment": payload.DestEnvironment,
			"components":       payload.Components,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIPostSyncErrorMsg)
	}

	if res.StatusCode != 201 {
		return empty, res.HandleFailure(output)
	}

	var sr types.SyncRequest
	err = json.Unmarshal(res.Body, &sr)
	if err != nil {
		return empty, err
	}

	return sr, nil
}

func PostSyncRequestUseLatestBackup(creds types.Keylink, output string, payload types.PostSyncRequestParams) (types.RestoreRequest, error) {
	empty := types.RestoreRequest{}
	req := &Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + payload.SubscriptionID + "/sync-request?use-latest-backup=true",
		MapStringPayload: map[string]interface{}{
			"restore_strategy": payload.RestoreStrategy,
			"src_environment":  payload.SrcEnvironment,
			"dest_environment": payload.DestEnvironment,
			"components":       payload.Components,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIPostSyncErrorMsg)
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

func GetSubscriptionSync(creds types.Keylink, output string, subAliasOrHashedID, syncName string) (types.SyncRequest, error) {
	empty := types.SyncRequest{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/sync-requests/" + syncName,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSyncErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var sr types.SyncRequest
	err = json.Unmarshal(res.Body, &sr)
	if err != nil {
		return empty, err
	}

	return sr, nil
}

func GetSubscriptionSyncRequests(creds types.Keylink, output string, subAliasOrHashedID string) ([]types.SyncRequest, error) {
	empty := []types.SyncRequest{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/sync-requests",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSyncErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var srs []types.SyncRequest
	err = json.Unmarshal(res.Body, &srs)
	if err != nil {
		return empty, err
	}

	return srs, nil
}
