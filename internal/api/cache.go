package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"encoding/json"
	"github.com/pkg/errors"
)

func GetEnvironmentCacheInvalidations(creds types.Keylink, subHashOrAlias, envHashOrAlias string) ([]types.CacheInvalidation, error) {
	empty := []types.CacheInvalidation{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/cache-invalidations",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetCacheInvalidationErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var cis []types.CacheInvalidation
	err = json.Unmarshal(res.Body, &cis)
	if err != nil {
		return empty, err
	}

	return cis, nil
}

func GetEnvironmentCacheInvalidation(creds types.Keylink, subHashOrAlias, envHashOrAlias, invalidationName string) (types.CacheInvalidation, error) {
	empty := types.CacheInvalidation{}
	req := &Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/cache-invalidations/" + invalidationName,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIGetCacheInvalidationErrorMsg)
	}

	if res.StatusCode != 200 {
		return empty, res.HandleFailure()
	}

	var ci types.CacheInvalidation
	err = json.Unmarshal(res.Body, &ci)
	if err != nil {
		return empty, err
	}

	return ci, nil
}

func PostEnvironmentCacheInvalidation(creds types.Keylink, subHashOrAlias, envHashOrAlias string) (types.CacheInvalidation, error) {
	empty := types.CacheInvalidation{}
	req := &Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/subscription/" + subHashOrAlias + "/environment/" + envHashOrAlias + "/cache-invalidation",
		MapStringPayload: map[string]interface{}{
			"objects": []string{"*"},
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return empty, errors.Wrap(err, errs.APIPostCacheInvalidationErrorMsg)
	}

	if res.StatusCode != 201 {
		return empty, res.HandleFailure()
	}

	var ci types.CacheInvalidation
	err = json.Unmarshal(res.Body, &ci)
	if err != nil {
		return empty, err
	}

	return ci, nil
}
