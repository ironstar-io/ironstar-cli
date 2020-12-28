package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
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
			"name":        payload.Name,
			"kind":        payload.Kind,
			"components":  payload.Components,
			"lock_tables": payload.LockTables,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIPostBackupErrorMsg)
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

func DeleteBackup(creds types.Keylink, payload types.DeleteBackupParams) error {
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "DELETE",
		Path:             "/subscription/" + payload.SubscriptionID + "/backups/" + payload.Name,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APIDeleteBackupErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure()
	}

	return nil
}

func GetSubscriptionBackupIterations(creds types.Keylink, subAliasOrHashedID, backupType string) ([]types.BackupIteration, error) {
	var qs string
	if backupType != "" {
		qs = "?backup-type=" + backupType
	}

	empty := []types.BackupIteration{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/backups" + qs,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetBackupErrorMsg)
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

func DownloadEnvironmentBackupComponent(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, backupName, savePath string, buComp types.BackupIterationComponent) error {
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/backups/" + backupName + "/download?component=" + buComp.Name,
		MapStringPayload: map[string]interface{}{},
	}

	resp, err := req.ArimaDownload(savePath, buComp.Name)
	if err != nil {
		return errors.Wrap(err, errs.APIGetBackupErrorMsg)
	}

	if resp.StatusCode != 200 {
		return errors.New(string(resp.Body))
	}

	return nil
}

func GetEnvironmentBackupIterations(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, backupType string) ([]types.BackupIteration, error) {
	var qs string
	if backupType != "" {
		qs = "?backup-type=" + backupType
	}

	empty := []types.BackupIteration{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/backups" + qs,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetBackupErrorMsg)
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

func GetLatestEnvironmentBackupIteration(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID string) (types.BackupIteration, error) {
	empty := types.BackupIteration{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subAliasOrHashedID + "/environment/" + envNameOrHashedID + "/backup-iterations?latest=true",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetBackupErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var bi types.BackupIteration
	err = json.Unmarshal(res.Body, &bi)
	if err != nil {
		return empty, err
	}

	return bi, nil
}

func GetEnvironmentBackup(creds types.Keylink, subAliasOrHashedID, envNameOrHashedID, backupName, errorOutput string) (types.Backup, error) {
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
		return empty, errors.Wrap(err, errs.APIGetBackupErrorMsg)
	}

	if res.StatusCode != 200 {
		if errorOutput == constants.SKIP_ERRORS {
			return empty, nil
		}
		return empty, res.HandleFailure()
	}

	var b types.Backup
	err = json.Unmarshal(res.Body, &b)
	if err != nil {
		return empty, err
	}

	return b, nil
}
