package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
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
	err = yaml.Unmarshal(res.Body, &br)
	if err != nil {
		return empty, err
	}

	return br, nil
}
