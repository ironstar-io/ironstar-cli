package api

import (
	"encoding/json"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

func GetUserOrganisationControls(creds types.Keylink) (*types.OrganisationControls, error) {
	req := &Request{
		RunTokenRefresh:  false,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/user/organisation-controls",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APIGetUserOrgControlsErrorMsg)
	}

	if res.StatusCode == 204 {
		// User doesn't belong to an organisation
		return nil, nil
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	uoc := &types.OrganisationControls{}
	err = json.Unmarshal(res.Body, uoc)
	if err != nil {
		return nil, err
	}

	return uoc, nil
}
