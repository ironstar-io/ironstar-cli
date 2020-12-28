package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func PutNewRelicApplicationConfig(creds types.Keylink, subId, envId string, payload types.PutNewRelicParams) error {
	req := &Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "PUT",
		Path:            "/subscription/" + subId + "/environment/" + envId + "/application-config",
		MapStringPayload: map[string]interface{}{
			"new_relic": payload,
		},
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
