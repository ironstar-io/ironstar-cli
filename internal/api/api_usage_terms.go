package api

import (
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"

	"github.com/pkg/errors"
)

func GetCurrentAPIUsageTerms() (types.APIUsageTerms, error) {
	empty := types.APIUsageTerms{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      types.Keylink{},
		Method:           "GET",
		Path:             "/api-usage-terms",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetAntivirusScanErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var aut types.APIUsageTerms
	err = json.Unmarshal(res.Body, &aut)
	if err != nil {
		return empty, err
	}

	return aut, nil
}

func PostAcceptAPIUsageTerms(authToken string) (*RawResponse, error) {
	req := &Request{
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: authToken,
		},
		Method:           "POST",
		Path:             "/user/accept-api-usage-terms",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return res, nil
}
