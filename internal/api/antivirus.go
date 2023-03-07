package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetEnvironmentAntivirusScans(creds types.Keylink, subHashOrAlias, envHashOrAlias string) ([]types.AntivirusScan, error) {
	empty := []types.AntivirusScan{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/antivirus-scans",
		MapStringPayload: map[string]interface{}{},
		Retries:          3,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetAntivirusScanErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var avs []types.AntivirusScan
	err = json.Unmarshal(res.Body, &avs)
	if err != nil {
		return empty, err
	}

	return avs, nil
}
