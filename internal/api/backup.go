package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"
	"github.com/pkg/errors"
)

func PostBackupRequest(creds types.Keylink, payload types.PostBackupRequestParams) (types.BackupRequest, error) {
	empty := types.BackupRequest{}
	req := &Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + payload.SubscriptionID + "/environment/" + payload.EnvironmentID + "/backup-request",
		MapStringPayload: map[string]interface{}{
			"name":       payload.Name,
			"kind":       payload.Kind,
			"components": payload.Components,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 201 {
		return empty, res.HandleFailure()
	}

	var br types.BackupRequest
	err = json.Unmarshal(res.Body, &br)
	if err != nil {
		return empty, err
	}

	return br, nil
}

func GetSubscriptionBackupIterations(creds types.Keylink, subAliasOrHashedID string) ([]types.BackupIteration, error) {
	empty := []types.BackupIteration{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/backups",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var bis []types.BackupIteration
	err = json.Unmarshal(res.Body, &bis)
	if err != nil {
		return empty, err
	}

	return bis, nil
}

func GetEnvironmentBackupIterations(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string) ([]types.BackupIteration, error) {
	empty := []types.BackupIteration{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/backup-iterations",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var bis []types.BackupIteration
	err = json.Unmarshal(res.Body, &bis)
	if err != nil {
		return empty, err
	}

	return bis, nil
}

func GetEnvironmentBackup(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, backupName string) (types.Backup, error) {
	empty := types.Backup{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/backups/" + backupName,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetSubscriptionErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var b types.Backup
	err = json.Unmarshal(res.Body, &b)
	if err != nil {
		return empty, err
	}

	return b, nil
}
