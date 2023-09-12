package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetDeployment(creds types.Keylink, output, deployID string) (types.Deployment, error) {
	empty := types.Deployment{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/deployment/" + deployID,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var d types.Deployment
	err = json.Unmarshal(res.Body, &d)
	if err != nil {
		return empty, err
	}

	return d, nil
}

func GetDeploymentActivity(creds types.Keylink, output, deployID string) ([]types.DeploymentActivityResponse, error) {
	empty := []types.DeploymentActivityResponse{}
	req := &Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/deployment/" + deployID + "/activity",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure(output)
	}

	var d []types.DeploymentActivityResponse
	err = json.Unmarshal(res.Body, &d)
	if err != nil {
		return empty, err
	}

	return d, nil
}
